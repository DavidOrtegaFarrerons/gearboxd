CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username text UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool DEFAULT false NOT NULL,
    version int DEFAULT 1 NOT NULL,
    created_at timestamptz DEFAULT NOW()
);