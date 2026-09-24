# Document de Suivi et Explications du Code

Ce document a pour but d'expliquer les choix techniques et le fonctionnement de chaque partie du projet, pour faciliter la compréhension et la préparation à la soutenance.

## 1. Architecture du Projet
Le projet suit une architecture claire et standard en Go :
- `cmd/` : Contient les points d'entrée des applications (le serveur HTTP, le client CLI, le CLI admin).
- `internal/` : Contient le code métier privé (modèles, base de données, requêtes, etc.). C'est une convention Go pour empêcher l'importation de ces packages depuis l'extérieur.
- `migrations/` : Contient les scripts SQL (`init.sql`) pour initialiser la structure de la base de données.

## 2. Base de Données (PostgreSQL)
**Contrainte respectée :** Utilisation exclusive de `database/sql` sans ORM (comme GORM). 
- Le driver utilisé est `github.com/lib/pq`.
- **`internal/database/database.go`** : Gère la connexion pure à Postgres via une fonction `InitDB`. Elle utilise un DSN (Data Source Name) et fait un `Ping()`.

## 3. Le Pattern Repository (Dépôt)
Pour éviter de mélanger le code SQL avec la logique de l'API, j'ai mis en place le "Pattern Repository".
- **`internal/repositories/user_repo.go`** : Contient une structure `UserRepository` qui possède un pointeur vers la base de données (`*sql.DB`).
- Les méthodes comme `Create` ou `GetByEmail` s'occupent uniquement de faire les requêtes SQL (ex: `INSERT INTO users...`) et de mapper les résultats dans nos modèles Go (avec `.Scan()`).

## 4. L'Étape 2 : API d'Authentification (sans framework)
**Contrainte respectée :** J'ai utilisé le routage natif de Go (`net/http`) comme montré dans le Module 10 (`http.HandleFunc("POST /register", ...)`).
- **`internal/handlers/auth_handler.go`** : 
  - La fonction **`Register`** décode le JSON reçu (avec `json.NewDecoder`). Elle utilise ensuite le package `golang.org/x/crypto/bcrypt` pour sécuriser (hasher) le mot de passe avant de l'envoyer au `UserRepository`. Elle donne ensuite un `ConfirmationCode`.
  - La fonction **`Confirm`** permet de valider le compte (le champ `is_confirmed` passe à `TRUE` en BDD).
  - La fonction **`Login`** vérifie si l'utilisateur existe, s'il est bien confirmé, puis utilise `bcrypt.CompareHashAndPassword` pour valider le mot de passe reçu avec celui de la base de données. S'il correspond, on est connecté !
- L'utilisation des *structs* anonymes dans chaque fonction permet de définir la structure JSON attendue très simplement (ex: `var req struct { Email string... }`).
