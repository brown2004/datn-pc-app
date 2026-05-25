package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	BackendURL     string
	RequestTimeout time.Duration
	PollInterval   time.Duration
}

func Load() Config {
	backendURL := strings.TrimSpace(os.Getenv("PCAPP_BACKEND_URL"))
	if backendURL == "" {
		backendURL = "http://localhost:8081"
	}

	return Config{
		BackendURL:     strings.TrimRight(backendURL, "/"),
		RequestTimeout: durationFromEnv("PCAPP_BACKEND_TIMEOUT", 10*time.Second),
		PollInterval:   durationFromEnv("PCAPP_PAIRING_POLL_INTERVAL", 3*time.Second),
	}
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}
