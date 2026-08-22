package main

import (
	"log"
	"os/exec"
	"time"

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
	menu1.EnableBorder(true).SetTitle("fruit")

	menu2 := tui.NewMenu(getFruitType("apple")...)
	menu2.EnableBorder(true).SetTitle("type")

	shell := tui.NewShell()
	shell.EnableBorder(true).SetTitle("terminal")

	textbox := tui.NewText("hello world, this is some text that's likely to need to wrap all through the box. let's make this so long that it runs over its max length and forces the need to scroll a bit.  vertically\n\nthis is where the scrolling needs to happen.\nhopefully these newlines will accelerate the process.")
	textbox.EnableBorder(true).SetTitle("text")

	layout := tui.NewFlexLayout(
		tui.OrientationHorizontal,
		tui.NewSegment(1, menu1),
		tui.NewSegment(1, tui.NewFlexLayout(
			tui.OrientationVertical,
			tui.NewSegment(1, menu2),
			tui.NewSegment(1, textbox),
		)),
		tui.NewSegment(3, shell),
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
