# TinyGo EPD47 Driver Makefile
# Simplifies common development tasks

# Variables
PROJECT_NAME := tinygo-epd47-s3
TINYGO_TARGET := esp32-coreboard-v2
TINYGO_TARGET_ALT := esp32
BUILD_DIR := build
EXAMPLES_DIR := examples
TEST_PACKAGE := ./epd47

# Default target
.PHONY: help
help: ## Show this help message
	@echo "TinyGo EPD47 Driver - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make test          # Run all tests"
	@echo "  make build-all     # Build all examples"
	@echo "  make flash-simple  # Flash simple example to device"
	@echo "  make clean         # Clean build artifacts"

# Setup and Installation
.PHONY: install-deps
install-deps: ## Install required dependencies (TinyGo, esptool)
	@echo "Installing dependencies..."
	@if ! command -v tinygo >/dev/null 2>&1; then \
		echo "❌ TinyGo not found. Please install TinyGo first."; \
		echo "   See: https://tinygo.org/getting-started/install/"; \
		exit 1; \
	fi
	@if ! command -v esptool.py >/dev/null 2>&1; then \
		echo "Installing esptool..."; \
		pip install esptool; \
	fi
	@echo "✅ Dependencies check complete"

.PHONY: check-env
check-env: ## Check development environment
	@echo "Checking environment..."
	@echo "Go version: $$(go version)"
	@if command -v tinygo >/dev/null 2>&1; then \
		echo "TinyGo version: $$(tinygo version)"; \
	else \
		echo "❌ TinyGo not found"; \
	fi
	@if command -v esptool.py >/dev/null 2>&1; then \
		echo "esptool version: $$(esptool.py version)"; \
	else \
		echo "⚠️  esptool not found (needed for flashing)"; \
	fi
	@echo "Available TinyGo ESP32 targets:"
	@tinygo targets | grep esp32 || echo "No ESP32 targets found"

# Testing
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test -v $(TEST_PACKAGE)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out $(TEST_PACKAGE)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race -v $(TEST_PACKAGE)

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem $(TEST_PACKAGE)

# Building
.PHONY: build-dir
build-dir:
	@mkdir -p $(BUILD_DIR)

.PHONY: build-simple
build-simple: build-dir ## Build simple example
	@echo "Building simple example..."
	tinygo build -target=$(TINYGO_TARGET) -o $(BUILD_DIR)/lilygo_simple.bin $(EXAMPLES_DIR)/lilygo_simple.go
	@echo "✅ Built: $(BUILD_DIR)/lilygo_simple.bin"

.PHONY: build-advanced
build-advanced: build-dir ## Build advanced example
	@echo "Building advanced example..."
	tinygo build -target=$(TINYGO_TARGET) -o $(BUILD_DIR)/lilygo_advanced.bin $(EXAMPLES_DIR)/lilygo_advanced.go
	@echo "✅ Built: $(BUILD_DIR)/lilygo_advanced.bin"

.PHONY: build-patterns
build-patterns: build-dir ## Build TinyGo patterns example
	@echo "Building TinyGo patterns example..."
	tinygo build -target=$(TINYGO_TARGET) -o $(BUILD_DIR)/tinygo_patterns.bin $(EXAMPLES_DIR)/tinygo_patterns.go
	@echo "✅ Built: $(BUILD_DIR)/tinygo_patterns.bin"

.PHONY: build-pixel-demo
build-pixel-demo: build-dir ## Build pixel interface demo
	@echo "Building pixel interface demo..."
	tinygo build -target=$(TINYGO_TARGET) -o $(BUILD_DIR)/pixel_interface_demo.bin $(EXAMPLES_DIR)/pixel_interface_demo.go
	@echo "✅ Built: $(BUILD_DIR)/pixel_interface_demo.bin"

.PHONY: build-performance
build-performance: build-dir ## Build performance demo
	@echo "Building performance demo..."
	tinygo build -target=$(TINYGO_TARGET) -o $(BUILD_DIR)/performance_demo.bin $(EXAMPLES_DIR)/performance_demo.go
	@echo "✅ Built: $(BUILD_DIR)/performance_demo.bin"

.PHONY: build-all
build-all: build-simple build-advanced build-patterns build-pixel-demo build-performance ## Build all examples
	@echo "✅ All examples built successfully"
	@ls -la $(BUILD_DIR)/

# Alternative target builds (if primary target fails)
.PHONY: build-simple-alt
build-simple-alt: build-dir ## Build simple example with alternative target
	@echo "Building simple example with alternative target..."
	tinygo build -target=$(TINYGO_TARGET_ALT) -o $(BUILD_DIR)/lilygo_simple_alt.bin $(EXAMPLES_DIR)/lilygo_simple.go
	@echo "✅ Built: $(BUILD_DIR)/lilygo_simple_alt.bin"

# Flashing
.PHONY: flash-simple
flash-simple: ## Flash simple example to device
	@echo "Flashing simple example..."
	tinygo flash -target=$(TINYGO_TARGET) $(EXAMPLES_DIR)/lilygo_simple.go

.PHONY: flash-advanced
flash-advanced: ## Flash advanced example to device
	@echo "Flashing advanced example..."
	tinygo flash -target=$(TINYGO_TARGET) $(EXAMPLES_DIR)/lilygo_advanced.go

