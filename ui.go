package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

const Red string = "\033[31m"
const Green string = "\033[32m"
const Yellow string = "\033[33m"
const ColorReset string = "\033[0m"

const GUESSES string = "guesses"
const STATUS string = "status"

type Game struct {
	User     uint64
	Word     string
	Input    string
	Guesses  []string
	Finished bool
}

func setupKeybinds(g *gocui.Gui) {
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, reset)
	g.SetKeybinding("", gocui.KeyBackspace2, gocui.ModNone, backSpace)
	g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, enter)
	bindLetters(g)
}

func backSpace(_ *gocui.Gui, v *gocui.View) error {
	length := len(game.Input)
	if length > 0 {
		game.Input = game.Input[0 : length-1]
		updateGuessView(v)
	}
	return nil
}

func enter(g *gocui.Gui, v *gocui.View) error {
	if !game.Finished && len(game.Guesses) < MAX_GUESS && len(game.Input) == len(game.Word) {
		if game.Input == game.Word {
			game.Finished = true
			setEndGame(true, g)
		}
		game.Guesses = append(game.Guesses, colorWord(game.Input, game.Word))
		game.Input = ""
		if len(game.Guesses) == MAX_GUESS {
			setEndGame(false, g)
		}
		return updateGuessView(v)
	}
	return nil
}

func setEndGame(isWinner bool, g *gocui.Gui) {
	v, err := g.View(STATUS)
	if err != nil {
		panic(err)
	}
	if isWinner {
		updateStatus(v, "WINNER")
	} else {
		updateStatus(v, Red+game.Word+ColorReset)
	}
}

func readChar(v *gocui.View, c rune) error {
	if !game.Finished && len(game.Input) < len(game.Word) {
		game.Input += string(c)
		return updateGuessView(v)
	}
	return nil
}

func updateGuessView(v *gocui.View) error {
	v.Clear()
	for _, g := range game.Guesses {
		fmt.Fprintln(v, g)
	}
	fmt.Fprint(v, game.Input)
	for i := 0; i < len(game.Word)-len(game.Input); i++ {
		fmt.Fprint(v, "_")
	}
	fmt.Fprintln(v)
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(GUESSES, maxX/2-3, maxY/2, maxX/2+3, maxY/2+8, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		if _, err := g.SetCurrentView(GUESSES); err != nil {
			return err
		}
		updateGuessView(v)
	}
	if v, err := g.SetView(STATUS, maxX/2-5, maxY/2-4, maxX/2+5, maxY/2-2, 0); err != nil {
		updateStatus(v, "WORDY")
	}

	return nil
}

func updateStatus(v *gocui.View, status string) error {
	v.Clear()
	width, _ := v.Size()
	padding := (width - len(status)) / 2
	for i := 0; i < padding; i++ {
		fmt.Fprint(v, " ")
	}
	fmt.Fprintln(v, status)
	return nil
}

func bindLetters(g *gocui.Gui) {
	for _, c := range "abcdefghijklmnopqrstuvwxyz" {
		if err := g.SetKeybinding("", c, gocui.ModNone, keyHandler(c)); err != nil {
			log.Panicln(err)
		}
	}
}

func keyHandler(c rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return readChar(v, c)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func reset(g *gocui.Gui, _ *gocui.View) error {
	game = newGame(words)
	v, _ := g.View(GUESSES)
	updateGuessView(v)
	v, _ = g.View(STATUS)
	updateStatus(v, "WORDY")
	return nil
}
