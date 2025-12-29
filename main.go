package main

import (
	"fmt"
	"os"

	"diff-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewModel()

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
