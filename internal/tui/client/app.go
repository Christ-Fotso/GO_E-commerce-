package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"ecommerce-cli/internal/models"
)

const defaultListWidth = 48

// Session holds the logged-in user (mock for now, token later).
type Session struct {
	User  *models.User
	Token string
}

// Model is the root TUI model: navigation + global state.
type Model struct {
	screen Screen
	width  int
	height int

	session  Session
	cart     models.Cart
	orders   []models.Order
	products []models.Product
	accounts []mockAccount

	statusOK  string
	statusErr string

	selectedProduct *models.Product
	cartCursor      int
	ordersCursor    int

	menu         list.Model
	catalogList  list.Model
	registerForm *huh.Form
	confirmForm  *huh.Form
	loginForm    *huh.Form
	resetForm    *huh.Form
	paymentForm  *huh.Form

	api *APIClient

	quitting bool
	nextID   int
}

type menuItem struct {
	title  string
	screen Screen
	quit   bool
	logout bool
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

func newMenuList(items []list.Item) list.Model {
	l := list.New(items, menuDelegate{}, defaultListWidth, len(items)+2)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)
	return l
}

func guestMenuItems() []list.Item {
	return []list.Item{
		menuItem{title: "Register", screen: ScreenRegister},
		menuItem{title: "Confirm", screen: ScreenConfirm},
		menuItem{title: "Login", screen: ScreenLogin},
		menuItem{title: "Reset password", screen: ScreenReset},
		menuItem{title: "Quitter", quit: true},
	}
}

func shopMenuItems() []list.Item {
	return []list.Item{
		menuItem{title: "Catalogue / Recherche", screen: ScreenCatalog},
		menuItem{title: "Panier", screen: ScreenCart},
		menuItem{title: "Mes commandes", screen: ScreenOrders},
		menuItem{title: "Deconnexion", logout: true},
		menuItem{title: "Quitter", quit: true},
	}
}

// NewModel returns the initial client TUI model.
func NewModel() Model {
	products := mockProducts()
	m := Model{
		screen:   ScreenMenu,
		products: products,
		orders:   nil,
		accounts: seedAccounts(),
		cart: models.Cart{
			ID:     1,
			UserID: 0,
			Items:  nil,
		},
		menu:   newMenuList(guestMenuItems()),
		api:    NewAPIClient(""),
		nextID: 2,
	}
	m.catalogList = newCatalogList(products)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) clearStatus() {
	m.statusOK = ""
	m.statusErr = ""
}

func (m *Model) refreshMenu() {
	if m.session.User != nil {
		m.menu = newMenuList(shopMenuItems())
	} else {
		m.menu = newMenuList(guestMenuItems())
	}
	if m.width > 0 {
		m.menu.SetWidth(m.width)
	}
}

func (m *Model) goTo(s Screen) tea.Cmd {
	m.clearStatus()
	m.screen = s
	switch s {
	case ScreenRegister:
		m.registerForm = newRegisterForm()
		return m.registerForm.Init()
	case ScreenConfirm:
		m.confirmForm = newConfirmForm()
		return m.confirmForm.Init()
	case ScreenLogin:
		m.loginForm = newLoginForm()
		return m.loginForm.Init()
	case ScreenReset:
		m.resetForm = newResetForm()
		return m.resetForm.Init()
	case ScreenCatalog:
		m.catalogList = newCatalogList(m.products)
		if m.width > 0 {
			m.catalogList.SetWidth(m.width)
			m.catalogList.SetHeight(12)
		}
		return nil
	case ScreenPayment:
		m.paymentForm = newPaymentForm()
		return m.paymentForm.Init()
	default:
		return nil
	}
}

func (m *Model) goMenu() tea.Cmd {
	m.screen = ScreenMenu
	m.clearStatus()
	return nil
}

func (m Model) isFormScreen() bool {
	switch m.screen {
	case ScreenRegister, ScreenConfirm, ScreenLogin, ScreenReset, ScreenPayment:
		return true
	default:
		return false
	}
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
		m.catalogList.SetWidth(msg.Width)
		m.catalogList.SetHeight(12)
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if m.screen == ScreenMenu && key == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		// Sur les formulaires huh, laisser esc au formulaire (abort).
		if !m.isFormScreen() && m.screen != ScreenMenu && (key == "esc" || key == "q") {
			if m.screen == ScreenProduct {
				return m, m.goTo(ScreenCatalog)
			}
			cmd := m.goMenu()
			return m, cmd
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
			if item.logout {
				m.session = Session{}
				m.cart.Items = nil
				m.cartCursor = 0
				m.orders = nil
				m.refreshMenu()
				m.statusOK = "Deconnecte"
				return m, nil
			}
			return m, m.goTo(item.screen)
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
	b.WriteString(titleStyle.Render("E-Commerce CLI - Client"))
	b.WriteString("\n")

	if m.session.User != nil {
		b.WriteString(infoStyle.Render(fmt.Sprintf("Connecte: %s (%s)", m.session.User.Username, m.session.User.Email)))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Panier: %d ligne(s) | Total TTC: %.2f EUR\n\n", len(m.cart.Items), cartTotalTTC(m.cart)))
	} else {
		b.WriteString(subtitleStyle.Render("Non connecte - inscris-toi ou connecte-toi"))
		b.WriteString("\n\n")
	}

	b.WriteString(m.menu.View())
	b.WriteString(helpStyle.Render("\nup/down | enter | q quitter"))
	b.WriteString("\n")
	return b.String()
}

func screenHeader(title string) string {
	return titleStyle.Render(title) + "\n" +
		helpStyle.Render("esc retour (formulaires: abort huh)") + "\n\n"
}

func applyFormUpdate(form *huh.Form, msg tea.Msg) (*huh.Form, tea.Cmd) {
	if form == nil {
		return nil, nil
	}
	model, cmd := form.Update(msg)
	if f, ok := model.(*huh.Form); ok {
		return f, cmd
	}
	return form, cmd
}
