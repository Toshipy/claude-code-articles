package model

import "time"

type Source struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	FeedURL       *string    `json:"feed_url,omitempty"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type SourceResponse struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	IconURL         *string    `json:"icon_url,omitempty"`
	ArticleCount    int        `json:"article_count"`
	LastCollectedAt *time.Time `json:"last_collected_at,omitempty"`
}
