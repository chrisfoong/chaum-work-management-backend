package config

import (
	"strings"
	"testing"
)

func TestRuntimeFailClosed(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://fake", "LINE_WEB_CHANNEL_ID": "web", "LINE_WORKER_CHANNEL_ID": "worker", "LINE_PROVIDER_ID": "provider"}
	get := func(k string) string { return env[k] }
	if _, e := LoadRuntime(get); e != nil {
		t.Fatal(e)
	}
	delete(env, "LINE_WEB_CHANNEL_ID")
	if _, e := LoadRuntime(get); e == nil || !strings.Contains(e.Error(), "LINE_WEB_CHANNEL_ID") {
		t.Fatal("missing web audience accepted")
	}
	env["LINE_WEB_CHANNEL_ID"] = "web"
	env["QR_SIGNING_SECRET"] = "short"
	if _, e := LoadRuntime(get); e == nil {
		t.Fatal("weak QR secret accepted")
	}
}
