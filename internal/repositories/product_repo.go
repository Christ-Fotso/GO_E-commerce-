package repositories

import (
	"database/sql"
	"ecommerce-cli/internal/models"
	"fmt"
	"strings"
)

type ProductRepository struct {
	Database *sql.DB
}

func NewProductRepository(database *sql.DB) *ProductRepository {
	return &ProductRepository{Database: database}
}

// Le prix TTC n'est pas stocké : il est calculé par PostgreSQL à partir du prix HT et de la TVA
const priceTTCExpression = `ROUND(products.price * (1 + products.tax_rate / 100), 2)`

// Colonnes préfixées par "products." pour pouvoir les réutiliser dans des jointures (panier, commandes)
const productColumns = `products.id, products.reference, products.name, COALESCE(products.description, ''),
	products.category, products.price, products.tax_rate, ` + priceTTCExpression + `, products.stock, products.created_at`

func scanProduct(row rowScanner) (*models.Product, error) {
	product := &models.Product{}
	err := row.Scan(&product.ID, &product.Reference, &product.Name, &product.Description, &product.Category,
		&product.Price, &product.TaxRate, &product.PriceTTC, &product.Stock, &product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// ProductFilter contient les critères de recherche (un champ vide ou à zéro est ignoré)
type ProductFilter struct {
	Query       string // cherche dans le nom OU la description
	Name        string
	Description string
	Category    string
	MinPrice    float64 // prix HT
	MaxPrice    float64
	MinPriceTTC float64
	MaxPriceTTC float64
}

// Search construit la clause WHERE au fur et à mesure selon les filtres remplis.
// Les valeurs passent toujours par des paramètres $1, $2... (jamais concaténées) pour éviter l'injection SQL.
func (repo *ProductRepository) Search(filter ProductFilter) ([]models.Product, error) {
	conditions := []string{}
	arguments := []any{}

	// addCondition ajoute une condition avec le bon numéro de paramètre ($1, $2, ...)
	addCondition := func(sqlCondition string, value any) {
		arguments = append(arguments, value)
		conditions = append(conditions, fmt.Sprintf(sqlCondition, len(arguments)))
	}

	if filter.Query != "" {
		addCondition("(name ILIKE $%[1]d OR description ILIKE $%[1]d)", "%"+filter.Query+"%")
	}
	if filter.Name != "" {
		addCondition("name ILIKE $%d", "%"+filter.Name+"%")
	}
	if filter.Description != "" {
		addCondition("description ILIKE $%d", "%"+filter.Description+"%")
	}
	if filter.Category != "" {
		addCondition("category ILIKE $%d", filter.Category)
	}
	if filter.MinPrice > 0 {
		addCondition("price >= $%d", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		addCondition("price <= $%d", filter.MaxPrice)
	}
	if filter.MinPriceTTC > 0 {
		addCondition(priceTTCExpression+" >= $%d", filter.MinPriceTTC)
	}
	if filter.MaxPriceTTC > 0 {
		addCondition(priceTTCExpression+" <= $%d", filter.MaxPriceTTC)
	}

	query := `SELECT ` + productColumns + ` FROM products`
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY name`

	rows, err := repo.Database.Query(query, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, rows.Err()
}

func (repo *ProductRepository) GetByID(productID int) (*models.Product, error) {
	return scanProduct(repo.Database.QueryRow(`SELECT `+productColumns+` FROM products WHERE id = $1`, productID))
}

func (repo *ProductRepository) GetByReference(reference string) (*models.Product, error) {
	return scanProduct(repo.Database.QueryRow(`SELECT `+productColumns+` FROM products WHERE reference = $1`, reference))
}

// Create ajoute un produit et lui génère son identifiant métier PDT-XXXXXX
func (repo *ProductRepository) Create(product *models.Product) error {
	reference, err := NewReference("PDT")
	if err != nil {
		return err
	}
	product.Reference = reference

	query := `
		INSERT INTO products (reference, name, description, category, price, tax_rate, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, ` + priceTTCExpression + `, created_at`

	return repo.Database.QueryRow(query, product.Reference, product.Name, product.Description, product.Category,
		product.Price, product.TaxRate, product.Stock).Scan(&product.ID, &product.PriceTTC, &product.CreatedAt)
}
