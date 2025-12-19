ALTER TABLE receipts
ALTER COLUMN income_id DROP NOT NULL;

ALTER TABLE receipts
DROP CONSTRAINT IF EXISTS receipts_income_id_key;
