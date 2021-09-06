package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/mac"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
)

type TrayMenu struct {
	runtime              *wails.Runtime
	defaultTrayMenu      *menu.TrayMenu
	updateMenuItem       *menu.MenuItem
	startAtLoginMenuItem *menu.MenuItem
	versionMenuItem      *menu.MenuItem
}

func (tray *TrayMenu) createTray() *menu.TrayMenu {
	tray.defaultTrayMenu = &menu.TrayMenu{
		Label: "ytd",
		Menu:  tray.createTrayMenu(),
	}

	return tray.defaultTrayMenu
}

func (tray *TrayMenu) reRenderTray(callback func()) *menu.TrayMenu {
	tray.runtime.Menu.DeleteTrayMenu(tray.defaultTrayMenu)
	m := tray.createTray()
	callback()
	tray.runtime.Menu.SetTrayMenu(tray.defaultTrayMenu)
	return m
}

func (tray *TrayMenu) createTrayMenu() *menu.Menu {
	m := &menu.Menu{}
	m.Append(&menu.MenuItem{
		Type:    menu.CheckboxType,
		Label:   "Clipboard watch",
		Checked: appState.Config.ClipboardWatch,
		Click: func(ctx *menu.CallbackData) {
			watch, err := appState.ReadSettingBoolValue("ClipboardWatch")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray clipboard watch error: %s", err))
			}
			appState.SaveSettingBoolValue("ClipboardWatch", !watch)
			ctx.MenuItem.Checked = appState.Config.ClipboardWatch
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
		},
	})
	m.Append(&menu.MenuItem{
		Type:    menu.CheckboxType,
		Label:   "Run in background on close", // hide window on close
		Checked: appState.Config.RunInBackgroundOnClose,
		Click: func(ctx *menu.CallbackData) {
			watch, err := appState.ReadSettingBoolValue("RunInBackgroundOnClose")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray RunInBackgroundOnClose error: %s", err))
			}
			appState.SaveSettingBoolValue("RunInBackgroundOnClose", !watch)
			ctx.MenuItem.Checked = appState.Config.RunInBackgroundOnClose
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
		},
	})
	tray.updateMenuItem = &menu.MenuItem{
		Type:    menu.CheckboxType,
		Label:   "Check for updates", // hide window on close
		Checked: appState.Config.CheckForUpdates,
		Click: func(ctx *menu.CallbackData) {
			watch, err := appState.ReadSettingBoolValue("CheckForUpdates")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray CheckForUpdates error: %s", err))
			}
			appState.SaveSettingBoolValue("CheckForUpdates", !watch)
			ctx.MenuItem.Checked = appState.Config.CheckForUpdates
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
		},
	}
	m.Append(tray.updateMenuItem)
	tray.startAtLoginMenuItem = &menu.MenuItem{
		Type:     menu.CheckboxType,
		Label:    "Start at login (system startup)",
		Checked:  appState.canStartAtLogin && appState.Config.StartAtLogin,
		Disabled: !appState.canStartAtLogin,
		Click: func(ctx *menu.CallbackData) {
			enabled, err := appState.ReadSettingBoolValue("StartAtLogin")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray StartAtLogin error: %s", err))
				return
			}

			notAvailable := mac.StartAtLogin(!enabled)
			if notAvailable != nil {
				tray.reRenderTray(func() {
					tray.startAtLoginMenuItem.Label = "⚠ Start at Login unavailable"
					tray.startAtLoginMenuItem.Disabled = true
				})
				appState.runtime.Dialog.Message(&dialog.MessageDialog{
					Type:         dialog.ErrorDialog,
					Title:        "Cannot enable start at login",
					Message:      notAvailable.Error(),
					Buttons:      []string{"OK"},
					CancelButton: "OK",
				})
				return
			}

			ctx.MenuItem.Checked = !enabled
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
			appState.SaveSettingBoolValue("StartAtLogin", !enabled)
		},
	}
	m.Append(tray.startAtLoginMenuItem)
	m.Append(menu.Separator())
	m.Append(&menu.MenuItem{
		Type:  menu.TextType,
		Label: "Settings",
		Click: func(ctx *menu.CallbackData) {
			appState.ShowWindow()
			appState.runtime.Events.Emit("ytd:show:dialog:settings")
		},
	})
	m.Append(&menu.MenuItem{
		Type:  menu.TextType,
		Label: "Show app",
		Click: func(ctx *menu.CallbackData) {
			appState.ShowWindow()
		},
	})
	tray.versionMenuItem = &menu.MenuItem{
		Type:     menu.TextType,
		Disabled: true,
		Label:    fmt.Sprintf("ytd (%s)", version),
	}
	m.Append(tray.versionMenuItem)
	m.Append(&menu.MenuItem{
		Type:        menu.TextType,
		Label:       "Quit app",
		Accelerator: keys.CmdOrCtrl("q"),
		Hidden:      false,
		Click: func(ctx *menu.CallbackData) {
			appState.ForceQuit()
		},
	})
	return m
}
