package epd47

// PinOut sets an output pin to the specified logic level.
// Implementations must be fast and non-allocating.
type PinOut func(level bool)

// PinGet reads a pin (reserved for future use, not needed currently).
type PinGet func() bool

// SleepUS blocks for approximately the given number of microseconds.
// Provide a platform-specific implementation via Config.
type SleepUS func(us int)

// Config provides all parameters and pin bindings needed to drive the ED047TC1 panel
// on the LilyGo T5 4.7" with ESP32-S3 in 8-bit parallel mode.
// This package uses functional PinOut to avoid importing machine in the driver.
type Config struct {
	Width  int
	Height int

	// Config shift-register pins (CFG_*).
	CFG_DATA PinOut
	CFG_CLK  PinOut
	CFG_STR  PinOut

	// Control lines.
	CKV PinOut // Vertical gate clock
	STH PinOut // Start/enable (input enable)
	CKH PinOut // Write strobe

	// Data bus pins D0..D7 (LSB..MSB).
	D0 PinOut
	D1 PinOut
	D2 PinOut
	D3 PinOut
	D4 PinOut
	D5 PinOut
	D6 PinOut
	D7 PinOut

	// Sleep function in microseconds.
	SleepUS SleepUS
}

// Device represents the ED047TC1 e-paper panel interface.
// It maintains a shadow of the configuration register and preallocated line buffers
// to avoid heap churn in hot paths.
type Device struct {
	bus      *parallelBus
	w, h     int
	cfg      reg
	dataMask uint8

	// Preallocated line buffers to avoid per-row allocations.
	// 1bpp: width/8 bytes per line.
	line1b [MaxWidthBytes1bpp]byte
	// 4bpp panel-lane formatted output: width/2 bytes per line.
	line4b [MaxWidthBytes4bpp]byte

	// LUT for 4bpp conversion - moved from global to instance
	// Use smaller LUT for constrained targets
	convLUT [1 << 12]byte // 4KB instead of 64KB

	// Sparse pixel buffers for interface compliance
	// Only stores non-zero/non-false pixels to minimize memory usage
	pixelBuffer     map[uint32]bool  // 1bpp pixels (key: y<<16|x)
	grayscaleBuffer map[uint32]uint8 // 4bpp pixels (key: y<<16|x)
}

// Hardware/format limits for this panel.
const (
	MaxWidth          = 960
	MaxHeight         = 540
	MaxWidthBytes1bpp = MaxWidth / 8
	MaxWidthBytes4bpp = MaxWidth / 2
	Frames4bpp        = 15
)

// parallelBus handles the low-level parallel communication with the display
type parallelBus struct {
	// Config shift register pins
	cfgData, cfgClk, cfgStr PinOut
	// Control pins
	ckv, sth, ckh PinOut
	// Data bus D0..D7
	dataPins [8]PinOut
	// Sleep function
	sleepUS SleepUS
}

// DefaultConfig returns a baseline configuration for the T5 4.7" panel.
// You still need to bind PinOut functions and SleepUS before calling New.
func DefaultConfig() Config {
	return Config{
		Width:  MaxWidth,
		Height: MaxHeight,
	}
}

// New creates a new Device, validates and stores pins, preallocates buffers,
// and pushes the initial (powered-off) configuration register.
// It does not perform any panel power transitions; call Configure() next.
func New(cfg Config) *Device {
	// Width/height sanity
	w := cfg.Width
	h := cfg.Height
	if w <= 0 || w > MaxWidth {
		w = MaxWidth
	}
	if h <= 0 || h > MaxHeight {
		h = MaxHeight
	}

	// Sleep function default (no-op) if not provided.
	sl := cfg.SleepUS
	if sl == nil {
		sl = func(us int) {}
	}

	// Collect data pins in D0..D7 order and compute a safety mask.
	dataPins := [8]PinOut{cfg.D0, cfg.D1, cfg.D2, cfg.D3, cfg.D4, cfg.D5, cfg.D6, cfg.D7}
	var mask uint8
	for i := 0; i < 8; i++ {
		if dataPins[i] != nil {
			mask |= 1 << uint(i)
		}
	}

	// Create parallel bus
	bus := &parallelBus{
		cfgData:  cfg.CFG_DATA,
		cfgClk:   cfg.CFG_CLK,
		cfgStr:   cfg.CFG_STR,
		ckv:      cfg.CKV,
		sth:      cfg.STH,
		ckh:      cfg.CKH,
		dataPins: dataPins,
		sleepUS:  sl,
	}

	d := &Device{
		bus:      bus,
		w:        w,
		h:        h,
		dataMask: mask,
		cfg: reg{
			epLatchEnable:   false,
			powerDisable:    true,
			posPowerEnable:  false,
			negPowerEnable:  false,
			epSTV:           true,
			epScanDirection: true,
			epMode:          false,
			epOutputEnable:  false,
		},
	}

	return d
}

