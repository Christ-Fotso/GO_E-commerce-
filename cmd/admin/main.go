package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"ecommerce-cli/internal/tui/admin"
)

func main() {
	p := tea.NewProgram(admin.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur du panel administrateur: %v\n", err)
		os.Exit(1)
	}
}
