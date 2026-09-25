package models

import "time"

// UserRole distingue un client d'un administrateur (CLIs séparés dans le sujet).
type UserRole string

const (
	RoleClient UserRole = "client"
	RoleAdmin  UserRole = "admin"
)

// User représente un utilisateur de la plateforme e-commerce.
// Inscription, connexion, confirmation par code et reset de mot de passe
// se font à partir de l'email (sujet : Authentification).
type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	Role             UserRole  `json:"role"`
	ConfirmationCode string    `json:"-"`
	IsConfirmed      bool      `json:"is_confirmed"`
	ResetCode        string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (r UserRole) IsValid() bool {
	return r == RoleClient || r == RoleAdmin
}
