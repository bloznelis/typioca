package cmd

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
)

var lineLenLimit int
var minLineLen int = 5
var maxLineLen int = 40
var resultsStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(5).
	PaddingRight(5)

func wrapWithCursor(shouldWrap bool, line string, stringStyle StringStyle) string {
	cursor := " "
	cursorClose := " "
	if shouldWrap {
		cursor = style(">", stringStyle)
		cursorClose = style("<", stringStyle)
	}

	return fmt.Sprintf("%s %s%s", cursor, line, cursorClose)
}

func (m model) View() string {
	var s string

	termWidth, termHeight := m.width, m.height

	reactiveLimit := (termWidth * 6) / 10
	lineLenLimit = int(math.Min(float64(maxLineLen), math.Max(float64(minLineLen), float64(reactiveLimit))))

	switch state := m.state.(type) {
	case MainMenu:
		typioca := style("  typioca", m.styles.faintGreen)
		typioca = lipgloss.NewStyle().PaddingBottom(1).Render(typioca)

		var choices []string
		choiceStyle := lipgloss.NewStyle().PaddingTop(1)
		for i, choice := range state.selections {
			choiceShow := choice.show(m.styles)
			if !choice.Enabled() {
				choiceShow = style(dropAnsiCodes(choiceShow), m.styles.toEnter)
			}

			choiceShow = wrapWithCursor(state.cursor == i, choiceShow, m.styles.runningTimer)
			choiceShow = choiceStyle.Render(choiceShow)
			choices = append(choices, choiceShow)
		}

		joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{typioca}, choices...)...)
		s = lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)

		return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	case ConfigView:
		absolutePad := longestStringLen(names(state.config.WordLists)) + 2
		var view string
		header := fmt.Sprintf("%s%*s%s/%s\n\n", "  wordlist", absolutePad-11, " ", "synced", "enabled")
		view += header

		for idx, elem := range state.config.EmbededWordLists {
			var enabled string
			if elem.Enabled {
				enabled = "x"
			} else {
				enabled = " "
			}
			enabled = style(enabled, m.styles.greener)

			toPad := absolutePad - len(elem.Name)
			line := fmt.Sprintf("%s%*s     [%s] ", style(elem.Name, m.styles.greener), toPad, "", enabled)

			view += wrapWithCursor(idx == state.cursor, line, m.styles.runningTimer)
			view += "\n"
		}
		view += "\n"

		maxAmmountToShow := m.height / 7
		total := len(state.config.WordLists)
		curs := state.cursor - len(state.config.EmbededWordLists) + 1
		lower := floor(curs - maxAmmountToShow)
		upper := int(math.Min(math.Max(float64(curs), float64(maxAmmountToShow)), float64(total)))

		wordListsToShow := state.config.WordLists[lower:upper]
		cursorWidget := fmt.Sprintf("  [%d-%d:%d]", lower+1, upper, total)

		view += style(cursorWidget, m.styles.toEnter)
		view += "\n"

		for idx, elem := range wordListsToShow {
			var synced string
			if elem.synced {
				synced = "x"
			} else {
				synced = " "
			}
			if !elem.isLocal {
				synced = style(synced, m.styles.greener)
			} else {
				synced = style(synced, m.styles.toEnter)
			}

			var enabled string
			if elem.Enabled {
				enabled = "x"
			} else {
				enabled = " "
			}

			if !elem.isLocal {
				enabled = style(enabled, m.styles.greener)
			} else {
				enabled = style(enabled, m.styles.toEnter)
			}

			toPad := absolutePad - len(elem.Name)
			line := fmt.Sprintf("%s%*s[%s]  [%s] ", style(elem.Name, m.styles.greener), toPad, "", synced, enabled)
			if !elem.syncOK {
				line = style(dropAnsiCodes(line), m.styles.mistakes)
			}

			view += wrapWithCursor(int(lower)+idx+len(defaultConfig().EmbededWordLists) == state.cursor, line, m.styles.runningTimer)
			view += "\n"
		}

		help := style("s sync/delete, e enable/disable, ctrl+q to menu", m.styles.toEnter)
		cursorWidget = lipgloss.NewStyle().Align(lipgloss.Left).Render(cursorWidget)
		help = lipgloss.NewStyle().Align(lipgloss.Center).Padding(1).Render(help)
		view = lipgloss.NewStyle().Align(lipgloss.Left).Render(view)

		all := lipgloss.JoinVertical(lipgloss.Center, view, help)

		return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, all)

	case TimerBasedTestResults:
		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
    deltaWpm := "Δavg: " + style(fmt.Sprintf("%s%.2f%%", plusIfPositive(state.results.deltaWpm), math.Min(state.results.deltaWpm, 100.0)), m.styles.greener)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		miscStatsLine1 := fmt.Sprintf("%s %s %s %s", accuracy, deltaWpm, rawWpmShow, givenTime)
		miscStatsLine2 := words

		miscStatsLine1Len := len(dropAnsiCodes(miscStatsLine1))
		plotData := append(state.wpmEachSecond, float64(state.results.wpm))
		wpmsPlot := plotWpms(plotData, miscStatsLine1Len-2)

		fullParagraph := lipgloss.JoinVertical(lipgloss.Center, resultsStyle.Padding(1).Render(wpm), wpmsPlot, resultsStyle.Padding(0).Render(miscStatsLine1), resultsStyle.Render(miscStatsLine2))
		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)

	case WordCountTestResults:
		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		//cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
    deltaWpm := "Δavg: " + style(fmt.Sprintf("%s%.2f%%", plusIfPositive(state.results.deltaWpm), math.Min(state.results.deltaWpm, 100.0)), m.styles.greener)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		wordCnt := "cnt: " + style(strconv.Itoa(state.wordCnt), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		miscStatsLine1 := fmt.Sprintf("%s %s %s %s", accuracy, deltaWpm, rawWpmShow, givenTime)
		miscStatsLine2 := wordCnt + " " + words

		miscStatsLine1Len := len(dropAnsiCodes(miscStatsLine1))

		plotData := append(state.wpmEachSecond, float64(state.results.wpm))
		wpmsPlot := plotWpms(plotData, miscStatsLine1Len-2)

		fullParagraph := lipgloss.JoinVertical(lipgloss.Center, resultsStyle.Padding(1).Render(wpm), wpmsPlot, resultsStyle.Padding(0).Render(miscStatsLine1), resultsStyle.Render(miscStatsLine2))
		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)

	case SentenceCountTestResults:
		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		//cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
    deltaWpm := "Δavg: " + style(fmt.Sprintf("%s%.2f%%", plusIfPositive(state.results.deltaWpm), math.Min(state.results.deltaWpm, 100.0)), m.styles.greener)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		sentenceCnt := "cnt: " + style(strconv.Itoa(state.sentenceCnt), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "sentences: " + style(state.results.wordList, m.styles.greener)

		miscStatsLine1 := fmt.Sprintf("%s %s %s %s", accuracy, deltaWpm, rawWpmShow, givenTime)
		miscStatsLine2 := sentenceCnt + " " + words

		miscStatsLine1Len := len(dropAnsiCodes(miscStatsLine1))
		plotData := append(state.wpmEachSecond, float64(state.results.wpm))
		wpmsPlot := plotWpms(plotData, miscStatsLine1Len-2)

		fullParagraph := lipgloss.JoinVertical(lipgloss.Center, resultsStyle.Padding(1).Render(wpm), wpmsPlot, resultsStyle.Padding(0).Render(miscStatsLine1), resultsStyle.Render(miscStatsLine2))
		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)

	case TimerBasedTest:
		var coloredTimer string
		if state.timer.isRunning {
			coloredTimer = style(state.timer.timer.View(), m.styles.runningTimer)
		} else {
			coloredTimer = style(state.timer.timer.View(), m.styles.stoppedTimer)
		}

		paragraphView := state.base.paragraphView(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(lines, state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		s += positionVerticaly(termHeight)
		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

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

		s += positionVerticaly(termHeight)
		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

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
		cursorLine := findCursorLine(lines, state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

		s += positionVerticaly(termHeight)
		s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.stopwatch.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

	}

	return s
}

func plusIfPositive(f float64) string {
  if f > 0.0 {
    return "+"
  } else {
    return ""
  }
}

func positionVerticaly(termHeight int) string {
	var acc strings.Builder

	for i := 0; i < termHeight/2-3; i++ {
		acc.WriteRune('\n')
	}

	return acc.String()
}

func plotWpms(wpms []float64, width int) string {
	wpmGraph := asciigraph.Plot(
		wpms,
		asciigraph.Precision(0),
		asciigraph.Height(5),
		asciigraph.Width(width),
		asciigraph.CaptionColor(2),
		asciigraph.LabelColor(2),
	)

	return lipgloss.NewStyle().Padding(1).Render(wpmGraph)
}

func averageLineLenFast(lines []string) int {
	linesLen := len(lines)
	linesToConsider := int(math.Min(float64(linesLen), 3))
	return averageStringLen(lines[:linesToConsider])
}

func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1] //Drop last line, as it might skew up average length
	}

	return averageStringLen(lines)
}

