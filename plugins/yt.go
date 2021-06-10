package plugins

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	ytDownloader "github.com/kkdai/youtube/v2"
)

type Yt struct {
	Name string
	dir  string
}

func (yt *Yt) GetName() string {
	return "youtube"
}

func (yt *Yt) Initialize() error {
	fmt.Println("Initializing from yt...")
	return nil
}

func (yt *Yt) SetDir(dir string) {
	yt.dir = dir
	if _, err := os.Stat(yt.dir); os.IsNotExist(err) {
		fmt.Printf("Creating %s directory for youtube plugin", yt.dir)
		os.MkdirAll(yt.dir, os.ModePerm)
	}
}

func (yt *Yt) Fetch(url string) error {
	fmt.Println("Fetching from yt...")

	y := ytDownloader.Client{Debug: true}
	/* 	video, err := y.GetVideo(url)
	   	if err != nil {
	   		return err
	   	} */

	playlist, err := y.GetPlaylist(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(playlist)
	return nil

	/* 	for _, format := range video.Formats.Type("audio/webm") {
	   		fmt.Printf("%d | %d | %s | %s | %s | %s | %s \n", format.ItagNo, format.AudioChannels, format.AudioQuality, format.AudioSampleRate, format.Quality, format.QualityLabel, format.MimeType)
	   	}
	   	audioFormats := video.Formats.Type("audio/webm")
	   	stream, _, err := y.GetSśtream(video, &audioFormats[0])
	   	if err != nil {
	   		return err
	   	}

	   	file, err := os.Create(fmt.Sprintf("%s/%s.webm", yt.dir, video.ID))
	   	if err != nil {
	   		return err
	   	}
	   	defer file.Close()

	   	_, err = io.Copy(file, stream)
	   	if err != nil {
	   		return err
	   	}
	   	return nil */
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
