package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSON envoie une réponse JSON avec le code HTTP donné
func writeJSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(data)
}

// writeError envoie une erreur au format {"error": "..."} : les CLI n'ont qu'un seul format à lire
func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}

// writeMessage envoie une confirmation simple au format {"message": "..."}
func writeMessage(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"message": message})
}

// decodeJSON lit le corps JSON de la requête dans "destination".
// En cas d'erreur, la réponse 400 est déjà envoyée : le handler n'a plus qu'à faire return.
func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) bool {
	if err := json.NewDecoder(request.Body).Decode(destination); err != nil {
		writeError(writer, http.StatusBadRequest, "Requête invalide : JSON mal formé")
		return false
	}
	return true
}

// pathID lit un identifiant numérique dans l'URL, ex : "GET /admin/users/{id}"
func pathID(writer http.ResponseWriter, request *http.Request, name string) (int, bool) {
	identifier, err := strconv.Atoi(request.PathValue(name))
	if err != nil || identifier <= 0 {
		writeError(writer, http.StatusBadRequest, "Identifiant invalide dans l'URL")
		return 0, false
	}
	return identifier, true
}
