package models

import "time"

// User représente un utilisateur de la plateforme e-commerce
type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"` // "-" cache ce champ en JSON
	ConfirmationCode string    `json:"-"` 
	IsConfirmed      bool      `json:"is_confirmed"`
	CreatedAt        time.Time `json:"created_at"`
}
