package models

import "time"

// Product représente un produit vendu
type Product struct {
	ID          int       `json:"id"`
	Reference   string    `json:"reference"` // identifiant métier PDT-XXXXXX
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`     // prix HT
	TaxRate     float64   `json:"tax_rate"`  // TVA en %
	PriceTTC    float64   `json:"price_ttc"` // calculé : price * (1 + tax_rate/100)
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}
