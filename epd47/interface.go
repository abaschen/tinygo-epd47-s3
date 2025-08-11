package epd47

// Displayer defines the interface for e-paper display operations.
// This follows TinyGo driver patterns for display interfaces.
type Displayer interface {
	// Configure initializes the display hardware
	Configure() error
	
	// Size returns the display dimensions
	Size() (x, y int16)
	
	// Display updates the screen with the current buffer contents
	Display() error
	
	// ClearDisplay clears the entire display
	ClearDisplay() error
	
	// SetPixel sets a single pixel (for 1bpp displays)
	SetPixel(x, y int16, c bool)
	
	// GetPixel gets a single pixel state
	GetPixel(x, y int16) bool
}

// GrayscaleDisplayer extends Displayer for grayscale operations
type GrayscaleDisplayer interface {
	Displayer
	
	// SetGrayscalePixel sets a pixel with grayscale value (0-15)
	SetGrayscalePixel(x, y int16, c uint8)
	
	// GetGrayscalePixel gets a grayscale pixel value
	GetGrayscalePixel(x, y int16) uint8
}

// Ensure Device implements the interfaces
var _ Displayer = (*Device)(nil)
var _ GrayscaleDisplayer = (*Device)(nil)