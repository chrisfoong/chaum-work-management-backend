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

## Publication update — 2026-10-08

- [x] GitHub Write permission is now available; atomic non-force push of all 13 feature branches succeeded. Prior 403 push blockers above are historical and resolved.
- [x] First hosted CI exposed an existing timezone-dependent overnight assertion (UTC representation versus Bangkok date). Fixed the assertion to parse instants, convert to Bangkok and verify both expected local dates/times plus exactly eight hours.
- [x] Full isolated local test run repeated with PostgreSQL session timezone=UTC: 198 passed, zero failed/skipped. Runtime shift calculation unchanged.
- [ ] Draft PR publication: GitHub connector returns 403 Resource not accessible by integration. GitHub browser is signed out; waiting for user login. All branch compare links are in historical/pr/FEATURE_PRS.md. No PR exists from this attempt.
- [ ] Hosted CI rerun must be checked after this correction; do not treat the first failed run as passing.

Auto-review rejected a proposed Git credential-helper extraction for API publication; that script was not executed and was removed. No credential was extracted/displayed. Existing Git push succeeded through its normal credential handling. No .env file inspected; no shared Supabase data/schema changed.

## Detailed Use Case audit — feature/usecase-completion (2026-10-08)

- [x] Read substantive SA tab descriptions for all 22 use cases; requested 5.2 QUERY TABLE tab is reference SQL, not the descriptions. Added USECASE_AUDIT.md mapping functions/routes/conditions/limits.
- [x] Preserved existing Foundation and exported Supabase schema; inspected changed columns/types against metadata. No .env access, DDL, shared Supabase writes or server/jobs startup.
- [x] Added schedule address/stored future status; shortage reason <=500; pending additional inspection and additional-only approval guards; correct no_purchase reviewer/time and inactive-equipment decline.
- [x] Corrected 2A zero acquisition versus 7A positive quantity/initial approval/follow-up funding; added shared-validation purchase preview with no business writes and weighted-price review.
- [x] Added labor/material report detail arrays that reconcile to exact totals, selected-assignment continuation flags, and new-payroll project/area summary notification.
- [x] Added guarded notification-only resend, decision/delivery business_saved + delivery status, bounded transport retries and stable LINE retry keys. No expense/quantity/evidence replay; no durable outbox claimed.
- [x] Baseline 198 tests/subtests passed before edits. Final gofmt, go test -count=1 ./..., go vet ./..., go build ./..., diff whitespace passed with explicitly isolated loopback PostgreSQL17 timezone=UTC: 219 passed, 0 failed, 0 skipped; statement coverage 62.2%.
- [x] Confirmed parent hosted GitHub CI succeeds: run 37799511198 at 90f6660. This is parent evidence; new feature hosted CI is unverified until publication.
- [x] Superseded by user-confirmed 9A MVP: persisted summary/send/read history is outside scope, not a blocker. Current reports and manual notices require no new tables.
- [ ] Live LINE/Storage/QR/GPS end-to-end and Thai PDF typography remain pending. Durable persistence facilities require separately approved design, not implicit schema changes.
- [ ] Further push held while completeness limits above remain explicit; previous 13 remote branches are unchanged. New feature is locally reviewable; no merge/deploy.

Evidence artifacts: backend-usecase-tests.jsonl and backend-usecase-coverage.out in the task workspace. Local simulation/isolated integration is not deployed acceptance. See USECASE_AUDIT.md and API.md for exact functions/payloads.

## 9A design / live checks — 2026-10-08

- [x] Proposed two-table snapshot/recipient queue design with Assistant-safe content and optional open/ack tracking in 9A_DESIGN_AND_LIVE_TEST_PLAN.md. No DDL implemented/executed. Read-receipt history is optional; 9A does not explicitly require it and its absence alone should not block all publication.
- [x] Existing localhost server returned backend JSON for readiness 200 and four missing/invalid-token Web/Worker checks 401. No .env inspection, server start or business writes. This does not establish successful LINE login or intended DB/code revision.
- [ ] Real Web/Worker login awaits nonsecret login/LIFF URLs and human login. Mutation flows await an explicitly isolated test environment; no shared Supabase writes authorized.

## Confirmed 9A MVP — feature/assistant-operations-summary

