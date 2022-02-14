package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/awesome-gocui/gocui"
)

type Game struct {
	User    uint64
	Word    string
	Input   string
	Guesses []string
}

const GUESSES string = "guesses"

var game Game

func readWordList() (string, error) {
	b, err := ioutil.ReadFile("/usr/share/dict/words")
	if err == nil {
		return string(b), nil
	}
	return "", err
}

func isAlpha(word string) bool {
	for _, r := range word {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func buildWordList(length int) []string {
	wordList, err := readWordList()
	if err != nil {
		panic(err)
	}

	words := make([]string, 0)
	for _, w := range strings.Split(wordList, "\n") {
		if len(w) == length && isAlpha(w) {
			words = append(words, w)
		}
	}
	return words
}

func randomWord(words []string) string {
	i := rand.Int() % len(words)
	return words[i]
}

func readInput() (string, error) {
	var word string
	_, err := fmt.Scan(&word)
	return word, err
}

const Green string = "\033[32m"
const Yellow string = "\033[33m"
const ColorReset string = "\033[0m"

func colorPrint(text string, color string) {
	fmt.Print(color, text, ColorReset)
}

func colorWord(guess string, word string) string {
	result := ""
	for i := 0; i < len(guess); i++ {
		if guess[i] == word[i] {
			result += Green + guess[i:i+1] + ColorReset
		} else if strings.Contains(word, guess[i:i+1]) {
			result += Yellow + guess[i:i+1] + ColorReset
		} else {
			result += guess[i : i+1]
		}
	}
	return result
}

func gameLoop() {
	wordLength := 5
	rand.Seed(time.Now().Unix())
	words := buildWordList(wordLength)
	word := randomWord(words)
	//fmt.Println(word)
	for i := 0; i < 6; i++ {
		fmt.Print("Guess: ")
		var guess string
		var err error
		for len(guess) != len(word) {
			guess, err = readInput()
			if err != nil {
				panic(err)
			}
		}
		if guess == word {
			fmt.Println("Correct")
		}
		fmt.Println(colorWord(guess, word))
	}
	fmt.Println(word)
}

func backSpace(_ *gocui.Gui, v *gocui.View) error {
	length := len(game.Input)
	if length > 0 {
		game.Input = game.Input[0 : length-1]
		updateGuessView(v)

	}
	return nil

}

func readChar(v *gocui.View, c rune) error {
	if len(game.Input) < len(game.Word) {
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
	fmt.Fprintln(v, game.Input)
	return nil
}

func main() {
	wordLength := 5
	rand.Seed(time.Now().Unix())
	words := buildWordList(wordLength)
	word := randomWord(words)
	game = Game{
		User:    0,
		Word:    word,
		Input:   "",
		Guesses: make([]string, 0),
	}
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	g.SetKeybinding("", gocui.KeyBackspace2, gocui.ModNone, backSpace)
	bindLetters(g)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)

	}
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
		fmt.Fprintln(v, game.Word)

	}

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
