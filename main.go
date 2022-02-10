package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func readDict() (string, error) {
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
	dict, err := readDict()
	if err == nil {
		words := make([]string, 0)
		for _, w := range strings.Split(dict, "\n") {
			if len(w) == length && isAlpha(w) {
				words = append(words, w)
			}
		}
		return words
	} else {
		panic(err)
	}
}

func randomWord(words []string) string {
	i := rand.Int() % len(words)
	return words[i]
}

func main() {
	rand.Seed(time.Now().Unix())
	words := buildWordList(5)
	fmt.Println(randomWord(words))
}
