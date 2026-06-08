package grpc

import (
	"github.com/gauas/account-service/grpc/actions"
	"github.com/gauas/account-service/grpc/runtime"
	"github.com/gauas/account-service/service"
)

type Server = runtime.Server

func Register(port string, service *service.Service) *Server {
	return runtime.Register(port, "account-service", func(server runtime.ServiceRegistrar) {
		actions.Register(server, service)
	})
}
