package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) updatePayment(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goTo(ScreenCart)
	}

	var cmd tea.Cmd
	m.paymentForm, cmd = applyFormUpdate(m.paymentForm, msg)

	if m.paymentForm != nil && m.paymentForm.State == huh.StateCompleted {
		if len(m.cart.Items) == 0 {
			m.statusErr = "Rien à payer"
			return m, m.goTo(ScreenCart)
		}
		card := m.paymentForm.GetString("card")
		expiry := m.paymentForm.GetString("expiry")
		cvc := m.paymentForm.GetString("cvc")
		m.loading = true
		return m, m.cmdPay(card, expiry, cvc)
	}

	if m.paymentForm != nil && m.paymentForm.State == huh.StateAborted {
		return m, m.goTo(ScreenCart)
	}

	return m, cmd
}

func (m Model) viewPayment() string {
	var body string
	if m.paymentForm != nil {
		body = m.paymentForm.View()
	}
	return screenHeader("Paiement") +
		subtitleStyle.Render(fmt.Sprintf("POST /cart/pay | Montant TTC: %.2f EUR", cartTotalTTC(m.cart))) + "\n\n" +
		body + "\n" +
		helpStyle.Render("Carte test API : 4242 4242 4242 4242 / 12/30 / 123") + "\n"
}
