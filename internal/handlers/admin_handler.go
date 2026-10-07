package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"ecommerce-cli/internal/middleware"
	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"
)

// Toutes les routes /admin/... passent par middleware.RequireAdmin
type AdminHandler struct {
	Users    *repositories.UserRepository
	Products *repositories.ProductRepository
	Orders   *repositories.OrderRepository
}

func NewAdminHandler(users *repositories.UserRepository, products *repositories.ProductRepository, orders *repositories.OrderRepository) *AdminHandler {
	return &AdminHandler{Users: users, Products: products, Orders: orders}
}

// ---------- Utilisateurs ----------

// GET /admin/users
func (handler *AdminHandler) ListUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := handler.Users.List()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture des utilisateurs")
		return
	}
	writeJSON(writer, http.StatusOK, users)
}

// GET /admin/users/{id}
func (handler *AdminHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	user, found := handler.findUser(writer, request)
	if !found {
		return
	}
	writeJSON(writer, http.StatusOK, user)
}

// POST /admin/users {"username", "email", "password", "role", "is_confirmed"}
func (handler *AdminHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Role        string `json:"role"`
		IsConfirmed bool   `json:"is_confirmed"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}
	if message := validateUserFields(body.Username, body.Email, body.Password); message != "" {
		writeError(writer, http.StatusBadRequest, message)
		return
	}
	if body.Role == "" {
		body.Role = models.RoleCustomer
	}
	if !isValidRole(body.Role) {
		writeError(writer, http.StatusBadRequest, "Rôle invalide (customer ou admin)")
		return
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}
	user := &models.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: hash,
		Role:         body.Role,
		IsConfirmed:  body.IsConfirmed,
	}
	// Un compte non confirmé reçoit un code, comme lors d'une inscription classique
	if !user.IsConfirmed {
		if user.ConfirmationCode, err = repositories.NewCode(); err != nil {
			writeError(writer, http.StatusInternalServerError, "Erreur interne")
			return
		}
	}

	if err := handler.Users.Create(user); err != nil {
		if repositories.IsUniqueViolation(err) {
			writeError(writer, http.StatusConflict, "Email ou pseudo déjà pris")
			return
		}
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la création")
		return
	}
	writeJSON(writer, http.StatusCreated, user)
}

// PUT /admin/users/{id} : seuls les champs envoyés sont modifiés
func (handler *AdminHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	user, found := handler.findUser(writer, request)
	if !found {
		return
	}

	// Des pointeurs : nil = champ absent du JSON = on garde l'ancienne valeur
	var body struct {
		Username    *string `json:"username"`
		Email       *string `json:"email"`
		Password    *string `json:"password"`
		Role        *string `json:"role"`
		IsConfirmed *bool   `json:"is_confirmed"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	if body.Username != nil {
		user.Username = *body.Username
	}
	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.Role != nil {
		if !isValidRole(*body.Role) {
			writeError(writer, http.StatusBadRequest, "Rôle invalide (customer ou admin)")
			return
		}
		user.Role = *body.Role
	}
	if body.IsConfirmed != nil {
		user.IsConfirmed = *body.IsConfirmed
	}

	if message := firstError(validateUsername(user.Username), validateEmail(user.Email)); message != "" {
		writeError(writer, http.StatusBadRequest, message)
		return
	}
	// Le mot de passe n'est vérifié et re-hashé que s'il a été envoyé
	if body.Password != nil {
		if message := validatePassword(*body.Password); message != "" {
			writeError(writer, http.StatusBadRequest, message)
			return
		}
		hash, err := hashPassword(*body.Password)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, "Erreur interne")
			return
		}
		user.PasswordHash = hash
	}

	if err := handler.Users.Update(user); err != nil {
		if repositories.IsUniqueViolation(err) {
			writeError(writer, http.StatusConflict, "Email ou pseudo déjà pris")
			return
		}
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la modification")
		return
	}
	writeJSON(writer, http.StatusOK, user)
}

