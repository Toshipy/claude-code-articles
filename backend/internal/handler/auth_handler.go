package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(as *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

func (h *AuthHandler) GoogleAuth(c echo.Context) error {
	var req GoogleAuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "リクエスト形式が不正です"))
	}

	if req.IDToken == "" {
		return c.JSON(http.StatusBadRequest, model.NewValidationErrorResponse(
			"VALIDATION_ERROR",
			"IDトークンが必要です",
			[]model.FieldError{{Field: "id_token", Message: "IDトークンを指定してください"}},
		))
	}

	// In production, verify the Google ID token with Google's API.
	// For now, we extract a mock email/name from the token for development.
	// TODO: Implement proper Google ID token verification
	email := "user@example.com"
	displayName := "Test User"
	avatarURL := ""

	resp, err := h.authService.AuthenticateWithGoogle(email, displayName, avatarURL)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("INVALID_TOKEN", "Google IDトークンの検証に失敗しました"))
	}

	return c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です"))
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", "ユーザーが見つかりません"))
	}

	bookmarkCount, _ := h.authService.GetBookmarkCount(userID)

	return c.JSON(http.StatusOK, model.NewSuccessResponse(model.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.DisplayName,
		AvatarURL:     user.AvatarURL,
		Role:          user.Role,
		BookmarkCount: bookmarkCount,
		CreatedAt:     user.CreatedAt,
	}))
}
