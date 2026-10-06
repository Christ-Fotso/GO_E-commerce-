package models

import "time"

// Rôles possibles d'un utilisateur
const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

// User représente un utilisateur de la plateforme e-commerce
type User struct {
	ID               int        `json:"id"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"` // "-" cache ce champ en JSON
	ConfirmationCode string     `json:"-"`
	IsConfirmed      bool       `json:"is_confirmed"`
	Role             string     `json:"role"`
	ResetCode        string     `json:"-"`
	ResetExpiresAt   *time.Time `json:"-"`
	CreatedAt        time.Time  `json:"created_at"`
}
