package main

import (
	_ "embed"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

//go:embed embedables/words/common-english.txt
var commonEnglish string

//go:embed embedables/words/dorian-gray.txt
var dorianGray string

//go:embed embedables/words/frankenstein.txt
var frankenstein string

//go:embed embedables/words/pride-and-prejudice.txt
var prideAndPrejudice string

//go:embed embedables/sentences/dorian-gray.txt
var dorianGraySentences string

//go:embed embedables/sentences/pride-and-prejudice.txt
var prideAndPrejudiceSentences string

//go:embed embedables/sentences/frankenstein.txt
var frankensteinSentences string

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

type WordsGenerator struct {
	Count int
	pools map[string]string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makePool(content string) []string {
	words := strings.Split(content, "\n")

	return words
}

func NewGenerator() (g WordsGenerator) {
	g.Count = 300
	g.pools = map[string]string{
		"common-words":                  commonEnglish,
		"dorian-gray-words":             dorianGray,
		"dorian-gray-sentences":         dorianGraySentences,
		"frankenstein-words":            frankenstein,
		"frankenstein-sentences":        frankensteinSentences,
		"pride-and-prejudice-words":     prideAndPrejudice,
		"pride-and-prejudice-sentences": prideAndPrejudiceSentences,
	}
	return g
}

func (this WordsGenerator) Generate(poolKey string) string {
	pool := makePool(this.pools[poolKey])
	acc := []string{}
	poolLength := len(pool)
	for i := 0; i < this.Count; i++ {
		word := pool[rand.Int()%poolLength]
		word = regexp.MustCompile("\r|\n").ReplaceAllString(word, "")
		acc = append(acc, word)
	}

	return strings.Join(acc, " ")
}
