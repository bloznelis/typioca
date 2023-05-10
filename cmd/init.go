package cmd

import (
	"time"

	"github.com/bloznelis/typioca/cmd/words"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

// todo: clean these up. Maybe we could reuse filtering by enabled and synce, because now it's redundant
type WordsSelection struct {
	name         string
	generatorKey string
}

func filterEnabledWordSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && !elem.Sentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}
	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && !elem.IsSentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

func filterEnabledSentenceSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && elem.Sentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && elem.IsSentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}
	return acc
}

func filterEnabledSelections(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

func initTimerBasedTest(settings TimerBasedTestSettings, mainMenu MainMenu) TimerBasedTest {
	return TimerBasedTest{
		settings: settings,
		timer: myTimer{
			timer:     timer.NewWithInterval(settings.timeSelections[settings.timeCursor], time.Second),
			duration:  settings.timeSelections[settings.timeCursor],
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter: mainMenu.timeBasedGenerator.Generate(settings.wordListSelections[settings.wordListCursor].generatorKey),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initWordCountBasedTest(settings WordCountBasedTestSettings, mainMenu MainMenu) WordCountBasedTest {
	mainMenu.wordCountGenerator.Count = settings.wordCountSelections[settings.wordCountCursor]
	return WordCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: mainMenu.wordCountGenerator.Generate(settings.wordListSelections[settings.wordListCursor].generatorKey),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initSentenceCountBasedTest(settings SentenceCountBasedTestSettings, mainMenu MainMenu) SentenceCountBasedTest {
	mainMenu.sentenceCountGenerator.Count = settings.sentenceCountSelections[settings.sentenceCountCursor]
	return SentenceCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: mainMenu.sentenceCountGenerator.Generate(settings.sentenceListSelections[settings.sentenceListCursor].generatorKey),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initTestSettingCursors() TestSettingCursors {
	return TestSettingCursors{
		TimerTimeCursor:             2,
		TimerWordlistCursor:         0,
		WordCountCursor:             2,
		WordCountWordlistCursor:     0,
		SentenceCountCursor:         2,
		SentenceCountWordlistCursor: 0,
	}
}

func (cursors *TestSettingCursors) resetWordlistCursors() {
    cursors.TimerWordlistCursor = 0
    cursors.WordCountWordlistCursor = 0
    cursors.SentenceCountWordlistCursor = 0
}

func initTimerBasedTestSettings(config Config, words []WordsSelection) TimerBasedTestSettings {
	return TimerBasedTestSettings{
		timeSelections:     []time.Duration{time.Second * 120, time.Second * 60, time.Second * 30, time.Second * 15},
		timeCursor:         config.TestSettingCursors.TimerTimeCursor,
		wordListSelections: words,
		wordListCursor:     config.TestSettingCursors.TimerWordlistCursor,
		cursor:             0,
		enabled:            len(words) > 0,
	}
}

func initWordCountBasedTestSettings(config Config, words []WordsSelection) WordCountBasedTestSettings {
	return WordCountBasedTestSettings{
		wordCountSelections: []int{100, 50, 25, 10},
		wordCountCursor:     config.TestSettingCursors.WordCountCursor,
		wordListSelections:  words,
		wordListCursor:      config.TestSettingCursors.WordCountWordlistCursor,
		cursor:              0,
		enabled:             len(words) > 0,
	}
}

func initSentenceCountBasedTestSettings(config Config, words []WordsSelection) SentenceCountBasedTestSettings {
	return SentenceCountBasedTestSettings{
		sentenceCountSelections: []int{30, 15, 5, 1},
		sentenceCountCursor:     config.TestSettingCursors.SentenceCountCursor,
		sentenceListSelections:  words,
		sentenceListCursor:      config.TestSettingCursors.SentenceCountWordlistCursor,
		cursor:                  0,
		enabled:                 len(words) > 0,
	}
}

func initConfigView(config Config, mainMenu MainMenu) ConfigView {
	configView := ConfigView{
		config:   config,
		mainMenu: mainMenu,
	}
	return configView
}

func initConfigViewSelection() ConfigViewSelection {
	return ConfigViewSelection{}
}

func initMainMenu() MainMenu {
	config := ReadConfig()
	timeBasedWordSelections := filterEnabledSelections(config)
	countBasedWordSelections := filterEnabledWordSelection(config)
	countBasedSentenceSelections := filterEnabledSentenceSelection(config)
	return MainMenu{
		config: config,
		selections: []MainMenuSelection{
			initTimerBasedTestSettings(config, timeBasedWordSelections),
			initWordCountBasedTestSettings(config, countBasedWordSelections),
			initSentenceCountBasedTestSettings(config, countBasedSentenceSelections),
			initConfigViewSelection(),
		},
		cursor:                 0,
		timeBasedGenerator:     words.NewGenerator(paths(timeBasedWordSelections)),
		wordCountGenerator:     words.NewGenerator(paths(countBasedWordSelections)),
		sentenceCountGenerator: words.NewGenerator(paths(countBasedSentenceSelections)),
	}
}

func paths(selections []WordsSelection) []string {
	var acc []string
	for _, elem := range selections {
		// XXX: don't to this at home
		if elem.generatorKey != elem.name {
			acc = append(acc, elem.generatorKey)
		}
	}
	return acc
}

func initialModel(profile termenv.Profile, fore termenv.Color, width, height int) model {
	return model{
		width:  width,
		height: height,
		state:  initMainMenu(),
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
	}
}
