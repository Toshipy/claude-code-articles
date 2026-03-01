BEGIN;

CREATE TABLE article_tags (
    article_id BIGINT NOT NULL,
    tag_id     BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

COMMIT;