// DELETE /admin/users/{id}
func (handler *AdminHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	userID, valid := pathID(writer, request, "id")
	if !valid {
		return
	}
	if userID == middleware.CurrentUser(request).ID {
		writeError(writer, http.StatusBadRequest, "Un administrateur ne peut pas supprimer son propre compte")
		return
	}

	err := handler.Users.Delete(userID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Utilisateur introuvable")
		return
	}
	if repositories.IsForeignKeyViolation(err) {
		writeError(writer, http.StatusConflict, "Impossible de supprimer : cet utilisateur a des commandes")
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la suppression")
		return
	}
	writeMessage(writer, http.StatusOK, "Utilisateur supprimé")
}

// POST /admin/users/{id}/confirm
func (handler *AdminHandler) ConfirmUser(writer http.ResponseWriter, request *http.Request) {
	user, found := handler.findUser(writer, request)
	if !found {
		return
	}
	if err := handler.Users.ConfirmUser(user.Email); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la confirmation")
		return
	}
	user.IsConfirmed = true
	writeJSON(writer, http.StatusOK, user)
}

func (handler *AdminHandler) findUser(writer http.ResponseWriter, request *http.Request) (*models.User, bool) {
	userID, valid := pathID(writer, request, "id")
	if !valid {
		return nil, false
	}
	user, err := handler.Users.GetByID(userID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Utilisateur introuvable")
		return nil, false
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture de l'utilisateur")
		return nil, false
	}
	return user, true
}

func isValidRole(role string) bool {
	return role == models.RoleCustomer || role == models.RoleAdmin
}

// ---------- Produits ----------

// GET /admin/products : tous les produits (même recherche que GET /products, sans filtre)
func (handler *AdminHandler) ListProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := handler.Products.Search(repositories.ProductFilter{})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture des produits")
		return
	}
	writeJSON(writer, http.StatusOK, products)
}

