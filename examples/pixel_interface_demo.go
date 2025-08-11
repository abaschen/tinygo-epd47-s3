// Demonstration of the improved pixel-level interface
package main

import (
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

func main() {
	// Create display
	display := epd47.NewLilyGoT547()
	display.Initialize()
	
	println("Pixel Interface Demo - Improved Implementation")
	
	// Get display dimensions
	w, h := display.Size()
	println("Display size:", int(w), "x", int(h))
	
	// Demo 1: Individual pixel drawing
	println("Demo 1: Drawing individual pixels...")
	
	// Draw a diagonal line using SetPixel
	for i := int16(0); i < 100; i++ {
		display.SetPixel(100+i, 100+i, true)
	}
	
	// Draw a cross pattern
	centerX, centerY := int16(200), int16(150)
	for i := int16(-20); i <= 20; i++ {
		display.SetPixel(centerX+i, centerY, true)   // Horizontal line
		display.SetPixel(centerX, centerY+i, true)   // Vertical line
	}
	
	// Render the 1bpp pixels
	display.Display()
	println("1bpp pixels rendered")
	
	time.Sleep(2 * time.Second)
	
	// Demo 2: Grayscale pixel drawing
	println("Demo 2: Drawing grayscale pixels...")
	
	// Draw a grayscale gradient using individual pixels
	for x := int16(300); x < 400; x++ {
		for y := int16(100); y < 200; y++ {
			// Create a gradient based on position
			shade := uint8(((x - 300) * 15) / 99)
			display.SetGrayscalePixel(x, y, shade)
		}
	}
	
	// Draw some grayscale circles
	for radius := uint8(1); radius <= 15; radius++ {
		// Simple circle approximation
		for angle := 0; angle < 360; angle += 10 {
			x := int16(500 + int(float32(radius*3) * 0.017453292 * float32(angle))) // cos approximation
			y := int16(150 + int(float32(radius*2) * 0.017453292 * float32(angle))) // sin approximation
			if x >= 0 && x < w && y >= 0 && y < h {
				display.SetGrayscalePixel(x, y, radius)
			}
		}
	}
	
	// Render the 4bpp pixels
	display.Display()
	println("4bpp pixels rendered")
	
	time.Sleep(2 * time.Second)
	
	// Demo 3: Mixed drawing - combine pixel interface with high-level methods
	println("Demo 3: Mixed drawing methods...")
	
	// Use high-level method for background
	display.DrawRectangle(50, 250, 200, 100, false)
	
	// Use pixel interface to add details inside
	for x := int16(60); x < 240; x += 5 {
		for y := int16(260); y < 340; y += 5 {
			display.SetPixel(x, y, true)
		}
	}
	
	// Add some grayscale accents
	for i := int16(0); i < 50; i++ {
		shade := uint8((i * 15) / 49)
		display.SetGrayscalePixel(270+i, 270, shade)
		display.SetGrayscalePixel(270+i, 280, shade)
		display.SetGrayscalePixel(270+i, 290, shade)
	}
	
	// Render everything
	display.Display()
	println("Mixed drawing completed")
	
	time.Sleep(2 * time.Second)
	
	// Demo 4: Pixel manipulation
	println("Demo 4: Pixel reading and manipulation...")
	
	// Set some test pixels
	display.SetPixel(400, 300, true)
	display.SetGrayscalePixel(410, 300, 8)
	
	// Read them back
	pixel1 := display.GetPixel(400, 300)
	pixel2 := display.GetGrayscalePixel(410, 300)
	
	println("Read pixel values:", pixel1, int(pixel2))
	
	// Modify based on current values
	if pixel1 {
		// Create a small pattern around the pixel
		for dx := int16(-2); dx <= 2; dx++ {
			for dy := int16(-2); dy <= 2; dy++ {
				display.SetPixel(400+dx, 300+dy, (dx+dy)%2 == 0)
			}
		}
	}
	
	if pixel2 > 0 {
		// Create a grayscale pattern
		for dx := int16(-3); dx <= 3; dx++ {
			for dy := int16(-3); dy <= 3; dy++ {
				distance := dx*dx + dy*dy
				if distance <= 9 { // Within circle
					shade := uint8(15 - distance)
					display.SetGrayscalePixel(410+dx, 300+dy, shade)
				}
			}
		}
	}
	
	display.Display()
	println("Pixel manipulation completed")
	
	time.Sleep(3 * time.Second)
	
	// Demo 5: Performance comparison
	println("Demo 5: Performance demonstration...")
	
	// Time pixel-by-pixel drawing
	start := time.Now()
	for i := int16(0); i < 1000; i++ {
		display.SetPixel(100+(i%100), 400+(i/100), i%2 == 0)
	}
	pixelTime := time.Since(start)
	
	// Time high-level drawing
	start = time.Now()
	display.DrawCheckerboard(200, 400, 100, 10, 5)
	highLevelTime := time.Since(start)
	
	println("Pixel-by-pixel time:", pixelTime.String())
	println("High-level method time:", highLevelTime.String())
	
	// Render both
	display.Display()
	
	time.Sleep(2 * time.Second)
	
	// Final message
	display.DrawText(50, 450, "Pixel Interface Demo Complete!", 8, 16)
	display.Display()
	
	println("\nDemo Summary:")
	println("✅ Individual pixel control (SetPixel/GetPixel)")
	println("✅ Grayscale pixel control (SetGrayscalePixel/GetGrayscalePixel)")
	println("✅ Sparse buffer implementation (memory efficient)")
	println("✅ Automatic rendering with Display()")
	println("✅ Mixed drawing methods")
	println("✅ Full display area support")
	
	time.Sleep(3 * time.Second)
	display.Shutdown()
	
	for {
		time.Sleep(time.Hour)
	}
}