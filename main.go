package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	babble "github.com/Beartime234/babble"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type myTimer struct {
	timer     timer.Model
	duration  time.Duration
	isRunning bool // Inner is running is being handled weirdly.
	timedout  bool
}

type StringStyle func(string) termenv.Style

type styles struct {
	mistakes StringStyle
	correct  StringStyle
	toEnter  StringStyle
}

type model struct {
	styles       styles
	timer        myTimer
	wordsToEnter string
	inputBuffer  []rune
	mistakesAt   map[int]bool
	completed    bool
	cursor       int
}

func initialModel() model {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = 20

	testDuration := time.Second * 120

	textToEnter := babbler.Babble()

	profile := termenv.ColorProfile()

	return model{
		styles: styles{
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1"))
			},
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10"))
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
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

func toKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
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

			//Delete mistakes
			_, ok := m.mistakesAt[len(m.inputBuffer)]
			if ok {
				delete(m.mistakesAt, len(m.inputBuffer))
			}

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

				letterToInput := m.wordsToEnter[len(m.inputBuffer)-1 : len(m.inputBuffer)]
				inputLetter := currentInput[floor(len(currentInput)-1):]

				if letterToInput != inputLetter {
					m.mistakesAt[len(m.inputBuffer)-1] = true
				}

			}

			return m, tea.Batch(commands...)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}

func style(str string, style StringStyle) string {
	return style(str).String()
}

func styleAllRunes(runes []rune, style StringStyle) string {
	acc := ""

	for _, char := range runes {
		acc += style(string(char)).String()
	}

	return acc
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

		mistakes := toKeysSlice(m.mistakesAt)
		sort.Ints(mistakes)

		paragraph := ""

		coloredInput := ""

		if len(mistakes) == 0 {

			coloredInput += styleAllRunes(m.inputBuffer, m.styles.correct)

		} else {

			previousMistake := -1

			for _, mistakeAt := range mistakes {
				sliceUntilMistake := m.inputBuffer[previousMistake+1 : mistakeAt]
				mistakeSlice := m.wordsToEnter[mistakeAt : mistakeAt+1]

				coloredInput += styleAllRunes(sliceUntilMistake, m.styles.correct)
				coloredInput += style(mistakeSlice, m.styles.mistakes)

				previousMistake = mistakeAt
			}

			inputAfterLastMistake := m.inputBuffer[previousMistake+1:]
			coloredInput += styleAllRunes(inputAfterLastMistake, m.styles.correct)
		}

		remainingWordsToEnterWithoutCursorLetter := m.wordsToEnter[len(m.inputBuffer)+1:]
		cursorLetter := m.wordsToEnter[len(m.inputBuffer) : len(m.inputBuffer)+1]

		paragraph += coloredInput
		paragraph += termenv.String(cursorLetter).Underline().String()
		paragraph += remainingWordsToEnterWithoutCursorLetter

		return lipgloss.NewStyle().Width(80).Padding(10).Align(lipgloss.Center).Render(paragraph)
	}

	// Send the UI for rendering
	return s
}

func main() {
	termenv.ClearScreen()
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
