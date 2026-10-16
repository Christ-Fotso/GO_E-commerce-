# Manuel de Test et Guide d'Utilisation

Ce guide détaille pas à pas la procédure complète pour tester l'ensemble des fonctionnalités du projet e-commerce conformément au barème de notation.

---

## 1. Préparation de l'environnement

Ouvrez un premier terminal pour démarrer les services d'infrastructure :

### 1.1 Lancer la base de données PostgreSQL (Docker)
```bash
docker-compose up -d
```
> Vérifiez que le conteneur tourne avec `docker ps`. Vous devez voir `ecommerce_db` sur le port `5432`.

### 1.2 Démarrer le serveur HTTP Backend
```bash
go run ./cmd/server
```
Le serveur démarre sur `http://localhost:8080`.
Au premier lancement, il initialise automatiquement le compte administrateur par défaut :
- **Email** : `admin@shop.local`
- **Mot de passe** : `admin123`

---

## 2. Test de l'Interface Client (TUI Charmbracelet)

Ouvrez un **deuxième terminal** pour exécuter le client boutique.

```bash
go run ./cmd/client
```

### 2.1 Inscription (Authentification - 1 point)
1. Dans le menu d'accueil, sélectionnez **Inscription** (flèche bas puis `Entrée`).
2. Remplissez le formulaire :
   - **Pseudo** : `alice`
   - **Email** : `alice@example.com`
   - **Mot de passe** : `Secret123!`
3. Validez avec `Entrée`.
4. Regardez dans le terminal du serveur (`cmd/server`) : un message s'affiche avec le code de confirmation aléatoire simulé (ex: `680966`).
5. Notez ce code.

### 2.2 Confirmation du compte
1. Dans le menu, sélectionnez **Confirmation**.
2. Entrez l'adresse email (`alice@example.com`) et le code reçu.
3. Validez : le compte est désormais activé.

### 2.3 Connexion
1. Sélectionnez **Connexion**.
2. Entrez `alice@example.com` et `Secret123!`.
3. Validez : un message de bienvenue apparaît et le menu s'enrichit des options boutique (Catalogue, Panier, Commandes, Déconnexion).

### 2.4 Recherche et consultation de produits (2 points)
1. Sélectionnez **Catalogue**.
2. Vous observez les produits pré-remplis en base avec leurs informations complètes :
   - Nom, description, catégorie
   - Prix HT, taux de TVA et prix calculé TTC
   - Stock disponible
3. Appuyez sur `/` pour filtrer en temps réel (ex: tapez `clavier`, `livres`, ou `informatique`).
4. Sélectionnez un produit avec `Entrée` pour afficher sa fiche détaillée.

### 2.5 Ajout au panier et gestion des quantités (2 points)
1. Sur la fiche du produit sélectionné, appuyez sur la touche `a` pour l'ajouter au panier.
2. Revenez au menu (`Échap`) puis ouvrez **Panier**.
3. Vous constatez :
   - L'identifiant métier unique du panier (format `BSK-XXXXXX`).
   - Les articles présents avec quantité et sous-total.
   - Les frais de livraison indiqués comme gratuits.
   - Le montant total TTC calculé.
4. Utilisez les touches pour modifier la quantité (`+` / `-`) ou retirer un article (`d`).

### 2.6 Paiement sécurisé simulé (2 points)
1. Depuis le panier, appuyez sur `p` pour accéder au paiement.
2. Renseignez les coordonnées bancaires de test :
   - **Numéro de carte** : `4242 4242 4242 4242`
   - **Expiration** : `12/30`
   - **CVC** : `123`
3. Validez : le paiement est validé, le panier est vidé et une commande payée est générée.

### 2.7 Suivi des commandes (2 points)
1. Retournez au menu puis sélectionnez **Commandes**.
2. Vous retrouvez la commande créée avec :
   - Son identifiant métier unique (format `CMD-XXXXXX`).
   - La date et l'heure d'enregistrement.
   - Le montant total payé.
   - Le statut en vert : `[PAID]`.

---

## 3. Test de l'Interface Administrateur (TUI Charmbracelet)

Ouvrez un **troisième terminal** pour lancer la console d'administration.

```bash
go run ./cmd/admin
```
L'interface se connecte automatiquement avec les identifiants administrateur du backend.

### 3.1 Gestion des Utilisateurs (1 point)
1. Appuyez sur `1` ou `u` pour accéder à la **Gestion des Utilisateurs**.
2. Vous observez la liste des utilisateurs enregistrés (`admin`, `alice`, etc.) avec leur rôle et leur statut de confirmation.
3. **Créer un utilisateur** :
   - Appuyez sur `c`.
   - Remplissez le formulaire (pseudo, email, mot de passe, rôle `customer` ou `admin`, confirmation).
   - Validez : le nouvel utilisateur apparaît dans la liste.
4. **Modifier un utilisateur** :
   - Placez le curseur sur un utilisateur avec les flèches `haut`/`bas`.
   - Appuyez sur `e` pour modifier ses champs.
5. **Confirmer un compte** :
   - Sélectionnez un compte au statut `[EN ATTENTE]`.
   - Appuyez sur `v` pour le valider immédiatement.
6. **Supprimer un utilisateur** :
   - Appuyez sur `d` sur un compte sans commande pour le supprimer.
