package models

import "time"

// Cart représente le panier d'un utilisateur
type Cart struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Items     []CartItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// CartItem représente un élément dans le panier
type CartItem struct {
	ID        int     `json:"id"`
	CartID    int     `json:"cart_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   *Product `json:"product,omitempty"`
}
