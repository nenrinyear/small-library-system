package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Middleware struct {
	adminToken string
}

func New(adminToken string) Middleware {
	return Middleware{adminToken: strings.TrimSpace(adminToken)}
}

func (m Middleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if m.adminToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "admin API token is not configured",
			})
		}

		authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "missing bearer token",
			})
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" || token != m.adminToken {
			return c.JSON(http.StatusForbidden, map[string]any{
				"error": "invalid token",
			})
		}

		return next(c)
	}
}
