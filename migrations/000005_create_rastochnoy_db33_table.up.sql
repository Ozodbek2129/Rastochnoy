CREATE TABLE IF NOT EXISTS rastochnoy_writedb33 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(30) UNIQUE NOT NULL,
    offsett FLOAT UNIQUE NOT NULL,
    value float NOT NULL
);