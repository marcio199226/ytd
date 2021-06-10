package clipboard

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ErrFormatNotSupported struct {
	msg        string
	formatType int
}

func (e *ErrFormatNotSupported) Error() string {
	return fmt.Sprint("This format is unsupported")
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return readAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return writeAll(text)
}

// Unsupported might be set true during clipboard init, to help callers decide
// whether or not to offer clipboard options.
var Unsupported bool

// Monitor starts monitoring the clipboard for changes. When
// a change is detected, it is sent over the channel.
func MonitorClipboard(interval time.Duration, ctx context.Context, wg *sync.WaitGroup, stopCh <-chan struct{}, changes chan<- string) error {
	currentValue, err := ReadAll()
	if err != nil {
		if _, ok := err.(*ErrFormatNotSupported); ok {
			fmt.Println(err)
		} else {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Monitor exiting...")
			close(changes)
			time.Sleep(time.Millisecond * 10)
			wg.Done()
			return nil
		default:
			if newValue, err := ReadAll(); err == nil {
				if newValue != currentValue {
					currentValue = newValue
					changes <- currentValue
				}
			} else {
				fmt.Println(err)
			}
		}
		time.Sleep(interval)
	}
}
