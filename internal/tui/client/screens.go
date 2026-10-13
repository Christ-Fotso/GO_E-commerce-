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
		return "Inscription"
	case ScreenConfirm:
		return "Confirmation"
	case ScreenLogin:
		return "Connexion"
	case ScreenReset:
		return "Réinitialisation"
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
