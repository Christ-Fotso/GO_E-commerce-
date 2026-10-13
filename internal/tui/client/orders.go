package client

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateOrders(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up", "k":
		if m.ordersCursor > 0 {
			m.ordersCursor--
		}
	case "down", "j":
		if m.ordersCursor < len(m.orders)-1 {
			m.ordersCursor++
		}
	case "r":
		m.loading = true
		return m, m.cmdLoadOrders()
	}
	return m, nil
}

func (m Model) viewOrders() string {
	var b strings.Builder
	b.WriteString(screenHeader("Mes commandes"))
	b.WriteString(subtitleStyle.Render("GET /orders | pending|paid|shipping|delivered|cancelled"))
	b.WriteString("\n\n")

	if len(m.orders) == 0 {
		b.WriteString(subtitleStyle.Render("Aucune commande") + "\n")
		return b.String()
	}

	for i, o := range m.orders {
		ref := o.Reference
		if ref == "" {
			ref = fmt.Sprintf("#%d", o.ID)
		}
		line := fmt.Sprintf("%s  %.2f EUR  [%s]  %s",
			ref, o.Total, o.Status, o.CreatedAt.Format("2006-01-02 15:04"))
		if i == m.ordersCursor {
			b.WriteString(selectedItemStyle.Render("> "+line) + "\n")
		} else {
			b.WriteString(itemStyle.Render("  "+line) + "\n")
		}
	}

	b.WriteString(helpStyle.Render("\nup/down | r refresh | esc menu"))
	b.WriteString("\n")
	return b.String()
}
