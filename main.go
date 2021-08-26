package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"sync"
	"time"

	. "ytd/clipboard"
	. "ytd/models"
	. "ytd/plugins"

	_ "embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

var plugins []Plugin = []Plugin{&Yt{Name: "youtube"}}

//go:embed frontend/dist/assets/*
var static embed.FS

var appState *AppState
var newEntries = make(chan GenericEntry)

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do your cors stuff
		// return if you do not want the FileServer handle a specific request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func main() {
	app := &AppState{}
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
						app.AddToDownload(change, true)
					}
				} else {
					log.Println("channel has been closed. exiting...")
					time.Sleep(time.Millisecond)
				}
			}
		}
	}()

	go func() {
		fs := http.StripPrefix("/static/", http.FileServer(http.FS(static)))
		http.Handle("/tracks/", http.StripPrefix("/tracks/", http.FileServer(http.Dir(currentDir))))
		http.Handle("/static/", cors(fs))
		http.ListenAndServe(":8080", nil)
	}()

	app.PreWailsInit()
	err := wails.Run(&options.App{
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		StartHidden:       app.Config.StartAtLogin,
		HideWindowOnClose: app.Config.RunInBackgroundOnClose,
		DisableResize:     false,
		Fullscreen:        false,
		Startup:           app.WailsInit,
		Shutdown:          app.WailsShutdown,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			TitleBar:                      mac.TitleBarHiddenInset(),
			ActivationPolicy:              mac.NSApplicationActivationPolicyAccessory,
		},
		Title: "ytd",
		Bind: []interface{}{
			app,
		},
		Frameless: false,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
