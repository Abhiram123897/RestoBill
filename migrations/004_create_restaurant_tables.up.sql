CREATE TABLE restaurant_tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INT UNIQUE NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved')),
    created_at TIMESTAMP DEFAULT NOW()
);
