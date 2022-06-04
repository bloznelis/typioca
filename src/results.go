package main

import (
	"math"
	"strings"
)

func (m TimerBasedTest) calculateResults() Results {
	elapsedMinutes := m.timer.duration.Minutes()
	return Results{
		wpm:      int(m.base.calculateNormalizedWpm(elapsedMinutes)),
		accuracy: m.base.calculateAccuracy(),
		rawWpm:   int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:      m.base.calculateCpm(elapsedMinutes),
		time:     m.timer.duration,
		wordList: m.settings.wordListSelections[m.settings.wordListCursor].show,
	}
}

func (m WordCountBasedTest) calculateResults() Results {
	elapsedMinutes := m.stopwatch.stopwatch.Elapsed().Minutes()
	return Results{
		wpm:      int(m.base.calculateNormalizedWpm(elapsedMinutes)),
		accuracy: m.base.calculateAccuracy(),
		rawWpm:   int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:      m.base.calculateCpm(elapsedMinutes),
		time:     m.stopwatch.stopwatch.Elapsed(),
		wordList: m.settings.wordListSelections[m.settings.wordListCursor].show,
	}
}

func (m SentenceCountBasedTest) calculateResults() Results {
	elapsedMinutes := m.stopwatch.stopwatch.Elapsed().Minutes()
	return Results{
		wpm:      int(m.base.calculateNormalizedWpm(elapsedMinutes)),
		accuracy: m.base.calculateAccuracy(),
		rawWpm:   int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:      m.base.calculateCpm(elapsedMinutes),
		time:     m.stopwatch.stopwatch.Elapsed(),
		wordList: m.settings.sentenceListSelections[m.settings.sentenceListCursor].show,
	}
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
