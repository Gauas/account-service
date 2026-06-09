package session

import (
	"net/http"
	"time"

	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/service"
	"github.com/labstack/echo/v4"
)

func Write(c echo.Context, cfg *config.Config, session *service.Session) {
	setCookie(c, cfg, "access_token", session.AccessToken, time.Until(session.ExpiresAt))
	setCookie(c, cfg, "refresh_token", session.RefreshToken, time.Until(session.RefreshExpiresAt))
}

func Response(session *service.Session) map[string]interface{} {
	return map[string]interface{}{
		"access_token":  session.AccessToken,
		"refresh_token": session.RefreshToken,
		"expires_in":    session.ExpiresIn,
	}
}

func Clear(c echo.Context, cfg *config.Config) {
	setCookie(c, cfg, "access_token", "", 0)
	setCookie(c, cfg, "refresh_token", "", 0)
}

func setCookie(c echo.Context, cfg *config.Config, name, value string, ttl time.Duration) {
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   cfg.Cookie.DomainName,
		Secure:   cfg.Cookie.Secure,
		HttpOnly: cfg.Cookie.HttpOnly,
		SameSite: cfg.Cookie.SameSite,
		MaxAge:   int(ttl.Seconds()),
		Expires:  time.Now().Add(ttl),
	})
}
