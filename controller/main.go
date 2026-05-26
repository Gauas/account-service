package controller

import (
	"github.com/gauas/account-service/config"
	auth "github.com/gauas/account-service/controller/auth"
	"github.com/gauas/account-service/controller/profile"
	"github.com/gauas/account-service/service"
)

type Controller struct {
	Authentication *auth.Handler
	Profile        *profile.Handler
	//Verification   *verification.Handler
}

func New(svc *service.Service, cfg *config.Config) *Controller {
	return &Controller{
		Authentication: auth.New(svc, cfg),
		Profile:        profile.New(svc, cfg),
		//Verification:   verification.New(svc, cfg),
	}
}
