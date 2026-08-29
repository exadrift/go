package main

import (
	"time"

	"github.com/exadrift/go/statusline"
)

func main() {
	s := statusline.NewStatusLine(statusline.WithSpinnerStyling("\x1b[92m"))
	defer s.Stop()

	s.Status("calculating...")
	time.Sleep(3 * time.Second)
	s.Emit("1 + 1 is 2")
	s.Status("verifying...")
	time.Sleep(2 * time.Second)
	s.Status("confirming...")
	time.Sleep(5 * time.Second)
	s.Emit("it's definitely the case")
}
