package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"ecommerce-cli/internal/middleware"
	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// Durée de validité d'un code de réinitialisation du mot de passe
const resetCodeDuration = 15 * time.Minute

const minPasswordLength = 6

type AuthHandler struct {
	Users    *repositories.UserRepository
	Sessions *repositories.SessionRepository
}

func NewAuthHandler(users *repositories.UserRepository, sessions *repositories.SessionRepository) *AuthHandler {
	return &AuthHandler{Users: users, Sessions: sessions}
}

// Les fonctions validate... renvoient un message d'erreur, ou "" si la valeur est valide

func validateUsername(username string) string {
	if strings.TrimSpace(username) == "" {
		return "Le pseudo est obligatoire"
	}
	return ""
}

func validateEmail(email string) string {
	if !strings.Contains(email, "@") {
		return "Email invalide"
	}
	return ""
}

func validatePassword(password string) string {
	if len(password) < minPasswordLength {
		return "Le mot de passe doit contenir au moins 6 caractères"
	}
	return ""
}

// firstError renvoie le premier message d'erreur non vide, ou "" si tout est valide
func firstError(messages ...string) string {
	for _, message := range messages {
		if message != "" {
			return message
		}
	}
	return ""
}

func validateUserFields(username, email, password string) string {
	return firstError(validateUsername(username), validateEmail(email), validatePassword(password))
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// POST /register
func (handler *AuthHandler) Register(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}
	if message := validateUserFields(body.Username, body.Email, body.Password); message != "" {
		writeError(writer, http.StatusBadRequest, message)
		return
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}
	confirmationCode, err := repositories.NewCode()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}

	user := &models.User{
		Username:         body.Username,
		Email:            body.Email,
		PasswordHash:     hash,
		ConfirmationCode: confirmationCode,
		IsConfirmed:      false,
		Role:             models.RoleCustomer,
	}
	if err := handler.Users.Create(user); err != nil {
		if repositories.IsUniqueViolation(err) {
			writeError(writer, http.StatusConflict, "Email ou pseudo déjà pris")
			return
		}
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la création du compte")
		return
	}

	// Pas de vrai serveur mail : on simule l'envoi de l'email dans la console du serveur
	// et on renvoie le code dans la réponse pour que le client CLI puisse l'afficher.
	log.Printf("[EMAIL SIMULÉ] à %s : ton code de confirmation est %s", user.Email, confirmationCode)

	writeJSON(writer, http.StatusCreated, map[string]any{
		"message":           "Compte créé. Confirme-le avec le code reçu.",
		"user":              user,
		"confirmation_code": confirmationCode,
	})
}

// POST /confirm
func (handler *AuthHandler) Confirm(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	user, err := handler.Users.GetByEmail(body.Email)
	if err != nil {
		writeError(writer, http.StatusNotFound, "Utilisateur introuvable")
		return
	}
	if user.IsConfirmed {
		writeMessage(writer, http.StatusOK, "Compte déjà confirmé")
		return
	}
	if body.Code == "" || user.ConfirmationCode != body.Code {
		writeError(writer, http.StatusUnauthorized, "Code de confirmation invalide")
		return
	}

	if err := handler.Users.ConfirmUser(body.Email); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la confirmation")
		return
	}
	writeMessage(writer, http.StatusOK, "Compte confirmé avec succès")
}

// POST /login : renvoie un token à envoyer ensuite dans "Authorization: Bearer <token>"
func (handler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	// Même message si l'email ou le mot de passe est faux : on ne révèle pas quels emails existent
	user, err := handler.Users.GetByEmail(body.Email)
	if err != nil {
		writeError(writer, http.StatusUnauthorized, "Email ou mot de passe incorrect")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		writeError(writer, http.StatusUnauthorized, "Email ou mot de passe incorrect")
		return
	}
	if !user.IsConfirmed {
		writeError(writer, http.StatusForbidden, "Compte non confirmé")
		return
	}

	token, err := handler.Sessions.Create(user.ID)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la création de la session")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

// POST /logout (connecté)
func (handler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	if err := handler.Sessions.Delete(token); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors de la déconnexion")
		return
	}
	writeMessage(writer, http.StatusOK, "Déconnecté")
}

// GET /me (connecté) : renvoie l'utilisateur lié au token
func (handler *AuthHandler) Me(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, middleware.CurrentUser(request))
}

// POST /password/forgot : génère un code de réinitialisation (sans être connecté)
func (handler *AuthHandler) ForgotPassword(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}

	if _, err := handler.Users.GetByEmail(body.Email); err != nil {
		writeError(writer, http.StatusNotFound, "Aucun compte avec cet email")
		return
	}

	resetCode, err := repositories.NewCode()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}
	if err := handler.Users.SetResetCode(body.Email, resetCode, time.Now().Add(resetCodeDuration)); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}

	log.Printf("[EMAIL SIMULÉ] à %s : ton code de réinitialisation est %s (valable 15 min)", body.Email, resetCode)

	writeJSON(writer, http.StatusOK, map[string]any{
		"message":    "Code de réinitialisation envoyé (valable 15 minutes)",
		"reset_code": resetCode,
	})
}

// POST /password/reset : change le mot de passe avec le code reçu
func (handler *AuthHandler) ResetPassword(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if !decodeJSON(writer, request, &body) {
		return
	}
	if message := validatePassword(body.NewPassword); message != "" {
		writeError(writer, http.StatusBadRequest, message)
		return
	}

	user, err := handler.Users.GetByEmail(body.Email)
	if err != nil || body.Code == "" || user.ResetCode != body.Code {
		writeError(writer, http.StatusUnauthorized, "Code de réinitialisation invalide")
		return
	}
	if user.ResetExpiresAt == nil || time.Now().After(*user.ResetExpiresAt) {
		writeError(writer, http.StatusUnauthorized, "Code de réinitialisation expiré")
		return
	}

	hash, err := hashPassword(body.NewPassword)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur interne")
		return
	}
	if err := handler.Users.ResetPassword(body.Email, hash); err != nil {
		writeError(writer, http.StatusInternalServerError, "Erreur lors du changement de mot de passe")
		return
	}
	writeMessage(writer, http.StatusOK, "Mot de passe modifié, tu peux te connecter")
}