- [x] User superseded the persistence proposal: only 16 existing tables, live operational report, no financials/snapshots/send/read/ack history. Earlier automatic Assistant payroll summary is superseded.
- [x] Added OperationsSummary and NotifyOperationsSummary; retained all-area Assistant scope and backend contract scheduling guards.
- [x] Removed automatic Assistant summary from PayrollBatch/CloseSummary; existing Worker payslip notices remain.
- [x] Added manual Supervisor review/closed-period/source-data/payroll/HTTPS guards, short LINE project/date Dashboard link and honest transport status.
- [x] Updated API, requirements, frontend handoff and acceptance docs; no .env inspection, server startup or shared database/schema writes.
- [x] Regenerated OpenAPI. Full gofmt, go test -count=1 ./..., go vet ./..., go build ./... and diff whitespace checks passed on explicitly isolated loopback PostgreSQL17 timezone=UTC: 233 tests/subtests passed, 0 failed, 0 skipped; statement coverage 63.5%. Evidence: backend-9a-tests.jsonl and backend-9a-coverage.out in task workspace.
- [ ] Real successful LINE login, private Storage, GPS/QR, actual manual LINE delivery and Dashboard navigation remain live-unverified. No new tables needed for 9A; frontend work stays in its own repo.

## Real LINE role checks — 2026-10-08

- [x] User authorized shared SA_DB role changes for USER 11111111-1111-1111-1111-111111111111 only. Switched supervisor -> assistant -> worker -> supervisor; each change returned one active row; restored supervisor verified. No Worker mapping or business data added, no .env access, server startup/restart or schema edits.
- [x] Real Web LINE /me 200 for target Supervisor, then Assistant with same verified token and DB-resolved role. Assistant payroll/locations 403; missing/invalid tokens and wrong Worker audience 401.
- [ ] Live read suites each 14 pass/2 fail: assignments and workers/available return 404 non-JSON from current process/gateway despite routes existing in latest local code. Running revision unverified; not counted as complete live acceptance.
- [ ] Worker desktop helper login encounters LINE 400 before Backend; direct Worker LIFF opens smartphone QR. No valid Worker token obtained; worker_rows=0. Successful Worker login/reads remain unverified, not passed.

Detailed sanitized evidence and helper changes: line-login-test/LIVE_RESULTS_2026-10-08.md and ROLE_TEST_PLAN.md in task workspace. Historical live-pending entries above now superseded for Supervisor/Assistant login only; Storage/GPS/QR/messages and Worker remain pending.
## All-role live recheck — 2026-10-09

- [x] Supervisor and Assistant read-only LINE suites each 16 pass, 0 fail on user-restarted backend. Target /me identity verified. Previously missing assignments/workers-available routes now return expected results, including Assistant 200 versus Supervisor 403 for available workers.
- [x] Temporarily switched the same authorized USER role to assistant/worker and restored supervisor; each scoped UPDATE returned one active row. No Worker mapping added; precheck count=0. No .env, schema or business workflow changes.
- [x] Worker role correctly denied Web API 403 with real Web token. A clean Worker login attempt still reaches LINE 400 before Backend; not counted as successful Worker authentication.
- [x] Read-only LINE Console identified Worker Developing Endpoint still default LINE page and openid unchecked. Scope editor cancelled without save. Requested confirmation for endpoint/openid change and separate Worker mapping preparation.
- [ ] Successful real Worker login/read suite awaits approved LINE configuration and Worker mapping; simulated tests are separate evidence. Storage/GPS/QR/messages/write workflows remain unverified.

Evidence: line-login-test/LIVE_RESULTS_2026-10-09.md in task workspace. Earlier 404 failures are historical and resolved for Web roles. No server was started/restarted by this check.
## Final Use Case verification and live cleanup — 2026-10-09

