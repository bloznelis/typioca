package main

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/muesli/termenv"
)

type myTimer struct {
	timer     timer.Model
	duration  time.Duration
	isRunning bool // Inner is running is being handled weirdly.
	timedout  bool
}

type myStopWatch struct {
	stopwatch stopwatch.Model
	isRunning bool
}

type mistakes struct {
	mistakesAt     map[int]bool
	rawMistakesCnt int // Should never be reduced
}

type StringStyle func(string) termenv.Style

type styles struct {
	correct      StringStyle
	toEnter      StringStyle
	mistakes     StringStyle
	cursor       StringStyle
	runningTimer StringStyle
	stoppedTimer StringStyle
	greener      StringStyle
	magenta      StringStyle
}

type model struct {
	state  State
	styles styles
}

type Results struct {
	wpm      int
	accuracy float64
	rawWpm   int
	cpm      int
	time     time.Duration
	wordList string
}

type State interface{}

type MainMenuSelection interface {
	show(s styles) string
}

type TimerBasedTestSettings struct {
	timeSelections     []time.Duration
	timeCursor         int
	wordListSelections []string
	wordListCursor     int
	cursor             int
}

type WordCountBasedTestSettings struct {
	wordCountSelections []int
	wordCountCursor     int
	wordListSelections  []string
	wordListCursor      int
	cursor              int
}

type MainMenu struct {
	choices []MainMenuSelection
	cursor  int
}

type TimerBasedTest struct {
	settings     TimerBasedTestSettings
	timer        myTimer
	wordsToEnter string
	inputBuffer  []rune
	rawInputCnt  int // Should not be reduced
	mistakes     mistakes
	completed    bool
	cursor       int
}

type TimerBasedTestResults struct {
	settings TimerBasedTestSettings
	results  Results
}

type WordCountBasedTest struct {
	settings     WordCountBasedTestSettings
	stopwatch    myStopWatch
	wordsToEnter string
	inputBuffer  []rune
	rawInputCnt  int // Should not be reduced
	mistakes     mistakes
	completed    bool
	cursor       int
}

type WordCountTestResults struct {
	settings WordCountBasedTestSettings
	wordCnt  int
	results  Results
}
