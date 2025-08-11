// +build !tinygo

package epd47

import (
	"time"
)

// Mock pin type for testing
type mockPin int

func (p mockPin) Configure(config interface{}) {}
func (p mockPin) High()                        {}
func (p mockPin) Low()                         {}

// LilyGoT547 represents a preconfigured EPD47 device for the LilyGo T5 4.7" ESP32-S3 board.
// This provides a convenient way to initialize the display with the correct pin mappings.
type LilyGoT547 struct {
	*Device
}

// NewLilyGoT547 creates a new EPD47 device preconfigured for the LilyGo T5 4.7" ESP32-S3 board.
// All pins are automatically configured according to the board's hardware layout.
// This is a mock version for testing without TinyGo.
func NewLilyGoT547() *LilyGoT547 {
	// Mock pin adapter function
	pinOut := func(p mockPin) PinOut {
		return func(level bool) {
			// Mock implementation - does nothing
		}
	}

	// Mock sleep function
	sleepUS := func(us int) {
		if us <= 0 {
			return
		}
		// Use nanoseconds for faster testing
		time.Sleep(time.Duration(us) * time.Nanosecond)
	}

	// LilyGo T5 4.7" ESP32-S3 pin configuration (mock)
	cfg := Config{
		Width:  960,
		Height: 540,

		// Configuration shift register pins
		CFG_DATA: pinOut(mockPin(13)),
		CFG_CLK:  pinOut(mockPin(12)),
		CFG_STR:  pinOut(mockPin(0)),

		// Control lines
		CKV: pinOut(mockPin(38)), // Vertical gate clock
		STH: pinOut(mockPin(40)), // Start/enable (input enable)
		CKH: pinOut(mockPin(41)), // Write strobe

		// 8-bit data bus D0..D7 (LSB..MSB)
		D0: pinOut(mockPin(8)),
		D1: pinOut(mockPin(1)),
		D2: pinOut(mockPin(2)),
		D3: pinOut(mockPin(3)),
		D4: pinOut(mockPin(4)),
		D5: pinOut(mockPin(5)),
		D6: pinOut(mockPin(6)),
		D7: pinOut(mockPin(7)),

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

// GetDisplayInfo returns information about the display.
func (d *LilyGoT547) GetDisplayInfo() (width, height int, model string) {
	return d.Width(), d.Height(), "LilyGo T5 4.7\" ESP32-S3 (ED047TC1) - Mock"
}