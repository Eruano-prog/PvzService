CREATE TABLE IF NOT EXISTS pvzs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT NOT NULL,
    registration_time TIMESTAMP NOT NULL
)