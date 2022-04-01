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
)

type myTimer struct {
	timer     timer.Model
	duration  time.Duration
	isRunning bool // Inner is running is being handled weirdly.
	timedout  bool
}

type model struct {
	timer        myTimer
	wordsToEnter string
	textInput    input.Model
	completed    bool
}

func initialModel() model {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = 100

	testDuration := time.Second * 30

	textToEnter := babbler.Babble()

	inputModel := input.NewModel()
	inputModel.Focus()
	inputModel.CursorStyle.Blink(true)
	inputModel.SetCursorMode(input.CursorBlink)
	inputModel.SetValue(textToEnter[:1])
	inputModel.SetCursor(inputModel.Cursor() - 1)
	inputModel.Prompt = "  " // Try adding some padding instead

	return model{
		timer: myTimer{
			timer:     timer.NewWithInterval(testDuration, time.Second),
			duration:  testDuration,
			isRunning: false,
			timedout:  false,
		},
		wordsToEnter: textToEnter,
		textInput:    inputModel,
		completed:    false,
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

func getCorrectWords(m model) []string {
	wordsToEnter := strings.Split(m.wordsToEnter, " ")
	enteredWords := strings.Split(m.textInput.Value(), " ")

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
			// XXX: If we catch backspace here, it does not ge propagated to be handled by input field
			if len(m.textInput.Value()) > 1 {
				textInputUpdate, cmdUpdate := m.textInput.Update(msg)
				m.textInput = textInputUpdate
				commands = append(commands, cmdUpdate)

				currentInput := m.textInput.Value()
				bareInput := dropLast(currentInput)
				nextLetter := m.wordsToEnter[floor(len(currentInput)-1):len(currentInput)]
				inputWithNext := fmt.Sprintf("%s%s", bareInput, nextLetter)
				m.textInput.SetValue(inputWithNext)
			}

			return m, tea.Batch(commands...)

		default:

			if !m.completed {
				textInputUpdate, cmdUpdate := m.textInput.Update(msg)
				m.textInput = textInputUpdate
				commands = append(commands, cmdUpdate)
			} else {
				break
			}

			currentInput := m.textInput.Value()

			if len(currentInput)-1 == len(m.wordsToEnter) {
				m.completed = true
			} else {
				// Having letter to input as the last one
				// and checking whether it matches or not.
				letterToInput := currentInput[floor(len(currentInput)-1):]
				inputLetter := currentInput[floor(len(currentInput)-2):floor(len(currentInput)-1)]
				nextLetter := m.wordsToEnter[floor(len(currentInput)-1):len(currentInput)]

				if !m.timer.isRunning {
					commands = append(commands, m.timer.timer.Init())
					m.timer.isRunning = true
				}

				if letterToInput == inputLetter {
					bareInput := dropLast(currentInput)
					inputWithNext := fmt.Sprintf("%s%s", bareInput, nextLetter)

					m.textInput.SetValue(inputWithNext)
				} else {
					bareInput := dropLastN(2, currentInput) // Drop last 2, because we replace wrong input with X
					inputWithWrongAndNext := fmt.Sprintf("%s%s%s", bareInput, "X", nextLetter)

					m.textInput.SetValue(inputWithWrongAndNext)
				}

			}

			return m, tea.Batch(commands...)
		}
	}

	// Remaining key strokes and blink messages passed here
	if !m.completed {
		textInputUpdate, cmdUpdate := m.textInput.Update(msg)
		m.textInput = textInputUpdate
		commands = append(commands, cmdUpdate)
	}

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

		remainingWordsToEnter := m.wordsToEnter[len(m.textInput.Value()):]

		s += m.textInput.View()
		s += remainingWordsToEnter
		s += "\n\n"
	}

	// Send the UI for rendering
	return s
}

func main() {

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
