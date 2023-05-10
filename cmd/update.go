package cmd

import (
	"github.com/bloznelis/typioca/cmd/words"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {

	// Update window size
	case tea.WindowSizeMsg:
		if msg.Width == 0 && msg.Height == 0 {
			return m, nil
		} else {
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		}

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
	case MainMenu:
		m.state = state.selections[state.cursor].handleInput(msg, state)
        WriteConfig(state.config)
		return m.quitOn(msg, "ctrl+q")

	case ConfigView:
		m.state = state.handleInput(msg, state)
		return m, nil

	case TimerBasedTestResults:
		m.state = state.handleInput(msg, state)
		return m, nil

	case WordCountTestResults:
		m.state = state.handleInput(msg, state)
		return m, nil

	case SentenceCountTestResults:
		m.state = state.handleInput(msg, state)
		return m, nil

	case TimerBasedTest:
		switch msg := msg.(type) {

		case timer.TickMsg:
			timerUpdate, cmdUpdate := state.timer.timer.Update(msg)
			state.timer.timer = timerUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinutes := state.timer.duration.Minutes() - state.timer.timer.Timeout.Minutes()
			if elapsedMinutes != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinutes))
			}

			m.state = state

			if state.timer.timer.Timedout() {
				termenv.Reset() // get rid of faintness
				state.timer.timedout = true
				m.state = TimerBasedTestResults{
					settings:      state.settings,
					wpmEachSecond: state.base.wpmEachSecond,
					results:       state.calculateResults(),
					mainMenu:      state.mainMenu,
				}
			}

		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "tab":

			case "ctrl+q":
				m.state = state.mainMenu
				return m, nil

			case "ctrl+r":
				m.state = initTimerBasedTest(state.settings, state.mainMenu)
				return m, nil

			case "backspace", "ctrl+h":
				state.base = state.base.handleBackspace()
				m.state = state

			case "ctrl+w":
				state.base = state.base.handleCtrlW()
				m.state = state

			case " ":
				state.base = state.base.handleSpace()
				m.state = state

			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.timer.isRunning {
						commands = append(commands, state.timer.timer.Init())
						state.timer.isRunning = true
					}
					state.base = state.base.handleRunes(msg)
					m.state = state
				}
			}
		}

	case WordCountBasedTest:
		switch msg := msg.(type) {

		case stopwatch.StartStopMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state

		case stopwatch.TickMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinutes := state.stopwatch.stopwatch.Elapsed().Minutes()

			if elapsedMinutes != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinutes))
			}

			m.state = state

		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "tab":

			case "ctrl+q":
				m.state = state.mainMenu
				return m, nil

			case "ctrl+r":
				m.state = initWordCountBasedTest(state.settings, state.mainMenu)
				return m, nil

			case "backspace", "ctrl+h":
				state.base = state.base.handleBackspace()
				m.state = state

			case "ctrl+w":
				state.base = state.base.handleCtrlW()
				m.state = state

			case " ":
				state.base = state.base.handleSpace()
				m.state = state

			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.stopwatch.isRunning {
						commands = append(commands, state.stopwatch.stopwatch.Init())
						state.stopwatch.isRunning = true
					}
					state.base = state.base.handleRunes(msg)
					m.state = state

				}
			}
		}

		// Finished?
		if len(state.base.wordsToEnter) == len(state.base.inputBuffer) {
			termenv.Reset() // get rid of faintness
			m.state = WordCountTestResults{
				settings:      state.settings,
				wpmEachSecond: state.base.wpmEachSecond,
				wordCnt:       state.settings.wordCountSelections[state.settings.wordCountCursor],
				results:       state.calculateResults(),
				mainMenu:      state.mainMenu,
			}
		}

	case SentenceCountBasedTest:
		switch msg := msg.(type) {

		case stopwatch.StartStopMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state

		case stopwatch.TickMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinutes := state.stopwatch.stopwatch.Elapsed().Minutes()
			if elapsedMinutes != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinutes))
			}

			m.state = state

		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "tab":

			case "ctrl+q":
				m.state = state.mainMenu
				return m, nil

			case "ctrl+r":
				m.state = initSentenceCountBasedTest(state.settings, state.mainMenu)
				return m, nil

			case "backspace", "ctrl+h":
				state.base = state.base.handleBackspace()
				m.state = state

			case "ctrl+w":
				state.base = state.base.handleCtrlW()
				m.state = state

			case " ":
				state.base = state.base.handleSpace()
				m.state = state

			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.stopwatch.isRunning {
						commands = append(commands, state.stopwatch.stopwatch.Init())
						state.stopwatch.isRunning = true
					}
					state.base = state.base.handleRunes(msg)
					m.state = state

				}
			}
		}

		//Finished?
		if len(state.base.wordsToEnter) == len(state.base.inputBuffer) {
			termenv.Reset() // get rid of faintness
			m.state = SentenceCountTestResults{
				settings:      state.settings,
				wpmEachSecond: state.base.wpmEachSecond,
				sentenceCnt:   state.settings.sentenceCountSelections[state.settings.sentenceCountCursor],
				results:       state.calculateResults(),
				mainMenu:      state.mainMenu,
			}
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, tea.Batch(commands...)
}

