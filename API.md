# API E-Commerce

Ne jamais renommer une URL de ce fichier.

- URL : `http://localhost:8080`
- Format : JSON
- Erreur : `{"error": "..."}`
- Routes connectées : header `Authorization: Bearer <token>` (token obtenu avec `POST /login`)
- Admin par défaut : `admin@shop.local` / `admin123`

## Objets

```json
User    {"id", "username", "email", "is_confirmed", "role", "created_at"}
Product {"id", "reference", "name", "description", "category", "price", "tax_rate", "price_ttc", "stock", "created_at"}
Cart    {"id", "reference", "user_id", "items", "total_ttc", "created_at"}
Order   {"id", "reference", "user_id", "total", "status", "cancel_reason", "items", "created_at"}
```

- `price` est HT, `price_ttc` et `total` sont TTC
- `status` : `pending`, `paid`, `shipping`, `delivered`, `cancelled`

## Public

### Authentification
- `POST /register` : `{"username", "email", "password"}` → `{"message", "user", "confirmation_code"}`
- `POST /confirm` : `{"email", "code"}`
- `POST /login` : `{"email", "password"}` → `{"token", "user"}`
- `POST /password/forgot` : `{"email"}` → `{"message", "reset_code"}`
- `POST /password/reset` : `{"email", "code", "new_password"}`

### Produits
- `GET /products` → `[Product]`
  - filtres : `q`, `name`, `description`, `category`, `min_price`, `max_price`, `min_price_ttc`, `max_price_ttc`
- `GET /products/{id}` → `Product` (id ou `PDT-XXXXXX`)

## Client connecté

### Compte
- `GET /me` → `User`
- `POST /logout`

### Panier
- `GET /cart` → `Cart`
- `POST /cart/items` : `{"product_id", "quantity"}` → `Cart`
- `PUT /cart/items/{product_id}` : `{"quantity"}` → `Cart` (0 = retirer)
- `DELETE /cart/items/{product_id}` → `Cart`
- `DELETE /cart` → `Cart` vide
- `POST /cart/pay` : `{"card_number", "expiry", "cvc"}` → `Order` payée
  - carte de test : `4242 4242 4242 4242`, `12/30`, `123`

### Commandes
- `GET /orders` → `[Order]`
- `GET /orders/{id}` → `Order` (id ou `CMD-XXXXXX`)

## Administrateur

### Utilisateurs
- `GET /admin/users` → `[User]`
- `POST /admin/users` : `{"username", "email", "password", "role", "is_confirmed"}` → `User`
- `GET /admin/users/{id}` → `User`
- `PUT /admin/users/{id}` : champs à modifier uniquement → `User`
- `DELETE /admin/users/{id}`
- `POST /admin/users/{id}/confirm` → `User`

### Produits
- `GET /admin/products` → `[Product]`
- `POST /admin/products` : `{"name", "description", "category", "price", "tax_rate", "stock"}` → `Product`

### Commandes
- `GET /admin/orders` → `[Order]` (filtres : `status`, `user_id`)
- `POST /admin/orders` : `{"user_id", "items": [{"product_id", "quantity"}]}` → `Order`
- `GET /admin/orders/{id}` → `Order`
- `PUT /admin/orders/{id}/status` : `{"status", "reason"}` → `Order` (`reason` obligatoire si `cancelled`)

Transitions autorisées :

```
pending → paid → shipping → delivered
pending, paid, shipping → cancelled
```
