package tui

import (
	"fmt"
	"regexp"
)

var ansiEscStripper = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

const (
	Black         int = 30
	Red           int = 31
	Green         int = 32
	Yellow        int = 33
	Blue          int = 34
	Magenta       int = 35
	Cyan          int = 36
	White         int = 37
	BlackBright   int = 90
	RedBright     int = 91
	GreenBright   int = 92
	YellowBright  int = 93
	BlueBright    int = 94
	MagentaBright int = 95
	CyanBright    int = 96
	WhiteBright   int = 97
)

const (
	UpArrow   = "\x1b[A"
	DownArrow = "\x1b[B"
	Enter     = "\r"
	CtrlC     = string(rune(3))
	Tab       = "\x09"
	ShiftTab  = "\x1b[Z"

	RenderFullCode = "\x1b[RENDER"
)

func StyleFg(color int, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
}

func StyleFgBg(fgColor int, bgColor int, text string) string {
	return fmt.Sprintf("\x1b[%d;%dm%s\x1b[0m", fgColor, bgColor+10, text)
}

func StripAnsiCodes(text string) string {
	return ansiEscStripper.ReplaceAllString(text, "")
}
