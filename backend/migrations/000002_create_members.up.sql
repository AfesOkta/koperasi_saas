-- 000002_create_members.up.sql

CREATE TABLE IF NOT EXISTS members (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    user_id         BIGINT REFERENCES users(id) ON DELETE SET NULL,
    member_number   VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    nik             VARCHAR(50) NOT NULL,
    address         TEXT,
    phone           VARCHAR(20),
    status          VARCHAR(20) DEFAULT 'pending',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(organization_id, member_number),
    UNIQUE(organization_id, nik)
);

CREATE INDEX idx_members_organization_id ON members(organization_id);
CREATE INDEX idx_members_user_id ON members(user_id);
CREATE INDEX idx_members_status ON members(status);
CREATE INDEX idx_members_deleted_at ON members(deleted_at);

-- Member Documents (KYC)
CREATE TABLE IF NOT EXISTS member_documents (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    member_id       BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    file_url        VARCHAR(500) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_member_documents_member_id ON member_documents(member_id);

-- Member Cards
CREATE TABLE IF NOT EXISTS member_cards (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    member_id       BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    card_number     VARCHAR(100) NOT NULL UNIQUE,
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_member_cards_member_id ON member_cards(member_id);
CREATE INDEX idx_member_cards_card_number ON member_cards(card_number);
