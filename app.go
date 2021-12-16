package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "ytd/constants"
	. "ytd/db"
	. "ytd/models"
	. "ytd/plugins"

	"github.com/denisbrodbeck/machineid"
	"github.com/leonelquinteros/gotext"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/mac"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
	"github.com/xujiajun/nutsdb"
)

var i18nPath string

var wailsRuntime *wails.Runtime

type HostInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AppState struct {
	runtime          *wails.Runtime
	db               *nutsdb.DB
	plugins          []Plugin
	Entries          []GenericEntry    `json:"entries"`
	OfflinePlaylists []OfflinePlaylist `json:"offlinePlaylists"`
	Config           *AppConfig        `json:"config"`
	Stats            *AppStats         `json:"stats"`
	AppVersion       string            `json:"appVersion"`
	Host             HostInfo          `json:"host"`
	PwaUrl           string            `json:"pwaUrl"`
	context.Context  `json:"-"`

	isInForeground         bool
	canStartAtLogin        bool
	tray                   TrayMenu
	updater                *Updater
	offlinePlaylistService *OfflinePlaylistService
	ngrok                  *NgrokService
	convertQueue           chan GenericEntry
}

func (state *AppState) PreWailsInit(ctx context.Context) {
	state.db = InitializeDb()
	state.Entries = DbGetAllEntries()
	state.OfflinePlaylists = DbGetAllOfflinePlaylists()
	state.offlinePlaylistService = &OfflinePlaylistService{}
	state.ngrok = &NgrokService{}
	state.ngrok.Context = ctx
	state.Config = state.Config.Init()
	state.AppVersion = version
	state.convertQueue = make(chan GenericEntry)

	// configure i18n paths
	gotext.Configure("/Users/oskarmarciniak/projects/golang/ytd/i18n", state.Config.Language, "default")

	if _, err := mac.StartsAtLogin(); err == nil {
		state.canStartAtLogin = true
	}

	if machineid, err := machineid.ProtectedID("ytd"); err == nil {
		state.Host.ID = machineid
	}
	if user, err := user.Current(); err == nil {
		state.Host.Username = user.Username
	}
}

