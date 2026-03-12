package main

import (
	_ "embed"
	"fmt"
	"sync/atomic"

	"golang.design/x/hotkey"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed icon.png
var trayIcon []byte

//go:embed menu_icon_18x18.png
var menuIcon1x []byte

//go:embed menu_icon_36x36@2x.png
var menuIcon2x []byte

// overlayVisible tracks window state for toggle logic
var overlayVisible atomic.Bool

// appInstance is set during initShell so platform-specific tray callbacks can reach the App.
var appInstance *App

// ===== LIFECYCLE =====

// initShell starts tray icon and global hotkey in background goroutines.
// Called from App.startup after all core services are initialized.
func (a *App) initShell() {
	appInstance = a
	a.startTray() // platform-specific (shell_tray_*.go)
	go a.runHotkey()
}

// cleanupShell tears down tray icon. Called from App.shutdown.
func (a *App) cleanupShell() {
	a.stopTray() // platform-specific
}

// ===== GLOBAL HOTKEY =====

func (a *App) runHotkey() {
	hk := newGlobalHotkey()
	if err := hk.Register(); err != nil {
		fmt.Println("[shell] global hotkey registration failed:", err)
		fmt.Println("[shell] On macOS: grant Accessibility permission in System Settings → Privacy & Security → Accessibility")
		return
	}
	fmt.Println("[shell] global hotkey registered: Ctrl+Shift+Y")
	defer hk.Unregister()

	for range hk.Keydown() {
		a.ToggleOverlay()
	}
}

// newGlobalHotkey is defined per-platform in shell_hotkey_*.go files
// All platforms: Ctrl+Shift+Y
var _ = func() *hotkey.Hotkey { return newGlobalHotkey() } // compile-time check

// ===== OVERLAY CONTROL =====

// ShowOverlay makes the overlay visible, always-on-top, centered, and emits event for frontend focus.
func (a *App) ShowOverlay() {
	overlayVisible.Store(true)
	wailsRuntime.WindowShow(a.ctx)
	wailsRuntime.WindowSetAlwaysOnTop(a.ctx, true)
	wailsRuntime.WindowCenter(a.ctx)
	wailsRuntime.EventsEmit(a.ctx, "overlay:shown")
}

// HideOverlay hides the overlay and updates state.
func (a *App) HideOverlay() {
	overlayVisible.Store(false)
	wailsRuntime.WindowHide(a.ctx)
}

// ToggleOverlay shows or hides the overlay based on current state.
func (a *App) ToggleOverlay() {
	if overlayVisible.Load() {
		a.HideOverlay()
	} else {
		a.ShowOverlay()
	}
}

// quitApp terminates the application. Called from tray callbacks.
func (a *App) quitApp() {
	wailsRuntime.Quit(a.ctx)
}
