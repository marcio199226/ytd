package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

//go:embed frontend/dist/index.html
var html string

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
		fmt.Println("EMIT YTD:ONLOAD")
		runtime.Events.Emit("ytd:onload", state)
	})

	for _, plugin := range plugins {
		plugin.SetWailsRuntime(runtime)
		plugin.SetAppConfig(state.Config)
		plugin.SetAppStats(state.Stats)
	}
	fmt.Println("APP STATE INITIALIZED")

	go func() {
		for {
			time.Sleep(10 * time.Second)
			state.checkForTracksToDownload()
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
	for _, t := range state.Entries {
		var freeSlots uint = state.Config.MaxParrallelDownloads - state.Stats.DownloadingCount
		// auto download only tracks with processing status
		// if track has pending/failed status it means that something goes wrong so user have to download it manually from UI
		if t.Type == "track" && t.Track.Status == TrackStatusProcessing {
			if freeSlots == 0 {
				return nil
			}

			// start download for track
			fmt.Printf("Found %s to download\n", t.Track.Name)
			plugin := getPluginFor(t.Source)
			// make chan GenericEntry
			// goriutine scrive li dentro
			if plugin != nil {
				/* 				go func() {
					storedEntry := state.GetEntryById(t)
					fmt.Printf("Stored entry %v\n", storedEntry)
					entry := plugin.StartDownload(&t)
					storedEntry.Track = entry.Track
				}() */
			}
			freeSlots--
		}
	}

	// qui un for che legge dal channel e ogni volta che riceve una entry downloadata
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

	error := DbWriteSetting(name, val)
	if err != nil {
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
			err = os.Remove(fmt.Sprintf("%s/%s/%s.webm", appState.Config.BaseSaveDir, entry.Source, entry.Track.ID))
			if os.IsNotExist(err) {
				return nil
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
		HTML:             html,
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
