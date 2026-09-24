package repositories

import (
	"database/sql"
	"ecommerce-cli/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(u *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, confirmation_code, is_confirmed) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at`
	
	err := r.DB.QueryRow(query, u.Username, u.Email, u.PasswordHash, u.ConfirmationCode, u.IsConfirmed).Scan(&u.ID, &u.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, username, email, password_hash, confirmation_code, is_confirmed, created_at FROM users WHERE email = $1`
	
	err := r.DB.QueryRow(query, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.ConfirmationCode, &u.IsConfirmed, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return u, nil
}

func (r *UserRepository) ConfirmUser(email string) error {
	query := `UPDATE users SET is_confirmed = TRUE, confirmation_code = '' WHERE email = $1`
	_, err := r.DB.Exec(query, email)
	return err
}

func (r *UserRepository) UpdatePassword(email, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE email = $2`
	_, err := r.DB.Exec(query, newPasswordHash, email)
	return err
}
