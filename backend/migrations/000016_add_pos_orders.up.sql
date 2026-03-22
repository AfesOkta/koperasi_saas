-- 000016_add_pos_orders.up.sql

CREATE TABLE pos_orders (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    shift_id BIGINT NOT NULL REFERENCES pos_shifts(id),
    reference_number VARCHAR(100) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    final_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(50) NOT NULL, -- cash, transfer, qris
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, completed, cancelled
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(organization_id, reference_number)
);

CREATE TABLE pos_order_items (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES pos_orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL DEFAULT 0,
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,
    kds_status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, preparing, ready, served
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_pos_orders_org_shift ON pos_orders(organization_id, shift_id);
CREATE INDEX idx_pos_order_items_order ON pos_order_items(order_id);
