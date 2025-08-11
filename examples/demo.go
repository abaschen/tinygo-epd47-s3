// Demo example showing EPD47 driver usage patterns
// This example can be built for various targets for testing
package main

import (
	"machine"
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

// Pin adapter function
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

// Sleep function
func sleepUS(us int) {
	if us <= 0 {
		return
	}
	time.Sleep(time.Duration(us) * time.Microsecond)
}

// Create a checkerboard pattern in 1bpp format
func createCheckerboard(w, h, blockSize int) []byte {
	data := make([]byte, (w+7)/8*h)
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			// Create checkerboard pattern
			blockR := r / blockSize
			blockC := c / blockSize
			on := (blockR&1)^(blockC&1) == 1
			
			if on {
				byteIdx := r*((w+7)/8) + (c >> 3)
				bit := 7 - (c & 7)
				data[byteIdx] |= (1 << bit)
			}
		}
	}
	return data
}

// Create a horizontal gradient in 4bpp format
func createGradient(w, h int) []byte {
	data := make([]byte, (w/2+w%2)*h)
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			// Create horizontal gradient (0-15)
			shade := byte((c * 15) / (w - 1))
			
			byteIdx := r*((w/2)+(w%2)) + (c >> 1)
			if (c & 1) == 0 {
				// Even column - upper nibble
				data[byteIdx] = (shade << 4) | (data[byteIdx] & 0x0F)
			} else {
				// Odd column - lower nibble
				data[byteIdx] = (data[byteIdx] & 0xF0) | (shade & 0x0F)
			}
		}
	}
	return data
}

func main() {
	// Configure for ESP32-S3 LilyGo T5 4.7"
	cfg := epd47.Config{
		Width:  960,
		Height: 540,

		// Configuration pins
		CFG_DATA: pinOut(machine.Pin(13)),
		CFG_CLK:  pinOut(machine.Pin(12)),
		CFG_STR:  pinOut(machine.Pin(0)),

		// Control pins
		CKV: pinOut(machine.Pin(38)),
		STH: pinOut(machine.Pin(40)),
		CKH: pinOut(machine.Pin(41)),

		// Data bus D0-D7
		D0: pinOut(machine.Pin(8)),
		D1: pinOut(machine.Pin(1)),
		D2: pinOut(machine.Pin(2)),
		D3: pinOut(machine.Pin(3)),
		D4: pinOut(machine.Pin(4)),
		D5: pinOut(machine.Pin(5)),
		D6: pinOut(machine.Pin(6)),
		D7: pinOut(machine.Pin(7)),

		SleepUS: sleepUS,
	}

	// Initialize display
	display := epd47.New(cfg)
	display.Configure()
	
	// Power on sequence
	display.PowerOn()
	time.Sleep(200 * time.Millisecond)

	// Clear the display
	display.Clear(2)
	time.Sleep(500 * time.Millisecond)

	// Demo 1: Draw checkerboard pattern (1bpp)
	cbWidth, cbHeight := 200, 120
	cbX := (display.Width() - cbWidth) / 2
	cbY := 50
	
	checkerboard := createCheckerboard(cbWidth, cbHeight, 16)
	display.Draw1bpp(cbX, cbY, cbWidth, cbHeight, checkerboard, 10)
	
	time.Sleep(1 * time.Second)

	// Demo 2: Draw gradient (4bpp)
	gradWidth, gradHeight := 300, 80
	gradX := (display.Width() - gradWidth) / 2
	gradY := cbY + cbHeight + 40
	
	gradient := createGradient(gradWidth, gradHeight)
	display.DrawImage4bpp(gradX, gradY, gradWidth, gradHeight, gradient, epd47.BlackOnWhite)
	
	time.Sleep(2 * time.Second)

	// Demo 3: Draw some text-like patterns
	textWidth, textHeight := 400, 60
	textX := (display.Width() - textWidth) / 2
	textY := gradY + gradHeight + 40
	
	// Create a simple text-like pattern
	textPattern := make([]byte, (textWidth+7)/8*textHeight)
	for r := 0; r < textHeight; r++ {
		for c := 0; c < textWidth; c++ {
			// Create horizontal lines with gaps (text-like)
			if r%10 < 6 && c%20 < 15 {
				byteIdx := r*((textWidth+7)/8) + (c >> 3)
				bit := 7 - (c & 7)
				textPattern[byteIdx] |= (1 << bit)
			}
		}
	}
	
	display.Draw1bpp(textX, textY, textWidth, textHeight, textPattern, 10)
	
	time.Sleep(3 * time.Second)

	// Power off
	display.PowerOffAll()

	// Sleep forever
	for {
		time.Sleep(time.Hour)
	}
}