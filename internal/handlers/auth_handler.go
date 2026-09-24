package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"
	
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo *repositories.UserRepository
}

func NewAuthHandler(repo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Hasher le mot de passe
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username:         req.Username,
		Email:            req.Email,
		PasswordHash:     string(hash),
		ConfirmationCode: "CODE123", // Code statique pour l'instant (à générer aléatoirement)
		IsConfirmed:      false,
	}

	if err := h.Repo.Create(user); err != nil {
		http.Error(w, "Erreur lors de la création (email ou pseudo déjà pris ?)", http.StatusConflict)
		return
	}

	// Message clair, façon Module 10
	fmt.Fprintf(w, "Utilisateur %s créé ! Ton code de confirmation est : %s\n", user.Username, user.ConfirmationCode)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	if !user.IsConfirmed {
		http.Error(w, "Compte non confirmé", http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Bienvenue %s ! Tu es connecté.\n", user.Username)
}

func (h *AuthHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByEmail(req.Email)
	if err != nil || user.ConfirmationCode != req.Code {
		http.Error(w, "Code invalide ou utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	if err := h.Repo.ConfirmUser(req.Email); err != nil {
		http.Error(w, "Erreur lors de la confirmation", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Ton compte %s a été confirmé avec succès !\n", req.Email)
}
