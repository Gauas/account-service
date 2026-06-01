package infra

import (
	"log"

	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/packages/uploader"
)

func connectUpload(cfg *config.Config) *uploader.Client {
	client := uploader.New(cfg.UploadURL, cfg.SecretKey)

	log.Println("infra: upload sdk connected")

	return client
}
