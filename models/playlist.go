package models

type GenericPlaylist struct {
	ID          string   `json:"id"`
	Url         string   `json:"url"`
	TracksCount int      `json:"count"`
	Name        string   `json:"name"`
	Thumbnails  string   `json:"thumbnail"`
	TracksIds   []string `json:"tracksIds"`
}

func NewGenericPlaylist(id string, name string, playlistSize int, tracks []string) GenericPlaylist {
	return GenericPlaylist{
		ID:          id,
		Name:        name,
		TracksCount: playlistSize,
		TracksIds:   tracks,
	}
}

func (s GenericPlaylist) isEmpty() bool {
	return s.ID == ""
}
