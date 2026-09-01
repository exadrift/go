package main

import (
	"fmt"
	"log"
	"os"

	"github.com/exadrift/go/pshell"
	"golang.org/x/term"
)

func main() {
	ps := pshell.New(pshell.WithPasswordCallback(func() (string, error) {
		fmt.Printf("give me password input > ")
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(pass), nil
	}))

	output, err := ps.ExecuteCommands("ls -la\nsudo ls .\nls /\nsudo ls .\nls /")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
