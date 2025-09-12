CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id SERIAL REFERENCES orders(id),
    product_id SERIAL NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE stock_reservations (
    id SERIAL PRIMARY KEY,
    order_id SERIAL REFERENCES orders(id),
    product_id SERIAL NOT NULL,
    warehouse_id SERIAL REFERENCES warehouses(id),
    quantity INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);