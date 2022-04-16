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

func initialModel() model {
	generator := NewGenerator()
	generator.Count = 200

	testDuration := time.Second * 15

	textToEnter := generator.Generate()

	profile := termenv.ColorProfile()

	fore := termenv.ForegroundColor()

	return model{
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {

	case timer.TickMsg:
		timerUpdate, cmdUpdate := m.timer.timer.Update(msg)
		m.timer.timer = timerUpdate
		commands = append(commands, cmdUpdate)
		if m.timer.timer.Timedout() {
			m.timer.timedout = true
			m.completed = true
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter", "tab	":

		case "backspace":
			m.inputBuffer = dropLastRune(m.inputBuffer)

			//Delete mistakes
			inputLength := len(m.inputBuffer)
			_, ok := m.mistakes.mistakesAt[inputLength]
			if ok {
				delete(m.mistakes.mistakesAt, inputLength)
			}

		default:

			if !m.completed {
				m.inputBuffer = append(m.inputBuffer, msg.Runes...)
				m.rawInputCnt += len(msg.Runes)
			} else {
				break
			}

			if !m.timer.isRunning {
				commands = append(commands, m.timer.timer.Init())
				m.timer.isRunning = true
			}

			//todo this should be moved to non-time gamemode
			if len(m.inputBuffer) == len(m.wordsToEnter) {
				m.completed = true
			}

			currentInput := string(m.inputBuffer)

			if len(currentInput)-1 == len(m.wordsToEnter) {
				m.completed = true
			} else {

				letterToInput := m.wordsToEnter[len(m.inputBuffer)-1 : len(m.inputBuffer)]
				inputLetter := currentInput[floor(len(currentInput)-1):]

				if letterToInput != inputLetter {
					m.mistakes.mistakesAt[len(m.inputBuffer)-1] = true
					m.mistakes.rawMistakesCnt = m.mistakes.rawMistakesCnt + 1
				}

			}

			return m, tea.Batch(commands...)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}
