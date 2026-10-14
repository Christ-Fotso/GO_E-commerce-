package admin

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenUsers
	ScreenUserCreate
	ScreenUserEdit
	ScreenProducts
	ScreenProductCreate
	ScreenOrders
	ScreenOrderDetails
	ScreenOrderStatus
	ScreenOrderCreate
)

func (s Screen) String() string {
	switch s {
	case ScreenMenu:
		return "Menu Principal"
	case ScreenUsers:
		return "Gestion des Utilisateurs"
	case ScreenUserCreate:
		return "Créer un Utilisateur"
	case ScreenUserEdit:
		return "Modifier un Utilisateur"
	case ScreenProducts:
		return "Gestion des Produits"
	case ScreenProductCreate:
		return "Ajouter un Produit"
	case ScreenOrders:
		return "Gestion des Commandes"
	case ScreenOrderDetails:
		return "Détails de la Commande"
	case ScreenOrderStatus:
		return "Changer le Statut"
	case ScreenOrderCreate:
		return "Créer une Commande"
	default:
		return "Administration"
	}
}
