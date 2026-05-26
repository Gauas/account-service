package middlewares

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/infra"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Middleware struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(cfg *config.Config, infra *infra.Infra) *Middleware {
	return &Middleware{Config: cfg, Infra: infra}
}

func (m *Middleware) RegisterGlobal(server *echo.Echo) {
	server.Use(echoMiddleware.Recover())
	server.Use(echoMiddleware.Logger())
	server.Use(echoMiddleware.RequestID())
}
