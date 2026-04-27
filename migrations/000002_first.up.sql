CREATE TABLE app.requests (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,

    city VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);