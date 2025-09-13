CREATE TABLE loans (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       amount NUMERIC(15,2) NOT NULL,
                       issued_at TIMESTAMP NOT NULL DEFAULT NOW()
);
