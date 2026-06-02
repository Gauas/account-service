package controller

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller/auth"
	"github.com/gauas/account-service/controller/info"
	"github.com/gauas/account-service/controller/mfas/totp"
	"github.com/gauas/account-service/service"
)

type Controller struct {
	Authentication *auth.Handler
	Info           *info.Handler
	TOTP           *totp.Handler
	//Verification   *verification.Handler
}

func New(svc *service.Service, cfg *config.Config) *Controller {
	return &Controller{
		Authentication: auth.New(svc, cfg),
		Info:           info.New(svc, cfg),
		TOTP:           totp.New(svc, cfg),
		//Verification:   verification.New(svc, cfg),
	}
}