7. Revenez au menu principal avec `Échap`.

### 3.2 Gestion des Produits (1 point)
1. Appuyez sur `2` ou `p` pour ouvrir la **Gestion des Produits**.
2. Vous visualisez la liste complète des références produits.
3. **Ajouter un produit** :
   - Appuyez sur `c`.
   - Renseignez les informations demandées :
     - **Nom** : `Webcam 4K Ultra HD`
     - **Description** : `Webcam avec microphone réducteur de bruit`
     - **Catégorie** : `informatique`
     - **Prix HT** : `89.90`
     - **TVA (%)** : `20.0`
     - **Stock** : `15`
   - Validez : le produit est inséré et un identifiant métier unique (ex: `PDT-K8N2F4`) est généré automatiquement par le système.
4. Revenez au menu avec `Échap`.

### 3.3 Gestion des Commandes (1 point)
1. Appuyez sur `3` ou `o` pour ouvrir la **Gestion des Commandes**.
2. Vous observez l'ensemble des commandes passées sur la plateforme.
3. **Consulter les détails** :
   - Appuyez sur `Entrée` sur une commande pour afficher les lignes d'articles, les quantités et le total.
   - Revenez avec `Échap`.
4. **Changer le statut d'une commande** :
   - Sélectionnez une commande au statut `paid`.
   - Appuyez sur `s`.
   - Choisissez le nouveau statut (ex: `shipping` pour expédition, puis `delivered` pour livraison).
   - Si vous sélectionnez `cancelled`, le champ raison devient obligatoire.
   - Validez : le badge de statut se met à jour en temps réel.
5. **Créer une commande pour le compte d'un client** :
   - Appuyez sur `c`.
   - Entrez l'ID du client (ex: `2`), l'ID d'un produit (ex: `1`), la quantité et le statut initial (`pending` ou `paid`).
   - Validez : la commande est créée et rattachée au compte client.

---

## 4. Test de l'Interface Distante via SSH (2 points)

Ce composant permet à n'importe quel client ou administrateur de se connecter au service sans avoir installé l'application en local.

### 4.1 Démarrer le serveur SSH
Ouvrez un terminal dédié :
```bash
go run ./cmd/ssh
```
Le serveur SSH écoute sur le port `2222`.

### 4.2 Connexion avec le menu interactif
Depuis un autre terminal (ou une autre machine du réseau) :
```bash
ssh -p 2222 localhost
```
Un portail d'accueil interactif apparaît dans le terminal :
```text
PORTAIL ACCES DISTANT SSH - E-COMMERCE

Veuillez sélectionner l'interface que vous souhaitez exécuter :
  [1] Interface Boutique & Achats (Client TUI)
  [2] Interface Administration & Gestion (Admin TUI)
  [q] Quitter la session SSH
```
- Appuyez sur `1` : l'interface Boutique s'ouvre dans la session SSH.
- Appuyez sur `2` : l'interface d'Administration s'ouvre dans la session SSH.

### 4.3 Connexion directe par commande
Vous pouvez également cibler directement l'interface souhaitée :
```bash
# Accès direct à l'administration
ssh -p 2222 localhost admin

# Accès direct à la boutique client
ssh -p 2222 localhost client
```

---

## 5. Grille de conformité avec les consignes

| Fonctionnalité demandée | Barème | Statut | Emplacement du code |
| :--- | :---: | :---: | :--- |
| **Authentification** (register, code aléatoire, confirm, login, reset) | 1 pt | Validé | `internal/handlers/auth_handler.go`, `cmd/client` |
| **Recherche produit** (nom, prix, catégorie, description, TTC) | 2 pts | Validé | `internal/repositories/product_repo.go`, `internal/tui/client` |
| **Panier** (ajout, quantité, suppression, total TTC, ref `BSK-`) | 2 pts | Validé | `internal/repositories/cart_repo.go`, `cmd/client` |
| **Paiement** (CB simulée, commande paid, vidage panier) | 2 pts | Validé | `internal/handlers/payment.go`, `cmd/client` |
| **Commandes client** (historique, statut, date, ref `CMD-`) | 2 pts | Validé | `internal/repositories/order_repo.go`, `cmd/client` |
| **Gestion commandes admin** (changement statut, création client) | 1 pt | Validé | `internal/handlers/admin_handler.go`, `cmd/admin` |
| **Gestion utilisateurs admin** (créer, modifier, supprimer, confirmer) | 1 pt | Validé | `internal/handlers/admin_handler.go`, `cmd/admin` |
| **Gestion produits admin** (ajout produit, ref auto `PDT-`) | 1 pt | Validé | `internal/handlers/admin_handler.go`, `cmd/admin` |
| **Interfaces CLI séparées & Serveur HTTP `net/http`** | 2 pts | Validé | `cmd/client`, `cmd/admin`, `cmd/server` |
| **Persistance PostgreSQL & `database/sql` sans ORM** | 1 pt | Validé | `internal/database/database.go`, `docker-compose.yml` |
| **Charmbracelet** (`bubbletea`, `huh`, `lipgloss`) | 3 pts | Validé | `internal/tui/client`, `internal/tui/admin` |
| **Interface SSH** (Wish, port 2222, accès client & admin) | 2 pts | Validé | `cmd/ssh/main.go` |
| **Total** | **20 / 20** | | |
