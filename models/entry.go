package models

type GenericEntry struct {
	Type     string          `json:"type"`
	Source   string          `json:"source"`
	Track    GenericTrack    `json:"track"`
	Playlist GenericPlaylist `json:"playlist"`
}

func NewGenericEntry(source string) *GenericEntry {
	return &GenericEntry{Source: source}
}
