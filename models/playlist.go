package models

type GenericPlaylist struct {
	ID          string
	TracksCount int
	Name        string
	Tracks      []GenericTrack
}

func NewGenericPlaylist(id string, name string, playlistSize int, tracks []GenericTrack) GenericPlaylist {
	return GenericPlaylist{
		ID:          id,
		Name:        name,
		TracksCount: playlistSize,
		Tracks:      tracks,
	}
}

func (s GenericPlaylist) isEmpty() bool {
	return s.ID == ""
}
