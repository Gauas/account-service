package service

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repository   *repository.Registry
	Config       config.Config
	MessageQueue *amqp.Channel
	Cache        *redis.Client
}

func New(repo *repository.Registry, cfg config.Config, mq *amqp.Channel, cache *redis.Client) *Service {
	return &Service{
		Repository:   repo,
		Config:       cfg,
		MessageQueue: mq,
		Cache:        cache,
	}
}
