package config

import "os"

type Config struct {
	Addr   string
	APIKey string
}

func Load() Config {
	return Config{
		Addr:   getenv("ADDR", ":8080"),
		APIKey: getenv("APIKey", "dev-secret-key"),
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
