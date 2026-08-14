package main

import (
	"fmt"
	"os"

	"github.com/exadrift/go/ansi"
	"golang.org/x/term"
)

func main() {
	// Enter raw mode and save the previous terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	// Ensure cleanup runs when main exits
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error restoring terminal: %v\n", err)
		}
	}()

	fmt.Print("Test keystrokes here. Press 'q' to quit.\r\n")

	byteBuf := make([]byte, 100)
	for {
		nBytes, err := os.Stdin.Read(byteBuf)
		if err != nil {
			return
		}

		seq := byteBuf[:nBytes]
		if seq[0] == 'q' {
			return
		}

		ansiCode := ansi.EscapedAnsiString(string(seq))
		var humanName string
		var keyCodeAnsi string
		keyCombo, err := ansi.ParseAnsiCode(string(seq))
		if err != nil {
			humanName = ""
			keyCodeAnsi = ""
		} else {
			humanName = keyCombo.HumanName
			keyCodeAnsi = keyCombo.EscapedAnsiString()
		}
		fmt.Printf("%s    %s    %s\r\n", ansiCode, humanName, keyCodeAnsi)
	}
}
