CREATE TABLE IF NOT EXISTS productions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID NOT NULL REFERENCES receptions(id),
    type TEXT NOT NULL,
    date_time TIMESTAMP NOT NULL
)