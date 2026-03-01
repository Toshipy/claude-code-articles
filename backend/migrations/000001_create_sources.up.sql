BEGIN;

CREATE TABLE sources (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    url             TEXT NOT NULL UNIQUE,
    feed_url        TEXT,
    last_fetched_at TIMESTAMPTZ,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
