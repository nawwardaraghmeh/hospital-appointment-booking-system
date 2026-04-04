package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configurable values loaded from the environment
type Config struct {
	AuthPort              string
	BookingPort           string
	BookingURL            string
	DBPath                string
	AdminRegistrationCode string
	SessionDuration       int
}

// Load reads configuration from a .env file
func Load() (*Config, error) {
	loadDotEnv(".env")

	sessionDuration := 3600
	if v := os.Getenv("SESSION_DURATION"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			sessionDuration = n
		}
	}

	cfg := &Config{
		AuthPort:              getEnv("AUTH_PORT", "8080"),
		BookingPort:           getEnv("BOOKING_PORT", "8081"),
		BookingURL:            getEnv("BOOKING_URL", "http://localhost:8081"),
		DBPath:                getEnv("DB_PATH", "./data/abs.db"),
		AdminRegistrationCode: getEnv("ADMIN_REGISTRATION_CODE", ""),
		SessionDuration:       sessionDuration,
	}

	if cfg.AdminRegistrationCode == "" {
		return nil, fmt.Errorf("ADMIN_REGISTRATION_CODE is not set in .env or environment")
	}

	return cfg, nil
}

// loadDotEnv parses a .env file and sets each key=value as an environment variable
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
