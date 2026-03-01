package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
)

type SourceHandler struct {
	sourceRepo *postgres.SourceRepository
}

func NewSourceHandler(sr *postgres.SourceRepository) *SourceHandler {
	return &SourceHandler{sourceRepo: sr}
}

func (h *SourceHandler) ListSources(c echo.Context) error {
	sources, err := h.sourceRepo.ListWithArticleCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "ソース一覧の取得に失敗しました"))
	}
	if sources == nil {
		sources = []model.SourceResponse{}
	}

	return c.JSON(http.StatusOK, model.NewSuccessResponse(sources))
}
