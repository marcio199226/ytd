package constants

const (
	TrackWithoutAudioFormat = constError("Track has not audio/webm format available")
	CannotDownloadTrack     = ""
)

type constError string

func (e constError) Error() string {
	return string(e)
}
