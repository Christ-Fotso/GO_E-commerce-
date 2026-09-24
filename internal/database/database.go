package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initialise la connexion à la base de données PostgreSQL
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture de la db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("impossible de ping la db: %w", err)
	}

	log.Println("Connexion à la base de données réussie")
	return db, nil
}
