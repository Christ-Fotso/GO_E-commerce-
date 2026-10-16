package admin

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"ecommerce-cli/internal/models"
)

type Model struct {
	screen   Screen
	width    int
	height   int
	api      *APIClient
	adminUser *models.User

	users        []models.User
	userCursor   int
	products     []models.Product
	prodCursor   int
	orders       []models.Order
	orderCursor  int
	currentOrder *models.Order

	// Formulaires
	activeForm   *huh.Form
	userCreate   userCreateData
	userEdit     userEditData
	prodCreate   prodCreateData
	orderStatus  orderStatusData
	orderCreate  orderCreateData

	statusMsg    string
	statusIsErr  bool
	loading      bool
}

type userCreateData struct {
	username    string
	email       string
	password    string
	role        string
	isConfirmed bool
}

type userEditData struct {
	id          int
	username    string
	email       string
	password    string
	role        string
	isConfirmed bool
}

type prodCreateData struct {
	name        string
	description string
	category    string
	price       string
	taxRate     string
	stock       string
}

type orderStatusData struct {
	orderID int
	status  string
	reason  string
}

type orderCreateData struct {
	userID    string
	productID string
	quantity  string
	status    string
}

func NewModel() *Model {
	return &Model{
		screen: ScreenMenu,
		api:    NewAPIClient(""),
	}
}

// NewAdminModel est un alias conservé pour la compatibilité
func NewAdminModel() *Model {
	return NewModel()
}

// Messages Bubble Tea internes
type usersLoadedMsg struct{ users []models.User; err error }
type prodsLoadedMsg struct{ products []models.Product; err error }
type ordersLoadedMsg struct{ orders []models.Order; err error }
type actionDoneMsg struct{ message string }
type actionErrMsg struct{ err error }

func (m *Model) Init() tea.Cmd {
	return m.cmdEnsureAuth()
}

func (m *Model) cmdEnsureAuth() tea.Cmd {
	return func() tea.Msg {
		user, err := m.api.Login("admin@shop.local", "admin123")
		if err != nil {
			return actionErrMsg{err: fmt.Errorf("impossible de se connecter à l'API admin : %w", err)}
		}
		return actionDoneMsg{message: fmt.Sprintf("Connecté en tant que %s (%s)", user.Username, user.Email)}
	}
}

func (m *Model) cmdLoadUsers() tea.Cmd {
	return func() tea.Msg {
		users, err := m.api.ListUsers()
		return usersLoadedMsg{users: users, err: err}
	}
}

func (m *Model) cmdLoadProducts() tea.Cmd {
	return func() tea.Msg {
		prods, err := m.api.ListProducts()
		return prodsLoadedMsg{products: prods, err: err}
	}
}

func (m *Model) cmdLoadOrders() tea.Cmd {
	return func() tea.Msg {
		orders, err := m.api.ListOrders(0, "")
		return ordersLoadedMsg{orders: orders, err: err}
	}
}

func (m *Model) cmdConfirmUser(id int) tea.Cmd {
	return func() tea.Msg {
		user, err := m.api.ConfirmUser(id)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Utilisateur %s (%d) confirmé avec succès", user.Username, user.ID)}
	}
}

func (m *Model) cmdDeleteUser(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.api.DeleteUser(id)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Utilisateur #%d supprimé", id)}
	}
}

