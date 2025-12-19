ALTER TABLE receipts
DROP CONSTRAINT IF EXISTS fk_receipts_expense;

ALTER TABLE receipts
DROP CONSTRAINT IF EXISTS chk_receipt_owner;

ALTER TABLE receipts
DROP COLUMN IF EXISTS expense_id;

DROP INDEX IF EXISTS idx_receipts_expense_id;

DROP TABLE IF EXISTS expenses;

