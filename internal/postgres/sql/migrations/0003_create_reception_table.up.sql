CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL REFERENCES pvzs(id),
    datetime TIMESTAMP NOT NULL,
    status TEXT NOT NULL
)