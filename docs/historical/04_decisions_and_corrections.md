# 04 — Decisions, corrections and open items

## 1. Adopted decisions (do not re-ask)

### 1.1 Identity: `user_id = worker_id`
A worker is a user. `worker_id` and `user_id` are the same value everywhere (schedule, attendance, leave, evidence, deduction, payroll). Do not create a second ID space. If the physical schema has both columns, report it and keep them equal; do not migrate without approval. WORKER holds only worker-specific data (`is_available`); role comes from `USER.role`.
Leave requests keyed by `user_id` are valid. The WORKER table's key column is `worker_id` (= `users.user_id`); payroll's recipient column is `user_id` (data dictionary).

### 1.2 Platforms and repos
Backend repo first (Gin + PostgreSQL/Supabase). Frontend separate (Next.js). Supervisor/Assistant on desktop web; Worker only on the LINE Mini App.

### 1.3 Diagram corrections already adopted
These are corrections to the *diagrams*. They say what the code must do; they do not prove the database already has the columns.

| Case | Correction |
|---|---|
| 7W payslip | Read `base_wage`, `total_deduction`, `net_wage`, `is_paid`. Payslip shows base_wage |
| 7W deductions | Read `DEDUCTION_TRANSACTION.penalty_amount` |
| 2A equipment name | JOIN REQUISITION_ITEM to EQUIPMENT on `equipment_id` to read `equipment_name` |
| 6S expense join | `ec.requisition_id = er.requisition_id`; many claims per requisition |
| 5W initial status | New additional requests start as `pending_survey` (so they appear in the 5A/6A queue). No transition step added |
| 6S material cost | Requisitions with status in (`approved`, `pending_supervisor`, `pending_procurement`, `completed`); only `expense_type = 'actual_expense'`; fund transfers excluded |
| 5S deduction | `createDeduction(workerId, attendanceId, penaltyAmount, penaltyReason)` inserts `worker_id, attendance_id, penalty_amount, reason` |
| 5S payroll | Inserts `user_id, managed_by_id, payroll_slip_no, period_start, period_end, base_wage, total_deduction, net_wage, is_paid=false` (column names per the data dictionary; payer per D1) |
| 6S labor | Sums `net_wage` (was net_pay) |

Note: `docs/references/uc.txt` (5S, 6S, 7W) still contains the pre-correction names (`deduction`, `net_pay`, `p.status`), and its 6S material query has no requisition-status filter; the class diagrams still use `calculateNetPay`, `createPayroll(…, netPay)`, `createDeduction` without attendanceId and `total_wage`. The corrections above stand; the class-diagram naming conflict is C16 (§2.2).

## 2. Open decisions — do NOT guess; stub and report
Mark code `TODO(decision-N)` and say which workflow is blocked. Unrelated workflows continue.

