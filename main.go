package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"sync"
	"time"

	. "ytd/clipboard"
	. "ytd/constants"
	. "ytd/db"
	. "ytd/models"
	. "ytd/plugins"

	_ "embed"

	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails"
	"github.com/xujiajun/nutsdb"
)

var wailsRuntime *wails.Runtime
var plugins []Plugin = []Plugin{&Yt{Name: "youtube"}}

//go:embed frontend/dist/main.js
var js string

//go:embed frontend/dist/styles.css
var css string

var appState *AppState
var newEntries = make(chan GenericEntry)

type AppState struct {
	log     *wails.CustomLogger
	runtime *wails.Runtime
	db      *nutsdb.DB
	plugins []Plugin
	Entries []GenericEntry `json:"entries"`
	Config  *AppConfig     `json:"config"`
	Stats   *AppStats
}

func (state *AppState) WailsInit(runtime *wails.Runtime) error {
	// Save runtime
	state.runtime = runtime
	state.log = runtime.Log.New("AppState")
	// Do some other initialisation

	state.db = InitializeDb()
	state.Entries = DbGetAllEntries()
	state.Config = state.Config.Init()
	state.Stats = &AppStats{}
	appState = state

	// this is sync so it blocks until finished and wails:loaded are not dispatched until this finishes
	runtime.Events.On("wails:loaded", func(...interface{}) {
		// entries := DbGetAllEntries()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("EMIT YTD:ONLOAD")
		runtime.Events.Emit("ytd:onload", state)
	})

	for _, plugin := range plugins {
		plugin.SetWailsRuntime(runtime)
		plugin.SetAppConfig(state.Config)
		plugin.SetAppStats(state.Stats)
	}
	fmt.Println("APP STATE INITIALIZED")

	/* 	go func() {
		for {
			time.Sleep(10 * time.Second)
			state.checkForTracksToDownload()
		}
	}() */

	go func() {
		restart := make(chan int)
		time.Sleep(3 * time.Second)
		go state.convertToMp3(restart)

		for {
			select {
			case <-restart:
				fmt.Println("Restart converting...")
				go state.convertToMp3(restart)
			}

		}
	}()

	return nil
}

func (state *AppState) GetAppConfig() *AppConfig {
	return state.Config
}

func (state *AppState) SelectDirectory() string {
	selectedDirectory := state.runtime.Dialog.SelectDirectory()
	return selectedDirectory
}

func (state *AppState) GetEntryById(entry GenericEntry) *GenericEntry {
	for _, t := range state.Entries {
		if t.Type == "track" && t.Track.ID == entry.Track.ID {
			return &t
		}
	}
	return nil
}

func (state *AppState) checkForTracksToDownload() error {
	fmt.Printf("Check for tracks to start downloads...%d/%d\n", state.Stats.DownloadingCount, state.Config.MaxParrallelDownloads)
	if state.Stats.DownloadingCount >= state.Config.MaxParrallelDownloads {
		return nil
	}

	fmt.Printf("Checking...%d entries\n", len(state.Entries))
	// range over DbGetAllEntries()
	for _, t := range DbGetAllEntries() {
		var freeSlots uint = state.Config.MaxParrallelDownloads - state.Stats.DownloadingCount
		entry := t
		// auto download only tracks with processing status
		// if track has pending/failed status it means that something goes wrong so user have to download it manually from UI
		if entry.Type == "track" && entry.Track.Status == TrackStatusProcessing {
			if freeSlots == 0 {
				return nil
			}

			// start download for track
			fmt.Printf("Found %s to download\n", entry.Track.Name)
			plugin := getPluginFor(entry.Source)
			// make chan GenericEntry
			// goriutine scrive li dentro
			if plugin != nil {
				freeSlots--
				go func(entry GenericEntry) {
					// storedEntry := state.GetEntryById(entry)
					// fmt.Printf("Stored entry %v\n", storedEntry)
					// fmt.Println(entry.Track.Url)
					// fmt.Println(storedEntry.Track.Url)
					plugin.StartDownload(&entry)
					// storedEntry.Track = entry.Track
				}(entry)
			}
		}
	}

	// qui un for che legge dal channel e ogni volta che riceve una entry downloadata
	return nil
}

