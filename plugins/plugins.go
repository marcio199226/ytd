package plugins

type Plugin interface {
	GetName() string
	Initialize() error
	SetDir(dir string)
	Fetch(url string) error
	GetFilename() error
	Supports(address string) bool
}