// POST /admin/products {"name", "description", "category", "price", "tax_rate", "stock"}
// L'identifiant métier PDT-XXXXXX est généré par le serveur.
func (handler *AdminHandler) CreateProduct(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Price       float64  `json:"price"`
		TaxRate     *float64 `json:"tax_rate"` // pointeur : absent = TVA par défaut de 20 %
		Stock       int      `json:"stock"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	taxRate := 20.0
	if body.TaxRate != nil {
		taxRate = *body.TaxRate
	}
	switch {
	case strings.TrimSpace(body.Name) == "":
		writeError(writer, http.StatusBadRequest, "Le nom est obligatoire")
		return
	case body.Price <= 0:
		writeError(writer, http.StatusBadRequest, "Le prix HT doit être supérieur à 0")
		return
	case taxRate < 0 || taxRate > 100:
		writeError(writer, http.StatusBadRequest, "La TVA doit être comprise entre 0 et 100")
		return
	case body.Stock < 0:
		writeError(writer, http.StatusBadRequest, "Le stock ne peut pas être négatif")
		return
	}
	if strings.TrimSpace(body.Category) == "" {
		body.Category = "divers"
	}

	product := &models.Product{
		Name:        body.Name,
		Description: body.Description,
		Category:    strings.ToLower(body.Category),
		Price:       body.Price,
		TaxRate:     taxRate,
		Stock:       body.Stock,
	}
	if err := handler.Products.Create(product); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la création du produit")
		return
	}
	writeJSON(writer, http.StatusCreated, product)
}

// ---------- Commandes ----------

// Changements de statut autorisés : pending → paid → shipping → delivered,
// et l'annulation est possible tant que la commande n'est pas livrée.
var allowedStatusChanges = map[models.OrderStatus][]models.OrderStatus{
	models.StatusPending:   {models.StatusPaid, models.StatusCancelled},
	models.StatusPaid:      {models.StatusShipping, models.StatusCancelled},
	models.StatusShipping:  {models.StatusDelivered, models.StatusCancelled},
	models.StatusDelivered: {},
	models.StatusCancelled: {},
}

func canChangeStatus(currentStatus, newStatus models.OrderStatus) bool {
	return slices.Contains(allowedStatusChanges[currentStatus], newStatus)
}

// GET /admin/orders?status=paid&user_id=3 (filtres optionnels)
func (handler *AdminHandler) ListOrders(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()

	userID := 0
	if rawUserID := params.Get("user_id"); rawUserID != "" {
		var err error
		if userID, err = strconv.Atoi(rawUserID); err != nil {
			writeError(writer, http.StatusBadRequest, "user_id doit être un nombre")
			return
		}
	}
	status := models.OrderStatus(params.Get("status"))
	if _, known := allowedStatusChanges[status]; status != "" && !known {
		writeError(writer, http.StatusBadRequest, "Statut inconnu")
		return
	}

	orders, err := handler.Orders.List(userID, status)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture des commandes")
		return
	}
	writeJSON(writer, http.StatusOK, orders)
}

// GET /admin/orders/{id} : id numérique ou CMD-XXXXXX, commande de n'importe quel client
func (handler *AdminHandler) GetOrder(writer http.ResponseWriter, request *http.Request) {
	order, found := findOrder(writer, request, handler.Orders)
	if !found {
		return
	}
	writeJSON(writer, http.StatusOK, order)
}

// POST /admin/orders {"user_id": 3, "items": [{"product_id": 1, "quantity": 2}], "status": "pending"}
// Crée une commande pour le compte de n'importe quel client.
func (handler *AdminHandler) CreateOrder(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		UserID int                      `json:"user_id"`
		Items  []repositories.OrderLine `json:"items"`
		Status models.OrderStatus       `json:"status"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	if _, err := handler.Users.GetByID(body.UserID); err != nil {
		writeError(writer, http.StatusNotFound, "Client introuvable")
		return
	}
	if len(body.Items) == 0 {
		writeError(writer, http.StatusBadRequest, "La commande doit contenir au moins un produit")
		return
	}
	for _, line := range body.Items {
		if line.Quantity <= 0 {
			writeError(writer, http.StatusBadRequest, "Chaque quantité doit être supérieure à 0")
			return
		}
	}
	if body.Status == "" {
		body.Status = models.StatusPending
	}
	if body.Status != models.StatusPending && body.Status != models.StatusPaid {
		writeError(writer, http.StatusBadRequest, "Une nouvelle commande est pending ou paid")
		return
	}

	order, err := handler.Orders.Create(body.UserID, body.Items, body.Status)
	if errors.Is(err, repositories.ErrStockUnavailable) {
		writeError(writer, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la création de la commande")
		return
	}
	writeJSON(writer, http.StatusCreated, order)
}

// PUT /admin/orders/{id}/status {"status": "cancelled", "reason": "Rupture de stock"}
func (handler *AdminHandler) UpdateOrderStatus(writer http.ResponseWriter, request *http.Request) {
	order, found := findOrder(writer, request, handler.Orders)
	if !found {
		return
	}
	var body struct {
		Status models.OrderStatus `json:"status"`
		Reason string             `json:"reason"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	if !canChangeStatus(order.Status, body.Status) {
		writeError(writer, http.StatusConflict,
			"Changement de statut impossible : "+string(order.Status)+" → "+string(body.Status))
		return
	}
	reason := strings.TrimSpace(body.Reason)
	if body.Status == models.StatusCancelled && reason == "" {
		writeError(writer, http.StatusBadRequest, "La raison est obligatoire pour annuler une commande")
		return
	}
	if body.Status != models.StatusCancelled {
		reason = ""
	}

	if err := handler.Orders.UpdateStatus(order.ID, body.Status, reason); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors du changement de statut")
		return
	}
	order.Status = body.Status
	order.CancelReason = reason
	writeJSON(writer, http.StatusOK, order)
}
