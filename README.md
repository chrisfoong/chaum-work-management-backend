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

Current schema lacks durable outbox, dedicated delivery history, payment timestamp/reference and historical wage/correction audit. Delivery evidence is supported through WORK_EVIDENCE descriptions. TOR payroll allocation uses actual workdays and is estimated. Basic PDF financial export is implemented; live integration remains unverified.

## Latest Use Case Description behavior

Business authority is SA-Group5 Use Case Description; physical schema remains exported Supabase metadata. See docs/REQUIREMENTS.md, docs/API.md and docs/TASKS.md. Backend only; unchanged database. Assistant reviews leave (same-day emergency supported), procurement can span rounds, delivery is recorded in existing work_evidence, payroll supports atomic batches and profit includes paid payroll only. Historical migrations must not run on Supabase. Never inspect .env. Integration requires isolated loopback PostgreSQL, not the shared Supabase project.

API changes include expected_actual_qty on each purchase item and Assistant ownership of leave review. Financial closing confirmation is a snapshot (persisted=false). Notifications have no durable outbox. PDF financial export uses TOR identifiers/English fields; Thai names remain in JSON/CSV.

## Final frontend guide coverage

Supabase remains the application database. Added private PNG contract confirmation, Assistant replacement candidates/leave details, filtered lists with joined project/shift data, contract workflow status, and eligible delivery schedules with transactional date/quantity checks. See docs/FRONTEND_HANDOFF.md for the workflow-to-API mapping. No advance replacement acceptance, check-in photo, nine-hour shift, notification inbox or payment timestamp is required by the adopted guide. No production schema or data is changed by this implementation.

## Automated and live acceptance

GitHub Actions in .github/workflows/backend.yml runs the full isolated PostgreSQL17 suite and rejects skipped tests. No Supabase secrets are used. See [backend acceptance](docs/BACKEND_ACCEPTANCE.md) for Web/Worker login diagnosis, configuration ownership and pending live checks. Historical shared evidence reports Supervisor login success; Worker/QR/Storage/message delivery remain separately unverified.
