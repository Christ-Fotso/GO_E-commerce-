package repositories

import (
	"database/sql"
	"ecommerce-cli/internal/models"
	"errors"
	"fmt"
)

var (
	ErrEmptyCart        = errors.New("le panier est vide")
	ErrStockUnavailable = errors.New("produit introuvable ou stock insuffisant")
)

type OrderRepository struct {
	Database *sql.DB
}

func NewOrderRepository(database *sql.DB) *OrderRepository {
	return &OrderRepository{Database: database}
}

// OrderLine est une ligne demandée : quel produit et combien
type OrderLine struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CreateFromCart transforme le panier de l'utilisateur en commande payée, puis supprime le panier.
// Tout se fait dans une transaction : si une étape échoue, rien n'est enregistré (Rollback).
func (repo *OrderRepository) CreateFromCart(userID int) (*models.Order, error) {
	transaction, err := repo.Database.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback() // sans effet si Commit a déjà été appelé

	query := `
		SELECT cart_items.product_id, cart_items.quantity
		FROM cart_items
		JOIN carts ON carts.id = cart_items.cart_id
		WHERE carts.user_id = $1`
	rows, err := transaction.Query(query, userID)
	if err != nil {
		return nil, err
	}
	lines := []OrderLine{}
	for rows.Next() {
		var line OrderLine
		if err := rows.Scan(&line.ProductID, &line.Quantity); err != nil {
			rows.Close()
			return nil, err
		}
		lines = append(lines, line)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, ErrEmptyCart
	}

	order, err := createOrderInTransaction(transaction, userID, lines, models.StatusPaid)
	if err != nil {
		return nil, err
	}

	// Le panier payé est supprimé : le prochain panier aura un nouvel identifiant BSK-XXXXXX
	if _, err := transaction.Exec(`DELETE FROM carts WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}

	return order, transaction.Commit()
}

// Create crée une commande pour n'importe quel client (utilisé par l'admin)
func (repo *OrderRepository) Create(userID int, lines []OrderLine, status models.OrderStatus) (*models.Order, error) {
	transaction, err := repo.Database.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	order, err := createOrderInTransaction(transaction, userID, lines, status)
	if err != nil {
		return nil, err
	}
	return order, transaction.Commit()
}

// createOrderInTransaction décrémente le stock de chaque produit, puis enregistre la commande et ses lignes
func createOrderInTransaction(transaction *sql.Tx, userID int, lines []OrderLine, status models.OrderStatus) (*models.Order, error) {
	order := &models.Order{UserID: userID, Status: status}

	// "AND stock >= $1" : la mise à jour n'a lieu que s'il reste assez de stock.
	// Si aucune ligne n'est mise à jour, Scan renvoie sql.ErrNoRows.
	stockQuery := `
		UPDATE products SET stock = stock - $1
		WHERE id = $2 AND stock >= $1
		RETURNING name, ` + priceTTCExpression

	for _, line := range lines {
		item := models.OrderItem{ProductID: line.ProductID, Quantity: line.Quantity}
		err := transaction.QueryRow(stockQuery, line.Quantity, line.ProductID).Scan(&item.ProductName, &item.PriceAtTime)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w (produit %d)", ErrStockUnavailable, line.ProductID)
		}
		if err != nil {
			return nil, err
		}
		order.Total += item.PriceAtTime * float64(item.Quantity)
		order.Items = append(order.Items, item)
	}
	order.Total = roundCents(order.Total)

	reference, err := NewReference("CMD")
	if err != nil {
		return nil, err
	}
	order.Reference = reference

	orderQuery := `INSERT INTO orders (reference, user_id, total, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err = transaction.QueryRow(orderQuery, order.Reference, userID, order.Total, order.Status).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price_at_time) VALUES ($1, $2, $3, $4) RETURNING id`
	for position := range order.Items {
		item := &order.Items[position]
		item.OrderID = order.ID
		err := transaction.QueryRow(itemQuery, order.ID, item.ProductID, item.Quantity, item.PriceAtTime).Scan(&item.ID)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

const orderColumns = `id, reference, user_id, total, status, COALESCE(cancel_reason, ''), created_at`

func scanOrder(row rowScanner) (*models.Order, error) {
	order := &models.Order{}
	err := row.Scan(&order.ID, &order.Reference, &order.UserID, &order.Total, &order.Status, &order.CancelReason, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID renvoie une commande avec ses produits
func (repo *OrderRepository) GetByID(orderID int) (*models.Order, error) {
	order, err := scanOrder(repo.Database.QueryRow(`SELECT `+orderColumns+` FROM orders WHERE id = $1`, orderID))
	if err != nil {
		return nil, err
	}
	return order, repo.loadItems(order)
}

// GetByReference renvoie une commande à partir de son identifiant métier (CMD-XXXXXX)
func (repo *OrderRepository) GetByReference(reference string) (*models.Order, error) {
	order, err := scanOrder(repo.Database.QueryRow(`SELECT `+orderColumns+` FROM orders WHERE reference = $1`, reference))
	if err != nil {
		return nil, err
	}
	return order, repo.loadItems(order)
}

// List renvoie les commandes, de la plus récente à la plus ancienne.
// userID = 0 : toutes les commandes (admin). status = "" : tous les statuts.
func (repo *OrderRepository) List(userID int, status models.OrderStatus) ([]models.Order, error) {
	query := `
		SELECT ` + orderColumns + ` FROM orders
		WHERE ($1 = 0 OR user_id = $1) AND ($2 = '' OR status::text = $2)
		ORDER BY created_at DESC, id DESC`
	rows, err := repo.Database.Query(query, userID, string(status))
	if err != nil {
		return nil, err
	}

	orders := []models.Order{}
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		orders = append(orders, *order)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// On charge les produits après avoir fermé "rows" (une seule requête à la fois par connexion)
	for position := range orders {
		if err := repo.loadItems(&orders[position]); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

// loadItems remplit les produits d'une commande (avec leur nom)
func (repo *OrderRepository) loadItems(order *models.Order) error {
	query := `
		SELECT order_items.id, order_items.product_id, products.name, order_items.quantity, order_items.price_at_time
		FROM order_items
		JOIN products ON products.id = order_items.product_id
		WHERE order_items.order_id = $1
		ORDER BY order_items.id`
	rows, err := repo.Database.Query(query, order.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	order.Items = []models.OrderItem{}
	for rows.Next() {
		item := models.OrderItem{OrderID: order.ID}
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.Quantity, &item.PriceAtTime); err != nil {
			return err
		}
		order.Items = append(order.Items, item)
	}
	return rows.Err()
}

// UpdateStatus change le statut d'une commande (la raison n'est gardée que pour une annulation).
// Une commande annulée remet ses produits en stock.
func (repo *OrderRepository) UpdateStatus(orderID int, status models.OrderStatus, cancelReason string) error {
	transaction, err := repo.Database.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	query := `UPDATE orders SET status = $1, cancel_reason = NULLIF($2, '') WHERE id = $3`
	result, err := transaction.Exec(query, status, cancelReason, orderID)
	if err != nil {
		return err
	}
	if err := checkRowsAffected(result); err != nil {
		return err
	}

	if status == models.StatusCancelled {
		// On additionne d'abord les quantités par produit (un produit peut apparaître sur plusieurs lignes)
		restockQuery := `
			UPDATE products SET stock = products.stock + returned.quantity
			FROM (
				SELECT product_id, SUM(quantity) AS quantity
				FROM order_items WHERE order_id = $1
				GROUP BY product_id
			) AS returned
			WHERE products.id = returned.product_id`
		if _, err := transaction.Exec(restockQuery, orderID); err != nil {
			return err
		}
	}
	return transaction.Commit()
}
