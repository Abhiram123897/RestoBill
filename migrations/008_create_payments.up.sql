CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id UUID NOT NULL REFERENCES bills(id),
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    method VARCHAR(20) NOT NULL CHECK (method IN ('cash', 'card', 'upi')),
    reference TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
