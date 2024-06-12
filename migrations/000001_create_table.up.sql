CREATE TABLE If not EXISTS payment(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id uuid ,
    amount DECIMAL,
    payment_method VARCHAR(100),
    payment_status VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at bigint DEFAULT 0
);
