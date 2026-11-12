package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	envPort       = "MAILCHECKER_PORT"
	envDBPath     = "MAILCHECKER_DB_PATH"
	defaultPort   = 8080
	defaultDBPath = "data/mail-checker.db"
	minPort       = 1
	maxPort       = 65535
)

// Config is read once at startup; any invalid value stops the app (fail fast).
type Config struct {
	Port   int
	DBPath string
}

func LoadConfig(getenv func(string) string) (Config, error) {
	port, err := readPort(getenv)
	if err != nil {
		return Config{}, err
	}
	dbPath, err := readDBPath(getenv)
	if err != nil {
		return Config{}, err
	}
	return Config{Port: port, DBPath: dbPath}, nil
}

func readPort(getenv func(string) string) (int, error) {
	raw := getenv(envPort)
	if raw == "" {
		return defaultPort, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port < minPort || port > maxPort {
		return 0, fmt.Errorf("%s=%q: must be a number between %d and %d", envPort, raw, minPort, maxPort)
	}
	return port, nil
}

func readDBPath(getenv func(string) string) (string, error) {
	raw := getenv(envDBPath)
	if raw == "" {
		return defaultDBPath, nil
	}
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("%s=%q: must not be blank", envDBPath, raw)
	}
	return raw, nil
}
