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
	price := productUnitTTC(&p.product)
	return fmt.Sprintf("%.2f EUR TTC | stock %d | %s", price, p.product.Stock, p.product.Description)
}
func (p productItem) FilterValue() string {
	return p.product.Name + " " + p.product.Description + " " + p.product.Category
}

type productDelegate struct{}

func (d productDelegate) Height() int                             { return 2 }
func (d productDelegate) Spacing() int                            { return 0 }
func (d productDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d productDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(productItem)
	if !ok {
		return
	}
	title := item.Title()
	desc := item.Description()
	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("> "+title)+"\n"+itemStyle.Render("    "+desc))
		return
	}
	fmt.Fprint(w, itemStyle.Render("  "+title)+"\n"+helpStyle.Render("    "+desc))
}

func newCatalogList(products []models.Product) list.Model {
	items := make([]list.Item, 0, len(products))
	for _, p := range products {
		items = append(items, productItem{product: p})
	}
	l := list.New(items, productDelegate{}, 70, 20)
	l.Title = "Catalogue des Produits"
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
					m.loading = true
					return m, m.cmdAddToCart(item.product.ID, 1)
				}
			case "r":
				m.loading = true
				return m, m.cmdLoadProducts("")
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
					m.loading = true
					return m, m.cmdAddToCart(m.selectedProduct.ID, 1)
				}
			}
		}
	}
	return m, nil
}

func (m Model) viewCatalog() string {
	if m.screen == ScreenProduct {
		return m.viewProduct()
	}

	var b strings.Builder
	b.WriteString(screenHeader("Catalogue / Recherche"))
	b.WriteString(subtitleStyle.Render("/ filtrer local | r recharger API | enter fiche | a panier"))
	b.WriteString("\n\n")
	b.WriteString(m.catalogList.View())
	b.WriteString("\n")
	return b.String()
}

func (m Model) viewProduct() string {
	var b strings.Builder
	b.WriteString(screenHeader("Fiche produit"))

	if m.selectedProduct == nil {
		b.WriteString(errorStyle.Render("Aucun produit sélectionné") + "\n")
		return b.String()
	}

	p := m.selectedProduct
	b.WriteString(sectionStyle.Render(p.Name) + "\n")
	b.WriteString(p.Description + "\n\n")
	b.WriteString(fmt.Sprintf("ID: %d | Ref: %s | Cat: %s\n", p.ID, p.Reference, p.Category))
	b.WriteString(fmt.Sprintf("Prix HT: %.2f | TTC: %.2f EUR (TVA %.0f%%)\n", p.Price, productUnitTTC(p), p.TaxRate))
	b.WriteString(fmt.Sprintf("Stock: %d\n", p.Stock))
	b.WriteString(fmt.Sprintf("Lignes panier: %d | Total TTC: %.2f EUR\n", len(m.cart.Items), cartTotalTTC(m.cart)))
	b.WriteString(helpStyle.Render("\na/enter ajouter | b catalogue | esc menu"))
	b.WriteString("\n")
	return b.String()
}
