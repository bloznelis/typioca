package main

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

type WordsGenerator struct {
	Count int
	Pool  []string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makePool(file string) []string {
	data, err := os.ReadFile(file)
	check(err)

	words := strings.Split(string(data), "\n")

	return words
}

func randomWordsPool() []string {
	return makePool("words/lists/dorian-gray.txt")
}

func NewGenerator() (g WordsGenerator) {
	g.Count = 30
	g.Pool = randomWordsPool()
	return g
}

func (this WordsGenerator) Generate() string {
	acc := []string{}
	poolLength := len(this.Pool)
	for i := 0; i < this.Count; i++ {
		acc = append(acc, this.Pool[rand.Int()%poolLength])
	}

	return strings.Join(acc, " ")
}
