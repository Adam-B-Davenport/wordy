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

const MAX_GUESS int = 6

var words []string
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

func newGame(words []string) Game {
	word := randomWord(words)
	return Game{
		User:     0,
		Word:     word,
		Input:    "",
		Guesses:  make([]string, 0),
		Finished: false,
	}
}

func setupGame() {
	wordLength := 5
	rand.Seed(time.Now().Unix())
	words = buildWordList(wordLength)
	game = newGame(words)
}

func main() {
	setupGame()

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	setupKeybinds(g)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)

	}
}
