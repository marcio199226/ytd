package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"sync"
	"time"

	. "ytd/clipboard"
	. "ytd/db"
	. "ytd/models"
	. "ytd/plugins"

	_ "embed"

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

var appConfig AppConfig

type AppState struct {
	log     *wails.CustomLogger
	runtime *wails.Runtime
	db      *nutsdb.DB
	plugins []Plugin
	Entries []GenericEntry `json:"entries"`
	Config  AppConfig      `json:"config"`
	Stats   struct {
		DownloadingCount uint
	}
}

func (state *AppState) WailsInit(runtime *wails.Runtime) error {
	// Save runtime
	state.runtime = runtime
	state.log = runtime.Log.New("AppState")
	// Do some other initialisation

	state.db = InitializeDb()

	for _, plugin := range plugins {
		plugin.SetWailsRuntime(runtime)
	}

	state.Entries = DbGetAllEntries()
	// this is sync so it blocks until finished and wails:loaded are not dispatched until this finishes
	runtime.Events.On("wails:loaded", func(...interface{}) {
		// entries := DbGetAllEntries()
		fmt.Println("EMIT YTD:ONLOAD")
		runtime.Events.Emit("ytd:onload", state)
	})

	appConfig = state.Config.Init()
	state.Config = appConfig
	fmt.Println("APP STATE INITIALIZED")

	go func() {
		for {
			time.Sleep(10 * time.Second)
			state.checkForTracksToDownload()
		}
	}()

	return nil
}

func (state *AppState) GetAppConfig() AppConfig {
	return state.Config
}

func (state *AppState) checkForTracksToDownload() error {
	if state.Stats.DownloadingCount >= state.Config.MaxParrallelDownloads {
		return nil
	}

	for _, t := range state.Entries {
		if t.Type == "track" && !t.Track.Downloaded {

		}
	}
	return nil
}

func saveSettingBoolValue(name string, val bool) error {
	return DbSaveSettingBoolValue(name, val)
}

func saveSettingValue(name string, val string) error {
	return DbWriteSetting(name, val)
}

func readSettingBoolValue(name string) (bool, error) {
	return DbReadSettingBoolValue(name)
}

func readSettingValue(name string) (string, error) {
	return DbReadSetting(name)
}

func addToDownload(url string) error {
	// usare lo logica di quando si riceve un copia
	fmt.Println(appConfig)
	return nil
	for _, plugin := range plugins {
		if support := plugin.Supports(url); support {
			if appConfig.ConcurrentDownloads {
				go func() {
					plugin.Fetch(url)
				}()
			} else {
				plugin.Fetch(url)
			}
			continue
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
		Width:  1024,
		Height: 768,
		Title:  "ytd",
		JS:     js,
		CSS:    css,
		//HTML:             html,
		Colour:           "#131313",
		DisableInspector: false,
	})
	app.Bind(&AppState{})
	app.Bind(saveSettingBoolValue)
	app.Bind(saveSettingValue)
	app.Bind(readSettingBoolValue)
	app.Bind(readSettingValue)
	app.Bind(addToDownload)

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
					addToDownload(change)
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
