// spinner provides visual feedback for command line applications

package spinner

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/leaanthony/synx"
	isatty "github.com/mattn/go-isatty"
)

// Specialise the type
type status int

// Status code constants.
const (
	errorStatus status = iota
	successStatus
)

// Gets the default spinner frames based on the operating system.
func getDefaultSpinnerFrames() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"|", "/", "-", "\\"}
	default:
		return []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	}
}

// Gets the default status symbols based on the operating system.
func getStatusSymbols() (successSymbol, errorSymbol string) {
	switch runtime.GOOS {
	case "windows":
		return ">", "!"
	default:
		return "✓", "✗"
	}
}

// Spinner defines our spinner data.
type Spinner struct {
	message       *synx.String      // message to display
	stopChan      chan struct{}     // exit channel
	speedUpdated  *synx.Bool        // Indicates speed has been updated
	exitStatus    status            // Status of exit
	successSymbol *synx.String      // Symbol printed when Success() called
	errorSymbol   *synx.String      // Symbol printed when Error() called
	spinFrames    *synx.StringSlice // Spinset frames
	frameNumber   int               // Current frame [default 0]
	termWidth     *synx.Int         // Terminal width
	termHeight    *synx.Int         // Terminal Height
	spinSpeed     *synx.Int         // Delay between spinner updates in milliseconds [default 100ms]
	currentLine   *synx.String      // The current line being displayed
	running       *synx.Bool        // Indicates if the spinner is running
	abortMessage  *synx.String      // Printed when handling ctrl-c interrupt
	isTerminal    *synx.Bool        // Flag indicating if we are outputting to terminal
}

// NewSpinner creates a new spinner and sets up the default values.
func NewSpinner(optionalMessage ...string) *Spinner {
	successSymbol, errorSymbol := getStatusSymbols()
	// Blank message by default
	message := ""
	if len(optionalMessage) > 0 {
		message = optionalMessage[0]
	}
	result := &Spinner{
		message:       synx.NewString(message),
		stopChan:      make(chan struct{}),
		speedUpdated:  synx.NewBool(true),
		successSymbol: synx.NewString(successSymbol),
		errorSymbol:   synx.NewString(errorSymbol),
		spinFrames:    synx.NewStringSlice(getDefaultSpinnerFrames()),
		spinSpeed:     synx.NewInt(100),
		termWidth:     synx.NewInt(1),
		termHeight:    synx.NewInt(1),
		abortMessage:  synx.NewString("Aborted."),
		frameNumber:   0,
		running:       synx.NewBool(false),
		isTerminal:    synx.NewBool(isatty.IsTerminal(os.Stdout.Fd())),
	}

	return result
}

// New is solely here to make code cleaner for importers.
// EG: spinner.New(...)
func New(message ...string) *Spinner {
	return NewSpinner(message...)
}

// SetSuccessSymbol sets the symbol displayed on success.
func (s *Spinner) SetSuccessSymbol(symbol string) {
	s.successSymbol.SetValue(symbol)
}

// getSuccessSymbol sets the symbol displayed on error.
func (s *Spinner) getSuccessSymbol() string {
	return s.successSymbol.GetValue()
}

// SetErrorSymbol sets the symbol displayed on error.
func (s *Spinner) SetErrorSymbol(symbol string) {
	s.errorSymbol.SetValue(symbol)
}

// getErrorSymbol sets the symbol displayed on error.
func (s *Spinner) getErrorSymbol() (symbol string) {
	return s.errorSymbol.GetValue()
}

// SetSpinFrames makes the spinner use the given characters.
func (s *Spinner) SetSpinFrames(frames []string) {
	s.spinFrames.SetValue(frames)
}

func (s *Spinner) getNextSpinnerFrame() (result string) {
	// Check if the current frame is valid. If not, loop to start
	s.frameNumber = s.frameNumber % s.spinFrames.Length()
	result = s.spinFrames.GetElement(s.frameNumber)
	s.frameNumber++
	return
}

func (s *Spinner) getCurrentSpinnerFrame() (result string) {
	s.frameNumber = s.frameNumber % s.spinFrames.Length()
	result = s.spinFrames.GetElement(s.frameNumber)
	return result
}

// SetSpinSpeed sets the speed of the spinner animation.
// The lower the value, the faster the spin.
func (s *Spinner) SetSpinSpeed(ms int) {
	// Floor to a speed of 1
	if ms < 1 {
		ms = 1
	}
	s.spinSpeed.SetValue(ms)
	s.speedUpdated.SetValue(true)
}

// getSpinSpeed gets the speed of the spinner animation.
func (s *Spinner) getSpinSpeed() (ms int) {
	return s.spinSpeed.GetValue()
}

