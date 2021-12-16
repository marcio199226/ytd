package constants

const (
	TrackStatusPending     = "pending"
	TrackStatusProcessing  = "processing"
	TrackStatusFailed      = "failed"
	TrackStatusDownloading = "downloading"
	TrackStatusDownladed   = "downloaded"
	TrakcConverting        = "converting"
	TrakcConverted         = "converted"
	TrakcConvertFailed     = "failed"
	TrakcConvertQueued     = "queued"

	NgrokStatusRunning = "running"
	NgrokStatusError   = "error"
	NgrokStatusTimeout = "timeout"
	NgrokStatusStopped = "killed"

	NotificationTypeSuccess = "success"
	NotificationTypeWarning = "warning"
	NotificationTypeError   = "error"
)
