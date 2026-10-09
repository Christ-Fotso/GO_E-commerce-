package client

import tea "github.com/charmbracelet/bubbletea"

func (m Model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return m, nil
}

func (m Model) viewLogin() string {
	return screenHeader("Login") +
		subtitleStyle.Render("Connexion (squelette)") + "\n\n" +
		"Formulaire huh pret: email / password\n" +
		"Apres login mock: session.User renseignee.\n"
}
