package models

import "time"

const OrderIDPrefix = "CMD"

// OrderStatus suit le cycle de vie d'une commande (sujet : Commande).
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipping  OrderStatus = "shipping"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

var ValidOrderStatuses = []OrderStatus{
	StatusPending,
	StatusPaid,
	StatusShipping,
	StatusDelivered,
	StatusCancelled,
}

// Order est une commande passée (ou créée par un admin pour un client).
// Identifiant métier obligatoire (ex: CMD-1F2S8B).
// L'annulation doit pouvoir porter une raison.
type Order struct {
	ID           int         `json:"id"`
	BusinessID   string      `json:"business_id"`
	UserID       int         `json:"user_id"`
	Total        float64     `json:"total"`
	Status       OrderStatus `json:"status"`
	CancelReason string      `json:"cancel_reason,omitempty"`
	Items        []OrderItem `json:"items,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// OrderItem fige le produit et son prix TTC au moment de la commande.
type OrderItem struct {
	ID          int      `json:"id"`
	OrderID     int      `json:"order_id"`
	ProductID   int      `json:"product_id"`
	Quantity    int      `json:"quantity"`
	PriceAtTime float64  `json:"price_at_time"`
	Product     *Product `json:"product,omitempty"`
}

func (s OrderStatus) IsValid() bool {
	for _, status := range ValidOrderStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// Label retourne le libellé français du sujet.
func (s OrderStatus) Label() string {
	switch s {
	case StatusPending:
		return "en attente"
	case StatusPaid:
		return "payé"
	case StatusShipping:
		return "en cours de livraison"
	case StatusDelivered:
		return "livré"
	case StatusCancelled:
		return "annulé"
	default:
		return string(s)
	}
}

// CanTransitionTo autorise les changements de statut réalistes
// (admin : gestion des commandes).
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	if s == next {
		return true
	}
	switch s {
	case StatusPending:
		return next == StatusPaid || next == StatusCancelled
	case StatusPaid:
		return next == StatusShipping || next == StatusCancelled
	case StatusShipping:
		return next == StatusDelivered || next == StatusCancelled
	default:
		return false
	}
}

func (o Order) IsCancelled() bool {
	return o.Status == StatusCancelled
}

func (o Order) RequiresCancelReason(next OrderStatus) bool {
	return next == StatusCancelled
}
