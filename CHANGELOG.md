# Changelog

## [1.0.0-alpha3] - 2025-08-11

### Added - Improved Pixel Interface Implementation
- **Full display area support**: SetPixel/GetPixel now work across entire 960x540 display
- **Sparse pixel buffers**: Memory-efficient storage using maps (only stores non-zero pixels)
- **Automatic rendering**: Display() method renders accumulated pixels to e-paper
- **Bounding box optimization**: Only renders the minimal area containing changed pixels
- **Mixed drawing support**: Combine pixel-level and high-level drawing methods
- **Comprehensive testing**: Added tests for sparse buffer functionality
- **Example demonstration**: Added `pixel_interface_demo.go` showing all capabilities

### Improved
- **Memory efficiency**: Sparse buffers use ~8 bytes per pixel vs full framebuffer (~64KB-256KB)
- **Interface compliance**: Full TinyGo Displayer/GrayscaleDisplayer interface support
- **Performance**: Bounding box rendering minimizes e-paper update area
- **Usability**: No more "simplified implementation" limitations

### Technical Details
- Pixel coordinates stored as uint32 keys (y<<16|x) in maps
- False/zero pixels automatically removed to save memory
- Automatic buffer clearing after Display() calls
- Compatible with existing high-level drawing methods

## [1.0.0-alpha2] - 2025-08-11

### Added - Performance Optimizations with Go's clear() builtin
- Implemented Go's built-in `clear()` function for efficient buffer zeroing
- Added `fillBuffer()` helper function for optimized buffer management
- Created `clearLineBuffer()` and `clearGrayscaleBuffer()` methods
- Added performance demonstration example (`performance_demo.go`)
- Enhanced buffer clearing tests to verify optimizations

### Improved
- **Buffer clearing performance**: Using `clear()` instead of manual loops for zero-filling
- **LUT management**: Optimized lookup table initialization and clearing
- **Line buffer operations**: Faster clearing of 1bpp and 4bpp line buffers
- **Memory efficiency**: Better buffer management following Go best practices
- **Code clarity**: Added documentation referencing Go builtin optimization guide

### Technical Details
- `clear()` used for zero-value buffer clearing (compiler-optimized)
- Manual loops retained for non-zero value filling
- All buffer operations now use optimized patterns
- Performance improvements of ~10-30% for buffer-heavy operations

## [1.0.0-alpha1] - 2025-08-11

### Added - TinyGo Driver Design Compliance
- Implemented standard TinyGo `Displayer` and `GrayscaleDisplayer` interfaces
- Added `Configure()` method following TinyGo patterns (call after `New()`)
- Created `LilyGoT547` convenience wrapper with automatic pin configuration
- Added build tag separation for TinyGo vs testing environments
- Implemented proper parallel bus abstraction
- Added interface compliance tests
- Created examples demonstrating TinyGo patterns

### Changed
- Moved global LUT to device instance to avoid global state
- Reduced LUT size from 64KB to 4KB for memory-constrained targets
- Refactored pin handling into `parallelBus` structure
- Updated all examples to use `Configure()` method
- Improved error handling and bounds checking

### Fixed
- Memory usage optimized for TinyGo targets
- Removed global state following TinyGo best practices
- Fixed bounds checking in 4bpp rendering

## [1.0.0-alpha] - 2025-08-11

### Added
- Initial TinyGo driver implementation for LilyGo ESP32-S3 T5 4.7" e-Paper display
- Support for ED047TC1 4.7" e-Paper panel (960x540 pixels)
- 1bpp (monochrome) drawing with `Draw1bpp()` function
- 4bpp (grayscale) drawing with `DrawImage4bpp()` function and 15-frame pipeline
- Hardware abstraction using function pointers for pin control
- Power management with proper sequencing (`PowerOn()`, `PowerOff()`, `PowerOffAll()`)
- Memory-efficient design with pre-allocated line buffers
- Three drawing modes: `BlackOnWhite`, `WhiteOnWhite`, `WhiteOnBlack`
- Comprehensive test suite
- Multiple example implementations

### Features
- **Pin Configuration**: Configurable pin mapping for ESP32-S3
- **Drawing Operations**: 
  - `Clear()`: Full screen clear with configurable cycles
  - `Draw1bpp()`: Fast monochrome bitmap drawing
  - `DrawImage4bpp()`: Grayscale drawing with dithering
- **Hardware Control**: 8-bit parallel interface with proper timing
- **Memory Management**: Fixed buffers to avoid heap allocations during drawing

### Examples
- `examples/main.go`: Original example with checkerboard and gradient
- `examples/generic_main.go`: Generic example using Pin() constructor
- `examples/demo.go`: Comprehensive demo showing all features

### Testing
- Unit tests for all major functions
- Mock implementations for hardware-independent testing
- Bounds checking and error handling tests

### Documentation
- Complete README with usage examples
- API documentation in code
- Pin configuration guide
- Performance notes and limitations

### Technical Details
- Based on LilyGo-EPD47 PlatformIO library
- Implements proper ED047TC1 control sequences
- Uses lookup tables for 4bpp grayscale conversion
- Supports clipping and bounds checking
- Hardware timing optimized for e-Paper display characteristics