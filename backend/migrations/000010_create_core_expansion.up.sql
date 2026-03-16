-- 000010_create_core_expansion.up.sql

-- Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    user_id         BIGINT REFERENCES users(id),
    action          VARCHAR(50) NOT NULL,
    resource        VARCHAR(100) NOT NULL,
    resource_id     VARCHAR(100),
    old_values      TEXT,
    new_values      TEXT,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_audit_logs_org ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource, resource_id);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    title           VARCHAR(255) NOT NULL,
    message         TEXT NOT NULL,
    type            VARCHAR(50) DEFAULT 'info',
    is_read         BOOLEAN DEFAULT FALSE,
    link            VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(is_read);

-- Billing & Subscriptions
CREATE TABLE IF NOT EXISTS subscription_plans (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    code            VARCHAR(50) NOT NULL UNIQUE,
    description     TEXT,
    price           DECIMAL(15,2) NOT NULL,
    max_users       INT DEFAULT 0,
    max_members     INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS org_subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL UNIQUE REFERENCES organizations(id),
    plan_id         BIGINT NOT NULL REFERENCES subscription_plans(id),
    start_date      TIMESTAMPTZ NOT NULL,
    end_date        TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    renewal_date    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- POS Shifts
CREATE TABLE IF NOT EXISTS pos_shifts (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    cashier_id      BIGINT NOT NULL REFERENCES users(id),
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ,
    start_balance   DECIMAL(15,2) NOT NULL,
    end_balance     DECIMAL(15,2),
    actual_cash     DECIMAL(15,2),
    difference      DECIMAL(15,2),
    notes           TEXT,
    status          VARCHAR(20) DEFAULT 'open',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- SHU (Profit Distribution)
CREATE TABLE IF NOT EXISTS shu_configs (
    id                  BIGSERIAL PRIMARY KEY,
    organization_id     BIGINT NOT NULL REFERENCES organizations(id),
    year                INT NOT NULL,
    total_shu           DECIMAL(15,2) NOT NULL,
    member_savings_pct  DECIMAL(5,2) NOT NULL,
    member_business_pct DECIMAL(5,2) NOT NULL,
    status              VARCHAR(20) DEFAULT 'draft',
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS shu_distributions (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    shu_config_id   BIGINT NOT NULL REFERENCES shu_configs(id) ON DELETE CASCADE,
    member_id       BIGINT NOT NULL REFERENCES members(id),
    savings_share   DECIMAL(15,2) NOT NULL,
    business_share  DECIMAL(15,2) NOT NULL,
    total_amount    DECIMAL(15,2) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_shu_dist_config ON shu_distributions(shu_config_id);
CREATE INDEX idx_shu_dist_member ON shu_distributions(member_id);
