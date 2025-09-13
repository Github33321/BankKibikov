CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
    otp_code TEXT,
    otp_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
    );
