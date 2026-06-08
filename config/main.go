package config

import (
	"net/http"

	"github.com/joho/godotenv"
)

type Cookie struct {
	HttpOnly   bool
	Secure     bool
	SameSite   http.SameSite
	DomainName string
}
type Config struct {
	Port      string
	GRPCPort  string
	DBUrl     string
	SecretKey string

	MemoryURL string
	QueueURL  string

	AuthorizationURL string
	UploadURL        string

	Cookie Cookie
}

func New() *Config {
	return fromEnv()
}

func fromEnv() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:      get("PORT", "8080"),
		GRPCPort:  get("GRPC_PORT", "9090"),
		DBUrl:     mustEnv("DB_URL"),
		SecretKey: mustEnv("SECRET_KEY"),

		MemoryURL: get("MEMORY_URL", "redis://localhost:6379/0"),
		QueueURL:  get("QUEUE_URL", "amqp://guest:guest@localhost:5672/"),

		AuthorizationURL: get("AUTHORIZATION_URL", "http://localhost:8082"),
		UploadURL:        get("UPLOAD_URL", "http://localhost:8081"),

		Cookie: Cookie{
			HttpOnly:   get("HTTP_ONLY", "false") == "true",
			Secure:     get("HTTP_SECURE", "false") == "true",
			SameSite:   parseSameSite(getEnvInt("SAME_SITE", 0)),
			DomainName: get("DOMAIN_NAME", ""),
		},
	}

	validate(cfg)
	return cfg
}
