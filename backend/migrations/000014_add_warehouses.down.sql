-- 000014_add_warehouses.down.sql

ALTER TABLE stock_movements DROP COLUMN IF EXISTS warehouse_id;
DROP TABLE IF EXISTS stock_transfers;
DROP TABLE IF EXISTS warehouse_items;
DROP TABLE IF EXISTS warehouses;
