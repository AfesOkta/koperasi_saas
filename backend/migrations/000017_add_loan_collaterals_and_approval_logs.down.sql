-- 000017_add_loan_collaterals_and_approval_logs.down.sql

ALTER TABLE loans DROP COLUMN IF EXISTS rejected_by;
ALTER TABLE loans DROP COLUMN IF EXISTS rejected_at;
ALTER TABLE loans DROP COLUMN IF EXISTS approved_by;
ALTER TABLE loans DROP COLUMN IF EXISTS disbursement_method;
ALTER TABLE loans DROP COLUMN IF EXISTS purpose;

DROP TABLE IF EXISTS approval_logs;
DROP TABLE IF EXISTS loan_collaterals;
