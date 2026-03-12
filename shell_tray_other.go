//go:build !darwin

package main

import (
	"github.com/energye/systray"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) startTray() {
	go func() {
		systray.Run(func() {
			systray.SetIcon(trayIcon)
			systray.SetTitle("")
			systray.SetTooltip("GoTube Downloader (Ctrl+Shift+Y)")

			mOpen := systray.AddMenuItem("Open GoTube", "Toggle overlay")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit GoTube", "Quit application")

			systray.SetOnClick(func(menu systray.IMenu) {
				a.ToggleOverlay()
			})
			systray.SetOnRClick(func(menu systray.IMenu) {
				menu.ShowMenu()
			})

			mOpen.Click(func() {
				a.ShowOverlay()
			})
			mQuit.Click(func() {
				wailsRuntime.Quit(a.ctx)
			})
		}, func() {})
	}()
}

func (a *App) stopTray() {
	systray.Quit()
}
