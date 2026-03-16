-- 000008_create_sales.up.sql

-- Orders Header
CREATE TABLE IF NOT EXISTS orders (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    member_id       BIGINT REFERENCES members(id),
    order_id        VARCHAR(100) NOT NULL UNIQUE, -- Human readable reference
    total_amount    DECIMAL(15,2) NOT NULL,
    discount        DECIMAL(15,2) DEFAULT 0,
    tax_amount      DECIMAL(15,2) DEFAULT 0,
    final_amount    DECIMAL(15,2) NOT NULL,
    payment_status  VARCHAR(20) DEFAULT 'unpaid', -- unpaid, paid, partial
    status          VARCHAR(20) DEFAULT 'completed',
    cashier_id      BIGINT NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_orders_org_id ON orders(organization_id);
CREATE INDEX idx_orders_member ON orders(member_id);

-- Order Items
CREATE TABLE IF NOT EXISTS order_items (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    order_id        BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL REFERENCES products(id),
    quantity        INT NOT NULL,
    unit_price      DECIMAL(15,2) NOT NULL,
    subtotal        DECIMAL(15,2) NOT NULL,
    discount        DECIMAL(15,2) DEFAULT 0,
    total_amount    DECIMAL(15,2) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_order_items_order ON order_items(order_id);

-- Order Payments
CREATE TABLE IF NOT EXISTS order_payments (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    order_id        BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    payment_method  VARCHAR(50) NOT NULL, -- cash, savings, transfer
    amount          DECIMAL(15,2) NOT NULL,
    reference_token VARCHAR(255),
    status          VARCHAR(20) DEFAULT 'completed',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_order_payments_order ON order_payments(order_id);
