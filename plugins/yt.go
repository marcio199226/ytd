package plugins

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	ytDownloader "github.com/kkdai/youtube/v2"
	"github.com/wailsapp/wails"

	. "ytd/constants"
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

func (yt *Yt) Fetch(url string) {
	fmt.Println("Fetching from yt...")
	// time.Sleep(60 * time.Second)

	if isPlaylist := strings.Contains(url, "playlist?"); isPlaylist {
		fmt.Println("Fetching playlist info...")
		ytEntry := &GenericEntry{Source: yt.Name, Type: "playlist"}
		playlist, err := yt.client.GetPlaylist(url)
		if err != nil {
			fmt.Printf("yt.client.GetPlaylist(url) error: %s", err)
			return
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
			/* 			// for testing purpose we limit to max 3 tracks to be downloaded from a playlist
			   			if k > 3 {
			   				break
			   			} */
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
			ytEntry.Track.Downloaded = true
			ytEntry.Track.Status = TrackStatusDownladed

			DbWriteEntry(ytEntry.Track.ID, ytEntry)
			yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		}
		return
	}

	yt.fetchTrack(url)
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

func (yt *Yt) fetchTrack(url string) {
	video, err := yt.client.GetVideo(url)
	if err != nil {
		fmt.Printf("yt.client.GetVideo(url) error: %s \n", err)
		ytEntry := GenericEntry{Source: yt.Name, Type: "track", Track: NewFailedTrack(url, err)}
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return
	}

	ytEntry := yt.addTrack(url, video)
	ytEntry.Track.Status = TrackStatusProcessing
	ytEntry.Track.Duration = video.Duration.Minutes()
	DbWriteEntry(ytEntry.Track.ID, ytEntry)           // write to rigth bucket
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry) // notify frontend

	err = yt.downloadTrack(video, ytEntry)
	if err != nil {
		fmt.Printf("yt.downloadTrack(video, video.ID) error: %s \n", err)
		ytEntry.Track.Status = TrackStatusFailed
		ytEntry.Track.StatusError = err.Error()
		DbWriteEntry(ytEntry.Track.ID, ytEntry)
		yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
		return
	}
	ytEntry.Track.Downloaded = true
	ytEntry.Track.Status = TrackStatusDownladed

	DbWriteEntry(ytEntry.Track.ID, ytEntry)
	yt.WailsRuntime.Events.Emit("ytd:track", ytEntry)
}

func (yt *Yt) downloadTrack(video *ytDownloader.Video, ytEntry *GenericEntry) error {
	for _, format := range video.Formats.Type("audio/webm") {
		fmt.Printf("%d | %d | %s | %s | %s | %s | %s \n", format.ItagNo, format.AudioChannels, format.AudioQuality, format.AudioSampleRate, format.Quality, format.QualityLabel, format.MimeType)
	}
	audioFormats := video.Formats.Type("audio/webm")

	if len(audioFormats) == 0 {
		return TrackWithoutAudioFormat
	}

	stream, size, err := yt.client.GetStream(video, &audioFormats[0])
	ytEntry.Track.FileSize = size
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
		track := GenericTrack{Url: url, Status: TrackStatusProcessing}
		return &GenericEntry{Source: yt.Name, Type: "track", Track: track}
	}
	track := NewGenericTrack(video.ID, video.Title, video.Author, video.ID, url)
	for _, thumbnail := range video.Thumbnails {
		track.Thumbnails = append(track.Thumbnails, thumbnail.URL)
	}
	ytEntry := &GenericEntry{Source: yt.Name, Type: "track", Track: track}
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
