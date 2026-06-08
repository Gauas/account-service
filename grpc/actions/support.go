package actions

import (
	"context"
	"errors"
	"strings"

	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/authorization-service/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func authenticate(ctx context.Context, authClient *sdk.Client) (string, error) {
	token := bearerToken(ctx)
	if token == "" {
		return "", status.Error(codes.Unauthenticated, "authorization token is required")
	}

	result, err := authClient.ValidateToken(ctx, token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	return result.UserID, nil
}

func bearerToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return ""
	}

	auth := values[0]
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}

func toStatusError(err error) error {
	var appErr *httpresp.Error
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case 400:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case 401:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case 403:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case 404:
			return status.Error(codes.NotFound, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}