func (state *AppState) WailsInit(runtime *wails.Runtime) {
	// Save runtime
	state.runtime = runtime
	state.offlinePlaylistService.runtime = runtime
	state.ngrok.runtime = runtime
	state.Config.SetRuntime(runtime)
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
		plugin.SetContext(state.Context)
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

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		pending := make(chan int, maxConvertJobs)

		for {
			select {
			case entry := <-state.convertQueue:
				go func() {
					if len(pending) == cap(pending) {
						state.runtime.Events.Emit("ytd:track:convert:queued", entry)
					}
					pending <- 1
					state.convertToMp3(ticker, &entry, true)
					<-pending
				}()
			case <-ticker.C:
				// get entries that could be converted
				for _, t := range DbGetAllEntries() {
					entry := t
					plugin := getPluginFor(entry.Source)

					if entry.Type == "track" && entry.Track.Status == TrackStatusDownladed && !entry.Track.IsConvertedToMp3 && plugin.IsTrackFileExists(entry.Track, "webm") {
						// skip tracks which has failed at least 3 times in a row
						if entry.Track.ConvertingStatus.Attempts >= 3 {
							fmt.Printf("Skipping audio extraction for %s(%s)...due to too many attempts\n", entry.Track.Name, entry.Track.ID)
							continue
						}

						go func() {
							// this will block until pending channel is full
							pending <- 1
							state.convertToMp3(ticker, &entry, false)
							<-pending
						}()
					}
				}
			default:
				fmt.Printf("Convert to mp3 nothing to do....max parralel %d | converting %d | pending %d entries\n\n", maxConvertJobs, state.Stats.ConvertingCount, len(pending))
				time.Sleep(1 * time.Second)
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

	if state.Config.PublicServer.Enabled {
		ShowLoader(state.runtime, "Starting public server...")
		result := state.ngrok.StartProcess(false)
		if result.err != nil {
			state.runtime.Events.Emit("ytd:ngrok", NgrokStateEventPayload{Status: NgrokStatusError, ErrCode: result.errCode})
			HideLoader(state.runtime)
			return
		}
		state.runtime.Events.Emit("ytd:ngrok", NgrokStateEventPayload{Status: result.status, Url: result.publicUrl})
		HideLoader(state.runtime)

		// monitor ngrok state
		go state.ngrok.MonitorNgrokProcess()
	}
}

func (state *AppState) WailsShutdown() {
	err := state.ngrok.KillProcess()
	if err != nil {
		fmt.Println("WailsShutdown state.ngrok.KillProcess() failed", err)
	}
	err = state.db.Merge()
	if err != nil {
		fmt.Println("WailsShutdown db.Merge() failed", err)
	}
	CloseDb()
}

func (state *AppState) ReloadNewLanguage() {
	gotext.Configure(i18nPath, state.Config.Language, "default")
	// re render tray to take effect for new language translations
	state.tray.reRenderTray(func() {})
	// do other stuff if needed to reload translations from some ui native elements
}

func (state *AppState) InitializeListeners() {
	state.runtime.Events.On("ytd:app:foreground", func(optionalData ...interface{}) {
		var json map[string]interface{} = optionalData[0].(map[string]interface{})
		if isInForeground, ok := json["isInForeground"]; ok {
			state.isInForeground = isInForeground.(bool)
		}
	})

	state.runtime.Events.On("ytd:offline:playlists:addedTrack", func(optionalData ...interface{}) {
		state.OfflinePlaylists, _ = state.offlinePlaylistService.GetPlaylists(true)
	})

	state.runtime.Events.On("ytd:offline:playlists:removedTrack", func(optionalData ...interface{}) {
		state.OfflinePlaylists, _ = state.offlinePlaylistService.GetPlaylists(true)
	})

	state.runtime.Events.On("ytd:offline:playlists:created", func(optionalData ...interface{}) {
		state.OfflinePlaylists, _ = state.offlinePlaylistService.GetPlaylists(true)
	})

	state.runtime.Events.On("ytd:offline:playlists:removed", func(optionalData ...interface{}) {
		state.OfflinePlaylists, _ = state.offlinePlaylistService.GetPlaylists(true)
	})

	state.runtime.Events.On("ngrok:configured", func(optionalData ...interface{}) {
		if state.Config.PublicServer.Enabled {
			ShowLoader(state.runtime, "Configuring public server...")
			result := state.ngrok.StartProcess(true)
			if result.err != nil {
				state.runtime.Events.Emit("ytd:ngrok", NgrokStateEventPayload{Status: NgrokStatusError, ErrCode: result.errCode})
				HideLoader(state.runtime)
				return
			}
			state.runtime.Events.Emit("ytd:ngrok", NgrokStateEventPayload{Status: result.status, Url: result.publicUrl})
			HideLoader(state.runtime)

			// monitor ngrok process
			go state.ngrok.MonitorNgrokProcess()
		}

		if !state.Config.PublicServer.Enabled {
			err := state.ngrok.KillProcess()
			if err != nil {
				SendNotification(state.runtime, NotificationEventPayload{Type: "error", Label: "Cannot shutdown public server"}, state.isInForeground)
			}
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
		Title:                gotext.Get("Choose directory"),
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

func (state *AppState) convertToMp3(restartTicker *time.Ticker, entry *GenericEntry, force bool) error {
	if !state.Config.ConvertToMp3 {
		// if option is not enabled restart check after 3s
		restartTicker.Stop()
		restartTicker.Reset(3 * time.Second)
		return nil
	}

	fmt.Println("Converting....")
	ffmpeg, _ := state.IsFFmpegInstalled()
	plugin := getPluginFor(entry.Source)

	if entry.Type == "track" && entry.Track.Status == TrackStatusDownladed && !entry.Track.IsConvertedToMp3 && plugin.IsTrackFileExists(entry.Track, "webm") {
		// skip tracks which has failed at least 3 times in a row
		if entry.Track.ConvertingStatus.Attempts >= 3 && !force {
			fmt.Printf("Skipping audio extraction for %s - (%s)...due to too many attempts\n", entry.Track.Name, entry.Track.ID)
			return nil
		}

		fmt.Printf("Extracting audio for %s - (%s)...\n", entry.Track.Name, entry.Track.ID)
		entry.Track.ConvertingStatus.Status = TrakcConverting
		DbWriteEntry(entry.Track.ID, entry)
		state.Stats.IncConvertCount()
		state.runtime.Events.Emit("ytd:track", entry)

		// ffmpeg -i "41qC3w3UUkU.webm" -vn -ab 128k -ar 44100 -y "41qC3w3UUkU.mp3"
		outputPath := fmt.Sprintf("%s/%s.mp3", plugin.GetDir(), entry.Track.ID)
		cmd := exec.CommandContext(
			state.Context,
			ffmpeg,
			"-loglevel", "quiet",
			"-i", fmt.Sprintf("%s/%s.webm", plugin.GetDir(), entry.Track.ID),
			"-vn",
			"-ab", "128k",
			"-ar", "44100",
			"-y", outputPath,
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
			state.Stats.DecConvertCount()
			state.runtime.Events.Emit("ytd:track", entry)
		} else {
			entry.Track.ConvertingStatus.Status = TrakcConverted
			entry.Track.IsConvertedToMp3 = true

			// check new filesize and save it
			fileInfo, err := os.Stat(outputPath)
			if err == nil {
				entry.Track.ConvertingStatus.Filesize = int(fileInfo.Size())
			}

			DbWriteEntry(entry.Track.ID, entry)
			state.Stats.DecConvertCount()
			state.runtime.Events.Emit("ytd:track", entry)

			// remove webm if needed
			if state.Config.CleanWebmFiles && plugin.IsTrackFileExists(entry.Track, "webm") {
				err = os.Remove(fmt.Sprintf("%s/%s.webm", plugin.GetDir(), entry.Track.ID))
				if err != nil && !os.IsNotExist(err) {
					fmt.Printf("Cannot remove %s.webm file after successfull converting to mp3\n", entry.Track.ID)
				}
			}
		}
		return nil
	}
	return nil
}

func (state *AppState) checkForUpdates() {
	if !state.Config.CheckForUpdates {
		return
	}

	state.updater = &Updater{
		CurrentVersion:              version,
		LatestReleaseGitHubEndpoint: "https://api.github.com/repos/marcio199226/ytd/releases",
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
			Title:        gotext.Get("Update check failed"),
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
			Title:        gotext.Get("You're up to date"),
			Message:      fmt.Sprintf(gotext.Get("%s is the latest version.", latest.TagName)),
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
		Title:        fmt.Sprintf(gotext.Get("New version available: %s", latest.TagName)),
		Message:      gotext.Get("Would you like to update?"),
		Buttons:      []string{"OK", "Changelog", gotext.Get("Update")},
		CancelButton: "OK",
	})

	state.tray.reRenderTray(func() {
		state.tray.versionMenuItem.Label = fmt.Sprintf(gotext.Get("⚠️ ytd (%s) (new version %s)", version, latest.TagName))
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
			Title:        gotext.Get("Update not successful"),
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
		Title:        gotext.Get("Update successful"),
		Message:      gotext.Get("Please restart ytd for the changes to take effect."),
		Buttons:      []string{"OK"},
		CancelButton: "OK",
	})
}

func (state *AppState) checkStartsAtLogin() {
	startsAtLogin, err := mac.StartsAtLogin()
	if err != nil {
		state.tray.reRenderTray(func() {
			state.tray.startAtLoginMenuItem.Label = gotext.Get("⚠ Start at Login unavailable")
			state.tray.startAtLoginMenuItem.Disabled = true
		})
	} else if startsAtLogin {
		mac.ShowNotification("Ytd", gotext.Get("App has been started in background"), "", "")
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
				if err != nil && !os.IsNotExist(err) {
					return err
				}
			}
			// remove mp3 if file has been already converted
			if plugin.IsTrackFileExists(entry.Track, "mp3") {
				err = os.Remove(fmt.Sprintf("%s/%s.mp3", plugin.GetDir(), entry.Track.ID))
				if err != nil && !os.IsNotExist(err) {
					return err
				}
			}

			// remove track from playlists
			for _, p := range state.OfflinePlaylists {
				for _, tid := range p.TracksIds {
					if tid == entry.Track.ID {
						state.offlinePlaylistService.RemoveTrackFromPlaylist(tid, p)
					}
				}
			}
			state.OfflinePlaylists, _ = state.offlinePlaylistService.GetPlaylists(true)
			return nil
		}
	}
	return err
}

func (state *AppState) AddToDownload(url string, isFromClipboard bool) error {
	for _, plugin := range plugins {
		if support := plugin.Supports(url); support {
			if appState.GetAppConfig().ConcurrentDownloads {
				go func() {
					newEntries <- plugin.Fetch(url, isFromClipboard)
				}()
			} else {
				newEntries <- plugin.Fetch(url, isFromClipboard)
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

func (state *AppState) AddToConvertQueue(entry GenericEntry) error {
	state.convertQueue <- entry
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
