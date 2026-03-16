-- 000012_add_plan_and_seed_superadmin.down.sql
DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE email = 'superadmin@koperasi.id');
DELETE FROM users WHERE email = 'superadmin@koperasi.id';
DELETE FROM roles WHERE name = 'SuperAdmin' AND organization_id = 1;
ALTER TABLE organizations DROP COLUMN IF EXISTS plan;
