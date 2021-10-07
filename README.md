# ytd

Dekstop app for downloading audio tracks from youtube built with wails & angular

## Main features:

## Screenshots

## Build from sources

#### Dev env

Angular

`npm run serve`

Wails

`wails dev --e "html"`

Open tab in chrome (preffered) and go to http://localhost:4200

#### Build binaries (Macos only at the moment)

`wails build --platform darwin/arm64 --clean --package --production`

`wails build --platform darwin/amd64 --clean --package --production --upx`

(*) upx doesn't work for apple m1 https://github.com/upx/upx/issues/446

## Roadmap
- [ ] Chrome extension so tracks may be downloaded without user interaction (even without copy yt links)
- [ ] Internalization
- [ ] Share tracks through telegram (user could subscribe to ytd bot and then will be able to send downloaded tracks to yourself telegram account)
