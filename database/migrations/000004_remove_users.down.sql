ALTER TABLE schedules ADD COLUMN user_id INT;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
