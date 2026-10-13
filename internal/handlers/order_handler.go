package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ecommerce-cli/internal/middleware"
	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"
)

type OrderHandler struct {
	Orders *repositories.OrderRepository
}

func NewOrderHandler(orders *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{Orders: orders}
}

// GET /orders : les commandes de l'utilisateur connecté, avec leur statut
func (handler *OrderHandler) List(writer http.ResponseWriter, request *http.Request) {
	orders, err := handler.Orders.List(middleware.CurrentUser(request).ID, "")
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture des commandes")
		return
	}
	writeJSON(writer, http.StatusOK, orders)
}

// GET /orders/{id} : accepte l'id numérique ou l'identifiant métier (CMD-XXXXXX)
func (handler *OrderHandler) Get(writer http.ResponseWriter, request *http.Request) {
	order, found := findOrder(writer, request, handler.Orders)
	if !found {
		return
	}
	// Un client ne voit que ses propres commandes (404 plutôt que 403 : on ne révèle pas qu'elle existe)
	if order.UserID != middleware.CurrentUser(request).ID {
		writeError(writer, http.StatusNotFound, "Commande introuvable")
		return
	}
	writeJSON(writer, http.StatusOK, order)
}

// findOrder cherche la commande désignée par {id} dans l'URL (partagé avec l'admin).
// En cas d'échec, la réponse d'erreur est déjà envoyée.
func findOrder(writer http.ResponseWriter, request *http.Request, orders *repositories.OrderRepository) (*models.Order, bool) {
	identifier := request.PathValue("id")

	var order *models.Order
	var err error
	if strings.HasPrefix(strings.ToUpper(identifier), "CMD-") {
		order, err = orders.GetByReference(strings.ToUpper(identifier))
	} else {
		orderID, convertError := strconv.Atoi(identifier)
		if convertError != nil {
			writeError(writer, http.StatusBadRequest, "Identifiant de commande invalide")
			return nil, false
		}
		order, err = orders.GetByID(orderID)
	}

	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Commande introuvable")
		return nil, false
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture de la commande")
		return nil, false
	}
	return order, true
}
