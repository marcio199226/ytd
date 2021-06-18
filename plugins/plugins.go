package plugins

import (
	"github.com/wailsapp/wails"
)

type Plugin interface {
	GetName() string
	Initialize() error
	SetDir(dir string)
	Fetch(url string)
	GetFilename() error
	Supports(address string) bool
	SetWailsRuntime(*wails.Runtime)
}
