package relationship

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/service"
)

type Handler struct {
	Service *service.Service
	Config  *config.Config
}

func New(svc *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		Service: svc,
		Config:  cfg,
	}
}
