package service

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/infra"
	"github.com/gauas/account-service/repository"
)

type Service struct {
	Repository *repository.Registry
	Config     config.Config
	Infra      *infra.Infra
}

func New(repo *repository.Registry, cfg config.Config, infra *infra.Infra) *Service {
	return &Service{
		Repository: repo,
		Config:     cfg,
		Infra:      infra,
	}
}
