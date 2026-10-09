package client

// Screen identifies the active TUI view.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenRegister
	ScreenConfirm
	ScreenLogin
	ScreenReset
	ScreenCatalog
	ScreenProduct
	ScreenCart
	ScreenPayment
	ScreenOrders
)

func (s Screen) String() string {
	switch s {
	case ScreenMenu:
		return "Menu"
	case ScreenRegister:
		return "Register"
	case ScreenConfirm:
		return "Confirm"
	case ScreenLogin:
		return "Login"
	case ScreenReset:
		return "Reset"
	case ScreenCatalog:
		return "Catalogue"
	case ScreenProduct:
		return "Fiche produit"
	case ScreenCart:
		return "Panier"
	case ScreenPayment:
		return "Paiement"
	case ScreenOrders:
		return "Commandes"
	default:
		return "Inconnu"
	}
}
