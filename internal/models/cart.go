package models

import "time"

const CartIDPrefix = "BSK"

// Cart est le panier d'un utilisateur (un panier par compte).
// Identifiant métier obligatoire (ex: BSK-1KH8E7).
// Les frais de livraison sont toujours gratuits : le total = somme TTC des lignes.
type Cart struct {
	ID         int        `json:"id"`
	BusinessID string     `json:"business_id"`
	UserID     int        `json:"user_id"`
	Items      []CartItem `json:"items,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CartItem est une ligne du panier (produit + quantité).
type CartItem struct {
	ID        int      `json:"id"`
	CartID    int      `json:"cart_id"`
	ProductID int      `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Product   *Product `json:"product,omitempty"`
}

// TotalTTC additionne les lignes. Livraison gratuite.
func (c Cart) TotalTTC() float64 {
	var total float64
	for _, item := range c.Items {
		if item.Product == nil {
			continue
		}
		total += item.Product.PriceTTC() * float64(item.Quantity)
	}
	return roundMoney(total)
}

func (c Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

func (c Cart) ItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}
