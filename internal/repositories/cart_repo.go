package repositories

import (
	"database/sql"
	"ecommerce-cli/internal/models"
	"errors"
	"math"
)

type CartRepository struct {
	Database *sql.DB
}

func NewCartRepository(database *sql.DB) *CartRepository {
	return &CartRepository{Database: database}
}

// GetOrCreate renvoie le panier de l'utilisateur (avec ses produits et son total TTC).
// S'il n'a pas encore de panier, on lui en crée un avec un identifiant métier BSK-XXXXXX.
func (repo *CartRepository) GetOrCreate(userID int) (*models.Cart, error) {
	cart, err := repo.findByUser(userID)
	if errors.Is(err, sql.ErrNoRows) {
		if err := repo.create(userID); err != nil {
			return nil, err
		}
		cart, err = repo.findByUser(userID)
	}
	if err != nil {
		return nil, err
	}

	if err := repo.loadItems(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (repo *CartRepository) findByUser(userID int) (*models.Cart, error) {
	cart := &models.Cart{}
	query := `SELECT id, reference, user_id, created_at FROM carts WHERE user_id = $1`
	err := repo.Database.QueryRow(query, userID).Scan(&cart.ID, &cart.Reference, &cart.UserID, &cart.CreatedAt)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (repo *CartRepository) create(userID int) error {
	reference, err := NewReference("BSK")
	if err != nil {
		return err
	}
	// ON CONFLICT DO NOTHING : si deux requêtes créent le panier en même temps (user_id est UNIQUE),
	// la seconde ne fait rien au lieu de planter
	query := `INSERT INTO carts (reference, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING`
	_, err = repo.Database.Exec(query, reference, userID)
	return err
}

// loadItems remplit les produits du panier et calcule le total TTC
func (repo *CartRepository) loadItems(cart *models.Cart) error {
	query := `
		SELECT cart_items.id, cart_items.quantity, ` + productColumns + `
		FROM cart_items
		JOIN products ON products.id = cart_items.product_id
		WHERE cart_items.cart_id = $1
		ORDER BY cart_items.id`

	rows, err := repo.Database.Query(query, cart.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	cart.Items = []models.CartItem{}
	cart.TotalTTC = 0
	for rows.Next() {
		item := models.CartItem{CartID: cart.ID, Product: &models.Product{}}
		product := item.Product
		err := rows.Scan(&item.ID, &item.Quantity,
			&product.ID, &product.Reference, &product.Name, &product.Description, &product.Category,
			&product.Price, &product.TaxRate, &product.PriceTTC, &product.Stock, &product.CreatedAt)
		if err != nil {
			return err
		}
		item.ProductID = product.ID
		item.LineTotalTTC = roundCents(product.PriceTTC * float64(item.Quantity))
		cart.TotalTTC += item.LineTotalTTC
		cart.Items = append(cart.Items, item)
	}
	cart.TotalTTC = roundCents(cart.TotalTTC)
	return rows.Err()
}

// AddItem ajoute un produit au panier. S'il y est déjà, on additionne les quantités.
func (repo *CartRepository) AddItem(userID, productID, quantity int) error {
	cart, err := repo.GetOrCreate(userID)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`
	_, err = repo.Database.Exec(query, cart.ID, productID, quantity)
	return err
}

// SetItemQuantity remplace la quantité d'un produit déjà présent dans le panier
func (repo *CartRepository) SetItemQuantity(userID, productID, quantity int) error {
	query := `
		UPDATE cart_items SET quantity = $1
		WHERE product_id = $2 AND cart_id = (SELECT id FROM carts WHERE user_id = $3)`
	result, err := repo.Database.Exec(query, quantity, productID, userID)
	if err != nil {
		return err
	}
	return checkRowsAffected(result)
}

// RemoveItem retire un produit du panier
func (repo *CartRepository) RemoveItem(userID, productID int) error {
	query := `DELETE FROM cart_items WHERE product_id = $1 AND cart_id = (SELECT id FROM carts WHERE user_id = $2)`
	result, err := repo.Database.Exec(query, productID, userID)
	if err != nil {
		return err
	}
	return checkRowsAffected(result)
}

// Clear vide le panier sans le supprimer
func (repo *CartRepository) Clear(userID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)`
	_, err := repo.Database.Exec(query, userID)
	return err
}

// roundCents arrondit un montant à 2 décimales (les centimes)
func roundCents(amount float64) float64 {
	return math.Round(amount*100) / 100
}
