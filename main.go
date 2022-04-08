package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	babble "github.com/Beartime234/babble"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type myTimer struct {
	timer     timer.Model
	duration  time.Duration
	isRunning bool // Inner is running is being handled weirdly.
	timedout  bool
}

type model struct {
	timer        myTimer
	colorProfile termenv.Profile
	wordsToEnter string
	inputBuffer  []rune
	mistakesAt   map[int]bool
	completed    bool
	cursor       int
}

func initialModel() model {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = 100

	testDuration := time.Second * 30

	textToEnter := babbler.Babble()

	return model{
		timer: myTimer{
			timer:     timer.NewWithInterval(testDuration, time.Second),
			duration:  testDuration,
			isRunning: false,
			timedout:  false,
		},
		colorProfile: termenv.ColorProfile(),
		wordsToEnter: textToEnter,
		inputBuffer:  make([]rune, 0),
		mistakesAt:   make(map[int]bool, 0),
		completed:    false,
		cursor:       0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink,
	)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func floor(value int) int32 {
	return int32(math.Max(0, float64(value)))
}

func dropLast(value string) string {
	return dropLastN(1, value)
}

func dropLastN(n int, value string) string {
	return value[:len(value)-n]
}

func dropLastRune(runes []rune) []rune {
	le := len(runes)
	if le != 0 {
		return runes[:le-1]
	} else {
		return runes
	}
}

func getCorrectWords(m model) []string {
	wordsToEnter := strings.Split(m.wordsToEnter, " ")
	enteredWords := strings.Split(string(m.inputBuffer), " ")

	var correctWords []string

	for _, enteredWord := range enteredWords {
		if contains(wordsToEnter, enteredWord) {
			correctWords = append(correctWords, enteredWord)
		}

	}

	return correctWords
}

func calculateWpm(m model) int {
	correctWords := getCorrectWords(m)
	testDuration := m.timer.duration

	return int(float64(len(correctWords)) / testDuration.Minutes())
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

		default:

			if !m.completed {
				m.inputBuffer = append(m.inputBuffer, msg.Runes...)
			} else {
				break
			}

			if !m.timer.isRunning {
				commands = append(commands, m.timer.timer.Init())
				m.timer.isRunning = true
			}

			if len(m.inputBuffer) == len(m.wordsToEnter) {
				m.completed = true
			}

			currentInput := string(m.inputBuffer)

			if len(currentInput)-1 == len(m.wordsToEnter) {
				m.completed = true
			} else {

				// abc lukas acc
				// abc z

				letterToInput := m.wordsToEnter[len(m.inputBuffer)-1 : len(m.inputBuffer)]
				inputLetter := currentInput[floor(len(currentInput)-1):]
				// nextLetter := m.wordsToEnter[floor(len(currentInput)-1):len(currentInput)]

				// println("letter to input ", letterToInput)
				// println("input letter ", inputLetter)

				if letterToInput != inputLetter {
					m.mistakesAt[len(m.inputBuffer)-1] = true
				}

			}

			return m, tea.Batch(commands...)
		}
	}

	// Remaining key strokes and blink messages passed here
	// if !m.completed {
	// 	textInputUpdate, cmdUpdate := m.textInput.Update(msg)
	// 	m.textInput = textInputUpdate
	// 	commands = append(commands, cmdUpdate)
	// }

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}

func (m model) View() string {
	s := ""

	if m.timer.timedout {
		s += "Timer timedout!"
		s += "\n\n"
		s += "WPM: "
		s += strconv.Itoa(calculateWpm(m))
		s += "\n\n"
		// s += m.wordsToEnter
		// s += "\n\n"
		// s += m.textInput.Value()

	} else if m.completed {
		s += "Out of words lol"
	} else {
		s += m.timer.timer.View()
		s += "\n\n"
		s += fmt.Sprintln(m.mistakesAt)
		s += "\n\n"

		remainingWordsToEnterWithoutCursorLetter := m.wordsToEnter[len(m.inputBuffer)+1:]
		cursorLetter := m.wordsToEnter[len(m.inputBuffer) : len(m.inputBuffer)+1]

		s += termenv.String(string(m.inputBuffer)).Foreground(m.colorProfile.Color("3")).String()
		s += termenv.String(cursorLetter).Underline().String()
		s += remainingWordsToEnterWithoutCursorLetter
		s += "\n\n"
	}

	// Send the UI for rendering
	return s
}

func main() {

	// runes := make([]rune, 0)
	// runes = append(runes, 'a')
	// runes = append(runes, 'b')
	// runes = append(runes, 'c')

	// println(string(runes))

	str := "abefcd"
	stri := "abe"
	println(str[len(stri)-1 : len(stri)])

	termenv.ShowCursor()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
