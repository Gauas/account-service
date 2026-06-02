package controller

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller/auth"
	"github.com/gauas/account-service/controller/info"
	"github.com/gauas/account-service/controller/mfas"
	"github.com/gauas/account-service/service"
)

type Controller struct {
	Authentication *auth.Handler
	Info           *info.Handler
	MFA            *mfas.Handler
	//Verification   *verification.Handler
}

func New(svc *service.Service, cfg *config.Config) *Controller {
	return &Controller{
		Authentication: auth.New(svc, cfg),
		Info:           info.New(svc, cfg),
		MFA:            mfas.New(svc, cfg),
		//Verification:   verification.New(svc, cfg),
	}
}
