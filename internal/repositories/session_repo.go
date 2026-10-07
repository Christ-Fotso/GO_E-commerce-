package repositories

import (
	"database/sql"
	"time"
)

// Durée de validité d'un token après le login
const sessionDuration = 24 * time.Hour

type SessionRepository struct {
	Database *sql.DB
}

func NewSessionRepository(database *sql.DB) *SessionRepository {
	return &SessionRepository{Database: database}
}

// Create génère un nouveau token pour l'utilisateur et l'enregistre en base
func (repo *SessionRepository) Create(userID int) (string, error) {
	token, err := NewToken()
	if err != nil {
		return "", err
	}
	query := `INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`
	if _, err := repo.Database.Exec(query, token, userID, time.Now().Add(sessionDuration)); err != nil {
		return "", err
	}
	return token, nil
}

// GetUserID renvoie l'utilisateur lié au token, si le token existe et n'a pas expiré
func (repo *SessionRepository) GetUserID(token string) (int, error) {
	var userID int
	query := `SELECT user_id FROM sessions WHERE token = $1 AND expires_at > NOW()`
	err := repo.Database.QueryRow(query, token).Scan(&userID)
	return userID, err
}

// Delete supprime le token (déconnexion)
func (repo *SessionRepository) Delete(token string) error {
	_, err := repo.Database.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}
