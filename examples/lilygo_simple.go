package main

import (
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

func main() {
	// Create a preconfigured LilyGo T5 4.7" device
	// All pins are automatically configured
	display := epd47.NewLilyGoT547()

	// Initialize the display (power on + clear)
	display.Initialize()

	// Get display information
	width, height, model := display.GetDisplayInfo()
	println("Display:", model)
	println("Resolution:", width, "x", height)

	// Demo 1: Draw a checkerboard pattern
	display.DrawCheckerboard(50, 50, 200, 120, 16)
	time.Sleep(1 * time.Second)

	// Demo 2: Draw some rectangles
	display.DrawRectangle(300, 50, 150, 80, false)  // Outline
	display.DrawRectangle(320, 70, 110, 40, true)   // Filled
	time.Sleep(1 * time.Second)

	// Demo 3: Draw gradients
	display.DrawGradient(50, 200, 200, 80, true, epd47.BlackOnWhite)   // Horizontal
	display.DrawGradient(300, 200, 80, 120, false, epd47.BlackOnWhite) // Vertical
	time.Sleep(1 * time.Second)

	// Demo 4: Draw some "text" (simple pattern)
	display.DrawText(50, 350, "Hello EPD47!", 12, 20)
	time.Sleep(2 * time.Second)

	// Demo 5: Full screen effects
	centerX := width / 2
	centerY := height / 2

	// Large centered rectangle
	display.DrawRectangle(centerX-100, centerY-60, 200, 120, false)
	time.Sleep(1 * time.Second)

	// Nested rectangles
	for i := 0; i < 5; i++ {
		size := 40 + i*20
		x := centerX - size/2
		y := centerY - size/2
		display.DrawRectangle(x, y, size, size, false)
		time.Sleep(500 * time.Millisecond)
	}

	// Final display time
	time.Sleep(5 * time.Second)

	// Shutdown the display
	display.Shutdown()

	// Sleep forever
	for {
		time.Sleep(time.Hour)
	}
}