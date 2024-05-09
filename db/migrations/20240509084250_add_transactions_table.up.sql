CREATE TABLE IF NOT EXISTS transactions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    paid int NOT NULL,
    change int NOT NULL,
    craeted_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (customer_id) REFERENCES customers(id)
    ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_trx_customer_id ON transactions(customer_id);