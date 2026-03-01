package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) List(params model.ArticleListParams) ([]model.Article, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if params.Tag != "" {
		conditions = append(conditions, fmt.Sprintf(
			"a.id IN (SELECT at2.article_id FROM article_tags at2 JOIN tags t ON t.id = at2.tag_id WHERE t.slug = $%d)", argIdx))
		args = append(args, params.Tag)
		argIdx++
	}
	if params.SourceID > 0 {
		conditions = append(conditions, fmt.Sprintf("a.source_id = $%d", argIdx))
		args = append(args, params.SourceID)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	sortCol := "a.published_at"
	if params.Sort == "collected_at" {
		sortCol = "a.created_at"
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articles a %s", where)
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}

	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT a.id, a.source_id, a.title, a.url, a.thumbnail_url, a.summary, a.source, a.published_at, a.created_at, a.updated_at
		FROM articles a %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortCol, order, argIdx, argIdx+1)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.SourceID, &a.Title, &a.URL, &a.ThumbnailURL, &a.Summary, &a.Source, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}

	return articles, total, nil
}

func (r *ArticleRepository) GetByID(id int64) (*model.Article, error) {
	query := `SELECT id, source_id, title, url, thumbnail_url, summary, source, published_at, created_at, updated_at
		FROM articles WHERE id = $1`

	var a model.Article
	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.SourceID, &a.Title, &a.URL, &a.ThumbnailURL, &a.Summary, &a.Source, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get article: %w", err)
	}
	return &a, nil
}

func (r *ArticleRepository) Search(params model.SearchParams) ([]model.Article, int, error) {
	tsQuery := strings.Join(strings.Fields(params.Query), " & ")

	countQuery := `SELECT COUNT(*) FROM articles WHERE search_vector @@ to_tsquery('simple', $1)`
	var total int
	if err := r.db.QueryRow(countQuery, tsQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count search: %w", err)
	}

	offset := (params.Page - 1) * params.PerPage
	query := `SELECT id, source_id, title, url, thumbnail_url, summary, source, published_at, created_at, updated_at
		FROM articles
		WHERE search_vector @@ to_tsquery('simple', $1)
		ORDER BY ts_rank(search_vector, to_tsquery('simple', $1)) DESC, published_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, tsQuery, params.PerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("search articles: %w", err)
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.SourceID, &a.Title, &a.URL, &a.ThumbnailURL, &a.Summary, &a.Source, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}

	return articles, total, nil
}

func (r *ArticleRepository) GetRelated(articleID int64, limit int) ([]model.RelatedArticle, error) {
	query := `SELECT a.id, a.title, a.thumbnail_url
		FROM articles a
		JOIN article_tags at1 ON a.id = at1.article_id
		WHERE at1.tag_id IN (SELECT tag_id FROM article_tags WHERE article_id = $1)
		AND a.id != $1
		GROUP BY a.id, a.title, a.thumbnail_url, a.published_at
		ORDER BY COUNT(*) DESC, a.published_at DESC
		LIMIT $2`

	rows, err := r.db.Query(query, articleID, limit)
	if err != nil {
		return nil, fmt.Errorf("get related: %w", err)
	}
	defer rows.Close()

	var related []model.RelatedArticle
	for rows.Next() {
		var r model.RelatedArticle
		if err := rows.Scan(&r.ID, &r.Title, &r.ThumbnailURL); err != nil {
			return nil, fmt.Errorf("scan related: %w", err)
		}
		related = append(related, r)
	}
	return related, nil
}

func (r *ArticleRepository) Create(a *model.Article) error {
	query := `INSERT INTO articles (source_id, title, url, thumbnail_url, summary, source, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (url, published_at) DO NOTHING
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, a.SourceID, a.Title, a.URL, a.ThumbnailURL, a.Summary, a.Source, a.PublishedAt).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil // duplicate, skip
	}
	if err != nil {
		return fmt.Errorf("create article: %w", err)
	}
	return nil
}

func (r *ArticleRepository) ExistsByURL(url string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE url = $1)", url).Scan(&exists)
	return exists, err
}
