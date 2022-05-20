package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

var avgLineLen int = 0
var lineLenLimit int = 40
var resultsStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(5).
	PaddingRight(5).
	BorderStyle(lipgloss.HiddenBorder()).
	BorderForeground(lipgloss.Color("2"))

func (m model) View() string {
	var s string

	termWidth, termHeight, _ := term.GetSize(0)

	reactiveLimit := (termWidth / 10) * 6
	if reactiveLimit < lineLenLimit {
		lineLenLimit = reactiveLimit
	}

	switch state := m.state.(type) {
	case MainMenu:
		s := style("  typioca", m.styles.faintGreen)
		s += "\n\n\n"

		for i, choice := range state.selections {
			cursor := " "
			cursorClose := " "
			if state.cursor == i {
				cursor = style(">", m.styles.runningTimer)
				cursorClose = style("<", m.styles.runningTimer)
			}

			// Render the row
			s += fmt.Sprintf("%s %s%s\n\n", cursor, choice.show(m.styles), cursorClose)
		}
		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.NewStyle().Align(lipgloss.Left).Render(s)

		return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	case TimerBasedTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime + " " + words

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, resultsStyle.Render(content))

	case WordCountTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		wordCnt := "cnt: " + style(strconv.Itoa(state.wordCnt), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime + " " + wordCnt + " " + words

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, resultsStyle.Render(content))

	case SentenceCountTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		sentenceCnt := "cnt: " + style(strconv.Itoa(state.sentenceCnt), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "sentences: " + style(state.results.wordList, m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime + " " + sentenceCnt + " " + words

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, resultsStyle.Render(content))

	case TimerBasedTest:
		var coloredTimer string
		if state.timer.isRunning {
			coloredTimer = style(state.timer.timer.View(), m.styles.runningTimer)
		} else {
			coloredTimer = style(state.timer.timer.View(), m.styles.stoppedTimer)
		}

		paragraphView := state.base.paragraphView(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		if avgLineLen == 0 {
			avgLineLen = averageLineLen(lines)
		}

		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredTimer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.timer.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

	case WordCountBasedTest:
		var coloredStopwatch string
		if state.stopwatch.isRunning {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.runningTimer)
		} else {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.stoppedTimer)
		}

		paragraphView := state.base.paragraphView(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		if avgLineLen == 0 {
			avgLineLen = averageLineLen(lines)
		}
		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.stopwatch.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

	case SentenceCountBasedTest:
		var coloredStopwatch string
		if state.stopwatch.isRunning {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.runningTimer)
		} else {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.stoppedTimer)
		}

		paragraphView := state.base.paragraphView(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		if avgLineLen == 0 {
			avgLineLen = averageLineLen(lines)
		}
		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.stopwatch.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

	}

	return s
}

func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1] //Drop last line, as it might skew up average length
	}

	return averageStringLen(lines)
}

func (selection TimerBasedTestSettings) show(styles Styles) string {
	selections := []string{selection.timeSelections[selection.timeCursor].String(), selection.wordListSelections[selection.wordListCursor].show}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Timer run", selectionsStr)
}

func (selection WordCountBasedTestSettings) show(styles Styles) string {
	selections := []string{fmt.Sprint(selection.wordCountSelections[selection.wordCountCursor]), selection.wordListSelections[selection.wordListCursor].show}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Word count run", selectionsStr)
}

func (selection SentenceCountBasedTestSettings) show(styles Styles) string {
	selections := []string{fmt.Sprint(selection.sentenceCountSelections[selection.sentenceCountCursor]), selection.sentenceListSelections[selection.sentenceListCursor].show}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Sentence count run", selectionsStr)
}

func showSelections(selections []string, cursor int, styles Styles) string {
	var selectionsStr string
	for i, option := range selections {
		if i+1 == cursor {
			selectionsStr += "[" + style(option, styles.runningTimer) + "]"
		} else {
			selectionsStr += "[" + style(option, styles.greener) + "]"
		}
		selectionsStr += " "
	}
	return selectionsStr
}

func getLinesAroundCursor(lines []string, cursorLine int) []string {
	cursor := cursorLine

	// 3 lines to show
	if cursorLine == 0 {
		cursor += 3
	} else {
		cursor += 2
	}

	low := int(math.Max(0, float64(cursorLine-1)))
	high := int(math.Min(float64(len(lines)), float64(cursor)))

	return lines[low:high]
}

func dropAnsiCodes(colored string) string {
	m := regexp.MustCompile("\x1b\\[[0-9;]*m")

	return m.ReplaceAllString(colored, "")
}

func (m model) indent(block string, indentBy uint) string {
	indentedBlock := indent.String(block, indentBy) // this crashes on small windows

	return indentedBlock
}

func (base TestBase) paragraphView(lineLimit int, styles Styles) string {
	paragraph := base.colorInput(styles)
	paragraph += base.colorCursor(styles)
	paragraph += base.colorWordsToEnter(styles)

	wrapped := wrapStyledParagraph(paragraph, lineLimit)

	return wrapped
}

func (base TestBase) colorInput(styles Styles) string {
	mistakes := toKeysSlice(base.mistakes.mistakesAt)
	sort.Ints(mistakes)

	coloredInput := ""

	if len(mistakes) == 0 {

		coloredInput += styleAllRunes(base.inputBuffer, styles.correct)

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := base.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := base.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput += styleAllRunes(sliceUntilMistake, styles.correct)
			coloredInput += style(mistakeSlice, styles.mistakes)

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := base.inputBuffer[previousMistake+1:]
		coloredInput += styleAllRunes(inputAfterLastMistake, styles.correct)
	}

	return coloredInput
}

func (base TestBase) colorCursor(styles Styles) string {
	cursorLetter := base.wordsToEnter[len(base.inputBuffer) : len(base.inputBuffer)+1]

	return style(cursorLetter, styles.cursor)
}

func (base TestBase) colorWordsToEnter(styles Styles) string {
	wordsToEnter := base.wordsToEnter[len(base.inputBuffer)+1:] // without cursor

	return style(wordsToEnter, styles.toEnter)
}

func wrapStyledParagraph(paragraph string, lineLimit int) string {
	// XXX: Replace spaces, because wordwrap trims them out at the ends
	paragraph = strings.Replace(paragraph, " ", "·", -1)

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'·'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))
	f.Close()

	paragraph = strings.Replace(f.String(), "·", " ", -1)

	return paragraph
}

func findCursorLine(lines []string, cursorAt int) int {
	lenAcc := 0
	cursorLine := 0

	for _, line := range lines {
		lineLen := len(dropAnsiCodes(line))

		lenAcc += lineLen

		if cursorAt <= lenAcc-1 {
			return cursorLine
		} else {
			cursorLine += 1
		}
	}

	return cursorLine
}

func style(str string, style StringStyle) string {
	return style(str).String()
}

func styleAllRunes(runes []rune, style StringStyle) string {
	acc := ""

	for idx, char := range runes {
		_ = idx
		acc += style(string(char)).String()
		// if idx == 0 {
		// 	acc += style(string(char)).String()
		// } else {
		// 	acc += string(char)
		// }
	}

	return acc
}
