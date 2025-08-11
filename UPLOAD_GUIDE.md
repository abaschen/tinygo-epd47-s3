# Upload Guide: ESP32-S3 LilyGo EPD 4.7" S3 Device

This guide explains how to compile and upload TinyGo firmware to the LilyGo T5 4.7" ESP32-S3 e-Paper display board.

## Quick Reference (For Experienced Users)

```bash
# 1. Install TinyGo (Linux ARM64 example)
wget https://github.com/tinygo-org/tinygo/releases/download/v0.38.0/tinygo0.38.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf tinygo0.38.0.linux-arm64.tar.gz
export PATH=$PATH:/usr/local/tinygo/bin

# 2. Install esptool
pip install esptool

# 3. Flash firmware (one command)
tinygo flash -target=esp32 ./examples/lilygo_simple.go
# If esp32 target has issues, try:
# tinygo flash -target=esp32-coreboard-v2 ./examples/lilygo_simple.go

# 4. Monitor output
screen /dev/ttyUSB0 115200
```

## Prerequisites

### 1. Install TinyGo

#### Linux (ARM64)
```bash
# Download and install TinyGo
wget https://github.com/tinygo-org/tinygo/releases/download/v0.38.0/tinygo0.38.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf tinygo0.38.0.linux-arm64.tar.gz
export PATH=$PATH:/usr/local/tinygo/bin

# Verify installation
tinygo version
```

#### Linux (x86_64)
```bash
wget https://github.com/tinygo-org/tinygo/releases/download/v0.38.0/tinygo0.38.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf tinygo0.38.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/tinygo/bin
```

#### macOS
```bash
# Using Homebrew
brew tap tinygo-org/tools
brew install tinygo

# Or download manually
wget https://github.com/tinygo-org/tinygo/releases/download/v0.38.0/tinygo0.38.0.darwin-amd64.tar.gz
sudo tar -C /usr/local -xzf tinygo0.38.0.darwin-amd64.tar.gz
```

