# 04 — Decisions, corrections and open items

## 1. Adopted decisions (do not re-ask)

### 1.1 Identity: `user_id = worker_id`
A worker is a user. `worker_id` and `user_id` are the same value everywhere (schedule, attendance, leave, evidence, deduction, payroll). Do not create a second ID space. If the physical schema has both columns, report it and keep them equal; do not migrate without approval. WORKER holds only worker-specific data (`is_available`); role comes from `USER.role`.
Leave requests keyed by `user_id` are valid. Payroll `worker_id` is the same user.

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
| 5S payroll | Inserts `worker_id, period_start, period_end, base_wage, total_deduction, net_wage, is_paid=false` |
| 6S labor | Sums `net_wage` (was net_pay) |

## 2. Open decisions — do NOT guess; stub and report
Mark code `TODO(decision-N)` and say which workflow is blocked. Unrelated workflows continue.

| # | Decision | Blocks |
|---|---|---|
| D1 | USER PAYS PAYROLL: who is the payer, which identifier, authorization. Corrected SQL keeps the existing worker_id recipient mapping | 5S, 7W, payroll auth |
| D2 | Evidence link: store `schedule_id`, `assignment_id`, or infer assignment through the schedule (4W, 8A insert schedule_id + worker_id) | 4W, 8A |
| D3 | How/when a deduction's `payroll_id` is populated (APPLIES). 5S now writes attendance_id, but the payroll link is not shown | 5S, 7W |
| D4 | Approval evidence: 7A relies on reviewed_by/reviewed_at/prior approval; shown 3S approve only updates status; 6A writes review fields on reject but not on escalation to `pending_supervisor` | 3S, 6A, 7A |
| D5 | Funded additional requests: 2S sets `pending_procurement`; 7A lists `approved` or a specific `pending_supervisor` funded path. Resolve without inventing a new approval stage | 2S, 7A |
| D6 | 6S labor allocation per project: the labor query sums all paid payroll in the interval with no project condition | 6S |
| D7 | 6S reporting scope: revenue uses `billing_month`, other totals use date ranges | 6S |
| D8 | 5S deduction period: sums by `created_at`, not attendance work date; and a duplicate guard on re-run (5S must be idempotent per the team's requirement, but the shown insert has no guard) | 5S |
| D9 | Where payroll becomes `is_paid = true`. 5S sets false; 6S selects paid; 7W assumes paid. No SD shows the payment step | 6S, 7W |
| D10 | Negative net wage, refunds, cancellation/reversal behavior | payroll, 6S |
| D11 | Attendance UPSERT key `(schedule_id, worker_id)` vs real uniqueness; replacement attendance flow (WORKER REPLACES) | 3W, 4S |
| D12 | Authorization: how MANAGES (USER–CONTRACT_TOR) becomes per-assignment/area permission. A role value alone is not enough | all Assistant flows |
| D13 | How Supervisor/Assistant authenticate on the web (reuse whatever the repo already has; do not add a new method) | auth middleware |
| D14 | Physical key for WORKER SUBMITS LEAVE_REQUEST (2W/4S SQL still uses user_id; valid under 1.1) | 2W, 4S, 4A |

### 2.1 Decided (user, 2026-10-05)
| # | Answer | Effect |
|---|---|---|
| D1 | No payer column; USER PAYS PAYROLL is not stored | `payroll` has only `worker_id` (recipient). Payroll authorization is still open |
| D2 | Evidence links by `assignment_id` only, per the team's data dictionary | `work_evidence(assignment_id, worker_id)`, no `schedule_id`. 4W/8A must resolve the assignment from the schedule |
| D3 | `deduction_transaction.payroll_id` references the payroll primary key | Nullable FK, NULL until the deduction is applied to a payroll. *When* 5S sets it is still to be planned in P8 |
| D11 | Attendance key is `(schedule_id, worker_id)`; a replacement is stored only once the replacement worker accepts | `UNIQUE (schedule_id, worker_id)`; nullable `attendance.replacement_worker_id`. The population flow is planned in P7 |
| D13 | Web login = Supabase Auth (chosen by Claude at the user's request) | Bearer access token verified via the project JWKS (ES256/RS256, issuer `<SUPABASE_URL>/auth/v1`, expiry required). `users.user_id` = Supabase Auth user id; no new column |
| D14 | Use `user_id` for now. **To discuss:** the user thinks it should be `schedule_id` | `leave_request.user_id`. Revisit before P6 |

Also decided: UUID keys (`github.com/google/uuid`); stored roles `supervisor`/`assistant`/`worker` and attendance statuses `on_time`/`late`/`leave`/`absent`, both enforced with CHECK constraints. Worker login: LIFF ID token verified by POST to `https://api.line.me/oauth2/v2.1/verify`, `sub` matched to `users.line_id`.

Still open: D4, D5, D6, D7, D8, D9, D10, D12.

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
