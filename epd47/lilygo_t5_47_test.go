// +build tinygo

package epd47

import (
	"testing"
)

func TestLilyGoT547Creation(t *testing.T) {
	display := NewLilyGoT547()
	if display == nil {
		t.Fatal("NewLilyGoT547() returned nil")
	}
	
	if display.Device == nil {
		t.Fatal("Device is nil")
	}
	
	// Check dimensions
	if display.Width() != 960 {
		t.Errorf("Expected width 960, got %d", display.Width())
	}
	
	if display.Height() != 540 {
		t.Errorf("Expected height 540, got %d", display.Height())
	}
}

func TestLilyGoT547Info(t *testing.T) {
	display := NewLilyGoT547()
	
	width, height, model := display.GetDisplayInfo()
	
	if width != 960 {
		t.Errorf("Expected width 960, got %d", width)
	}
	
	if height != 540 {
		t.Errorf("Expected height 540, got %d", height)
	}
	
	expectedModel := "LilyGo T5 4.7\" ESP32-S3 (ED047TC1)"
	if model != expectedModel {
		t.Errorf("Expected model %q, got %q", expectedModel, model)
	}
}

func TestLilyGoT547DrawingMethods(t *testing.T) {
	display := NewLilyGoT547()
	
	// Test all drawing methods (should not panic)
	display.DrawCheckerboard(0, 0, 100, 100, 10)
	display.DrawRectangle(0, 0, 50, 50, false)
	display.DrawRectangle(0, 0, 50, 50, true)
	display.DrawGradient(0, 0, 100, 100, true, BlackOnWhite)
	display.DrawGradient(0, 0, 100, 100, false, BlackOnWhite)
	display.DrawText(0, 0, "Test", 8, 16)
}

func TestLilyGoT547PowerManagement(t *testing.T) {
	display := NewLilyGoT547()
	
	// Test power management (should not panic)
	display.Initialize()
	display.Shutdown()
}

func TestLilyGoT547BoundaryConditions(t *testing.T) {
	display := NewLilyGoT547()
	
	// Test with zero/negative dimensions (should not panic)
	display.DrawCheckerboard(0, 0, 0, 0, 10)
	display.DrawRectangle(0, 0, -10, -10, false)
	display.DrawGradient(0, 0, 0, 100, true, BlackOnWhite)
	display.DrawText(0, 0, "", 8, 16)
	
	// Test with very large dimensions (should be clipped)
	display.DrawCheckerboard(0, 0, 2000, 1000, 10)
	display.DrawRectangle(0, 0, 2000, 1000, false)
}