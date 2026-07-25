package tui

import (
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/exadrift/go/tui/internal/terminal"
	"golang.org/x/term"
)

type Application struct {
	inputHandler func(string) string
	inFocus      Widget
	ctrlCExit    bool
	termFd       int
	onStart      func() error
	onExit       func()
	sigChan      chan os.Signal
	closeLock    sync.Mutex
	closeChan    chan struct{}
	wg           *sync.WaitGroup
	root         Widget
}

func WithApplicationOptionInputHandler(handleInput func(string) string) *Option {
	return &Option{
		optionType: ApplicationOptionInputHandler,
		data:       handleInput,
	}
}

func WithApplicationOptionOnStart(onStart func() error) *Option {
	return &Option{
		optionType: ApplicationOptionWithOnStart,
		data:       onStart,
	}
}

func WithApplicationOptionOnExit(onExit func()) *Option {
	return &Option{
		optionType: ApplicationOptionWithOnExit,
		data:       onExit,
	}
}

func WithApplicationOptionExitSignals(signals ...syscall.Signal) *Option {
	exitSignals := make([]syscall.Signal, len(signals))
	copy(exitSignals, signals)
	return &Option{
		optionType: ApplicationOptionExitSignals,
		data:       exitSignals,
	}
}

// New constructs and returns a new Application with a root widget specified.
func New(root Widget, options ...Option) *Application {
	var exitSignals []os.Signal

	app := &Application{
		wg:        &sync.WaitGroup{},
		root:      root,
		sigChan:   make(chan os.Signal, 1),
		closeChan: make(chan struct{}, 1),
	}

	for _, option := range options {
		switch option.optionType {
		case ApplicationOptionWithOnStart:
			app.onStart = option.data.(func() error)
		case ApplicationOptionWithOnExit:
			app.onExit = option.data.(func())
		case ApplicationOptionExitSignals:
			exitSignals = option.data.([]os.Signal)
		case ApplicationOptionInputHandler:
			app.inputHandler = option.data.(func(string) string)
		default:
			panic("unknown application option")
		}
	}

	if exitSignals == nil {
		exitSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	if slices.Contains(exitSignals, os.Signal(syscall.SIGINT)) {
		app.ctrlCExit = true
	}
	signal.Notify(app.sigChan, exitSignals...)

	return app
}

func (a *Application) Exit() {
	a.closeLock.Lock()
	defer a.closeLock.Unlock()

	select {
	case _, ok := <-a.closeChan:
		if !ok {
			return
		}

	default:
		close(a.closeChan)
	}
}

// GetWaitGroup returns the underlying wait group, allowing adjacent threads to adjust the application's wait group,
// allowing for clean thread syncrhonization on exit
func (a *Application) GetWaitGroup() *sync.WaitGroup {
	return a.wg
}

func (a *Application) SetFocus(obj Widget) *Application {
	a.inFocus = obj
	return a
}

func (a *Application) findFocusNext() Widget {
	focus := a.root.FindInFocus(a.root, nil, a.inFocus, nil)
	if focus.Next != nil {
		return focus.Next
	}

	return focus.First
}

func (a *Application) findFocusPrev() Widget {
	focus := a.root.FindInFocus(a.root, nil, a.inFocus, nil)
	if focus.Prev != nil {
		return focus.Prev
	}

	return focus.Last
}

// handleInput directs input to the widget with focus
func (a *Application) handleInput(input string) string {
	return a.inFocus.CaptureInput(input)
}

func (a *Application) renderAll() {
	terminal.HideCursor()
	a.root.Render(RenderModeAll, a.inFocus)
	terminal.ShowCursor()
}

func (a *Application) renderAllRefocus() {
	terminal.HideCursor()
	a.root.Render(RenderModeBorder, a.inFocus)
	terminal.ShowCursor()
}

func (a *Application) renderFocused() {
	terminal.HideCursor()
	a.inFocus.Render(RenderModeContent, a.inFocus)
	terminal.ShowCursor()
}

// Start starts the application loop
func (a *Application) Start() error {
	defer a.wg.Wait()

	// this exit will signal prior to the waitgroup wait function, allowing any blocking threads
	// an opportunity to gracefully exit
	defer a.Exit()
	// main application loop

	termFd, restore, err := terminal.Start()
	if err != nil {
		return err
	}
	a.termFd = termFd
	defer restore()

	if a.onExit != nil {
		defer a.onExit()
	}

	if a.onStart != nil {
		if err := a.onStart(); err != nil {
			return err
		}
	}

	inputChan := make(chan string, 100)
	byteBuf := make([]byte, 100)
	go func() {
		defer close(inputChan)
		for {
			nBytes, err := os.Stdin.Read(byteBuf)
			if err != nil {
				return
			}

			inputChan <- string(byteBuf[:nBytes])
		}
	}()

	termSizeChan := make(chan os.Signal, 100)
	signal.Notify(termSizeChan, syscall.SIGWINCH)

	width, height, err := term.GetSize(a.termFd)
	if err != nil {
		return err
	}
	a.root.SetDimensions(0, 0, width, height)
	a.root.Render(RenderModeAll, a.inFocus)

	// application event loop
	for {
		select {
		case <-a.sigChan:
			// process received a signal to exit
			return nil
		case <-a.closeChan:
			// application is being closed internally
			return nil
		case <-termSizeChan:
			// resize event occurred
			width, height, err := term.GetSize(a.termFd)
			if err != nil {
				return err
			}
			a.root.SetDimensions(0, 0, width, height)
			a.renderAll()
		case c := <-inputChan:
			switch c {
			case CtrlC:
				if a.ctrlCExit {
					return nil
				}
			case Tab:
				if a.inFocus != nil {
					a.SetFocus(a.findFocusNext())
					a.renderAllRefocus()
				}
			case ShiftTab:
				if a.inFocus != nil {
					a.SetFocus(a.findFocusPrev())
					a.renderAllRefocus()
				}
			default:
				if a.inputHandler != nil {
					c = a.inputHandler(c)
					if c == RenderFullCode {
						a.renderAll()
						continue
					}
				}

				c = a.handleInput(c)
				if c == RenderFullCode {
					a.renderAll()
				} else {
					// Render just the in-focus thing
					a.renderFocused()
				}
			}
		}
	}
}
