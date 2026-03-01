package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
)

type SourceRepository struct {
	db *sql.DB
}

func NewSourceRepository(db *sql.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) List() ([]model.Source, error) {
	query := `SELECT id, name, url, feed_url, last_fetched_at, is_active, created_at, updated_at
		FROM sources ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list sources: %w", err)
	}
	defer rows.Close()

	var sources []model.Source
	for rows.Next() {
		var s model.Source
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.FeedURL, &s.LastFetchedAt, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan source: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func (r *SourceRepository) ListWithArticleCount() ([]model.SourceResponse, error) {
	query := `SELECT s.id, s.name, s.url, COUNT(a.id) AS article_count, s.last_fetched_at
		FROM sources s
		LEFT JOIN articles a ON s.id = a.source_id
		WHERE s.is_active = true
		GROUP BY s.id, s.name, s.url, s.last_fetched_at
		ORDER BY s.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list sources with count: %w", err)
	}
	defer rows.Close()

	var sources []model.SourceResponse
	for rows.Next() {
		var s model.SourceResponse
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.ArticleCount, &s.LastCollectedAt); err != nil {
			return nil, fmt.Errorf("scan source: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func (r *SourceRepository) GetByID(id int64) (*model.Source, error) {
	query := `SELECT id, name, url, feed_url, last_fetched_at, is_active, created_at, updated_at
		FROM sources WHERE id = $1`

	var s model.Source
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.URL, &s.FeedURL, &s.LastFetchedAt, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get source: %w", err)
	}
	return &s, nil
}

func (r *SourceRepository) ListActive() ([]model.Source, error) {
	query := `SELECT id, name, url, feed_url, last_fetched_at, is_active, created_at, updated_at
		FROM sources WHERE is_active = true ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list active sources: %w", err)
	}
	defer rows.Close()

	var sources []model.Source
	for rows.Next() {
		var s model.Source
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.FeedURL, &s.LastFetchedAt, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan source: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func (r *SourceRepository) UpdateLastFetched(id int64, t time.Time) error {
	_, err := r.db.Exec(`UPDATE sources SET last_fetched_at = $1, updated_at = NOW() WHERE id = $2`, t, id)
	return err
}
