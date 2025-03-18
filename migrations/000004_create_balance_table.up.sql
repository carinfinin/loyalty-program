CREATE TABLE IF NOT EXISTS balance (
                        id SERIAL PRIMARY KEY,
                        current NUMERIC(20, 2) DEFAULT 0.00,
                        withdrawn NUMERIC(20, 2) DEFAULT 0.00,
                        user_id INT REFERENCES users(id) ON DELETE SET NULL UNIQUE,
                        created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_balance_user_id ON balance (user_id);