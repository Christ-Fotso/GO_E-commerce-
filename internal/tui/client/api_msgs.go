package client

import "ecommerce-cli/internal/models"

// SourceAPI / SourceMock indicate where data came from.
const (
	SourceAPI  = "api"
	SourceMock = "mock"
)

type registerDoneMsg struct {
	username string
	email    string
	password string
	code     string
	source   string
	err      error
}

type confirmDoneMsg struct {
	email  string
	source string
	err    error
}

type loginDoneMsg struct {
	user   models.User
	token  string
	source string
	err    error
}

type resetDoneMsg struct {
	message string
	source  string
	err     error
}

type productsDoneMsg struct {
	products []models.Product
	source   string
	err      error
}

type cartDoneMsg struct {
	cart    models.Cart
	source  string
	okMsg   string
	err     error
	goPay   bool
}

type payDoneMsg struct {
	order  models.Order
	source string
	err    error
}

type ordersDoneMsg struct {
	orders []models.Order
	source string
	err    error
}

type logoutDoneMsg struct {
	source string
	err    error
}

func sourceLabel(source string) string {
	if source == SourceAPI {
		return "API"
	}
	return "mock"
}
