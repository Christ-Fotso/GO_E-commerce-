package client

// APIClient will call the HTTP API described in API.md.
// Stub only for now: screens use mock data until Dev 1 exposes routes.
type APIClient struct {
	BaseURL string
}

// NewAPIClient returns a client pointed at the local server by default.
func NewAPIClient(baseURL string) *APIClient {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &APIClient{BaseURL: baseURL}
}

// TODO: GET /products, POST /register, POST /login, etc.
