package actions

import (
	"context"

	"github.com/gauas/account-service/dto/response"
	"github.com/gauas/account-service/grpc/schema"
	"github.com/gauas/account-service/model"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	infoServiceName       = schema.ServiceName
	getUserInfoMethodName = schema.GetUserInfoMethodName
)

func registerInfo(server grpc.ServiceRegistrar, svc *service.Service) {
	action := &Action{service: svc}

	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: infoServiceName,
		HandlerType: (*serviceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: getUserInfoMethodName,
				Handler:    getUserInfoHandler,
			},
		},
	}, action)
}

func (a *Action) GetUserInfo(ctx context.Context, request *dynamicpb.Message) (*dynamicpb.Message, error) {
	key := schema.GetUserInfoRequestKey(request)

	currentUserKey, err := authenticate(ctx, a.service.Infra.AuthSDK)
	if err != nil {
		return nil, err
	}

	if key == "" {
		key = currentUserKey
	}

	user, err := a.service.GetProfileByKey(ctx, key)
	if err != nil {
		return nil, toStatusError(err)
	}

	profile := httpresp.Refine[*model.User, response.ProfileResponse](user)
	return schema.NewGetUserInfoResponse(profile), nil
}

func getUserInfoHandler(server interface{}, ctx context.Context, decode func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	request := schema.NewGetUserInfoRequest()
	if err := decode(request); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return server.(serviceServer).GetUserInfo(ctx, request)
	}

	info := &grpc.UnaryServerInfo{
		Server:     server,
		FullMethod: "/" + infoServiceName + "/" + getUserInfoMethodName,
	}

	handler := func(ctx context.Context, request interface{}) (interface{}, error) {
		return server.(serviceServer).GetUserInfo(ctx, request.(*dynamicpb.Message))
	}

	return interceptor(ctx, request, info, handler)
}
