# Remote Claude Code Desktop Build Configuration

.PHONY: all deps build build-all clean

# Default target
all: deps build

# Install dependencies
deps:
	cd frontend && npm install
	wails generate module

# Development mode
dev:
	wails dev

# Build for current platform
build:
	wails build

# Build for all platforms
build-all: build-mac build-windows build-linux

# Build for macOS
build-mac:
	wails build -platform darwin/amd64
	wails build -platform darwin/arm64

# Build for Windows
build-windows:
	wails build -platform windows/amd64

# Build for Linux
build-linux:
	wails build -platform linux/amd64

# Build with backend embedded
build-standalone:
	@echo "Building standalone application with embedded backend..."
	@# First build the backend
	cd ../backend && go build -o ../desktop/build/bin/backend ./cmd/server
	@# Then build the desktop app
	wails build

# Clean build artifacts
clean:
	rm -rf build/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

# Run in development mode with hot reload
run: dev

# Package for distribution
package-mac:
	@echo "Creating macOS DMG..."
	@if command -v create-dmg >/dev/null 2>&1; then \
		create-dmg build/bin/Remote\ Claude\ Code.app; \
	else \
		echo "create-dmg not found. Install with: brew install create-dmg"; \
	fi

# Generate icons
icons:
	@echo "Generating application icons..."
	@# You'll need to provide a source image
	@if [ -f "build/appicon.png" ]; then \
		wails generate icon; \
	else \
		echo "Please place an icon image at build/appicon.png (512x512 PNG)"; \
	fi

# Help
help:
	@echo "Remote Claude Code Desktop Build System"
	@echo ""
	@echo "Targets:"
	@echo "  deps          Install dependencies"
	@echo "  dev           Run in development mode"
	@echo "  build         Build for current platform"
	@echo "  build-all     Build for all platforms"
	@echo "  build-mac     Build for macOS (amd64 + arm64)"
	@echo "  build-windows Build for Windows"
	@echo "  build-linux   Build for Linux"
	@echo "  clean         Remove build artifacts"
	@echo "  help          Show this help message"
