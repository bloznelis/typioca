package main

import (
	_ "embed"
	"math/rand"
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

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

type WordsGenerator struct {
	Count int
	Pools []string
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
	g.Pools = []string{commonEnglish, dorianGray, frankenstein, prideAndPrejudice}
	return g
}

func (this WordsGenerator) Generate() string {
	randomIndex := rand.Intn(len(this.Pools))
	pool := makePool(this.Pools[randomIndex])
	acc := []string{}
	poolLength := len(pool)
	for i := 0; i < this.Count; i++ {
		acc = append(acc, pool[rand.Int()%poolLength])
	}

	return strings.Join(acc, " ")
}
