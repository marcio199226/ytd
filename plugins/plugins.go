package plugins

import (
	. "ytd/models"

	"github.com/wailsapp/wails"
)

type Plugin interface {
	GetName() string
	Initialize() error
	SetDir(dir string)
	Fetch(url string) (*GenericEntry, error)
	GetFilename() error
	Supports(address string) bool
	SetWailsRuntime(*wails.Runtime)
}

type GenericEntry struct {
	Type     string          `json:"type"`
	Source   string          `json:"source"`
	Track    GenericTrack    `json:"track"`
	Playlist GenericPlaylist `json:"playlist"`
}

func NewGenericEntry(source string) *GenericEntry {
	return &GenericEntry{Source: source}
}
