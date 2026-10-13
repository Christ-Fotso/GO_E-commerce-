package models

import "time"

// Product représente un produit vendu
type Product struct {
	ID          int       `json:"id"`
	Reference   string    `json:"reference"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`     // HT
	TaxRate     float64   `json:"tax_rate"`  // %
	PriceTTC    float64   `json:"price_ttc"` // TTC
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}
