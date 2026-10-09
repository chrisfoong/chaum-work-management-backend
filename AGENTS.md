# Chaum Backend agent instructions

Continue the existing Go/Gin backend. Read docs/TASKS.md, docs/REQUIREMENTS.md and docs/API.md. Business authority is the latest SA-Group5 Use Case Description adopted by the user on 2026-10-08; it supersedes previous Supervisor leave, 24h notice and single-purchase policies. Backend only; preserve the existing Foundation.

Never open, read, inspect, copy or modify .env files, directly or via scripts/diagnostics/other agents. Never log credentials, tokens, DSNs or bank details. .env.example is placeholders only. Normal server config loading requires separate user authorization and does not authorize secret inspection.

Physical authority: docs/schema/columns.json and constraints.csv. Never run migrations, DDL, AutoMigrate, seeds or destructive SQL on Supabase. Migrations are historical. Quote public."USER". Distinct worker_id/user_id; payroll.user_id. Detect duplicate Worker/User mappings, never choose an arbitrary row. Generated to_buy_qty is read-only. Parameterize SQL, transactionally lock related writes, use exact decimal or integer satang.

All roles use verified LINE subject resolved by USER.line_id; reject unknown/inactive users and wrong audience. Roles/ownership from DB. Web supervisor/assistant, Mini App worker. Assistant handles all assignments and reviews leave; procurement approval/funding, payroll/payments and profit remain Supervisor-only.

Follow REQUIREMENTS.md for same-day emergency leave, advance-only absence exemption, atomic review/replacement, equipment readiness/minimum staffing, partial repeated procurement, 2A latest-price vs 7A weighted-price, actual delivery evidence/notification and paid-only profit. Retain eight-hour overnight shifts, GPS+QR, Bangkok timestamps and existing absence cutoff where descriptions do not specify it.

Delivery records use a documented WORK_EVIDENCE description convention because no requisition FK exists. Closing confirmation is a snapshot with persisted=false, never terminate a contract to imitate closing. Do not claim durable outbox, item-round history, wage snapshots or replacement relations absent from schema.

Use feature/<name> branches; preserve uncommitted work. User authorized local feature commits/push/Draft PRs, never merge/deploy. Run gofmt, go test ./..., go vet ./..., go build ./.... Integration writes require explicitly isolated loopback PostgreSQL. Skipped means unverified. Keep docs/TASKS.md honest.

Latest frontend guide: exclude advance replacement acceptance; confirm today's shift with both GPS/QR and exactly eight hours. Checkout photos, LINE Chat equipment results, actual schema-based payslips, contract status dashboard. TOR PNG <=5MiB is verified server-side. Purchase receipt photos for base/additional, funding may include PDF. Keep candidate/read-model filters and ownership in SQL. No mock Figma statuses/times are requirements.

9A MVP confirmed by user: operational current-data report only, no financials for Assistant, no persisted snapshot/send/read history or new tables. Supervisor manually sends after reviewing 5S/6S; no automatic Assistant summary notice from payroll or profit confirmation. Report absence means no source records for selected dates, not no notification. Assistant still manages all areas under prior agreement. Frontend uses can_continue; Backend scheduling guards remain authoritative.
