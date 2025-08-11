package epd47

import (
	"testing"
	"time"
)

// Mock pin implementation for testing
func mockPinOut(level bool) {
	// Mock implementation - does nothing
}

func mockSleep(us int) {
	if us <= 0 {
		return
	}
	// Use a much shorter sleep for testing
	time.Sleep(time.Duration(us) * time.Nanosecond)
}

func TestDeviceCreation(t *testing.T) {
	cfg := Config{
		Width:  960,
		Height: 540,

		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,

		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,

		D0: mockPinOut,
		D1: mockPinOut,
		D2: mockPinOut,
		D3: mockPinOut,
		D4: mockPinOut,
		D5: mockPinOut,
		D6: mockPinOut,
		D7: mockPinOut,

		SleepUS: mockSleep,
	}

	d := New(cfg)
	if d == nil {
		t.Fatal("New() returned nil")
	}
	
	err := d.Configure()
	if err != nil {
		t.Fatalf("Configure() failed: %v", err)
	}

	if d.Width() != 960 {
		t.Errorf("Expected width 960, got %d", d.Width())
	}

	if d.Height() != 540 {
		t.Errorf("Expected height 540, got %d", d.Height())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Width != MaxWidth {
		t.Errorf("Expected width %d, got %d", MaxWidth, cfg.Width)
	}
	if cfg.Height != MaxHeight {
		t.Errorf("Expected height %d, got %d", MaxHeight, cfg.Height)
	}
}

func TestPowerSequence(t *testing.T) {
	cfg := Config{
		Width:  960,
		Height: 540,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test power on/off sequence (should not panic)
	d.PowerOn()
	d.PowerOff()
	d.PowerOffAll()
}

func TestDrawing1bpp(t *testing.T) {
	cfg := Config{
		Width:  960,
		Height: 540,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test 1bpp drawing (should not panic)
	w, h := 100, 50
	src := make([]byte, (w+7)/8*h)
	d.Draw1bpp(0, 0, w, h, src, 10)
	
	// Test clear
	d.Clear(1)
}

func TestDrawing4bpp(t *testing.T) {
	cfg := Config{
		Width:  960,
		Height: 540,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test 4bpp drawing (should not panic)
	w, h := 100, 50
	src := make([]byte, (w/2+w%2)*h)
	d.DrawImage4bpp(0, 0, w, h, src, BlackOnWhite)
}

func TestBounds(t *testing.T) {
	cfg := Config{
		Width:  100,  // Smaller for testing
		Height: 100,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test drawing outside bounds (should not panic)
	w, h := 50, 50
	src1 := make([]byte, (w+7)/8*h)
	src4 := make([]byte, (w/2+w%2)*h)
	
	// These should be clipped or ignored gracefully
	d.Draw1bpp(-10, -10, w, h, src1, 10)  // Negative position
	d.Draw1bpp(200, 200, w, h, src1, 10)  // Outside bounds
	d.DrawImage4bpp(-10, -10, w, h, src4, BlackOnWhite)
	d.DrawImage4bpp(200, 200, w, h, src4, BlackOnWhite)
}

func TestDisplayerInterface(t *testing.T) {
	cfg := Config{
		Width:  100,
		Height: 100,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test Displayer interface
	var disp Displayer = d
	
	w, h := disp.Size()
	if w != 100 || h != 100 {
		t.Errorf("Expected size 100x100, got %dx%d", w, h)
	}
	
	// Test pixel operations (only first line supported in this implementation)
	disp.SetPixel(10, 0, true)
	if !disp.GetPixel(10, 0) {
		t.Error("Pixel should be set")
	}
	
	disp.SetPixel(10, 0, false)
	if disp.GetPixel(10, 0) {
		t.Error("Pixel should be cleared")
	}
	
	// Test bounds
	disp.SetPixel(-1, -1, true) // Should not panic
	disp.SetPixel(200, 200, true) // Should not panic
	
	// Test interface methods
	err := disp.Display()
	if err != nil {
		t.Errorf("Display() failed: %v", err)
	}
	
	err = disp.ClearDisplay()
	if err != nil {
		t.Errorf("ClearDisplay() failed: %v", err)
	}
}

func TestGrayscaleDisplayerInterface(t *testing.T) {
	cfg := Config{
		Width:  100,
		Height: 100,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test GrayscaleDisplayer interface
	var gdisp GrayscaleDisplayer = d
	
	// Test grayscale pixel operations (only first line supported in this implementation)
	gdisp.SetGrayscalePixel(10, 0, 8)
	if gdisp.GetGrayscalePixel(10, 0) != 8 {
		t.Errorf("Expected grayscale value 8, got %d", gdisp.GetGrayscalePixel(10, 0))
	}
	
	// Test bounds
	gdisp.SetGrayscalePixel(-1, -1, 15) // Should not panic
	gdisp.SetGrayscalePixel(200, 200, 15) // Should not panic
	
	// Test value clamping
	gdisp.SetGrayscalePixel(20, 0, 20) // Should clamp to 15
	if gdisp.GetGrayscalePixel(20, 0) != 15 {
		t.Errorf("Expected clamped value 15, got %d", gdisp.GetGrayscalePixel(20, 0))
	}
}

func TestBufferClearing(t *testing.T) {
	cfg := Config{
		Width:  100,
		Height: 100,
		CFG_DATA: mockPinOut,
		CFG_CLK:  mockPinOut,
		CFG_STR:  mockPinOut,
		CKV: mockPinOut,
		STH: mockPinOut,
		CKH: mockPinOut,
		D0: mockPinOut, D1: mockPinOut, D2: mockPinOut, D3: mockPinOut,
		D4: mockPinOut, D5: mockPinOut, D6: mockPinOut, D7: mockPinOut,
		SleepUS: mockSleep,
	}

	d := New(cfg)
	d.Configure()
	
	// Test fillBuffer function
	testBuf := make([]byte, 100)
	
	// Fill with non-zero value
	fillBuffer(testBuf, 0xAA)
	for i, v := range testBuf {
		if v != 0xAA {
			t.Errorf("Expected 0xAA at index %d, got 0x%02X", i, v)
		}
	}
	
	// Clear with zero (should use clear() internally)
	fillBuffer(testBuf, 0x00)
	for i, v := range testBuf {
		if v != 0x00 {
			t.Errorf("Expected 0x00 at index %d, got 0x%02X", i, v)
		}
	}
	
	// Test device buffer clearing methods
	d.clearLineBuffer()
	for i, v := range d.line1b[:10] { // Check first 10 bytes
		if v != 0 {
			t.Errorf("line1b not cleared at index %d, got 0x%02X", i, v)
		}
	}
	
	d.clearGrayscaleBuffer()
	for i, v := range d.line4b[:10] { // Check first 10 bytes
		if v != 0 {
			t.Errorf("line4b not cleared at index %d, got 0x%02X", i, v)
		}
	}
}