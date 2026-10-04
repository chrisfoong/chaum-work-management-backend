// Package config loads runtime settings from environment variables.
package config

import (
	"fmt"
	"strings"
)

// Config holds runtime settings. Values are never logged.
type Config struct {
	Port          string
	DatabaseURL   string
	SupabaseURL   string // https://<project-ref>.supabase.co, no trailing slash
	LineChannelID string // LINE Login channel ID that owns the LIFF app
}

// Load reads the configuration through getenv (os.Getenv in production).
// It reports every missing required variable by name, never by value.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:          getenv("PORT"),
		DatabaseURL:   getenv("DATABASE_URL"),
		SupabaseURL:   strings.TrimRight(getenv("SUPABASE_URL"), "/"),
		LineChannelID: getenv("LINE_LOGIN_CHANNEL_ID"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	required := []struct {
		name  string
		value string
	}{
		{"DATABASE_URL", cfg.DatabaseURL},
		{"SUPABASE_URL", cfg.SupabaseURL},
		{"LINE_LOGIN_CHANNEL_ID", cfg.LineChannelID},
	}
	var missing []string
	for _, r := range required {
		if r.value == "" {
			missing = append(missing, r.name)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}
