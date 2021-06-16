package plugins

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	ytDownloader "github.com/kkdai/youtube/v2"
	"github.com/wailsapp/wails"

	. "ytd/db"
	. "ytd/models"
)

type Yt struct {
	Name         string
	WailsRuntime *wails.Runtime
	dir          string
	client       ytDownloader.Client
}

type YtEntry struct {
	Type     string
	Track    GenericTrack
	Playlist GenericPlaylist
}

func (yt *Yt) SetWailsRuntime(runtime *wails.Runtime) {
	yt.WailsRuntime = runtime
}

func (yt *Yt) GetName() string {
	return "youtube"
}

func (yt *Yt) Initialize() error {
	fmt.Println("Initializing yt client...")
	yt.client = ytDownloader.Client{Debug: true}
	return nil
}

func (yt *Yt) SetDir(dir string) {
	yt.dir = dir
	if _, err := os.Stat(yt.dir); os.IsNotExist(err) {
		fmt.Printf("Creating %s directory for youtube plugin", yt.dir)
		os.MkdirAll(yt.dir, os.ModePerm)
	}
}

func (yt *Yt) Fetch(url string) (*GenericEntry, error) {
	fmt.Println("Fetching from yt...")
	// time.Sleep(60 * time.Second)

	if isPlaylist := strings.Contains(url, "playlist?"); isPlaylist {
		fmt.Println("Fetching playlist info...")
		ytEntry := &GenericEntry{Source: yt.Name, Type: "playlist"}
		playlist, err := yt.client.GetPlaylist(url)
		if err != nil {
			return &GenericEntry{}, err
		}

		var playlistTracks []GenericTrack
		playlistEntry := NewGenericPlaylist(playlist.ID, playlist.Title, len(playlist.Videos), playlistTracks)
		ytEntry.Playlist = playlistEntry
		for k, v := range playlist.Videos {
			fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
			video, err := yt.client.VideoFromPlaylistEntry(v)
			if err != nil {
				return &GenericEntry{}, err
			}

			track := NewGenericTrack(video.ID, video.Title, video.Author, video.ID, url)
			playlistTracks = append(playlistTracks, track)

			err = yt.downloadTrack(video, fmt.Sprintf("%s_%s", playlist.ID, video.ID))
			if err != nil {
				return &GenericEntry{}, err
			}
			track.Downloaded = true
		}
		return &GenericEntry{}, nil
	}

	return yt.fetchTrack(url)
}

func (yt *Yt) GetFilename() error {
	return nil
}

func (yt *Yt) Supports(address string) bool {
	u, err := url.Parse(address)
	if err != nil {
		return false
	}
	return strings.Contains(u.Hostname(), "youtube")
}

func (yt *Yt) fetchTrack(url string) (*GenericEntry, error) {
	video, err := yt.client.GetVideo(url)
	if err != nil {
		return &GenericEntry{}, err
	}
	track := NewGenericTrack(video.ID, video.Title, video.Author, video.ID, url)
	for _, thumbnail := range video.Thumbnails {
		track.Thumbnails = append(track.Thumbnails, thumbnail.URL)
	}
	ytEntry := &GenericEntry{Source: yt.Name, Type: "track", Track: track}
	DbWriteEntry(track.ID, ytEntry)
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)

	err = yt.downloadTrack(video, video.ID)
	if err != nil {
		return &GenericEntry{}, err
	}
	ytEntry.Track.Downloaded = true

	DbWriteEntry(track.ID, ytEntry)
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
	return ytEntry, nil
}

func (yt *Yt) downloadTrack(video *ytDownloader.Video, filename string) error {
	for _, format := range video.Formats.Type("audio/webm") {
		fmt.Printf("%d | %d | %s | %s | %s | %s | %s \n", format.ItagNo, format.AudioChannels, format.AudioQuality, format.AudioSampleRate, format.Quality, format.QualityLabel, format.MimeType)
	}
	audioFormats := video.Formats.Type("audio/webm")
	stream, _, err := yt.client.GetStream(video, &audioFormats[0])
	if err != nil {
		return err
	}
	return yt.saveTrack(stream, filename)
}

func (yt *Yt) saveTrack(stream io.ReadCloser, filename string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.webm", yt.dir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}
	return nil
}
