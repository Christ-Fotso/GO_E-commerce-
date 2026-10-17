package client

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func customKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Input.Next = key.NewBinding(key.WithKeys("enter", "tab", "down"), key.WithHelp("enter/↓", "suivant"))
	km.Input.Prev = key.NewBinding(key.WithKeys("shift+tab", "up"), key.WithHelp("↑", "précédent"))
	km.Confirm.Next = key.NewBinding(key.WithKeys("enter", "tab", "down"), key.WithHelp("enter/↓", "suivant"))
	km.Confirm.Prev = key.NewBinding(key.WithKeys("shift+tab", "up"), key.WithHelp("↑", "précédent"))
	return km
}

func newRegisterForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("username").Title("Pseudo").Placeholder("demo"),
			huh.NewInput().Key("email").Title("Email").Placeholder("demo@example.com"),
			huh.NewInput().Key("password").Title("Mot de passe").EchoMode(huh.EchoModePassword),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newConfirmForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("email").Title("Email").Placeholder("demo@example.com"),
			huh.NewInput().Key("code").Title("Code de confirmation").Placeholder(mockConfirmCode),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newLoginForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("email").Title("Email").Placeholder("demo@example.com"),
			huh.NewInput().Key("password").Title("Mot de passe").EchoMode(huh.EchoModePassword),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newResetForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("email").Title("Email").Placeholder("demo@example.com"),
			huh.NewInput().Key("code").Title("Code reset (vide = demander un code)").Placeholder(mockResetCode),
			huh.NewInput().Key("new_password").Title("Nouveau mot de passe").EchoMode(huh.EchoModePassword),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newPaymentForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("card").
				Title("Numéro de carte").
				Placeholder("4242 4242 4242 4242").
				Validate(validateCardNumber),
			huh.NewInput().
				Key("expiry").
				Title("Date d'expiration (MM/YY)").
				Placeholder("12/30").
				Validate(validateExpiry),
			huh.NewInput().
				Key("cvc").
				Title("CVC").
				Placeholder("123").
				EchoMode(huh.EchoModePassword).
				Validate(validateCVC),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func validateCardNumber(s string) error {
	d := digitsOnly(s)
	if len(d) < 13 || len(d) > 19 {
		return fmt.Errorf("le numéro doit contenir 13 à 19 chiffres")
	}
	return nil
}

func validateExpiry(s string) error {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return fmt.Errorf("format attendu: MM/YY")
	}
	month, errM := strconv.Atoi(parts[0])
	_, errY := strconv.Atoi(parts[1])
	if errM != nil || errY != nil || month < 1 || month > 12 || len(parts[1]) != 2 {
		return fmt.Errorf("date invalide (MM/YY)")
	}
	return nil
}

func validateCVC(s string) error {
	d := digitsOnly(s)
	if len(d) < 3 || len(d) > 4 {
		return fmt.Errorf("CVC: 3 ou 4 chiffres")
	}
	return nil
}
