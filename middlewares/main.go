package middlewares

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/service"
)

type Middleware struct {
	config  config.Config
	service *service.Service
}

func New(cfg config.Config, svc *service.Service) *Middleware {
	return &Middleware{config: cfg, service: svc}
}

func (m *Middleware) RegisterGlobal(server *echo.Echo) {
	server.Use(echoMiddleware.Recover())
	server.Use(echoMiddleware.Logger())
	server.Use(echoMiddleware.RequestID())
}
