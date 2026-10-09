package client

import (
	"time"

	"ecommerce-cli/internal/models"
)

// Données en dur pour développer l'UI avant le branchement API.

func mockProducts() []models.Product {
	now := time.Now()
	return []models.Product{
		{ID: 1, Name: "Clavier mecanique", Description: "Switchs silencieux, layout AZERTY", Price: 89.90, Stock: 12, CreatedAt: now},
		{ID: 2, Name: "Souris sans fil", Description: "Ergonomique, 2.4 GHz", Price: 39.50, Stock: 30, CreatedAt: now},
		{ID: 3, Name: "Ecran 27 pouces", Description: "IPS 144 Hz", Price: 249.00, Stock: 5, CreatedAt: now},
		{ID: 4, Name: "Casque USB", Description: "Micro integre, confortable", Price: 59.99, Stock: 18, CreatedAt: now},
	}
}

func mockOrders(userID int) []models.Order {
	now := time.Now()
	return []models.Order{
		{
			ID: 1001, UserID: userID, Total: 129.40, Status: models.StatusPaid, CreatedAt: now.Add(-48 * time.Hour),
			Items: []models.OrderItem{
				{ID: 1, OrderID: 1001, ProductID: 1, Quantity: 1, PriceAtTime: 89.90},
				{ID: 2, OrderID: 1001, ProductID: 2, Quantity: 1, PriceAtTime: 39.50},
			},
		},
		{
			ID: 1002, UserID: userID, Total: 249.00, Status: models.StatusShipping, CreatedAt: now.Add(-24 * time.Hour),
			Items: []models.OrderItem{
				{ID: 3, OrderID: 1002, ProductID: 3, Quantity: 1, PriceAtTime: 249.00},
			},
		},
		{
			ID: 1003, UserID: userID, Total: 59.99, Status: models.StatusPending, CreatedAt: now,
			Items: []models.OrderItem{
				{ID: 4, OrderID: 1003, ProductID: 4, Quantity: 1, PriceAtTime: 59.99},
			},
		},
	}
}

func mockUser() *models.User {
	return &models.User{
		ID:          1,
		Username:    "demo",
		Email:       "demo@example.com",
		IsConfirmed: true,
		CreatedAt:   time.Now(),
	}
}

const mockConfirmCode = "CODE123"
