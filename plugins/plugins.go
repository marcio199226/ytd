package plugins

import (
	"context"

	"github.com/wailsapp/wails/v2"

	offline "ytd/internal/offline"
	. "ytd/models"
)

type Plugin interface {
	GetName() string
	Initialize() error
	GetDir() string
	SetDir(dir string)
	IsTrackFileExists(track GenericTrack, fileType string) bool
	Fetch(url string, isFromClipboard bool) *GenericEntry
	StartDownload(ytEntry *GenericEntry) GenericEntry
	GetFilename() error
	Supports(address string) bool
	SetWailsRuntime(*wails.Runtime)
	SetContext(context.Context)
	SetAppConfig(config *AppConfig)
	SetQueue(queue chan GenericEntry) error
	SetOfflineService(service *offline.OfflinePlaylistService)
	SetAppStats(stats *AppStats)
}
