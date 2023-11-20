package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bloznelis/typioca/cmd/words"
)

type TestType = string
type NumericSetting = int
type WordListName = string
type AllPersistedResults = map[TestType]map[NumericSetting]map[WordListName][]PersistentResultsNode

func PersistResults(results Results) PersistentResults {
	var resultsFile = getResultsPath()
	var persistentResults PersistentResults

	//File does not exist?
	if _, err := os.Stat(resultsFile); os.IsNotExist(err) {
		persistentResults = defaultPersistentResults()
	} else {
		readResults(&persistentResults)
		// XXX: Once needed, version check should happen here
	}

	persistentResults.addResults(results)

	writeResults(persistentResults)

	return persistentResults
}

func defaultPersistentResults() PersistentResults {
	return PersistentResults{
		Results: AllPersistedResults{},
		Version: 1,
	}
}

func ReadResults(i ResultsIdentifier) []PersistentResultsNode {
	var resultsFile = getResultsPath()
	var persistentResults PersistentResults

	//File does not exist?
	if _, err := os.Stat(resultsFile); os.IsNotExist(err) {
		persistentResults = defaultPersistentResults()
	} else {
		readResults(&persistentResults)
		// XXX: Once needed, version check should happen here
	}
	var res = persistentResults.Results[i.testType][i.numeric][i.words]
	if res == nil {
		return make([]PersistentResultsNode, 0)
	}
	return res
}

func readResults(results *PersistentResults) {
	var resultsFilePath = getResultsPath()
	fh, err := os.Open(resultsFilePath)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	decoder.Decode(results)
}

func (p *PersistentResults) addResults(results Results) {
	var limit = 25 // XXX: Should be configurable eventually

	var node = PersistentResultsNode{
		Wpm:           results.wpm,
		Accuracy:      results.accuracy,
		DeltaWpm:      results.deltaWpm,
		RawWpm:        results.rawWpm,
		Cpm:           results.cpm,
		WpmEachSecond: results.wpmEachSecond,
	}
	var i = results.identifier

	if p.Results[i.testType] == nil {
		p.Results[i.testType] = map[NumericSetting]map[WordListName][]PersistentResultsNode{}
	}

	if p.Results[i.testType][i.numeric] == nil {
		p.Results[i.testType][i.numeric] = map[WordListName][]PersistentResultsNode{}
	}

	var nodes = append(p.Results[i.testType][i.numeric][i.words], node)

	if len(nodes) > limit {
		nodes = nodes[len(nodes)-limit:]
	}

	p.Results[i.testType][i.numeric][i.words] = nodes
}

func writeResults(results PersistentResults) {
	var resultsFilePath = getResultsPath()
	words.EnsureDir(resultsFilePath)
	fh, err := os.Create(resultsFilePath)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	encoder.SetIndent("", "\t")
	encoder.Encode(results)
}

func getResultsPath() string {
	var cachePath = getCachePath()
	var resultsFilePath = filepath.Join(cachePath, "results.json")
	return resultsFilePath
}
