# Troubleshooting Checklist

Quick reference for common issues with the EPD47 driver and ESP32-S3 LilyGo board.

## ❌ Build Issues

### "package machine is not in std"
```bash
# ✅ Solution: Use TinyGo, not regular Go
tinygo build -target=esp32 ./examples/lilygo_simple.go
# Not: go build ./examples/lilygo_simple.go
```

### "undefined: machine.GPIO13"
```bash
# ✅ Solution: Use Pin() constructor or LilyGoT547 wrapper
# Instead of: machine.GPIO13
# Use: machine.Pin(13)
# Or better: epd47.NewLilyGoT547()
```

### "value is too big (66155 bytes)"
```bash
# ✅ Solution: Already fixed in v1.1.0+
# LUT size reduced from 64KB to 4KB
# Update to latest version
```

## ❌ Upload Issues

### Device not detected
```bash
# ✅ Check USB connection
lsusb | grep -i "ch340\|cp210"

# ✅ Check permissions (Linux)
sudo usermod -a -G dialout $USER
# Then log out and back in

# ✅ Try with sudo
sudo tinygo flash -target=esp32 ./examples/lilygo_simple.go
```

### Flash timeout/errors
```bash
# ✅ Try lower baud rate
esptool.py --chip esp32s3 --port /dev/ttyUSB0 --baud 115200 write_flash 0x0 firmware.bin

# ✅ Erase flash first
esptool.py --chip esp32s3 --port /dev/ttyUSB0 erase_flash

# ✅ Put board in download mode manually
# Hold BOOT button, press RESET, release BOOT
```

### Wrong serial port
```bash
# ✅ List available ports
# Linux/macOS:
ls /dev/tty*
# Look for /dev/ttyUSB0 or /dev/ttyACM0

# Windows: Check Device Manager
# Look for COM ports under "Ports (COM & LPT)"

# ✅ Specify port explicitly
tinygo flash -target=esp32 -port=/dev/ttyUSB0 ./examples/lilygo_simple.go
```

## ❌ Runtime Issues

### Display not updating
```bash
# ✅ Check power sequence
display := epd47.NewLilyGoT547()
display.Initialize()  // Must call this!

# ✅ Check if Configure() was called
device := epd47.New(config)
device.Configure()    // Required after New()
device.PowerOn()
```

### Garbled display output
```bash
# ✅ Check pin configuration
# Use LilyGoT547 for automatic pin setup
display := epd47.NewLilyGoT547()

# ✅ Or verify manual pin mapping matches hardware
# D0-D7 should be GPIO 8,1,2,3,4,5,6,7
```

### Memory issues
```bash
# ✅ Use pre-allocated drawing methods
display.DrawCheckerboard(x, y, w, h, blockSize)  // Good
# Instead of creating large byte arrays

# ✅ Check available memory
tinygo build -target=esp32 -print-sizes ./examples/lilygo_simple.go
```

## ❌ Development Issues

### Tests failing with machine package
```bash
# ✅ Tests use mock implementation automatically
go test ./epd47  # Works without TinyGo

# ✅ Build tags separate TinyGo from test code
# Real implementation: // +build tinygo
# Mock implementation: // +build !tinygo
```

### Interface compliance errors
```bash
# ✅ Device implements standard interfaces
var display epd47.Displayer = device
var grayscale epd47.GrayscaleDisplayer = device

# ✅ Use interface methods
w, h := display.Size()
display.SetPixel(x, 0, true)  // Note: only line 0 supported
```

## ❌ Hardware Issues

### Board not powering on
- ✅ Check USB-C cable (must support data, not just power)
- ✅ Try different USB port
- ✅ Check for driver installation (CH340/CP2102)

### Display stays blank
- ✅ Verify it's a LilyGo T5 4.7" ESP32-S3 (not other variants)
- ✅ Check display connection (should be built-in)
- ✅ Try the simple example first: `./examples/lilygo_simple.go`

### Partial display updates
- ✅ e-Paper displays are slow (several seconds for full update)
- ✅ 4bpp drawing takes longer than 1bpp (15-frame pipeline)
- ✅ Wait for operations to complete before power off

## 🔧 Quick Diagnostic Commands

```bash
# Check TinyGo installation
tinygo version

# List available targets
tinygo targets | grep esp32

# Check serial devices
ls /dev/tty* | grep -E "(USB|ACM)"

# Test build without flashing
tinygo build -target=esp32 -o test.bin ./examples/lilygo_simple.go

# Monitor serial output
screen /dev/ttyUSB0 115200
```

## 📚 Getting Help

1. **Start with simple example**: `./examples/lilygo_simple.go`
2. **Check hardware**: Verify board model and connections
3. **Test build**: Ensure code compiles before flashing
4. **Check serial output**: Look for error messages
5. **Review documentation**: [UPLOAD_GUIDE.md](UPLOAD_GUIDE.md) and [README.md](README.md)

## ✅ Working Configuration

If everything works, you should see:
```bash
$ tinygo flash -target=esp32 ./examples/lilygo_simple.go
   code    data     bss |   flash     ram
  xxxxx   xxxxx   xxxxx |   xxxxx   xxxxx

$ screen /dev/ttyUSB0 115200
Display: LilyGo T5 4.7" ESP32-S3 (ED047TC1)
Resolution: 960 x 540
# Display should show patterns and text
```