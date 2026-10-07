# TASKS.md — project state tracker

**Read this first in every session. Update it at the end of every task.**
Nothing is committed in this project, so git history does NOT show progress — this file is the state. Keep it truthful: only tick a box for work you actually verified, and say what you ran.

Status legend: `[ ]` not started · `[~]` in progress · `[x]` done and verified · `[!]` blocked (name the D#) · `[-]` skipped by user decision

You may edit this file without a separate plan (it is bookkeeping), but show the diff in your report. Do not change task scope or order without the user's approval. Never tick a box based on files merely existing.

---

## Current state (overwrite each session)
- Last updated: 2026-10-05
- Active phase: P3 foundation done; next is P4
- Repo: **rebuilt from scratch at user request** (old Gin skeleton deleted: main.go, routes/, handlers/, cmd/api/; the committed parts can be recovered from commit cf538a6)
- Schema: designed from docs/02 in migrations/0001_init.*.sql — **NOT applied, not syntax-checked against PostgreSQL**
- Uncommitted changes present: yes (tracked main.go/routes/handlers deleted; go.mod/go.sum/.gitignore modified; cmd/, internal/, migrations/, .env.example, docs/, README.md untracked)
- Build / vet / test: 2026-10-05, go1.26.0 — build exit 0, vet exit 0, `go test -count=1 ./...` all 9 packages ok; gofmt clean
- Next action: **1S built — stopped for user review** (P4 runs 1S → 1A → 9A with a stop after each). Database option (B): no database; P4 SQL **unverified**; integration tests skip unless TEST_DATABASE_URL (local only) is set.

## Phases and tasks

### P0 — Verify existing repo (read-only) — Prompt 0
- [~] Rules and decisions restated back to user and confirmed (restated 2026-10-05; user went ahead without correcting it, but never confirmed explicitly)
- [x] docs/07 §A checklist run; table of expected / found / match recorded under "Findings" (schema items 4–7 unverifiable: no schema in repo)
- [x] Baseline `go build ./...`, `go vet ./...`, `go test ./...` recorded
- [x] What the user already initialised is listed (structure, packages, config, routes, migrations, tests)
- [x] Contradictions with docs listed

### P1 — Schema reconciliation (read-only) — Prompt 1
No schema existed, so reconciliation was replaced by schema design (user decision 2026-10-05).
- [~] migrations/0001_init.up.sql / .down.sql written. **Not applied, not syntax-checked** (local psql/postgres binaries exist; needs approval)
- [x] Static review of the migration against docs/02 (2026-10-05, read-only; see Findings "P1 migration review"):
  - [x] Entity map: 16 entities → table, PK, FKs, unique constraints, types
  - [x] Canonical names: all docs/02 §1 attributes present and spelled exactly (incl. exat_deduction_amount); no total_wage / net_pay / payroll status / expense_id on requisition
  - [x] Identity: one key (`worker.user_id` PK + FK → users); no worker surrogate key
  - [x] Status value sets recorded; open sets flagged (contract_tor, shift_status, company_invoice)
  - [x] Constraint gaps recorded (leave per date, FK indexes, R16 1:1, RLS, payroll re-run guard D8)
  - [x] Relationships R01–R24 + U1 mapped; R24 not implemented (D1)
- [x] Review fixes round 1 (2026-10-05, approved): H1 RLS on all 16 tables (no policies), L2 net_wage CHECK (TODO(decision-10)), L3 fund_transfer needs transfer_ref_no, L1 12 FK indexes (each justified by an SD query). Verified statically: 16 tables = 16 RLS statements, every index column exists, no other lines changed. **Still not executed**
- [x] M4, M5, L4, L6 recorded as service-layer notes (below)
- [ ] Waiting on user: M1 (D8), M2 (R16 / replacement), M3 (duplicate leave), L5, L7, L8, requisition status `rejected`
- [x] Migration round 2 (data dictionary alignment): approved and applied to the FILE 2026-10-05 (not to any database). Static check: 16 tables = 16 RLS statements; all REFERENCES hit a primary key; 11 indexes on existing columns; every datadict column present in all 16 tables; only extra column `reject_reason` (deliberate); named unique constraints uq_contract_tor_contract_no, uq_location_location_name, uq_equipment_equipment_name, uq_equipment_requisition_requisition_no for 409 mapping. **Not executed / not syntax-checked by PostgreSQL.** Down file unchanged (same table names). Covers role values unchanged (`assistant`, per user 2026-10-05), requisition status set, `requested_by`, payroll `user_id` + `managed_by_id`, `worker.worker_id` PK, leave → users + UNIQUE (user_id, leave_date), attendance `substitute_worker_id` + UNIQUE (schedule_id), expense_claim rules, `line_id` NOT NULL, billing_month 'YYYY-MM', status sets for contract/invoice/shift, DD types and NOT NULLs. Resolves M2, M3, L5 (billing_month), L7 (attendance_id nullable) once applied

#### Service-layer notes (from P1 review; no SQL)
- M4: create-user with role `worker` must insert the `worker` row in the same transaction; never create a `worker` row for another role.
- M5: `users.line_id` stores the LINE **userId** = `sub` returned by ID-token verification (not the public LINE ID). Workers and supervisors may both have one (LINE notifications).
- L4: services set `updated_at = now()` on every UPDATE of `users` and `contract_tor` (no trigger).
- L6: when writing attendance, copy `work_date` from the referenced `work_schedule`; never take it from client input.

#### Supabase setup notes (from H1; Supabase docs checked 2026-10-05)
- Backend `DATABASE_URL` connects as `postgres` (table owner, BYPASSRLS). Never give clients the DB connection or the service_role key.
- Table RLS does not cover Storage. Use **private** buckets for receipt photos and work evidence (public buckets are readable by anyone with the URL). The backend uploads with the service_role key (server-only), stores only the object path, and serves files through signed URLs.
- Do not add RLS policies or expose tables through the Data API without a decision.
- [-] docs/decisions/schema-mapping.md: not created (the user did not ask for it; findings live in this file)

### P2 — Architecture and skeleton plan (plan only) — Prompt 2
- [x] Layout proposed and approved (rebuild plan, 2026-10-05): cmd/server, internal/{config,db,apperr,httpx,auth,notify,storage,health,server}, migrations/. Workflow modules are added per slice
- [ ] Workflow matrix (use case → endpoint → service → repository → entities → permission → transaction): do at the start of each slice
- [x] Dependencies approved: gin, pgx/v5, google/uuid, golang-jwt/jwt/v5, MicahParks/keyfunc/v3, go-playground/validator/v10. shopspring/decimal is added with the first money slice
- [x] Slice order: step by step, one prompt at a time (user)

### P3 — Foundation — Prompt 3
- [x] Config loading (internal/config, table tests)
- [~] DB pool (internal/db.NewPool): compiles; never run against a database
- [x] Transaction helper (+ tests: commit, rollback, commit error, begin error, panic)
- [x] Error model (+ tests) — internal/apperr + httpx.WriteError, internal detail never sent
- [x] Validation helper (httpx.BindJSON, json field names, tests)
- [~] Logging: slog JSON + request logger (route template only, no query/body); no automated test
- [x] Health endpoint (/health, 200 / 503, tests)
- [x] Auth + role middleware (+ tests). Web = Supabase JWT via JWKS, Worker = LINE ID token verify (D13 decided). Not tested against real Supabase/LINE
- [!] Area-permission middleware — D12 open; role checks only
- [x] Notifier interface + no-op (always reports `skipped`, never `sent`)
- [x] Storage interface + in-memory fake

### P4 — Contracts and survey — Prompt 4
- [~] 1S Create TOR + scope — code + unit/handler tests done 2026-10-05; **SQL unverified** (never run against PostgreSQL)
  - Endpoints: POST /api/web/contracts/info, /contracts/scope (check only), /contracts/confirm (201), GET /api/web/locations?q= (supervisor only)
  - Verified by tests (fake store): call order per Collaboration_1S; one tx; 23505 by constraint name → 409 (contract_no, location_name); FK 23503 → 400; requisition_no REQ-YYYYMMDD-NNN per Bangkok day with whole-tx retry (max 3) → 409 requisition_no_busy; equipment create-on-conflict-do-nothing then find; Q6 duplicate rules; field names prefixed contract./scope. on confirm; role access (401/403); notification after commit in background, own 10 s timeout, logs sent/skipped/failed, failure keeps the 201
  - Integration tests written (5), skipped without TEST_DATABASE_URL: unique 409 names, FK 400 name, equipment upsert, numbering + unique violation
  - Stubs: TODO(decision-file) contract file; TODO(decision-12) notification recipients = all active assistants
- [ ] 1A Site survey (to_buy_qty, validations, status)
- [ ] 9A Work continuation check

### P5 — Procurement and funding — Prompt 5
- [ ] 2S Transfer funds (duplicate transfer_ref_no, fund_transfer claim, status → pending_procurement)
- [ ] 2A Procurement (equipment_name via EQUIPMENT join, cumulative qty, actual_expense, receipt photo)
- [ ] 3S / 3.1S Review requisition — [!] approval-evidence part waits on D4

### P6 — Scheduling and leave — Prompt 6
- [ ] 3A Create schedule
- [ ] 1W View my schedule
- [ ] 2W Submit leave request
- [ ] 4A Leave review and replacement — D12 for area permission, D14 for leave key

### P7 — Attendance and evidence — Prompt 7
- [x] Decide D2 (evidence link) and D11 (attendance key / replacement) before coding: decided 2026-10-05 (docs/04 §2.1); the replacement population flow is still to be planned
- [ ] 3W Substitute check-in
- [ ] 4W Work report + check-out
- [ ] 4S Daily attendance job (re-runnable)

### P8 — Payroll and payslip — Prompt 8
- [ ] [!] D1, D3, D8, D9, D10 decided or explicitly stubbed (D1, D3 decided 2026-10-05; D8, D9, D10 open)
- [ ] 5S Deductions (attendance_id, tiers)
- [ ] 5S Payroll (base_wage, total_deduction, net_wage, is_paid=false), idempotent
- [ ] 7W Payslip (own data only)

### P9 — Extra equipment — Prompt 9
- [ ] [!] D2, D4, D5 decided or stubbed
- [ ] 5W Worker request (pending_survey)
- [ ] 6W Worker views result
- [ ] 5A / 6A Review, forward or reject
- [ ] 7A Additional purchase (locks, no duplicate submit)
- [ ] 8A Delivery evidence

### P10 — Financial report — Prompt 10
- [ ] [!] D6, D7 decided
- [ ] 6S Revenue, labor, material, net profit; PDF export
- [ ] Fixture with a fund_transfer and an actual_expense on one requisition (no double count)

### P11 — Notifications and jobs wiring — Prompt 11
- [ ] Real LINE notifier (current official docs; config from user)
- [ ] Scheduler for 4S and 5S
- [ ] Recipient checks per use case

### P12 — Review pass (read-only) — Prompt 12
- [ ] Review findings listed by severity
- [ ] Nothing committed confirmed

---

## Findings (append; newest last)
### P0 — 2026-10-05 (read-only verification)
| Item | Expected | Found | Match? |
|---|---|---|---|
| Branch / history | clean known state | `develop`, 3 commits (Initial commit → initial gin → "Added backend structure with no class modeling") | n/a |
| Working tree | — | M .gitignore (adds `CLAUDE.md` to ignore list); untracked README.md, cmd/, docs/, .DS_Store | n/a |
| Go / module | Go + Gin | go1.26.0; module `chrisfoong/chaum-work-management-backend`; gin v1.12.0 | yes |
| go.mod hygiene | direct deps listed | gin marked `// indirect` though imported directly (not tidied); mongo-driver/v2 present only transitively via gin/binding | minor |
| Entry point | one | two: `main.go` (package main → routes.SetupRouter) and `cmd/api/main.go` (package `chaumworkmanagementbackend`, not `main`; unreachable /ping; not gofmt'd) | no |
| Config / .env.example | config loading, placeholders | none; no `config/` (routes.go references `config/database.go`, missing) | no |
| DB / migrations / schema | PostgreSQL (Supabase), 16 tables | none; no .sql, no migrations, no DB driver | no (cannot verify items 4–7) |
| Database target in comments | PostgreSQL | handlers/routes comments say MongoDB collections + `bson.M` | **contradiction** |
| Identity | user_id = worker_id | only a comment in workers_handler.go ("shares its user_id as PK"); separate /users and /workers CRUD | consistent in intent; unverified physically |
| Canonical columns (docs/07 §A.6) | present in schema | no schema → unverifiable | unknown |
| Status value sets | in schema/code | none used in code | unknown |
| Auth / permissions | reuse existing | none (no middleware, no LINE verification) | absent (D13) |
| Routes | SD-driven workflows (docs/06) | generic CRUD per entity (15 groups + /health) + nested relationship GETs; all handlers return TODO JSON, no DB | shape differs |
| Relationship names in routes | docs/02 §2 | "CLEARS" for PAYROLL→DEDUCTION (ER: APPLIES R23); "HAS_EXPENSE" CONTRACT_TOR→EXPENSE_CLAIM (not in ER; claims link via requisition); "REQUEST" (ER: REQUESTS); no routes for PAYS, REVIEWS, REPLACES, requisition GENERATES expense claims | no |
| Tests | table-driven unit tests | none | n/a (baseline) |
| Baseline | build/vet/test | all exit 0; `go test`: no test files | yes |

### P1 migration review — 2026-10-05 (static, read-only; migration not executed)
Files: `migrations/0001_init.up.sql` (defines the schema), `migrations/0001_init.down.sql` (drops it in FK-safe order). No other schema files in the repo.
- Attributes: all 16 entities complete and spelled exactly. USER → table `users` (`user` is a reserved word). Extra columns: `equipment_requisition.reject_reason` (documented physical column) and `attendance.replacement_worker_id` (R18, D11). All other non-ER columns are FKs implementing relationships.
- Relationships: R01–R23 all implemented by FKs; R24 PAYS not implemented (D1); U1 = `worker.user_id` PK+FK. R16 is enforced only as 1:M (unique is (schedule_id, worker_id), not schedule_id) → see M2.
- Decided items: identity, D1, D2, D3, D11, canonical payroll/deduction columns, expense_claim.requisition_id with no requisition.expense_id: all reflected.
- Types: 11 money fields are all NUMERIC; 14 timestamps are all timestamptz; no float/real/money columns; `shift_start_time` is `time` (wall clock, Asia/Bangkok by convention).
- Enforced: unique contract_no, transfer_ref_no, (worker_id, work_date), location_name, equipment_name, all *_no columns. Not enforced: unique leave per worker per date (deliberate).
- ON DELETE: none declared → NO ACTION everywhere (deletes are blocked by references; no cascades).
- Indexes: only PK/UNIQUE. No secondary indexes on FKs used by SD queries.
- Status: CHECK on role, attendance.status, leave.status, requisition.status/type, expense_type. Free text (value set not specified): contract_tor.status, work_schedule.shift_status, company_invoice.status. `rejected` in requisition.status is implied by the reject flows, not observed.
- Supabase Auth: `users.user_id` = auth.users.id by convention only (no FK). `line_id`: text UNIQUE, nullable, on users (all roles).

Risks (smallest fix to consider; nothing applied):
| # | Sev | Risk | Smallest fix |
|---|---|---|---|
| H1 | High | RLS not enabled on any table. If Supabase's Data API exposes `public`, the anon/authenticated keys could read and write every table | `ENABLE ROW LEVEL SECURITY` on all 16 tables with no policies (backend owner role bypasses it), or don't expose `public` in the Data API. Verify against current Supabase docs |
| M1 | Med | No DB guard for 5S re-runs (no unique on payroll (worker_id, period_start, period_end) or on deduction attendance_id) | Decide D8, then add a unique key |
| M2 | Med | R16 1:1 is not enforced, and attendance.worker_id is not tied to work_schedule.worker_id | Decide whether replacement rows share the schedule; then UNIQUE(schedule_id) or a composite FK (schedule_id, worker_id) → work_schedule |
| M3 | Med | Duplicate leave for the same worker+date is not enforced (2W rule) | Partial unique index on (user_id, leave_date) WHERE status <> 'rejected', once re-submission is decided |
| M4 | Med | Nothing ensures a users row with role='worker' has a worker row, or that worker rows are role='worker' | Enforce in the service at user creation (no schema change), or a trigger |
| M5 | Med | line_id must hold the LINE userId (`sub` from token verify), not the public LINE ID; the column name doesn't say so | Comment/doc only |
| L1 | Low | No FK indexes (work_schedule.assignment_id, work_date; attendance.worker_id; leave(user_id, leave_date); deduction worker_id/payroll_id; payroll.worker_id; requisition assignment_id/status; requisition_item.requisition_id; expense_claim.requisition_id) | Add a few indexes once query shapes are known |
| L2 | Low | net_wage is not checked against the formula | CHECK (net_wage = base_wage - total_deduction) — formula is fixed; does not decide D10 |
| L3 | Low | fund_transfer claim may lack transfer_ref_no | CHECK (expense_type <> 'fund_transfer' OR transfer_ref_no IS NOT NULL) |
| L4 | Low | updated_at has no trigger | Set in the service or add a trigger |
| L5 | Low | billing_month day=1 CHECK, integer quantities, required_workers >= 0 are Claude's assumptions | Confirm |
| L6 | Low | attendance.work_date can differ from the schedule's work_date (docs/04 §3) | Service sets it from the schedule |
| L7 | Low | deduction.attendance_id NOT NULL makes TRIGGERS mandatory (ER gives max cardinality only) | Confirm every deduction comes from attendance |
| L8 | Low | payroll/schedule/leave FKs point to worker, not users, so only workers can have payroll | Matches identity decision; confirm |

### Repository layout notes (2026-10-05)
- `docs/references/` is git-ignored and exists only on the user's machine; it may be absent on a fresh clone. `docs/00`–`07` and this file are the tracked summary. If a reference file is missing: say so and ask; never guess.
- `CLAUDE.md` is now tracked: removed from `.gitignore` (user will `git add` it). Checked first: no secrets, connection strings or tokens (only rule text mentioning them). `.env` and `docs/references/` stay ignored.

### Flagged to the team (undefined screens/flows)
- Contract list screen ("รายการสัญญา", 1S step 2): no use case, no query, no class-diagram operation.
- Contract detail page (1S step 11): content undefined; for now the POST /contracts 201 response carries the data.
- Location search (1S step 6): GET /api/web/locations?q= added by the user's decision; not in the use case or class diagram — team to confirm.

### Class diagrams intake — 2026-10-05 (read-only)
`docs/references/Chaum-Diagrams-PlantUML.txt`: 4,496 lines, 55 diagrams, all read in full (Chaum_1–13, Chaum_A4_Summary, Data Flow, Component, 22 Collaboration_*_Normal_Case, Class_Interactions_1–15, State). No 3.1S, no 7S, no sequence-diagram code (images only, per the file header).
- Entities (Chaum_11–13) match datadict.txt column-for-column; they show types only (no NOT NULL/defaults). They differ from user decisions: role `assistant` (diagram `asst_supervisor`), requisition status set (diagram `pending_fund/pending_approval/cancelled`), `reject_reason` (absent in diagram).
- Operation names match docs/03 almost entirely. Differences: 3S notifies directly (no BackgroundJob, no 3.1S); 9A `checkAreaAccessPermission` moved to ContractRepository and one call `loadContractAndArea` does load + evaluation; new `WebUI.submitContractForm(..., contractFile)`; 5A `selectRequisition()`; 6A `openOriginalRequest()`; 6W `EquipmentResultController.loadEquipmentRequestResult/receiveEquipmentResult`; several params now named (status, message, nextStep).
- Still uses pre-correction names (conflict with docs/04 §1.3): `calculateNetPay`, `createDeduction(workerId, penaltyAmount, penaltyReason)`, `createPayroll(..., netPay)`, `showDigitalPayslipCard(..., total_wage, ...)`. Open — blocks P8, not P4.
- State diagram covers LEAVE_REQUEST only (pending → approved/rejected; replacement schedule in the same transaction). No state diagrams for requisition, contract, attendance, payroll.
- Proposed P4 plan changes (not applied): endpoints named after operations; confirmContract and processSiteSurvey make no extra re-check queries (rely on unique/FK constraints and a conditional UPDATE); 9A becomes one call; 1S notification after the response.

## Decisions log (user answers to D1–D14 and anything new)
| Date | Decision | Answer | Source |
|---|---|---|---|
| — | user_id = worker_id | Adopted: one identity | user |
| 2026-10-05 | Rebuild | Delete pre-existing code; rebuild step by step | user |
| 2026-10-05 | IDs | UUID, google/uuid | user |
| 2026-10-05 | D1 | No payer column | user |
| 2026-10-05 | D2 | work_evidence links by assignment_id only (data dictionary) | user |
| 2026-10-05 | D3 | deduction.payroll_id = FK to the payroll PK, nullable until applied | user |
| 2026-10-05 | D11 | Attendance key (schedule_id, worker_id); replacement stored only once accepted | user |
| 2026-10-05 | D13 | Supabase Auth JWT (JWKS); users.user_id = auth user id | Claude, at user's request |
| 2026-10-05 | D14 | leave_request.user_id for now — **TO DISCUSS: user prefers schedule_id** | user |
| 2026-10-05 | Value sets | role + attendance status lowercase, CHECK constraints | user |
| 2026-10-05 | D1 (revised) | payroll.managed_by_id FK users NOT NULL, set by 5S from the logged-in supervisor | user |
| 2026-10-05 | Role value | ~~'asst_supervisor'~~ → **keep 'assistant'** (user reversed; data dictionary says asst_supervisor — deliberate deviation) | user |
| 2026-10-05 | Requisition status | pending_survey, pending_procurement, pending_supervisor, approved, rejected, completed | user |
| 2026-10-05 | Names | requisition.requested_by; payroll.user_id; worker PK worker_id; leave_request.user_id → users | user, data dictionary |
| 2026-10-05 | expense_claim | transfer_ref_no required only for fund_transfer; expense_no always set; receipt_photo_url NOT NULL; total_amount > 0 | user |
| 2026-10-05 | line_id | LINE userId (sub), NOT NULL, UNIQUE | user |
| 2026-10-05 | D14 | keep user_id; UNIQUE (user_id, leave_date); no re-submission of leave | user, uc 2W Q2W.2 |
| 2026-10-05 | D11 (revised) | attendance UNIQUE (schedule_id) + UNIQUE (schedule_id, worker_id); substitute_worker_id nullable, unused; replacement = new work_schedule row | user, data dictionary, uc 4A Q4A.7 |
| 2026-10-05 | reject_reason | keep, nullable, until the team confirms | user |
| 2026-10-05 | Status sets | contract registered/active/complete/cancelled (default registered); invoice pending/paid (default pending); shift scheduled/completed/cancelled | user, data dictionary |
| 2026-10-05 | billing_month | varchar(7) 'YYYY-MM' with format CHECK | user, data dictionary |
| 2026-10-05 | Penalties | late <1h 300, 1–3h 400, >3h 1,500; absent 1,500 | user, uc 5S Q5S.2 |
| 2026-10-05 | requisition_no | `REQ-YYYYMMDD-NNN`, sequence per Bangkok day, generated in the transaction, retry on unique conflict | user (Claude's recommendation) |
| 2026-10-05 | 1S items | per area: one tor_base requisition + items per TOR_LOCATION_ASSIGNMENT (requisition has assignment_id FK) | user |
| 2026-10-05 | 1S notification | LINE to assistants (ผู้ดูแลงาน); which assistants → all active `assistant` users with line_id until D12 | user; recipient scope TODO(decision-12) |
| 2026-10-05 | 1A list (C15) | Q1A.0 adds `AND er.requisition_type = 'tor_base'` (5W additional requests go to 5A only) | user |
| 2026-10-05 | 1S duplicates (Q6) | reject the same location twice in one contract and the same equipment name twice in one area | user |
| 2026-10-05 | Q1–Q4, Q6 | accepted as recorded above; Q2 confirmed by Collaboration_1S (one initial requisition per area) | user |
| 2026-10-05 | Role (final) | keep 'assistant'; deliberate difference from datadict + class diagram ('asst_supervisor'); record in docs/04 + frontend mapping note | user |
| 2026-10-05 | Q5a | add GET /api/web/locations?q= (supervisor only, ≤20 matches); mark in docs/03 as an addition for team confirmation | user |
| 2026-10-05 | Q5b | POST /contracts 201 response carries the detail-page data; GET detail only if the team specifies the page | user |
| 2026-10-05 | Q5c | contract list screen: not built; flagged to the team as undefined (no use case, no query, no class-diagram operation) | user |
| 2026-10-05 | D12, D16, D17, D-file | stay open; stub with TODO markers | user |
| 2026-10-05 | Throwaway DB | not yet; no database use until the user supplies a local connection string | user |
| 2026-10-05 | C-P4-1..8 | all accepted with notes: 9A one GET (summary null until D17); checkAreaAccessPermission(assignmentId, userId) stub; no confirm re-check queries, 23505 by constraint name → 409, FK → 400, equipment insert-on-conflict-do-nothing then select, /info and /scope keep early checks and store nothing; 1A status update last, explicit value, conditional on pending_survey; every item exactly once + belongs to requisition (deviations from diagram); background notify with own context/timeout, no outbox; contracts saved without file until D-file; route-registration test for /requisitions/pending-survey vs /:id | user |
| 2026-10-05 | Reference intake | approved with corrections and applied to CLAUDE.md, docs/00, 02, 04, 05 | user |
| 2026-10-05 | D1 / D1b | D1: payer = the supervisor who runs 5S manually (uc 5S steps 1–2). D1b (open): automatic 5S start from 4S (uc 4S step 6) has no payer → not built until the team decides | user |
| 2026-10-05 | C5 | uc SQL uses w.worker_id and w.user_id (5S Q5S.1, 4A Q4A.1); datadict has only worker_id → read w.user_id as worker_id | user |
| 2026-10-05 | Round 2 | keep the migration as is; not applied to any database | user |
| 2026-10-05 | Docs D-1, D-2, D-3 | approved and applied (docs/04 role row + frontend notes; docs/03 class-diagram section + additions; docs/00 rows for PlantUML, datadict, uc, ER + precedence) | user |
(Add a row whenever the user settles a D#; also update docs/04 §2.)

## Files changed but not committed (append per task)
(Record paths so nothing is lost: since nothing is committed, this is the change list.)
- 2026-10-05 P0: docs/TASKS.md (bookkeeping only). No code changed.
- 2026-10-05 Rebuild: DELETED main.go, routes/routes.go, handlers/*.go (16; recoverable from cf538a6), cmd/api/main.go (never committed; gone), .DS_Store.
  ADDED cmd/server/main.go; internal/{config,db,apperr,httpx,health,auth,notify,storage,server}/*.go (+ _test.go); migrations/0001_init.up.sql, 0001_init.down.sql; .env.example.
  MODIFIED go.mod, go.sum (regenerated), .gitignore (+.DS_Store), docs/04_decisions_and_corrections.md (§2.1 decided), docs/TASKS.md.
- 2026-10-05 Handoff: docs/handoff/ created then DELETED at user request (Claude Desktop prompts are given in chat, not as files).
- 2026-10-05 P1 migration review: docs/TASKS.md only (read-only review; migration unchanged).
- 2026-10-05 Docs D-1..D-3: MODIFIED docs/00_INDEX.md, docs/03_workflows_operations.md, docs/04_decisions_and_corrections.md.
- 2026-10-05 Migration round 2: MODIFIED migrations/0001_init.up.sql (down unchanged).
- 2026-10-05 P4/1S: ADDED internal/contract/{model,validate,repository,service,handler}.go + validate_test.go, handler_test.go, repository_integration_test.go. MODIFIED internal/db/db.go (+DBTX, TxRunner, PoolTx, UniqueViolation, ForeignKeyViolation) + db_test.go; internal/apperr/apperr.go (+WithFields); internal/notify/notify.go (+Disabled, non-recording) + notify_test.go; cmd/server/main.go (wire 1S, notify.Disabled).
- 2026-10-05 Reference intake: MODIFIED CLAUDE.md, docs/00_INDEX.md, docs/02_domain_model.md, docs/04_decisions_and_corrections.md, docs/05_business_rules.md.
- 2026-10-05 Local-only references: MODIFIED CLAUDE.md, docs/00_INDEX.md (references local-only rule); .gitignore (CLAUDE.md line removed); docs/TASKS.md.
- 2026-10-05 P1 fixes round 1: MODIFIED migrations/0001_init.up.sql (H1, L1, L2, L3), migrations/0001_init.down.sql (comment only), docs/TASKS.md.

## Session log (append; keep each entry to 3–5 lines)
(date · task · what was verified · what is left)
- 2026-10-05 · P0 verify repo · git state, go toolchain, tree, routes/handlers, baseline build/vet/test (all pass, no tests) · no schema/migrations/config/auth exist; MongoDB comments contradict PostgreSQL; two entry points · left: user confirms restatement, supplies schema source → P1
- 2026-10-05 · Rebuild + P3 foundation + schema files · build/vet/test green (9 pkgs), gofmt clean · migration not applied or syntax-checked; D12 area check open; next P4
- 2026-10-05 · P1 migration review (static) · 16 entities, R01–R24/U1, decided items, types, constraints checked · 1 high (RLS), 5 medium, 8 low; nothing applied or executed · left: user picks fixes; optional syntax check on disposable DB
- 2026-10-05 · P1 fixes round 1 (H1, L1, L2, L3; notes M4, M5, L4, L6) · static re-check: 16/16 RLS, 12 indexes on existing columns, diff additive except one comment · migration still never executed; M1–M3, L5, L7, L8 wait on user
- 2026-10-05 · References intake (datadict.txt, uc.txt, er-latest.png read in full) + migration round 2 plan · plan only; decisions logged · waiting: approval of migration plan; docs 00/02/04/05 + CLAUDE.md update; Go role constant
- 2026-10-05 · Class diagrams intake (55 PlantUML diagrams read in full) · compared with docs/03, datadict and P4 plan · waiting: approvals (round 2, Q5, P4 changes, docs/00 + docs/03 diffs)
- 2026-10-05 · Decisions Q1–Q6, role, Q5a–c recorded · round 2 + P4 NOT started: conditional approval, class-diagram conflicts C-P4-1..8 shown to user · waiting for per-item answers and doc diffs approval
- 2026-10-05 · Docs D-1..D-3 applied; migration round 2 applied to file + static check (16/16 tables, datadict columns all present) · stopped at checkpoint before P4 · waiting: user go for P4; reference-intake doc diff approval
- 2026-10-05 · references marked local-only (CLAUDE.md, docs/00); CLAUDE.md un-ignored after secret scan · reference-intake doc diff prepared, not applied · waiting: user review of intake diff, then round-2 checkpoint (already shown), then P4
- 2026-10-05 · Reference intake applied (D1b open, C5 reworded) · round 2 kept · P4 scope listed, not started · waiting: "go P4"
- 2026-10-05 · P4 step 1 (1S) built · go build/vet ok, gofmt clean, go test ./... 10 packages ok (56 contract test cases incl. subtests pass, 5 integration skipped) · SQL unverified (no DB) · waiting: user review before 1A
