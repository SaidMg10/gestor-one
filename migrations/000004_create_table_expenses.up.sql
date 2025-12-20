CREATE TABLE IF NOT EXISTS expenses (
    id BIGSERIAL PRIMARY KEY,
    amount NUMERIC(12,2) NOT NULL,
    description VARCHAR(255),
    date TIMESTAMP NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

ALTER TABLE expenses
ADD CONSTRAINT fk_expenses_created_by
    FOREIGN KEY (created_by)
    REFERENCES users(id);

ALTER TABLE receipts
ADD COLUMN expense_id BIGINT NULL;

ALTER TABLE receipts
ADD CONSTRAINT fk_receipts_expense
    FOREIGN KEY (expense_id)
    REFERENCES expenses(id)
    ON DELETE CASCADE;

CREATE INDEX idx_receipts_expense_id
ON receipts(expense_id);

ALTER TABLE receipts
ADD CONSTRAINT chk_receipt_owner
CHECK (
    (income_id IS NOT NULL AND expense_id IS NULL)
 OR (income_id IS NULL AND expense_id IS NOT NULL)
);

