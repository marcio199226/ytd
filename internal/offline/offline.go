package offline

import (
	"fmt"
	"io"
	"os"
	db "ytd/db"
	. "ytd/models"

	"github.com/leonelquinteros/gotext"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2"
)

type OfflinePlaylistService struct {
	Runtime *wails.Runtime
	config  AppConfig
}

func (p *OfflinePlaylistService) SetConfig(config AppConfig) {
	p.config = config
}

func (p *OfflinePlaylistService) CreateNewPlaylist(name string) (OfflinePlaylist, error) {
	playlist := NewOfflinePlaylist(name, []string{})
	err := db.DbAddOfflinePlaylist(playlist.UUID, playlist, false)
	if err != nil {
		return OfflinePlaylist{}, err
	}
	return playlist, nil
}

func (p *OfflinePlaylistService) CreateNewPlaylistWithTracks(name string, tracks []string) (OfflinePlaylist, error) {
	playlist := NewOfflinePlaylist(name, tracks)
	err := db.DbAddOfflinePlaylist(playlist.UUID, playlist, false)
	if err != nil {
		return OfflinePlaylist{}, err
	}
	p.Runtime.Events.Emit("ytd:offline:playlists:created")
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
	playlist.TracksIds = append(playlist.TracksIds[:idx], playlist.TracksIds[idx+1:]...)
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

func (p *OfflinePlaylistService) ExportPlaylist(uuid string, path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}

	playlist := p.GetPlaylistByUUID(uuid)
	copied := 0
	for idx, id := range playlist.TracksIds {
		err := copyFile(fmt.Sprintf("%s/%s/%s.mp3", p.config.BaseSaveDir, "youtube", id), fmt.Sprintf("%s/%s.mp3", path, id))
		copied++
		ShowLoader(p.Runtime, fmt.Sprintf(
			gotext.Get("Exporting...%d/%d", idx+1, len(playlist.TracksIds))),
		)
		if err != nil {
			fmt.Printf("Error while copying track: %s \n\n", err)
			copied--
		}
	}

	return copied == len(playlist.TracksIds), nil
}

func (p *OfflinePlaylistService) GetPlaylists(emitEvent bool) ([]OfflinePlaylist, error) {
	playlists := db.DbGetAllOfflinePlaylists()
	if emitEvent {
		p.Runtime.Events.Emit("ytd:offline:playlists", playlists)
	}
	return playlists, nil
}

func (p *OfflinePlaylistService) GetPlaylistByUUID(uuid string) OfflinePlaylist {
	playlists, _ := p.GetPlaylists(false)
	for _, playlist := range playlists {
		if playlist.UUID == uuid {
			return playlist
		}
	}
	return OfflinePlaylist{}
}

func copyFile(src string, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return errors.Wrap(err, "copyFile os.Open(src)")
	}
	defer srcfd.Close()

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		// dst does not exist so create it
		if dstfd, err = os.Create(dst); err != nil {
			return errors.Wrap(err, "copyFile os.Create(dst)")
		}
		defer dstfd.Close()
	} else {
		if dstfd, err = os.Open(dst); err != nil {
			return errors.Wrap(err, "copyFile os.Open(dst)")
		}
		defer dstfd.Close()
	}

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return errors.Wrap(err, "copyFile io.Copy(dstfd, srcfd)")
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return errors.Wrap(err, "copyFile os.Stat(src)")
	}
	if err = os.Chmod(dst, srcinfo.Mode()); err != nil {
		return errors.Wrap(err, "copyFile os.Chmod(dst, srcinfo.Mode())")
	}
	return nil
}
