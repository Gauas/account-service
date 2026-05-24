package config

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("config: %s is required", key)
	return ""
}

func mustInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("config: invalid int %q", v)
	}
	return i
}

func validate(cfg *Config) {
	if cfg.Port == "" {
		log.Fatal("config: PORT is required")
	}
	if cfg.DBUrl == "" {
		log.Fatal("config: DB_URL is required")
	}
	if cfg.SecretKey == "" {
		log.Fatal("config: SECRET_KEY is required")
	}
}

func parseSameSite(v int) http.SameSite {
	switch v {
	case 1:
		return http.SameSiteLaxMode
	case 2:
		return http.SameSiteStrictMode
	case 3:
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return i
}
