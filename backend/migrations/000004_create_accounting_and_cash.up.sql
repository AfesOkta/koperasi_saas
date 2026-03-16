-- 000004_create_accounting_and_cash.up.sql

-- Chart of Accounts
CREATE TABLE IF NOT EXISTS accounts (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL, -- Asset, Liability, Equity, Revenue, Expense
    normal_balance  VARCHAR(20) NOT NULL, -- debit, credit
    parent_id       BIGINT REFERENCES accounts(id),
    is_active       BOOLEAN DEFAULT true,
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_accounts_org_id ON accounts(organization_id);

-- Journal Entries
CREATE TABLE IF NOT EXISTS journal_entries (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  BIGINT NOT NULL REFERENCES organizations(id),
    reference_number VARCHAR(100) NOT NULL UNIQUE,
    date             DATE NOT NULL,
    description      TEXT,
    status           VARCHAR(20) DEFAULT 'posted',
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_journal_entries_org_date ON journal_entries(organization_id, date);

-- Journal Entry Lines
CREATE TABLE IF NOT EXISTS journal_entry_lines (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  BIGINT NOT NULL REFERENCES organizations(id),
    journal_entry_id BIGINT NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id       BIGINT NOT NULL REFERENCES accounts(id),
    description      TEXT,
    debit            DECIMAL(15,2) DEFAULT 0,
    credit           DECIMAL(15,2) DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_journal_entry_lines_je_id ON journal_entry_lines(journal_entry_id);
CREATE INDEX idx_journal_entry_lines_account_id ON journal_entry_lines(account_id);

-- Cash Registers
CREATE TABLE IF NOT EXISTS cash_registers (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL, -- cash, bank, e-wallet
    balance         DECIMAL(15,2) DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'active',
    account_id      BIGINT REFERENCES accounts(id),
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_cash_registers_org_id ON cash_registers(organization_id);

-- Cash Transactions
CREATE TABLE IF NOT EXISTS cash_transactions (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    cash_register_id  BIGINT NOT NULL REFERENCES cash_registers(id) ON DELETE CASCADE,
    reference_number  VARCHAR(100) NOT NULL UNIQUE,
    type              VARCHAR(50) NOT NULL, -- in, out, transfer
    amount            DECIMAL(15,2) NOT NULL,
    balance_after     DECIMAL(15,2) NOT NULL,
    category          VARCHAR(100) NOT NULL,
    description       TEXT,
    related_entity    VARCHAR(100),
    related_entity_id BIGINT,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_cash_txns_org_register ON cash_transactions(organization_id, cash_register_id);
