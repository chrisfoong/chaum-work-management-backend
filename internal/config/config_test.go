package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	full := map[string]string{
		"DATABASE_URL":          "postgres://example",
		"SUPABASE_URL":          "https://ref.supabase.co/",
		"LINE_LOGIN_CHANNEL_ID": "123",
	}

	tests := []struct {
		name        string
		env         map[string]string
		wantErr     string
		wantPort    string
		wantSupaURL string
	}{
		{name: "all set, default port", env: full, wantPort: "8080", wantSupaURL: "https://ref.supabase.co"},
		{name: "explicit port", env: with(full, "PORT", "9000"), wantPort: "9000", wantSupaURL: "https://ref.supabase.co"},
		{name: "missing one", env: with(full, "DATABASE_URL", ""), wantErr: "DATABASE_URL"},
		{name: "missing all", env: map[string]string{}, wantErr: "DATABASE_URL, SUPABASE_URL, LINE_LOGIN_CHANNEL_ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load(func(k string) string { return tt.env[k] })
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Port != tt.wantPort || cfg.SupabaseURL != tt.wantSupaURL {
				t.Fatalf("got port %q url %q", cfg.Port, cfg.SupabaseURL)
			}
		})
	}
}

func with(base map[string]string, k, v string) map[string]string {
	out := make(map[string]string, len(base)+1)
	for bk, bv := range base {
		out[bk] = bv
	}
	out[k] = v
	return out
}
