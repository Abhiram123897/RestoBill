CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID REFERENCES restaurant_tables(id),
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'served', 'billed', 'paid')),
    order_type VARCHAR(20) NOT NULL CHECK (order_type IN ('dine_in', 'takeaway', 'delivery')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
