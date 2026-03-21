# Project Variables
APP_NAME := bdoPF
BUILD_TAGS := webkit2_41
OUTPUT_DIR := build/bin

# Default target: builds everything
.PHONY: all
all: clean win-amd64 win-arm64 linux-amd64 linux-arm64

# Display help information
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          Clean and build all platforms"
	@echo "  win-amd64    Build for Windows x64"
	@echo "  win-arm64    Build for Windows ARM64"
	@echo "  linux-amd64  Build for Linux x64"
# 	@echo "  linux-arm64  Build for Linux ARM64"
	@echo "  clean        Remove all build artifacts"

# Build for Windows AMD64
.PHONY: win-amd64
win-amd64:
	@echo "Building Windows AMD64..."
	wails build -platform windows/amd64 -o $(APP_NAME)_windows_amd64.exe

# Build for Windows ARM64
.PHONY: win-arm64
win-arm64:
	@echo "Building Windows ARM64..."
	wails build -platform windows/arm64 -o $(APP_NAME)_windows_arm64.exe

# Build for Linux AMD64
.PHONY: linux-amd64
linux-amd64:
	@echo "Building Linux AMD64..."
	wails build -tags $(BUILD_TAGS) -platform linux/amd64 -o $(APP_NAME)_linux_amd64

# Build for Linux ARM64
# .PHONY: linux-arm64
# linux-arm64:
# 	@echo "Building Linux ARM64..."
# 	# We set CC to the cross-compiler we just installed
# 	CC=aarch64-linux-gnu-gcc wails build -tags $(BUILD_TAGS) -platform linux/arm64 -o $(APP_NAME)_linux_arm64

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build directory..."
	rm -rf $(OUTPUT_DIR)/*