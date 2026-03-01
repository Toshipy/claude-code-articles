BEGIN;

-- articles indexes
CREATE INDEX idx_articles_published_at ON articles (published_at DESC);
CREATE INDEX idx_articles_source_id ON articles (source_id);

-- Full-text search vector column and index
ALTER TABLE articles ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(summary, ''))
    ) STORED;
CREATE INDEX idx_articles_search ON articles USING GIN (search_vector);

-- article_tags indexes
CREATE INDEX idx_article_tags_tag_id ON article_tags (tag_id);

-- bookmarks indexes
CREATE INDEX idx_bookmarks_article_id ON bookmarks (article_id);

-- search_history indexes
CREATE INDEX idx_search_history_user_id_created ON search_history (user_id, created_at DESC);

-- sources indexes
CREATE INDEX idx_sources_is_active ON sources (is_active) WHERE is_active = true;

COMMIT;
