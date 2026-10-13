package client

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"ecommerce-cli/internal/models"
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
		for _, acc := range m.accounts {
			if strings.EqualFold(acc.User.Email, email) || strings.EqualFold(acc.User.Username, username) {
				m.statusErr = "Email ou username deja pris"
				m.registerForm = newRegisterForm()
				return m, m.registerForm.Init()
			}
		}

		acc := mockAccount{
			User: models.User{
				ID:               m.nextID,
				Username:         username,
				Email:            email,
				IsConfirmed:      false,
				ConfirmationCode: mockConfirmCode,
				CreatedAt:        time.Now(),
			},
			Password: password,
		}
		m.nextID++
		m.accounts = append(m.accounts, acc)

		cmd := m.goMenu()
		m.statusOK = fmt.Sprintf("Compte cree, code = %s", mockConfirmCode)
		return m, cmd
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
	return screenHeader("Register") +
		subtitleStyle.Render("Inscription (mock local)") + "\n\n" +
		body
}
