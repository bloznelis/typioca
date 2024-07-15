package cmd

import (
	"time"

	"github.com/bloznelis/typioca/cmd/words"
	"github.com/charmbracelet/bubbles/stopwatch"
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

type myStopWatch struct {
	stopwatch stopwatch.Model
	isRunning bool
}

type mistakes struct {
	mistakesAt     map[int]bool
	rawMistakesCnt int // Should never be reduced
}

type StringStyle func(string) termenv.Style

type Styles struct {
	correct      StringStyle
	toEnter      StringStyle
	mistakes     StringStyle
	cursor       StringStyle
	runningTimer StringStyle
	stoppedTimer StringStyle
	greener      StringStyle
	faintGreen   StringStyle
}

type model struct {
	state  State
	styles Styles
	width  int
	height int
}

type Results struct {
	identifier    ResultsIdentifier
	wpm           int
	accuracy      float64
	deltaWpm      float64
	rawWpm        int
	cpm           int
	time          time.Duration
	wordList      string
	wpmEachSecond []float64
}

type PersistentResults struct {
	Results AllPersistedResults
	Version int
}

type ResultsIdentifier struct {
	testType TestType
	numeric  NumericSetting
	words    WordListName
}

type PersistentResultsNode struct {
	Wpm           int
	Accuracy      float64
	DeltaWpm      float64
	RawWpm        int
	Cpm           int
	WpmEachSecond []float64
}

type WordListSelection struct {
	key  string
	show string
}

type State interface{}

type MainMenuSelection interface {
	Enabled() bool
	show(s Styles) string
	handleInput(msg tea.Msg, menu MainMenu) State
}

type TimerBasedTestSettings struct {
	timeSelections     []time.Duration
	timeCursor         int
	wordListSelections []WordsSelection
	wordListCursor     int
	cursor             int
	enabled            bool
}

func (s TimerBasedTestSettings) Enabled() bool {
	return s.enabled
}

type WordCountBasedTestSettings struct {
	wordCountSelections []int
	wordCountCursor     int
	wordListSelections  []WordsSelection
	wordListCursor      int
	cursor              int
	enabled             bool
}

func (s WordCountBasedTestSettings) Enabled() bool {
	return s.enabled
}

type SentenceCountBasedTestSettings struct {
	sentenceCountSelections []int
	sentenceCountCursor     int
	sentenceListSelections  []WordsSelection
	sentenceListCursor      int
	cursor                  int
	enabled                 bool
}

func (s SentenceCountBasedTestSettings) Enabled() bool {
	return s.enabled
}

type ConfigViewSelection struct{}

func (s ConfigViewSelection) Enabled() bool {
	return true
}

type MainMenu struct {
	config                 Config
	selections             []MainMenuSelection
	cursor                 int
	timeBasedGenerator     words.WordsGenerator
	wordCountGenerator     words.WordsGenerator
	sentenceCountGenerator words.WordsGenerator
}

type TestBase struct {
	wordsToEnter  []rune
	inputBuffer   []rune
	wpmEachSecond []float64
	rawInputCnt   int // Should not be reduced
	mistakes      mistakes
	cursor        int
}

type TimerBasedTest struct {
	settings  TimerBasedTestSettings
	timer     myTimer
	base      TestBase
	completed bool
	mainMenu  MainMenu
}

type TimerBasedTestResults struct {
	settings      TimerBasedTestSettings
	wpmEachSecond []float64
	results       Results
	mainMenu      MainMenu
}

type WordCountBasedTest struct {
	settings  WordCountBasedTestSettings
	stopwatch myStopWatch
	base      TestBase
	completed bool
	mainMenu  MainMenu
}

type WordCountTestResults struct {
	settings      WordCountBasedTestSettings
	wpmEachSecond []float64
	wordCnt       int
	results       Results
	mainMenu      MainMenu
}

type SentenceCountBasedTest struct {
	settings  SentenceCountBasedTestSettings
	stopwatch myStopWatch
	base      TestBase
	completed bool
	mainMenu  MainMenu
}

type SentenceCountTestResults struct {
	settings      SentenceCountBasedTestSettings
	wpmEachSecond []float64
	sentenceCnt   int
	results       Results
	mainMenu      MainMenu
}

type ConfigView struct {
	mainMenu MainMenu
	config   Config
	cursor   int
}

type TestSettingCursors struct {
	TimerTimeCursor     int
	TimerWordlistCursor int

	WordCountCursor         int
	WordCountWordlistCursor int

	SentenceCountCursor         int
	SentenceCountWordlistCursor int
}

type LayoutFile struct {
	Name      string
	Path      string
	RemoteURI string
	synced    bool
}

type Layout struct {
	Name     string        `json:"name"`
	Mappings map[rune]rune `json:"mappings"`
}

type WordList struct {
	Sentences bool
	Name      string
	Path      string
	RemoteURI string
	isLocal   bool
	Enabled   bool
	synced    bool
	syncOK    bool
}

func (wordList *WordList) toggleEnabled() {
	if !wordList.isLocal {
		wordList.Enabled = !wordList.Enabled
	}
}

type EmbededWordList struct {
	Name        string
	IsSentences bool
	Enabled     bool
}

func (embeded *EmbededWordList) toggleEnabled() {
	embeded.Enabled = !embeded.Enabled
}

type Config struct {
	TestSettingCursors TestSettingCursors
	EmbededWordLists   []EmbededWordList
	WordLists          []WordList
	LayoutFiles        []LayoutFile
	Layout             Layout
	Version            int
}

type LocalConfig struct {
	Words []WordList
}

func (cfg Config) configTotalSelectionsCount() int {
	return len(cfg.WordLists) + len(cfg.EmbededWordLists) + len(cfg.LayoutFiles)
}

type Toggleable interface {
	toggle()
}
