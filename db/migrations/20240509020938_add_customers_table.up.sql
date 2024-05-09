CREATE TABLE IF NOT EXISTS customers(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_customer_phone_number ON customers(phone_number);
