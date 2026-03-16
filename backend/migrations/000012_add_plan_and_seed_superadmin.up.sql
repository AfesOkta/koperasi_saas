-- 000012_add_plan_and_seed_superadmin.up.sql

ALTER TABLE organizations ADD COLUMN IF NOT EXISTS plan VARCHAR(50) DEFAULT 'basic';

-- Seed SuperAdmin role for Org ID 1
INSERT INTO roles (organization_id, name, description, is_system)
VALUES (1, 'SuperAdmin', 'Platform Global Administrator', true)
ON CONFLICT (organization_id, name) DO NOTHING;

-- Seed SuperAdmin User for Org ID 1 (password: password123)
INSERT INTO users (organization_id, name, email, password_hash, status)
VALUES (1, 'Super Admin', 'superadmin@koperasi.id', '$2a$10$46skW2R2smbjN2Ox5HFqNO8HhTnXzHjBC4NZeX4IWz622bnMZ5D3W', 'active')
ON CONFLICT (organization_id, email) DO NOTHING;

-- Link SuperAdmin role to SuperAdmin user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id 
FROM users u, roles r 
WHERE u.email = 'superadmin@koperasi.id' 
  AND r.name = 'SuperAdmin' 
  AND r.organization_id = 1
ON CONFLICT DO NOTHING;

-- Grant all permissions to SuperAdmin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'SuperAdmin' AND r.organization_id = 1
ON CONFLICT DO NOTHING;
