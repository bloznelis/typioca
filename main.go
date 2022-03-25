package main

import (
	"fmt"
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

	inputModel := input.NewModel()
	inputModel.Focus()
	inputModel.CursorStyle.Blink(true)
	// inputModel.CursorStyle = inputModel.CursorStyle.Copy().Blink(true)
	inputModel.Prompt = "  " // Try adding some padding instead
	// inputModel.Placeholder = babbler.Babble()

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// case "enter":
		// 	words := m.babbler.Babble()
		// 	m.wordsToEnter = words

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		default:
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, cmd
}

func (m model) View() string {

	s := "\n"
	// s += "\n\n"

	remainingWordsToEnter := m.wordsToEnter[len(m.textInput.Value())+1:]

	s += m.textInput.View()
	s += remainingWordsToEnter
	s += "\n\n"

	// Send the UI for rendering
	return s
}

func main() {

	// fmt.Println("abcd"[:2])

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