- [x] Freshly exported SA-Group5 source; rechecked all 22 substantive Use Case Descriptions against routes/implementation/tests and confirmed user overrides. No missing Backend flow identified; external acceptance limits remain explicit. See FINAL_USECASE_REVIEW.md.
- [x] Real LINE read suites completed: Supervisor 16/16, Assistant 16/16, Worker 11/11. Endpoint/openid issues and Worker mapping absence above are historical, resolved for the live Worker test.
- [x] User authorized adding exactly one Worker mapping, then deleting it and restoring Supervisor. Cleanup checked related Worker FK records and removed only worker_id fb3629c2-65cd-48dc-b72f-c802a0ef7bda. Verified target USER active=true, role=supervisor, worker_rows=0. This bounded shared-data exception is not authorization for other mutation tests.
- [x] Fresh gofmt/test/vet/build/diff-check passed on final branch: 233 tests/subtests, 0 failed, 0 skipped, 16 packages. Statement coverage 63.1%. Isolated PostgreSQL17 database chaum_final_audit_20261009 at loopback:55439 only; schema initialization there, never on Supabase.
- [ ] 5S wording reconciliation in USECASE_AUDIT.md deferred by user on 2026-10-09. Reverted only that wording edit; existing manual 9A runtime behavior and verification evidence remain unchanged.
- [ ] Live private Storage/upload, physical GPS+QR check-in/out, actual LINE message delivery and Assistant dashboard navigation remain unverified. Empty-list reads do not validate populated business workflows.
- [ ] Thai PDF typography remains limited; unchanged-schema persistence limits remain documented, with 9A history explicitly outside MVP.

Evidence: task workspace backend-final-tests.jsonl, backend-final-coverage.out and line-login-test/LIVE_RESULTS_2026-10-09.md. No .env access, Supabase DDL/migrations/seeds, normal server/jobs startup, new push/PR/merge/deploy.

## Current handoff — 2026-10-10
- [x] Paired startup/acceptance guide: SERVER_AND_TEST_STATUS.md (Backend 8080, Frontend 5500, LINE tunnel).
- [x] Latest rerun: 235 tests/subtests, 0 failures/test skips on isolated PostgreSQL; vet/build passed. Live LINE roles and QR issuance/expiry verified; shared test mapping deleted and Supervisor restored.
- [ ] Populated procurement/funding/delivery, Storage transfer, physical GPS+camera attendance/checkout, payroll/payment/export and actual LINE recipient delivery remain live-unverified. Historical pending-login notes are superseded by this handoff; mock tests are separate evidence.

## Dev / Docker — 2026-10-10

- เพิ่ม Dockerfile แบบ multi-stage, non-root runtime, Compose healthcheck และคู่มือ DEV_DOCKER.md
- Docker build ผ่าน; container backend health HTTP 200 บนพอร์ตทดสอบ 8081
- ใช้ config Supabase เดิม ไม่มี database container/migration/reset; ผล build/health ไม่ใช่ผลทดสอบครบทุก use case
- ตรวจซ้ำ go test ./..., go vet ./..., go build ./... ผ่าน; integration ที่ skip ไม่ถือว่าผ่านจริง

## Live export / 9A check — 2026-10-10

- Profit PDF/CSV ดาวน์โหลดจาก UI จริง เปิดไฟล์ ตรวจ TOR/ช่วงวันที่/ยอด 0, 0, 300, -300 ตรงกับรายงานรอบ 1–10 ต.ค.
- 9A manual POST รอบ TEST 16–30 ก.ย. HTTP 200 แต่ notification failed ไม่ใช่ delivery success; active Assistant 1 บัญชีมี LINE ID ผิดรูปแบบ
- ผู้ใช้มีเฉพาะ LINE บัญชีเดิมซึ่งคืน Supervisor แล้ว จึงยังไม่มีผู้รับ Assistant จริงแยก; ไม่ข้าม authorization หรือแก้ recipient role เพื่อให้ผลผ่าน
- Physical phone TEST shift เช็กอิน 17:10:10 / เช็กเอาท์ 17:12:35 เวลาไทย, completed; checkout evidence PNG อ่านจาก Storage HTTP 200

## 2026-10-10 Dashboard and acceptance audit

- Dashboard only counts active contracts within their date range, returns contract numbers and filters terminal requisitions before selecting pending work.
- Added integration regression covering future/expired/registered contract counts and completed requests not hiding pending work.
- 242 tests/subtests passed, no failures or skipped tests, on isolated localhost PostgreSQL; vet/build passed.
- Live exports and physical check-in/out are documented separately from mocked transport tests. See [UI_ACCEPTANCE_20261010.md](UI_ACCEPTANCE_20261010.md) for all 22 use cases and outstanding live verification.
- No schema changes, no authentication bypass, no assertion that LINE acceptance means real recipient delivery.

