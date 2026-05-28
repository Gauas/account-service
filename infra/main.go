package infra

import (
	"github.com/gauas/account-service/config"
	auth "github.com/gauas/authorization-service/sdk"
	upload "github.com/gauas/upload-service/sdk"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Infra struct {
	DB      *gorm.DB
	Memory  *redis.Client
	Queue   *amqp.Channel
	AuthSDK *auth.Client
	Upload  *upload.Client
}

func New(cfg *config.Config) *Infra {
	return &Infra{
		DB:     connectDatabase(cfg.DBUrl),
		Memory: connectMemory(cfg),
		Queue:  connectQueue(cfg.QueueURL),
		AuthSDK: auth.New(auth.Options{
			BaseURL:   cfg.AuthorizationURL,
			SecretKey: cfg.SecretKey,
		}),
		Upload: connectUpload(cfg),
	}
}
