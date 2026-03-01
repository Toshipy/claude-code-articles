package model

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TagWithCount struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int    `json:"article_count"`
}
