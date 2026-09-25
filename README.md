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
go run ./cmd/server
```

3. Tester (dans **un autre** terminal PowerShell, le serveur doit rester ouvert) :
```powershell
curl.exe http://localhost:8080/health
```

Sur PowerShell, `curl` tout court est un alias de `Invoke-WebRequest` : utiliser `curl.exe`.

Si le port 8080 est pris :
```powershell
$env:PORT=8081; go run ./cmd/server
curl.exe http://localhost:8081/health
```

Les URLs de l'équipe sont figées dans [API.md](API.md) et `internal/api/routes.go`.