.PHONY: flash-patterns
flash-patterns: ## Flash TinyGo patterns example to device
	@echo "Flashing TinyGo patterns example..."
	tinygo flash -target=$(TINYGO_TARGET) $(EXAMPLES_DIR)/tinygo_patterns.go

.PHONY: flash-pixel-demo
flash-pixel-demo: ## Flash pixel interface demo to device
	@echo "Flashing pixel interface demo..."
	tinygo flash -target=$(TINYGO_TARGET) $(EXAMPLES_DIR)/pixel_interface_demo.go

.PHONY: flash-performance
flash-performance: ## Flash performance demo to device
	@echo "Flashing performance demo..."
	tinygo flash -target=$(TINYGO_TARGET) $(EXAMPLES_DIR)/performance_demo.go

# Manual flashing with esptool
.PHONY: flash-manual
flash-manual: build-simple ## Flash using esptool (manual method)
	@echo "Flashing with esptool..."
	@if [ -z "$(PORT)" ]; then \
		echo "Usage: make flash-manual PORT=/dev/ttyUSB0"; \
		echo "Available ports:"; \
		ls /dev/tty* | grep -E "(USB|ACM)" || echo "No USB/ACM ports found"; \
		exit 1; \
	fi
	esptool.py --chip esp32s3 --port $(PORT) --baud 921600 write_flash 0x0 $(BUILD_DIR)/lilygo_simple.bin

# Serial monitoring
.PHONY: monitor
monitor: ## Monitor serial output (requires PORT variable)
	@if [ -z "$(PORT)" ]; then \
		echo "Usage: make monitor PORT=/dev/ttyUSB0"; \
		echo "Available ports:"; \
		ls /dev/tty* | grep -E "(USB|ACM)" || echo "No USB/ACM ports found"; \
		exit 1; \
	fi
	@echo "Monitoring $(PORT) - Press Ctrl+A then K to exit"
	screen $(PORT) 115200

# Development
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ go vet passed"

.PHONY: lint
lint: ## Run golint (if available)
	@echo "Running golint..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not found, skipping..."; \
	fi

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)
	@echo "✅ All checks passed"

# Size analysis
.PHONY: size
size: build-simple ## Show binary size information
	@echo "Binary size analysis:"
	@ls -lh $(BUILD_DIR)/lilygo_simple.bin
	@echo ""
	@echo "TinyGo size breakdown:"
	tinygo build -target=$(TINYGO_TARGET) -print-sizes $(EXAMPLES_DIR)/lilygo_simple.go

# Documentation
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all ./epd47 > docs.txt
	@echo "✅ Documentation generated: docs.txt"

# Git operations
.PHONY: git-status
git-status: ## Show git status
	@git status

.PHONY: git-clean-check
git-clean-check: ## Check if git working directory is clean
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "❌ Working directory is not clean"; \
		git status --short; \
		exit 1; \
	else \
		echo "✅ Working directory is clean"; \
	fi

# Release helpers
.PHONY: version
version: ## Show current version from git tags
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@echo "All tags:"
	@git tag -l --sort=version:refname

.PHONY: changelog
changelog: ## Show recent changelog entries
	@echo "Recent changes:"
	@head -20 CHANGELOG.md

# Cleanup
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html docs.txt
	@echo "✅ Clean complete"

.PHONY: clean-all
clean-all: clean ## Clean everything including go cache
	@echo "Cleaning go cache..."
	go clean -cache -testcache -modcache
	@echo "✅ Deep clean complete"

# Utility targets
.PHONY: list-examples
list-examples: ## List all available examples
	@echo "Available examples:"
	@ls -1 $(EXAMPLES_DIR)/*.go | sed 's/.*\//  /' | sed 's/\.go//'

.PHONY: list-targets
list-targets: ## List available TinyGo targets
	@echo "Available TinyGo ESP32 targets:"
	@tinygo targets | grep esp32

.PHONY: device-info
device-info: ## Show connected device information
	@echo "USB devices:"
	@lsusb | grep -i "ch340\|cp210\|esp32" || echo "No ESP32-related devices found"
	@echo ""
	@echo "Serial ports:"
	@ls /dev/tty* | grep -E "(USB|ACM)" || echo "No USB/ACM ports found"

# Quick development workflow
.PHONY: dev
dev: check build-simple ## Quick development check (test + build)
	@echo "✅ Development check complete"

.PHONY: quick-test
quick-test: build-simple flash-simple ## Quick test: build and flash simple example
	@echo "✅ Quick test complete - check your device!"

# Help for common workflows
.PHONY: workflow-help
workflow-help: ## Show common development workflows
	@echo "Common Development Workflows:"
	@echo ""
	@echo "🚀 First time setup:"
	@echo "  make install-deps check-env"
	@echo ""
	@echo "🔧 Development cycle:"
	@echo "  make dev                    # Test and build"
	@echo "  make flash-simple           # Flash to device"
	@echo "  make monitor PORT=/dev/ttyUSB0  # Monitor output"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test                   # Run tests"
	@echo "  make test-coverage          # With coverage"
	@echo ""
	@echo "📦 Building:"
	@echo "  make build-all              # Build all examples"
	@echo "  make size                   # Check binary size"
	@echo ""
	@echo "🔍 Debugging:"
	@echo "  make device-info            # Check connected devices"
	@echo "  make list-targets           # Show TinyGo targets"
	@echo ""
	@echo "For more commands: make help"

# Default target when no arguments provided
.DEFAULT_GOAL := help