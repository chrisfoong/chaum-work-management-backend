# Backend verification and publication

2026-10-08: feature/delivery-validation contains the full stacked implementation. FRONTEND_HANDOFF.md maps every adopted guide use case to implemented backend APIs; frontend belongs to a separate repository. Mock Figma advance-acceptance, nine-hour shifts and GPS-or-QR are excluded.

Actual final-repo verification: gofmt has no outstanding formatting; go test -count=1 ./... with coverage passes 177 tests/subtests, zero skipped/failures, including isolated localhost PostgreSQL17; go vet ./... and go build ./... pass. Tests do not connect to production Supabase or load .env. The temporary test server is stopped after checks.

Statement coverage: total 60.2%; Auth 88.9%, Contract 86.3%, Work 59.1%. Coverage is execution of instrumented statements, not proof of every business scenario. Normal package coverage also excludes calls instrumented only in other packages. Deployed LINE audiences/provider/login, private Supabase Storage, real message delivery and frontend end-to-end remain unverified. No claim of 100% coverage or live readiness. Persistence limitations remain documented in REQUIREMENTS.md.

Publication attempted: git push --atomic -u origin for all twelve feature branches; no force push, merge or deploy. GitHub rejected with HTTP403: permission to chrisfoong/chaum-work-management-backend denied to Sxthxwit. No branch was pushed by this attempt. Resolve repository Write access or use an authorized Git identity; do not put GitHub tokens in source/.env or messages. Prepared PR descriptions are available in the task workspace BACKEND_GUIDE_DRAFT_PRS.md.

Branches: feature/schema-line-auth, feature/catalog-contracts, feature/scheduling-leave, feature/attendance, feature/procurement, feature/payroll, feature/reports, feature/runtime-api, feature/usecase-description, feature/tor-file-upload, feature/frontend-api, feature/delivery-validation. All are ancestors of the complete final branch and preserve the existing develop Foundation. Latest verification documentation is committed on feature/delivery-validation.

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
- [ ] Draft PR publication: GitHub connector returns 403 Resource not accessible by integration. GitHub browser is signed out; waiting for user login. All branch compare links are in FEATURE_PRS.md. No PR exists from this attempt.
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
- [ ] 9A persisted summary-sent/read history remains unavailable in the existing schema. Assignment continuation and actual summary transport exist, but do not claim historical notice list/no-message state.
- [ ] Live LINE/Storage/QR/GPS end-to-end and Thai PDF typography remain pending. Durable persistence facilities require separately approved design, not implicit schema changes.
- [ ] Further push held while completeness limits above remain explicit; previous 13 remote branches are unchanged. New feature is locally reviewable; no merge/deploy.

Evidence artifacts: backend-usecase-tests.jsonl and backend-usecase-coverage.out in the task workspace. Local simulation/isolated integration is not deployed acceptance. See USECASE_AUDIT.md and API.md for exact functions/payloads.
