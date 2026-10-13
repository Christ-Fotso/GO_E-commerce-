package client

import (
	"fmt"
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
		newPass := m.resetForm.GetString("new_password")

		if email == "" {
			m.statusErr = "Email requis"
			m.resetForm = newResetForm()
			return m, m.resetForm.Init()
		}

		for i := range m.accounts {
			if strings.EqualFold(m.accounts[i].User.Email, email) {
				if newPass != "" {
					m.accounts[i].Password = newPass
				}
				break
			}
		}

		cmd := m.goMenu()
		m.statusOK = fmt.Sprintf("Email de reset envoye a %s", email)
		return m, cmd
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
	return screenHeader("Reset password") +
		subtitleStyle.Render("Reinitialisation (mock)") + "\n\n" +
		body
}
