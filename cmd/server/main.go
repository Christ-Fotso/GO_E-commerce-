package main

import (
	"fmt"
	"log"
	"net/http"

	"ecommerce-cli/internal/database"
	"ecommerce-cli/internal/handlers"
	"ecommerce-cli/internal/repositories"
)

func main() {
	// 1. Connexion à la base de données
	// dsn = Data Source Name (url de connexion)
	dsn := "host=localhost port=5432 user=admin password=password dbname=ecommerce sslmode=disable"
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB: %v", err)
	}
	defer db.Close() // defer permet de fermer la DB à la toute fin (vu dans le Module 3)

	// 2. Initialisation des Repositories
	userRepo := repositories.NewUserRepository(db)

	// 3. Initialisation des Handlers
	authHandler := handlers.NewAuthHandler(userRepo)

	// 4. Définition des routes (Style Module 10 avec Go 1.22+)
	http.HandleFunc("POST /register", authHandler.Register)
	http.HandleFunc("POST /login", authHandler.Login)
	http.HandleFunc("POST /confirm", authHandler.Confirm)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur l'API E-Commerce ! L'environnement est prêt.")
	})

	// 5. Lancement du serveur
	port := ":8080"
	fmt.Printf("Démarrage du serveur HTTP sur http://localhost%s...\n", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erreur au démarrage du serveur: %v\n", err)
	}
}
