package main

import (
	"fmt"
	"os"
	"time"

	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	termenv.ClearScreen()
	termenv.SetWindowTitle("typioca")

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	termenv.Reset()
	println("bye!")
}

func initialTimerBasedTest() TimerBasedTest {
	generator := NewGenerator()
	generator.Count = 300

	testDuration := time.Second * 30

	textToEnter := generator.Generate()

	return TimerBasedTest{
		timer: myTimer{
			timer:     timer.NewWithInterval(testDuration, time.Second),
			duration:  testDuration,
			isRunning: false,
			timedout:  false,
		},
		wordsToEnter: textToEnter,
		inputBuffer:  make([]rune, 0),
		rawInputCnt:  0,
		mistakes: mistakes{
			mistakesAt:     make(map[int]bool, 0),
			rawMistakesCnt: 0,
		},
		completed: false,
		cursor:    0,
	}
}

func initialModel() model {
	profile := termenv.ColorProfile()
	fore := termenv.ForegroundColor()

	return model{
		state: initialTimerBasedTest(),
		styles: styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink, //we probably don't need this anymore
	)
}
