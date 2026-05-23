package config

import (
	"net/http"
	"os"
)

type Cookie struct {
	HttpOnly   bool
	Secure     bool
	SameSite   http.SameSite
	DomainName string
}
type Config struct {
	Port      string
	DBUrl     string
	SecretKey string

	MemoryURL string
	QueueURL  string

	AuthorizationURL string
	UploadURL        string

	Cookie Cookie
}

func New() *Config {
	cfg := &Config{
		Port:      get("PORT", "8080"),
		DBUrl:     get("DB_URL", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"),
		SecretKey: os.Getenv("SECRET_KEY"),

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
