package repositories

import (
	"database/sql"
	"ecommerce-cli/internal/models"
	"time"
)

type UserRepository struct {
	Database *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{Database: database}
}

// Colonnes lues à chaque SELECT, dans le même ordre que scanUser.
// COALESCE remplace les NULL par ” pour pouvoir les lire dans un string.
const userColumns = `id, username, email, password_hash, COALESCE(confirmation_code, ''), is_confirmed,
	role, COALESCE(reset_code, ''), reset_expires_at, created_at`

// rowScanner est implémenté à la fois par *sql.Row (une ligne) et *sql.Rows (plusieurs lignes)
type rowScanner interface {
	Scan(destinations ...any) error
}

func scanUser(row rowScanner) (*models.User, error) {
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.ConfirmationCode, &user.IsConfirmed,
		&user.Role, &user.ResetCode, &user.ResetExpiresAt, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Create(user *models.User) error {
	if user.Role == "" {
		user.Role = models.RoleCustomer
	}
	query := `
		INSERT INTO users (username, email, password_hash, confirmation_code, is_confirmed, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return repo.Database.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.ConfirmationCode, user.IsConfirmed, user.Role).
		Scan(&user.ID, &user.CreatedAt)
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	return scanUser(repo.Database.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = $1`, email))
}

func (repo *UserRepository) GetByID(userID int) (*models.User, error) {
	return scanUser(repo.Database.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = $1`, userID))
}

func (repo *UserRepository) List() ([]models.User, error) {
	rows, err := repo.Database.Query(`SELECT ` + userColumns + ` FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, rows.Err()
}

func (repo *UserRepository) ConfirmUser(email string) error {
	query := `UPDATE users SET is_confirmed = TRUE, confirmation_code = '' WHERE email = $1`
	_, err := repo.Database.Exec(query, email)
	return err
}

func (repo *UserRepository) UpdatePassword(email, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE email = $2`
	_, err := repo.Database.Exec(query, newPasswordHash, email)
	return err
}

// SetResetCode enregistre le code de réinitialisation et sa date d'expiration
func (repo *UserRepository) SetResetCode(email, resetCode string, expiresAt time.Time) error {
	query := `UPDATE users SET reset_code = $1, reset_expires_at = $2 WHERE email = $3`
	_, err := repo.Database.Exec(query, resetCode, expiresAt, email)
	return err
}

// ResetPassword change le mot de passe et invalide le code de réinitialisation
func (repo *UserRepository) ResetPassword(email, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1, reset_code = NULL, reset_expires_at = NULL WHERE email = $2`
	_, err := repo.Database.Exec(query, newPasswordHash, email)
	return err
}

// Update modifie un utilisateur existant (utilisé par l'admin)
func (repo *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, role = $4, is_confirmed = $5 WHERE id = $6`
	result, err := repo.Database.Exec(query, user.Username, user.Email, user.PasswordHash, user.Role, user.IsConfirmed, user.ID)
	if err != nil {
		return err
	}
	return checkRowsAffected(result)
}

func (repo *UserRepository) Delete(userID int) error {
	result, err := repo.Database.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	return checkRowsAffected(result)
}

// checkRowsAffected renvoie sql.ErrNoRows si aucune ligne n'a été modifiée (id inexistant)
func checkRowsAffected(result sql.Result) error {
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
