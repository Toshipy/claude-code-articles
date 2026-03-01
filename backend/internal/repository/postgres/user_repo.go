package postgres

import (
	"database/sql"
	"fmt"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	query := `SELECT id, email, display_name, avatar_url, role, created_at FROM users WHERE id = $1`
	var u model.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, display_name, avatar_url, role, created_at FROM users WHERE email = $1`
	var u model.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Create(u *model.User) error {
	query := `INSERT INTO users (email, display_name, avatar_url, role) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query, u.Email, u.DisplayName, u.AvatarURL, u.Role).Scan(&u.ID, &u.CreatedAt)
}

func (r *UserRepository) GetBookmarkCount(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}
