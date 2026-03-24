-- 000017_add_loan_collaterals_and_approval_logs.up.sql

-- Loan Collaterals (for Professional/Enterprise plans)
CREATE TABLE IF NOT EXISTS loan_collaterals (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    loan_id         BIGINT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    description     TEXT,
    document_url    VARCHAR(500),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_loan_collaterals_loan_id ON loan_collaterals(loan_id);

-- Approval Logs (multi-level approval tracking)
CREATE TABLE IF NOT EXISTS approval_logs (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    loan_id         BIGINT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    approver_id     BIGINT NOT NULL,
    role            VARCHAR(20) NOT NULL,
    action          VARCHAR(20) NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_approval_logs_loan_id ON approval_logs(loan_id);

-- Add purpose field to loans table
ALTER TABLE loans ADD COLUMN IF NOT EXISTS purpose TEXT;

-- Add disbursement_method to loans table
ALTER TABLE loans ADD COLUMN IF NOT EXISTS disbursement_method VARCHAR(20) DEFAULT 'transfer';

-- Add rejected columns if not present
ALTER TABLE loans ADD COLUMN IF NOT EXISTS approved_by BIGINT;
ALTER TABLE loans ADD COLUMN IF NOT EXISTS rejected_at TIMESTAMPTZ;
ALTER TABLE loans ADD COLUMN IF NOT EXISTS rejected_by BIGINT;
