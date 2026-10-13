package client

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goMenu()
	}

	var cmd tea.Cmd
	m.loginForm, cmd = applyFormUpdate(m.loginForm, msg)

	if m.loginForm != nil && m.loginForm.State == huh.StateCompleted {
		email := strings.TrimSpace(m.loginForm.GetString("email"))
		password := m.loginForm.GetString("password")
		m.loading = true
		return m, m.cmdLogin(email, password)
	}

	if m.loginForm != nil && m.loginForm.State == huh.StateAborted {
		return m, m.goMenu()
	}

	return m, cmd
}

func (m Model) viewLogin() string {
	var body string
	if m.loginForm != nil {
		body = m.loginForm.View()
	}
	return screenHeader("Connexion") +
		subtitleStyle.Render("POST /login | mock : demo@example.com / demo123") + "\n\n" +
		body
}
