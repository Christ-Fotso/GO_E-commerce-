package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ecommerce-cli/internal/api"
	"ecommerce-cli/internal/database"
	"ecommerce-cli/internal/handlers"
	"ecommerce-cli/internal/repositories"
)

func main() {
	dsn := "host=localhost port=5432 user=admin password=password dbname=ecommerce sslmode=disable"
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)

	// Routes déjà implémentées
	http.HandleFunc(api.RouteHealth, handlers.Health)
	http.HandleFunc(api.RouteRegister, authHandler.Register)
	http.HandleFunc(api.RouteLogin, authHandler.Login)
	http.HandleFunc(api.RouteConfirm, authHandler.Confirm)

	// Contrat figé : 501 tant que Dev 1 ne les implémente pas
	for _, route := range []string{
		api.RouteResetPassword,
		api.RouteProductList,
		api.RouteProductGet,
		api.RouteCartGet,
		api.RouteCartAddItem,
		api.RouteCartUpdate,
		api.RouteCartDelete,
		api.RouteCartPay,
		api.RouteOrderList,
		api.RouteOrderGet,
		api.RouteAdminUserList,
		api.RouteAdminUserCreate,
		api.RouteAdminUserUpdate,
		api.RouteAdminUserDelete,
		api.RouteAdminUserConfirm,
		api.RouteAdminProductCreate,
		api.RouteAdminOrderList,
		api.RouteAdminOrderCreate,
		api.RouteAdminOrderStatus,
	} {
		http.HandleFunc(route, handlers.NotImplemented)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = api.DefaultPort
	}

	addr := ":" + port
	fmt.Printf("Démarrage du serveur HTTP sur http://localhost%s ...\n", addr)
	fmt.Println("Test : curl.exe http://localhost" + addr + "/health")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Port %s déjà utilisé. Ferme l'autre terminal (Ctrl+C) ou relance avec :\n  $env:PORT=8081; go run ./cmd/server\nDétail: %v", port, err)
	}
}
