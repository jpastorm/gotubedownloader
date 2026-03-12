//go:build darwin

package main

/*
#cgo LDFLAGS: -framework Cocoa

extern void nativeSetupTray(const void *icon1x, int icon1xLen,
                            const void *icon2x, int icon2xLen);
extern void nativeTeardownTray(void);
*/
import "C"

import "unsafe"

// startTray creates the macOS menu-bar icon via native Cocoa APIs.
func (a *App) startTray() {
	C.nativeSetupTray(
		unsafe.Pointer(&menuIcon1x[0]), C.int(len(menuIcon1x)),
		unsafe.Pointer(&menuIcon2x[0]), C.int(len(menuIcon2x)),
	)
}

// stopTray removes the menu-bar icon.
func (a *App) stopTray() {
	C.nativeTeardownTray()
}

//export goTrayToggle
func goTrayToggle() {
	if appInstance != nil {
		// MUST dispatch to a goroutine — calling Wails runtime directly
		// from a CGo/ObjC callback crashes (wrong stack/thread context).
		go appInstance.ToggleOverlay()
	}
}

//export goTrayQuit
func goTrayQuit() {
	if appInstance != nil {
		go appInstance.quitApp()
	}
}
