package middlewares

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type contextKey string

const DeviceIDKey contextKey = "device_id"

func (m *Middleware) DeviceRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			deviceID := c.Request().Header.Get("X-Device-ID")
			if deviceID == "" {
				return echo.NewHTTPError(
					http.StatusBadRequest,
					"X-Device-ID header is required",
				)
			}

			ctx := context.WithValue(
				c.Request().Context(),
				DeviceIDKey,
				deviceID,
			)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
