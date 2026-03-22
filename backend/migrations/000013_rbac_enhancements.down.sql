-- 000013_rbac_enhancements.down.sql
ALTER TABLE roles DROP COLUMN IF EXISTS version;
ALTER TABLE permissions DROP COLUMN IF EXISTS scope;
DROP INDEX IF EXISTS idx_permissions_resource_action;
