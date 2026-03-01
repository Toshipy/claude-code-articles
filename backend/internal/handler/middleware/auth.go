package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/service"
)

func JWTAuth(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です"))
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("INVALID_TOKEN", "トークン形式が不正です"))
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, model.NewErrorResponse("INVALID_TOKEN", "トークンが無効または期限切れです"))
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role != "admin" {
				return c.JSON(http.StatusForbidden, model.NewErrorResponse("FORBIDDEN", "管理者権限が必要です"))
			}
			return next(c)
		}
	}
}
