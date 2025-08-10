package main

import (
	"time"

	"github.com/abaschen/epd47"
)

// Adapt machine.Pin to epd47.PinOut
func pinOut(p machine.Pin) epd47.PinOut {
	p.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return func(level bool) {
		if level {
			p.High()
		} else {
			p.Low()
		}
	}
}

func sleepUS(us int) {
	if us <= 0 {
		return
	}
	time.Sleep(time.Duration(us) * time.Microsecond)
}

func main() {
	// ESP32-S3 LilyGo T5 4.7" pins from ed047tc1.h
	cfg := epd47.Config{
		Width:  960,
		Height: 540,

		CFG_DATA: pinOut(machine.GPIO13),
		CFG_CLK:  pinOut(machine.GPIO12),
		CFG_STR:  pinOut(machine.GPIO0),

		CKV: pinOut(machine.GPIO38),
		STH: pinOut(machine.GPIO40),
		CKH: pinOut(machine.GPIO41),

		// D7..D0 = 7,6,5,4,3,2,1,8
		D0: pinOut(machine.GPIO8),
		D1: pinOut(machine.GPIO1),
		D2: pinOut(machine.GPIO2),
		D3: pinOut(machine.GPIO3),
		D4: pinOut(machine.GPIO4),
		D5: pinOut(machine.GPIO5),
		D6: pinOut(machine.GPIO6),
		D7: pinOut(machine.GPIO7),

		SleepUS: sleepUS,
	}
	d := epd47.New(cfg)

	d.PowerOn()
	time.Sleep(200 * time.Millisecond)

	// Clear and show 1bpp checkerboard
	d.Clear(2)
	w := 200
	h := 120
	x := (d.Width() - w) / 2
	y := (d.Height() - h) / 2
	src1 := make([]byte, (w+7)/8*h)
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			on := ((r>>4)&1)^((c>>4)&1) == 1
			if on {
				byteIdx := r*((w+7)/8) + (c >> 3)
				bit := 7 - (c & 7)
				src1[byteIdx] |= (1 << bit)
			}
		}
	}
	d.Draw1bpp(x, y, w, h, src1, 10)

	// 4bpp gradient
	gw, gh := 240, 120
	gx := (d.Width() - gw) / 2
	gy := y + h + 20
	src4 := make([]byte, (gw/2+gw%2)*gh)
	for r := 0; r < gh; r++ {
		for c := 0; c < gw; c++ {
			shade := byte((c * 15) / (gw - 1)) // 0..15
			byteIdx := r*((gw/2)+(gw%2)) + (c >> 1)
			if (c & 1) == 0 {
				src4[byteIdx] = (shade << 4) | (src4[byteIdx] & 0x0F)
			} else {
				src4[byteIdx] = (src4[byteIdx] & 0xF0) | (shade & 0x0F)
			}
		}
	}
	d.DrawImage4bpp(gx, gy, gw, gh, src4, epd47.BlackOnWhite)

	time.Sleep(3 * time.Second)
	d.PowerOffAll()

	for {
		time.Sleep(time.Hour)
	}
}