// Configure initializes the device and pushes the initial configuration.
// This should be called after New() and before any drawing operations.
func (d *Device) Configure() error {
	// Clear line buffers for clean state
	clear(d.line1b[:])
	clear(d.line4b[:])
	clear(d.convLUT[:])
	
	// Initialize pixel buffers (lazy initialization - they'll be created when needed)
	d.pixelBuffer = nil
	d.grayscaleBuffer = nil
	
	// Push initial (all-safed) config.
	d.pushCfg()
	return nil
}

// Width returns the panel width in pixels.
func (d *Device) Width() int { return d.w }

// Height returns the panel height in pixels.
func (d *Device) Height() int { return d.h }

// clearLineBuffer efficiently clears the 1bpp line buffer
func (d *Device) clearLineBuffer() {
	clear(d.line1b[:])
}

// clearGrayscaleBuffer efficiently clears the 4bpp line buffer
func (d *Device) clearGrayscaleBuffer() {
	clear(d.line4b[:])
}

// fillBuffer efficiently fills a byte slice with a specific value.
// Uses Go's built-in clear() for zero values (more efficient), 
// manual loop for non-zero values.
// Reference: https://zetcode.com/golang/builtins-clear/
func fillBuffer(buf []byte, value byte) {
	if value == 0 {
		// Use clear() for zero values - optimized by the compiler
		clear(buf)
	} else {
		// Manual fill for non-zero values
		for i := range buf {
			buf[i] = value
		}
	}
}

// renderPixelBuffer converts sparse pixel buffer to display format and renders it
func (d *Device) renderPixelBuffer() error {
	if len(d.pixelBuffer) == 0 {
		return nil
	}
	
	// Find bounding box to minimize rendering area
	minX, minY := int16(d.w), int16(d.h)
	maxX, maxY := int16(0), int16(0)
	
	for key := range d.pixelBuffer {
		x := int16(key & 0xFFFF)
		y := int16(key >> 16)
		if x < minX { minX = x }
		if x > maxX { maxX = x }
		if y < minY { minY = y }
		if y > maxY { maxY = y }
	}
	
	// Calculate render area
	w := int(maxX - minX + 1)
	h := int(maxY - minY + 1)
	
	// Create bitmap for the bounding box
	stride := (w + 7) / 8
	bitmap := make([]byte, stride*h)
	
	// Fill bitmap from sparse buffer
	for key, value := range d.pixelBuffer {
		if !value { continue } // Skip false pixels
		
		x := int16(key & 0xFFFF) - minX
		y := int16(key >> 16) - minY
		
		byteIdx := int(y)*stride + int(x>>3)
		bit := 7 - (int(x) & 7)
		bitmap[byteIdx] |= (1 << bit)
	}
	
	// Render to display
	d.Draw1bpp(int(minX), int(minY), w, h, bitmap, 10)
	
	// Clear the buffer after rendering
	clear(d.pixelBuffer)
	
	return nil
}

