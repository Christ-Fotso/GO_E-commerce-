package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"ecommerce-cli/internal/models"
)

const defaultListWidth = 42

// Session holds the "logged in" user (mock for now, token later).
type Session struct {
	User  *models.User
	Token string
}

// Model is the root TUI model: navigation + global state.
type Model struct {
	screen Screen
	width  int
	height int

	session Session
	cart    models.Cart
	orders  []models.Order
	products []models.Product

	statusOK  string
	statusErr string

	// Catalog selection (shared with product detail).
	selectedProduct *models.Product
	catalogCursor   int
	cartCursor      int
	ordersCursor    int

	menu list.Model
	api  *APIClient

	quitting bool
}

type menuItem struct {
	title  string
	screen Screen
	quit   bool
}

func (i menuItem) FilterValue() string { return i.title }

type menuDelegate struct{}

func (d menuDelegate) Height() int                             { return 1 }
func (d menuDelegate) Spacing() int                            { return 0 }
func (d menuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d menuDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("> "+i.title))
		return
	}
	fmt.Fprint(w, itemStyle.Render("  "+i.title))
}

// NewModel returns the initial client TUI model.
func NewModel() Model {
	items := []list.Item{
		menuItem{title: "Register", screen: ScreenRegister},
		menuItem{title: "Confirm", screen: ScreenConfirm},
		menuItem{title: "Login", screen: ScreenLogin},
		menuItem{title: "Reset password", screen: ScreenReset},
		menuItem{title: "Catalogue / Recherche", screen: ScreenCatalog},
		menuItem{title: "Panier", screen: ScreenCart},
		menuItem{title: "Mes commandes", screen: ScreenOrders},
		menuItem{title: "Quitter", quit: true},
	}

	l := list.New(items, menuDelegate{}, defaultListWidth, len(items)+2)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)

	products := mockProducts()
	return Model{
		screen:   ScreenMenu,
		products: products,
		orders:   mockOrders(1),
		cart: models.Cart{
			ID:     1,
			UserID: 0,
			Items:  nil,
		},
		menu: l,
		api:  NewAPIClient(""),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) clearStatus() {
	m.statusOK = ""
	m.statusErr = ""
}

func (m *Model) goTo(s Screen) {
	m.clearStatus()
	m.screen = s
}

func (m *Model) goMenu() {
	m.goTo(ScreenMenu)
}

func cartTotalTTC(cart models.Cart) float64 {
	var total float64
	for _, it := range cart.Items {
		if it.Product != nil {
			total += it.Product.Price * float64(it.Quantity)
		}
	}
	return total
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu.SetWidth(msg.Width)
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		// q quitte depuis le menu ; ailleurs esc / q reviennent au menu.
		if m.screen == ScreenMenu && key == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		if m.screen != ScreenMenu && (key == "esc" || key == "q") {
			m.goMenu()
			return m, nil
		}
	}

	switch m.screen {
	case ScreenMenu:
		return m.updateMenu(msg)
	case ScreenRegister:
		return m.updateRegister(msg)
	case ScreenConfirm:
		return m.updateConfirm(msg)
	case ScreenLogin:
		return m.updateLogin(msg)
	case ScreenReset:
		return m.updateReset(msg)
	case ScreenCatalog, ScreenProduct:
		return m.updateCatalog(msg)
	case ScreenCart:
		return m.updateCart(msg)
	case ScreenPayment:
		return m.updatePayment(msg)
	case ScreenOrders:
		return m.updateOrders(msg)
	default:
		return m, nil
	}
}

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		if item, ok := m.menu.SelectedItem().(menuItem); ok {
			if item.quit {
				m.quitting = true
				return m, tea.Quit
			}
			m.goTo(item.screen)
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return "A bientot !\n"
	}

	var body string
	switch m.screen {
	case ScreenMenu:
		body = m.viewMenu()
	case ScreenRegister:
		body = m.viewRegister()
	case ScreenConfirm:
		body = m.viewConfirm()
	case ScreenLogin:
		body = m.viewLogin()
	case ScreenReset:
		body = m.viewReset()
	case ScreenCatalog, ScreenProduct:
		body = m.viewCatalog()
	case ScreenCart:
		body = m.viewCart()
	case ScreenPayment:
		body = m.viewPayment()
	case ScreenOrders:
		body = m.viewOrders()
	default:
		body = "Ecran inconnu"
	}

	var b strings.Builder
	b.WriteString(body)
	if m.statusOK != "" {
		b.WriteString(successStyle.Render("\n"+m.statusOK) + "\n")
	}
	if m.statusErr != "" {
		b.WriteString(errorStyle.Render("\n"+m.statusErr) + "\n")
	}
	return b.String()
}

func (m Model) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Hello / Menu - E-Commerce CLI"))
	b.WriteString("\n")

	if m.session.User != nil {
		b.WriteString(infoStyle.Render(fmt.Sprintf("Connecte: %s (%s)", m.session.User.Username, m.session.User.Email)))
		b.WriteString("\n")
	} else {
		b.WriteString(subtitleStyle.Render("Non connecte (donnees mock)"))
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("Panier: %d article(s) | Total TTC: %.2f EUR\n\n", len(m.cart.Items), cartTotalTTC(m.cart)))
	b.WriteString(m.menu.View())
	b.WriteString(helpStyle.Render("\nup/down | enter ouvrir | q quitter"))
	b.WriteString("\n")
	return b.String()
}

func screenHeader(title string) string {
	return titleStyle.Render(title) + "\n" +
		helpStyle.Render("esc/q retour menu") + "\n\n"
}
