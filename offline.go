package main

import (
	"fmt"
	db "ytd/db"
	. "ytd/models"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
)

type OfflinePlaylistService struct {
	runtime *wails.Runtime
}

func (p *OfflinePlaylistService) CreateNewPlaylist(name string) (OfflinePlaylist, error) {
	playlist := NewOfflinePlaylist(name, []string{})
	err := db.DbAddOfflinePlaylist(playlist.UUID, playlist, false)
	if err != nil {
		return OfflinePlaylist{}, err
	}
	return playlist, nil
}

func (p *OfflinePlaylistService) RemovePlaylist(uuid string) (bool, error) {
	err := db.DbRemoveOfflinePlaylist(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *OfflinePlaylistService) RemoveTrackFromPlaylist(tid string, playlist OfflinePlaylist) (OfflinePlaylist, error) {
	var idx int
	for k, id := range playlist.TracksIds {
		if id == tid {
			idx = k
			break
		}
	}
	fmt.Println("============PLAYLIST BEFORE REMOVE==============")
	fmt.Println(playlist)
	playlist.TracksIds = append(playlist.TracksIds[:idx], playlist.TracksIds[idx+1:]...)
	fmt.Println("============PLAYLIST==============")
	fmt.Println(playlist)
	err := db.DbAddOfflinePlaylist(playlist.UUID, playlist, true)
	if err != nil {
		return OfflinePlaylist{}, err
	}
	return playlist, nil
}

func (p *OfflinePlaylistService) AddTrackToPlaylist(payload []map[string]interface{}) (bool, error) {
	var err error
	var playlists []OfflinePlaylist
	err = mapstructure.Decode(payload, &playlists)
	if err != nil {
		return false, err
	}
	for _, p := range playlists {
		err = db.DbAddOfflinePlaylist(p.UUID, p, true)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (p *OfflinePlaylistService) ExportPlaylist(uuid string) (string, error) {
	selectedDirectory, err := p.runtime.Dialog.OpenDirectory(&dialog.OpenDialog{
		AllowFiles:           false,
		CanCreateDirectories: true,
		AllowDirectories:     true,
		Title:                "Choose directory",
	})

	if err != nil {
		return "", errors.Wrap(err, "OfflinePlaylistService ExportPlaylist()")
	}

	return selectedDirectory, nil
}

func (p *OfflinePlaylistService) GetPlaylists(emitEvent bool) ([]OfflinePlaylist, error) {
	playlists := db.DbGetAllOfflinePlaylists()
	if emitEvent {
		p.runtime.Events.Emit("ytd:offline:playlists", playlists)
	}
	return playlists, nil
}
