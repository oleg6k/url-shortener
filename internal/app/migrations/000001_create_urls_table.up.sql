CREATE TABLE IF NOT EXISTS urls (
    id SERIAL primary key,
    short_url VARCHAR UNIQUE NOT NULL,
    original_url VARCHAR NOT NULL
)