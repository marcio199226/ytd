// +build windows

package clipboard

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                  = syscall.MustLoadDLL("user32")
	openClipboard           = user32.MustFindProc("OpenClipboard")
	closeClipboard          = user32.MustFindProc("CloseClipboard")
	emptyClipboard          = user32.MustFindProc("EmptyClipboard")
	getClipboardData        = user32.MustFindProc("GetClipboardData")
	setClipboardData        = user32.MustFindProc("SetClipboardData")
	setClipboardViewer      = user32.MustFindProc("SetClipboardViewer")
	CountClipboardFormats   = user32.MustFindProc("CountClipboardFormats")
	EnumClipboardFormats    = user32.MustFindProc("EnumClipboardFormats")
	GetClipboardFormatNameW = user32.MustFindProc("GetClipboardFormatNameW")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

const gmemMoveable = 0x0002

const (
	CF_UNICODETEXT = 13
	CF_TEXT        = 1
	CF_BITMAP      = 2
	CF_OEMTEXT     = 7
	CF_DIB         = 8
	CF_DIBV5       = 17
)

var supportedFormats [4]int = [4]int{CF_UNICODETEXT, CF_TEXT, CF_OEMTEXT}

func isSupportedFormat(format int) bool {
	for i := 0; i < len(supportedFormats); i++ {
		if supportedFormats[i] == format {
			return true
		}
	}
	return false
}

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	var error error
	for time.Now().Before(limit) {
		r, _, error = openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return error
}

func readAll() (string, error) {
	err := waitOpenClipboard()
	if err != nil {
		return "", err
	}
	defer closeClipboard.Call()

	formatsCount, _, _ := CountClipboardFormats.Call(0)
	fmt.Printf("Available formats: %d\n", formatsCount)
	availableFormat, _, _ := EnumClipboardFormats.Call(0)
	fmt.Print("Available content formats: ")

	hasSupportedFormat := false
	for i := 0; i < (int)(formatsCount); i++ {
		index := 0
		if i > 0 {
			index = i
		}
		format, _, _ := EnumClipboardFormats.Call(uintptr(index))
		if format != 0 {
			if isSupported := isSupportedFormat((int)(format)); isSupported {
				hasSupportedFormat = true
				break
			}
		}
	}

	if !hasSupportedFormat {
		return "", &ErrFormatNotSupported{formatType: (int)(availableFormat)}
	}

	h, _, err := getClipboardData.Call(CF_UNICODETEXT)
	if h == 0 {
		return "", err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return "", err
	}

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])
	r, _, err := globalUnlock.Call(h)
	if r == 0 {
		return "", err
	}

	return text, nil
}

func writeAll(text string) error {
	err := waitOpenClipboard()
	if err != nil {
		return err
	}
	defer closeClipboard.Call()

	r, _, err := emptyClipboard.Call(0)
	if r == 0 {
		return err
	}

	data := syscall.StringToUTF16(text)

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	h, _, err := globalAlloc.Call(gmemMoveable, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		return err
	}
	defer func() {
		if h != 0 {
			globalFree.Call(h)
		}
	}()

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return err
	}

	r, _, err = lstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		return err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		if err.(syscall.Errno) != 0 {
			return err
		}
	}

	r, _, err = setClipboardData.Call(CF_UNICODETEXT, h)
	if r == 0 {
		return err
	}
	h = 0 // suppress deferred cleanup
	return nil
}
