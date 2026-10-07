# Copilot review prompt

Review the existing Chaum Work Management Backend feature stack against develop. Read AGENTS.md, docs/TASKS.md, docs/REQUIREMENTS.md and docs/schema exports first. Use actual Supabase schema as source of truth; historical migrations are not deployable instructions.

Inspect git status and preserve uncommitted work. Do not rebuild Foundation, run migrations/DDL/seed, connect tests to live Supabase, expose credentials, merge or deploy. Review public."USER", distinct worker_id/user_id, payroll.user_id, generated columns, real constraints, transactional locks, duplicate Worker conflicts, LINE server verification/audiences, DB roles and ownership.

Check scheduling/leave, QR+GPS, server Bangkok timestamps, overnight shifts, exact satang penalties/payroll, procurement unit prices and report allocation without double-counting transfers. Prioritize concurrency, unauthorized access, invalid state transitions and rollback behavior. Do not mark mocks/skips as real integration passes.

Run gofmt check, go test ./..., go vet ./..., go build ./.... Database integration requires a separately prepared disposable localhost database plus TEST_DATABASE_ISOLATED=yes. Report findings with concrete file/line, failing scenario and proposed fix. Read docs/PR_PLAN.md for stacked branch order. Do not claim PRs exist unless verified. Ask for missing deployment configuration only for affected modules.
