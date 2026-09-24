package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID    int
	Email string
	Name  string
}

func main() {
	users := []User{
		{
			ID:    1,
			Email: "first@domain.com",
			Name:  "First User",
		},
		{
			ID:    2,
			Email: "second@domain.com",
			Name:  "Second User",
		},
		{
			ID:    3,
			Email: "third@domain.com",
			Name:  "Third User",
		},
	}

	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if len(strings.Trim(name, " ")) == 0 {
			fmt.Fprintln(w, "Bonjour, inconnu !")
			return
		}

		fmt.Fprintf(w, "Bonjour, %s !", name)
	})

	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		value := r.PathValue("id")
		id, err := strconv.Atoi(value)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(w, "Erreur, mauvais identifiant utilisateur, reçu: %s\n", value)

			return
		}

		for _, user := range users {
			if user.ID == id {
				fmt.Fprintf(w, "Le nom de l'utilisateur est %s\n", user.Name)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Utilisateur avec identifiant %d introuvable\n", id)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur l'API")
	})

	fmt.Println("Serveur démarré et prêt à écouter les requêtes HTTP")

	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur sur http://0.0.0.0:8000 (%s)", err.Error())
	}
}
