package postgres

import (
	"database/sql"
	"fmt"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
)

type BookmarkRepository struct {
	db *sql.DB
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Add(userID, articleID int64) (*model.Bookmark, error) {
	query := `INSERT INTO bookmarks (user_id, article_id) VALUES ($1, $2)
		ON CONFLICT (user_id, article_id) DO NOTHING
		RETURNING user_id, article_id, created_at`

	var b model.Bookmark
	err := r.db.QueryRow(query, userID, articleID).Scan(&b.UserID, &b.ArticleID, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("duplicate")
	}
	if err != nil {
		return nil, fmt.Errorf("add bookmark: %w", err)
	}
	return &b, nil
}

func (r *BookmarkRepository) Delete(userID, articleID int64) error {
	result, err := r.db.Exec(`DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`, userID, articleID)
	if err != nil {
		return fmt.Errorf("delete bookmark: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not_found")
	}
	return nil
}

func (r *BookmarkRepository) ListByUser(userID int64, page, perPage int) ([]model.Bookmark, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bookmarks: %w", err)
	}

	offset := (page - 1) * perPage
	query := `SELECT user_id, article_id, created_at FROM bookmarks
		WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		if err := rows.Scan(&b.UserID, &b.ArticleID, &b.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, b)
	}

	return bookmarks, total, nil
}
