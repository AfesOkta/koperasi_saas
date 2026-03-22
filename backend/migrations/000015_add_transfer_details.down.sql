-- 000015_add_transfer_details.down.sql

ALTER TABLE stock_transfers DROP COLUMN IF EXISTS product_id;
ALTER TABLE stock_transfers DROP COLUMN IF EXISTS quantity;
