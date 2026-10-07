package api

import (
	"database/sql"
	"net/http"

	"ecommerce-cli/internal/handlers"
	"ecommerce-cli/internal/middleware"
	"ecommerce-cli/internal/repositories"
)

// NewRouter déclare toutes les routes de l'API (voir API.md).
// Routage natif de Go 1.22+ : "MÉTHODE /chemin/{paramètre}", sans framework.
func NewRouter(database *sql.DB) http.Handler {
	// Repositories : l'accès SQL
	users := repositories.NewUserRepository(database)
	sessions := repositories.NewSessionRepository(database)
	products := repositories.NewProductRepository(database)
	carts := repositories.NewCartRepository(database)
	orders := repositories.NewOrderRepository(database)

	// Handlers : la logique HTTP
	authHandler := handlers.NewAuthHandler(users, sessions)
	productHandler := handlers.NewProductHandler(products)
	cartHandler := handlers.NewCartHandler(carts, products, orders)
	orderHandler := handlers.NewOrderHandler(orders)
	adminHandler := handlers.NewAdminHandler(users, products, orders)

	// Middleware : vérifie le token (RequireUser) et le rôle admin (RequireAdmin)
	auth := middleware.NewAuth(sessions, users)

	router := http.NewServeMux()

	// Public
	router.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write([]byte(`{"status":"ok"}`))
	})
	router.HandleFunc("POST /register", authHandler.Register)
	router.HandleFunc("POST /confirm", authHandler.Confirm)
	router.HandleFunc("POST /login", authHandler.Login)
	router.HandleFunc("POST /password/forgot", authHandler.ForgotPassword)
	router.HandleFunc("POST /password/reset", authHandler.ResetPassword)
	router.HandleFunc("GET /products", productHandler.List)
	router.HandleFunc("GET /products/{id}", productHandler.Get)

	// Utilisateur connecté
	router.HandleFunc("POST /logout", auth.RequireUser(authHandler.Logout))
	router.HandleFunc("GET /me", auth.RequireUser(authHandler.Me))
	router.HandleFunc("GET /cart", auth.RequireUser(cartHandler.Get))
	router.HandleFunc("DELETE /cart", auth.RequireUser(cartHandler.Clear))
	router.HandleFunc("POST /cart/items", auth.RequireUser(cartHandler.AddItem))
	router.HandleFunc("PUT /cart/items/{product_id}", auth.RequireUser(cartHandler.UpdateItem))
	router.HandleFunc("DELETE /cart/items/{product_id}", auth.RequireUser(cartHandler.RemoveItem))
	router.HandleFunc("POST /cart/pay", auth.RequireUser(cartHandler.Pay))
	router.HandleFunc("GET /orders", auth.RequireUser(orderHandler.List))
	router.HandleFunc("GET /orders/{id}", auth.RequireUser(orderHandler.Get))

	// Administrateur
	router.HandleFunc("GET /admin/users", auth.RequireAdmin(adminHandler.ListUsers))
	router.HandleFunc("POST /admin/users", auth.RequireAdmin(adminHandler.CreateUser))
	router.HandleFunc("GET /admin/users/{id}", auth.RequireAdmin(adminHandler.GetUser))
	router.HandleFunc("PUT /admin/users/{id}", auth.RequireAdmin(adminHandler.UpdateUser))
	router.HandleFunc("DELETE /admin/users/{id}", auth.RequireAdmin(adminHandler.DeleteUser))
	router.HandleFunc("POST /admin/users/{id}/confirm", auth.RequireAdmin(adminHandler.ConfirmUser))
	router.HandleFunc("GET /admin/products", auth.RequireAdmin(adminHandler.ListProducts))
	router.HandleFunc("POST /admin/products", auth.RequireAdmin(adminHandler.CreateProduct))
	router.HandleFunc("GET /admin/orders", auth.RequireAdmin(adminHandler.ListOrders))
	router.HandleFunc("POST /admin/orders", auth.RequireAdmin(adminHandler.CreateOrder))
	router.HandleFunc("GET /admin/orders/{id}", auth.RequireAdmin(adminHandler.GetOrder))
	router.HandleFunc("PUT /admin/orders/{id}/status", auth.RequireAdmin(adminHandler.UpdateOrderStatus))

	return router
}
