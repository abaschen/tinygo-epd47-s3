package main

import (
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

// Create a more complex pattern
func createComplexPattern(w, h int) []byte {
	data := make([]byte, (w+7)/8*h)
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			// Create a complex mathematical pattern
			x := float32(c - w/2)
			y := float32(r - h/2)
			
			// Distance from center
			dist := x*x + y*y
			
			// Create ripple effect
			ripple := int(dist/100) % 4
			
			// Create radial lines
			angle := int((x*16/y + 100)) % 8
			
			if ripple == 0 || angle < 2 {
				byteIdx := r*((w+7)/8) + (c >> 3)
				bit := 7 - (c & 7)
				data[byteIdx] |= (1 << bit)
			}
		}
	}
	return data
}

// Create a mandala-like pattern in 4bpp
func createMandala(w, h int) []byte {
	data := make([]byte, (w/2+w%2)*h)
	
	centerX := w / 2
	centerY := h / 2
	
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			x := c - centerX
			y := r - centerY
			
			// Distance from center
			dist := x*x + y*y
			
			// Create concentric circles with varying intensity
			intensity := (int(dist/50) % 16)
			
			// Add some angular variation
			if x != 0 || y != 0 {
				angle := int((float32(y)/float32(x)*8 + 16)) % 16
				intensity = (intensity + angle) % 16
			}
			
			shade := byte(intensity)
			
			byteIdx := r*((w/2)+(w%2)) + (c >> 1)
			if (c & 1) == 0 {
				data[byteIdx] = (shade << 4) | (data[byteIdx] & 0x0F)
			} else {
				data[byteIdx] = (data[byteIdx] & 0xF0) | (shade & 0x0F)
			}
		}
	}
	return data
}

func main() {
	// Create the preconfigured device
	display := epd47.NewLilyGoT547()
	
	// Initialize
	display.Initialize()
	
	width, height, model := display.GetDisplayInfo()
	println("Starting advanced demo on", model)
	println("Resolution:", width, "x", height)

	// Scene 1: Geometric patterns
	println("Scene 1: Geometric patterns")
	
	// Multiple checkerboards with different sizes
	display.DrawCheckerboard(50, 50, 150, 100, 8)
	display.DrawCheckerboard(250, 50, 150, 100, 16)
	display.DrawCheckerboard(450, 50, 150, 100, 24)
	
	// Nested rectangles
	for i := 0; i < 8; i++ {
		size := 20 + i*15
		x := 700 + (120-size)/2
		y := 50 + (100-size)/2
		display.DrawRectangle(x, y, size, size, i%2 == 0)
	}
	
	time.Sleep(2 * time.Second)

	// Scene 2: Complex mathematical patterns
	println("Scene 2: Complex patterns")
	
	// Clear previous content
	display.Clear(1)
	
	// Draw complex pattern
	patternW, patternH := 300, 200
	patternX := (width - patternW) / 2
	patternY := 50
	
	complexData := createComplexPattern(patternW, patternH)
	display.Draw1bpp(patternX, patternY, patternW, patternH, complexData, 10)
	
	// Add some text labels
	display.DrawText(patternX, patternY-30, "Complex Pattern", 8, 16)
	
	time.Sleep(3 * time.Second)

	// Scene 3: Grayscale mandala
	println("Scene 3: Grayscale mandala")
	
	display.Clear(1)
	
	// Create and draw mandala
	mandalaW, mandalaH := 400, 300
	mandalaX := (width - mandalaW) / 2
	mandalaY := (height - mandalaH) / 2
	
	mandalaData := createMandala(mandalaW, mandalaH)
	display.DrawImage4bpp(mandalaX, mandalaY, mandalaW, mandalaH, mandalaData, epd47.BlackOnWhite)
	
	// Add title
	display.DrawText(mandalaX, mandalaY-40, "Grayscale Mandala", 10, 20)
	
	time.Sleep(4 * time.Second)

	// Scene 4: Mixed content layout
	println("Scene 4: Mixed content layout")
	
	display.Clear(1)
	
	// Title
	display.DrawText(50, 30, "EPD47 Driver Demo", 16, 24)
	
	// Left column - 1bpp content
	display.DrawText(50, 80, "1bpp Graphics:", 8, 16)
	display.DrawCheckerboard(50, 110, 100, 80, 12)
	display.DrawRectangle(50, 210, 100, 60, false)
	display.DrawRectangle(60, 220, 80, 40, true)
	
	// Center column - gradients
	display.DrawText(200, 80, "4bpp Gradients:", 8, 16)
	display.DrawGradient(200, 110, 150, 40, true, epd47.BlackOnWhite)
	display.DrawGradient(200, 160, 40, 100, false, epd47.BlackOnWhite)
	display.DrawGradient(250, 160, 100, 100, true, epd47.WhiteOnBlack)
	
	// Right column - patterns
	display.DrawText(400, 80, "Patterns:", 8, 16)
	for i := 0; i < 6; i++ {
		y := 110 + i*25
		display.DrawRectangle(400, y, 150, 20, i%2 == 0)
	}
	
	// Bottom section - large gradient
	display.DrawText(50, 300, "Full Width Gradient:", 8, 16)
	display.DrawGradient(50, 330, width-100, 60, true, epd47.BlackOnWhite)
	
	// Footer
	display.DrawText(50, 450, "LilyGo T5 4.7\" ESP32-S3 - TinyGo Driver", 6, 12)
	
	time.Sleep(5 * time.Second)

	// Scene 5: Animation-like sequence
	println("Scene 5: Animation sequence")
	
	display.Clear(1)
	
	// Animated expanding rectangles
	centerX := width / 2
	centerY := height / 2
	
	for frame := 0; frame < 10; frame++ {
		if frame > 0 {
			display.Clear(1)
		}
		
		// Draw expanding rectangles
		for i := 0; i <= frame; i++ {
			size := 20 + i*30
			x := centerX - size/2
			y := centerY - size/2
			
			if x >= 0 && y >= 0 && x+size <= width && y+size <= height {
				display.DrawRectangle(x, y, size, size, i%2 == 0)
			}
		}
		
		// Add frame counter
		display.DrawText(20, 20, "Frame:", 8, 16)
		// Simple number display (just show pattern based on frame)
		display.DrawRectangle(100, 20, frame*10, 16, true)
		
		time.Sleep(500 * time.Millisecond)
	}
	
	time.Sleep(2 * time.Second)

	// Final scene
	println("Demo complete!")
	display.DrawText(centerX-100, centerY-10, "Demo Complete!", 12, 20)
	
	time.Sleep(3 * time.Second)
	
	// Shutdown
	display.Shutdown()
	
	println("Display powered off. Sleeping...")
	for {
		time.Sleep(time.Hour)
	}
}