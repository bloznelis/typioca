package words

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

type Metadata struct {
	Name       string
	Size       int
	PackagedAt string //Use some kind of date type here?
	Version    int
}

type WordSource struct {
	Metadata Metadata
	Words    []string
}

//go:embed embedables/words/common-english.json
var commonEnglish string

//go:embed embedables/words/dorian-gray.json
var dorianGray string

//go:embed embedables/words/frankenstein.json
var frankenstein string

//go:embed embedables/words/pride-and-prejudice.json
var prideAndPrejudice string

//go:embed embedables/sentences/dorian-gray.json
var dorianGraySentences string

//go:embed embedables/sentences/pride-and-prejudice.json
var prideAndPrejudiceSentences string

//go:embed embedables/sentences/frankenstein.json
var frankensteinSentences string

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

type WordsGenerator struct {
	Count     int
	pools     map[string]string
	poolsJson map[string]WordSource
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func unmarshalSources(sources []string) map[string]WordSource {
	acc := make(map[string]WordSource, len(sources))
	for _, source := range sources {
		wordSource := unmarshalSource(source)
		acc[wordSource.Metadata.Name] = wordSource
	}

	return acc
}

func unmarshalSource(sourceRaw string) WordSource {
	var wordSource WordSource
	err := json.Unmarshal([]byte(sourceRaw), &wordSource)

	check(err)
	return wordSource
}

func NewGenerator() (g WordsGenerator) {
	g.Count = 300
	g.poolsJson = unmarshalSources([]string{
		commonEnglish,
		dorianGray,
		dorianGraySentences,
		frankenstein,
		frankensteinSentences,
		prideAndPrejudice,
		prideAndPrejudiceSentences,
	})
	return g
}

func (this WordsGenerator) Generate(listName string) []rune {
	pool := this.poolsJson[listName]
	acc := []string{}
	poolLength := pool.Metadata.Size
	for i := 0; i < this.Count; i++ {
		word := pool.Words[rand.Int()%poolLength]
		acc = append(acc, word)
	}

	return []rune(strings.Join(acc, " "))
}
