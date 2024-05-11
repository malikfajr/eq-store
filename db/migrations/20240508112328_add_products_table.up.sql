CREATE TABLE IF NOT EXISTS products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) NOT NULL,
    sku VARCHAR(30) NOT NULL,
    category VARCHAR(11) NOT NULL,
    image_url TEXT NOT NULL,
    notes VARCHAR(200) NOT NULL,
    price INT NOT NULL CHECK(price >= 1),
    stock INT NOT NULL CHECK(stock >= 0 AND stock <= 100000),
    location VARCHAR(200) NOT NULL,
    is_available BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

CREATE INDEX idx_product_name ON products USING GIN(name);
CREATE INDEX idx_product_sku ON products (sku);
CREATE INDEX idx_product_category ON products (category);