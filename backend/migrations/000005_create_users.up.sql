BEGIN;

CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    avatar_url   TEXT,
    role         VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
