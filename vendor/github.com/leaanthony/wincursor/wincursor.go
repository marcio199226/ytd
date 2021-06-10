// +build windows

package wincursor

// Some parts from https://github.com/Azure/go-ansiterm

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32              = syscall.NewLazyDLL("kernel32.dll")
	getConsoleCursorInfoProc = modkernel32.NewProc("GetConsoleCursorInfo")
	setConsoleCursorInfoProc = modkernel32.NewProc("SetConsoleCursorInfo")
)

// ConsoleCursorInfo holds the cursor information data
type ConsoleCursorInfo struct {
	Size    uint32
	Visible int32
}

// checkError evaluates the results of a Windows API call and returns the error if it failed.
func checkError(r1, r2 uintptr, err error) error {
	// Windows APIs return non-zero to indicate success
	if r1 != 0 {
		return nil
	}

	// Return the error if provided, otherwise default to EINVAL
	if err != nil {
		return err
	}
	return syscall.EINVAL
}

// getConsoleCursorInfo - get the cursor information
func getConsoleCursorInfo(handle uintptr, cursorInfo *ConsoleCursorInfo) error {
	r1, r2, err := getConsoleCursorInfoProc.Call(handle, uintptr(unsafe.Pointer(cursorInfo)), 0)
	return checkError(r1, r2, err)
}

// setConsoleCursorInfo - set the cursor information
func setConsoleCursorInfo(handle uintptr, cursorInfo *ConsoleCursorInfo) error {
	r1, r2, err := setConsoleCursorInfoProc.Call(handle, uintptr(unsafe.Pointer(cursorInfo)), 0)
	return checkError(r1, r2, err)
}

// setCursorVisible - takes a boolean and sets the visibility of
// the cursor to this value
func setCursorVisible(visible bool) error {
	var cursor ConsoleCursorInfo
	// -11 is stdout on windows
	fd, err := syscall.GetStdHandle(-11)
	if err != nil {
		return err
	}
	err = getConsoleCursorInfo(uintptr(fd), &cursor)
	if err != nil {
		return err
	}
	if visible {
		cursor.Visible = 1
	} else {
		cursor.Visible = 0
	}
	err = setConsoleCursorInfo(uintptr(fd), &cursor)
	return err
}

// Show shows the windows console cursor
func Show() error {
	return setCursorVisible(true)
}

// Hide hides the windows console cursor
func Hide() error {
	return setCursorVisible(false)
}
