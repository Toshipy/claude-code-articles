package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

type AdminHandler struct {
	collectorService *service.CollectorService
}

func NewAdminHandler(cs *service.CollectorService) *AdminHandler {
	return &AdminHandler{collectorService: cs}
}

type CollectRequest struct {
	SourceID *int64 `json:"source_id,omitempty"`
}

func (h *AdminHandler) TriggerCollect(c echo.Context) error {
	var req CollectRequest
	_ = c.Bind(&req)

	if req.SourceID != nil {
		result, err := h.collectorService.CollectBySourceID(c.Request().Context(), *req.SourceID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "記事収集に失敗しました"))
		}
		if result == nil {
			return c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", "指定されたソースが見つかりません"))
		}

		return c.JSON(http.StatusOK, model.NewSuccessResponse(map[string]interface{}{
			"status":         "completed",
			"target_sources": []string{result.SourceName},
			"started_at":     time.Now().UTC(),
			"result":         result,
		}))
	}

	results, err := h.collectorService.CollectAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_ERROR", "記事収集に失敗しました"))
	}

	sourceNames := make([]string, len(results))
	for i, r := range results {
		sourceNames[i] = r.SourceName
	}

	return c.JSON(http.StatusOK, model.NewSuccessResponse(map[string]interface{}{
		"status":         "completed",
		"target_sources": sourceNames,
		"started_at":     time.Now().UTC(),
		"results":        results,
	}))
}
