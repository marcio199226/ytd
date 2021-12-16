package plugins

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	ytDownloader "github.com/kkdai/youtube/v2"
	"github.com/wailsapp/wails/v2"

	. "ytd/constants"
	. "ytd/db"
	. "ytd/models"
)

type Yt struct {
	Name         string
	WailsRuntime *wails.Runtime
	AppConfig    *AppConfig
	AppStats     *AppStats
	NewEntries   chan GenericEntry
	dir          string
	client       ytDownloader.Client
	ctx          context.Context
}

type YtEntry struct {
	Type     string
	Track    GenericTrack
	Playlist GenericPlaylist
}

func (yt *Yt) SetWailsRuntime(runtime *wails.Runtime) {
	yt.WailsRuntime = runtime
}

func (yt *Yt) SetContext(ctx context.Context) {
	yt.ctx = ctx
}

func (yt *Yt) SetAppConfig(config *AppConfig) {
	yt.AppConfig = config
}

func (yt *Yt) SetAppStats(stats *AppStats) {
	yt.AppStats = stats
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

func (yt *Yt) GetDir() string {
	return yt.dir
}

func (yt *Yt) IsTrackFileExists(track GenericTrack, fileType string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s.%s", yt.dir, track.ID, fileType)); os.IsNotExist(err) {
		return false
	}
	return true
}

func (yt *Yt) Fetch(url string, isFromClipboard bool) *GenericEntry {
	fmt.Printf("Fetching from yt %s...", url)

	if isPlaylist := strings.Contains(url, "playlist?"); isPlaylist {
		fmt.Println("Fetching playlist info...")
		ytEntry := &GenericEntry{Source: yt.Name, Type: "playlist"}
		playlist, err := yt.client.GetPlaylist(url)
		if err != nil {
			fmt.Printf("yt.client.GetPlaylist(url) error: %s", err)
			return nil
		}

		var playlistTracks []string
		playlistEntry := NewGenericPlaylist(playlist.ID, playlist.Title, len(playlist.Videos), playlistTracks)
		ytEntry.Playlist = playlistEntry
		ytEntry.Playlist.Url = url
		for _, v := range playlist.Videos {
			ytEntry.Playlist.TracksIds = append(ytEntry.Playlist.TracksIds, v.ID)
		}
		err = DbWriteEntry(ytEntry.Playlist.ID, ytEntry)
		if err != nil {
			fmt.Printf("Error while saving playlist: %s", err)
		}
		yt.WailsRuntime.Events.Emit("ytd:playlist", ytEntry) // notify frontend
		for k, v := range playlist.Videos {
			fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
			video, err := yt.client.VideoFromPlaylistEntry(v)
			if err != nil {
				fmt.Printf("yt.client.VideoFromPlaylistEntry(v) error: %s", err)
				continue
			}

			track := NewGenericTrack(video.ID, video.Title, video.Author, video.ID, url)
			for _, thumbnail := range video.Thumbnails {
				track.Thumbnails = append(track.Thumbnails, thumbnail.URL)
			}
			ytEntry := &GenericEntry{Source: yt.Name, Type: "track", Track: track}
			ytEntry.Track.Status = TrackStatusProcessing
			ytEntry.Track.Duration = video.Duration.Minutes()
			ytEntry.Track.PlaylistID = playlist.ID
			DbWriteEntry(ytEntry.Track.ID, ytEntry)           // write to rigth bucket
			yt.WailsRuntime.Events.Emit("ytd:track", ytEntry) // notify frontend

			err = yt.downloadTrack(video, ytEntry)
			if err != nil {
				fmt.Printf("yt.downloadTrack(video, video.ID) error: %s \n", err)
				ytEntry.Track.Status = TrackStatusFailed
				ytEntry.Track.StatusError = err.Error()
				DbWriteEntry(ytEntry.Track.ID, ytEntry)
				yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
				continue
			}
			ytEntry.Track.Status = TrackStatusDownladed

			DbWriteEntry(ytEntry.Track.ID, ytEntry)
			yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		}
		return ytEntry
	}

	return yt.fetchTrack(url, isFromClipboard)
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

func (yt *Yt) fetchTrack(url string, isFromClipboard bool) *GenericEntry {
	video, err := yt.client.GetVideoContext(yt.ctx, url)
	if err != nil {
		fmt.Printf("yt.client.GetVideo(url) error: %s \n", err)
		ytEntry := GenericEntry{Source: yt.Name, Type: "track", Track: NewFailedTrack(url, err)}
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return nil
	}

	ytEntry := yt.addTrack(url, video)
	DbWriteEntry(ytEntry.Track.ID, ytEntry)           // write to rigth bucket
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry) // notify frontend

	if yt.AppConfig.ClipboardWatch && !yt.AppConfig.DownloadOnCopy && isFromClipboard {
		return nil
	}

	if yt.AppStats.DownloadingCount >= yt.AppConfig.MaxParrallelDownloads {
		fmt.Println("MaxParrallelDownloads has been reached")
		return nil
	}

	yt.AppStats.IncDndCount()
	err = yt.downloadTrack(video, ytEntry)
	yt.AppStats.DecDndCount()
	if err != nil {
		ytEntry.Track.Status = TrackStatusFailed
		ytEntry.Track.StatusError = err.Error()
		DbWriteEntry(ytEntry.Track.ID, ytEntry)
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return nil
	}
	ytEntry.Track.Status = TrackStatusDownladed
	DbWriteEntry(ytEntry.Track.ID, ytEntry)
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
	return ytEntry
}

