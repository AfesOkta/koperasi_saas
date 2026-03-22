-- 000013_rbac_enhancements.up.sql
-- Add version column to roles for cache busting
ALTER TABLE roles ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1;

-- Add scope column to permissions for resource:action:scope pattern
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS scope VARCHAR(10) NOT NULL DEFAULT 'any';

-- Create index for permission lookup
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action, scope);