func (m model) quitOn(msg tea.Msg, strokes ...string) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, elem := range strokes {
			if elem == msg.String() {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (settings TimerBasedTestSettings) handleInput(msg tea.Msg, menu MainMenu) State {
	cursorToSave := menu.cursor

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if settings.enabled {
				return initTimerBasedTest(settings, menu)
			}
		case "left", "h":
			if settings.cursor > 0 {
				settings.cursor--
			}
		case "right", "l", "tab":
			if settings.cursor < 2 {
				settings.cursor++
			} else {
				settings.cursor = 0
			}
		case "up", "k":
			switch settings.cursor {
			case 0:
				if menu.cursor > 0 {
					menu.cursor--
				} else {
					menu.cursor = len(menu.selections) - 1
				}
			case 1:
				if settings.timeCursor > 0 {
					settings.timeCursor--
				} else {
					settings.timeCursor = len(settings.timeSelections) - 1
				}
			case 2:
				if settings.wordListCursor > 0 {
					settings.wordListCursor--
				} else {
					settings.wordListCursor = len(settings.wordListSelections) - 1
				}
			}
		case "down", "j":
			switch settings.cursor {
			case 0:
				if menu.cursor < len(menu.selections)-1 {
					menu.cursor++
				}
			case 1:
				if settings.timeCursor < len(settings.timeSelections)-1 {
					settings.timeCursor++
				} else {
					settings.timeCursor = 0
				}
			case 2:
				if settings.wordListCursor < len(settings.wordListSelections)-1 {
					settings.wordListCursor++
				} else {
					settings.wordListCursor = 0
				}
			}
		}
		menu.selections[cursorToSave] = settings
	}

    menu.config.TestSettingCursors.TimerTimeCursor = settings.timeCursor
    menu.config.TestSettingCursors.TimerWordlistCursor = settings.wordListCursor

	return menu
}

func (settings WordCountBasedTestSettings) handleInput(msg tea.Msg, menu MainMenu) State {
	cursorToSave := menu.cursor

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if settings.enabled {
				return initWordCountBasedTest(settings, menu)
			}
		case "left", "h":
			if settings.cursor > 0 {
				settings.cursor--
			}
		case "right", "l", "tab":
			if settings.cursor < 2 {
				settings.cursor++
			} else {
				settings.cursor = 0
			}
		case "up", "k":
			switch settings.cursor {
			case 0:
				if menu.cursor > 0 {
					menu.cursor--
				}
			case 1:
				if settings.wordCountCursor > 0 {
					settings.wordCountCursor--
				} else {
					settings.wordCountCursor = len(settings.wordCountSelections) - 1
				}
			case 2:
				if settings.wordListCursor > 0 {
					settings.wordListCursor--
				} else {
					settings.wordListCursor = len(settings.wordListSelections) - 1
				}
			}
		case "down", "j":
			switch settings.cursor {
			case 0:
				if menu.cursor < len(menu.selections)-1 {
					menu.cursor++
				}
			case 1:
				if settings.wordCountCursor < len(settings.wordCountSelections)-1 {
					settings.wordCountCursor++
				} else {
					settings.wordCountCursor = 0
				}
			case 2:
				if settings.wordListCursor < len(settings.wordListSelections)-1 {
					settings.wordListCursor++
				} else {
					settings.wordListCursor = 0
				}
			}
		}
		menu.selections[cursorToSave] = settings
	}

    menu.config.TestSettingCursors.WordCountCursor = settings.wordCountCursor
    menu.config.TestSettingCursors.WordCountWordlistCursor = settings.wordListCursor

	return menu
}

