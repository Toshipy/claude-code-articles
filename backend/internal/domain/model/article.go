package model

import "time"

type Article struct {
	ID           int64     `json:"id"`
	SourceID     int64     `json:"source_id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Summary      *string   `json:"summary,omitempty"`
	Source       string    `json:"source"`
	PublishedAt  time.Time `json:"published_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ArticleSummaryResponse struct {
	ID           int64           `json:"id"`
	Title        string          `json:"title"`
	Summary      *string         `json:"summary,omitempty"`
	URL          string          `json:"url"`
	ThumbnailURL *string         `json:"thumbnail_url,omitempty"`
	PublishedAt  time.Time       `json:"published_at"`
	CollectedAt  time.Time       `json:"collected_at"`
	Source       SourceBrief     `json:"source"`
	Tags         []TagResponse   `json:"tags"`
}

type SourceBrief struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	IconURL *string `json:"icon_url,omitempty"`
}

type ArticleDetailResponse struct {
	ID              int64                  `json:"id"`
	Title           string                 `json:"title"`
	Summary         *string                `json:"summary,omitempty"`
	URL             string                 `json:"url"`
	ThumbnailURL    *string                `json:"thumbnail_url,omitempty"`
	PublishedAt     time.Time              `json:"published_at"`
	CollectedAt     time.Time              `json:"collected_at"`
	Source          SourceDetailBrief      `json:"source"`
	Tags            []TagResponse          `json:"tags"`
	RelatedArticles []RelatedArticle       `json:"related_articles"`
}

type SourceDetailBrief struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	URL     string  `json:"url"`
	IconURL *string `json:"icon_url,omitempty"`
}

type RelatedArticle struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
}

type ArticleListParams struct {
	Page     int
	PerPage  int
	Tag      string
	SourceID int64
	Sort     string
	Order    string
}

type SearchParams struct {
	Query   string
	Page    int
	PerPage int
}
