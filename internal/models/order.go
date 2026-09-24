package models

import "time"

// OrderStatus représente les différents états possibles d'une commande
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipping  OrderStatus = "shipping"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// Order représente une commande passée par un utilisateur
type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	Items     []OrderItem `json:"items,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// OrderItem représente un produit au sein d'une commande
type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	PriceAtTime float64 `json:"price_at_time"`
}
