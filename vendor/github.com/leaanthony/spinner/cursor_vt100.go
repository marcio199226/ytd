// +build !windows

package spinner

import (
	"fmt"
)

func init() {

}

func showCursor() {
	fmt.Printf("\033[?25h")
}

func hideCursor() {
	fmt.Printf("\033[?25l")
}

func (s *Spinner) clearCurrentLine() {
	fmt.Printf("\r\033[0K")
}
