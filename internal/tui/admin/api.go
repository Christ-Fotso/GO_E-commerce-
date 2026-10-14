package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"ecommerce-cli/internal/models"
)

type APIClient struct {
	BaseURL string
	HTTP    *http.Client
	Token   string
}

func NewAPIClient(baseURL string) *APIClient {
	if baseURL == "" {
		baseURL = os.Getenv("ECOMMERCE_API_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &APIClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *APIClient) SetToken(token string) {
	c.Token = token
}

type apiErrorResponse struct {
	Error string `json:"error"`
}

func (c *APIClient) doJSON(method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("API inaccessible : %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errResp apiErrorResponse
		if json.Unmarshal(data, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("%s", errResp.Error)
		}
		msg := strings.TrimSpace(string(data))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("%s", msg)
	}

	if out == nil || len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, out)
}

func (c *APIClient) Login(email, password string) (*models.User, error) {
	var out struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}
	err := c.doJSON(http.MethodPost, "/login", map[string]string{
		"email":    email,
		"password": password,
	}, &out)
	if err != nil {
		return nil, err
	}
	c.Token = out.Token
	return &out.User, nil
}

// --- Utilisateurs ---

func (c *APIClient) ListUsers() ([]models.User, error) {
	var users []models.User
	err := c.doJSON(http.MethodGet, "/admin/users", nil, &users)
	return users, err
}

func (c *APIClient) CreateUser(username, email, password, role string, isConfirmed bool) (*models.User, error) {
	var user models.User
	err := c.doJSON(http.MethodPost, "/admin/users", map[string]any{
		"username":     username,
		"email":        email,
		"password":     password,
		"role":         role,
		"is_confirmed": isConfirmed,
	}, &user)
	return &user, err
}

func (c *APIClient) UpdateUser(id int, fields map[string]any) (*models.User, error) {
	var user models.User
	err := c.doJSON(http.MethodPut, fmt.Sprintf("/admin/users/%d", id), fields, &user)
	return &user, err
}

func (c *APIClient) DeleteUser(id int) error {
	return c.doJSON(http.MethodDelete, fmt.Sprintf("/admin/users/%d", id), nil, nil)
}

func (c *APIClient) ConfirmUser(id int) (*models.User, error) {
	var user models.User
	err := c.doJSON(http.MethodPost, fmt.Sprintf("/admin/users/%d/confirm", id), nil, &user)
	return &user, err
}

// --- Produits ---

func (c *APIClient) ListProducts() ([]models.Product, error) {
	var products []models.Product
	err := c.doJSON(http.MethodGet, "/admin/products", nil, &products)
	return products, err
}

func (c *APIClient) CreateProduct(name, description, category string, price, taxRate float64, stock int) (*models.Product, error) {
	var prod models.Product
	err := c.doJSON(http.MethodPost, "/admin/products", map[string]any{
		"name":        name,
		"description": description,
		"category":    category,
		"price":       price,
		"tax_rate":    taxRate,
		"stock":       stock,
	}, &prod)
	return &prod, err
}

// --- Commandes ---

type OrderItemReq struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (c *APIClient) ListOrders(userID int, status string) ([]models.Order, error) {
	path := "/admin/orders"
	var params []string
	if userID > 0 {
		params = append(params, fmt.Sprintf("user_id=%d", userID))
	}
	if status != "" {
		params = append(params, fmt.Sprintf("status=%s", status))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	var orders []models.Order
	err := c.doJSON(http.MethodGet, path, nil, &orders)
	return orders, err
}

func (c *APIClient) GetOrder(id int) (*models.Order, error) {
	var order models.Order
	err := c.doJSON(http.MethodGet, fmt.Sprintf("/admin/orders/%d", id), nil, &order)
	return &order, err
}

func (c *APIClient) CreateOrder(userID int, items []OrderItemReq, status string) (*models.Order, error) {
	var order models.Order
	err := c.doJSON(http.MethodPost, "/admin/orders", map[string]any{
		"user_id": userID,
		"items":   items,
		"status":  status,
	}, &order)
	return &order, err
}

func (c *APIClient) UpdateOrderStatus(orderID int, status, reason string) (*models.Order, error) {
	var order models.Order
	err := c.doJSON(http.MethodPut, fmt.Sprintf("/admin/orders/%d/status", orderID), map[string]string{
		"status": status,
		"reason": reason,
	}, &order)
	return &order, err
}
