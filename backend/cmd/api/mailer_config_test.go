package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestLoadConfig_ProductionRequiresTheResendSettings(t *testing.T) {
	cases := []struct {
		name    string
		env     map[string]string
		missing string
	}{
		{"without key", map[string]string{"MAILCHECKER_ENV": "production", "MAILCHECKER_MAIL_FROM": "Fokko <codigos@example.com>"}, "MAILCHECKER_RESEND_API_KEY"},
		{"without sender", map[string]string{"MAILCHECKER_ENV": "production", "MAILCHECKER_RESEND_API_KEY": "re_key"}, "MAILCHECKER_MAIL_FROM"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadConfig(fakeEnv(tc.env))

			if err == nil || !strings.Contains(err.Error(), tc.missing) {
				t.Errorf("LoadConfig error = %v, want it to name %s", err, tc.missing)
			}
		})
	}
}

func TestLoadConfig_ReadsTheResendSettingsInProduction(t *testing.T) {
	got, err := LoadConfig(fakeEnv(map[string]string{
		"MAILCHECKER_ENV":            "production",
		"MAILCHECKER_RESEND_API_KEY": "re_key",
		"MAILCHECKER_MAIL_FROM":      "Fokko <codigos@example.com>",
	}))

	want := MailerConfig{ResendAPIKey: "re_key", MailFrom: "Fokko <codigos@example.com>"}
	if err != nil || !got.IsProduction || got.Mailer != want || got.Mailer.UsesConsole() {
		t.Errorf("LoadConfig = %+v, %v; want production with %+v and Resend", got, err, want)
	}
}

func TestLoadConfig_DevelopmentWithoutKeyUsesTheConsole(t *testing.T) {
	for _, mode := range []string{"", "development"} {
		t.Run(mode, func(t *testing.T) {
			got, err := LoadConfig(fakeEnv(map[string]string{"MAILCHECKER_ENV": mode}))

			if err != nil || got.IsProduction || !got.Mailer.UsesConsole() {
				t.Errorf("LoadConfig = %+v, %v; want development with the console sender", got, err)
			}
		})
	}
}

func TestLoadConfig_RejectsAnUnknownMode(t *testing.T) {
	_, err := LoadConfig(fakeEnv(map[string]string{"MAILCHECKER_ENV": "prod"}))

	if err == nil || !strings.Contains(err.Error(), "MAILCHECKER_ENV") {
		t.Errorf("LoadConfig error = %v, want it to name MAILCHECKER_ENV", err)
	}
}

func TestConfig_NeverPrintsTheKey(t *testing.T) {
	config := Config{Port: 8080, IsProduction: true, Mailer: MailerConfig{ResendAPIKey: "re_very_secret", MailFrom: "codigos@example.com"}}
	var logs bytes.Buffer
	slog.New(slog.NewJSONHandler(&logs, nil)).Info("config", "config", config, "key", config.Mailer.ResendAPIKey)

	for _, printed := range []string{fmt.Sprint(config), fmt.Sprintf("%+v", config), fmt.Sprintf("%#v", config), logs.String()} {
		if strings.Contains(printed, "re_very_secret") {
			t.Errorf("printed %q contains the key", printed)
		}
	}
}
