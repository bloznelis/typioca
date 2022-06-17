package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	RootCmd = &cobra.Command{
		Use:  "typioca",
		Long: "typioca is a typing test program.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Println("typioca ", Version)
				return nil
			} else {
				termenv.SetWindowTitle("typioca")
				defer println("bye!")

				termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
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
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show typioca version")
	RootCmd.AddCommand(serveCmd)
}
