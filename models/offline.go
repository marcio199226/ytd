package models

import guuid "github.com/google/uuid"

type OfflinePlaylist struct {
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	TracksIds []string `json:"tracksIds"`
}

func NewOfflinePlaylist(name string, tracks []string) OfflinePlaylist {
	id := guuid.NewString()
	return OfflinePlaylist{
		UUID:      id,
		Name:      name,
		TracksIds: tracks,
	}
}
