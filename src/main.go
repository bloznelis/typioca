package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	OsInit()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
