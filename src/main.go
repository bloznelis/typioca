package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	serverBind    = ""
	serverPort    = 2229
	serverKeyPath = ""
	showVersion   = false
)

var (
	Version = "dev"
	rootCmd = &cobra.Command{
		Use:  "typioca",
		Long: "typioca is a typing test program.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Println("typioca ", Version)
				return nil
			} else {
				termenv.SetWindowTitle("typioca")
				defer println("bye!")

				termWidth, termHeight, _ := term.GetSize(0)
				p := tea.NewProgram(
					initialModel(
						termenv.ColorProfile(),
						termenv.ForegroundColor(),
						termWidth,
						termHeight,
					),
					tea.WithAltScreen(),
				)

				return p.Start()
			}
		},
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve the typioca server",
		Long:  "serve starts the typioca SSH server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", serverBind, serverPort)),
				wish.WithHostKeyPath(serverKeyPath),
				wish.WithMiddleware(
					lm.Middleware(),
					bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
						pty, _, active := s.Pty()
						if !active {
							wish.Fatal(s, fmt.Errorf("not a tty"))
							return nil, nil
						}
						return initialModel(
								termenv.ANSI256,
								termenv.ANSIWhite,
								pty.Window.Width,
								pty.Window.Height,
							),
							[]tea.ProgramOption{tea.WithAltScreen()}
					}),
				),
			)
			if err != nil {
				return err
			}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			log.Printf("Starting server on %s:%d", serverBind, serverPort)
			go func() {
				if err := s.ListenAndServe(); err != nil {
					log.Fatalln(err)
				}
			}()

			<-done

			log.Printf("Stopping SSH server on %s:%d", serverBind, serverPort)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer func() { cancel() }()
			if err := s.Shutdown(ctx); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	serveCmd.Flags().StringVarP(&serverKeyPath, "key", "k", "typioca", "path to the server key")
	serveCmd.Flags().StringVarP(&serverBind, "bind", "b", "", "address to bind on")
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 2229, "port to serve on")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show typioca version")
	rootCmd.AddCommand(serveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initTimerBasedTest(settings TimerBasedTestSettings) TimerBasedTest {
	return TimerBasedTest{
		settings: settings,
		timer: myTimer{
			timer:     timer.NewWithInterval(settings.timeSelections[settings.timeCursor], time.Second),
			duration:  settings.timeSelections[settings.timeCursor],
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter: NewGenerator().Generate(settings.wordListSelections[settings.wordListCursor].key),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initWordCountBasedTest(settings WordCountBasedTestSettings) WordCountBasedTest {
	generator := NewGenerator()
	generator.Count = settings.wordCountSelections[settings.wordCountCursor]
	return WordCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: strings.TrimSpace(generator.Generate(settings.wordListSelections[settings.wordListCursor].key)),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initSentenceCountBasedTest(settings SentenceCountBasedTestSettings) SentenceCountBasedTest {
	generator := NewGenerator()
	generator.Count = 40
	generator.Count = settings.sentenceCountSelections[settings.sentenceCountCursor]
	return SentenceCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: strings.TrimSpace(generator.Generate(settings.sentenceListSelections[settings.sentenceListCursor].key)),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initTimerBasedTestSelection() TimerBasedTestSettings {
	return TimerBasedTestSettings{
		timeSelections: []time.Duration{time.Second * 120, time.Second * 60, time.Second * 30, time.Second * 15},
		timeCursor:     2,
		wordListSelections: []WordListSelection{
			{
				key:  "dorian-gray-words",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-words",
				show: "frankenstein",
			},
			{
				key:  "common-words",
				show: "common-words",
			},
			{
				key:  "pride-and-prejudice-words",
				show: "pride-and-prejudice",
			},
			{
				key:  "dorian-gray-sentences",
				show: "dorian-gray-sentences",
			},
			{
				key:  "frankenstein-sentences",
				show: "frankenstein-sentences",
			},
			{
				key:  "pride-and-prejudice-sentences",
				show: "pride-and-prejudice-sentences",
			},
		},
		wordListCursor: 2,
		cursor:         0,
	}
}

func initWordCountBasedTestSelection() WordCountBasedTestSettings {
	return WordCountBasedTestSettings{
		wordCountSelections: []int{100, 50, 25, 10},
		wordCountCursor:     2,
		wordListSelections: []WordListSelection{
			{
				key:  "dorian-gray-words",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-words",
				show: "frankenstein",
			},
			{
				key:  "common-words",
				show: "common-words",
			},
			{
				key:  "pride-and-prejudice-words",
				show: "pride-and-prejudice",
			},
		},
		wordListCursor: 2,
		cursor:         0,
	}
}

func initSentenceCountBasedTestSelection() SentenceCountBasedTestSettings {
	return SentenceCountBasedTestSettings{
		sentenceCountSelections: []int{30, 15, 5, 1},
		sentenceCountCursor:     2,
		sentenceListSelections: []WordListSelection{
			{
				key:  "dorian-gray-sentences",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-sentences",
				show: "frankenstein",
			},
			{
				key:  "pride-and-prejudice-sentences",
				show: "pride-and-prejudice",
			},
		},
		sentenceListCursor: 1,
		cursor:             0,
	}
}

func initMainMenu() MainMenu {
	return MainMenu{
		selections: []MainMenuSelection{
			initTimerBasedTestSelection(),
			initWordCountBasedTestSelection(),
			initSentenceCountBasedTestSelection(),
		},
		cursor: 0,
	}
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink, //we probably don't need this anymore
	)
}
