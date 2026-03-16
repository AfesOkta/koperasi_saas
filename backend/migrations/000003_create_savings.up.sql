-- 000003_create_savings.up.sql

-- Saving Products
CREATE TABLE IF NOT EXISTS saving_products (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    status          VARCHAR(20) DEFAULT 'active',
    is_withdrawble  BOOLEAN DEFAULT false,
    interest_rate   DECIMAL(5,2) DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_saving_products_org_id ON saving_products(organization_id);

-- Saving Accounts
CREATE TABLE IF NOT EXISTS saving_accounts (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    member_id         BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    saving_product_id BIGINT NOT NULL REFERENCES saving_products(id),
    account_number    VARCHAR(100) NOT NULL UNIQUE,
    balance           DECIMAL(15,2) DEFAULT 0,
    status            VARCHAR(20) DEFAULT 'active',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    UNIQUE(organization_id, member_id, saving_product_id)
);

CREATE INDEX idx_saving_accounts_org_member ON saving_accounts(organization_id, member_id);

-- Saving Transactions
CREATE TABLE IF NOT EXISTS saving_transactions (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    saving_account_id BIGINT NOT NULL REFERENCES saving_accounts(id) ON DELETE CASCADE,
    reference_number  VARCHAR(100) NOT NULL UNIQUE,
    type              VARCHAR(50) NOT NULL, -- deposit, withdrawal
    amount            DECIMAL(15,2) NOT NULL,
    balance_after     DECIMAL(15,2) NOT NULL,
    description       TEXT,
    status            VARCHAR(20) DEFAULT 'completed',
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_saving_txns_org_account ON saving_transactions(organization_id, saving_account_id);
CREATE INDEX idx_saving_txns_type ON saving_transactions(type);
