CREATE TABLE withdrawals (
                        id SERIAL PRIMARY KEY,
                        sum INT,
                        order_number VARCHAR(50) UNIQUE NOT NULL,
                        user_id INT REFERENCES users(id) ON DELETE SET NULL,
                        processed_at TIMESTAMP DEFAULT NOW()
);
