package client

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"ecommerce-cli/internal/models"
)

func (m Model) updatePayment(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, m.goTo(ScreenCart)
	}

	var cmd tea.Cmd
	m.paymentForm, cmd = applyFormUpdate(m.paymentForm, msg)

	if m.paymentForm != nil && m.paymentForm.State == huh.StateCompleted {
		if len(m.cart.Items) == 0 {
			m.statusErr = "Rien a payer"
			return m, m.goTo(ScreenCart)
		}

		userID := 0
		if m.session.User != nil {
			userID = m.session.User.ID
		}

		order := models.Order{
			ID:        2000 + len(m.orders) + 1,
			UserID:    userID,
			Total:     cartTotalTTC(m.cart),
			Status:    models.StatusPaid,
			CreatedAt: time.Now(),
		}
		for _, it := range m.cart.Items {
			price := 0.0
			if it.Product != nil {
				price = it.Product.Price
			}
			order.Items = append(order.Items, models.OrderItem{
				ProductID:   it.ProductID,
				Quantity:    it.Quantity,
				PriceAtTime: price,
			})
		}

		m.orders = append([]models.Order{order}, m.orders...)
		m.cart.Items = nil
		m.cartCursor = 0
		m.ordersCursor = 0
		m.screen = ScreenOrders
		m.clearStatus()
		m.statusOK = fmt.Sprintf("Paiement OK - commande #%d [%s]", order.ID, order.Status)
		return m, nil
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
		subtitleStyle.Render(fmt.Sprintf("Montant TTC: %.2f EUR", cartTotalTTC(m.cart))) + "\n\n" +
		body + "\n" +
		helpStyle.Render("Validation locale: 13-19 chiffres, MM/YY, CVC 3-4") + "\n"
}
