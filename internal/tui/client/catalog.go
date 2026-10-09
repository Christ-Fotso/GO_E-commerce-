package client

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ecommerce-cli/internal/models"
)

func (m Model) updateCatalog(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	k := strings.ToLower(key.String())

	switch m.screen {
	case ScreenCatalog:
		switch k {
		case "up", "k":
			if m.catalogCursor > 0 {
				m.catalogCursor--
			}
		case "down", "j":
			if m.catalogCursor < len(m.products)-1 {
				m.catalogCursor++
			}
		case "enter":
			if len(m.products) > 0 {
				p := m.products[m.catalogCursor]
				m.selectedProduct = &p
				m.goTo(ScreenProduct)
			}
		case "a":
			// Raccourci: ajouter depuis la liste, sans ouvrir la fiche.
			if len(m.products) > 0 {
				p := m.products[m.catalogCursor]
				m.addToCart(p, 1)
				m.statusOK = fmt.Sprintf("Ajoute au panier: %s (panier: %d)", p.Name, len(m.cart.Items))
			}
		}
	case ScreenProduct:
		switch k {
		case "b", "backspace":
			m.goTo(ScreenCatalog)
		case "a", "enter":
			if m.selectedProduct != nil {
				m.addToCart(*m.selectedProduct, 1)
				m.statusOK = fmt.Sprintf("Ajoute au panier: %s (panier: %d)", m.selectedProduct.Name, len(m.cart.Items))
			}
		}
	}
	return m, nil
}

func (m *Model) addToCart(p models.Product, qty int) {
	for i := range m.cart.Items {
		if m.cart.Items[i].ProductID == p.ID {
			m.cart.Items[i].Quantity += qty
			return
		}
	}
	prod := p
	m.cart.Items = append(m.cart.Items, models.CartItem{
		ID:        len(m.cart.Items) + 1,
		CartID:    m.cart.ID,
		ProductID: p.ID,
		Quantity:  qty,
		Product:   &prod,
	})
}

func (m Model) viewCatalog() string {
	if m.screen == ScreenProduct {
		return m.viewProduct()
	}

	var b strings.Builder
	b.WriteString(screenHeader("Catalogue / Recherche"))
	b.WriteString(subtitleStyle.Render("enter = fiche produit | a = ajouter au panier"))
	b.WriteString("\n\n")

	for i, p := range m.products {
		line := fmt.Sprintf("%s - %.2f EUR (stock %d)", p.Name, p.Price, p.Stock)
		if i == m.catalogCursor {
			b.WriteString(selectedItemStyle.Render("> "+line) + "\n")
		} else {
			b.WriteString(itemStyle.Render("  "+line) + "\n")
		}
	}

	b.WriteString(helpStyle.Render("\nup/down | enter fiche | a panier | esc menu"))
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewProduct() string {
	var b strings.Builder
	b.WriteString(screenHeader("Fiche produit"))

	if m.selectedProduct == nil {
		b.WriteString(errorStyle.Render("Aucun produit selectionne") + "\n")
		return b.String()
	}

	p := m.selectedProduct
	b.WriteString(sectionStyle.Render(p.Name) + "\n")
	b.WriteString(p.Description + "\n\n")
	b.WriteString(fmt.Sprintf("Prix: %.2f EUR\n", p.Price))
	b.WriteString(fmt.Sprintf("Stock: %d\n", p.Stock))
	b.WriteString(fmt.Sprintf("Articles dans le panier: %d\n", len(m.cart.Items)))
	b.WriteString(helpStyle.Render("\na ou enter = ajouter | b = retour catalogue | esc = menu"))
	b.WriteString("\n")
	return b.String()
}
