package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/exadrift/go/ansi/style"
	"github.com/exadrift/go/tui"
)

var fruits = map[string][]string{
	"apple": {
		"delicious",
		"macintosh",
		"fuji",
		"gala",
		"honeycrisp",
		"pink lady",
		"cosmic crisp",
		"mila zagoras piliou",
		"cripps pink",
		"maçã bravo de esmolfe",
		"ambrosia",
		"firiki piliou",
	},
	"orange": {
		"navel",
		"valencia",
		"blood",
		"cara cara",
		"clementine",
	},
	"strawberry": {
		"albion",
		"jewel",
		"mara des bois",
		"allstar",
		"honeoye",
	},
	"watermelon": {
		"crimson sweet",
		"sugar baby",
		"moon and stars",
		"orangelo",
		"yellow baby",
		"black beauty",
		"8424",
	},
}

func getFruits() []string {
	var fs []string
	for f := range fruits {
		fs = append(fs, f)
	}

	return fs
}

func getFruitType(fruit string) []string {
	return fruits[fruit]
}

func main() {
	menu1 := tui.NewMenu(getFruits()...)
	menu1.SetTitle("fruit")

	menu2 := tui.NewMenu(getFruitType("apple")...)
	menu2.SetTitle("type")

	shell := tui.NewShell()
	shell.SetTitle("terminal")

	textbox := tui.NewText(style.B(
		"hello world, this is some text that's likely to need to wrap all through the box. let's make this so long that it runs over its max length and forces the need to scroll a bit.  vertically",
		"",
		"this is where the scrolling needs to happen.",
		"hopefully these newlines will accelerate the process.",
	))
	textbox.SetTitle("text")

	topBar := tui.NewText(style.B(style.T("example program is the best and this is too")))
	topBar.SetFocusable(false)

	horizTestText1 := tui.NewText(style.B("how many things?", "maybe many?", "lots of them??", "so many things in this box?"))
	horizTestText2 := tui.NewText(style.B(style.S("how many things?", style.Yellow.Fg()), "maybe many?"))
	horizTestText2.SetTitle("horiz")

	finalHorizLayout := tui.NewFlexLayout(
		tui.OrientationHorizontal,
		1,
		tui.NewSegment(1, horizTestText1),
		tui.NewSegment(1, horizTestText2),
	)

	shellLayout := tui.NewFlexLayout(
		tui.OrientationVertical,
		1,
		tui.NewSegment(3, shell),
		tui.NewSegment(1, finalHorizLayout),
	)

	focusLayout := tui.NewFlexLayout(
		tui.OrientationHorizontal,
		1,
		tui.NewSegment(1, menu1),
		tui.NewSegment(1, tui.NewFlexLayout(
			tui.OrientationVertical,
			1,
			tui.NewSegment(1, menu2),
			tui.NewSegment(1, textbox),
		)),
		tui.NewSegment(3, shellLayout),
	)

	layout := tui.NewFlexLayout(
		tui.OrientationVertical,
		0,
		tui.NewSegment(1, topBar, tui.WithSegmentOptionMinChars(1)),
		tui.NewSegment(1000, focusLayout),
	)

	app := tui.New(layout).SetFocus(menu1)

	c := exec.Command("/bin/bash")
	if err := shell.Start(app, c); err != nil {
		log.Fatal(err)
	}

	menu1.SetSelectHandler(
		func(selectedIndex int, selectedItem string) any {
			time.Sleep(time.Second * 4)
			return getFruitType(selectedItem)
		},
		tui.WithBusyModal("loading fruits...",
			func(a any) {
				loadedFruits := a.([]string)
				menu2.SetContents(loadedFruits...)
				app.SetFocus(menu2)
			},
		),
	)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
