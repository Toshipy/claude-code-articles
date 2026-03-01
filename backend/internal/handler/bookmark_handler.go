package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

type BookmarkHandler struct {
	bookmarkRepo   *postgres.BookmarkRepository
	articleService *service.ArticleService
}

func NewBookmarkHandler(br *postgres.BookmarkRepository, as *service.ArticleService) *BookmarkHandler {
	return &BookmarkHandler{bookmarkRepo: br, articleService: as}
}

type AddBookmarkRequest struct {
	ArticleID int64 `json:"article_id"`
}

func (h *BookmarkHandler) AddBookmark(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です"))
	}

	var req AddBookmarkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "リクエスト形式が不正です"))
	}

	if req.ArticleID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, model.NewValidationErrorResponse(
			"VALIDATION_ERROR",
			"記事IDが必要です",
			[]model.FieldError{{Field: "article_id", Message: "記事IDを指定してください"}},
		))
	}

	bookmark, err := h.bookmarkRepo.Add(userID, req.ArticleID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.JSON(http.StatusConflict, model.NewErrorResponse("DUPLICATE_RESOURCE", "この記事は既にブックマークされています"))
		}
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "ブックマークの追加に失敗しました"))
	}

	return c.JSON(http.StatusCreated, model.NewSuccessResponse(map[string]interface{}{
		"article_id": bookmark.ArticleID,
		"created_at": bookmark.CreatedAt,
	}))
}

func (h *BookmarkHandler) DeleteBookmark(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です"))
	}

	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "記事IDの形式が不正です"))
	}

	if err := h.bookmarkRepo.Delete(userID, articleID); err != nil {
		if strings.Contains(err.Error(), "not_found") {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", "指定されたブックマークが見つかりません"))
		}
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "ブックマークの削除に失敗しました"))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BookmarkHandler) ListBookmarks(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です"))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	bookmarks, total, err := h.bookmarkRepo.ListByUser(userID, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "ブックマーク一覧の取得に失敗しました"))
	}

	responses := make([]model.BookmarkResponse, 0, len(bookmarks))
	for _, b := range bookmarks {
		article, _ := h.articleService.GetArticle(b.ArticleID)
		var articleSummary model.ArticleSummaryResponse
		if article != nil {
			articleSummary = model.ArticleSummaryResponse{
				ID:           article.ID,
				Title:        article.Title,
				Summary:      article.Summary,
				URL:          article.URL,
				ThumbnailURL: article.ThumbnailURL,
				PublishedAt:  article.PublishedAt,
				CollectedAt:  article.CollectedAt,
				Tags:         article.Tags,
			}
		}

		responses = append(responses, model.BookmarkResponse{
			ArticleID:    b.ArticleID,
			BookmarkedAt: b.CreatedAt,
			Article:      articleSummary,
		})
	}

	pagination := model.Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	}

	return c.JSON(http.StatusOK, model.NewListResponse(responses, pagination))
}
