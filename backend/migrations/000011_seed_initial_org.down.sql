-- 000011_seed_initial_org.down.sql
DELETE FROM permissions WHERE resource = 'organization';
DELETE FROM organizations WHERE id = 1;
