package models

import (
	"math"
	"time"
)

const (
	ProductIDPrefix = "PDT"
	DefaultVATRate  = 0.20
)

// Product représente un produit vendu.
// Recherche prévue : nom, prix HT, description, catégorie, prix TTC.
// Identifiant métier obligatoire (ex: PDT-7D2K8N).
type Product struct {
	ID          int       `json:"id"`
	BusinessID  string    `json:"business_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	PriceHT     float64   `json:"price_ht"`
	VATRate     float64   `json:"vat_rate"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PriceTTC calcule le prix TTC à partir du HT et du taux de TVA.
func (p Product) PriceTTC() float64 {
	rate := p.VATRate
	if rate <= 0 {
		rate = DefaultVATRate
	}
	return roundMoney(p.PriceHT * (1 + rate))
}

func (p Product) HasStock(quantity int) bool {
	return quantity > 0 && p.Stock >= quantity
}

func roundMoney(amount float64) float64 {
	return math.Round(amount*100) / 100
}
