package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"
)

type ProductHandler struct {
	Products *repositories.ProductRepository
}

func NewProductHandler(products *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{Products: products}
}

// GET /products?q=&name=&description=&category=&min_price=&max_price=&min_price_ttc=&max_price_ttc=
func (handler *ProductHandler) List(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	filter := repositories.ProductFilter{
		Query:       params.Get("q"),
		Name:        params.Get("name"),
		Description: params.Get("description"),
		Category:    params.Get("category"),
	}

	// Chaque filtre de prix est optionnel ; s'il est présent, il doit être un nombre
	priceFilters := map[string]*float64{
		"min_price":     &filter.MinPrice,
		"max_price":     &filter.MaxPrice,
		"min_price_ttc": &filter.MinPriceTTC,
		"max_price_ttc": &filter.MaxPriceTTC,
	}
	for paramName, destination := range priceFilters {
		value, valid := parsePriceParam(params, paramName)
		if !valid {
			writeError(writer, http.StatusBadRequest, "Le paramètre "+paramName+" doit être un nombre")
			return
		}
		*destination = value
	}

	products, err := handler.Products.Search(filter)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la recherche")
		return
	}
	writeJSON(writer, http.StatusOK, products)
}

// parsePriceParam lit un prix dans l'URL. Paramètre absent = 0 (filtre ignoré).
func parsePriceParam(params url.Values, paramName string) (float64, bool) {
	rawValue := params.Get(paramName)
	if rawValue == "" {
		return 0, true
	}
	value, err := strconv.ParseFloat(strings.Replace(rawValue, ",", ".", 1), 64)
	if err != nil || value < 0 {
		return 0, false
	}
	return value, true
}

// GET /products/{id} : accepte l'id numérique ou l'identifiant métier (PDT-XXXXXX)
func (handler *ProductHandler) Get(writer http.ResponseWriter, request *http.Request) {
	identifier := request.PathValue("id")

	var product *models.Product
	var err error
	if strings.HasPrefix(strings.ToUpper(identifier), "PDT-") {
		product, err = handler.Products.GetByReference(strings.ToUpper(identifier))
	} else {
		productID, convertError := strconv.Atoi(identifier)
		if convertError != nil {
			writeError(writer, http.StatusBadRequest, "Identifiant de produit invalide")
			return
		}
		product, err = handler.Products.GetByID(productID)
	}

	if errors.Is(err, sql.ErrNoRows) {
		writeError(writer, http.StatusNotFound, "Produit introuvable")
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la lecture du produit")
		return
	}
	writeJSON(writer, http.StatusOK, product)
}
