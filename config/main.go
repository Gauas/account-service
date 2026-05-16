package config

import (
	"log"
	"os"
)

type Config struct {
	Port       string
	DBUrl      string
	SecretKey  string
	PrivateKey string

	MemoryURL string
	QueueURL  string

	AuthorizationURL string
	UploadURL        string

	CookieDomain string
	DomainName   string
}

func New() Config {
	cfg := Config{
		Port:       getEnv("PORT", "8080"),
		SecretKey:  os.Getenv("SECRET_KEY"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),

		DBUrl:     getEnv("DB_URL", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"),
		MemoryURL: getEnv("MEMORY_URL", "redis://localhost:6379/0"),
		QueueURL:  getEnv("QUEUE_URL", "amqp://guest:guest@localhost:5672/"),

		AuthorizationURL: getEnv("AUTHORIZATION_URL", "http://localhost:8082"),
		UploadURL:        getEnv("UPLOAD_URL", "http://localhost:8081"),

		CookieDomain: os.Getenv("GLOBAL_DOMAIN"),
		DomainName:   getEnv("DOMAIN_NAME", "gauas.online"),
	}

	validate(cfg)
	return cfg
}

func validate(cfg Config) {
	if cfg.Port == "" {
		log.Fatal("config: PORT is required")
	}
	if cfg.DBUrl == "" {
		log.Fatal("config: DB_URL is required")
	}
	if cfg.SecretKey == "" {
		log.Fatal("config: SECRET_KEY is required")
	}
	if cfg.PrivateKey == "" {
		log.Fatal("config: PRIVATE_KEY is required")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
