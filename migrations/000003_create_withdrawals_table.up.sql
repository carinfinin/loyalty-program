CREATE TABLE withdrawals (
                        id SERIAL PRIMARY KEY,
                        sum NUMERIC(20, 2) NOT NULL,
                        order_number VARCHAR(50) UNIQUE NOT NULL,
                        user_id INT REFERENCES users(id) ON DELETE SET NULL,
                        processed_at TIMESTAMPTZ DEFAULT NOW()
);