func (selection TimerBasedTestSettings) show(styles Styles) string {
	var wordListSelection string
	if selection.enabled {
		wordListSelection = selection.wordListSelections[selection.wordListCursor].name
	} else {
		wordListSelection = "no wordlist enabled"
	}

	selections := []string{selection.timeSelections[selection.timeCursor].String(), wordListSelection}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Timer run", selectionsStr)
}

func (selection WordCountBasedTestSettings) show(styles Styles) string {
	var wordListSelection string
	if selection.enabled {
		wordListSelection = selection.wordListSelections[selection.wordListCursor].name
	} else {
		wordListSelection = "no wordlist enabled"
	}

	selections := []string{fmt.Sprint(selection.wordCountSelections[selection.wordCountCursor]), wordListSelection}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Word count run", selectionsStr)
}

func (selection SentenceCountBasedTestSettings) show(styles Styles) string {
	var wordListSelection string
	if selection.enabled {
		wordListSelection = selection.sentenceListSelections[selection.sentenceListCursor].name
	} else {
		wordListSelection = "no wordlist enabled"
	}
	selections := []string{fmt.Sprint(selection.sentenceCountSelections[selection.sentenceCountCursor]), wordListSelection}
	selectionsStr := showSelections(selections, selection.cursor, styles)
	return fmt.Sprintf("%s %s", "Sentence count run", selectionsStr)
}

