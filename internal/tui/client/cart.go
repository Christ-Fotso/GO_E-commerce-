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
			it := m.cart.Items[m.cartCursor]
			m.loading = true
			return m, m.cmdUpdateCartQty(it.ProductID, it.Quantity+1)
		}
	case "-", "_":
		if len(m.cart.Items) > 0 && m.cart.Items[m.cartCursor].Quantity > 1 {
			it := m.cart.Items[m.cartCursor]
			m.loading = true
			return m, m.cmdUpdateCartQty(it.ProductID, it.Quantity-1)
		}
	case "d":
		if len(m.cart.Items) > 0 {
			it := m.cart.Items[m.cartCursor]
			m.loading = true
			return m, m.cmdDeleteCartItem(it.ProductID)
		}
	case "r":
		m.loading = true
		return m, m.cmdLoadCart("")
	case "p", "enter":
		if len(m.cart.Items) == 0 {
			m.statusErr = "Panier vide"
			return m, nil
		}
		return m, m.goTo(ScreenPayment)
	}
	return m, nil
}

func (m Model) viewCart() string {
	var b strings.Builder
	b.WriteString(screenHeader("Panier"))
	b.WriteString(subtitleStyle.Render("GET/PUT/DELETE /cart (API.md)"))
	b.WriteString("\n\n")

	if len(m.cart.Items) == 0 {
		b.WriteString(subtitleStyle.Render("Panier vide — ajoute des produits depuis le catalogue"))
		b.WriteString("\n")
		return b.String()
	}

	for i, it := range m.cart.Items {
		name := fmt.Sprintf("produit #%d", it.ProductID)
		lineTotal := it.LineTotalTTC
		if it.Product != nil {
			name = it.Product.Name
			if lineTotal == 0 {
				lineTotal = productUnitTTC(it.Product) * float64(it.Quantity)
			}
		}
		line := fmt.Sprintf("%s  x%d  = %.2f EUR TTC", name, it.Quantity, lineTotal)
		if i == m.cartCursor {
			b.WriteString(selectedItemStyle.Render("> "+line) + "\n")
		} else {
			b.WriteString(itemStyle.Render("  "+line) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render(fmt.Sprintf("Total TTC: %.2f EUR", cartTotalTTC(m.cart))))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("\n+/- qté | d supprimer | r refresh | p/enter payer | esc menu"))
	b.WriteString("\n")
	return b.String()
}
