CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    confirmation_code VARCHAR(100),
    is_confirmed BOOLEAN DEFAULT FALSE,
    role VARCHAR(20) NOT NULL DEFAULT 'customer' CHECK (role IN ('customer', 'admin')),
    reset_code VARCHAR(100),
    reset_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Un token est créé à chaque login et envoyé dans le header "Authorization: Bearer <token>"
CREATE TABLE IF NOT EXISTS sessions (
    token VARCHAR(64) PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(20) UNIQUE NOT NULL, -- identifiant métier PDT-XXXXXX
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL DEFAULT 'divers',
    price DECIMAL(10, 2) NOT NULL,                  -- prix HT
    tax_rate DECIMAL(5, 2) NOT NULL DEFAULT 20.00,  -- TVA en %
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(20) UNIQUE NOT NULL, -- identifiant métier BSK-XXXXXX
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE, -- un seul panier par utilisateur
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_id)
);

CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipping', 'delivered', 'cancelled');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(20) UNIQUE NOT NULL, -- identifiant métier CMD-XXXXXX
    user_id INT REFERENCES users(id) ON DELETE RESTRICT,
    total DECIMAL(10, 2) NOT NULL,  -- total TTC
    status order_status DEFAULT 'pending',
    cancel_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_time DECIMAL(10, 2) NOT NULL -- prix unitaire TTC au moment de la commande
);

-- Quelques produits de démo pour tester la recherche et le panier
INSERT INTO products (reference, name, description, category, price, tax_rate, stock) VALUES
    ('PDT-7D2K8N', 'Clavier mécanique', 'Clavier mécanique AZERTY rétroéclairé', 'informatique', 79.99, 20.00, 25),
    ('PDT-3F9Q1A', 'Souris sans fil', 'Souris ergonomique Bluetooth', 'informatique', 29.90, 20.00, 50),
    ('PDT-8H4M2C', 'Écran 27 pouces', 'Écran IPS 27 pouces 1440p', 'informatique', 249.00, 20.00, 10),
    ('PDT-5K1P7T', 'Le Petit Prince', 'Roman de Saint-Exupéry, édition poche', 'livres', 6.90, 5.50, 100),
    ('PDT-2B6V9R', 'Café en grains 1kg', 'Café arabica torréfié', 'alimentation', 18.50, 5.50, 40);
