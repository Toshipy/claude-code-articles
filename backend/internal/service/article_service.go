package service

import (
	"fmt"
	"math"

	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
)

type ArticleService struct {
	articleRepo *postgres.ArticleRepository
	tagRepo     *postgres.TagRepository
	sourceRepo  *postgres.SourceRepository
}

func NewArticleService(ar *postgres.ArticleRepository, tr *postgres.TagRepository, sr *postgres.SourceRepository) *ArticleService {
	return &ArticleService{articleRepo: ar, tagRepo: tr, sourceRepo: sr}
}

func (s *ArticleService) ListArticles(params model.ArticleListParams) ([]model.ArticleSummaryResponse, *model.Pagination, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 20
	}

	articles, total, err := s.articleRepo.List(params)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]model.ArticleSummaryResponse, 0, len(articles))
	for _, a := range articles {
		tags, _ := s.tagRepo.GetTagsForArticle(a.ID)
		if tags == nil {
			tags = []model.TagResponse{}
		}

		src, _ := s.sourceRepo.GetByID(a.SourceID)
		var srcBrief model.SourceBrief
		if src != nil {
			srcBrief = model.SourceBrief{ID: src.ID, Name: src.Name}
		}

		responses = append(responses, model.ArticleSummaryResponse{
			ID:           a.ID,
			Title:        a.Title,
			Summary:      a.Summary,
			URL:          a.URL,
			ThumbnailURL: a.ThumbnailURL,
			PublishedAt:  a.PublishedAt,
			CollectedAt:  a.CreatedAt,
			Source:       srcBrief,
			Tags:         tags,
		})
	}

	pagination := &model.Pagination{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PerPage))),
	}

	return responses, pagination, nil
}

func (s *ArticleService) GetArticle(id int64) (*model.ArticleDetailResponse, error) {
	a, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, nil
	}

	tags, _ := s.tagRepo.GetTagsForArticle(a.ID)
	if tags == nil {
		tags = []model.TagResponse{}
	}

	src, _ := s.sourceRepo.GetByID(a.SourceID)
	var srcDetail model.SourceDetailBrief
	if src != nil {
		srcDetail = model.SourceDetailBrief{ID: src.ID, Name: src.Name, URL: src.URL}
	}

	related, _ := s.articleRepo.GetRelated(a.ID, 5)
	if related == nil {
		related = []model.RelatedArticle{}
	}

	return &model.ArticleDetailResponse{
		ID:              a.ID,
		Title:           a.Title,
		Summary:         a.Summary,
		URL:             a.URL,
		ThumbnailURL:    a.ThumbnailURL,
		PublishedAt:     a.PublishedAt,
		CollectedAt:     a.CreatedAt,
		Source:          srcDetail,
		Tags:            tags,
		RelatedArticles: related,
	}, nil
}

func (s *ArticleService) SearchArticles(params model.SearchParams) ([]model.ArticleSummaryResponse, *model.Pagination, error) {
	if len(params.Query) < 2 {
		return nil, nil, fmt.Errorf("validation: query must be at least 2 characters")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 20
	}

	articles, total, err := s.articleRepo.Search(params)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]model.ArticleSummaryResponse, 0, len(articles))
	for _, a := range articles {
		tags, _ := s.tagRepo.GetTagsForArticle(a.ID)
		if tags == nil {
			tags = []model.TagResponse{}
		}

		src, _ := s.sourceRepo.GetByID(a.SourceID)
		var srcBrief model.SourceBrief
		if src != nil {
			srcBrief = model.SourceBrief{ID: src.ID, Name: src.Name}
		}

		responses = append(responses, model.ArticleSummaryResponse{
			ID:           a.ID,
			Title:        a.Title,
			Summary:      a.Summary,
			URL:          a.URL,
			ThumbnailURL: a.ThumbnailURL,
			PublishedAt:  a.PublishedAt,
			CollectedAt:  a.CreatedAt,
			Source:       srcBrief,
			Tags:         tags,
		})
	}

	pagination := &model.Pagination{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.PerPage))),
	}

	return responses, pagination, nil
}
