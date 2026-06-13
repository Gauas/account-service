package service

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/infra"
	"github.com/gauas/account-service/publisher"
	"github.com/gauas/account-service/repository"
)

type Service struct {
	Repository *repository.Registry
	Publisher  *publisher.Registry
	Config     *config.Config
	Infra      *infra.Infra
}

func New(repo *repository.Registry, pub *publisher.Registry, cfg *config.Config, infra *infra.Infra) *Service {
	return &Service{
		Repository: repo,
		Publisher:  pub,
		Config:     cfg,
		Infra:      infra,
	}
}
