package client

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ecommerce-cli/internal/models"
)

func (m Model) cmdRegister(username, email, password string) tea.Cmd {
	api := m.api
	accounts := m.accounts
	nextID := m.nextID
	return func() tea.Msg {
		res, err := api.Register(username, email, password)
		if err == nil {
			return registerDoneMsg{
				username: username,
				email:    email,
				password: password,
				code:     res.ConfirmationCode,
				source:   SourceAPI,
			}
		}
		// Fallback mock
		for _, acc := range accounts {
			if strings.EqualFold(acc.User.Email, email) || strings.EqualFold(acc.User.Username, username) {
				return registerDoneMsg{err: fmt.Errorf("email ou username déjà pris (mock); API: %v", err)}
			}
		}
		_ = nextID
		return registerDoneMsg{
			username: username,
			email:    email,
			password: password,
			code:     mockConfirmCode,
			source:   SourceMock,
		}
	}
}

func (m Model) cmdConfirm(email, code string) tea.Cmd {
	api := m.api
	accounts := append([]mockAccount(nil), m.accounts...)
	return func() tea.Msg {
		if err := api.Confirm(email, code); err == nil {
			return confirmDoneMsg{email: email, source: SourceAPI}
		} else {
			apiErr := err
			idx := -1
			for i, acc := range accounts {
				if strings.EqualFold(acc.User.Email, email) {
					idx = i
					break
				}
			}
			if idx < 0 {
				return confirmDoneMsg{err: fmt.Errorf("utilisateur introuvable (mock); API: %v", apiErr)}
			}
			if code != mockConfirmCode && code != accounts[idx].User.ConfirmationCode {
				return confirmDoneMsg{err: fmt.Errorf("code invalide (mock); API: %v", apiErr)}
			}
			return confirmDoneMsg{email: email, source: SourceMock}
		}
	}
}

func (m Model) cmdLogin(email, password string) tea.Cmd {
	api := m.api
	accounts := append([]mockAccount(nil), m.accounts...)
	return func() tea.Msg {
		res, err := api.Login(email, password)
		if err == nil {
			return loginDoneMsg{user: res.User, token: res.Token, source: SourceAPI}
		}
		apiErr := err
		var found *mockAccount
		for i := range accounts {
			if strings.EqualFold(accounts[i].User.Email, email) {
				found = &accounts[i]
				break
			}
		}
		if found == nil {
			return loginDoneMsg{err: fmt.Errorf("utilisateur introuvable (mock); API: %v", apiErr)}
		}
		if !found.User.IsConfirmed {
			return loginDoneMsg{err: fmt.Errorf("compte non confirmé (mock); API: %v", apiErr)}
		}
		if found.Password != password {
			return loginDoneMsg{err: fmt.Errorf("mot de passe incorrect (mock); API: %v", apiErr)}
		}
		return loginDoneMsg{user: found.User, token: "mock-token", source: SourceMock}
	}
}

func (m Model) cmdReset(email, code, newPassword string) tea.Cmd {
	api := m.api
	accounts := m.accounts
	return func() tea.Msg {
		// Pas de code => demander un code (forgot)
		if strings.TrimSpace(code) == "" {
			res, err := api.ForgotPassword(email)
			if err == nil {
				return resetDoneMsg{
					message: fmt.Sprintf("%s (code=%s)", res.Message, res.ResetCode),
					source:  SourceAPI,
				}
			}
			return resetDoneMsg{
				message: fmt.Sprintf("Email de reset envoyé à %s (code mock=%s)", email, mockResetCode),
				source:  SourceMock,
			}
		}
		if newPassword == "" {
			return resetDoneMsg{err: fmt.Errorf("nouveau mot de passe requis pour reset")}
		}
		if err := api.ResetPassword(email, code, newPassword); err == nil {
			return resetDoneMsg{message: "Mot de passe modifié, tu peux te connecter", source: SourceAPI}
		}
		for i := range accounts {
			if strings.EqualFold(accounts[i].User.Email, email) {
				accounts[i].Password = newPassword
				break
			}
		}
		return resetDoneMsg{
			message: fmt.Sprintf("Mot de passe modifié (mock) pour %s", email),
			source:  SourceMock,
		}
	}
}

func (m Model) cmdLogout() tea.Cmd {
	api := m.api
	token := m.session.Token
	return func() tea.Msg {
		api.SetToken(token)
		err := api.Logout()
		if err != nil {
			api.ClearToken()
			return logoutDoneMsg{source: SourceMock, err: err}
		}
		return logoutDoneMsg{source: SourceAPI}
	}
}

func (m Model) cmdLoadProducts(query string) tea.Cmd {
	api := m.api
	return func() tea.Msg {
		products, err := api.ListProducts(query)
		if err == nil {
			return productsDoneMsg{products: products, source: SourceAPI}
		}
		all := mockProducts()
		if strings.TrimSpace(query) == "" {
			return productsDoneMsg{products: all, source: SourceMock}
		}
		q := strings.ToLower(query)
		var filtered []models.Product
		for _, p := range all {
			if strings.Contains(strings.ToLower(p.Name), q) || strings.Contains(strings.ToLower(p.Description), q) {
				filtered = append(filtered, p)
			}
		}
		return productsDoneMsg{products: filtered, source: SourceMock}
	}
}

func (m Model) cmdLoadCart(okMsg string) tea.Cmd {
	api := m.api
	token := m.session.Token
	local := m.cart
	return func() tea.Msg {
		api.SetToken(token)
		cart, err := api.GetCart()
		if err == nil {
			return cartDoneMsg{cart: *cart, source: SourceAPI, okMsg: okMsg}
		}
		return cartDoneMsg{cart: local, source: SourceMock, okMsg: okMsg}
	}
}

