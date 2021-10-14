# ytd

Dekstop app for downloading audio tracks from youtube built with wails & angular

**(*) please note this is an alpha version in case of malfunction please open an issue**

## Install
[Download from releases page](https://github.com/marcio199226/ytd/releases)

Supported platfroms for now
- Macos (tested on 11.5.x apple m1)

## Main features:
- check for updates & app update
- clipboard watch (once yt link is copied it will be automatically downloaded, **can be disabled**)
- run in bg on close (run app in bg even if you closed app window)
- convert webm files to mp3
- system tray with fast settings
- in app player for single tracks and playlists (for now only for offline created playlists)
- clean and simple UI (I hope ;))
- create offline playlist
  - playback playlist
  - add/remove tracks
  - export to any external devices (pen drive , external hd etc...) or any folder
- Made with :green_heart: with golang & angular in my spare time

## Screenshots
<img width="1552" alt="home" title="Home" src="https://user-images.githubusercontent.com/10244404/136397513-8bde9053-0d1f-4257-83f5-71322e166cb1.png">
<img width="1552" alt="settings_dialog" title="Settings" src="https://user-images.githubusercontent.com/10244404/136397782-3962a502-0b26-4048-adc0-50e871f2e6f4.png">
<img width="1552" alt="downloading_track" title="Downloading track" src="https://user-images.githubusercontent.com/10244404/136397721-6ccebfaa-b7b7-4d53-94de-1cbdd0ac8ea3.png">
<img width="1552" alt="track_playback" title="Playback track" src="https://user-images.githubusercontent.com/10244404/136397829-eaae8d9f-f5a4-422d-aba1-a23bba2a2501.png">
<img width="1552" alt="playlist_playback" title="Playback playlist" src="https://user-images.githubusercontent.com/10244404/136397770-79cfdef2-8daf-4bd4-a1e4-3b969c2031f9.png">


## Build from sources

Wails requirements: https://wails.io/docs/gettingstarted/installation

#### Dev env

Angular

`cd frontend && npm install && npm run serve`

Wails

`wails dev --e "html"`

Open tab in chrome (preffered) and go to http://localhost:4200

Extract translation for golang side:

`xgotext -exclude "vendor,frontend" -in "/Users/oskarmarciniak/projects/golang/ytd" -out "/Users/oskarmarciniak/projects/golang/ytd/i18n"`

#### Build binaries (Macos only at the moment)

`wails build --platform darwin/arm64 --clean --package --production`

`wails build --platform darwin/amd64 --clean --package --production --upx`

(*) upx doesn't work for apple m1 https://github.com/upx/upx/issues/446

## Roadmap
- [ ] Chrome extension so tracks may be downloaded without user interaction (even without copy yt links)
- [ ] Internalization
- [ ] Share tracks through telegram (user could subscribe to ytd bot and then will be able to send downloaded tracks to yourself telegram account)
- [ ] Download playlists from yt (exports them, search for playlist, playback playlist etc...)
