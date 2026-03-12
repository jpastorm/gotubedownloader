package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Minimal menu — global hotkey (Ctrl+Shift+Y) handles toggle, not the menu
	appMenu := menu.NewMenu()
	appMenu.Append(menu.AppMenu())  // macOS app menu (About, Quit)
	appMenu.Append(menu.EditMenu()) // REQUIRED for ⌘C/⌘V/⌘X/⌘A in frameless windows

	windowMenu := appMenu.AddSubmenu("Window")
	windowMenu.AddText("Close", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {
		app.HideOverlay()
	})

	err := wails.Run(&options.App{
		Title:             "GoTube",
		Width:             680,
		Height:            480,
		MinWidth:          680,
		MinHeight:         480,
		MaxWidth:          680,
		MaxHeight:         480,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         true,
		StartHidden:       true,
		HideWindowOnClose: true,
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Menu:              appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               true,
				FullSizeContent:            true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "GoTube Downloader",
				Message: "Fast YouTube Downloader\nv1.0.0",
			},
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
