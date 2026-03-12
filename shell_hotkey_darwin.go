//go:build darwin

package main

import "golang.design/x/hotkey"

// newGlobalHotkey returns Ctrl+Shift+Y (universal shortcut, Y = YouTube)
func newGlobalHotkey() *hotkey.Hotkey {
	return hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyY)
}
