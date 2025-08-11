[![Built with Kiro](https://img.shields.io/badge/built_with_Kiro-9046ff.svg?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyMCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDIwIDI0IiBmaWxsPSJub25lIj48cGF0aCBkPSJNMy44MDA4MSAxOC41NjYxQzEuMzIzMDYgMjQuMDU3MiA2LjU5OTA0IDI1LjQzNCAxMC40OTA0IDIyLjIyMDVDMTEuNjMzOSAyNS44MjQyIDE1LjkyNiAyMy4xMzYxIDE3LjQ2NTIgMjAuMzQ0NUMyMC44NTc4IDE0LjE5MTUgMTkuNDg3NyA3LjkxNDU5IDE5LjEzNjEgNi42MTk4OEMxNi43MjQ0IC0yLjIwOTcyIDQuNjcwNTUgLTIuMjE4NTIgMi41OTU4MSA2LjY2NDlDMi4xMTEzNiA4LjIxOTQ2IDIuMTAyODQgOS45ODc1MiAxLjgyODQ2IDExLjgyMzNDMS42OTAxMSAxMi43NDkgMS41OTI1OCAxMy4zMzk4IDEuMjM0MzYgMTQuMzEzNUMxLjAyODQxIDE0Ljg3MzMgMC43NDUwNDMgMTUuMzcwNCAwLjI5OTgzMyAxNi4yMDgyQy0wLjM5MTU5NCAxNy41MDk1IC0wLjA5OTg4MDIgMjAuMDIxIDMuNDYzOTcgMTguNzE4NlYxOC43MTk1TDMuODAwODEgMTguNTY2MVoiIGZpbGw9IndoaXRlIj48L3BhdGg+PHBhdGggZD0iTTEwLjk2MTQgMTAuNDQxM0M5Ljk3MjAyIDEwLjQ0MTMgOS44MjQyMiA5LjI1ODkzIDkuODI0MjIgOC41NTQwN0M5LjgyNDIyIDcuOTE3OTEgOS45MzgyNCA3LjQxMjQgMTAuMTU0MiA3LjA5MTk3QzEwLjM0NDEgNi44MTAwMyAxMC42MTU4IDYuNjY2OTkgMTAuOTYxNCA2LjY2Njk5QzExLjMwNzEgNi42NjY5OSAxMS42MDM2IDYuODEyMjggMTEuODEyOCA3LjA5ODkyQzEyLjA1MTEgNy40MjU1NCAxMi4xNzcgNy45Mjg2MSAxMi4xNzcgOC41NTQwN0MxMi4xNzcgOS43MzU5MSAxMS43MjI2IDEwLjQ0MTMgMTAuOTYxNiAxMC40NDEzSDEwLjk2MTRaIiBmaWxsPSJibGFjayI+PC9wYXRoPjxwYXRoIGQ9Ik0xNS4wMzE4IDEwLjQ0MTNDMTQuMDQyMyAxMC40NDEzIDEzLjg5NDUgOS4yNTg5MyAxMy44OTQ1IDguNTU0MDdDMTMuODk0NSA3LjkxNzkxIDE0LjAwODYgNy40MTI0IDE0LjIyNDUgNy4wOTE5N0MxNC40MTQ0IDYuODEwMDMgMTQuNjg2MSA2LjY2Njk5IDE1LjAzMTggNi42NjY5OUMxNS4zNzc0IDYuNjY2OTkgMTUuNjczOSA2LjgxMjI4IDE1Ljg4MzEgNy4wOTg5MkMxNi4xMjE0IDcuNDI1NTQgMTYuMjQ3NCA3LjkyODYxIDE2LjI0NzQgOC41NTQwN0MxNi4yNDc0IDkuNzM1OTEgMTUuNzkzIDEwLjQ0MTMgMTUuMDMxOSAxMC40NDEzSDE1LjAzMThaIiBmaWxsPSJibGFjayI+PC9wYXRoPjwvc3ZnPg==)](https://kiro.dev)

# TinyGo EPD47 Driver for LilyGo ESP32-S3 T5 4.7"

A TinyGo driver for the 4.7" e-Paper display (ED047TC1) on the LilyGo ESP32-S3 T5 board. This implementation follows TinyGo driver design patterns and is based on the original [LilyGo-EPD47](https://github.com/Xinyuan-LilyGO/LilyGo-EPD47) PlatformIO library.

## Features

- **TinyGo driver patterns**: Implements standard `Displayer` and `GrayscaleDisplayer` interfaces
- **1bpp (monochrome) drawing**: Fast black and white graphics
- **4bpp (grayscale) drawing**: 16-level grayscale with proper dithering
- **Hardware abstraction**: Pin-agnostic design using function pointers
- **Memory efficient**: Pre-allocated line buffers to avoid heap allocations
- **Power management**: Proper power sequencing for the e-paper display
- **Build tag support**: Separate implementations for TinyGo and testing

## Hardware Support

- **Display**: ED047TC1 4.7" e-Paper panel (960x540 pixels)
- **Controller**: ESP32-S3 with 8-bit parallel interface
- **Board**: LilyGo T5 4.7" ESP32-S3

## Interfaces

The driver implements standard TinyGo display interfaces:

- `Displayer`: Basic display operations (Size, SetPixel, GetPixel, Display, ClearDisplay)
- `GrayscaleDisplayer`: Extends Displayer with grayscale operations

```go
// Use as standard TinyGo display interface
var display epd47.Displayer = device
w, h := display.Size()
display.SetPixel(10, 0, true)  // Note: simplified implementation
```

## Pin Configuration

The driver uses the following pins on the ESP32-S3:

```go
// Configuration shift register
CFG_DATA: GPIO13
CFG_CLK:  GPIO12  
CFG_STR:  GPIO0

// Control signals
CKV: GPIO38  // Vertical gate clock
STH: GPIO40  // Start/enable
CKH: GPIO41  // Write strobe

// 8-bit data bus (D0-D7)
D0: GPIO8
D1: GPIO1
D2: GPIO2
D3: GPIO3
D4: GPIO4
D5: GPIO5
D6: GPIO6
D7: GPIO7
```

## Usage

### Quick Start (Recommended)

For the LilyGo T5 4.7" ESP32-S3 board, use the preconfigured device:

```go
package main

import (
    "time"
    "github.com/abaschen/tinygo-epd47-s3/epd47"
)

func main() {
    // Create preconfigured device (all pins automatically set)
    display := epd47.NewLilyGoT547()
    
    // Initialize (power on + clear)
    display.Initialize()
    
    // Draw a checkerboard
    display.DrawCheckerboard(100, 100, 200, 150, 16)
    
    // Draw a gradient
    display.DrawGradient(350, 100, 200, 150, true, epd47.BlackOnWhite)
    
    // Draw some "text"
    display.DrawText(100, 300, "Hello EPD47!", 12, 20)
    
    time.Sleep(5 * time.Second)
    
    // Shutdown
    display.Shutdown()
    
    for {
        time.Sleep(time.Hour)
    }
}
```

### Manual Configuration

For custom pin configurations or other boards:

```go
package main

import (
    "machine"
    "time"
    "github.com/abaschen/tinygo-epd47-s3/epd47"
)

func pinOut(p machine.Pin) epd47.PinOut {
    p.Configure(machine.PinConfig{Mode: machine.PinOutput})
    return func(level bool) {
        if level {
            p.High()
        } else {
            p.Low()
        }
    }
}

func sleepUS(us int) {
    if us <= 0 {
        return
    }
    time.Sleep(time.Duration(us) * time.Microsecond)
}

func main() {
    cfg := epd47.Config{
        Width:  960,
        Height: 540,
        
        // Pin assignments
        CFG_DATA: pinOut(machine.Pin(13)),
        CFG_CLK:  pinOut(machine.Pin(12)),
        CFG_STR:  pinOut(machine.Pin(0)),
        
        CKV: pinOut(machine.Pin(38)),
        STH: pinOut(machine.Pin(40)),
        CKH: pinOut(machine.Pin(41)),
        
        D0: pinOut(machine.Pin(8)),
        D1: pinOut(machine.Pin(1)),
        D2: pinOut(machine.Pin(2)),
        D3: pinOut(machine.Pin(3)),
        D4: pinOut(machine.Pin(4)),
        D5: pinOut(machine.Pin(5)),
        D6: pinOut(machine.Pin(6)),
        D7: pinOut(machine.Pin(7)),
        
        SleepUS: sleepUS,
    }
    
    d := epd47.New(cfg)
    
    // Configure the device (required after New())
    err := d.Configure()
    if err != nil {
        // Handle error
    }
    
    // Power on the display
    d.PowerOn()
    time.Sleep(200 * time.Millisecond)
    
    // Clear the screen
    d.Clear(2)
    
    // Your drawing code here...
    
    // Power off when done
    d.PowerOffAll()
}
```

### Preconfigured Device Methods

The `LilyGoT547` device provides convenient methods:

```go
display := epd47.NewLilyGoT547()

// Initialization and shutdown
display.Initialize()                    // Power on + clear
display.Shutdown()                      // Power off all

// Drawing helpers
display.DrawCheckerboard(x, y, w, h, blockSize)
display.DrawRectangle(x, y, w, h, filled)
display.DrawGradient(x, y, w, h, horizontal, mode)
display.DrawText(x, y, text, charWidth, charHeight)

// Get display info
width, height, model := display.GetDisplayInfo()
```

### Low-Level Drawing Operations

#### 1bpp (Monochrome) Drawing

```go
// Create a 1bpp bitmap (MSB first, packed)
w, h := 200, 120
src := make([]byte, (w+7)/8*h)

// Draw a checkerboard pattern
for r := 0; r < h; r++ {
    for c := 0; c < w; c++ {
        on := ((r>>4)&1)^((c>>4)&1) == 1
        if on {
            byteIdx := r*((w+7)/8) + (c >> 3)
            bit := 7 - (c & 7)
            src[byteIdx] |= (1 << bit)
        }
    }
}

// Draw at position (x, y) with 10μs pulse width
d.Draw1bpp(x, y, w, h, src, 10)
```

#### 4bpp (Grayscale) Drawing

```go
// Create a 4bpp bitmap (2 pixels per byte)
w, h := 240, 120
src := make([]byte, (w/2+w%2)*h)

// Draw a gradient
for r := 0; r < h; r++ {
    for c := 0; c < w; c++ {
        shade := byte((c * 15) / (w - 1)) // 0-15 grayscale
        byteIdx := r*((w/2)+(w%2)) + (c >> 1)
        if (c & 1) == 0 {
            src[byteIdx] = (shade << 4) | (src[byteIdx] & 0x0F)
        } else {
            src[byteIdx] = (src[byteIdx] & 0xF0) | (shade & 0x0F)
        }
    }
}

// Draw with black-on-white mode
d.DrawImage4bpp(x, y, w, h, src, epd47.BlackOnWhite)
```

### Drawing Modes

The driver supports three drawing modes for 4bpp images:

- `epd47.BlackOnWhite`: Standard black text on white background
- `epd47.WhiteOnWhite`: White text on white background (erasing)
- `epd47.WhiteOnBlack`: White text on black background (inverse)

## Building and Uploading

### Quick Start
```bash
# Build and flash in one command
tinygo flash -target=esp32 ./examples/lilygo_simple.go

# Or build manually
tinygo build -target=esp32 -o firmware.bin ./examples/lilygo_simple.go
```

### Detailed Instructions
- **[UPLOAD_GUIDE.md](UPLOAD_GUIDE.md)**: Complete setup and flashing instructions
- **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)**: Quick fixes for common issues

Key steps:
1. Install TinyGo and esptool
2. Connect board via USB-C
3. Run `tinygo flash -target=esp32 ./examples/lilygo_simple.go`
4. Monitor with `screen /dev/ttyUSB0 115200`

## Architecture

The driver is organized into several modules:

- `device.go`: Core device structure and configuration
- `ed047tc1.go`: Hardware control and power management
- `bus_parallel.go`: 8-bit parallel bus communication
- `grayscale.go`: 4bpp grayscale rendering with LUT
- `examples/`: Usage examples
  - `lilygo_simple.go`: **Recommended** - Simple example using preconfigured device
  - `lilygo_advanced.go`: Advanced demo with complex patterns and animations
  - `main.go`: Original example with checkerboard and gradient
  - `generic_main.go`: Generic example using Pin() constructor
  - `demo.go`: Comprehensive demo showing all features

## Performance Notes

- **1bpp drawing**: Fast, suitable for text and simple graphics
- **4bpp drawing**: Slower (15 frame pipeline), but provides smooth grayscales
- **Memory usage**: Fixed buffers avoid heap allocations during drawing
- **Power management**: Proper sequencing prevents display damage

## Limitations

- Currently requires TinyGo with ESP32-S3 support
- 8-bit parallel interface only (no SPI mode)
- Fixed 960x540 resolution
- No partial refresh support

## License

This driver is based on the LilyGo-EPD47 library and follows similar licensing terms.