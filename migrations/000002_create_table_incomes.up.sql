CREATE TABLE IF NOT EXISTS incomes (
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

ALTER TABLE incomes
ADD CONSTRAINT fk_incomes_created_by
    FOREIGN KEY (created_by)
    REFERENCES users(id);
