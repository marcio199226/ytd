package models

import (
	"fmt"

	"github.com/leonelquinteros/gotext"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
)

type LoaderEventPayload struct {
	Label        string `json:"label"`
	TemplateName string `json:"templateName"`
}

type NotificationEventPayload struct {
	Label string `json:"label"`
	Type  string `json:"type"` // success | warning | error
}

type NgrokStateEventPayload struct {
	ErrCode string `json:"errCode"`
	Status  string `json:"status"`
	Url     string `json:"url"`
}

func ShowLoader(r *wails.Runtime, label string) {
	r.Events.Emit("ytd:loader:show", LoaderEventPayload{Label: label, TemplateName: "card"})
}

func HideLoader(r *wails.Runtime) {
	r.Events.Emit("ytd:loader:hide")
}

func SendNotification(r *wails.Runtime, payload NotificationEventPayload, isForeground bool) {
	if isForeground {
		r.Events.Emit("ytd:notification", payload)
		return
	}

	types := map[string]dialog.DialogType{
		"error":   dialog.ErrorDialog,
		"warning": dialog.WarningDialog,
		"success": dialog.InfoDialog,
	}
	_, err := r.Dialog.Message(&dialog.MessageDialog{
		Type:         types[payload.Type],
		Title:        gotext.Get("An error occured"),
		Message:      payload.Label,
		Buttons:      []string{"OK"},
		CancelButton: "OK",
	})
	if err != nil {
		fmt.Println("SendNotification error", err)
		return
	}
}
