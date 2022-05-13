package main

import "strings"

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
	return calculateWpm(m, len(m.inputBuffer)/5)
}

func calculateRawWpm(m TimerBasedTest) int {
	return calculateWpm(m, len(strings.Split(string(m.inputBuffer), " ")))
}

func calculateWpm(m TimerBasedTest, wordCnt int) int {
	grossWpm := float64(wordCnt) / m.timer.duration.Minutes()
	netWpm := grossWpm - float64(len(m.mistakes.mistakesAt))/m.timer.duration.Minutes()

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
	return m.calculateWpm(len(m.inputBuffer) / 5)
}

func (m WordCountBasedTest) calculateRawWpm() int {
	return m.calculateWpm(len(strings.Split(string(m.inputBuffer), " ")))
}

func (m WordCountBasedTest) calculateWpm(wordCnt int) int {
	grossWpm := float64(wordCnt) / m.stopwatch.stopwatch.Elapsed().Minutes()
	netWpm := grossWpm - float64(len(m.mistakes.mistakesAt))/m.stopwatch.stopwatch.Elapsed().Minutes()

	return int(netWpm)
}

func (m WordCountBasedTest) calculateCpm() int {
	return int(float64(m.rawInputCnt) / m.stopwatch.stopwatch.Elapsed().Minutes())
}

func (m WordCountBasedTest) calculateAccuracy() float64 {
	mistakesRate := float64(m.mistakes.rawMistakesCnt*100) / float64(m.rawInputCnt)
	accuracy := 100 - mistakesRate
	return accuracy
}
