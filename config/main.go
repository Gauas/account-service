package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gauas/config-service/sdk"
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
	DBUrl     string
	SecretKey string

	MemoryURL string
	QueueURL  string

	AuthorizationURL string
	UploadURL        string

	Cookie Cookie
}

const localConfigPath = ".config/config.json"

func New() *Config {
	if cfg, ok := fromFile(); ok {
		return cfg
	}
	return fromSDK()
}

func fromFile() (*Config, bool) {
	data, err := os.ReadFile(localConfigPath)
	if err != nil {
		return nil, false
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatalf("config: invalid %s: %v", localConfigPath, err)
	}

	reader := func(key string) string {
		v := raw[key]
		if v == "" {
			log.Fatalf("config: %s is required in %s", key, localConfigPath)
		}
		return v
	}

	cfg := &Config{
		Port:      reader("PORT"),
		DBUrl:     reader("DB_URL"),
		SecretKey: reader("SECRET_KEY"),

		MemoryURL: reader("MEMORY_URL"),
		QueueURL:  reader("QUEUE_URL"),

		AuthorizationURL: reader("AUTHORIZATION_URL"),
		UploadURL:        reader("UPLOAD_URL"),

		Cookie: Cookie{
			HttpOnly:   reader("HTTP_ONLY") == "true",
			Secure:     reader("HTTP_SECURE") == "true",
			SameSite:   parseSameSite(mustInt(reader("SAME_SITE"))),
			DomainName: reader("DOMAIN_NAME"),
		},
	}

	validate(cfg)
	return cfg, true
}

func fromSDK() *Config {
	_ = godotenv.Load()

	client := sdk.New(sdk.Options{
		BaseURL:   mustEnv("CONFIG_SERVICE_URL"),
		SecretKey: mustEnv("SECRET_KEY"),
	})

	remote, err := client.Get("account-service", mustEnv("ENVIRONMENT"))
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cfg := &Config{
		Port:      remote.GetString("PORT", "8080"),
		DBUrl:     remote.GetString("DB_URL", ""),
		SecretKey: mustEnv("SECRET_KEY"),

		MemoryURL: remote.GetString("MEMORY_URL", "redis://localhost:6379/0"),
		QueueURL:  remote.GetString("QUEUE_URL", "amqp://guest:guest@localhost:5672/"),

		AuthorizationURL: remote.GetString("AUTHORIZATION_URL", "http://localhost:8082"),
		UploadURL:        remote.GetString("UPLOAD_URL", "http://localhost:8081"),

		Cookie: Cookie{
			HttpOnly:   remote.GetString("HTTP_ONLY", "false") == "true",
			Secure:     remote.GetString("HTTP_SECURE", "false") == "true",
			SameSite:   parseSameSite(int(remote.GetFloat64("SAME_SITE", 0))),
			DomainName: remote.GetString("DOMAIN_NAME", ""),
		},
	}

	validate(cfg)
	return cfg
}
