BEGIN;

CREATE TABLE bookmarks (
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, article_id)
);

COMMIT;
