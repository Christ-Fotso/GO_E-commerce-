package client

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) updateRegister(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goMenu()
	}

	var cmd tea.Cmd
	m.registerForm, cmd = applyFormUpdate(m.registerForm, msg)

	if m.registerForm != nil && m.registerForm.State == huh.StateCompleted {
		username := strings.TrimSpace(m.registerForm.GetString("username"))
		email := strings.TrimSpace(m.registerForm.GetString("email"))
		password := m.registerForm.GetString("password")

		if username == "" || email == "" || password == "" {
			m.statusErr = "Tous les champs sont requis"
			m.registerForm = newRegisterForm()
			return m, m.registerForm.Init()
		}
		m.loading = true
		m.statusOK = "Inscription en cours..."
		return m, m.cmdRegister(username, email, password)
	}

	if m.registerForm != nil && m.registerForm.State == huh.StateAborted {
		return m, m.goMenu()
	}

	return m, cmd
}

func (m Model) viewRegister() string {
	var body string
	if m.registerForm != nil {
		body = m.registerForm.View()
	}
	return screenHeader("Inscription") +
		subtitleStyle.Render("POST /register (API.md) + fallback mock") + "\n\n" +
		body
}
