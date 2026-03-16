-- 000006_create_inventory.up.sql

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, name)
);

CREATE INDEX idx_categories_org_id ON categories(organization_id);

-- Products
CREATE TABLE IF NOT EXISTS products (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    category_id     BIGINT REFERENCES categories(id),
    sku             VARCHAR(100) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price           DECIMAL(15,2) NOT NULL,
    cost_price      DECIMAL(15,2) NOT NULL,
    stock           INT NOT NULL DEFAULT 0,
    min_stock       INT NOT NULL DEFAULT 0,
    unit            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, sku)
);

CREATE INDEX idx_products_org_id ON products(organization_id);
CREATE INDEX idx_products_category ON products(category_id);

-- Stock Movements
CREATE TABLE IF NOT EXISTS stock_movements (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    product_id        BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    reference_number  VARCHAR(100) NOT NULL,
    type              VARCHAR(50) NOT NULL, -- in, out, adj
    quantity          INT NOT NULL,
    balance_after     INT NOT NULL,
    notes             TEXT,
    related_entity    VARCHAR(100),
    related_entity_id BIGINT,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_stock_movements_org_product ON stock_movements(organization_id, product_id);
CREATE INDEX idx_stock_movements_ref ON stock_movements(reference_number);
