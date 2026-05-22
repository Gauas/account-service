package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func userIDFromContext(c echo.Context) (uuid.UUID, error) {
	raw := c.Get("user_id")
	if raw == nil {
		return uuid.Nil, fmt.Errorf("user_id not in context")
	}

	switch v := raw.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	}

	return uuid.Nil, fmt.Errorf("invalid user_id type")
}

func deviceIDFromRequest(c echo.Context, fromBody string) string {
	if fromBody != "" {
		return fromBody
	}

	return c.Request().Header.Get("X-Device-ID")
}

func resolveIdentifier(req loginRequest) (string, string) {
	switch {
	case req.Username != nil:
		return "username", *req.Username
	case req.Email != nil:
		return "email", *req.Email
	case req.Phone != nil:
		return "phone", *req.Phone
	default:
		return "", ""
	}
}

func (ctrl *Controller) setAccessCookie(c echo.Context, token string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   ctrl.cookieDomain,
		HttpOnly: true,
	})
}

func (ctrl *Controller) setRefreshCookie(c echo.Context, token string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   ctrl.cookieDomain,
		HttpOnly: true,
	})
}
