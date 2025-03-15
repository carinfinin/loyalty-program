CREATE TABLE balance (
                        id SERIAL PRIMARY KEY,
                        current NUMERIC(20, 2) NOT NULL,
                        withdrawn NUMERIC(20, 2) NOT NULL,
                        user_id INT REFERENCES users(id) ON DELETE SET NULL UNIQUE,
                        created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_balance_user_id ON balance (user_id);