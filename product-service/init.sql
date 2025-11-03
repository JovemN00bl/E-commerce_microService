CREATE TABLE IF NOT EXISTS Products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL ,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK(price >= 0),
    stock_quantity INT NOT NULL DEFAULT 0 CHECK(stock_quantity >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')

);

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);