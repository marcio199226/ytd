package models

type AppConfig struct {
	ClipboardWatch              bool
	DownloadOnCopy              bool
	ConcurrentDownloads         bool
	ConcurrentPlaylistDownloads bool
	MaxParrallelDownloads       uint8
	BaseSaveDir                 string
	Proxy                       interface{}
}

func NewAppConfig(watch bool, dldOnCopy bool, cDownloads bool, cPlaylistDownloads bool, baseSaveDir string) AppConfig {
	return AppConfig{
		ClipboardWatch:              watch,
		DownloadOnCopy:              dldOnCopy,
		ConcurrentDownloads:         cDownloads,
		ConcurrentPlaylistDownloads: cPlaylistDownloads,
		BaseSaveDir:                 baseSaveDir,
	}
}
