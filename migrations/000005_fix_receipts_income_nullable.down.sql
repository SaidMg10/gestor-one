ALTER TABLE receipts
ALTER COLUMN income_id SET NOT NULL;

ALTER TABLE receipts
ADD CONSTRAINT receipts_income_id_key UNIQUE (income_id);
