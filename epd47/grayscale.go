package epd47

// DrawMode like C
type DrawMode uint8

const (
	BlackOnWhite DrawMode = iota
	WhiteOnWhite
	WhiteOnBlack
)

// Contrast cycles from C driver
var contrast4 = [Frames4bpp]int{30, 30, 20, 20, 30, 30, 30, 40, 40, 50, 50, 50, 100, 200, 300}
var contrast4White = [Frames4bpp]int{10, 10, 8, 8, 8, 8, 8, 10, 10, 15, 15, 20, 20, 100, 300}

func (d *Device) resetLUT(mode DrawMode) {
	fill := byte(0x55)
	if mode == WhiteOnBlack || mode == WhiteOnWhite {
		fill = 0xAA
	}
	for i := 0; i < len(d.convLUT); i++ {
		d.convLUT[i] = fill
	}
}

func (d *Device) updateLUT(k uint8, mode DrawMode) {
	kk := k
	if mode == BlackOnWhite || mode == WhiteOnWhite {
		kk = 15 - k
	}
	// Simplified LUT update for smaller table
	lutSize := uint32(len(d.convLUT))
	for l := uint32(kk); l < lutSize; l += 16 {
		if l < lutSize {
			d.convLUT[l] &= 0xFC
		}
	}
	for l := uint32(kk) << 4; l < lutSize; l += (1 << 8) {
		for p := uint32(0); p < 16 && l+p < lutSize; p++ {
			d.convLUT[l+p] &= 0xF3
		}
	}
}

// calcEPDInput4bpp: fill line4b from v1..v4 blocks.
// For simplicity we assume lane order b1,b2,b3,b4; swap if needed after hardware test.
func (d *Device) calcEPDInput4bpp(v []uint16, outLen int) {
	oi := 0
	lutMask := uint16(len(d.convLUT) - 1)
	for j := 0; oi < outLen && j+3 < len(v); j += 4 {
		b1 := d.convLUT[v[j+0]&lutMask]
		b2 := d.convLUT[v[j+1]&lutMask]
		b3 := d.convLUT[v[j+2]&lutMask]
		b4 := d.convLUT[v[j+3]&lutMask]
		d.line4b[oi+0] = b1
		d.line4b[oi+1] = b2
		d.line4b[oi+2] = b3
		d.line4b[oi+3] = b4
		oi += 4
	}
}

// expand4bppLine creates v1..v4 uint16 sequence across full width from a subspan.
// No allocations: uses a local static workspace and copies minimal slices.
func (d *Device) expand4bppLine(src []byte, x, w int, v []uint16) {
	// full-width packed 4bpp scratch: reuse line4b as temporary byte buffer, since we overwrite it later from calc
	full := d.line4b[:d.w/2]
	// clear only the part needed at edges
	for i := 0; i < len(full); i++ {
		full[i] = 0
	}
	// place src at x/2
	start := x / 2
	copy(full[start:start+len(src)], src)

	// Build v from pairs of bytes
	vi := 0
	for i := 0; i < len(full); i += 2 {
		if vi >= len(v) {
			break
		}
		v[vi] = uint16(full[i]) | (uint16(full[i+1]) << 8)
		vi++
	}
}

// Public 4bpp draw: 15-frame pipeline.
func (d *Device) DrawImage4bpp(x, y, w, h int, data []byte, mode DrawMode) {
	if w <= 0 || h <= 0 {
		return
	}
	if x < 0 || y < 0 || x+w > d.w || y+h > d.h {
		// add full clipping later
		return
	}

	lut := contrast4[:]
	if mode == WhiteOnBlack {
		lut = contrast4White[:]
	}
	srcStride := (w / 2) + (w % 2)

	// v buffer: 4 uint16 per 2 bytes across full width/2 -> (Width/2)/2 = Width/4 entries
	var v [MaxWidth / 4]uint16
	outLen := d.w / 2

	d.resetLUT(mode)
	for k := 0; k < Frames4bpp; k++ {
		d.updateLUT(uint8(k), mode)
		d.StartFrame()
		for row := 0; row < d.h; row++ {
			if row < y || row >= y+h {
				d.SkipRow()
				continue
			}
			sr := data[(row-y)*srcStride : (row-y+1)*srcStride]
			d.expand4bppLine(sr, x, w, v[:d.w/4])
			d.calcEPDInput4bpp(v[:d.w/4], outLen)
			d.latchRow()
			d.pulseCKV(lut[k], 50)
			d.writeLineBytes(d.line4b[:outLen])
			d.pulseCKV(1, 1)
		}
		d.EndFrame()
		// small settle
		d.bus.sleepUS(5_000)
	}
}

// 1bpp helpers

// Clear flashes the full screen dark/white for cycles.
func (d *Device) Clear(cycles int) {
	if cycles <= 0 {
		cycles = 2
	}
	n := d.w / 8
	for c := 0; c < cycles; c++ {
		// dark
		for i := 0; i < n; i++ {
			d.line1b[i] = 0x00
		}
		d.StartFrame()
		for row := 0; row < d.h; row++ {
			d.outputRow1bpp(n, 10)
		}
		d.EndFrame()

		// white
		for i := 0; i < n; i++ {
			d.line1b[i] = 0xFF
		}
		d.StartFrame()
		for row := 0; row < d.h; row++ {
			d.outputRow1bpp(n, 10)
		}
		d.EndFrame()
	}
}

// Draw1bpp draws a packed MSB-first 1bpp image at x,y.
func (d *Device) Draw1bpp(x, y, w, h int, src []byte, pulseUS int) {
	if w <= 0 || h <= 0 {
		return
	}
	if x < 0 || y < 0 || x+w > d.w || y+h > d.h {
		return
	}
	srcStride := (w + 7) / 8
	dstStride := d.w / 8

	d.StartFrame()
	for row := 0; row < d.h; row++ {
		if row < y || row >= y+h {
			d.SkipRow()
			continue
		}
		// zero line
		for i := 0; i < dstStride; i++ {
			d.line1b[i] = 0
		}
		// blit row bits into position
		sr := src[(row-y)*srcStride : (row-y+1)*srcStride]
		for col := 0; col < w; col++ {
			sbyte := sr[col>>3]
			sbit := 7 - (col & 7)
			on := (sbyte>>sbit)&1 == 1
			if on {
				dbitpos := x + col
				dbyte := dbitpos >> 3
				dbit := 7 - (dbitpos & 7)
				d.line1b[dbyte] |= (1 << dbit)
			}
		}
		d.outputRow1bpp(dstStride, pulseUS)
	}
	d.EndFrame()
}
