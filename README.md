# E-Commerce CLI

Projet de création d'un outil CLI et d'un serveur e-commerce en Go.

## Structure du projet

- `cmd/server/` : Le serveur HTTP API.
- `cmd/client/` : Le client TUI pour les clients (boutique, panier, commandes).
- `cmd/admin/` : Le CLI TUI d'administration (utilisateurs, produits, commandes).
- `cmd/ssh/` : Le serveur SSH (Wish) permettant l'accès distant aux TUIs.
- `internal/` : Logique métier (handlers, repositories, modèles, TUI).
- `migrations/` : Scripts SQL d'initialisation de la base de données.

## Prérequis
- Go 1.22+
- Docker & Docker Compose

## Démarrage rapide

### 1. Base de données
```bash
docker-compose up -d
```

### 2. Serveur Backend (API)
```bash
go run ./cmd/server
```
Le serveur démarre sur `http://localhost:8080`. Un compte administrateur initial est créé par défaut (`admin@shop.local` / `admin123`).

### 3. Interface Client (Boutique & Achats)
Dans un nouveau terminal :
```bash
go run ./cmd/client
```

### 4. Interface Administrateur (Gestion locale)
Dans un nouveau terminal :
```bash
go run ./cmd/admin
```

### 5. Accès distant via SSH
1. Démarrer le serveur SSH :
```bash
go run ./cmd/ssh
```
2. Se connecter depuis un terminal distant :
```bash
ssh -p 2222 localhost
```
Un menu interactif permet de choisir entre l'interface Boutique et l'interface Administration. On peut aussi lancer directement une interface spécifique :
```bash
ssh -p 2222 localhost admin
ssh -p 2222 localhost client
```

## Guide de test et évaluation

Pour dérouler le protocole complet de test pas-à-pas couvrant l'ensemble des critères du sujet (authentification, recherche produit, panier, paiement avec algorithme de Luhn, commandes, administration et accès SSH), consultez le fichier [`MANUEL_TESTS.md`](MANUEL_TESTS.md).

