package main

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/muesli/termenv"
)

type myTimer struct {
	timer     timer.Model
	duration  time.Duration
	isRunning bool // Inner is running is being handled weirdly.
	timedout  bool
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
}

type model struct {
	state        StateType
	styles       styles
	timer        myTimer
	wordsToEnter string
	inputBuffer  []rune
	rawInputCnt  int // Should not be reduced
	mistakes     mistakes
	completed    bool
	cursor       int
}

type Results struct {
	wpm      int
	accuracy float64
	rawWpm   int
	cpm      int
	time     time.Duration
}

type StateType interface{}
type TimerBasedTest struct{}
type TimerBasedTestResults struct {
	results Results
}
