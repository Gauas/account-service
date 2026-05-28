package infra

import (
	"log"

	"github.com/gauas/account-service/config"
	upload "github.com/gauas/upload-service/sdk"
)

func connectUpload(cfg *config.Config) *upload.Client {
	client, err := upload.NewClient(upload.Config{
		BaseURL:   cfg.UploadURL,
		SecretKey: cfg.SecretKey,
	})
	if err != nil {
		log.Fatalf("infra: failed to init upload sdk: %v", err)
	}

	log.Println("infra: upload sdk connected")

	return client
}