#### Windows
1. Download `tinygo0.38.0.windows-amd64.zip` from [TinyGo releases](https://github.com/tinygo-org/tinygo/releases)
2. Extract to `C:\tinygo`
3. Add `C:\tinygo\bin` to your PATH

### 2. Install Required Tools

#### esptool (for flashing)
```bash
# Python/pip method (recommended)
pip install esptool

# Or using package manager
# Ubuntu/Debian:
sudo apt install esptool
# macOS:
brew install esptool
```

#### USB-to-Serial Drivers
The LilyGo T5 4.7" uses a CH340 or CP2102 USB-to-serial chip. Install the appropriate driver:

- **CH340**: [Download from manufacturer](http://www.wch-ic.com/downloads/CH341SER_EXE.html)
- **CP2102**: [Download from Silicon Labs](https://www.silabs.com/developers/usb-to-uart-bridge-vcp-drivers)

## Hardware Setup

### 1. Board Identification
- **Board**: LilyGo T5 4.7" ESP32-S3
- **Display**: ED047TC1 4.7" e-Paper (960x540)
- **MCU**: ESP32-S3-WROOM-1
- **USB**: USB-C connector

### 2. Physical Connections
1. Connect the board to your computer using a USB-C cable
2. Ensure the cable supports data transfer (not just charging)
3. The board should appear as a serial device (e.g., `/dev/ttyUSB0`, `/dev/ttyACM0`, or `COM3`)

### 3. Check Device Detection
```bash
# Linux/macOS - list serial devices
ls /dev/tty*

# Look for devices like:
# /dev/ttyUSB0 (CH340)
# /dev/ttyACM0 (CP2102)

# Verify USB device is detected
lsusb | grep -i "ch340\|cp210"

# Windows - check Device Manager
# Look for "Silicon Labs CP210x" or "CH340" under Ports (COM & LPT)
```

### 4. Hardware Verification
Before flashing, verify you have the correct board:

- **Model**: LilyGo T5 4.7" ESP32-S3 (not ESP32 or other variants)
- **Display**: 4.7" e-Paper display should be attached
- **Markings**: Look for "T5-4.7" or similar on the PCB
- **USB**: USB-C connector (not micro-USB)

**⚠️ Important**: This driver is specifically for the ESP32-S3 variant with the ED047TC1 display. Other LilyGo boards may have different pin configurations.

**📝 Note**: Some TinyGo versions may have issues with ESP32 targets. If you encounter build errors, try:
1. Updating to the latest TinyGo version
2. Using `esp32-coreboard-v2` target as an alternative
3. Building for `arduino` target for syntax verification (won't run on hardware)

## Building and Uploading

### Method 1: Using TinyGo Flash Command (Recommended)

#### 1. Build and Flash in One Step
```bash
# Navigate to your project directory
cd /path/to/tinygo-epd47-s3

# Flash the simple example
tinygo flash -target=esp32 ./examples/lilygo_simple.go

# Flash the advanced demo
tinygo flash -target=esp32 ./examples/lilygo_advanced.go

# Flash with custom patterns
tinygo flash -target=esp32 ./examples/tinygo_patterns.go
```

#### 2. Specify Serial Port (if needed)
```bash
# Linux/macOS
tinygo flash -target=esp32 -port=/dev/ttyUSB0 ./examples/lilygo_simple.go

# Windows
tinygo flash -target=esp32 -port=COM3 ./examples/lilygo_simple.go
```

### Method 2: Manual Build and Flash

#### 1. Build the Firmware
```bash
# Build firmware binary
tinygo build -target=esp32 -o firmware.bin ./examples/lilygo_simple.go

# Verify the binary was created
ls -la firmware.bin
```

#### 2. Put Board in Download Mode
1. Hold the **BOOT** button (if present)
2. Press and release the **RESET** button
3. Release the **BOOT** button
4. The board is now in download mode

#### 3. Flash Using esptool
```bash
# Flash the firmware
esptool.py --chip esp32s3 --port /dev/ttyUSB0 --baud 921600 write_flash 0x0 firmware.bin

# Windows example
esptool.py --chip esp32s3 --port COM3 --baud 921600 write_flash 0x0 firmware.bin
```

#### 4. Reset the Board
Press the **RESET** button or disconnect/reconnect USB to run the new firmware.

## Troubleshooting

### Common Issues

#### 1. Device Not Detected
```bash
# Check if device is connected
lsusb | grep -i "ch340\|cp210"

# Check permissions (Linux)
sudo usermod -a -G dialout $USER
# Log out and back in for changes to take effect

# Or use sudo temporarily
sudo tinygo flash -target=esp32 ./examples/lilygo_simple.go
```

#### 2. Flash Errors
```bash
# Try lower baud rate
esptool.py --chip esp32s3 --port /dev/ttyUSB0 --baud 115200 write_flash 0x0 firmware.bin

# Erase flash first
esptool.py --chip esp32s3 --port /dev/ttyUSB0 erase_flash
```

#### 3. Board Not Entering Download Mode
- Some boards enter download mode automatically
- Try holding BOOT while connecting USB
- Check if your USB cable supports data transfer

#### 4. Build Errors
```bash
# Ensure TinyGo is properly installed
tinygo version

# Check available targets
tinygo targets | grep esp32

# Try building without flashing first
tinygo build -target=esp32 -o test.bin ./examples/lilygo_simple.go
```

### Serial Monitor

#### View Debug Output
```bash
# Using screen (Linux/macOS)
screen /dev/ttyUSB0 115200

# Using minicom (Linux)
minicom -D /dev/ttyUSB0 -b 115200

# Using PuTTY (Windows)
# Set connection type to Serial, speed to 115200
```

## Example Projects

### 1. Simple Demo
```bash
# Basic functionality test
tinygo flash -target=esp32 ./examples/lilygo_simple.go
```

### 2. Advanced Demo
```bash
# Complex patterns and animations
tinygo flash -target=esp32 ./examples/lilygo_advanced.go
```

### 3. TinyGo Patterns
```bash
# Demonstrates TinyGo interface usage
tinygo flash -target=esp32 ./examples/tinygo_patterns.go
```

### 4. Custom Application
```go
package main

import (
    "time"
    "github.com/abaschen/tinygo-epd47-s3/epd47"
)

func main() {
    display := epd47.NewLilyGoT547()
    display.Initialize()
    
    // Your custom code here
    display.DrawText(100, 200, "Hello World!", 16, 24)
    
    time.Sleep(10 * time.Second)
    display.Shutdown()
    
    for {
        time.Sleep(time.Hour)
    }
}
```

## Performance Tips

### 1. Optimize Build Size
```bash
# Build with size optimization
tinygo build -target=esp32 -opt=s -o firmware.bin ./examples/lilygo_simple.go
```

### 2. Enable Debug Info
```bash
# Build with debug information
tinygo build -target=esp32 -debug -o firmware.bin ./examples/lilygo_simple.go
```

### 3. Monitor Memory Usage
```bash
# Check binary size
ls -lh firmware.bin

# TinyGo will show memory usage during build
tinygo build -target=esp32 -print-sizes ./examples/lilygo_simple.go
```

## Board-Specific Notes

### LilyGo T5 4.7" ESP32-S3 Specifications
- **MCU**: ESP32-S3-WROOM-1 (Dual-core Xtensa LX7, 240MHz)
- **Flash**: 16MB
- **PSRAM**: 8MB
- **Display**: 4.7" e-Paper, 960x540, 16 grayscale levels
- **Interface**: 8-bit parallel
- **Power**: USB-C, 3.7V Li-Po battery connector

### Pin Configuration (Automatic with LilyGoT547)
- **Data Bus**: D0-D7 (GPIO 8,1,2,3,4,5,6,7)
- **Control**: CKV(38), STH(40), CKH(41)
- **Config**: CFG_DATA(13), CFG_CLK(12), CFG_STR(0)

### Power Management
- The driver handles proper power sequencing
- Battery life depends on refresh frequency
- e-Paper displays retain image when powered off

## Additional Resources

- [TinyGo ESP32 Documentation](https://tinygo.org/docs/reference/microcontrollers/esp32/)
- [LilyGo T5 4.7" Hardware Repository](https://github.com/Xinyuan-LilyGO/LilyGo-EPD47)
- [ESP32-S3 Technical Reference](https://www.espressif.com/sites/default/files/documentation/esp32-s3_technical_reference_manual_en.pdf)
- [esptool Documentation](https://docs.espressif.com/projects/esptool/en/latest/)

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Verify hardware connections and drivers
3. Test with the simple examples first
4. Check TinyGo and esptool versions
5. Review the project's GitHub issues

Happy coding with your e-Paper display! 🎉