## 2026-10-14: LINE invitation onboarding (no ER change)

- Replaced the prepared registration-request approach with Supervisor-issued invitations using existing USER/WORKER columns only; no Supabase DDL or new table.
- Added creation/renewal, verified-LINE preview/claim, 24-hour expiry, one-time consumption and inactive-account protection. Pending users have no business access or Worker mapping.
- Frontend creates/copies LIFF links, preserves invitation fragments through redirects, removes manual LINE subject entry, and displays pending LINE binding separately from suspension.
- Verified backend go test ./... against isolated loopback PostgreSQL, go vet/build; frontend lint/typecheck/format, 35 tests; production Docker builds and both services healthy (5501/8081).
- Preview browser check passed on desktop/mobile: three roles and no LINE subject field. This is visual evidence only. A real separate LINE recipient still needs to complete the mobile claim flow before claiming live end-to-end acceptance.
- Full instructions: docs/ONBOARDING.md. No push/merge performed.

## 2026-10-14: Cancel LINE invitations

Supervisor can cancel the current unclaimed invitation without deleting personnel or changing ER. Cancelled accounts can be re-invited; linked accounts cannot be deactivated through cancellation. Added transaction/row-lock race coverage and UI confirmation/status labels. Tests run on isolated loopback PostgreSQL, not shared Supabase.


## New TOR location coordinates — 2026-10-14

Extended new_location with optional paired latitude/longitude for backward compatibility. Normalization retains them; scope validates finite geographic bounds; location insertion stores both atomically with TOR in existing columns. No DDL/ER changes. Full go test ./... with isolated loopback PostgreSQL, go vet and go build passed. New validation and isolated persistence regression tests passed without skips. No live Supabase TOR was created for this change.


## TOR-scoped Worker equipment read — 2026-10-14

Implemented assignment_id filter, ownership guard and DISTINCT-by-EXISTS initial TOR equipment read. Full isolated Backend tests, go vet and build passed; new regression covers exclusion of unrelated equipment, ownership denial, invalid assignment, pagination and inactive equipment. No Supabase mutation test or schema change.


## Verified onboarding platform routing — 2026-10-14

Registered onboarding responses include platform derived from current database role (Worker -> worker, Supervisor/Assistant -> web). Inactive/unknown responses do not supply an authorized destination. Identity middleware and application role/audience guards remain unchanged. Full isolated tests, go vet/build and Docker production build passed; regression checks all three roles and inactive users. No schema change or Supabase mutation.


## Early checkout defect — 2026-10-14

Live read confirmed requested Worker checked in at 20:07 and checked out at 20:39 for 18:00-02:00 shift; historical record preserved as evidence. Added scheduled-end server guard and retained unclosed overnight read. Full isolated tests/vet/build passed; regression rejects 23:00 and one nanosecond before 04:00 for 20:00-04:00 shift, verifies no evidence/status/checkout writes on rejection, and accepts exact end. No shared DB mutation. User reports TOR equipment filter, real GPS/QR and additional request test exercised; live request is pending_procurement. New real-mobile checkout policy still requires verification.


## Supabase URL aliases — 2026-10-15

Both backend config loaders now resolve canonical SUPABASE_URL and fallback NEXT_PUBLIC_SUPABASE_URL through one helper. Trim surrounding whitespace and trailing slashes, accept matching aliases, and reject conflicting project URLs with variable names only. Service role key remains backend-only; frontend has no direct Supabase client and needs neither URL. Full isolated tests, vet/build and Docker build passed. Tests cover both loaders, alias-only, empty canonical, matching values and mismatches. No private env values modified or logged.


## Workflow route refactor — 2026-10-28

`internal/work/handler.go` now wires the existing role groups to domain registration helpers: shared reads, catalog, scheduling/attendance, procurement, and reporting/payroll. All 62 source route registrations and original handler bodies were moved without changing HTTP methods, paths, role middleware, response payloads, services or SQL.

Verification: generated OpenAPI before/after the refactor is identical (70 paths, 81 operations). Existing authorization and full isolated-loopback database tests, gofmt, go vet and go build passed. No Supabase/schema/config changes were made for this refactor.