// renderGrayscaleBuffer converts sparse grayscale buffer to display format and renders it
func (d *Device) renderGrayscaleBuffer() error {
	if len(d.grayscaleBuffer) == 0 {
		return nil
	}
	
	// Find bounding box
	minX, minY := int16(d.w), int16(d.h)
	maxX, maxY := int16(0), int16(0)
	
	for key := range d.grayscaleBuffer {
		x := int16(key & 0xFFFF)
		y := int16(key >> 16)
		if x < minX { minX = x }
		if x > maxX { maxX = x }
		if y < minY { minY = y }
		if y > maxY { maxY = y }
	}
	
	// Calculate render area
	w := int(maxX - minX + 1)
	h := int(maxY - minY + 1)
	
	// Create 4bpp bitmap for the bounding box
	stride := (w + 1) / 2
	bitmap := make([]byte, stride*h)
	
	// Fill bitmap from sparse buffer
	for key, value := range d.grayscaleBuffer {
		if value == 0 { continue } // Skip zero pixels
		
		x := int16(key & 0xFFFF) - minX
		y := int16(key >> 16) - minY
		
		byteIdx := int(y)*stride + int(x>>1)
		if int(x)&1 == 0 {
			// Even column - upper nibble
			bitmap[byteIdx] = (value << 4) | (bitmap[byteIdx] & 0x0F)
		} else {
			// Odd column - lower nibble
			bitmap[byteIdx] = (bitmap[byteIdx] & 0xF0) | (value & 0x0F)
		}
	}
	
	// Render to display
	d.DrawImage4bpp(int(minX), int(minY), w, h, bitmap, BlackOnWhite)
	
	// Clear the buffer after rendering
	clear(d.grayscaleBuffer)
	
	return nil
}
// Size returns the display dimensions as required by Displayer interface
func (d *Device) Size() (x, y int16) {
	return int16(d.w), int16(d.h)
}

// Display updates the screen with accumulated pixels from SetPixel/SetGrayscalePixel calls
func (d *Device) Display() error {
	// Render 1bpp pixels if any exist
	if d.pixelBuffer != nil && len(d.pixelBuffer) > 0 {
		err := d.renderPixelBuffer()
		if err != nil {
			return err
		}
	}
	
	// Render 4bpp pixels if any exist
	if d.grayscaleBuffer != nil && len(d.grayscaleBuffer) > 0 {
		err := d.renderGrayscaleBuffer()
		if err != nil {
			return err
		}
	}
	
	return nil
}

// ClearDisplay clears the entire display and pixel buffers
func (d *Device) ClearDisplay() error {
	// Clear the physical display
	d.Clear(2)
	
	// Clear pixel buffers
	if d.pixelBuffer != nil {
		clear(d.pixelBuffer)
	}
	if d.grayscaleBuffer != nil {
		clear(d.grayscaleBuffer)
	}
	
	return nil
}

// SetPixel sets a single pixel in the internal buffer (1bpp)
// This implementation uses a sparse pixel buffer to avoid full framebuffer memory usage.
// Call Display() to render accumulated pixels to the e-paper display.
func (d *Device) SetPixel(x, y int16, c bool) {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return
	}
	
	// Initialize pixel buffer if needed
	if d.pixelBuffer == nil {
		d.pixelBuffer = make(map[uint32]bool)
	}
	
	// Use a single uint32 key to store x,y coordinates
	key := (uint32(y) << 16) | uint32(x)
	
	if c {
		d.pixelBuffer[key] = true
	} else {
		delete(d.pixelBuffer, key) // Remove false pixels to save memory
	}
}

// GetPixel gets a single pixel state from the internal buffer
func (d *Device) GetPixel(x, y int16) bool {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return false
	}
	
	if d.pixelBuffer == nil {
		return false
	}
	
	key := (uint32(y) << 16) | uint32(x)
	return d.pixelBuffer[key] // Returns false if key doesn't exist
}

// SetGrayscalePixel sets a pixel with grayscale value (0-15)
// This implementation uses a sparse pixel buffer to avoid full framebuffer memory usage.
// Call Display() to render accumulated pixels to the e-paper display.
func (d *Device) SetGrayscalePixel(x, y int16, c uint8) {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return
	}
	
	if c > 15 {
		c = 15
	}
	
	// Initialize grayscale buffer if needed
	if d.grayscaleBuffer == nil {
		d.grayscaleBuffer = make(map[uint32]uint8)
	}
	
	key := (uint32(y) << 16) | uint32(x)
	
	if c == 0 {
		delete(d.grayscaleBuffer, key) // Remove zero pixels to save memory
	} else {
		d.grayscaleBuffer[key] = c
	}
}

// GetGrayscalePixel gets a grayscale pixel value
func (d *Device) GetGrayscalePixel(x, y int16) uint8 {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return 0
	}
	
	if d.grayscaleBuffer == nil {
		return 0
	}
	
	key := (uint32(y) << 16) | uint32(x)
	return d.grayscaleBuffer[key] // Returns 0 if key doesn't exist
}