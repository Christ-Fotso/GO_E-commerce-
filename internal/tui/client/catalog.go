package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"ecommerce-cli/internal/models"
)

type productItem struct {
	product models.Product
}

func (p productItem) Title() string { return p.product.Name }
func (p productItem) Description() string {
	return fmt.Sprintf("%.2f EUR | stock %d | %s", p.product.Price, p.product.Stock, p.product.Description)
}
func (p productItem) FilterValue() string {
	return p.product.Name + " " + p.product.Description
}

type productDelegate struct{}

func (d productDelegate) Height() int                               { return 2 }
func (d productDelegate) Spacing() int                              { return 1 }
func (d productDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd   { return nil }
func (d productDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(productItem)
	if !ok {
		return
	}
	title := item.Title()
	desc := item.Description()
	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("> "+title)+"\n"+itemStyle.Render("  "+desc))
		return
	}
	fmt.Fprint(w, itemStyle.Render("  "+title)+"\n"+itemStyle.Render("  "+desc))
}

func newCatalogList(products []models.Product) list.Model {
	items := make([]list.Item, 0, len(products))
	for _, p := range products {
		items = append(items, productItem{product: p})
	}
	l := list.New(items, productDelegate{}, 60, 14)
	l.Title = "Produits"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.SetShowPagination(true)
	return l
}

func (m Model) updateCatalog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case ScreenCatalog:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(key.String()) {
			case "enter":
				if item, ok := m.catalogList.SelectedItem().(productItem); ok {
					p := item.product
					m.selectedProduct = &p
					m.screen = ScreenProduct
					m.clearStatus()
					return m, nil
				}
			case "a":
				if item, ok := m.catalogList.SelectedItem().(productItem); ok {
					m.addToCart(item.product, 1)
					m.statusOK = fmt.Sprintf("Ajoute: %s (lignes panier: %d)", item.product.Name, len(m.cart.Items))
					return m, nil
				}
			}
		}
		var cmd tea.Cmd
		m.catalogList, cmd = m.catalogList.Update(msg)
		return m, cmd

	case ScreenProduct:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(key.String()) {
			case "b", "backspace":
				return m, m.goTo(ScreenCatalog)
			case "a", "enter":
				if m.selectedProduct != nil {
					m.addToCart(*m.selectedProduct, 1)
					m.statusOK = fmt.Sprintf("Ajoute: %s (lignes panier: %d)", m.selectedProduct.Name, len(m.cart.Items))
				}
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
	b.WriteString(subtitleStyle.Render("/ pour filtrer | enter fiche | a panier"))
	b.WriteString("\n\n")
	b.WriteString(m.catalogList.View())
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
	b.WriteString(fmt.Sprintf("ID: %d\n", p.ID))
	b.WriteString(fmt.Sprintf("Prix: %.2f EUR\n", p.Price))
	b.WriteString(fmt.Sprintf("Stock: %d\n", p.Stock))
	b.WriteString(fmt.Sprintf("Lignes panier: %d | Total TTC: %.2f EUR\n", len(m.cart.Items), cartTotalTTC(m.cart)))
	b.WriteString(helpStyle.Render("\na/enter ajouter | b catalogue | esc menu"))
	b.WriteString("\n")
	return b.String()
}
