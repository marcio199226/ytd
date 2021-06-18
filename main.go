package main

import (
	"context"
	"fmt"
	"log"
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

//go:embed frontend/dist/my-app/main.js
var js string

//go:embed frontend/dist/my-app/styles.css
var css string

type AppState struct {
	log     *wails.CustomLogger
	runtime *wails.Runtime
	db      *nutsdb.DB
	plugins []Plugin
	Entries []*GenericEntry
	Config  AppConfig
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

	// this is sync so it blocks until finished and wails:loaded are not dispatched until this finishes
	runtime.Events.Once("wails:loaded", func(...interface{}) {
		entries := DbGetAllEntries()
		runtime.Events.Emit("ytd:onload", entries)
	})

	fmt.Println("APP STATE INITIALIZED")
	return nil
}

func saveSettingBoolValue(name string, val bool) error {
	var v string
	if val {
		v = "1"
	} else {
		v = "0"
	}
	return DbWriteSetting(name, v)
}

func saveSettingValue(name string, val string) error {
	return DbWriteSetting(name, val)
}

func readSettingBoolValue(name string) (bool, error) {
	val, err := DbReadSetting(name)
	if err != nil {
		return false, err
	}
	if val == "1" {
		return true, nil
	}
	return false, nil
}

func readSettingValue(name string) (string, error) {
	return DbReadSetting(name)
}

func addToDownload(url string) error {
	// usare lo logica di quando si riceve un copia
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
		Colour: "#131313",
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
					for _, plugin := range plugins {
						if support := plugin.Supports(change); support {
							go func() {
								plugin.Fetch(change)
							}()
							continue
						}
					}

				} else {
					log.Println("channel has been closed. exiting...")
					time.Sleep(time.Millisecond)
				}
			}
		}
	}()

	app.Run()
}
