package models

import "github.com/wailsapp/wails/v2"

type LoaderEventPayload struct {
	Label        string `json:"label"`
	TemplateName string `json:"templateName"`
}

type NotificationEventPayload struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type NgrokStateEventPayload struct {
	ErrCode string `json:"errCode"`
	Status  string `json:"status"`
	Url     string `json:"url "`
}

func ShowLoader(r *wails.Runtime, label string) {
	r.Events.Emit("ytd:loader:show", LoaderEventPayload{Label: label, TemplateName: "card"})
}

func HideLoader(r *wails.Runtime) {
	r.Events.Emit("ytd:loader:hide")
}
