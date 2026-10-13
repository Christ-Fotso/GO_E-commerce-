package client

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) updateReset(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goMenu()
	}

	var cmd tea.Cmd
	m.resetForm, cmd = applyFormUpdate(m.resetForm, msg)

	if m.resetForm != nil && m.resetForm.State == huh.StateCompleted {
		email := strings.TrimSpace(m.resetForm.GetString("email"))
		code := strings.TrimSpace(m.resetForm.GetString("code"))
		newPass := m.resetForm.GetString("new_password")
		if email == "" {
			m.statusErr = "Email requis"
			m.resetForm = newResetForm()
			return m, m.resetForm.Init()
		}
		m.loading = true
		return m, m.cmdReset(email, code, newPass)
	}

	if m.resetForm != nil && m.resetForm.State == huh.StateAborted {
		return m, m.goMenu()
	}

	return m, cmd
}

func (m Model) viewReset() string {
	var body string
	if m.resetForm != nil {
		body = m.resetForm.View()
	}
	return screenHeader("Réinitialiser le mot de passe") +
		subtitleStyle.Render("POST /password/forgot puis /password/reset (API.md)") + "\n\n" +
		body
}
