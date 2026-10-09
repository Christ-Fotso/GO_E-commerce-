package client

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateCart(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up", "k":
		if m.cartCursor > 0 {
			m.cartCursor--
		}
	case "down", "j":
		if m.cartCursor < len(m.cart.Items)-1 {
			m.cartCursor++
		}
	case "+", "=":
		if len(m.cart.Items) > 0 {
			m.cart.Items[m.cartCursor].Quantity++
		}
	case "-", "_":
		if len(m.cart.Items) > 0 && m.cart.Items[m.cartCursor].Quantity > 1 {
			m.cart.Items[m.cartCursor].Quantity--
		}
	case "d":
		if len(m.cart.Items) > 0 {
			i := m.cartCursor
			m.cart.Items = append(m.cart.Items[:i], m.cart.Items[i+1:]...)
			if m.cartCursor >= len(m.cart.Items) && m.cartCursor > 0 {
				m.cartCursor--
			}
		}
	case "p":
		if len(m.cart.Items) == 0 {
			m.statusErr = "Panier vide"
			return m, nil
		}
		m.goTo(ScreenPayment)
	}
	return m, nil
}

func (m Model) viewCart() string {
	var b strings.Builder
	b.WriteString(screenHeader("Panier"))

	if len(m.cart.Items) == 0 {
		b.WriteString(subtitleStyle.Render("Panier vide - ajoute des produits depuis le catalogue"))
		b.WriteString("\n")
		return b.String()
	}

	for i, it := range m.cart.Items {
		name := fmt.Sprintf("produit #%d", it.ProductID)
		price := 0.0
		if it.Product != nil {
			name = it.Product.Name
			price = it.Product.Price
		}
		line := fmt.Sprintf("%s  x%d  (%.2f EUR)", name, it.Quantity, price*float64(it.Quantity))
		if i == m.cartCursor {
			b.WriteString(selectedItemStyle.Render("> "+line) + "\n")
		} else {
			b.WriteString(itemStyle.Render("  "+line) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render(fmt.Sprintf("Total TTC: %.2f EUR", cartTotalTTC(m.cart))))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("\n+/- qte | d supprimer | p payer | esc menu"))
	b.WriteString("\n")
	return b.String()
}
