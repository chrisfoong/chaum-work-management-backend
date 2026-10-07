# Backend task tracker

Updated 2026-10-08. Source: user agreement, exported metadata, existing code/docs. Continue develop Foundation. Old tracker: historical/TASKS-2026-10-05.md.

## Implemented and verified locally

- [x] public."USER", distinct worker_id/user_id, payroll.user_id, actual constraint names and read-only generated quantity.
- [x] Read-only startup schema/enum compatibility checks; no live DDL/migrations.
- [x] LINE Web/Mini App audiences/issuer/expiry, DB roles, unknown/inactive rejection and duplicate Worker conflicts.
- [x] Existing TOR preview/confirm, numbering concurrency and application name guards.
- [x] User/Worker/catalog/contract management and ownership.
- [x] Assistant schedules,24h leave notice, Supervisor review, unpaid approved leave/no absence penalty, replacement.
- [x] QR+GPS200m/accuracy50m, eight-hour overnight shifts, evidence checkout, rerunnable absent finalization.
- [x] Requisition survey/review/funding/purchase, one purchase/item, actual unit prices, separate material/transfer accounting.
- [x] Half-month payroll, exact satang, penalty tiers/cap, concurrency/overlap guards and payslips.
- [x] Invoice receipts, estimated TOR payroll allocation, profit JSON/CSV and dashboards.
- [x] Private Storage proxy, evidence validation, best-effort LINE notifications with explicit failures.
- [x] Shared entry point, safe logging/errors, request IDs, explicit CORS, timeouts/graceful shutdown.
- [x] Current agent rules, requirements, historical docs, placeholders and API documentation.

## Evidence

Baseline develop build passed; test/vet failed on missing testPool and LockRequisitionNumbering. These preexisting failures are repaired.

Final prepared stack: gofmt, go test -count=1 ./..., go vet ./..., go build ./... passed 2026-10-08. Database integration ran on fresh loopback PostgreSQL17 with isolated fixture: concurrent check-in/payroll, evidence rollback, leave/replacement/absence, procurement retries, invoices/profit and duplicate Worker conflicts. LINE/Storage tests use local HTTP fixtures.

User authorized local commits, feature pushes and Draft PRs on 2026-10-08; no merge/deployment. Publication in progress. Stacked features must be reviewed/merged in dependency order.

## Blocked / integration unverified

- [ ] Live Supabase schema/enum check: no live credentials used; fixture labels are expectations.
- [ ] Real LINE tokens and same-provider membership: confirm channels in LINE Developers.
- [ ] Existing private Storage bucket/policies and real evidence roundtrip.
- [ ] Real messaging delivery/recipient eligibility: API acceptance does not prove delivery.

## Planned / unchanged-schema limits

- [ ] PDF export, richer list filters and frontend end-to-end testing.
- [ ] Durable retry/outbox, wage snapshots and correction audit need an approved persistence design.
- [ ] Delivery history and payment timestamp/reference unavailable in schema.

Application locks cover backend writers, not direct DB writers. Replacement retains original schedule without a durable replacement relation. Uploads can leave unused objects. Funding references/invoice receipts allow matching-payload retries; repeat purchases conflict. No automatic financial transfers.
