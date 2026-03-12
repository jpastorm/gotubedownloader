# ─────────────────────────────────────────────────
# GoTube Downloader – Makefile
# Cross-compile for macOS, Windows & Linux
# ─────────────────────────────────────────────────

APP_NAME    := GotubeDownloader
VERSION     := 1.0.0
BUILD_DIR   := build/bin

# ── Default: build for current platform ──────────
.PHONY: all build dev clean \
        build-darwin build-darwin-arm64 build-darwin-amd64 \
        build-windows build-linux \
        build-all icon-windows \
        open kill dmg

all: build

# ── Development mode ─────────────────────────────
dev:
	wails dev

# ── Build for current OS/arch ────────────────────
build:
	wails build

# ── macOS builds ─────────────────────────────────
build-darwin: build-darwin-arm64

build-darwin-arm64:
	@echo "▸ Building macOS (arm64)..."
	wails build -platform darwin/arm64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME).app"

build-darwin-amd64:
	@echo "▸ Building macOS (amd64)..."
	wails build -platform darwin/amd64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME).app"

build-darwin-universal: build-darwin-arm64
	@echo "▸ Building macOS (amd64) for universal..."
	wails build -platform darwin/amd64 -o $(APP_NAME)-amd64
	@echo "▸ Creating universal binary..."
	lipo -create -output "$(BUILD_DIR)/$(APP_NAME)-universal" \
		"$(BUILD_DIR)/$(APP_NAME).app/Contents/MacOS/$(APP_NAME)" \
		"$(BUILD_DIR)/$(APP_NAME)-amd64.app/Contents/MacOS/$(APP_NAME)-amd64"
	@echo "✓ Universal binary created"

# ── Windows build ────────────────────────────────
build-windows:
	@echo "▸ Building Windows (amd64)..."
	wails build -platform windows/amd64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME).exe"

build-windows-arm64:
	@echo "▸ Building Windows (arm64)..."
	wails build -platform windows/arm64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME).exe"

# ── Linux build ──────────────────────────────────
build-linux:
	@echo "▸ Building Linux (amd64)..."
	wails build -platform linux/amd64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME)"

build-linux-arm64:
	@echo "▸ Building Linux (arm64)..."
	wails build -platform linux/arm64
	@echo "✓ $(BUILD_DIR)/$(APP_NAME)"

# ── Build all platforms ──────────────────────────
build-all: build-darwin-arm64 build-windows build-linux
	@echo ""
	@echo "════════════════════════════════════════"
	@echo "  All platforms built successfully!"
	@echo "════════════════════════════════════════"

# ── Generate Windows .ico from icon.png ──────────
icon-windows:
	@echo "▸ Generating Windows icon..."
	magick icon.png -define icon:auto-resize=256,128,64,48,32,16 build/windows/icon.ico
	@echo "✓ build/windows/icon.ico"

# ── Generate bindings ────────────────────────────
bindings:
	wails generate module

# ── Clean build artifacts ────────────────────────
clean:
	@echo "▸ Cleaning..."
	rm -rf $(BUILD_DIR)/*
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite
	@echo "✓ Clean"

# ── Open / Kill the app (macOS) ──────────────────
open:
	@echo "▸ Opening $(APP_NAME).app..."
	@open "$(BUILD_DIR)/$(APP_NAME).app"

kill:
	@echo "▸ Killing $(APP_NAME)..."
	@pkill -x $(APP_NAME) 2>/dev/null && echo "✓ Killed" || echo "⚠ Not running"

# ── Create .dmg installer (macOS) ────────────────
dmg:
	@echo "▸ Creating DMG..."
	@rm -f "$(BUILD_DIR)/$(APP_NAME).dmg"
	@mkdir -p "$(BUILD_DIR)/$(APP_NAME).iconset"
	@sips -z 16 16     icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_16x16.png"    >/dev/null
	@sips -z 32 32     icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_16x16@2x.png" >/dev/null
	@sips -z 32 32     icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_32x32.png"    >/dev/null
	@sips -z 64 64     icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_32x32@2x.png" >/dev/null
	@sips -z 128 128   icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_128x128.png"  >/dev/null
	@sips -z 256 256   icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_128x128@2x.png" >/dev/null
	@sips -z 256 256   icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_256x256.png"  >/dev/null
	@sips -z 512 512   icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_256x256@2x.png" >/dev/null
	@sips -z 512 512   icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_512x512.png"  >/dev/null
	@sips -z 1024 1024 icon.png --out "$(BUILD_DIR)/$(APP_NAME).iconset/icon_512x512@2x.png" >/dev/null
	@iconutil -c icns "$(BUILD_DIR)/$(APP_NAME).iconset" -o "$(BUILD_DIR)/$(APP_NAME).icns"
	@rm -rf "$(BUILD_DIR)/$(APP_NAME).iconset"
	@create-dmg \
		--volname "$(APP_NAME)" \
		--volicon "$(BUILD_DIR)/$(APP_NAME).icns" \
		--window-pos 200 120 \
		--window-size 600 400 \
		--icon-size 100 \
		--icon "$(APP_NAME).app" 150 185 \
		--hide-extension "$(APP_NAME).app" \
		--app-drop-link 450 185 \
		"$(BUILD_DIR)/$(APP_NAME).dmg" \
		"$(BUILD_DIR)/$(APP_NAME).app"
	@rm -f "$(BUILD_DIR)/$(APP_NAME).icns"
	@echo "✓ $(BUILD_DIR)/$(APP_NAME).dmg"

# ── Install frontend deps ───────────────────────
deps:
	cd frontend && npm install

# ── Help ─────────────────────────────────────────
help:
	@echo ""
	@echo "GoTube Downloader – Build targets"
	@echo "─────────────────────────────────────"
	@echo "  make dev                  Run in dev mode"
	@echo "  make build                Build for current platform"
	@echo "  make build-darwin-arm64   macOS Apple Silicon"
	@echo "  make build-darwin-amd64   macOS Intel"
	@echo "  make build-windows        Windows amd64"
	@echo "  make build-windows-arm64  Windows arm64"
	@echo "  make build-linux          Linux amd64"
	@echo "  make build-linux-arm64    Linux arm64"
	@echo "  make build-all            All main platforms"
	@echo "  make icon-windows         Regenerate Windows .ico"
	@echo "  make bindings             Regenerate Wails bindings"
	@echo "  make deps                 Install frontend npm deps"
	@echo "  make open                 Open the built .app (macOS)"
	@echo "  make kill                 Kill running app (macOS)"
	@echo "  make dmg                  Create .dmg installer (macOS)"
	@echo "  make clean                Remove build artifacts"
	@echo ""
