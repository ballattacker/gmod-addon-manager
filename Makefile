# Variables
GO := go
BUILD_DIR := build
BIN_NAME := gmod-addon-manager

# Default target (shows help)
all:
	@echo "Usage: make [linux|windows|release]"

# Build for Windows
windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-windows.exe

# Build for Linux
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-linux

# Create annotated release tag and push
release:
	@echo "Creating annotated release tag $(TAG)..."
	@git tag -a $(TAG) -m "Release $(TAG)"
	@git push origin $(TAG)
	@echo "Annotated tag $(TAG) created and pushed successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

# Show help
help:
	@echo "Available targets:"
	@echo "  make windows   - Build for Windows"
	@echo "  make linux     - Build for Linux"
	@echo "  make release   - Create and push an annotated release tag (requires TAG=vX.Y.Z)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help"
