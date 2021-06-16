package models

import "time"

type GenericTrack struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Duration   time.Duration `json:"duration"`
	Author     string        `json:"author"`
	Thumbnails []string      `json:"thumbnails"`
	Downloaded bool          `json:"downloaded"`
	filename   string
	url        string
}

func NewGenericTrack(id string, name string, author string, filename string, url string) GenericTrack {
	return GenericTrack{
		ID:         id,
		Name:       name,
		Author:     author,
		Downloaded: false,
		filename:   filename,
		url:        url,
	}
}

func (s GenericTrack) isEmpty() bool {
	return s.ID == ""
}
