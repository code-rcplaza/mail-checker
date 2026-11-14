package main

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

const (
	envPort       = "MAILCHECKER_PORT"
	envDBPath     = "MAILCHECKER_DB_PATH"
	envMode       = "MAILCHECKER_ENV"
	envResendKey  = "MAILCHECKER_RESEND_API_KEY"
	envMailFrom   = "MAILCHECKER_MAIL_FROM"
	modeDev       = "development"
	modeProd      = "production"
	defaultPort   = 8080
	defaultDBPath = "data/mail-checker.db"
	minPort       = 1
	maxPort       = 65535
)

// Config is read once at startup; any invalid value stops the app (fail fast).
type Config struct {
	Port         int
	DBPath       string
	IsProduction bool
	Mailer       MailerConfig
}

// Secret never shows its value when printed or logged.
type Secret string

const redacted = "[redacted]"

func (Secret) String() string               { return redacted }
func (Secret) GoString() string             { return redacted }
func (Secret) LogValue() slog.Value         { return slog.StringValue(redacted) }
func (Secret) MarshalText() ([]byte, error) { return []byte(redacted), nil }

type MailerConfig struct {
	ResendAPIKey Secret
	MailFrom     string
}

// UsesConsole is true in development without a key; production always has one.
func (m MailerConfig) UsesConsole() bool { return m.ResendAPIKey == "" }

func LoadConfig(getenv func(string) string) (Config, error) {
	port, err := readPort(getenv)
	if err != nil {
		return Config{}, err
	}
	dbPath, err := readDBPath(getenv)
	if err != nil {
		return Config{}, err
	}
	isProduction, err := readMode(getenv)
	if err != nil {
		return Config{}, err
	}
	mailer, err := readMailer(getenv, isProduction)
	if err != nil {
		return Config{}, err
	}
	return Config{Port: port, DBPath: dbPath, IsProduction: isProduction, Mailer: mailer}, nil
}

// readMode defaults to development: production is always set explicitly (Docker).
func readMode(getenv func(string) string) (bool, error) {
	switch raw := getenv(envMode); raw {
	case "", modeDev:
		return false, nil
	case modeProd:
		return true, nil
	default:
		return false, fmt.Errorf("%s=%q: must be %s or %s", envMode, raw, modeDev, modeProd)
	}
}

func readMailer(getenv func(string) string, isProduction bool) (MailerConfig, error) {
	mailer := MailerConfig{ResendAPIKey: Secret(strings.TrimSpace(getenv(envResendKey))), MailFrom: strings.TrimSpace(getenv(envMailFrom))}
	if !isProduction {
		return mailer, nil
	}
	if mailer.ResendAPIKey == "" {
		return MailerConfig{}, fmt.Errorf("%s is required in %s", envResendKey, modeProd)
	}
	if mailer.MailFrom == "" {
		return MailerConfig{}, fmt.Errorf("%s is required in %s", envMailFrom, modeProd)
	}
	return mailer, nil
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