func (state *AppState) convertToMp3(restart chan<- int) error {
	if !state.Config.ConvertToMp3 {
		// if option is not enabled restart check after 30s
		time.Sleep(30 * time.Second)
		restart <- 1
		return nil
	}

	fmt.Println("Converting....")
	ffmpeg, _ := isFFmpegInstalled()

	for _, t := range DbGetAllEntries() {
		entry := t
		plugin := getPluginFor(entry.Source)

		if entry.Type == "track" && entry.Track.Status == TrackStatusDownladed && !entry.Track.IsConvertedToMp3 && plugin.IsTrackFileExists(entry.Track, "webm") {
			fmt.Printf("Extracting audio for %s...\n", entry.Track.Name)

			// ffmpeg -i "41qC3w3UUkU.webm" -vn -ab 128k -ar 44100 -y "41qC3w3UUkU.mp3"
			cmd := exec.Command(
				ffmpeg,
				"-loglevel", "quiet",
				"-i", fmt.Sprintf("%s/%s.webm", plugin.GetDir(), entry.Track.ID),
				"-vn",
				"-ab", "128k",
				"-ar", "44100",
				"-y", fmt.Sprintf("%s/%s.mp3", plugin.GetDir(), entry.Track.ID),
			)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Println("Failed to extract audio:", err)
			} else {
				entry.Track.IsConvertedToMp3 = true
				DbWriteEntry(entry.Track.ID, entry)
				state.runtime.Events.Emit("ytd:track", entry) // track:converted:mp3

				// remove webm
				if state.Config.CleanWebmFiles && plugin.IsTrackFileExists(entry.Track, "webm") {
					err = os.Remove(fmt.Sprintf("%s/%s.webm", plugin.GetDir(), entry.Track.ID))
					if err != nil && !os.IsNotExist(err) {
						fmt.Printf("Cannot remove %s.webm file after successfull converting to mp3\n", entry.Track.ID)
					}
				}
			}
			restart <- 1
			return nil
		}
	}

	// if there are no tracks to convert delay between restart
	time.Sleep(30 * time.Second)
	restart <- 1
	return nil
}

func saveSettingBoolValue(name string, val bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovering from panic saveSettingValue:", r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

	error := DbSaveSettingBoolValue(name, val)
	if err != nil {
		return error
	}

	appState.Config.Set(name, val)
	return nil
}

func saveSettingValue(name string, val string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovering from panic saveSettingValue:", r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

	fmt.Println(val)
	fmt.Printf("%T", val)
	fmt.Println()
	error := DbWriteSetting(name, val)
	if err != nil {
		fmt.Println("ERRROR IN WRITE SETTINT")
		return error
	}

	appState.Config.Set(name, val)
	return nil
}

func readSettingBoolValue(name string) (bool, error) {
	return DbReadSettingBoolValue(name)
}

func readSettingValue(name string) (string, error) {
	return DbReadSetting(name)
}

func removeEntry(record map[string]interface{}) error {
	var err error
	var entry GenericEntry
	err = mapstructure.Decode(record, &entry)
	if err != nil {
		return err
	}

	if entry.Type == "track" {
		if err = DbDeleteEntry(entry.Track.ID); err == nil {
			plugin := getPluginFor(entry.Source)
			if plugin.IsTrackFileExists(entry.Track, "webm") {
				err = os.Remove(fmt.Sprintf("%s/%s.webm", plugin.GetDir(), entry.Track.ID))
				fmt.Println(err)
				if err != nil && !os.IsNotExist(err) {
					return err
				}
			}
			// remove mp3 if file has been already converted
			if plugin.IsTrackFileExists(entry.Track, "mp3") {
				err = os.Remove(fmt.Sprintf("%s/%s.mp3", plugin.GetDir(), entry.Track.ID))
				fmt.Println(err)
				if err != nil && !os.IsNotExist(err) {
					return err
				}
			}
		}
	}
	return err
}

