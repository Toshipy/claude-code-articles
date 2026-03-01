package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(as *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: as}
}

func (h *ArticleHandler) ListArticles(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	sourceID, _ := strconv.ParseInt(c.QueryParam("source_id"), 10, 64)

	params := model.ArticleListParams{
		Page:     page,
		PerPage:  perPage,
		Tag:      c.QueryParam("tag"),
		SourceID: sourceID,
		Sort:     c.QueryParam("sort"),
		Order:    c.QueryParam("order"),
	}

	articles, pagination, err := h.articleService.ListArticles(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "記事一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, model.NewListResponse(articles, *pagination))
}

func (h *ArticleHandler) GetArticle(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "IDの形式が不正です"))
	}

	article, err := h.articleService.GetArticle(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "記事の取得に失敗しました"))
	}
	if article == nil {
		return c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", "指定された記事が見つかりません"))
	}

	return c.JSON(http.StatusOK, model.NewSuccessResponse(article))
}

func (h *ArticleHandler) SearchArticles(c echo.Context) error {
	q := strings.TrimSpace(c.QueryParam("q"))
	if len(q) < 2 {
		return c.JSON(http.StatusUnprocessableEntity, model.NewValidationErrorResponse(
			"VALIDATION_ERROR",
			"検索キーワードは2文字以上入力してください",
			[]model.FieldError{{Field: "q", Message: "2文字以上入力してください"}},
		))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	params := model.SearchParams{
		Query:   q,
		Page:    page,
		PerPage: perPage,
	}

	articles, pagination, err := h.articleService.SearchArticles(params)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation:") {
			return c.JSON(http.StatusUnprocessableEntity, model.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "検索に失敗しました"))
	}

	return c.JSON(http.StatusOK, model.NewListResponse(articles, *pagination))
}
