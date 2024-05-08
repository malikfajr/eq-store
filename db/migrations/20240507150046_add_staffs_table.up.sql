CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    password CHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staff_phone_number ON staffs(phone_number);