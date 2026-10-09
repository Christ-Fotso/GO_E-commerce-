package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return m, nil
}

func (m Model) viewConfirm() string {
	return screenHeader("Confirm") +
		subtitleStyle.Render("Confirmation de compte (squelette)") + "\n\n" +
		fmt.Sprintf("Formulaire huh pret: email / code\nCode mock attendu: %s\n", mockConfirmCode)
}