func (yt *Yt) StartDownload(ytEntry *GenericEntry) GenericEntry {
	// update AppState.Entries
	video, err := yt.client.GetVideo(ytEntry.Track.Url)
	if err != nil {
		fmt.Printf("StartDownload error: %s \n", err)
		ytEntry.Track.Status = TrackStatusFailed
		ytEntry.Track.StatusError = err.Error()
		DbWriteEntry(ytEntry.Track.ID, ytEntry)
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return *ytEntry
	}

	yt.AppStats.IncDndCount()
	err = yt.downloadTrack(video, ytEntry)
	yt.AppStats.DecDndCount()
	if err != nil {
		ytEntry.Track.Status = TrackStatusFailed
		ytEntry.Track.StatusError = err.Error()
		DbWriteEntry(ytEntry.Track.ID, ytEntry)
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return *ytEntry
	}
	ytEntry.Track.Status = TrackStatusDownladed
	DbWriteEntry(ytEntry.Track.ID, ytEntry)
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
	return *ytEntry
}

func (yt *Yt) downloadTrack(video *ytDownloader.Video, ytEntry *GenericEntry) error {
	audioFormats := video.Formats.Type("audio/webm")

	if len(audioFormats) == 0 {
		return TrackWithoutAudioFormat
	}

	stream, size, err := yt.client.GetStream(video, &audioFormats[0])
	ytEntry.Track.FileSize = size
	ytEntry.Track.Status = TrackStatusDownloading
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
	if err != nil {
		return err
	}
	return yt.saveTrack(stream, video.ID, size)
}

func (yt *Yt) saveTrack(stream io.ReadCloser, filename string, filesize int64) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.webm", yt.dir, filename))
	if err != nil {
		fmt.Printf("Cannot create file for track %s | %s\n", filename, err)
		return err
	}
	defer file.Close()

	done := make(chan int64)
	progress := make(chan uint)
	go downloadProgress(done, progress, fmt.Sprintf("%s/%s.webm", yt.dir, filename), filesize)
	go func() { // read from progress channel as soon as downloadProgress writes to progress chan
		// loop until channel progress is closed by downloadProgress goroutine
		for p := range progress {
			fmt.Printf("Received progress: %d for video with ID %s\n", p, filename)
			yt.WailsRuntime.Events.Emit(
				"ytd:track:progress",
				struct {
					Id       string `json:"id"`
					Progress uint   `json:"progress"`
				}{Id: filename, Progress: p},
			)
		}
	}()

	w, err := io.Copy(file, stream)
	if err != nil {
		return err
	}

	done <- w

	return nil
}

func (yt *Yt) addTrack(url string, video *ytDownloader.Video) *GenericEntry {
	if video == nil {
		track := GenericTrack{Url: url, Status: TrackStatusPending}
		return &GenericEntry{Source: yt.Name, Type: "track", Track: track}
	}
	track := NewGenericTrack(video.ID, video.Title, video.Author, video.ID, url)
	for _, thumbnail := range video.Thumbnails {
		track.Thumbnails = append(track.Thumbnails, thumbnail.URL)
	}
	ytEntry := &GenericEntry{Source: yt.Name, Type: "track", Track: track}
	ytEntry.Track.Status = TrackStatusProcessing
	ytEntry.Track.Duration = video.Duration.Minutes()
	return ytEntry
}

func downloadProgress(done chan int64, progress chan<- uint, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}

			fi, err := file.Stat()
			if err != nil {
				fmt.Println(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100
			progress <- uint(percent)
		}

		if stop {
			close(progress)
			break
		}

		time.Sleep(time.Second)
	}
}
