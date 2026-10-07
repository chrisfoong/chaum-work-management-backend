package config

import (
	"fmt"
	"strings"
)

// Runtime is deployment configuration for the confirmed MVP.
type Runtime struct {
	Port, DatabaseURL, WebChannel, WorkerChannel, ProviderID        string
	QRSecret, StorageURL, StorageKey, StorageBucket, MessagingToken string
}

func LoadRuntime(get func(string) string) (Runtime, error) {
	c := Runtime{Port: get("PORT"), DatabaseURL: get("DATABASE_URL"), WebChannel: get("LINE_WEB_CHANNEL_ID"), WorkerChannel: get("LINE_WORKER_CHANNEL_ID"), ProviderID: get("LINE_PROVIDER_ID"), QRSecret: get("QR_SIGNING_SECRET"), StorageURL: strings.TrimRight(get("SUPABASE_URL"), "/"), StorageKey: get("SUPABASE_SERVICE_ROLE_KEY"), StorageBucket: get("STORAGE_BUCKET"), MessagingToken: get("LINE_MESSAGING_TOKEN")}
	if c.Port == "" {
		c.Port = "8080"
	}
	var missing []string
	for _, r := range []struct{ name, value string }{{"DATABASE_URL", c.DatabaseURL}, {"LINE_WEB_CHANNEL_ID", c.WebChannel}, {"LINE_WORKER_CHANNEL_ID", c.WorkerChannel}, {"LINE_PROVIDER_ID", c.ProviderID}} {
		if strings.TrimSpace(r.value) == "" {
			missing = append(missing, r.name)
		}
	}
	if len(missing) > 0 {
		return Runtime{}, fmt.Errorf("missing configuration: %s", strings.Join(missing, ", "))
	}
	if c.QRSecret != "" && len(c.QRSecret) < 32 {
		return Runtime{}, fmt.Errorf("QR_SIGNING_SECRET must contain at least 32 bytes")
	}
	return c, nil
}
