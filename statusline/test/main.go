package main

import (
	"fmt"
	"time"

	"github.com/exadrift/go/statusline"
)

func main() {
	s := statusline.NewStatusLine(statusline.WithSpinnerStyling("\x1b[92m"))
	defer s.Stop()

	s.Status("calculating things...")
	time.Sleep(1 * time.Second)
	s.Status("calculating...")
	time.Sleep(1 * time.Second)
	s.Emit("1 + 1 is 2")
	s.Status("verifying...")
	time.Sleep(2 * time.Second)
	s.Status("confirming...")
	time.Sleep(2 * time.Second)
	s.Emit("it's definitely the case")
	for i := 3; i > 0; i-- {
		s.Status("counting down to input %d", i)
		time.Sleep(time.Second)
	}
	s.Status("about to prompt")
	time.Sleep(time.Second)
	message := s.Prompt(false, "enter your message >")
	time.Sleep(2 * time.Second)
	s.Emit("you said \"%s\"", message)
	isConfirmed := s.Confirm(false, "confirm the thing")
	fmt.Println(isConfirmed)
}