func (m Model) cmdAddToCart(productID, qty int) tea.Cmd {
	api := m.api
	token := m.session.Token
	local := m.cart
	products := m.products
	return func() tea.Msg {
		api.SetToken(token)
		cart, err := api.AddCartItem(productID, qty)
		if err == nil {
			return cartDoneMsg{cart: *cart, source: SourceAPI, okMsg: "Ajouté au panier [API]"}
		}
		// mock fallback
		var prod *models.Product
		for i := range products {
			if products[i].ID == productID {
				p := products[i]
				prod = &p
				break
			}
		}
		if prod == nil {
			return cartDoneMsg{err: fmt.Errorf("produit introuvable; API: %v", err)}
		}
		cartMock := local
		found := false
		for i := range cartMock.Items {
			if cartMock.Items[i].ProductID == productID {
				cartMock.Items[i].Quantity += qty
				cartMock.Items[i].LineTotalTTC = productUnitTTC(cartMock.Items[i].Product) * float64(cartMock.Items[i].Quantity)
				found = true
				break
			}
		}
		if !found {
			p := *prod
			cartMock.Items = append(cartMock.Items, models.CartItem{
				ID:           len(cartMock.Items) + 1,
				CartID:       cartMock.ID,
				ProductID:    productID,
				Quantity:     qty,
				LineTotalTTC: productUnitTTC(&p) * float64(qty),
				Product:      &p,
			})
		}
		cartMock.TotalTTC = 0
		cartMock.TotalTTC = cartTotalTTC(cartMock)
		return cartDoneMsg{cart: cartMock, source: SourceMock, okMsg: fmt.Sprintf("Ajouté : %s [mock]", prod.Name)}
	}
}

func (m Model) cmdUpdateCartQty(productID, quantity int) tea.Cmd {
	api := m.api
	token := m.session.Token
	local := m.cart
	return func() tea.Msg {
		api.SetToken(token)
		cart, err := api.UpdateCartItem(productID, quantity)
		if err == nil {
			return cartDoneMsg{cart: *cart, source: SourceAPI, okMsg: "Quantité mise à jour [API]"}
		}
		cartMock := local
		for i := range cartMock.Items {
			if cartMock.Items[i].ProductID == productID {
				if quantity <= 0 {
					cartMock.Items = append(cartMock.Items[:i], cartMock.Items[i+1:]...)
				} else {
					cartMock.Items[i].Quantity = quantity
					cartMock.Items[i].LineTotalTTC = productUnitTTC(cartMock.Items[i].Product) * float64(quantity)
				}
				break
			}
		}
		cartMock.TotalTTC = cartTotalTTC(cartMock)
		return cartDoneMsg{cart: cartMock, source: SourceMock, okMsg: "Quantité mise à jour [mock]"}
	}
}

func (m Model) cmdDeleteCartItem(productID int) tea.Cmd {
	api := m.api
	token := m.session.Token
	local := m.cart
	return func() tea.Msg {
		api.SetToken(token)
		cart, err := api.DeleteCartItem(productID)
		if err == nil {
			return cartDoneMsg{cart: *cart, source: SourceAPI, okMsg: "Article retiré [API]"}
		}
		cartMock := local
		for i := range cartMock.Items {
			if cartMock.Items[i].ProductID == productID {
				cartMock.Items = append(cartMock.Items[:i], cartMock.Items[i+1:]...)
				break
			}
		}
		cartMock.TotalTTC = cartTotalTTC(cartMock)
		return cartDoneMsg{cart: cartMock, source: SourceMock, okMsg: "Article retiré [mock]"}
	}
}

func (m Model) cmdPay(card, expiry, cvc string) tea.Cmd {
	api := m.api
	token := m.session.Token
	localCart := m.cart
	userID := 0
	if m.session.User != nil {
		userID = m.session.User.ID
	}
	nOrders := len(m.orders)
	return func() tea.Msg {
		api.SetToken(token)
		order, err := api.PayCart(card, expiry, cvc)
		if err == nil {
			return payDoneMsg{order: *order, source: SourceAPI}
		}
		if len(localCart.Items) == 0 {
			return payDoneMsg{err: fmt.Errorf("panier vide; API: %v", err)}
		}
		mockOrder := models.Order{
			ID:        2000 + nOrders + 1,
			Reference: fmt.Sprintf("CMD-MOCK-%d", 2000+nOrders+1),
			UserID:    userID,
			Total:     cartTotalTTC(localCart),
			Status:    models.StatusPaid,
			CreatedAt: time.Now(),
		}
		for _, it := range localCart.Items {
			mockOrder.Items = append(mockOrder.Items, models.OrderItem{
				ProductID:   it.ProductID,
				Quantity:    it.Quantity,
				PriceAtTime: productUnitTTC(it.Product),
			})
		}
		return payDoneMsg{order: mockOrder, source: SourceMock}
	}
}

func (m Model) cmdLoadOrders() tea.Cmd {
	api := m.api
	token := m.session.Token
	local := m.orders
	userID := 0
	if m.session.User != nil {
		userID = m.session.User.ID
	}
	return func() tea.Msg {
		api.SetToken(token)
		orders, err := api.ListOrders()
		if err == nil {
			return ordersDoneMsg{orders: orders, source: SourceAPI}
		}
		if local == nil {
			local = mockOrders(userID)
		}
		return ordersDoneMsg{orders: local, source: SourceMock}
	}
}
