package contract

import (
	"chrisfoong/chaum-work-management-backend/internal/schema"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("isolated PostgreSQL not configured; integration unverified")
	}
	if os.Getenv("TEST_DATABASE_ISOLATED") != "yes" {
		t.Fatal("TEST_DATABASE_ISOLATED=yes required")
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil {
		t.Fatal("invalid test configuration")
	}
	switch cfg.ConnConfig.Host {
	case "localhost", "127.0.0.1", "::1":
	default:
		t.Fatal("local isolated database required")
	}
	pool, e := pgxpool.NewWithConfig(context.Background(), cfg)
	if e != nil {
		t.Fatal("test database unavailable")
	}
	t.Cleanup(pool.Close)
	if e = schema.Check(context.Background(), pool); e != nil {
		t.Fatal(e)
	}
	return pool
}