func (m *Model) cmdSubmitUserCreate() tea.Cmd {
	return func() tea.Msg {
		user, err := m.api.CreateUser(
			m.userCreate.username,
			m.userCreate.email,
			m.userCreate.password,
			m.userCreate.role,
			m.userCreate.isConfirmed,
		)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Utilisateur %s créé avec succès (#%d)", user.Username, user.ID)}
	}
}

func (m *Model) cmdSubmitUserEdit() tea.Cmd {
	return func() tea.Msg {
		fields := map[string]any{
			"username":     m.userEdit.username,
			"email":        m.userEdit.email,
			"role":         m.userEdit.role,
			"is_confirmed": m.userEdit.isConfirmed,
		}
		if strings.TrimSpace(m.userEdit.password) != "" {
			fields["password"] = m.userEdit.password
		}
		user, err := m.api.UpdateUser(m.userEdit.id, fields)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Utilisateur #%d (%s) modifié avec succès", user.ID, user.Username)}
	}
}

func (m *Model) cmdSubmitProductCreate() tea.Cmd {
	return func() tea.Msg {
		price, err := strconv.ParseFloat(strings.TrimSpace(m.prodCreate.price), 64)
		if err != nil {
			return actionErrMsg{err: fmt.Errorf("prix invalide")}
		}
		taxRate, err := strconv.ParseFloat(strings.TrimSpace(m.prodCreate.taxRate), 64)
		if err != nil {
			taxRate = 20.0
		}
		stock, err := strconv.Atoi(strings.TrimSpace(m.prodCreate.stock))
		if err != nil {
			return actionErrMsg{err: fmt.Errorf("stock invalide")}
		}
		prod, err := m.api.CreateProduct(
			m.prodCreate.name,
			m.prodCreate.description,
			m.prodCreate.category,
			price,
			taxRate,
			stock,
		)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Produit %s ajouté avec la référence %s", prod.Name, prod.Reference)}
	}
}

func (m *Model) cmdSubmitOrderStatus() tea.Cmd {
	return func() tea.Msg {
		order, err := m.api.UpdateOrderStatus(m.orderStatus.orderID, m.orderStatus.status, m.orderStatus.reason)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Commande %s passée au statut %s", order.Reference, order.Status)}
	}
}

func (m *Model) cmdSubmitOrderCreate() tea.Cmd {
	return func() tea.Msg {
		uid, err := strconv.Atoi(strings.TrimSpace(m.orderCreate.userID))
		if err != nil {
			return actionErrMsg{err: fmt.Errorf("id client invalide")}
		}
		pid, err := strconv.Atoi(strings.TrimSpace(m.orderCreate.productID))
		if err != nil {
			return actionErrMsg{err: fmt.Errorf("id produit invalide")}
		}
		qty, err := strconv.Atoi(strings.TrimSpace(m.orderCreate.quantity))
		if err != nil || qty <= 0 {
			return actionErrMsg{err: fmt.Errorf("quantité invalide (minimum 1)")}
		}
		items := []OrderItemReq{{ProductID: pid, Quantity: qty}}
		order, err := m.api.CreateOrder(uid, items, m.orderCreate.status)
		if err != nil {
			return actionErrMsg{err: err}
		}
		return actionDoneMsg{message: fmt.Sprintf("Commande %s créée pour le client #%d", order.Reference, uid)}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case actionDoneMsg:
		m.loading = false
		m.statusMsg = msg.message
		m.statusIsErr = false
		if m.screen == ScreenUserCreate || m.screen == ScreenUserEdit {
			m.screen = ScreenUsers
			return m, m.cmdLoadUsers()
		}
		if m.screen == ScreenProductCreate {
			m.screen = ScreenProducts
			return m, m.cmdLoadProducts()
		}
		if m.screen == ScreenOrderStatus || m.screen == ScreenOrderCreate {
			m.screen = ScreenOrders
			return m, m.cmdLoadOrders()
		}
		return m, nil

	case actionErrMsg:
		m.loading = false
		m.statusMsg = "Erreur : " + msg.err.Error()
		m.statusIsErr = true
		return m, nil

	case usersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.statusMsg = "Erreur de chargement : " + msg.err.Error()
			m.statusIsErr = true
		} else {
			m.users = msg.users
			if m.userCursor >= len(m.users) {
				m.userCursor = max(0, len(m.users)-1)
			}
		}
		return m, nil

	case prodsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.statusMsg = "Erreur de chargement : " + msg.err.Error()
			m.statusIsErr = true
		} else {
			m.products = msg.products
			if m.prodCursor >= len(m.products) {
				m.prodCursor = max(0, len(m.products)-1)
			}
		}
		return m, nil

	case ordersLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.statusMsg = "Erreur de chargement : " + msg.err.Error()
			m.statusIsErr = true
		} else {
			m.orders = msg.orders
			if m.orderCursor >= len(m.orders) {
				m.orderCursor = max(0, len(m.orders)-1)
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Gestion globale pour quitter
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Si un formulaire est actif
		if m.activeForm != nil {
			form, cmd := m.activeForm.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.activeForm = f
			}
			if m.activeForm.State == huh.StateCompleted {
				var submitCmd tea.Cmd
				switch m.screen {
				case ScreenUserCreate:
					submitCmd = m.cmdSubmitUserCreate()
				case ScreenUserEdit:
					submitCmd = m.cmdSubmitUserEdit()
				case ScreenProductCreate:
					submitCmd = m.cmdSubmitProductCreate()
				case ScreenOrderStatus:
					submitCmd = m.cmdSubmitOrderStatus()
				case ScreenOrderCreate:
					submitCmd = m.cmdSubmitOrderCreate()
				}
				m.activeForm = nil
				m.loading = true
				return m, submitCmd
			}
			if m.activeForm.State == huh.StateAborted || msg.String() == "esc" {
				m.activeForm = nil
				m.statusMsg = "Opération annulée"
				m.statusIsErr = false
				switch m.screen {
				case ScreenUserCreate, ScreenUserEdit:
					m.screen = ScreenUsers
				case ScreenProductCreate:
					m.screen = ScreenProducts
				case ScreenOrderStatus, ScreenOrderCreate:
					m.screen = ScreenOrders
				default:
					m.screen = ScreenMenu
				}
				return m, nil
			}
			return m, cmd
		}

		// Navigation des écrans
		switch m.screen {
		case ScreenMenu:
			return m.updateMenu(msg)
		case ScreenUsers:
			return m.updateUsers(msg)
		case ScreenProducts:
			return m.updateProducts(msg)
		case ScreenOrders:
			return m.updateOrders(msg)
		case ScreenOrderDetails:
			if msg.String() == "esc" || msg.String() == "b" || msg.String() == "enter" {
				m.screen = ScreenOrders
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1", "u":
		m.screen = ScreenUsers
		m.statusMsg = ""
		m.loading = true
		return m, m.cmdLoadUsers()
	case "2", "p":
		m.screen = ScreenProducts
		m.statusMsg = ""
		m.loading = true
		return m, m.cmdLoadProducts()
	case "3", "o", "c":
		m.screen = ScreenOrders
		m.statusMsg = ""
		m.loading = true
		return m, m.cmdLoadOrders()
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) updateUsers(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.userCursor > 0 {
			m.userCursor--
		}
	case "down", "j":
		if m.userCursor < len(m.users)-1 {
			m.userCursor++
		}
	case "c":
		m.userCreate = userCreateData{role: "customer", isConfirmed: true}
		m.activeForm = newUserCreateForm(
			&m.userCreate.username,
			&m.userCreate.email,
			&m.userCreate.password,
			&m.userCreate.role,
			&m.userCreate.isConfirmed,
		)
		m.screen = ScreenUserCreate
		return m, m.activeForm.Init()
	case "e":
		if len(m.users) == 0 {
			return m, nil
		}
		u := m.users[m.userCursor]
		m.userEdit = userEditData{
			id:          u.ID,
			username:    u.Username,
			email:       u.Email,
			role:        u.Role,
			isConfirmed: u.IsConfirmed,
		}
		m.activeForm = newUserEditForm(
			&m.userEdit.username,
			&m.userEdit.email,
			&m.userEdit.password,
			&m.userEdit.role,
			&m.userEdit.isConfirmed,
		)
		m.screen = ScreenUserEdit
		return m, m.activeForm.Init()
	case "v":
		if len(m.users) == 0 {
			return m, nil
		}
		u := m.users[m.userCursor]
		m.loading = true
		return m, m.cmdConfirmUser(u.ID)
	case "d", "x":
		if len(m.users) == 0 {
			return m, nil
		}
		u := m.users[m.userCursor]
		m.loading = true
		return m, m.cmdDeleteUser(u.ID)
	case "r":
		m.loading = true
		return m, m.cmdLoadUsers()
	case "esc", "b":
		m.screen = ScreenMenu
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) updateProducts(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.prodCursor > 0 {
			m.prodCursor--
		}
	case "down", "j":
		if m.prodCursor < len(m.products)-1 {
			m.prodCursor++
		}
	case "c", "a":
		m.prodCreate = prodCreateData{
			taxRate: "20.0",
			stock:   "10",
		}
		m.activeForm = newProductCreateForm(
			&m.prodCreate.name,
			&m.prodCreate.description,
			&m.prodCreate.category,
			&m.prodCreate.price,
			&m.prodCreate.taxRate,
			&m.prodCreate.stock,
		)
		m.screen = ScreenProductCreate
		return m, m.activeForm.Init()
	case "r":
		m.loading = true
		return m, m.cmdLoadProducts()
	case "esc", "b":
		m.screen = ScreenMenu
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) updateOrders(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.orderCursor > 0 {
			m.orderCursor--
		}
	case "down", "j":
		if m.orderCursor < len(m.orders)-1 {
			m.orderCursor++
		}
	case "enter", "d":
		if len(m.orders) == 0 {
			return m, nil
		}
		o := m.orders[m.orderCursor]
		m.currentOrder = &o
		m.screen = ScreenOrderDetails
		return m, nil
	case "s":
		if len(m.orders) == 0 {
			return m, nil
		}
		o := m.orders[m.orderCursor]
		m.orderStatus = orderStatusData{
			orderID: o.ID,
			status:  "paid",
		}
		m.activeForm = newOrderStatusForm(
			&m.orderStatus.status,
			&m.orderStatus.reason,
		)
		m.screen = ScreenOrderStatus
		return m, m.activeForm.Init()
	case "c":
		m.orderCreate = orderCreateData{
			quantity: "1",
			status:   "pending",
		}
		m.activeForm = newOrderCreateForm(
			&m.orderCreate.userID,
			&m.orderCreate.productID,
			&m.orderCreate.quantity,
			&m.orderCreate.status,
		)
		m.screen = ScreenOrderCreate
		return m, m.activeForm.Init()
	case "r":
		m.loading = true
		return m, m.cmdLoadOrders()
	case "esc", "b":
		m.screen = ScreenMenu
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("═══ PANEL D'ADMINISTRATION E-COMMERCE ═══"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Navigation : " + m.screen.String()))
	b.WriteString("\n\n")

	if m.statusMsg != "" {
		if m.statusIsErr {
			b.WriteString(errorStyle.Render("[ERREUR] "+m.statusMsg) + "\n\n")
		} else {
			b.WriteString(successStyle.Render("[SUCCES] "+m.statusMsg) + "\n\n")
		}
	}

	if m.loading {
		b.WriteString(subtitleStyle.Render("Chargement en cours...") + "\n\n")
		return b.String()
	}

	if m.activeForm != nil {
		b.WriteString(m.activeForm.View())
		b.WriteString("\n" + helpStyle.Render("tab/shift+tab: naviguer | enter: valider | esc: annuler"))
		return b.String()
	}

	switch m.screen {
	case ScreenMenu:
		b.WriteString(itemStyle.Render("[1] Gérer les utilisateurs (créer, modifier, supprimer, confirmer)"))
		b.WriteString("\n")
		b.WriteString(itemStyle.Render("[2] Gérer les produits (consulter et ajouter des produits PDT-...)"))
		b.WriteString("\n")
		b.WriteString(itemStyle.Render("[3] Gérer les commandes (changer les statuts, créer une commande client)"))
		b.WriteString("\n\n")
		b.WriteString(itemStyle.Render("[q] Quitter"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Appuyez sur 1, 2, 3 ou q"))

	case ScreenUsers:
		b.WriteString(m.renderUsers())

	case ScreenProducts:
		b.WriteString(m.renderProducts())

	case ScreenOrders:
		b.WriteString(m.renderOrders())

	case ScreenOrderDetails:
		b.WriteString(m.renderOrderDetails())
	}

	return b.String()
}

func (m *Model) renderUsers() string {
	var b strings.Builder
	if len(m.users) == 0 {
		b.WriteString(subtitleStyle.Render("Aucun utilisateur trouvé.") + "\n")
	} else {
		for i, u := range m.users {
			roleBadge := badgeCustomerStyle.Render("[CLIENT]")
			if u.Role == models.RoleAdmin {
				roleBadge = badgeAdminStyle.Render("[ADMIN]")
			}
			confirmBadge := successStyle.Render("[CONFIRMÉ]")
			if !u.IsConfirmed {
				confirmBadge = errorStyle.Render("[EN ATTENTE]")
			}

			line := fmt.Sprintf("#%d %s (%s) %s %s", u.ID, u.Username, u.Email, roleBadge, confirmBadge)
			if i == m.userCursor {
				b.WriteString(selectedItemStyle.Render("> " + line) + "\n")
			} else {
				b.WriteString(itemStyle.Render("  " + line) + "\n")
			}
		}
	}
	b.WriteString("\n" + helpStyle.Render("c: Créer | e: Modifier | v: Confirmer | d: Supprimer | r: Actualiser | esc: Menu"))
	return b.String()
}

func (m *Model) renderProducts() string {
	var b strings.Builder
	if len(m.products) == 0 {
		b.WriteString(subtitleStyle.Render("Aucun produit trouvé.") + "\n")
	} else {
		for i, p := range m.products {
			priceTTC := p.PriceTTC
			if priceTTC <= 0 {
				priceTTC = p.Price * (1 + p.TaxRate/100)
			}
			line := fmt.Sprintf("%-10s %-24s %7.2f€ HT (%5.1f%% TVA) -> %7.2f€ TTC | Stock: %3d | Cat: %s",
				p.Reference, p.Name, p.Price, p.TaxRate, priceTTC, p.Stock, p.Category)

			if i == m.prodCursor {
				b.WriteString(selectedItemStyle.Render("> " + line) + "\n")
			} else {
				b.WriteString(itemStyle.Render("  " + line) + "\n")
			}
		}
	}
	b.WriteString("\n" + helpStyle.Render("c: Ajouter un produit | r: Actualiser | esc: Menu"))
	return b.String()
}

func (m *Model) renderOrders() string {
	var b strings.Builder
	if len(m.orders) == 0 {
		b.WriteString(subtitleStyle.Render("Aucune commande trouvée.") + "\n")
	} else {
		for i, o := range m.orders {
			var badge string
			switch o.Status {
			case models.StatusPaid:
				badge = badgePaidStyle.Render("[PAID]")
			case models.StatusShipping:
				badge = badgeShippingStyle.Render("[SHIPPING]")
			case models.StatusDelivered:
				badge = badgeDeliveredStyle.Render("[DELIVERED]")
			case models.StatusCancelled:
				badge = badgeCancelledStyle.Render("[CANCELLED]")
			default:
				badge = badgePendingStyle.Render("[PENDING]")
			}

			line := fmt.Sprintf("%-10s Client #%-3d %8.2f€ TTC %s %s",
				o.Reference, o.UserID, o.Total, badge, o.CreatedAt.Format("2006-01-02 15:04"))

			if i == m.orderCursor {
				b.WriteString(selectedItemStyle.Render("> " + line) + "\n")
			} else {
				b.WriteString(itemStyle.Render("  " + line) + "\n")
			}
		}
	}
	b.WriteString("\n" + helpStyle.Render("enter: Détails | s: Changer Statut | c: Créer commande client | r: Actualiser | esc: Menu"))
	return b.String()
}

func (m *Model) renderOrderDetails() string {
	var b strings.Builder
	if m.currentOrder == nil {
		b.WriteString("Aucune commande sélectionnée.\n")
		return b.String()
	}
	o := m.currentOrder
	b.WriteString(fmt.Sprintf("Commande Référence : %s (ID #%d)\n", o.Reference, o.ID))
	b.WriteString(fmt.Sprintf("Client ID          : %d\n", o.UserID))
	b.WriteString(fmt.Sprintf("Statut             : %s\n", o.Status))
	if o.CancelReason != "" {
		b.WriteString(fmt.Sprintf("Motif d'annulation : %s\n", o.CancelReason))
	}
	b.WriteString(fmt.Sprintf("Montant Total TTC  : %.2f €\n", o.Total))
	b.WriteString(fmt.Sprintf("Date               : %s\n\n", o.CreatedAt.Format("2006-01-02 15:04:05")))

	b.WriteString("Articles :\n")
	if len(o.Items) == 0 {
		b.WriteString("  (Détail des lignes de commande disponible via l'API)\n")
	} else {
		for _, item := range o.Items {
			name := item.ProductName
			if name == "" {
				name = fmt.Sprintf("Produit #%d", item.ProductID)
			}
			b.WriteString(fmt.Sprintf("  - %s x%d : %.2f € unitaire TTC\n", name, item.Quantity, item.PriceAtTime))
		}
	}

	b.WriteString("\n" + helpStyle.Render("Appuyez sur 'esc' ou 'enter' pour revenir à la liste"))
	return b.String()
}
