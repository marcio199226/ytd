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
	. "ytd/plugins"

	_ "embed"

	"github.com/wailsapp/wails"
)

func basic() string {
	return "World!"
}

//go:embed frontend/dist/my-app/main.js
var js string

//go:embed frontend/dist/my-app/styles.css
var css string

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "ytd",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)

	var plugins []Plugin = []Plugin{&Yt{Name: "youtube"}}
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
							fetchErr := plugin.Fetch(change)
							if fetchErr != nil {
								fmt.Printf("Unable to download %s \n", change)
							}
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
