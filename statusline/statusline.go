package statusline

import (
	"fmt"
	"sync"
	"time"
)

type StatusLineOptions struct {
	SpinnerSet     []string
	SpinnerStyling string
}

type StatusMessage struct {
	Message string
	Emit    bool
}

type StatusLine struct {
	wg          sync.WaitGroup
	options     StatusLineOptions
	messageChan chan StatusMessage
	killChan    chan struct{}
}

var defaultSpinner = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

type Option func(opt *StatusLineOptions)

func WithSpinner(spinnerSet []string) Option {
	return func(opt *StatusLineOptions) {
		opt.SpinnerSet = spinnerSet
	}
}

func WithSpinnerStyling(style string) Option {
	return func(opt *StatusLineOptions) {
		opt.SpinnerStyling = style
	}
}

func NewStatusLine(options ...Option) *StatusLine {
	s := &StatusLine{
		messageChan: make(chan StatusMessage, 100),
		killChan:    make(chan struct{}, 1),
	}
	s.options.SpinnerSet = defaultSpinner

	for _, opt := range options {
		opt(&s.options)
	}

	s.wg.Add(1)
	go s.runLoop()

	return s
}

func (s *StatusLine) Stop() {
	close(s.killChan)
	s.wg.Wait()
}

func (s *StatusLine) displayStatusMessage(roller string, message string) {
	fmt.Printf("\r%s%s\x1b[0m %s\x1b[0m", s.options.SpinnerStyling, roller, message)
	fmt.Printf("\x1b[K")
}

func (s *StatusLine) emitMessage(message string) {
	fmt.Printf("\r%s\x1b[0m", message)
	fmt.Printf("\x1b[K\n")
}

func (s *StatusLine) runLoop() {
	defer s.wg.Done()
	defer fmt.Printf("\x1b[?25h")
	var statusMessage string
	var rollerIndex int

	ticker := time.NewTicker(100 * time.Millisecond)

	fmt.Printf("\x1b[?25l")

	for {
		select {
		case <-s.killChan:
			ticker.Stop()
			return
		case m := <-s.messageChan:
			switch m.Emit {
			case true:
				statusMessage = ""
				s.emitMessage(m.Message)
			case false:
				statusMessage = m.Message
				s.displayStatusMessage(s.options.SpinnerSet[rollerIndex], statusMessage)
			}
		case <-ticker.C:
			if statusMessage != "" {
				rollerIndex++
				if rollerIndex >= len(s.options.SpinnerSet) {
					rollerIndex = 0
				}
				s.displayStatusMessage(s.options.SpinnerSet[rollerIndex], statusMessage)
			}
		}
	}
}

// Status displays a persistent status message with a roller, indicating that the status is ongoing.
// Status will not advance the current line when updated and thus status can be updated with new
// messages on the same line as the status changees
func (s *StatusLine) Status(message string) {
	s.messageChan <- StatusMessage{message, false}
}

// Emit emits a new message in the style of a running log.  It will clear the current line and
// remove any status message which may be present.  Typically the use case is to emit a message
// once it enters the log record, in order to conclude the current ongoing operation
func (s *StatusLine) Emit(message string) {
	s.messageChan <- StatusMessage{message, true}
}
