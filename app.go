package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "ytd/constants"
	. "ytd/db"
	. "ytd/models"
	. "ytd/plugins"

	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/mac"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
	"github.com/xujiajun/nutsdb"

	tb "gopkg.in/tucnak/telebot.v2"
)

var wailsRuntime *wails.Runtime

type AppState struct {
	runtime          *wails.Runtime
	db               *nutsdb.DB
	plugins          []Plugin
	Entries          []GenericEntry    `json:"entries"`
	OfflinePlaylists []OfflinePlaylist `json:"offlinePlaylists"`
	Config           *AppConfig        `json:"config"`
	Stats            *AppStats
	AppVersion       string `json:"appVersion"`

	isInForeground         bool
	canStartAtLogin        bool
	tray                   TrayMenu
	updater                *Updater
	offlinePlaylistService *OfflinePlaylistService
}

func (state *AppState) PreWailsInit() {
	state.db = InitializeDb()
	state.Entries = DbGetAllEntries()
	state.OfflinePlaylists = DbGetAllOfflinePlaylists()
	state.offlinePlaylistService = &OfflinePlaylistService{}
	state.Config = state.Config.Init()
	state.AppVersion = version

	if _, err := mac.StartsAtLogin(); err == nil {
		state.canStartAtLogin = true
	}
}

