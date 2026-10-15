package config

import (
	"fmt"
	"strings"
)

// supabaseURL resolves the public project URL, never the database DSN or service key.
func supabaseURL(get func(string) string) (string, error) {
	normalize := func(value string) string { return strings.TrimRight(strings.TrimSpace(value), "/") }
	canonical := normalize(get("SUPABASE_URL"))
	alias := normalize(get("NEXT_PUBLIC_SUPABASE_URL"))
	if canonical != "" && alias != "" && canonical != alias {
		return "", fmt.Errorf("SUPABASE_URL and NEXT_PUBLIC_SUPABASE_URL must refer to the same project URL")
	}
	if canonical != "" {
		return canonical, nil
	}
	return alias, nil
}