// UpdateMessage sets the spinner message.
// Can be flickery if not appending so use with care.
func (s *Spinner) UpdateMessage(message string) {
	// Clear line if this isn't an append.
	// for smoother screen updates.
	if strings.Index(message, s.getMessage()) != 0 {
		s.clearCurrentLine()
	}
	s.setMessage(message)
}

// SetAbortMessage sets the message that gets printed when
// the user kills the spinners by pressing ctrl-c.
func (s *Spinner) SetAbortMessage(message string) {
	s.abortMessage.SetValue(message)
}

func (s *Spinner) getAbortMessage() string {
	return s.abortMessage.GetValue()
}

func (s *Spinner) setMessage(message string) {
	s.message.SetValue(message)
}

func (s *Spinner) getMessage() string {
	return s.message.GetValue()
}

func (s *Spinner) getRunning() bool {
	return s.running.GetValue()
}

func (s *Spinner) setRunning(value bool) {
	s.running.SetValue(value)
}

func (s *Spinner) printSuccess(message string, args ...interface{}) {
	color.HiGreen(message, args...)
}

// Start the spinner!
func (s *Spinner) Start(optionalMessage ...string) {

	// If we're trying to start an already running spinner,
	// add a slight delay and retry. This allows the spinner
	// to complete a previous stop command gracefully.
	count := 0
	maxCount := 10
	for s.getRunning() == true && count < maxCount {
		//
		time.Sleep(time.Millisecond * 50)
		count++
	}

	// Did we fail?
	if count == maxCount {
		s.Error("Tried to start a running spinner with message: " + s.getMessage())
		return
	}

	// If we have a message, set it
	if len(optionalMessage) > 0 {
		s.setMessage(optionalMessage[0])
	}

	// make it look tidier.
	hideCursor()

	// Store the fact we are now running.
	s.setRunning(true)

	// Handle ctrl-c
	go func(stopChan chan struct{}) {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		// Notify and clean up
		s.stopChan <- struct{}{}
		fmt.Println("")
		color.HiRed("\r%s %s", s.getErrorSymbol(), s.getAbortMessage())
		os.Exit(1)
	}(s.stopChan)

	// spawn off a goroutine to handle the animation.
	go func() {

		ticker := time.NewTicker(time.Millisecond * time.Duration(s.spinSpeed.GetValue()))

		// Let's go!
		for {
			select {
			// For each frame tick
			case <-ticker.C:
				// Rewind to start of line and print the current frame and message.
				// Note: We don't fully clear the line here as this causes flickering.
				fmt.Printf("\r")
				fmt.Printf("%s %s", s.getNextSpinnerFrame(), s.getMessage())

				// Do we need to update the ticker?
				if s.speedUpdated.GetValue() == true {
					ticker.Stop()
					ticker = time.NewTicker(time.Millisecond * time.Duration(s.spinSpeed.GetValue()))
				}

			// If we get a stop signal
			case <-s.stopChan:

				// Store the fact we aren't running
				s.setRunning(false)

				// Quit the animation
				return
			}
		}
	}()
}

// stop will stop the spinner.
// The final message will either be the current message
// or the optional, given message.
// Success status will print the message in green.
// Error status will print the message in red.
func (s *Spinner) stop(message ...string) {

	var finalMessage = s.getMessage()

	// If we have an optional message, save it.
	if len(message) > 0 {
		finalMessage = message[0]
	}

	// Ensure we are running before issuing stop signal.
	if s.running.GetValue() {
		// Issue stop signal to animation.
		s.stopChan <- struct{}{}
	}

	// Clear the line, because a new message may be shorter than the original.
	s.clearCurrentLine()

	// Output the symbol and message depending on the status code.
	if s.exitStatus == errorStatus {
		color.HiRed("\r%s %s", s.getErrorSymbol(), finalMessage)
	} else {
		color.HiGreen("\r%s %s", s.getSuccessSymbol(), finalMessage)
	}

	// Show the cursor again
	showCursor()
}

// Error stops the spinner and sets the status code to error.
// Optional message to print instead of current message.
func (s *Spinner) Error(message ...string) {
	s.exitStatus = errorStatus
	s.stop(message...)
}

// Errorf stops the spinner, formats and sets the status code to error.
// Formats and prints the given message instead of current message.
func (s *Spinner) Errorf(format string, args ...interface{}) {
	s.exitStatus = errorStatus
	message := fmt.Sprintf(format, args...)
	s.stop(message)
}

// Success stops the spinner and sets the status code to success.
// Optional message to print instead of current message.
func (s *Spinner) Success(message ...string) {
	s.exitStatus = successStatus
	s.stop(message...)
}

// Successf stops the spinner, formats and sets the status code to success.
// Formats and prints the given message instead of current message.
func (s *Spinner) Successf(format string, args ...interface{}) {
	s.exitStatus = successStatus
	message := fmt.Sprintf(format, args...)
	s.stop(message)
}
