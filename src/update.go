package main

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		}
	}

	switch state := m.state.(type) {

	case TimerBasedTestResults:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "ctrl+r":
				m.state = initialTimerBasedTest()
				return m, nil
			}
		}

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
			case "enter":

			case "ctrl+r":
				m.state = initialTimerBasedTest()
				return m, nil

			case "backspace":
				m.state = state.handleBackspace()

			case "ctrl+w":
				m.state = state.handleCtrlW()

			case " ":
				m.state = state.handleSpace()

			default:
				m.state = state.handleRunes(msg, &commands)
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}

func (state TimerBasedTest) handleBackspace() TimerBasedTest {
	state.inputBuffer = dropLastRune(state.inputBuffer)

	//Delete mistakes
	inputLength := len(state.inputBuffer)
	_, ok := state.mistakes.mistakesAt[inputLength]
	if ok {
		delete(state.mistakes.mistakesAt, inputLength)
	}

	return state
}

func (state TimerBasedTest) handleCtrlW() TimerBasedTest {
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

	return state
}

func dropUntilWsIdx(input []rune, wsIdx int) []rune {
	if wsIdx == 0 {
		return make([]rune, 0)
	} else {
		return input[:wsIdx+1]
	}
}

func (state TimerBasedTest) handleRunes(msg tea.KeyMsg, commands *[]tea.Cmd) TimerBasedTest {
	state.inputBuffer = append(state.inputBuffer, msg.Runes...)
	state.rawInputCnt += len(msg.Runes)

	if !state.timer.isRunning {
		*commands = append(*commands, state.timer.timer.Init())
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

	return state
}

func (state TimerBasedTest) handleSpace() TimerBasedTest {
	if len(state.inputBuffer) > 0 && state.wordsToEnter[state.cursor-1] != ' ' {
		nextSpaceIdx := findNextSpaceIndex(state.wordsToEnter, state.cursor)
		spacesToEnterCnt := (nextSpaceIdx - state.cursor) + 1
		spaces := make([]rune, spacesToEnterCnt)
		for i := range spaces {
			spaces[i] = ' '
		}

		if spacesToEnterCnt > 1 {
			//Mark mistakes
			for i := state.cursor; i < nextSpaceIdx; i++ {
				state.mistakes.mistakesAt[i] = true
			}

			state.mistakes.rawMistakesCnt = state.mistakes.rawMistakesCnt + 1
		}

		state.inputBuffer = append(state.inputBuffer, spaces...)
		state.cursor = len(state.inputBuffer)
		state.rawInputCnt += 1
	}

	return state
}

func findNextSpaceIndex(wordsToInput string, cursorAt int) int {
	trimmedWordsToInput := wordsToInput[cursorAt:]
	words := []rune(trimmedWordsToInput)

	var wsIdx int = 0
	for idx, value := range words {
		if value == ' ' {
			wsIdx = idx
			break
		}
	}

	return wsIdx + cursorAt
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
