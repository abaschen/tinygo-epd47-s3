package epd47

import "time"

// Clear flashes whole screen by driving a band across all rows: dark then white, a few cycles.
func (e *EPD) Clear(cycles int) {
	if cycles <= 0 {
		cycles = 2
	}
	fullLineBytes := EPDWidth / 8

	// Two presets: DARK and WHITE. We approximate with bytes patterns similar to CLEAR_BYTE/DARK_BYTE idea.
	dark := make([]byte, fullLineBytes)
	white := make([]byte, fullLineBytes)

	// dark: set bits 0 to produce a darkening pulse; white: set bits 1 for lightening.
	for i := 0; i < fullLineBytes; i++ {
		dark[i] = 0x00
		white[i] = 0xFF
	}

	for c := 0; c < cycles; c++ {
		// Dark cycle
		e.StartFrame()
		for row := 0; row < EPDHeight; row++ {
			e.OutputRow1bpp(dark, 10)
		}
		e.EndFrame()
		time.Sleep(20 * time.Millisecond)

		// White cycle
		e.StartFrame()
		for row := 0; row < EPDHeight; row++ {
			e.OutputRow1bpp(white, 10)
		}
		e.EndFrame()
		time.Sleep(20 * time.Millisecond)
	}
}

// PushBand approximates epd_push_pixels(): draw a horizontal area y..y+h with a solid color (1=white).
func (e *EPD) PushBand(y, h int, color bool, pulseUS int) {
	if y < 0 {
		y = 0
	}
	if h <= 0 {
		return
	}
	if y+h > EPDHeight {
		h = EPDHeight - y
	}
	fullLineBytes := EPDWidth / 8
	rowData := make([]byte, fullLineBytes)
	fill := byte(0x00)
	if color { // white
		fill = 0xFF
	}
	for i := range rowData {
		rowData[i] = fill
	}

	e.StartFrame()
	for row := 0; row < EPDHeight; row++ {
		if row < y {
			e.SkipRow()
		} else if row >= y && row < y+h {
			e.OutputRow1bpp(rowData, pulseUS)
		} else {
			e.SkipRow()
		}
	}
	// pipeline tail
	e.OutputRow1bpp(rowData, pulseUS)
	e.EndFrame()
}

// Draw1bpp draws a 1bpp bitmap area into the screen. src is tightly packed, MSB first.
// area: x,y,w,h. If x,y out of bounds, they are clipped.
func (e *EPD) Draw1bpp(x, y, w, h int, src []byte, pulseUS int) {
	if w <= 0 || h <= 0 {
		return
	}
	if x < 0 {
		// naive clipping: shift src horizontally is complex; for baseline require x>=0, else just clip left columns away
		clip := (-x)
		if clip >= w {
			return
		}
		x = 0
		w -= clip
	}
	if y < 0 {
		clipRows := -y
		if clipRows >= h {
			return
		}
		y = 0
		h -= clipRows
		// also advance src by clipRows*w bits
		bytesPerSrcLine := (w + 7) / 8
		src = src[clipRows*bytesPerSrcLine:]
	}
	if x+w > EPDWidth {
		w = EPDWidth - x
	}
	if y+h > EPDHeight {
		h = EPDHeight - y
	}
	if w <= 0 || h <= 0 {
		return
	}

	bytesPerSrcLine := (w + 7) / 8
	fullLineBytes := EPDWidth / 8

	e.StartFrame()
	for row := 0; row < EPDHeight; row++ {
		if row < y || row >= y+h {
			e.SkipRow()
			continue
		}
		// compose one full line buffer EPDWidth/8, filling left x pixels with off, right side off.
		line := make([]byte, fullLineBytes)
		// place src row bits at position x..x+w
		srcRow := src[(row-y)*bytesPerSrcLine : (row-y+1)*bytesPerSrcLine]

		// naive placement: iterate bits and set in destination
		for col := 0; col < w; col++ {
			byteIdx := col >> 3
			bit := 7 - (col & 7)
			on := (srcRow[byteIdx]>>bit)&1 == 1
			if on {
				destBitPos := x + col
				dByte := destBitPos >> 3
				dBit := 7 - (destBitPos & 7)
				line[dByte] |= (1 << dBit)
			}
		}
		e.OutputRow1bpp(line, pulseUS)
	}
	// tail
	tail := make([]byte, fullLineBytes)
	e.OutputRow1bpp(tail, pulseUS)
	e.EndFrame()
}

func (e *EPD) startFramePipeline() { e.StartFrame() }
func (e *EPD) endFramePipeline()   { e.EndFrame() }
