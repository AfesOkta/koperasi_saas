-- 000009_create_purchases.up.sql

-- Purchase Orders Header
CREATE TABLE IF NOT EXISTS purchase_orders (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    supplier_id     BIGINT NOT NULL REFERENCES suppliers(id),
    po_number       VARCHAR(100) NOT NULL UNIQUE,
    total_amount    DECIMAL(15,2) NOT NULL,
    discount        DECIMAL(15,2) DEFAULT 0,
    tax_amount      DECIMAL(15,2) DEFAULT 0,
    final_amount    DECIMAL(15,2) NOT NULL,
    payment_status  VARCHAR(20) DEFAULT 'unpaid', -- unpaid, paid, partial
    status          VARCHAR(20) DEFAULT 'pending', -- pending, ordered, received, cancelled
    notes           TEXT,
    received_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_purchase_orders_org_id ON purchase_orders(organization_id);
CREATE INDEX idx_purchase_orders_supplier ON purchase_orders(supplier_id);

-- Purchase Order Items
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    purchase_order_id BIGINT NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id        BIGINT NOT NULL REFERENCES products(id),
    quantity          INT NOT NULL,
    cost_price        DECIMAL(15,2) NOT NULL,
    subtotal          DECIMAL(15,2) NOT NULL,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_purchase_order_items_po ON purchase_order_items(purchase_order_id);

-- Purchase Payments
CREATE TABLE IF NOT EXISTS purchase_payments (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    purchase_order_id BIGINT NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    payment_method    VARCHAR(50) NOT NULL, -- cash, transfer
    amount            DECIMAL(15,2) NOT NULL,
    payment_date      TIMESTAMPTZ NOT NULL,
    reference_token   VARCHAR(255),
    status            VARCHAR(20) DEFAULT 'completed',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_purchase_payments_po ON purchase_payments(purchase_order_id);
