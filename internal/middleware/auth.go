package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"ecommerce-cli/internal/models"
	"ecommerce-cli/internal/repositories"
)

// Type privé pour la clé du contexte : évite tout conflit avec une autre clé "user"
type contextKey string

const userContextKey contextKey = "user"

// Auth vérifie le header "Authorization: Bearer <token>" avant d'appeler le vrai handler
type Auth struct {
	Sessions *repositories.SessionRepository
	Users    *repositories.UserRepository
}

func NewAuth(sessions *repositories.SessionRepository, users *repositories.UserRepository) *Auth {
	return &Auth{Sessions: sessions, Users: users}
}

// RequireUser laisse passer uniquement un utilisateur connecté.
// L'utilisateur trouvé est rangé dans le contexte de la requête pour le handler.
func (auth *Auth) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token, found := strings.CutPrefix(request.Header.Get("Authorization"), "Bearer ")
		if !found || token == "" {
			writeError(writer, http.StatusUnauthorized, "Token manquant (header Authorization: Bearer <token>)")
			return
		}

		userID, err := auth.Sessions.GetUserID(token)
		if err != nil {
			writeError(writer, http.StatusUnauthorized, "Token invalide ou expiré")
			return
		}

		user, err := auth.Users.GetByID(userID)
		if err != nil {
			writeError(writer, http.StatusUnauthorized, "Utilisateur introuvable")
			return
		}

		requestWithUser := request.WithContext(context.WithValue(request.Context(), userContextKey, user))
		next(writer, requestWithUser)
	}
}

// RequireAdmin laisse passer uniquement un utilisateur connecté ayant le rôle admin
func (auth *Auth) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return auth.RequireUser(func(writer http.ResponseWriter, request *http.Request) {
		if CurrentUser(request).Role != models.RoleAdmin {
			writeError(writer, http.StatusForbidden, "Accès réservé aux administrateurs")
			return
		}
		next(writer, request)
	})
}

// CurrentUser récupère l'utilisateur rangé dans le contexte par RequireUser
func CurrentUser(request *http.Request) *models.User {
	user, _ := request.Context().Value(userContextKey).(*models.User)
	return user
}

// Même format d'erreur que les handlers : {"error": "..."}
func writeError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(map[string]string{"error": message})
}
