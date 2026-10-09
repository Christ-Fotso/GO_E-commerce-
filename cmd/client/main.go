package main

import (
	"fmt"
	"os"

	"ecommerce-cli/internal/tui/client"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(client.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "erreur TUI client: %v\n", err)
		os.Exit(1)
	}
}
