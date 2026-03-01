BEGIN;

DROP INDEX IF EXISTS idx_sources_is_active;
DROP INDEX IF EXISTS idx_search_history_user_id_created;
DROP INDEX IF EXISTS idx_bookmarks_article_id;
DROP INDEX IF EXISTS idx_article_tags_tag_id;
DROP INDEX IF EXISTS idx_articles_search;
ALTER TABLE articles DROP COLUMN IF EXISTS search_vector;
DROP INDEX IF EXISTS idx_articles_source_id;
DROP INDEX IF EXISTS idx_articles_published_at;

COMMIT;