func (state *AppState) WailsInit(runtime *wails.Runtime) {
	// Save runtime
	state.runtime = runtime
	state.offlinePlaylistService.runtime = runtime
	// Do some other initialisation
	state.Stats = &AppStats{}
	appState = state

	// this is sync so it blocks until finished and wails:loaded are not dispatched until this finishes
	if runtime.System.AppType() == "default" { // wails serve & ng serve
		runtime.Events.On("wails:loaded", func(...interface{}) {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("EMIT YTD:ONLOAD")
			runtime.Events.Emit("ytd:onload", state)
		})
	} else { // dekstop build
		go func() {
			runtime.Events.Emit("ytd:onload", state)
		}()
	}

	// initialize plugins
	for _, plugin := range plugins {
		plugin.SetWailsRuntime(runtime)
		plugin.SetAppConfig(state.Config)
		plugin.SetAppStats(state.Stats)
	}

	fmt.Println("APP STATE INITIALIZED")
	state.InitializeListeners()
	// create app sys tray menu
	state.tray.runtime = runtime
	state.tray.createTray()
	state.runtime.Menu.SetTrayMenu(state.tray.defaultTrayMenu)
	// event emitted from fe if tray should be updated
	runtime.Events.On("ytd:app:tray:update", func(data ...interface{}) {
		state.tray.reRenderTray(func() {})
	})

	state.checkStartsAtLogin()

	/* 	go func() {
		for {
			time.Sleep(10 * time.Second)
			state.checkForTracksToDownload()
		}
	}() */

	/* 	go func() {
		for {
			restart := make(chan int)
			time.Sleep(3 * time.Second)
			go state.telegramShareTracks(restart)

			for {
				select {
				case <-restart:
					fmt.Println("Share tracks...")
					go state.telegramShareTracks(restart)
				}

			}
		}
	}() */

	go func() {
		restart := make(chan int)
		time.Sleep(3 * time.Second)
		go state.convertToMp3(restart)

		for {
			select {
			// @TODO
			// qui leggiamo da newEntries dove andiamo a scrivere o dentro yt.Fetch alla fine del metodo
			// oppure da dentro AddToDownload e facciamo tornare a plugin.fetch l'entry creata
			// cosi se legge prima da newEntries converte per prima quella altrimenti convertira le tracce che sono in attesa da prima
			/* 			case <-newEntries:
			fmt.Println("Converto new track...")
			go state.convertToMp3(restart, newEntries) */
			case <-restart:
				fmt.Println("Restart converting...")
				go state.convertToMp3(restart)
			}

		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		for {
			state.checkForUpdates()
			// check again in twelve hours
			time.Sleep(12 * time.Hour)
		}
	}()
}

func (state *AppState) WailsShutdown() {
	err := state.db.Merge()
	if err != nil {
		fmt.Println("WailsShutdown db.Merge() failed", err)
	}
	CloseDb()
}

func (state *AppState) InitializeListeners() {
	state.runtime.Events.On("ytd:app:foreground", func(optionalData ...interface{}) {
		var json map[string]interface{} = optionalData[0].(map[string]interface{})
		if isInForeground, ok := json["isInForeground"]; ok {
			state.isInForeground = isInForeground.(bool)
		}
	})
}

func (state *AppState) GetAppConfig() *AppConfig {
	return state.Config
}

func (state *AppState) SelectDirectory() (string, error) {
	selectedDirectory, err := state.runtime.Dialog.OpenDirectory(&dialog.OpenDialog{
		AllowFiles:           false,
		CanCreateDirectories: true,
		AllowDirectories:     true,
		Title:                "Choose directory",
	})
	return selectedDirectory, err
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

type recipientString string

func (r recipientString) Recipient() string {
	return r.Recipient()
}

func (state *AppState) telegramShareTracks(restart chan<- int) error {
	if !state.Config.Telegram.Share {
		// if option is not enabled restart check after 30s
		time.Sleep(30 * time.Second)
		restart <- 1
		return nil
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  "1903196088:AAHhWGvfhQfS_MlhvohFvQYnrg3z7GsBPOM",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	audio := &tb.Audio{File: tb.FromDisk("/Users/oskarmarciniak/songs/youtube/2vOU4nVI_DA.mp3"), FileName: "Track name", Title: "Track name title", Duration: 207, Caption: "Caption tracks"}
	b.Send(&tb.Chat{ID: 903612486}, audio)

	b.Handle("/hello", func(m *tb.Message) {
		fmt.Println(m.Chat)
		b.Send(m.Sender, "Hello World!")
		// cercare il primo msg che come username ha quello impostato nell'app e salvarsi la chat_id
	})

	b.Start()

	return nil
}

func (state *AppState) convertToMp3(restart chan<- int) error {
	if !state.Config.ConvertToMp3 {
		// if option is not enabled restart check after 30s
		time.Sleep(60 * time.Second)
		restart <- 1
		return nil
	}

	fmt.Println("Converting....")
	ffmpeg, _ := state.IsFFmpegInstalled()

	for _, t := range DbGetAllEntries() {
		entry := t
		plugin := getPluginFor(entry.Source)

		if entry.Type == "track" && entry.Track.Status == TrackStatusDownladed && !entry.Track.IsConvertedToMp3 && plugin.IsTrackFileExists(entry.Track, "webm") {
			fmt.Printf("Extracting audio for %s...\n", entry.Track.Name)
			entry.Track.ConvertingStatus.Status = TrakcConverting
			DbWriteEntry(entry.Track.ID, entry)
			state.runtime.Events.Emit("ytd:track", entry)

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
				entry.Track.ConvertingStatus.Status = TrakcConvertFailed
				entry.Track.ConvertingStatus.Err = err.Error()
				entry.Track.ConvertingStatus.Attempts += 1
				DbWriteEntry(entry.Track.ID, entry)
				state.runtime.Events.Emit("ytd:track", entry)

			} else {
				entry.Track.ConvertingStatus.Status = TrakcConverted
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
	time.Sleep(15 * time.Second)
	restart <- 1
	return nil
}

func (state *AppState) checkForUpdates() {
	if !state.Config.CheckForUpdates {
		return
	}

	state.updater = &Updater{
		CurrentVersion:              version,
		LatestReleaseGitHubEndpoint: "https://api.github.com/repos/marcio199226/ytd-binaries/releases",
		Client:                      &http.Client{Timeout: 10 * time.Minute},
		SelectAsset: func(release Release, asset Asset) bool {
			// look for the zip file
			return strings.Contains(asset.Name, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)) && filepath.Ext(asset.Name) == ".zip"
		},
		DownloadBytesLimit: 10_741_824, // 10MB
	}

	latest, hasUpdate, err := state.updater.HasUpdate()

	if err != nil {
		_, err := state.runtime.Dialog.Message(&dialog.MessageDialog{
			Type:         dialog.ErrorDialog,
			Title:        "Update check failed",
			Message:      err.Error(),
			Buttons:      []string{"OK"},
			CancelButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	if !hasUpdate {
		_, err := state.runtime.Dialog.Message(&dialog.MessageDialog{
			Type:         dialog.InfoDialog,
			Title:        "You're up to date",
			Message:      fmt.Sprintf("%s is the latest version.", latest.TagName),
			Buttons:      []string{"OK"},
			CancelButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	clickedAction, _ := state.runtime.Dialog.Message(&dialog.MessageDialog{
		Type:         dialog.InfoDialog,
		Title:        fmt.Sprintf("New version available: %s", latest.TagName),
		Message:      "Would you like to update?",
		Buttons:      []string{"OK", "Changelog", "Update"},
		CancelButton: "OK",
	})

	state.tray.reRenderTray(func() {
		state.tray.versionMenuItem.Label = fmt.Sprintf("⚠️ ytd (%s) (new version %s)", version, latest.TagName)
	})

	switch clickedAction {
	case "Update":
		// continue
	case "Changelog":
		state.runtime.Window.Show()
		state.runtime.Events.Emit("ytd:app:update:changelog", latest)
		return
	case "OK":
		state.runtime.Events.Emit("ytd:app:update:available", latest)
		return
	}

	state.Update(false)
}

func (state *AppState) Update(restart bool) {
	_, err := state.updater.Update()
	if err != nil {
		_, err := state.runtime.Dialog.Message(&dialog.MessageDialog{
			Type:         dialog.WarningDialog,
			Title:        "Update not successful",
			Message:      err.Error(),
			Buttons:      []string{"OK"},
			CancelButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return
		}
	}

	// restart for now does not work properly, it closes current ytd instance but not launch the new one
	// so notify user that update has been done and could restart app
	// err = u.Restart()
	state.runtime.Dialog.Message(&dialog.MessageDialog{
		Type:         dialog.InfoDialog,
		Title:        "Update successful",
		Message:      "Please restart ytd for the changes to take effect.",
		Buttons:      []string{"OK"},
		CancelButton: "OK",
	})
}

func (state *AppState) checkStartsAtLogin() {
	startsAtLogin, err := mac.StartsAtLogin()
	if err != nil {
		state.tray.reRenderTray(func() {
			state.tray.startAtLoginMenuItem.Label = "⚠ Start at Login unavailable"
			state.tray.startAtLoginMenuItem.Disabled = true
		})
	} else if startsAtLogin {
		mac.ShowNotification("Ytd", "App has been started in background", "", "")
		state.tray.reRenderTray(func() {
			state.tray.startAtLoginMenuItem.Checked = true
			state.tray.startAtLoginMenuItem.Disabled = false
		})
	}
}

func (state *AppState) SaveSettingBoolValue(name string, val bool) (err error) {
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

func (state *AppState) SaveSettingValue(name string, val string) (err error) {
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

func (state *AppState) ReadSettingBoolValue(name string) (bool, error) {
	return DbReadSettingBoolValue(name)
}

func (state *AppState) ReadSettingValue(name string) (string, error) {
	return DbReadSetting(name)
}

func (state *AppState) RemoveEntry(record map[string]interface{}) error {
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

func (state *AppState) AddToDownload(url string, isFromClipboard bool) error {
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

func (state *AppState) StartDownload(record map[string]interface{}) error {
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

func (state *AppState) IsSupportedUrl(url string) bool {
	for _, plugin := range plugins {
		if support := plugin.Supports(url); support {
			return true
		}
	}
	return false
}

func (state *AppState) IsFFmpegInstalled() (string, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if state.runtime.System.Platform() == "darwin" && err != nil {
		// on darwin check if ffmpeg is maybe installed by homebrew
		// (searching for ffmpeg only give wrong results if installed with homebrew)
		ffmpeg, err := exec.LookPath("/opt/homebrew/bin/ffmpeg")
		return ffmpeg, err
	}
	return ffmpeg, err
}

func (state *AppState) OpenUrl(url string) error {
	return state.runtime.Browser.Open(url)
}

func (state *AppState) ForceQuit() {
	state.runtime.Quit()
}

func (state *AppState) ShowWindow() {
	state.runtime.Window.Show()
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
