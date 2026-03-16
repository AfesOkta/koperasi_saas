-- 000005_create_loans.up.sql

-- Loan Products
CREATE TABLE IF NOT EXISTS loan_products (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    interest_rate   DECIMAL(5,2) NOT NULL,
    interest_type   VARCHAR(50) NOT NULL, -- flat, declining
    max_amount      DECIMAL(15,2) NOT NULL,
    max_term        INT NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_loan_products_org_id ON loan_products(organization_id);

-- Loans Applications / Active Loans
CREATE TABLE IF NOT EXISTS loans (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   BIGINT NOT NULL REFERENCES organizations(id),
    member_id         BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    loan_product_id   BIGINT NOT NULL REFERENCES loan_products(id),
    loan_number       VARCHAR(100) NOT NULL UNIQUE,
    principal_amount  DECIMAL(15,2) NOT NULL,
    interest_rate     DECIMAL(5,2) NOT NULL,
    term_months       INT NOT NULL,
    total_interest    DECIMAL(15,2) DEFAULT 0,
    expected_total    DECIMAL(15,2) DEFAULT 0,
    outstanding       DECIMAL(15,2) DEFAULT 0,
    status            VARCHAR(20) DEFAULT 'pending', -- pending, approved, active, paid, defaulted
    approved_at       TIMESTAMPTZ,
    disbursed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    updated_at        TIMESTAMPTZ DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_loans_org_member ON loans(organization_id, member_id);

-- Loan Repayment Schedules
CREATE TABLE IF NOT EXISTS loan_schedules (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  BIGINT NOT NULL REFERENCES organizations(id),
    loan_id          BIGINT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    period           INT NOT NULL,
    due_date         DATE NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_amount  DECIMAL(15,2) NOT NULL,
    total_amount     DECIMAL(15,2) NOT NULL,
    paid_amount      DECIMAL(15,2) DEFAULT 0,
    status           VARCHAR(20) DEFAULT 'unpaid', -- unpaid, partial, paid
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_loan_schedules_loan_id ON loan_schedules(loan_id);

-- Loan Payments (Actual payments made)
CREATE TABLE IF NOT EXISTS loan_payments (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  BIGINT NOT NULL REFERENCES organizations(id),
    loan_id          BIGINT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    reference_number VARCHAR(100) NOT NULL UNIQUE,
    amount           DECIMAL(15,2) NOT NULL,
    payment_date     TIMESTAMPTZ NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_loan_payments_loan_id ON loan_payments(loan_id);
