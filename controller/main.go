package controller

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller/auth"
	"github.com/gauas/account-service/controller/profile"
	"github.com/gauas/account-service/controller/totp"
	"github.com/gauas/account-service/service"
)

type Controller struct {
	Authentication *auth.Handler
	Profile        *profile.Handler
	MFA            *totp.Handler
	//Verification   *verification.Handler
}

func New(svc *service.Service, cfg *config.Config) *Controller {
	return &Controller{
		Authentication: auth.New(svc, cfg),
		Profile:        profile.New(svc, cfg),
		MFA:            totp.New(svc, cfg),
		//Verification:   verification.New(svc, cfg),
	}
}
