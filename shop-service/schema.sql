CREATE TABLE shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE shop_warehouses (
    shop_id UUID REFERENCES shops(id),
    warehouse_id UUID REFERENCES warehouses(id),
    PRIMARY KEY (shop_id, warehouse_id)
);