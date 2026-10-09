package client

import tea "github.com/charmbracelet/bubbletea"

func (m Model) updateReset(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return m, nil
}

func (m Model) viewReset() string {
	return screenHeader("Reset password") +
		subtitleStyle.Render("Reinitialisation (squelette)") + "\n\n" +
		"Formulaire huh pret: email\n"
}
