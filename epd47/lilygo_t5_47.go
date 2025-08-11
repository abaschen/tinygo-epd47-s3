// +build tinygo

package epd47

import (
	"machine"
	"time"
)

// LilyGoT547 represents a preconfigured EPD47 device for the LilyGo T5 4.7" ESP32-S3 board.
// This provides a convenient way to initialize the display with the correct pin mappings.
type LilyGoT547 struct {
	*Device
}

// NewLilyGoT547 creates a new EPD47 device preconfigured for the LilyGo T5 4.7" ESP32-S3 board.
// All pins are automatically configured according to the board's hardware layout.
func NewLilyGoT547() *LilyGoT547 {
	// Pin adapter function
	pinOut := func(p machine.Pin) PinOut {
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
	sleepUS := func(us int) {
		if us <= 0 {
			return
		}
		time.Sleep(time.Duration(us) * time.Microsecond)
	}

	// LilyGo T5 4.7" ESP32-S3 pin configuration
	cfg := Config{
		Width:  960,
		Height: 540,

		// Configuration shift register pins
		CFG_DATA: pinOut(machine.Pin(13)),
		CFG_CLK:  pinOut(machine.Pin(12)),
		CFG_STR:  pinOut(machine.Pin(0)),

		// Control lines
		CKV: pinOut(machine.Pin(38)), // Vertical gate clock
		STH: pinOut(machine.Pin(40)), // Start/enable (input enable)
		CKH: pinOut(machine.Pin(41)), // Write strobe

		// 8-bit data bus D0..D7 (LSB..MSB)
		// Pin mapping: D7..D0 = 7,6,5,4,3,2,1,8
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

	device := New(cfg)
	device.Configure()
	return &LilyGoT547{Device: device}
}

// Initialize performs the complete initialization sequence for the display.
// This includes power-on and initial clear operations.
func (d *LilyGoT547) Initialize() {
	// Power on the display
	d.PowerOn()
	time.Sleep(200 * time.Millisecond)

	// Clear the display to ensure clean state
	d.Clear(2)
	time.Sleep(100 * time.Millisecond)
}

// Shutdown performs a complete shutdown of the display.
// This should be called before the device goes to sleep or powers off.
func (d *LilyGoT547) Shutdown() {
	d.PowerOffAll()
}

// DrawText draws a simple text-like pattern for demonstration.
// This is a basic implementation - for real text rendering, use a font library.
func (d *LilyGoT547) DrawText(x, y int, text string, charWidth, charHeight int) {
	if len(text) == 0 {
		return
	}

	totalWidth := len(text) * (charWidth + 2) // 2px spacing between chars
	totalHeight := charHeight
	
	// Create a simple pattern for each character
	data := make([]byte, (totalWidth+7)/8*totalHeight)
	// Note: make() already zeros the slice, so no need to clear()
	
	for i, char := range text {
		charX := i * (charWidth + 2)
		
		// Create a simple pattern based on the character
		// This is just a demo - real text would use font data
		pattern := int(char) % 8
		
		for r := 0; r < charHeight; r++ {
			for c := 0; c < charWidth; c++ {
				// Simple pattern generation
				if (r+c+pattern)%3 == 0 {
					bitX := charX + c
					byteIdx := r*((totalWidth+7)/8) + (bitX >> 3)
					bit := 7 - (bitX & 7)
					if byteIdx < len(data) {
						data[byteIdx] |= (1 << bit)
					}
				}
			}
		}
	}
	
	d.Draw1bpp(x, y, totalWidth, totalHeight, data, 10)
}

// DrawRectangle draws a rectangle outline in 1bpp mode.
func (d *LilyGoT547) DrawRectangle(x, y, w, h int, filled bool) {
	if w <= 0 || h <= 0 {
		return
	}
	
	data := make([]byte, (w+7)/8*h)
	// Note: make() already zeros the slice, so no need to clear()
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			draw := false
			
			if filled {
				draw = true
			} else {
				// Only draw border
				if r == 0 || r == h-1 || c == 0 || c == w-1 {
					draw = true
				}
			}
			
			if draw {
				byteIdx := r*((w+7)/8) + (c >> 3)
				bit := 7 - (c & 7)
				data[byteIdx] |= (1 << bit)
			}
		}
	}
	
	d.Draw1bpp(x, y, w, h, data, 10)
}

// DrawGradient draws a gradient pattern in 4bpp mode.
func (d *LilyGoT547) DrawGradient(x, y, w, h int, horizontal bool, mode DrawMode) {
	if w <= 0 || h <= 0 {
		return
	}
	
	data := make([]byte, (w/2+w%2)*h)
	// Note: make() already zeros the slice, so no need to clear()
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			var shade byte
			
			if horizontal {
				// Horizontal gradient
				shade = byte((c * 15) / (w - 1))
			} else {
				// Vertical gradient
				shade = byte((r * 15) / (h - 1))
			}
			
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
	
	d.DrawImage4bpp(x, y, w, h, data, mode)
}

// DrawCheckerboard draws a checkerboard pattern in 1bpp mode.
func (d *LilyGoT547) DrawCheckerboard(x, y, w, h, blockSize int) {
	if w <= 0 || h <= 0 || blockSize <= 0 {
		return
	}
	
	data := make([]byte, (w+7)/8*h)
	// Note: make() already zeros the slice, so no need to clear()
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
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
	
	d.Draw1bpp(x, y, w, h, data, 10)
}

// GetDisplayInfo returns information about the display.
func (d *LilyGoT547) GetDisplayInfo() (width, height int, model string) {
	return d.Width(), d.Height(), "LilyGo T5 4.7\" ESP32-S3 (ED047TC1)"
}