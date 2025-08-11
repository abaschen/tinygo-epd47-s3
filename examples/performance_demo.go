// Performance demonstration showing the benefits of using clear() builtin
package main

import (
	"time"

	"github.com/abaschen/tinygo-epd47-s3/epd47"
)

// Benchmark function to demonstrate clear() performance
func benchmarkBufferClearing() {
	println("Performance Demo: Buffer Clearing Optimizations")
	
	// Create test buffers
	const bufSize = 4096 // Same size as our LUT
	testBuf1 := make([]byte, bufSize)
	testBuf2 := make([]byte, bufSize)
	
	// Method 1: Manual loop (old way)
	start := time.Now()
	for i := 0; i < 1000; i++ {
		for j := range testBuf1 {
			testBuf1[j] = 0
		}
	}
	manualTime := time.Since(start)
	
	// Method 2: Using clear() builtin (new way)
	start = time.Now()
	for i := 0; i < 1000; i++ {
		clear(testBuf2)
	}
	clearTime := time.Since(start)
	
	println("Manual loop time:", manualTime.String())
	println("clear() builtin time:", clearTime.String())
	
	if clearTime < manualTime {
		improvement := float64(manualTime-clearTime) / float64(manualTime) * 100
		println("Performance improvement: ~", int(improvement), "%")
	}
}

func main() {
	// Run performance demo
	benchmarkBufferClearing()
	
	// Create display and demonstrate optimized operations
	display := epd47.NewLilyGoT547()
	display.Initialize()
	
	println("\nDemonstrating optimized display operations...")
	
	// The following operations now use clear() internally for better performance:
	
	// 1. Display clearing (uses clear() for zero buffers)
	println("1. Clearing display (optimized with clear())")
	display.Clear(1)
	
	// 2. Drawing operations that clear buffers internally
	println("2. Drawing checkerboard (internal buffer clearing optimized)")
	display.DrawCheckerboard(100, 100, 200, 150, 16)
	
	// 3. Grayscale operations (LUT clearing optimized)
	println("3. Drawing gradient (LUT operations optimized)")
	display.DrawGradient(350, 100, 200, 150, true, epd47.BlackOnWhite)
	
	// 4. Text drawing (buffer initialization optimized)
	println("4. Drawing text (buffer management optimized)")
	display.DrawText(100, 300, "Optimized with clear()!", 8, 16)
	
	time.Sleep(3 * time.Second)
	
	println("\nAll operations completed using optimized buffer management!")
	println("Key improvements:")
	println("- clear() used for zero-filling buffers")
	println("- Efficient LUT management")
	println("- Optimized line buffer clearing")
	println("- Better memory performance")
	
	display.Shutdown()
	
	for {
		time.Sleep(time.Hour)
	}
}