package model

import "time"

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserResponse struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Role          string    `json:"role"`
	BookmarkCount int       `json:"bookmark_count,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
