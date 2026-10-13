package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"ecommerce-cli/internal/api"
	"ecommerce-cli/internal/database"
	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// getEnv lit une variable d'environnement, avec une valeur par défaut si elle est absente
func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 1. Connexion à la base de données
	// dataSourceName = url de connexion à PostgreSQL
	dataSourceName := getEnv("DATABASE_URL", "host=localhost port=5432 user=admin password=password dbname=ecommerce sslmode=disable")
	connection, err := database.InitDB(dataSourceName)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB: %v", err)
	}
	defer connection.Close() // defer permet de fermer la DB à la toute fin (vu dans le Module 3)

	// 2. Compte administrateur par défaut (sinon personne ne peut utiliser /admin/...)
	adminEmail := getEnv("ADMIN_EMAIL", "admin@shop.local")
	adminPassword := getEnv("ADMIN_PASSWORD", "admin123")
	if err := ensureAdminExists(repositories.NewUserRepository(connection), adminEmail, adminPassword); err != nil {
		log.Fatalf("Impossible de créer le compte admin: %v", err)
	}

	// 3. Routes (toutes déclarées dans internal/api/routes.go)
	router := api.NewRouter(connection)

	// 4. Lancement du serveur
	port := ":" + getEnv("PORT", "8080")
	fmt.Printf("Démarrage du serveur HTTP sur http://localhost%s...\n", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Erreur au démarrage du serveur: %v\n", err)
	}
}

// ensureAdminExists crée le compte admin au premier démarrage s'il n'existe pas encore
func ensureAdminExists(users *repositories.UserRepository, email, password string) error {
	_, err := users.GetByEmail(email)
	if err == nil {
		return nil // déjà créé
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := &models.User{
		Username:     "admin",
		Email:        email,
		PasswordHash: string(hash),
		IsConfirmed:  true,
		Role:         models.RoleAdmin,
	}
	if err := users.Create(admin); err != nil {
		return err
	}
	log.Printf("Compte admin créé : %s / %s", email, password)
	return nil
}
