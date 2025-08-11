// Example demonstrating TinyGo driver patterns and interfaces
package main

import (
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

// Example using the Displayer interface
func demonstrateDisplayerInterface(display epd47.Displayer) {
	// Get display size
	w, h := display.Size()
	println("Display size:", int(w), "x", int(h))
	
	// Clear the display
	display.ClearDisplay()
	
	// Set some pixels (simplified implementation - only first line)
	for x := int16(0); x < w && x < 100; x += 10 {
		display.SetPixel(x, 0, true)
	}
	
	// Update display
	display.Display()
}

// Example using the GrayscaleDisplayer interface
func demonstrateGrayscaleInterface(display epd47.GrayscaleDisplayer) {
	// Set grayscale pixels (simplified implementation - only first line)
	w, _ := display.Size()
	for x := int16(0); x < w && x < 100; x++ {
		shade := uint8((x * 15) / 99) // 0-15 gradient
		display.SetGrayscalePixel(x, 0, shade)
	}
	
	// Update display
	display.Display()
}

func main() {
	// Method 1: Use the convenience wrapper (recommended)
	println("Using LilyGoT547 wrapper...")
	lilygo := epd47.NewLilyGoT547()
	lilygo.Initialize()
	
	// High-level drawing operations
	lilygo.DrawCheckerboard(50, 50, 200, 100, 16)
	lilygo.DrawGradient(300, 50, 200, 100, true, epd47.BlackOnWhite)
	lilygo.DrawRectangle(50, 200, 150, 80, false)
	lilygo.DrawText(250, 220, "TinyGo EPD47", 8, 16)
	
	time.Sleep(3 * time.Second)
	
	// Method 2: Use standard TinyGo interfaces
	println("Using standard interfaces...")
	
	// Cast to interfaces
	var displayer epd47.Displayer = lilygo.Device
	var grayscaleDisplayer epd47.GrayscaleDisplayer = lilygo.Device
	
	// Use interface methods
	demonstrateDisplayerInterface(displayer)
	time.Sleep(1 * time.Second)
	
	demonstrateGrayscaleInterface(grayscaleDisplayer)
	time.Sleep(2 * time.Second)
	
	// Method 3: Manual configuration (for custom setups)
	println("Manual configuration example...")
	// This would be used for custom pin configurations
	// See other examples for manual setup
	
	// Shutdown
	lilygo.Shutdown()
	
	println("Demo complete!")
	
	// Sleep forever
	for {
		time.Sleep(time.Hour)
	}
}