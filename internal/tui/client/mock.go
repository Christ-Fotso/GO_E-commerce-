package client

import (
	"time"

	"ecommerce-cli/internal/models"
)

const mockConfirmCode = "CODE123"
const mockResetCode = "RESET123"

// mockAccount stocke un user mock + mot de passe en clair (fallback UI).
type mockAccount struct {
	User     models.User
	Password string
}

func mockProducts() []models.Product {
	now := time.Now()
	return []models.Product{
		{ID: 1, Reference: "PDT-000001", Name: "Clavier mécanique", Description: "Switchs silencieux, layout AZERTY", Category: "périphériques", Price: 74.92, TaxRate: 20, PriceTTC: 89.90, Stock: 12, CreatedAt: now},
		{ID: 2, Reference: "PDT-000002", Name: "Souris sans fil", Description: "Ergonomique, 2.4 GHz", Category: "périphériques", Price: 32.92, TaxRate: 20, PriceTTC: 39.50, Stock: 30, CreatedAt: now},
		{ID: 3, Reference: "PDT-000003", Name: "Écran 27 pouces", Description: "IPS 144 Hz", Category: "écrans", Price: 207.50, TaxRate: 20, PriceTTC: 249.00, Stock: 5, CreatedAt: now},
		{ID: 4, Reference: "PDT-000004", Name: "Casque USB", Description: "Micro intégré, confortable", Category: "audio", Price: 49.99, TaxRate: 20, PriceTTC: 59.99, Stock: 18, CreatedAt: now},
		{ID: 5, Reference: "PDT-000005", Name: "Webcam HD", Description: "1080p avec micro", Category: "vidéo", Price: 37.50, TaxRate: 20, PriceTTC: 45.00, Stock: 22, CreatedAt: now},
	}
}

func mockOrders(userID int) []models.Order {
	now := time.Now()
	return []models.Order{
		{ID: 1001, Reference: "CMD-001001", UserID: userID, Total: 129.40, Status: models.StatusPaid, CreatedAt: now.Add(-72 * time.Hour)},
		{ID: 1002, Reference: "CMD-001002", UserID: userID, Total: 249.00, Status: models.StatusShipping, CreatedAt: now.Add(-48 * time.Hour)},
		{ID: 1003, Reference: "CMD-001003", UserID: userID, Total: 59.99, Status: models.StatusPending, CreatedAt: now.Add(-24 * time.Hour)},
		{ID: 1004, Reference: "CMD-001004", UserID: userID, Total: 45.00, Status: models.StatusDelivered, CreatedAt: now.Add(-96 * time.Hour)},
		{ID: 1005, Reference: "CMD-001005", UserID: userID, Total: 39.50, Status: models.StatusCancelled, CancelReason: "demande client", CreatedAt: now.Add(-120 * time.Hour)},
	}
}

func seedAccounts() []mockAccount {
	u := models.User{
		ID:               1,
		Username:         "demo",
		Email:            "demo@example.com",
		IsConfirmed:      true,
		Role:             models.RoleCustomer,
		ConfirmationCode: mockConfirmCode,
		CreatedAt:        time.Now(),
	}
	return []mockAccount{{User: u, Password: "demo123"}}
}

func productUnitTTC(p *models.Product) float64 {
	if p == nil {
		return 0
	}
	if p.PriceTTC > 0 {
		return p.PriceTTC
	}
	return p.Price
}
