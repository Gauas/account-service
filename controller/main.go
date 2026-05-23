package controller

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/service"
)

type Controller struct {
	service *service.Service
	config  *config.Config
}

func New(svc *service.Service, config *config.Config) *Controller {
	return &Controller{service: svc, config: config}
}
