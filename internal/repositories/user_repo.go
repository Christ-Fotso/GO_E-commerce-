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
	if u.Role == "" {
		u.Role = models.RoleClient
	}

	query := `
		INSERT INTO users (username, email, password_hash, role, confirmation_code, is_confirmed, reset_code) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(
		query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.ConfirmationCode,
		u.IsConfirmed,
		u.ResetCode,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}
	query := `
		SELECT id, username, email, password_hash, role, confirmation_code, is_confirmed, reset_code, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.DB.QueryRow(query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.ConfirmationCode,
		&u.IsConfirmed,
		&u.ResetCode,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
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
