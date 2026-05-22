package middlewares

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization token is required")
			}

			if err := m.service.ValidateAccessToken(context.Background(), token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			parsed, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "malformed token")
			}

			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			if sub, ok := claims["sub"].(string); ok {
				c.Set("user_id", sub)
			} else if uid, ok := claims["user_id"].(string); ok {
				c.Set("user_id", uid)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user_id in token")
			}

			if perm, ok := claims["permission"].(string); ok {
				c.Set("permission", perm)
			}

			return next(c)
		}
	}
}
