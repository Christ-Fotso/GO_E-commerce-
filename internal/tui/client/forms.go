package client

import "github.com/charmbracelet/huh"

// Formulaires auth (huh) — prêts à être branchés sur les écrans.

func newRegisterForm(username, email, password *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Username").Value(username),
			huh.NewInput().Title("Email").Value(email),
			huh.NewInput().Title("Password").EchoMode(huh.EchoModePassword).Value(password),
		),
	)
}

func newConfirmForm(email, code *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Email").Value(email),
			huh.NewInput().Title("Code de confirmation").Value(code),
		),
	)
}

func newLoginForm(email, password *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Email").Value(email),
			huh.NewInput().Title("Password").EchoMode(huh.EchoModePassword).Value(password),
		),
	)
}

func newResetForm(email *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Email").Value(email),
		),
	)
}
