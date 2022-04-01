package main

import (
	"fmt"
	"log"
	"math"
	"os"

	babble "github.com/Beartime234/babble"
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices      []string         // items on the to-do list
	cursor       int              // which to-do list item our cursor is pointing at
	selected     map[int]struct{} // which to-do items are selected
	babbler      babble.Babbler   // Word generator
	wordsToEnter string
	textInput    input.Model
}

func initialModel() model {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = 10

	textToEnter := babbler.Babble()

	log.Println("text to enter", textToEnter)

	inputModel := input.NewModel()
	inputModel.Focus()
	inputModel.CursorStyle.Blink(true)
	inputModel.SetCursorMode(input.CursorBlink)
	inputModel.SetValue(textToEnter[:1])
	inputModel.SetCursor(inputModel.Cursor() - 1)
	inputModel.Prompt = "  " // Try adding some padding instead
	// inputModel.Placeholder = textToEnter

	return model{
		// Our shopping list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		// babbler:   babbler,
		wordsToEnter: textToEnter,
		textInput:    inputModel,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	// return nil
	return input.Blink //???
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "backspace":

		default:
			m.textInput, cmd = m.textInput.Update(msg)
			currentInput := m.textInput.Value()
			letterToInput := currentInput[floor(len(currentInput)-1):]
			inputLetter := currentInput[floor(len(currentInput)-2):floor(len(currentInput)-1)]
			nextLetter := m.wordsToEnter[floor(len(currentInput)-1):len(currentInput)]

			// log.Println("current input key", msg.String())

			// log.Println("current input", currentInput)
			// fmt.Println("letter to input", letterToInput)
			// fmt.Println("input letter", inputLetter)

			//todo maybe this would work?
			// having letter to input as the last one
			// and checking whether it matches or not.
			//We should never allow to put cursor on that "letter to input"

			// hello mom
			// hell
			//     o mom
			//     o

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

	}

	// Remaining key strokes and blink messages passed here
	// m.textInput, cmd = m.textInput.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, cmd
}

func (m model) View() string {

	s := m.wordsToEnter
	s += "\n"
	// s += "\n\n"

	remainingWordsToEnter := m.wordsToEnter[len(m.textInput.Value()):]

	s += m.textInput.View()
	s += remainingWordsToEnter
	s += "\n\n"

	// Send the UI for rendering
	return s
}

func main() {

	// first := "abcd"[2:]
	// second := first[:1]

	// fmt.Println(second)

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