func (selection ConfigViewSelection) show(styles Styles) string {
	return "Config "
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

	var coloredInput strings.Builder

	if len(mistakes) == 0 {

		coloredInput.WriteString(styleAllRunes(base.inputBuffer, styles.correct))

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := base.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := base.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput.WriteString(styleAllRunes(sliceUntilMistake, styles.correct))
			coloredInput.WriteString(style(string(mistakeSlice), styles.mistakes))

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := base.inputBuffer[previousMistake+1:]
		coloredInput.WriteString(styleAllRunes(inputAfterLastMistake, styles.correct))
	}

	return coloredInput.String()
}

func (base TestBase) colorCursor(styles Styles) string {
	cursorLetter := base.wordsToEnter[len(base.inputBuffer) : len(base.inputBuffer)+1]

	return style(string(cursorLetter), styles.cursor)
}

func (base TestBase) colorWordsToEnter(styles Styles) string {
	wordsToEnter := base.wordsToEnter[len(base.inputBuffer)+1:] // without cursor

	return style(string(wordsToEnter), styles.toEnter)
}

func wrapStyledParagraph(paragraph string, lineLimit int) string {
	// XXX: Replace spaces, because wordwrap trims them out at the ends
	paragraph = strings.ReplaceAll(paragraph, " ", "·")

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'·'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))
	f.Close()

	paragraph = strings.ReplaceAll(f.String(), "·", " ")

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
	var acc strings.Builder

	for idx, char := range runes {
		_ = idx
		acc.WriteString(style(string(char)).String())
	}

	return acc.String()
}
