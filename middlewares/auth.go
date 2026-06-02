package middlewares

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization token is required")
			}

			result, err := m.Infra.AuthSDK.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			ctx := c.Request().Context()

			ctx = context.WithValue(ctx, userIDKey, result.UserID)

			ctx = context.WithValue(ctx, permissionKey, result.Permission)

			ctx = context.WithValue(ctx, deviceIDKey, result.DeviceID)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