func addToDownload(url string, isFromClipboard bool) error {
	for _, plugin := range plugins {
		if support := plugin.Supports(url); support {
			if appState.GetAppConfig().ConcurrentDownloads {
				go func() {
					plugin.Fetch(url, isFromClipboard)
				}()
			} else {
				plugin.Fetch(url, isFromClipboard)
			}
			continue
		}
	}
	return nil
}

func startDownload(record map[string]interface{}) error {
	var err error
	var entry GenericEntry
	err = mapstructure.Decode(record, &entry)
	if err != nil {
		return err
	}

	if appState.Stats.DownloadingCount >= appState.Config.MaxParrallelDownloads {
		return errors.New("Max simultaneous downloads are reached please retry after some track finished downloading")
	}

	for _, plugin := range plugins {
		if plugin.GetName() == entry.Source && entry.Type == "track" {
			if appState.GetAppConfig().ConcurrentDownloads {
				go func() {
					plugin.StartDownload(&entry)
				}()
			} else {
				plugin.StartDownload(&entry)
			}
			continue
		}
	}
	return nil
}

func isSupportedUrl(url string) bool {
	for _, plugin := range plugins {
		if support := plugin.Supports(url); support {
			return true
		}
	}
	return false
}

func isFFmpegInstalled() (string, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	return ffmpeg, err
}

func getPluginFor(name string) Plugin {
	for _, plugin := range plugins {
		if plugin.GetName() == name {
			return plugin
		}
	}
	return nil
}

//WailsRuntime .
type WailsRuntime struct {
	runtime *wails.Runtime
}

func (s *WailsRuntime) WailsShutdown() {
	CloseDb()
}

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:            1024,
		Height:           768,
		Title:            "ytd",
		JS:               js,
		CSS:              css,
		Colour:           "#131313",
		DisableInspector: false,
	})
	app.Bind(&AppState{})
	app.Bind(saveSettingBoolValue)
	app.Bind(saveSettingValue)
	app.Bind(readSettingBoolValue)
	app.Bind(readSettingValue)
	app.Bind(removeEntry)
	app.Bind(addToDownload)
	app.Bind(startDownload)
	app.Bind(isSupportedUrl)
	app.Bind(isFFmpegInstalled)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancelCtx := context.WithCancel(context.Background())
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill)
	changes := make(chan string, 10)
	stopCh := make(chan struct{}, 1)
	go MonitorClipboard(time.Second, ctx, wg, stopCh, changes)

	go func(wg *sync.WaitGroup) {
		sig := <-exitChan
		fmt.Println("Received signal", sig)
		// send the shutdown signal through the context.Context
		cancelCtx()
		wg.Done()
		stopCh <- struct{}{}
		return
	}(wg)

	// initialize plugins
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/songs", usr.HomeDir)
	for _, plugin := range plugins {
		log.Println("download to dir=", fmt.Sprintf("%s/%s", currentDir, plugin.GetName()))
		plugin.SetDir(fmt.Sprintf("%s/%s", currentDir, plugin.GetName()))
		plugin.Initialize()
	}

	go func() {
		// Watch for changes of local clipboard
		for {
			select {
			case <-stopCh:
				fmt.Printf("stopped manually")
				time.Sleep(time.Second)
				wg.Wait()
				os.Exit(0)
				return
			default:
				change, ok := <-changes /*  */
				if ok && change != "" {
					log.Printf("change received: '%s'", change)
					if appState.Config.ClipboardWatch {
						addToDownload(change, true)
					}
				} else {
					log.Println("channel has been closed. exiting...")
					time.Sleep(time.Millisecond)
				}
			}
		}
	}()

	go func() {
		http.Handle("/", http.FileServer(http.Dir(currentDir)))
		http.ListenAndServe(":8080", nil)
	}()

	app.Run()
}
