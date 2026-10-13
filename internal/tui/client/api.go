package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"ecommerce-cli/internal/models"
)

// APIClient calls the HTTP API described in API.md.
type APIClient struct {
	BaseURL string
	HTTP    *http.Client
	Token   string
}

// NewAPIClient returns a client pointed at the local server by default.
// Override with env ECOMMERCE_API_URL.
func NewAPIClient(baseURL string) *APIClient {
	if baseURL == "" {
		baseURL = os.Getenv("ECOMMERCE_API_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &APIClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 4 * time.Second},
	}
}

func (c *APIClient) SetToken(token string) {
	c.Token = token
}

func (c *APIClient) ClearToken() {
	c.Token = ""
}

type apiErrorBody struct {
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
		return fmt.Errorf("api unreachable: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errBody apiErrorBody
		if json.Unmarshal(data, &errBody) == nil && errBody.Error != "" {
			return fmt.Errorf("%s", errBody.Error)
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
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("réponse JSON invalide: %w", err)
	}
	return nil
}

// Ping reports whether the API responds (GET /health per API routes).
func (c *APIClient) Ping() bool {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// --- Auth (API.md) ---

type registerResponse struct {
	Message          string      `json:"message"`
	User             models.User `json:"user"`
	ConfirmationCode string      `json:"confirmation_code"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type forgotResponse struct {
	Message   string `json:"message"`
	ResetCode string `json:"reset_code"`
}

func (c *APIClient) Register(username, email, password string) (*registerResponse, error) {
	var out registerResponse
	err := c.doJSON(http.MethodPost, "/register", map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}, &out)
	return &out, err
}

func (c *APIClient) Confirm(email, code string) error {
	var out messageResponse
	return c.doJSON(http.MethodPost, "/confirm", map[string]string{
		"email": email,
		"code":  code,
	}, &out)
}

func (c *APIClient) Login(email, password string) (*loginResponse, error) {
	var out loginResponse
	err := c.doJSON(http.MethodPost, "/login", map[string]string{
		"email":    email,
		"password": password,
	}, &out)
	if err == nil {
		c.SetToken(out.Token)
	}
	return &out, err
}

func (c *APIClient) Logout() error {
	err := c.doJSON(http.MethodPost, "/logout", nil, nil)
	c.ClearToken()
	return err
}

func (c *APIClient) ForgotPassword(email string) (*forgotResponse, error) {
	var out forgotResponse
	err := c.doJSON(http.MethodPost, "/password/forgot", map[string]string{
		"email": email,
	}, &out)
	return &out, err
}

func (c *APIClient) ResetPassword(email, code, newPassword string) error {
	var out messageResponse
	return c.doJSON(http.MethodPost, "/password/reset", map[string]string{
		"email":        email,
		"code":         code,
		"new_password": newPassword,
	}, &out)
}

// --- Products ---

func (c *APIClient) ListProducts(query string) ([]models.Product, error) {
	path := "/products"
	if strings.TrimSpace(query) != "" {
		path += "?q=" + url.QueryEscape(query)
	}
	var out []models.Product
	err := c.doJSON(http.MethodGet, path, nil, &out)
	return out, err
}

func (c *APIClient) GetProduct(id int) (*models.Product, error) {
	var out models.Product
	err := c.doJSON(http.MethodGet, fmt.Sprintf("/products/%d", id), nil, &out)
	return &out, err
}

// --- Cart ---

func (c *APIClient) GetCart() (*models.Cart, error) {
	var out models.Cart
	err := c.doJSON(http.MethodGet, "/cart", nil, &out)
	return &out, err
}

func (c *APIClient) AddCartItem(productID, quantity int) (*models.Cart, error) {
	var out models.Cart
	err := c.doJSON(http.MethodPost, "/cart/items", map[string]int{
		"product_id": productID,
		"quantity":   quantity,
	}, &out)
	return &out, err
}

func (c *APIClient) UpdateCartItem(productID, quantity int) (*models.Cart, error) {
	var out models.Cart
	err := c.doJSON(http.MethodPut, fmt.Sprintf("/cart/items/%d", productID), map[string]int{
		"quantity": quantity,
	}, &out)
	return &out, err
}

func (c *APIClient) DeleteCartItem(productID int) (*models.Cart, error) {
	var out models.Cart
	err := c.doJSON(http.MethodDelete, fmt.Sprintf("/cart/items/%d", productID), nil, &out)
	return &out, err
}

func (c *APIClient) PayCart(cardNumber, expiry, cvc string) (*models.Order, error) {
	var out models.Order
	err := c.doJSON(http.MethodPost, "/cart/pay", map[string]string{
		"card_number": cardNumber,
		"expiry":      expiry,
		"cvc":         cvc,
	}, &out)
	return &out, err
}

// --- Orders ---

func (c *APIClient) ListOrders() ([]models.Order, error) {
	var out []models.Order
	err := c.doJSON(http.MethodGet, "/orders", nil, &out)
	return out, err
}
