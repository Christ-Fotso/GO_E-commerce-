package client

import (
	"fmt"
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

		var found *mockAccount
		for i := range m.accounts {
			if strings.EqualFold(m.accounts[i].User.Email, email) {
				found = &m.accounts[i]
				break
			}
		}
		if found == nil {
			m.statusErr = "Utilisateur introuvable"
			m.loginForm = newLoginForm()
			return m, m.loginForm.Init()
		}
		if !found.User.IsConfirmed {
			m.statusErr = "Compte non confirme - utilise Confirm"
			m.loginForm = newLoginForm()
			return m, m.loginForm.Init()
		}
		if found.Password != password {
			m.statusErr = "Mot de passe incorrect"
			m.loginForm = newLoginForm()
			return m, m.loginForm.Init()
		}

		user := found.User
		m.session = Session{User: &user, Token: "mock-token"}
		m.cart.UserID = user.ID
		m.orders = mockOrders(user.ID)
		m.refreshMenu()
		cmd := m.goMenu()
		m.statusOK = fmt.Sprintf("Bienvenue %s !", user.Username)
		return m, cmd
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
	return screenHeader("Login") +
		subtitleStyle.Render("Compte demo: demo@example.com / demo123") + "\n\n" +
		body
}
