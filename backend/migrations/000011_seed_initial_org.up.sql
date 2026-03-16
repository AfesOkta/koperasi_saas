-- 000011_seed_initial_org.up.sql

INSERT INTO organizations (id, name, slug, status, settings) 
VALUES (1, 'Platform Owner', 'platform-owner', 'active', '{"enabled_modules": ["iam", "organization", "member", "savings", "loans"]}')
ON CONFLICT (id) DO NOTHING;

-- Seed some default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
('org:create', 'organization', 'create', 'Create new organization'),
('org:update', 'organization', 'update', 'Update organization'),
('org:read', 'organization', 'read', 'Read organization details'),
('org:settings', 'organization', 'settings', 'Manage organization settings')
ON CONFLICT (name) DO NOTHING;
