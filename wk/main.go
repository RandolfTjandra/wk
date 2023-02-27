package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel(IndexView), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting: %v", err)
		os.Exit(1)
	}
}
