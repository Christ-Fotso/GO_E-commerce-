package admin

import (
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

func newUserCreateForm(username, email, password, role *string, isConfirmed *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nom d'utilisateur").
				Value(username).
				Placeholder("ex: jean_dupont"),
			huh.NewInput().
				Title("Adresse email").
				Value(email).
				Placeholder("ex: jean@example.com"),
			huh.NewInput().
				Title("Mot de passe").
				Value(password).
				EchoMode(huh.EchoModePassword),
			huh.NewSelect[string]().
				Title("Rôle").
				Options(
					huh.NewOption("Client (customer)", "customer"),
					huh.NewOption("Administrateur (admin)", "admin"),
				).
				Value(role),
			huh.NewConfirm().
				Title("Compte immédiatement confirmé ?").
				Value(isConfirmed),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newUserEditForm(username, email, password, role *string, isConfirmed *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nom d'utilisateur").
				Value(username),
			huh.NewInput().
				Title("Adresse email").
				Value(email),
			huh.NewInput().
				Title("Nouveau mot de passe (laisser vide pour ne pas changer)").
				Value(password).
				EchoMode(huh.EchoModePassword),
			huh.NewSelect[string]().
				Title("Rôle").
				Options(
					huh.NewOption("Client (customer)", "customer"),
					huh.NewOption("Administrateur (admin)", "admin"),
				).
				Value(role),
			huh.NewConfirm().
				Title("Compte confirmé ?").
				Value(isConfirmed),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newProductCreateForm(name, description, category, price, taxRate, stock *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nom du produit").
				Value(name).
				Placeholder("ex: Clavier Mécanique RGB"),
			huh.NewInput().
				Title("Description").
				Value(description).
				Placeholder("ex: Clavier mécanique compact"),
			huh.NewInput().
				Title("Catégorie").
				Value(category).
				Placeholder("ex: informatique"),
			huh.NewInput().
				Title("Prix HT (€)").
				Value(price).
				Placeholder("ex: 79.99"),
			huh.NewInput().
				Title("Taux TVA (%)").
				Value(taxRate).
				Placeholder("20.0"),
			huh.NewInput().
				Title("Stock initial").
				Value(stock).
				Placeholder("15"),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newOrderStatusForm(status, reason *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Nouveau Statut").
				Options(
					huh.NewOption("Payée (paid)", "paid"),
					huh.NewOption("En cours d'expédition (shipping)", "shipping"),
					huh.NewOption("Livrée (delivered)", "delivered"),
					huh.NewOption("Annulée (cancelled)", "cancelled"),
				).
				Value(status),
			huh.NewInput().
				Title("Motif d'annulation (obligatoire si statut = cancelled)").
				Value(reason).
				Placeholder("ex: Rupture de stock ou demande client"),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}

func newOrderCreateForm(userID, productID, quantity, status *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("ID du client (user_id)").
				Value(userID).
				Placeholder("ex: 2"),
			huh.NewInput().
				Title("ID du produit (product_id)").
				Value(productID).
				Placeholder("ex: 1"),
			huh.NewInput().
				Title("Quantité").
				Value(quantity).
				Placeholder("ex: 1"),
			huh.NewSelect[string]().
				Title("Statut initial").
				Options(
					huh.NewOption("En attente (pending)", "pending"),
					huh.NewOption("Payée (paid)", "paid"),
				).
				Value(status),
		),
	).WithKeyMap(customKeyMap()).WithShowHelp(true)
}
