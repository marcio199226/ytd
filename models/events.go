package models

import "github.com/wailsapp/wails/v2"

type LoaderEventPayload struct {
	Label        string `json:"label"`
	TemplateName string `json:"templateName"`
}

func ShowLoader(r *wails.Runtime, label string) {
	r.Events.Emit("ytd:loader:show", LoaderEventPayload{Label: label, TemplateName: "card"})
}

func HideLoader(r *wails.Runtime, label string) {
	r.Events.Emit("ytd:loader:hide")
}
