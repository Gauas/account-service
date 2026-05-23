package infra

import (
	"github.com/gauas/account-service/config"
	"github.com/gauas/authorization-service/sdk"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Infra struct {
	DB      *gorm.DB
	Memory  *redis.Client
	Queue   *amqp.Channel
	AuthSDK *sdk.Client
}

func New(cfg *config.Config) *Infra {
	return &Infra{
		DB:     connectDatabase(cfg.DBUrl),
		Memory: connectMemory(cfg),
		Queue:  connectQueue(cfg.QueueURL),
		AuthSDK: sdk.New(sdk.Options{
			BaseURL:   cfg.AuthorizationURL,
			SecretKey: cfg.SecretKey,
		}),
	}
}
