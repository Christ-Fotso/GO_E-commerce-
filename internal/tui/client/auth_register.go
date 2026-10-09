package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateRegister(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Formulaire huh branche a l'etape auth mock.
	_ = msg
	return m, nil
}

func (m Model) viewRegister() string {
	return screenHeader("Register") +
		subtitleStyle.Render("Inscription (squelette)") + "\n\n" +
		fmt.Sprintf("Formulaire huh pret: username / email / password\nCode mock: %s\n", mockConfirmCode)
}
