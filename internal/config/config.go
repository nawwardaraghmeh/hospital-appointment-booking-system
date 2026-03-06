package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds all configurable application values loaded from the environment
type Config struct {
	ServerPort            string
	DBPath                string
	AdminRegistrationCode string
	SessionDuration       int
}

// Load reads configuration from a .env file
func Load() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		ServerPort:            getEnv("SERVER_PORT", "8080"),
		DBPath:                getEnv("DB_PATH", "./data/abs.db"),
		AdminRegistrationCode: getEnv("ADMIN_REGISTRATION_CODE", ""),
		SessionDuration:       3600,
	}

	if cfg.AdminRegistrationCode == "" {
		return nil, fmt.Errorf("ADMIN_REGISTRATION_CODE is not set in .env or environment")
	}

	return cfg, nil
}

// loadDotEnv parses a .env file
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
