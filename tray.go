package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
)

type TrayMenu struct {
	defaultTrayMenu *menu.TrayMenu
}

func (tray *TrayMenu) createTray() *menu.TrayMenu {
	tray.defaultTrayMenu = &menu.TrayMenu{
		Label: "ytd",
		Menu:  tray.createTrayMenu(),
	}

	return tray.defaultTrayMenu
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
	m.Append(&menu.MenuItem{
		Type:     menu.CheckboxType,
		Label:    "Check for updates", // hide window on close
		Disabled: true,
		Checked:  appState.Config.CheckForUpdates,
		Click: func(ctx *menu.CallbackData) {
			watch, err := appState.ReadSettingBoolValue("CheckForUpdates")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray CheckForUpdates error: %s", err))
			}
			appState.SaveSettingBoolValue("CheckForUpdates", !watch)
			ctx.MenuItem.Checked = appState.Config.CheckForUpdates
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
		},
	})
	m.Append(&menu.MenuItem{
		Type:     menu.CheckboxType,
		Label:    "Start at login (system startup)",
		Disabled: true,
		Checked:  appState.Config.StartAtLogin,
		Click: func(ctx *menu.CallbackData) {
			watch, err := appState.ReadSettingBoolValue("StartAtLogin")
			if err != nil {
				appState.runtime.Log.Error(fmt.Sprintf("Tray StartAtLogin error: %s", err))
			}
			appState.SaveSettingBoolValue("StartAtLogin", !watch)

			ctx.MenuItem.Checked = appState.Config.StartAtLogin
			appState.runtime.Events.Emit("ytd:app:config", appState.Config)
			// @TODO: update to latest commit or wait for next release
			// mac.StartAtLogin(ctx.MenuItem.Checked)
			appState.runtime.Dialog.Message(&dialog.MessageDialog{
				Type:         dialog.InfoDialog,
				Title:        "Update successful",
				Message:      "Please restart app for the changes to take effect.",
				Buttons:      []string{"OK"},
				CancelButton: "OK",
			})
		},
	})
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
	m.Append(&menu.MenuItem{
		Type:  menu.TextType,
		Label: fmt.Sprintf("ytd (%s)", version),
	})
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
