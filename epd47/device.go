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
	w, h int
	p    pins
	t    timeFns

	// Config register shadow.
	cfg reg

	// Data pin availability mask (bit i set => Di provided).
	dataMask uint8

	// Preallocated line buffers to avoid per-row allocations.
	// 1bpp: width/8 bytes per line.
	line1b [MaxWidthBytes1bpp]byte
	// 4bpp panel-lane formatted output: width/2 bytes per line.
	line4b [MaxWidthBytes4bpp]byte
}

// Hardware/format limits for this panel.
const (
	MaxWidth          = 960
	MaxHeight         = 540
	MaxWidthBytes1bpp = MaxWidth / 8
	MaxWidthBytes4bpp = MaxWidth / 2
	Frames4bpp        = 15
)

// Internal pin bundle (function pointers, no machine dependency).
type pins struct {
	// Config shift
	cfgData, cfgClk, cfgStr PinOut
	// Control
	ckv, sth, ckh PinOut
	// Data bus D0..D7
	d [8]PinOut
}

// Time function bundle.
type timeFns struct {
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
// It does not perform any panel power transitions; call PowerOn() next.
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
	pinsArr := [8]PinOut{cfg.D0, cfg.D1, cfg.D2, cfg.D3, cfg.D4, cfg.D5, cfg.D6, cfg.D7}
	var mask uint8
	for i := 0; i < 8; i++ {
		if pinsArr[i] != nil {
			mask |= 1 << uint(i)
		}
	}

	// Assemble pins bundle.
	p := pins{
		cfgData: cfg.CFG_DATA,
		cfgClk:  cfg.CFG_CLK,
		cfgStr:  cfg.CFG_STR,
		ckv:     cfg.CKV,
		sth:     cfg.STH,
		ckh:     cfg.CKH,
		d:       pinsArr,
	}

	d := &Device{
		w: w, h: h,
		p: p,
		t: timeFns{sleepUS: sl},
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
		dataMask: mask,
	}

	// Push initial (all-safed) config.
	d.pushCfg()
	return d
}

// Width returns the panel width in pixels.
func (d *Device) Width() int { return d.w }

// Height returns the panel height in pixels.
func (d *Device) Height() int { return d.h }
