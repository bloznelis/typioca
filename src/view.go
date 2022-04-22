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

func (m model) View() string {
	var s string

	switch state := m.state.(type) {
	case TimerBasedTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime

		var style = lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(5).
			PaddingRight(5).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.Color("2"))

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, style.Render(content))

	case TimerBasedTest:
		var lineLenLimit int = 40 // todo: calculate out of model. Have max lineLimit and lower taking term size in consideration

		var coloredTimer string
		if state.timer.isRunning {
			coloredTimer = style(state.timer.timer.View(), m.styles.runningTimer)
		} else {
			coloredTimer = style(state.timer.timer.View(), m.styles.stoppedTimer)
		}

		paragraphView := m.paragraphView(lineLenLimit, state)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		termWidth, termHeight, _ := term.GetSize(0)

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		avgLineLen := averageStringLen(lines[:len(lines)-1])
		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredTimer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
	}

	return s
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

func (m model) paragraphView(lineLimit int, test TimerBasedTest) string {
	paragraph := m.colorInput(test)
	paragraph += m.colorCursor(test)
	paragraph += m.colorWordsToEnter(test)

	wrapped := wrapStyledParagraph(paragraph, lineLimit)

	return wrapped
}

func (m model) colorInput(test TimerBasedTest) string {
	mistakes := toKeysSlice(test.mistakes.mistakesAt)
	sort.Ints(mistakes)

	coloredInput := ""

	if len(mistakes) == 0 {

		coloredInput += styleAllRunes(test.inputBuffer, m.styles.correct)

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := test.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := test.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput += styleAllRunes(sliceUntilMistake, m.styles.correct)
			coloredInput += style(mistakeSlice, m.styles.mistakes)

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := test.inputBuffer[previousMistake+1:]
		coloredInput += styleAllRunes(inputAfterLastMistake, m.styles.correct)
	}

	return coloredInput
}

func (m model) colorCursor(test TimerBasedTest) string {
	cursorLetter := test.wordsToEnter[len(test.inputBuffer) : len(test.inputBuffer)+1]

	return style(cursorLetter, m.styles.cursor)
}

func (m model) colorWordsToEnter(test TimerBasedTest) string {
	wordsToEnter := test.wordsToEnter[len(test.inputBuffer)+1:] // without cursor

	return style(wordsToEnter, m.styles.toEnter)
}

func wrapStyledParagraph(paragraph string, lineLimit int) string {

	// XXX: Replace spaces, because wordwrap trims them out at the ends
	paragraph = strings.Replace(paragraph, " ", "·", -1)

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'·'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))

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

	for _, char := range runes {
		acc += style(string(char)).String()
	}

	return acc
}
