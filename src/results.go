package main

import "strings"

func (m TimerBasedTest) calculateResults() Results {
	return Results{
		wpm:      calculateNormalizedWpm(m),
		accuracy: calculateAccuracy(m),
		rawWpm:   calculateRawWpm(m),
		cpm:      calculateCpm(m),
		time:     m.timer.duration,
	}
}

func calculateNormalizedWpm(m TimerBasedTest) int {
	return calculateWpm(m, len(m.inputBuffer)/5)
}

func calculateRawWpm(m TimerBasedTest) int {
	return calculateWpm(m, len(strings.Split(string(m.inputBuffer), " ")))
}

func calculateWpm(m TimerBasedTest, wordCnt int) int {
	grossWpm := float64(wordCnt) / m.timer.duration.Minutes()
	netWpm := grossWpm - (float64(len(m.mistakes.mistakesAt)) / m.timer.duration.Minutes())

	return int(netWpm)
}

func calculateCpm(m TimerBasedTest) int {
	return int(float64(m.rawInputCnt) / m.timer.duration.Minutes())
}

func calculateAccuracy(m TimerBasedTest) float64 {
	mistakesRate := float64(m.mistakes.rawMistakesCnt*100) / float64(m.rawInputCnt)
	accuracy := 100 - mistakesRate
	return accuracy
}