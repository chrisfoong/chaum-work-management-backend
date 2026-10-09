# Backend verification

Run gofmt, go test -count=1 ./..., go vet ./..., go build ./.... Integration requires TEST_DATABASE_URL to isolated loopback PostgreSQL plus TEST_DATABASE_ISOLATED=yes. Never inspect .env or write integration data to shared Supabase. Skipped means unverified.

Checks cover Assistant leave role, atomic replacement rollback, emergency/advance absence, purchase replay/concurrency and zero-cost policy, latest/weighted price, delivery multi-photo rollback/retry, batch payroll reuse, paid-only labor, PDF parsing/rendering, plus existing LINE audience/ownership/schema tests.

Live LINE/Storage/message delivery and remote PR publication are separate unverified steps. See TASKS.md for evidence and limits.
