package main

import (
	"math"
	"strings"
)

func (m TimerBasedTest) calculateResults() Results {
	return Results{
		wpm:      calculateNormalizedWpm(m),
		accuracy: calculateAccuracy(m),
		rawWpm:   calculateRawWpm(m),
		cpm:      calculateCpm(m),
		time:     m.timer.duration,
		wordList: m.settings.wordListSelections[m.settings.wordListCursor],
	}
}

func calculateNormalizedWpm(m TimerBasedTest) int {
	return calculateWpm(m, len(m.base.inputBuffer)/5)
}

func calculateRawWpm(m TimerBasedTest) int {
	return calculateWpm(m, len(strings.Split(string(m.base.inputBuffer), " ")))
}

func calculateWpm(m TimerBasedTest, wordCnt int) int {
	grossWpm := float64(wordCnt) / m.timer.duration.Minutes()
	netWpm := grossWpm - float64(len(m.base.mistakes.mistakesAt))/m.timer.duration.Minutes()

	return int(math.Max(0, netWpm))
}

func calculateCpm(m TimerBasedTest) int {
	return int(float64(m.base.rawInputCnt) / m.timer.duration.Minutes())
}

func calculateAccuracy(m TimerBasedTest) float64 {
	mistakesRate := float64(m.base.mistakes.rawMistakesCnt*100) / float64(m.base.rawInputCnt)
	accuracy := 100 - mistakesRate
	return accuracy
}
func (m WordCountBasedTest) calculateResults() Results {
	return Results{
		wpm:      m.calculateNormalizedWpm(),
		accuracy: m.calculateAccuracy(),
		rawWpm:   m.calculateRawWpm(),
		cpm:      m.calculateCpm(),
		time:     m.stopwatch.stopwatch.Elapsed(),
		wordList: m.settings.wordListSelections[m.settings.wordListCursor],
	}
}

func (m WordCountBasedTest) calculateNormalizedWpm() int {
	return m.calculateWpm(len(m.base.inputBuffer) / 5)
}

func (m WordCountBasedTest) calculateRawWpm() int {
	return m.calculateWpm(len(strings.Split(string(m.base.inputBuffer), " ")))
}

func (m WordCountBasedTest) calculateWpm(wordCnt int) int {
	grossWpm := float64(wordCnt) / m.stopwatch.stopwatch.Elapsed().Minutes()
	netWpm := grossWpm - float64(len(m.base.mistakes.mistakesAt))/m.stopwatch.stopwatch.Elapsed().Minutes()

	return int(math.Max(0, netWpm))
}

func (m WordCountBasedTest) calculateCpm() int {
	return int(float64(m.base.rawInputCnt) / m.stopwatch.stopwatch.Elapsed().Minutes())
}

func (m WordCountBasedTest) calculateAccuracy() float64 {
	mistakesRate := float64(m.base.mistakes.rawMistakesCnt*100) / float64(m.base.rawInputCnt)
	accuracy := 100 - mistakesRate
	return accuracy
}
