package middlewares

import (
	"context"

	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	userIDKey     contextKey = "user_id"
	permissionKey contextKey = "permission"
	deviceIDKey   contextKey = "device_id"
)

func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)

	return userID
}

func Permission(ctx context.Context) string {
	permission, _ := ctx.Value(permissionKey).(string)

	return permission
}

func DeviceID(ctx context.Context) string {
	deviceID, _ := ctx.Value(deviceIDKey).(string)

	return deviceID
}

func RefreshToken(ctx echo.Context) string {
	token := ctx.Request().Header.Get("X-Refresh-Token")

	if token != "" {
		return token
	}

	cookie, err := ctx.Cookie("refresh_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}
