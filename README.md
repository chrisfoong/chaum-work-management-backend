# Chaum Work Management Backend

Go/Gin continued from develop, reconciled against supplied Supabase exports (16 tables, 121 columns). Never run historical migrations on the existing database.

Read [requirements](docs/REQUIREMENTS.md), [progress](docs/TASKS.md), [API notes](docs/API.md), [OpenAPI](api/openapi.yaml), and [agent rules](AGENTS.md).

## Run

Copy .env.example to .env. Configure existing database and LINE Web/Mini App channel IDs/provider. Confirm both channels belong to the same provider in LINE Developers; provider ID is a deployment declaration, not an automatic lookup. Keep secrets out of Git.

Run: `go run ./cmd/server` (root entry point uses the same application). Startup checks columns/constraints/required enums read-only and rejects incompatibility without repairing the DB. Clients send Authorization: Bearer <LINE ID token>; server verifies configured audience and resolves public."USER".line_id, rejecting unknown/inactive users.

QR_SIGNING_SECRET needs at least32 bytes. Existing private Storage bucket and server-only service key enable evidence/receipts. Missing Storage blocks dependent operations; missing messaging token explicitly disables notifications. WEB_ALLOWED_ORIGINS is a comma-separated allowlist.

## Checks

Run gofmt -w ., go test ./..., go vet ./..., go build ./....

Integration requires TEST_DATABASE_URL and TEST_DATABASE_ISOLATED=yes pointing at a disposable database. Apply internal/schema/isolated-test-schema.sql ONLY there, then go test -count=1 ./.... Never target existing Supabase. Skipped means integration unverified; fixture enum labels still need live verification.

Absent finalization runs each minute after shift end +2h. Payroll is an explicit Supervisor action after period/cutoff closure. No automatic payments occur.

Current schema lacks durable outbox, delivery history, payment timestamp/reference and historical wage/correction audit. TOR payroll allocation uses actual workdays and is estimated. PDF export and live integration remain planned.
