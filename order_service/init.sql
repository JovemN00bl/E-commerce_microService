
CREATE TYPE order_status as ENUM (
    'PENDING',
    'PAID',
    'SHIPPED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,
    total_price NUMERIC(10,2) NOT NULL CHECK(total_price >= 0),
    order order_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS roder_items (
    id UUID PRIMARY KEY DEFFAULT gen_random_uuid(),

    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity INT NOT NULL CHECK(quantity > 0),
    price_at_time NUMERIC(10,2 ) NOT NULL CHECK(price_at_time >= 0)

);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)