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
		m.loading = true
		return m, m.cmdConfirm(email, code)
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
	return screenHeader("Confirmation") +
		subtitleStyle.Render(fmt.Sprintf("POST /confirm | code mock : %s", mockConfirmCode)) + "\n\n" +
		body
}
