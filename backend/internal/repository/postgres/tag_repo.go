package postgres

import (
	"database/sql"
	"fmt"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) ListWithCount() ([]model.TagWithCount, error) {
	query := `SELECT t.id, t.name, t.slug, COUNT(at2.article_id) AS article_count
		FROM tags t
		LEFT JOIN article_tags at2 ON t.id = at2.tag_id
		GROUP BY t.id, t.name, t.slug
		ORDER BY article_count DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []model.TagWithCount
	for rows.Next() {
		var t model.TagWithCount
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.ArticleCount); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) GetBySlug(slug string) (*model.Tag, error) {
	query := `SELECT id, name, slug, created_at FROM tags WHERE slug = $1`
	var t model.Tag
	err := r.db.QueryRow(query, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tag by slug: %w", err)
	}
	return &t, nil
}

func (r *TagRepository) GetTagsForArticle(articleID int64) ([]model.TagResponse, error) {
	query := `SELECT t.id, t.name, t.slug
		FROM tags t
		JOIN article_tags at2 ON t.id = at2.tag_id
		WHERE at2.article_id = $1`

	rows, err := r.db.Query(query, articleID)
	if err != nil {
		return nil, fmt.Errorf("get tags for article: %w", err)
	}
	defer rows.Close()

	var tags []model.TagResponse
	for rows.Next() {
		var t model.TagResponse
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) GetOrCreate(name, slug string) (int64, error) {
	var id int64
	err := r.db.QueryRow(`SELECT id FROM tags WHERE slug = $1`, slug).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = r.db.QueryRow(`INSERT INTO tags (name, slug) VALUES ($1, $2) ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name RETURNING id`, name, slug).Scan(&id)
	return id, err
}

func (r *TagRepository) AddArticleTag(articleID, tagID int64) error {
	_, err := r.db.Exec(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, articleID, tagID)
	return err
}
