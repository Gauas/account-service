package middlewares

import "context"

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
