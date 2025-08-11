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
	// Push initial (all-safed) config.
	d.pushCfg()
	return nil
}

// Width returns the panel width in pixels.
func (d *Device) Width() int { return d.w }

// Height returns the panel height in pixels.
func (d *Device) Height() int { return d.h }
// Size returns the display dimensions as required by Displayer interface
func (d *Device) Size() (x, y int16) {
	return int16(d.w), int16(d.h)
}

// Display updates the screen - for e-paper this is handled by drawing operations
func (d *Device) Display() error {
	// E-paper displays update immediately during draw operations
	// This is a no-op but satisfies the interface
	return nil
}

// ClearDisplay clears the entire display
func (d *Device) ClearDisplay() error {
	d.Clear(2)
	return nil
}

// SetPixel sets a single pixel in the internal buffer (1bpp)
// Note: This is a simplified implementation for interface compliance
// For actual drawing, use Draw1bpp() method
func (d *Device) SetPixel(x, y int16, c bool) {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return
	}
	
	// For this implementation, we'll use the line buffer as a single-line buffer
	// This is mainly for interface compliance and testing
	if int(y) == 0 { // Only support first line for simplicity
		byteIdx := int(x >> 3)
		bit := 7 - (int(x) & 7)
		
		if byteIdx < len(d.line1b) {
			if c {
				d.line1b[byteIdx] |= (1 << bit)
			} else {
				d.line1b[byteIdx] &^= (1 << bit)
			}
		}
	}
}

// GetPixel gets a single pixel state from the internal buffer
func (d *Device) GetPixel(x, y int16) bool {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return false
	}
	
	// For this implementation, we'll use the line buffer as a single-line buffer
	if int(y) == 0 { // Only support first line for simplicity
		byteIdx := int(x >> 3)
		bit := 7 - (int(x) & 7)
		
		if byteIdx < len(d.line1b) {
			return (d.line1b[byteIdx]>>bit)&1 == 1
		}
	}
	return false
}

// SetGrayscalePixel sets a pixel with grayscale value (0-15)
// Note: This is a simplified implementation for interface compliance
// For actual drawing, use DrawImage4bpp() method
func (d *Device) SetGrayscalePixel(x, y int16, c uint8) {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return
	}
	
	if c > 15 {
		c = 15
	}
	
	// For this implementation, we'll use the line buffer as a single-line buffer
	if int(y) == 0 { // Only support first line for simplicity
		byteIdx := int(x >> 1)
		if byteIdx < len(d.line4b) {
			if int(x)&1 == 0 {
				// Even column - upper nibble
				d.line4b[byteIdx] = (c << 4) | (d.line4b[byteIdx] & 0x0F)
			} else {
				// Odd column - lower nibble
				d.line4b[byteIdx] = (d.line4b[byteIdx] & 0xF0) | (c & 0x0F)
			}
		}
	}
}

// GetGrayscalePixel gets a grayscale pixel value
func (d *Device) GetGrayscalePixel(x, y int16) uint8 {
	if x < 0 || y < 0 || int(x) >= d.w || int(y) >= d.h {
		return 0
	}
	
	// For this implementation, we'll use the line buffer as a single-line buffer
	if int(y) == 0 { // Only support first line for simplicity
		byteIdx := int(x >> 1)
		if byteIdx < len(d.line4b) {
			if int(x)&1 == 0 {
				// Even column - upper nibble
				return (d.line4b[byteIdx] >> 4) & 0x0F
			} else {
				// Odd column - lower nibble
				return d.line4b[byteIdx] & 0x0F
			}
		}
	}
	return 0
}