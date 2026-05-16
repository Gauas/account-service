package controller

import "github.com/gauas/account-service/service"

type Controller struct {
	service      *service.Service
	cookieDomain string
}

func New(svc *service.Service, cookieDomain string) *Controller {
	return &Controller{service: svc, cookieDomain: cookieDomain}
}
