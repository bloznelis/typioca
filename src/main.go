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

func dropUntilWsIdx(input []rune, wsIdx int) []rune {
	if wsIdx == 0 {
		return make([]rune, 0)
	} else {
		return input[:wsIdx+1]
	}
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "enter", "ctrl+r":
			m.state = initialTimerBasedTest()
			return m, nil

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		}
	}

	switch state := m.state.(type) {

	case TimerBasedTestResults:
		break

	case TimerBasedTest:
		switch msg := msg.(type) {

		case timer.TickMsg:
			timerUpdate, cmdUpdate := state.timer.timer.Update(msg)
			state.timer.timer = timerUpdate
			commands = append(commands, cmdUpdate)

			m.state = state
			if state.timer.timer.Timedout() {
				state.timer.timedout = true
				m.state = TimerBasedTestResults{results: state.calculateResults()}
			}

		case tea.KeyMsg:

			switch msg.String() {

			case "backspace":
				state.inputBuffer = dropLastRune(state.inputBuffer)

				//Delete mistakes
				inputLength := len(state.inputBuffer)
				_, ok := state.mistakes.mistakesAt[inputLength]
				if ok {
					delete(state.mistakes.mistakesAt, inputLength)
				}

				m.state = state

			case "ctrl+w":
				state.inputBuffer = dropUntilWsIdx(state.inputBuffer, state.findLatestWsIndex())
				bufferLen := len(state.inputBuffer)
				state.cursor = bufferLen

				//Delete mistakes
				var newMistakes map[int]bool = make(map[int]bool, 0)
				for at := range state.mistakes.mistakesAt {
					if at < bufferLen {
						newMistakes[at] = true
					}
				}
				state.mistakes.mistakesAt = newMistakes

				m.state = state

			case "tab":
				m.state = TimerBasedTestResults{results: state.calculateResults()}

			default:
				state.inputBuffer = append(state.inputBuffer, msg.Runes...)
				state.rawInputCnt += len(msg.Runes)

				if !state.timer.isRunning {
					commands = append(commands, state.timer.timer.Init())
					state.timer.isRunning = true
				}

				inputLen := len(state.inputBuffer)
				inputLenDec := inputLen - 1

				letterToInput := state.wordsToEnter[inputLenDec:inputLen]
				inputLetter := string(state.inputBuffer[inputLenDec:])

				if letterToInput != inputLetter {
					state.mistakes.mistakesAt[inputLenDec] = true
					state.mistakes.rawMistakesCnt = state.mistakes.rawMistakesCnt + 1
				}

				//Set cursor
				state.cursor = inputLen

				m.state = state
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}

func (state TimerBasedTest) findLatestWsIndex() int {
	return findLatestWsIndex(state.wordsToEnter, state.cursor)
}

func findLatestWsIndex(wordsToInput string, cursorAt int) int {
	trimmedWordsToInput := wordsToInput[:cursorAt]
	reversed := reverse([]rune(trimmedWordsToInput))

	var wsIdx int = 0
	for idx, value := range reversed {
		if value == ' ' && idx != 0 {
			wsIdx = len(reversed) - 1 - idx
			break
		}
	}

	return int(floor(wsIdx))
}