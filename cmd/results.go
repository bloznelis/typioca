package cmd

import (
	"math"
	"strings"
)

func (m TimerBasedTest) calculateResults() Results {
	wordlist := m.settings.wordListSelections[m.settings.wordListCursor].name
	identifier := ResultsIdentifier{
		testType: "TimerBasedTest",
		numeric:  int(m.timer.duration),
		words:    wordlist,
	}

	elapsedMinutes := m.timer.duration.Minutes()
	wpm := m.base.calculateNormalizedWpm(elapsedMinutes)
	deltaWpm := calculateAverageWpmDeltaPercentage(wpm, ReadResults(identifier))

	return Results{
		identifier:    identifier,
		wpm:           int(wpm),
		accuracy:      m.base.calculateAccuracy(),
		deltaWpm:      deltaWpm,
		rawWpm:        int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:           m.base.calculateCpm(elapsedMinutes),
		time:          m.timer.duration,
		wordList:      wordlist,
		wpmEachSecond: m.base.wpmEachSecond,
	}
}

func (m WordCountBasedTest) calculateResults() Results {
	count := m.settings.wordCountSelections[m.settings.wordCountCursor]
	wordlist := m.settings.wordListSelections[m.settings.wordListCursor].name

	identifier := ResultsIdentifier{
		testType: "WordCountBasedTest",
		numeric:  count,
		words:    wordlist,
	}

	elapsedMinutes := m.stopwatch.stopwatch.Elapsed().Minutes()
	wpm := m.base.calculateNormalizedWpm(elapsedMinutes)
	deltaWpm := calculateAverageWpmDeltaPercentage(wpm, ReadResults(identifier))

	return Results{
		identifier:    identifier,
		wpm:           int(wpm),
		accuracy:      m.base.calculateAccuracy(),
		deltaWpm:      deltaWpm,
		rawWpm:        int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:           m.base.calculateCpm(elapsedMinutes),
		time:          m.stopwatch.stopwatch.Elapsed(),
		wordList:      wordlist,
		wpmEachSecond: m.base.wpmEachSecond,
	}
}

func (m SentenceCountBasedTest) calculateResults() Results {
	count := m.settings.sentenceCountSelections[m.settings.sentenceCountCursor]
	wordlist := m.settings.sentenceListSelections[m.settings.sentenceListCursor].name

	identifier := ResultsIdentifier{
		testType: "SentenceCountBasedTest",
		numeric:  count,
		words:    wordlist,
	}

	elapsedMinutes := m.stopwatch.stopwatch.Elapsed().Minutes()
	wpm := m.base.calculateNormalizedWpm(elapsedMinutes)
	deltaWpm := calculateAverageWpmDeltaPercentage(wpm, ReadResults(identifier))

	return Results{
		identifier:    identifier,
		wpm:           int(wpm),
		accuracy:      m.base.calculateAccuracy(),
		deltaWpm:      deltaWpm,
		rawWpm:        int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:           m.base.calculateCpm(elapsedMinutes),
		time:          m.stopwatch.stopwatch.Elapsed(),
		wordList:      wordlist,
		wpmEachSecond: m.base.wpmEachSecond,
	}
}

func calculateAverageWpmDeltaPercentage(wpm float64, previousResults []PersistentResultsNode) float64 {
	previousAvg := calcPreviousResultsAvgWpm(previousResults)

	return ((wpm - previousAvg) / math.Max(1.0, previousAvg)) * 100
}

func calcPreviousResultsAvgWpm(previousResults []PersistentResultsNode) float64 {
	if len(previousResults) == 0 {
		return 0
	}
	var sum int
	for _, v := range previousResults {
		sum += v.Wpm
	}

	return float64(sum) / float64(len(previousResults))
}

func (base TestBase) calculateNormalizedWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(base.inputBuffer)/5, elapsedMinutes)
}

func (base TestBase) calculateRawWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(strings.Split(string(base.inputBuffer), " ")), elapsedMinutes)
}

func (base TestBase) calculateWpm(wordCnt int, elapsedMinutes float64) float64 {
	if elapsedMinutes == 0 {
		return 0
	} else {
		grossWpm := float64(wordCnt) / elapsedMinutes
		netWpm := grossWpm - float64(len(base.mistakes.mistakesAt))/elapsedMinutes

		return math.Max(0, netWpm)
	}
}

func (base TestBase) calculateCpm(elapsedMinutes float64) int {
	return int(float64(base.rawInputCnt) / elapsedMinutes)
}

func (base TestBase) calculateAccuracy() float64 {
	mistakesRate := float64(base.mistakes.rawMistakesCnt*100) / float64(base.rawInputCnt)
	accuracy := 100 - mistakesRate
	return accuracy
}
