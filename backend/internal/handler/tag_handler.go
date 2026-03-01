package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

type TagHandler struct {
	tagRepo        *postgres.TagRepository
	articleService *service.ArticleService
}

func NewTagHandler(tr *postgres.TagRepository, as *service.ArticleService) *TagHandler {
	return &TagHandler{tagRepo: tr, articleService: as}
}

func (h *TagHandler) ListTags(c echo.Context) error {
	tags, err := h.tagRepo.ListWithCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "タグ一覧の取得に失敗しました"))
	}
	if tags == nil {
		tags = []model.TagWithCount{}
	}

	return c.JSON(http.StatusOK, model.NewSuccessResponse(tags))
}

func (h *TagHandler) GetArticlesByTag(c echo.Context) error {
	slug := c.Param("slug")

	tag, err := h.tagRepo.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "タグの取得に失敗しました"))
	}
	if tag == nil {
		return c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", "指定されたタグが見つかりません"))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	params := model.ArticleListParams{
		Page:    page,
		PerPage: perPage,
		Tag:     slug,
		Sort:    c.QueryParam("sort"),
		Order:   c.QueryParam("order"),
	}

	articles, pagination, err := h.articleService.ListArticles(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "記事一覧の取得に失敗しました"))
	}

	return c.JSON(http.StatusOK, model.NewListResponse(articles, *pagination))
}
