CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID NOT NULL REFERENCES receptions(id),
    type TEXT NOT NULL,
    datetime TIMESTAMP NOT NULL
)