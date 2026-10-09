package client

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ecommerce-cli/internal/models"
)

func (m Model) updatePayment(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	// Squelette: enter = paiement mock reussi.
	if key.String() == "enter" {
		if len(m.cart.Items) == 0 {
			m.statusErr = "Rien a payer"
			return m, nil
		}

		userID := 0
		if m.session.User != nil {
			userID = m.session.User.ID
		}

		order := models.Order{
			ID:        1000 + len(m.orders) + 1,
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
		m.goTo(ScreenOrders)
		m.statusOK = fmt.Sprintf("Paiement mock OK - commande #%d", order.ID)
	}
	return m, nil
}

func (m Model) viewPayment() string {
	return screenHeader("Paiement") +
		subtitleStyle.Render("Formulaire (squelette): n carte / date / CVC") + "\n\n" +
		fmt.Sprintf("Montant a payer: %.2f EUR\n\n", cartTotalTTC(m.cart)) +
		"Champs a brancher a l'etape mock UI.\n" +
		helpStyle.Render("\nenter payer (mock) | esc menu") + "\n"
}
