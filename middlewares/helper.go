package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func extractToken(c echo.Context) string {
	if auth := c.Request().Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
	}

	if cookie, err := c.Cookie("access_token"); err == nil {
		return cookie.Value
	}

	if q := c.QueryParam("access_token"); q != "" {
		return q
	}

	return ""
}
