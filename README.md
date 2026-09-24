# E-Commerce CLI

Projet de création d'un outil CLI et d'un serveur e-commerce en Go.

## Structure du projet

- `cmd/server/` : Le serveur HTTP.
- `cmd/client/` : Le client CLI.
- `cmd/admin/` : Le CLI d'administration.
- `internal/models/` : Les structures de données.
- `migrations/` : Scripts SQL d'initialisation de la base de données.

## Prérequis
- Go installé
- Docker et Docker Compose

## Démarrage rapide

1. Lancer la base de données :
```bash
docker-compose up -d
```

2. Lancer le serveur :
```bash
go run cmd/server/main.go
```
