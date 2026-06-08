package actions

import (
	"context"

	"github.com/gauas/account-service/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/dynamicpb"
)

type serviceServer interface {
	GetUserInfo(context.Context, *dynamicpb.Message) (*dynamicpb.Message, error)
}

type Action struct {
	service *service.Service
}

func Register(server grpc.ServiceRegistrar, svc *service.Service) {
	registerInfo(server, svc)
}
