package model

import "time"

type Bookmark struct {
	UserID    int64     `json:"user_id"`
	ArticleID int64     `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

type BookmarkResponse struct {
	ArticleID   int64                  `json:"article_id"`
	BookmarkedAt time.Time             `json:"bookmarked_at"`
	Article     ArticleSummaryResponse `json:"article"`
}
