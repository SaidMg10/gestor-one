CREATE TABLE IF NOT EXISTS receipts (
    id BIGSERIAL PRIMARY KEY,
    income_id BIGINT NOT NULL UNIQUE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    uploaded_by BIGINT,
    checksum VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE receipts
ADD CONSTRAINT fk_receipts_income
    FOREIGN KEY (income_id)
    REFERENCES incomes(id)
    ON DELETE CASCADE;

ALTER TABLE receipts
ADD CONSTRAINT fk_receipts_uploaded_by
    FOREIGN KEY (uploaded_by)
    REFERENCES users(id);