| # | Decision | Blocks |
|---|---|---|
| D1 | **Decided — see §2.1.** USER PAYS PAYROLL: who is the payer, which identifier, authorization. Corrected SQL keeps the existing worker_id recipient mapping | 5S, 7W, payroll auth |
| D2 | **Decided — see §2.1.** Evidence link: store `schedule_id`, `assignment_id`, or infer assignment through the schedule (4W, 8A insert schedule_id + worker_id) | 4W, 8A |
| D3 | **Decided — see §2.1.** How/when a deduction's `payroll_id` is populated (APPLIES). 5S now writes attendance_id, but the payroll link is not shown | 5S, 7W |
| D4 | Approval evidence: 7A relies on reviewed_by/reviewed_at/prior approval; shown 3S approve only updates status; 6A writes review fields on reject but not on escalation to `pending_supervisor` | 3S, 6A, 7A |
| D5 | Funded additional requests: 2S sets `pending_procurement`; 7A lists `approved` or a specific `pending_supervisor` funded path. Resolve without inventing a new approval stage | 2S, 7A |
| D6 | 6S labor allocation per project: the labor query sums all paid payroll in the interval with no project condition | 6S |
| D7 | 6S reporting scope: revenue uses `billing_month`, other totals use date ranges | 6S |
| D8 | 5S deduction period: sums by `created_at`, not attendance work date; and a duplicate guard on re-run (5S must be idempotent per the team's requirement, but the shown insert has no guard) | 5S |
| D9 | Where payroll becomes `is_paid = true`. 5S sets false; 6S selects paid; 7W assumes paid. No SD shows the payment step | 6S, 7W |
| D10 | Negative net wage, refunds, cancellation/reversal behavior | payroll, 6S |
| D11 | **Decided — see §2.1.** Attendance UPSERT key `(schedule_id, worker_id)` vs real uniqueness; replacement attendance flow (WORKER REPLACES) | 3W, 4S |
| D12 | Authorization: how MANAGES (USER–CONTRACT_TOR) becomes per-assignment/area permission. A role value alone is not enough | all Assistant flows |
| D13 | **Decided — see §2.1.** How Supervisor/Assistant authenticate on the web (reuse whatever the repo already has; do not add a new method) | auth middleware |
| D14 | **Decided — see §2.1.** Physical key for WORKER SUBMITS LEAVE_REQUEST (2W/4S SQL still uses user_id; valid under 1.1) | 2W, 4S, 4A |
| D1b | Who is the payer when 5S starts automatically from 4S (uc 4S step 6)? There is no logged-in supervisor; the automatic start is NOT built until the team decides | 4S → 5S link, P11 scheduler |
| D15 | When is the ATTENDANCE row created? 4S updates by attendance_id, but only 3W (substitute check-in) inserts rows; there is no normal check-in use case and no row for absent workers | 4S, 4W, 5S |
| D16 | Who changes CONTRACT_TOR.status from `registered` to `active`? No use case does; 9A requires `active` (until decided, 9A always answers "cannot continue") | 9A, 3A |
| D17 | Where is the 9A "summary from the supervisor" stored? No table or use case | 9A step 1 |
| D18 | `pending_supervisor` is used by 2S (from 2A shortfalls) and 3S (from 6A purchases); the 7A funded branch (pending_supervisor + fund_transfer + reviewed_by) cannot be reached because 2S sets `pending_procurement` and 3S never sets reviewed_by. Overlaps D4/D5 | 2S, 3S, 7A |
| D-file | Contract file: allowed type(s) (uc 1S says `.png` only; data-dictionary samples are `.pdf`) and upload timing (class diagram: with step 1) | 1S |

### 2.1 Decided (user, 2026-10-05)
| # | Answer | Effect |
|---|---|---|
| D1 | **Revised (manual 5S only):** the payer is the supervisor who runs 5S — uc 5S steps 1–2: opens the payroll menu, enters the period, clicks process; `payroll.managed_by_id` (FK users, NOT NULL) = that user | `payroll.user_id` = recipient, `managed_by_id` = payer. **Not fully settled:** the automatic start from 4S is open decision D1b and is not built |
| D2 | Evidence links by `assignment_id` only, per the team's data dictionary | `work_evidence(assignment_id, worker_id)`, no `schedule_id`. 4W/8A must resolve the assignment from the schedule |
| D3 | `deduction_transaction.payroll_id` references the payroll primary key | Nullable FK, NULL until the deduction is applied to a payroll. *When* 5S sets it is still to be planned in P8 |
| D11 | **Revised:** one attendance per schedule; a replacement gets a new WORK_SCHEDULE row (uc 4A Q4A.7) | `UNIQUE (schedule_id)` + `UNIQUE (schedule_id, worker_id)` (ON CONFLICT target for 3W); `substitute_worker_id` nullable, unused |
| D13 | Web login = Supabase Auth (chosen by Claude at the user's request) | Bearer access token verified via the project JWKS (ES256/RS256, issuer `<SUPABASE_URL>/auth/v1`, expiry required). `users.user_id` = Supabase Auth user id; no new column |
| D14 | Keep `user_id` (FK users); no re-submission of leave for the same date, whatever its status (uc 2W Q2W.2) | `UNIQUE (user_id, leave_date)` |
| Penalties | Late < 1h = 300; 1–3h inclusive = 400; > 3h = 1,500; absent = 1,500 | uc 5S Q5S.2 |
| Names | `equipment_requisition.requested_by`; `payroll.user_id`; WORKER key `worker_id`; `leave_request.user_id` → users | data dictionary |
| Requisition status | `pending_survey`, `pending_procurement`, `pending_supervisor`, `approved`, `rejected`, `completed` | user (data dictionary differs) |
| Status sets | contract `registered/active/complete/cancelled` (default registered); invoice `pending/paid`; shift `scheduled/completed/cancelled` | data dictionary |
| billing_month | `varchar(7)` 'YYYY-MM' with a format CHECK | data dictionary |
| expense_claim | `expense_no` always set; `transfer_ref_no` required only for fund_transfer; `receipt_photo_url` NOT NULL; `total_amount > 0` | user |
| line_id | LINE userId (`sub` of ID-token verification), NOT NULL, UNIQUE | user, uc 3W regex |
| reject_reason | kept (nullable) until the team confirms | user |
| 9A rule | continue only if contract status = `active` AND `end_date >=` today (Asia/Bangkok date) | uc 9A step 7 |
| actual_price | unit price | data dictionary |
| 1A result | all `to_buy_qty` = 0 → `completed`, else `pending_procurement`; list filters `requisition_type = 'tor_base'` | uc 1A step 11, user |
| 1S | one tor_base requisition + items per area; `requisition_no` = `REQ-YYYYMMDD-NNN` (per Bangkok day, retry on conflict); reject the same location twice or the same equipment name twice in one area; LINE notice to active assistants with a line_id until D12 | user |
| Role | Stored `users.role` values are `supervisor`, `assistant`, `worker`. Deliberate difference from datadict.txt and the class diagram (`asst_supervisor`), decided by the user 2026-10-05 | migrations CHECK; internal/auth RoleAssistant = "assistant" |

Also decided: UUID keys (`github.com/google/uuid`); stored roles `supervisor`/`assistant`/`worker` and attendance statuses `on_time`/`late`/`leave`/`absent`, both enforced with CHECK constraints. Worker login: LIFF ID token verified by POST to `https://api.line.me/oauth2/v2.1/verify`, `sub` matched to `users.line_id`.

Still open: D1b, D4, D5, D6, D7, D8, D9, D10, D12, D15, D16, D17, D18, D-file.

### 2.2 Source conflicts (data dictionary vs use cases vs class diagrams)
| # | Conflict | Status |
|---|---|---|
| C1 | role `asst_supervisor` (datadict, class diagram) vs `assistant` | resolved: `assistant` (deliberate difference) |
| C2 | requisition status sets differ | resolved: user set (§2.1) |
| C3 | `requested_by` vs uc `user_id` | resolved: `requested_by` |
| C4 | payroll payer `managed_by_id` vs old D1 | resolved: D1 revised |
| C5 | uc SQL uses both `w.worker_id` and `w.user_id` on WORKER (e.g. 5S Q5S.1, 4A Q4A.1); the data dictionary has only `worker_id` (= users.user_id) | resolved: read `w.user_id` as `w.worker_id` when adapting uc SQL; no `worker.user_id` column |
| C6 | `expense_no`/`transfer_ref_no` NOT NULL vs inserts that omit them | resolved: §2.1 expense_claim |
| C7 | `substitute_worker_id` purpose; 3W ON CONFLICT key | resolved: D11 revised |
| C8 | `line_id` sample values look like public LINE IDs vs 3W userId regex | resolved: LINE userId |
| C9 | `to_buy_qty` "computed" vs written by 1A/5W | resolved: plain column (deviation) |
| C10 | 2A overwrites `actual_price`; 7A keeps a weighted average | **open** |
| C11 | `reject_reason` absent from datadict | resolved: kept until the team confirms |
| C12 | contract file `.png` only (uc) vs `.pdf` samples | **open** (D-file) |
| C13 | IDs validated as integers in uc 3A/4A/5A–9A vs UUID | handled as UUID (datadict types); team to confirm |
| C14 | 1W shift labels ("pending confirm / confirmed / leave") vs `shift_status` values | **open** |
| C15 | 1A list includes 5W `additional` requests | resolved: filter `tor_base` |
| C16 | class diagrams use `calculateNetPay`, `createPayroll(…, netPay)`, `createDeduction` without attendanceId, `total_wage` vs docs/04 §1.3 | **open** (blocks P8) |

## 3. Other known inconsistencies in the sources (informational)
- 4S reads check_in by schedule_id and worker_id but its SELECT does not return attendance_id, while the update is by attendance_id. 4S also uses a worker parameter in a `user_id` predicate (valid under 1.1).
- `work_date` exists on both WORK_SCHEDULE and ATTENDANCE; consistency between them is not specified.
- Unique transfer reference, one schedule per worker per day and the attendance UPSERT key are use-case rules; whether the database enforces them is unknown.
- `exat_deduction_amount` looks like a typo but its intended replacement is unknown. Keep it; any rename needs an old/new mapping and a migration plan.
- Four ER drawing defects (CONTRACT_TOR.status connector, TRIGGERS endpoint, REQUESTS endpoint) are drawing issues, not missing business relationships. HOSTS was fixed.
- SDs that contain `...` in SQL are abbreviated documentation, not executable SQL. Omitted columns/joins cannot be verified.
- The full audit that produced these lists is in `docs/references/Production_ER_Conflict_Review.md`.

## 4. Rule for conflicts
Physical schema vs ER: report first, never silently choose the easier source. Never use an old field name just because it appears in the historical class diagram.

## 5. Notes for the frontend repo (Next.js)
- Role values from the API are `supervisor`, `assistant`, `worker`. Where team documents say `asst_supervisor`, it means `assistant`.
- Web endpoints are under `/api/web` (Supabase Auth bearer token); Worker endpoints under `/api/liff` (LIFF ID token).
- Contract list and contract detail screens are undefined by the team (see docs/TASKS.md). After 1S, render the detail page from the POST /api/web/contracts/confirm 201 response.
- Location search: GET /api/web/locations?q= (supervisor; up to 20 matches) — an addition pending team confirmation.
