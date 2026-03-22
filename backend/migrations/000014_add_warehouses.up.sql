-- 000014_add_warehouses.up.sql

-- Create Warehouses
CREATE TABLE IF NOT EXISTS warehouses (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    address         TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_warehouses_org ON warehouses(organization_id);

-- Create Warehouse Items
CREATE TABLE IF NOT EXISTS warehouse_items (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    warehouse_id    BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity        INT NOT NULL DEFAULT 0,
    min_stock       INT NOT NULL DEFAULT 5,
    reorder_point   INT NOT NULL DEFAULT 10,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(warehouse_id, product_id)
);

CREATE INDEX idx_warehouse_items_org_product ON warehouse_items(organization_id, product_id);

-- Create Stock Transfers
CREATE TABLE IF NOT EXISTS stock_transfers (
    id                   BIGSERIAL PRIMARY KEY,
    organization_id      BIGINT NOT NULL REFERENCES organizations(id),
    reference_number     VARCHAR(100) NOT NULL,
    from_warehouse_id    BIGINT NOT NULL REFERENCES warehouses(id),
    to_warehouse_id      BIGINT NOT NULL REFERENCES warehouses(id),
    status               VARCHAR(50) DEFAULT 'pending',
    notes                TEXT,
    shipped_at           TIMESTAMPTZ,
    received_at          TIMESTAMPTZ,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ,
    UNIQUE(organization_id, reference_number)
);

-- Alter Stock Movements
ALTER TABLE stock_movements ADD COLUMN warehouse_id BIGINT REFERENCES warehouses(id);

-- DATA MIGRATION
-- 1. Create Main Warehouse for each organization
INSERT INTO warehouses (organization_id, code, name, is_active, created_at, updated_at)
SELECT id, 'WH-MAIN', 'Main Warehouse', true, NOW(), NOW() 
FROM organizations 
ON CONFLICT (organization_id, code) DO NOTHING;

-- 2. Migrate stock from products to warehouse_items
INSERT INTO warehouse_items (organization_id, warehouse_id, product_id, quantity, min_stock, reorder_point, created_at, updated_at)
SELECT p.organization_id, w.id, p.id, p.stock, p.min_stock, p.min_stock * 2, NOW(), NOW()
FROM products p
JOIN warehouses w ON w.organization_id = p.organization_id AND w.code = 'WH-MAIN'
ON CONFLICT (warehouse_id, product_id) DO NOTHING;

-- 3. Update existing stock movements to point to Main Warehouse
UPDATE stock_movements sm
SET warehouse_id = w.id
FROM warehouses w
WHERE sm.organization_id = w.organization_id AND w.code = 'WH-MAIN'
AND sm.warehouse_id IS NULL;
