package client

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goMenu()
	}

	var cmd tea.Cmd
	m.confirmForm, cmd = applyFormUpdate(m.confirmForm, msg)

	if m.confirmForm != nil && m.confirmForm.State == huh.StateCompleted {
		email := strings.TrimSpace(m.confirmForm.GetString("email"))
		code := strings.TrimSpace(m.confirmForm.GetString("code"))

		idx := -1
		for i, acc := range m.accounts {
			if strings.EqualFold(acc.User.Email, email) {
				idx = i
				break
			}
		}
		if idx < 0 {
			m.statusErr = "Utilisateur introuvable"
			m.confirmForm = newConfirmForm()
			return m, m.confirmForm.Init()
		}
		if code != mockConfirmCode && code != m.accounts[idx].User.ConfirmationCode {
			m.statusErr = fmt.Sprintf("Code invalide (attendu: %s)", mockConfirmCode)
			m.confirmForm = newConfirmForm()
			return m, m.confirmForm.Init()
		}

		m.accounts[idx].User.IsConfirmed = true
		cmd := m.goMenu()
		m.statusOK = fmt.Sprintf("Compte %s confirme", email)
		return m, cmd
	}

	if m.confirmForm != nil && m.confirmForm.State == huh.StateAborted {
		return m, m.goMenu()
	}

	return m, cmd
}

func (m Model) viewConfirm() string {
	var body string
	if m.confirmForm != nil {
		body = m.confirmForm.View()
	}
	return screenHeader("Confirm") +
		subtitleStyle.Render(fmt.Sprintf("Code mock: %s", mockConfirmCode)) + "\n\n" +
		body
}
