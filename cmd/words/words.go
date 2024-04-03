package words

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"math/rand"
	"os"
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

func addEmbededSource(sources map[string]WordSource) map[string]WordSource {
	var wordSource WordSource
	err := json.Unmarshal([]byte(commonEnglish), &wordSource)
	check(err)
	var sentenceSource WordSource
	err = json.Unmarshal([]byte(frankensteinSentences), &sentenceSource)
	check(err)

	sources["Common words"] = wordSource
	sources["Frankenstein sentences"] = sentenceSource

	return sources
}

func unmarshalSources(paths []string) map[string]WordSource {
	acc := make(map[string]WordSource, len(paths))
	for _, sourceFilePath := range paths {
		var wordSource WordSource
		if strings.HasSuffix(sourceFilePath, ".json") {
			wordSource = readJsonSource(sourceFilePath)
		} else {
			wordSource = readNewLineSource(sourceFilePath)
		}

		acc[sourceFilePath] = wordSource
	}

	return acc
}

func readJsonSource(sourceFilePath string) WordSource {
	var wordSource WordSource

	fh, err := os.Open(sourceFilePath)
	defer fh.Close()
	check(err)

	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&wordSource)
	check(err)

	return wordSource
}

func readNewLineSource(sourceFilePath string) WordSource {
	fh, err := os.Open(sourceFilePath)
	defer fh.Close()
	check(err)

	var lines []string
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	metadata := Metadata{
		Name:       fh.Name(),
		Size:       len(lines),
		PackagedAt: "1970-01-01T00:00:00Z",
		Version:    1,
	}

	return WordSource{
		Metadata: metadata,
		Words:    lines,
	}
}

func NewGenerator(paths []string) (g WordsGenerator) {
	g.Count = 300
	g.poolsJson = unmarshalSources(paths)
	g.poolsJson = addEmbededSource(g.poolsJson)

	return g
}

func (this WordsGenerator) Generate(listName string) []rune {
	pool := this.poolsJson[listName].Words

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	takeAmount := min(this.Count, len(pool))
	words := pool[0:takeAmount]

	return []rune(strings.Join(words, " "))
}