func (selection ConfigViewSelection) handleInput(msg tea.Msg, menu MainMenu) State {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return initConfigView(menu.config, menu)
		case "up", "k":
			if menu.cursor > 0 {
				menu.cursor--
			}
		case "down", "j":
			if menu.cursor < len(menu.selections)-1 {
				menu.cursor++
			} else {
				menu.cursor = 0
			}

		}
	}

	return menu
}

func (settings SentenceCountBasedTestSettings) handleInput(msg tea.Msg, menu MainMenu) State {
	cursorToSave := menu.cursor

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if settings.enabled {
				return initSentenceCountBasedTest(settings, menu)
			}
		case "left", "h":
			if settings.cursor > 0 {
				settings.cursor--
			}
		case "right", "l", "tab":
			if settings.cursor < 2 {
				settings.cursor++
			} else {
				settings.cursor = 0
			}
		case "up", "k":
			switch settings.cursor {
			case 0:
				if menu.cursor > 0 {
					menu.cursor--
				}
			case 1:
				if settings.sentenceCountCursor > 0 {
					settings.sentenceCountCursor--
				} else {
					settings.sentenceCountCursor = len(settings.sentenceCountSelections) - 1
				}
			case 2:
				if settings.sentenceListCursor > 0 {
					settings.sentenceListCursor--
				} else {
					settings.sentenceListCursor = len(settings.sentenceListSelections) - 1
				}
			}
		case "down", "j":
			switch settings.cursor {
			case 0:
				if menu.cursor < len(menu.selections)-1 {
					menu.cursor++
				}
			case 1:
				if settings.sentenceCountCursor < len(settings.sentenceCountSelections)-1 {
					settings.sentenceCountCursor++
				} else {
					settings.sentenceCountCursor = 0
				}
			case 2:
				if settings.sentenceListCursor < len(settings.sentenceListSelections)-1 {
					settings.sentenceListCursor++
				} else {
					settings.sentenceListCursor = 0
				}
			}
		}
		menu.selections[cursorToSave] = settings
	}

    menu.config.TestSettingCursors.SentenceCountCursor = settings.sentenceCountCursor
    menu.config.TestSettingCursors.SentenceCountWordlistCursor = settings.sentenceListCursor

	return menu
}

func (configView ConfigView) handleInput(msg tea.Msg, state State) State {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			embededLength := len(configView.config.EmbededWordLists)
			if configView.cursor < embededLength {
				configView.config.EmbededWordLists[configView.cursor].toggleEnabled()
			} else {
				configView.config.WordLists[configView.cursor-embededLength].toggleEnabled()
			}

            // We might not have wordlist that config points to
            configView.config.TestSettingCursors.resetWordlistCursors()

			WriteConfig(configView.config)
			state = configView
		case "s":
			embededLength := len(configView.config.EmbededWordLists)
			if configView.cursor > embededLength-1 {
				configView.config.WordLists[configView.cursor-embededLength].toggleSynced()
			}

            // We might not have wordlist that config points to
            configView.config.TestSettingCursors.resetWordlistCursors()

			WriteConfig(configView.config)
			state = configView
		case "up", "k":
			if configView.cursor > 0 {
				configView.cursor--
			} else {
				configView.cursor = configView.config.wordListsCount() - 1
			}

			state = configView
		case "down", "j":
			if configView.cursor < configView.config.wordListsCount()-1 {
				configView.cursor++
			} else {
				configView.cursor = 0
			}

			state = configView
		case "ctrl+q":
			state = initMainMenu()
		}
	}

	return state
}

func (wl *WordList) toggleSynced() {
	if !wl.isLocal {
		var err error
		if wl.synced {
			err = words.DeleteWordList(wl.Path)
		} else {
			err = words.DownloadFile(wl.Path, wl.RemoteURI)
		}

		if err != nil {
			// log.Panicln(err)
			wl.syncOK = false
		} else {
			wl.synced = !wl.synced
			wl.syncOK = true
		}
	}
}

