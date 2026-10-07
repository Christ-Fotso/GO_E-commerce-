package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"ecommerce-cli/internal/middleware"
	"ecommerce-cli/internal/repositories"
)

// Toutes les routes du panier passent par middleware.RequireUser :
// le panier est toujours celui de l'utilisateur connecté.
type CartHandler struct {
	Carts    *repositories.CartRepository
	Products *repositories.ProductRepository
	Orders   *repositories.OrderRepository
}

func NewCartHandler(carts *repositories.CartRepository, products *repositories.ProductRepository, orders *repositories.OrderRepository) *CartHandler {
	return &CartHandler{Carts: carts, Products: products, Orders: orders}
}

// writeCart renvoie le panier à jour (après chaque modification, le client voit le nouveau total)
func (handler *CartHandler) writeCart(writer http.ResponseWriter, request *http.Request, status int) {
	cart, err := handler.Carts.GetOrCreate(middleware.CurrentUser(request).ID)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture du panier")
		return
	}
	writeJSON(writer, status, cart)
}

// GET /cart
func (handler *CartHandler) Get(writer http.ResponseWriter, request *http.Request) {
	handler.writeCart(writer, request, http.StatusOK)
}

// POST /cart/items {"product_id": 1, "quantity": 2}
func (handler *CartHandler) AddItem(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}
	if body.Quantity <= 0 {
		writeError(writer, http.StatusBadRequest, "La quantité doit être supérieure à 0")
		return
	}

	product, err := handler.Products.GetByID(body.ProductID)
	if err != nil {
		writeError(writer, http.StatusNotFound, "Produit introuvable")
		return
	}
	if product.Stock < body.Quantity {
		writeError(writer, http.StatusConflict, "Stock insuffisant")
		return
	}

	if err := handler.Carts.AddItem(middleware.CurrentUser(request).ID, body.ProductID, body.Quantity); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de l'ajout au panier")
		return
	}
	handler.writeCart(writer, request, http.StatusOK)
}

// PUT /cart/items/{product_id} {"quantity": 3} — une quantité de 0 retire le produit
func (handler *CartHandler) UpdateItem(writer http.ResponseWriter, request *http.Request) {
	productID, valid := pathID(writer, request, "product_id")
	if !valid {
		return
	}
	var body struct {
		Quantity int `json:"quantity"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}
	if body.Quantity < 0 {
		writeError(writer, http.StatusBadRequest, "La quantité ne peut pas être négative")
		return
	}

	userID := middleware.CurrentUser(request).ID
	var err error
	if body.Quantity == 0 {
		err = handler.Carts.RemoveItem(userID, productID)
	} else {
		err = handler.Carts.SetItemQuantity(userID, productID, body.Quantity)
	}
	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Ce produit n'est pas dans le panier")
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la modification du panier")
		return
	}
	handler.writeCart(writer, request, http.StatusOK)
}

// DELETE /cart/items/{product_id}
func (handler *CartHandler) RemoveItem(writer http.ResponseWriter, request *http.Request) {
	productID, valid := pathID(writer, request, "product_id")
	if !valid {
		return
	}
	err := handler.Carts.RemoveItem(middleware.CurrentUser(request).ID, productID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Ce produit n'est pas dans le panier")
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la suppression")
		return
	}
	handler.writeCart(writer, request, http.StatusOK)
}

// DELETE /cart : vide le panier
func (handler *CartHandler) Clear(writer http.ResponseWriter, request *http.Request) {
	if err := handler.Carts.Clear(middleware.CurrentUser(request).ID); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors du vidage du panier")
		return
	}
	handler.writeCart(writer, request, http.StatusOK)
}

// POST /cart/pay {"card_number": "4242 4242 4242 4242", "expiry": "12/27", "cvc": "123"}
// Paiement simulé (pas de Stripe) : on vérifie la carte, puis le panier devient une commande "paid".
func (handler *CartHandler) Pay(writer http.ResponseWriter, request *http.Request) {
	var card PaymentCard
	if !decodeJSON(writer, request, &card) {
		return
	}
	if message := card.Validate(); message != "" {
		writeError(writer, http.StatusPaymentRequired, message)
		return
	}

	// Les informations bancaires ne sont jamais enregistrées en base
	order, err := handler.Orders.CreateFromCart(middleware.CurrentUser(request).ID)
	if errors.Is(err, repositories.ErrEmptyCart) {
		writeError(writer, http.StatusBadRequest, "Le panier est vide")
		return
	}
	if errors.Is(err, repositories.ErrStockUnavailable) {
		writeError(writer, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors du paiement")
		return
	}
	writeJSON(writer, http.StatusCreated, order)
}
