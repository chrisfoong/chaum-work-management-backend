# Backend task tracker



Updated 2026-10-08. Business authority: latest SA-Group5 Use Case Descriptions adopted by user. Physical authority: Supabase exports. Existing develop Foundation retained. REQUIREMENTS.md records reconciliations; API.md documents routes.



## Implemented and verified locally



- [x] Existing schema mapping, LINE auth, ownership, TOR Foundation preserved.

- [x] 3A active contract/equipment readiness/minimum staffing/one worker-day guards.

- [x] 2W emergency/advance dates; 4A Assistant review plus atomic replacement and rollback.

- [x] 3W same-day early replacement confirmation; 4S advance-only leave exemption; emergency absence penalty; GPS/QR and overnight behavior retained.

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



## Final frontend guide  -  backend completion 2026-10-08



- [x] 1S: required uploaded private PNG <=5MiB at confirmation; verify owner/actual MIME/size before transaction; save existing contract_file_url; retain wizard and 10-digit/date/value validations.

- [x] 3A/4A: Assistant-only available-worker/leave-candidate/detail APIs; exclude original leaver, busy/inactive/unavailable/approved-leave/calculated-period workers; conflict on duplicate mappings before paging; recheck approved leave during schedule write.

- [x] 1W/9A: project/location/worker names and start/end timestamps, exact 8h overnight shifts; contract workflow_status/can_operate in lists/dashboard; existing enums unchanged.

- [x] Status/date/TOR/assignment/type filters applied in SQL before paging; preserve Worker ownership and abort missing/duplicate-worker middleware before route execution.

- [x] 2A/7A receipt photos when positive total; 2S funding retains PDF support. Purchase decision reason now available to Supervisor review.

- [x] 8A: eligible original-recipient delivery schedules; reject future dates and incomplete quantities; all purchased items handed over using guide's schedule/description/photos payload; photos/retries retain atomicity; LINE summary has actual items/quantities with bounded message chunks.

- [x] Adopt user's exclusions: no advance replacement accept/reject, no check-in photo requirement, no nine-hour mock shifts, no notification inbox, no required paid_at. Both GPS+QR retained; checkout photos and LINE Chat retained.

- [x] API/OpenAPI, requirements, handoff and README updated. Tests cover invalid filters, ownership, candidate conflicts, delivery guards/retry, overnight read model, roles, required TOR file, MIME/size/owner and notification chunk integrity.



Verification: unchanged baseline tests passed on an isolated loopback PostgreSQL17 database created for this task. Final implementation gofmt/test/vet/build passed: 177 tests/subtests, 0 failures, 0 skipped. No .env accessed, server/jobs not started, production Supabase/schema/data not changed. HTTP test servers verify Storage byte/type/size/ownership and existing LINE transport tests; deployed LINE/private Storage/Supabase end-to-end remains unverified.



Reviewable work is split into stacked local feature branches/commits for TOR file upload, frontend read APIs, and delivery validation. Draft PR publication remains blocked by previously verified GitHub write denial; no merge/deploy. See FRONTEND_HANDOFF.md for the final use-case-to-route contract.


## Pre-push audit — 2026-10-08

- [x] Rechecked all use-case API mappings in FRONTEND_HANDOFF.md against implemented routes.
- [x] Reran gofmt, go test -count=1 with coverage, go vet and go build on actual final branch. 177 tests/subtests passed; zero failures/skipped, including isolated-loopback DB tests.
- [x] Measured statement coverage: total 60.2%, Auth 88.9%, Contract 86.3%, Work 59.1%. Tests passing does not mean every branch/scenario is tested; live LINE/Storage/Supabase end-to-end remains unverified.
- [x] Attempted atomic, non-force push of all 12 feature branches to authorized origin.
- [ ] Push blocked: GitHub returned 403, permission denied to authenticated Git account Sxthxwit on chrisfoong/chaum-work-management-backend. No branches published by this attempt; no Draft PR/merge/deployment. Requires Write permission or an authorized Git identity before retry.

See VERIFICATION_AND_PUBLISH.md. No .env inspected or production data/schema changed.

## Backend completion audit — 2026-10-08

- [x] Baseline passed before changes on the existing complete branch using isolated loopback PostgreSQL17.
- [x] Added composed LINE verification + real USER lookup + Web/Worker role integration coverage. Actual LINE response is simulated, no live credentials used.
- [x] Added actual /me HTTP integration: Supervisor/Assistant/Worker, role changes without restart, wrong-role rejection, missing Worker mapping (404), duplicate mapping (409), inactive account (401).
- [x] Added GitHub Actions formatting/tests/vet/build with disposable PostgreSQL17; skipped tests fail CI; results/coverage retained seven days. Workflow YAML parsed locally.
- [x] Final local checks pass: 198 tests/subtests, zero fail/skip, gofmt, vet, build, diff whitespace check. Total statement coverage 61.0%.
- [x] Added BACKEND_ACCEPTANCE.md separating operator configuration, historical reported Supervisor login success and live checks still pending.
- [x] Corrected stale Assistant-all-areas comments, receipt rules and agent documentation paths. Local .env variants ignored; .env.example remains tracked. No .env content inspected.
- [ ] GitHub CI execution awaits publication; local YAML validation is not a hosted workflow run.
- [ ] Real Assistant/Worker LINE login, private Storage, QR/GPS and actual messaging remain unverified in this task. Shared conversation reports these configuration gaps; current values have not been inspected.

No shared Supabase schema/data changes, server startup, deployment or merge. Frontend remains in a separate repository. Persistence limits remain explicit in REQUIREMENTS.md.