func (results TimerBasedTestResults) handleInput(msg tea.Msg, state State) State {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+r":
			state = initTimerBasedTest(results.settings, results.mainMenu)
		case "ctrl+q":
			state = results.mainMenu
		}
	}

	return state
}

func (results WordCountTestResults) handleInput(msg tea.Msg, state State) State {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+r":
			state = initWordCountBasedTest(results.settings, results.mainMenu)
		case "ctrl+q":
			state = results.mainMenu
		}
	}

	return state
}

func (results SentenceCountTestResults) handleInput(msg tea.Msg, state State) State {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+r":
			state = initSentenceCountBasedTest(results.settings, results.mainMenu)
		case "ctrl+q":
			state = results.mainMenu
		}
	}

	return state
}

func (base TestBase) handleBackspace() TestBase {
	base.inputBuffer = dropLastRune(base.inputBuffer)

	//Delete mistakes
	inputLength := len(base.inputBuffer)
	_, ok := base.mistakes.mistakesAt[inputLength]
	if ok {
		delete(base.mistakes.mistakesAt, inputLength)
	}

	base.cursor = inputLength

	return base
}

func (base TestBase) handleCtrlW() TestBase {
	base.inputBuffer = dropUntilWsIdx(base.inputBuffer, base.findLatestWsIndex())
	bufferLen := len(base.inputBuffer)
	base.cursor = bufferLen

	//Delete mistakes
	var newMistakes map[int]bool = make(map[int]bool, 0)
	for at := range base.mistakes.mistakesAt {
		if at < bufferLen {
			newMistakes[at] = true
		}
	}
	base.mistakes.mistakesAt = newMistakes

	return base
}

func dropUntilWsIdx(input []rune, wsIdx int) []rune {
	if wsIdx == 0 {
		return make([]rune, 0)
	} else {
		return input[:wsIdx+1]
	}
}

func (base TestBase) handleRunes(msg tea.KeyMsg) TestBase {
	base.inputBuffer = append(base.inputBuffer, msg.Runes...)
	base.rawInputCnt += len(msg.Runes)

	inputLen := len(base.inputBuffer)
	inputLenDec := inputLen - 1

	letterToInput := base.wordsToEnter[inputLenDec]
	inputLetter := base.inputBuffer[inputLenDec]

	if letterToInput != inputLetter {
		base.mistakes.mistakesAt[inputLenDec] = true
		base.mistakes.rawMistakesCnt = base.mistakes.rawMistakesCnt + 1
	}

	//Set cursor
	base.cursor = inputLen

	return base
}

func (base TestBase) handleSpace() TestBase {
	if len(base.inputBuffer) > 0 && base.wordsToEnter[base.cursor-1] != ' ' {
		nextSpaceIdx := findNextSpaceIndex(base.wordsToEnter, base.cursor)
		var spacesToAppend []rune
		if nextSpaceIdx == base.cursor {
			spacesToAppend = []rune{' '}
		} else {
			spacesToEnterCnt := (nextSpaceIdx - base.cursor) + 1
			spaces := make([]rune, spacesToEnterCnt)
			for i := range spaces {
				spaces[i] = ' '
			}
			spacesToAppend = spaces

		}

		// Mark mistakes
		for i := base.cursor; i <= nextSpaceIdx; i++ {
			if base.wordsToEnter[i] != ' ' {
				base.mistakes.mistakesAt[i] = true
				base.mistakes.rawMistakesCnt = base.mistakes.rawMistakesCnt + 1
			}
		}

		base.inputBuffer = append(base.inputBuffer, spacesToAppend...)
		base.cursor = len(base.inputBuffer)
		base.rawInputCnt += 1
	}

	return base
}

func findNextSpaceIndex(wordsToInput []rune, cursorAt int) int {
	trimmedWordsToInput := wordsToInput[cursorAt:]
	words := trimmedWordsToInput

	var wsIdx int = 0
	for idx, value := range words {
		if value == ' ' {
			wsIdx = idx
			break
		}
	}

	return wsIdx + cursorAt
}

func (base TestBase) findLatestWsIndex() int {
	var wsIdx int = 0
	for idx, value := range base.wordsToEnter {
		if idx+1 >= base.cursor {
			break
		}
		if value == ' ' {
			wsIdx = idx
		}
	}

	return wsIdx
}
