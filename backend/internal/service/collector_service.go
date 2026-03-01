package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
)

type CollectorService struct {
	sourceRepo  *postgres.SourceRepository
	articleRepo *postgres.ArticleRepository
	tagRepo     *postgres.TagRepository
	parser      *gofeed.Parser
}

func NewCollectorService(sr *postgres.SourceRepository, ar *postgres.ArticleRepository, tr *postgres.TagRepository) *CollectorService {
	return &CollectorService{
		sourceRepo:  sr,
		articleRepo: ar,
		tagRepo:     tr,
		parser:      gofeed.NewParser(),
	}
}

type CollectResult struct {
	SourceName    string `json:"source_name"`
	ArticlesFound int    `json:"articles_found"`
	ArticlesSaved int    `json:"articles_saved"`
	Errors        int    `json:"errors"`
}

func (s *CollectorService) CollectAll(ctx context.Context) ([]CollectResult, error) {
	sources, err := s.sourceRepo.ListActive()
	if err != nil {
		return nil, err
	}

	var results []CollectResult
	for _, src := range sources {
		if src.FeedURL == nil {
			continue
		}
		result := s.collectFromSource(ctx, src)
		results = append(results, result)
	}
	return results, nil
}

func (s *CollectorService) CollectBySourceID(ctx context.Context, sourceID int64) (*CollectResult, error) {
	src, err := s.sourceRepo.GetByID(sourceID)
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, nil
	}
	if src.FeedURL == nil {
		return &CollectResult{SourceName: src.Name}, nil
	}

	result := s.collectFromSource(ctx, *src)
	return &result, nil
}

func (s *CollectorService) collectFromSource(ctx context.Context, src model.Source) CollectResult {
	result := CollectResult{SourceName: src.Name}

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := s.parser.ParseURLWithContext(*src.FeedURL, ctxTimeout)
	if err != nil {
		log.Printf("Error fetching feed %s: %v", src.Name, err)
		result.Errors = 1
		return result
	}

	result.ArticlesFound = len(feed.Items)

	for _, item := range feed.Items {
		exists, err := s.articleRepo.ExistsByURL(item.Link)
		if err != nil {
			result.Errors++
			continue
		}
		if exists {
			continue
		}

		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		var summary *string
		if item.Description != "" {
			desc := item.Description
			if len(desc) > 500 {
				desc = desc[:500] + "..."
			}
			summary = &desc
		}

		var thumbnailURL *string
		if item.Image != nil && item.Image.URL != "" {
			thumbnailURL = &item.Image.URL
		}

		article := &model.Article{
			SourceID:     src.ID,
			Title:        item.Title,
			URL:          item.Link,
			ThumbnailURL: thumbnailURL,
			Summary:      summary,
			Source:       src.Name,
			PublishedAt:  publishedAt,
		}

		if err := s.articleRepo.Create(article); err != nil {
			result.Errors++
			continue
		}

		if article.ID > 0 {
			result.ArticlesSaved++
			s.assignTags(article, item.Categories)
		}
	}

	now := time.Now()
	if err := s.sourceRepo.UpdateLastFetched(src.ID, now); err != nil {
		log.Printf("Error updating last_fetched_at for source %s: %v", src.Name, err)
	}

	return result
}

func (s *CollectorService) assignTags(article *model.Article, categories []string) {
	for _, cat := range categories {
		cat = strings.TrimSpace(cat)
		if cat == "" {
			continue
		}
		slug := strings.ToLower(strings.ReplaceAll(cat, " ", "-"))
		tagID, err := s.tagRepo.GetOrCreate(cat, slug)
		if err != nil {
			log.Printf("Error creating tag %s: %v", cat, err)
			continue
		}
		if err := s.tagRepo.AddArticleTag(article.ID, tagID); err != nil {
			log.Printf("Error linking tag %s to article %d: %v", cat, article.ID, err)
		}
	}
}
