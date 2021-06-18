package models

import (
	. "ytd/constants"
)

type GenericTrack struct {
	ID               string   `json:"id"`
	PlaylistID       string   `json:"playlistId"`
	Name             string   `json:"name"`
	Duration         float64  `json:"duration"`
	Author           string   `json:"author"`
	Thumbnails       []string `json:"thumbnails"`
	Downloaded       bool     `json:"downloaded"`
	DownloadProgress uint8    `json:"downloadProgress"`
	Status           string   `json:"status"`
	StatusError      string   `json:"statusError"`
	FileSize         int64    `json:"filesize"` // bytes
	filename         string
	Url              string `json:"url"`
}

func NewGenericTrack(id string, name string, author string, filename string, url string) GenericTrack {
	return GenericTrack{
		ID:         id,
		Name:       name,
		Author:     author,
		Downloaded: false,
		Status:     TrackStatusPending,
		filename:   filename,
		Url:        url,
	}
}

func NewFailedTrack(url string, err error) GenericTrack {
	return GenericTrack{
		Status:      TrackStatusFailed,
		StatusError: err.Error(),
		Url:         url,
	}
}

func (s GenericTrack) isEmpty() bool {
	return s.ID == ""
}

func (s GenericTrack) isFromPlaylist() bool {
	return s.PlaylistID != ""
}
