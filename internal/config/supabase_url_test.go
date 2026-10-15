package config

import (
	"strings"
	"testing"
)

func TestSupabaseURLAliasesAcrossLoaders(t *testing.T) {
	for _, tc := range []struct {
		name, canonical, alias, want string
		bad                          bool
	}{
		{"canonical", " https://project.supabase.co/// ", "", "https://project.supabase.co", false},
		{"alias", "", " https://project.supabase.co/ ", "https://project.supabase.co", false},
		{"same", "https://project.supabase.co/", "https://project.supabase.co", "https://project.supabase.co", false},
		{"blank canonical", "  ", "https://project.supabase.co/", "https://project.supabase.co", false},
		{"mismatch", "https://one.supabase.co", "https://two.supabase.co", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env := map[string]string{"SUPABASE_URL": tc.canonical, "NEXT_PUBLIC_SUPABASE_URL": tc.alias, "DATABASE_URL": "postgres://fake", "LINE_LOGIN_CHANNEL_ID": "web", "LINE_WEB_CHANNEL_ID": "web", "LINE_WORKER_CHANNEL_ID": "worker", "LINE_PROVIDER_ID": "provider"}
			get := func(k string) string { return env[k] }
			runtime, e1 := LoadRuntime(get)
			legacy, e2 := Load(get)
			if tc.bad {
				if e1 == nil || e2 == nil {
					t.Fatal("mismatched project URLs accepted")
				}
				if strings.Contains(e1.Error(), tc.canonical) || strings.Contains(e1.Error(), tc.alias) {
					t.Fatal("error exposed config values")
				}
				return
			}
			if e1 != nil || e2 != nil || runtime.StorageURL != tc.want || legacy.SupabaseURL != tc.want {
				t.Fatalf("loaders disagree: %v %v", e1, e2)
			}
		})
	}
}
