-- 000015_add_transfer_details.up.sql

ALTER TABLE stock_transfers ADD COLUMN product_id BIGINT NOT NULL REFERENCES products(id);
ALTER TABLE stock_transfers ADD COLUMN quantity INT NOT NULL DEFAULT 0;
