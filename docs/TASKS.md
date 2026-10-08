# Backend task tracker

Updated 2026-10-08. Business authority: latest SA-Group5 Use Case Descriptions adopted by user. Physical authority: Supabase exports. Existing develop Foundation retained. REQUIREMENTS.md records reconciliations; API.md documents routes.

## Implemented and verified locally

- [x] Existing schema mapping, LINE auth, ownership, TOR Foundation preserved.
- [x] 3A active contract/equipment readiness/minimum staffing/one worker-day guards.
- [x] 2W emergency/advance dates; 4A Assistant review plus atomic replacement and rollback.
- [x] 4S advance-only leave exemption; emergency absence penalty; GPS/QR and overnight behavior retained.
- [x] 5W active attendance gate; 5A TOR-wide inspection; 6A explicit purchase/no_purchase.
- [x] 1A base survey; 2S/3S funding/approval; 2A/7A partial repeated purchases, replay guard, latest/weighted prices and zero-cost rounds.
- [x] 8A matching recipient delivery/photos, transactional rollback/retry and post-delivery notification.
- [x] 5S atomic payroll batch/retries; 7W own paid slips and month filter.
- [x] 6S paid-only labor, TOR/date filters, JSON/CSV/PDF, confirmation snapshot; 9A active-contract continuation data.
- [x] Agent rules, requirements, API models, README and verification notes updated.

## Evidence

Baseline go test/vet/build passed after using writable GOCACHE; default cache had access denied. Baseline integration initially unverified.

Changed code: gofmt, go test -count=1 ./..., go vet ./..., go build ./... passed with explicitly isolated loopback PostgreSQL17. Tests cover existing auth/ownership and atomic leave rollback, emergency penalties, batch retries, active-attendance requests, decisions, simultaneous/partial purchases, latest/weighted prices, zero-cost rounds, delivery rollback/retries and paid labor. PDF parsed as three pages preserving all 25 test rows, rendered and visually inspected. 169 tests/subtests passed, zero skipped on final repository verification. No .env inspection or Supabase writes.

## Blocked / live integration unverified

- [ ] Deployed LINE/provider, private Storage and actual message delivery need end-to-end verification.
- [ ] Publish Draft PR: prior remote push rejected 403; Write access required.
- [ ] Persistent closing ledger, durable notification retries and item-round history need approved persistence design; no schema change authorized.
- [ ] PDF Thai project-name typography; current PDF identifies TOR by UUID and uses English fields.

## Unchanged-schema limits

UUIDs/real enum labels override illustrative document SQL. Payroll uses user_id; generated to_buy_qty is never written. Delivery evidence uses a description prefix without DB FK/delivery status. Closing explicitly returns persisted=false. No durable outbox, historical wage snapshots or replacement relation. Backend locks do not cover direct DB writers. None of these limitations are claimed as implemented persistent facilities.
