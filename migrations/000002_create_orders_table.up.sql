CREATE TABLE orders (
                       id SERIAL PRIMARY KEY,
                       number BIGINT UNIQUE NOT NULL,
                       status VARCHAR(50) NOT NULL,
                       accrual NUMERIC(20, 2),
                       user_id INT REFERENCES users(id) ON DELETE SET NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_number ON orders (number);