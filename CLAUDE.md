# CLAUDE.md — Backend (Go / Gin / PostgreSQL on Supabase)

Backend for **Chaum Resource Management**: TOR contracts, work locations, equipment procurement, staffing, leave, attendance, payroll and financial reporting. This repo is **backend only**; the Next.js frontend is a separate repo.

**All references live in this file and in `docs/`. There is nothing else to look up.** Start with `docs/00_INDEX.md`. `docs/00`–`07` and `docs/TASKS.md` are the tracked summary. **`docs/references/` is git-ignored and local-only: it may be absent on a fresh clone.** If a reference file (or anything else) you need is missing, say so and ask the user — never guess its content, never invent it, and never claim to have read a file that is not in the repo.

**State lives in `docs/TASKS.md`. Read it at the start of every session and update it at the end of every task** (nothing is committed, so git does not show progress).

Core docs (loaded with this file):
@docs/TASKS.md
@docs/00_INDEX.md
@docs/02_domain_model.md
@docs/04_decisions_and_corrections.md

Read on demand: `docs/01_project_overview.md`, `docs/03_workflows_operations.md`, `docs/05_business_rules.md`, `docs/06_architecture_and_conventions.md`, `docs/07_verification_checklist.md`, `docs/prompts/CLAUDE_Prompts_Backend.md`, `docs/references/` (local-only; may be absent — ask, never guess).

---

## 0. Hard rules (always apply)

1. **Plan first.** For any non-trivial task, write a plan and wait for approval before editing files.
2. **Verify before proceeding.** The repo was initialised by a human before you arrived. Never assume what exists. Inspect it (tree, `go.mod`, code, migrations, env example, tests, `git status`) and report what you found versus these docs before adding anything. Do not overwrite or re-initialise existing work.
3. **Do not commit, push, add, reset, stash, branch or tag.** Leave all changes uncommitted. Read-only git (`status`, `diff`, `log`) is fine.
4. **Never apply migrations or write to a real database** (local, staging, Supabase) without explicit approval in the current task. Writing migration *files* only after an approved plan.
5. **No new tables, no database over-engineering.** If a workflow seems to need a new table or column, stop and ask.
6. **Do not guess business rules.** Open decisions are listed in `docs/04` §2 (D1–D18) and §2.2 (conflicts C#). Stub with `TODO(decision-N)` and report.
7. **Never print, log or commit secrets** (`.env`, Supabase keys, LINE secrets, tokens). Only edit `.env.example` with placeholders.
8. **Stop and ask** when the repo, schema and docs disagree. Report the difference first; do not silently choose the easier source.
9. **Keep `docs/TASKS.md` current:** tick a box only for work you verified, record files changed (uncommitted), decisions made and what is blocked. You may edit that file without a separate plan; show the diff.
10. **Do not add** offline sync, extra approval stages, extra reports, extra queries/nested workflow steps, or alternative sign-in methods that no requirement asks for.

## 1. Stack and platforms
Backend Go + Gin. Database PostgreSQL (Supabase). Storage Supabase Storage. Supervisor and Assistant use **desktop web**; Worker uses **only the LINE Mini App (LIFF)** on mobile — never move Worker flows to desktop screens. Roles are values of `USER.role` (Supervisor, Assistant, Worker); they are not subclasses.

## 2. Adopted decisions (do not re-ask)
- **`user_id = worker_id`.** A worker is a user; one ID space; the WORKER table's key column is `worker_id` (= `users.user_id`). If the schema has both columns, report it and keep them equal; do not migrate without approval.
- **Physical names/types** follow `docs/references/datadict.txt` (local-only) with the deviations recorded in `docs/04`; role values are `supervisor`, `assistant`, `worker` (deliberate difference from the data dictionary's `asst_supervisor`).
- **Canonical names (production ER):** PAYROLL `base_wage`, `total_deduction`, `net_wage` (derived, may be stored), `is_paid`, recipient `user_id`, payer `managed_by_id`; DEDUCTION_TRANSACTION `penalty_amount` (+ `attendance_id` link); COMPANY_INVOICE `exat_deduction_amount` (keep the ER spelling); `equipment_name` lives on EQUIPMENT; EXPENSE_CLAIM links by `requisition_id` (the requisition has no `expense_id`). **Never introduce** `total_wage`, `net_pay`, `deduction`, a payroll `status` column, or `dt.amount`.
- **Corrections adopted:** 7W reads base_wage/total_deduction/net_wage/is_paid and penalty_amount; 2A gets equipment_name via JOIN on equipment_id; 6S joins `ec.requisition_id = er.requisition_id` and sums `net_wage`; 5W creates `pending_survey`; 6S material cost = `actual_expense` claims on requisitions in approved/pending_supervisor/pending_procurement/completed (fund transfers excluded); 5S inserts deductions with `worker_id, attendance_id, penalty_amount, reason` and payroll with `user_id, managed_by_id, payroll_slip_no, period_start, period_end, base_wage, total_deduction, net_wage, is_paid=false` (column names per the data dictionary; payer = the supervisor who runs 5S from the payroll menu, D1; an automatic 5S start is open D1b). These are diagram corrections — verify the physical columns.
- Full list and the open decisions: `docs/04_decisions_and_corrections.md`.

## 3. Business rules you must not alter (details in `docs/05`)
`to_buy_qty = GREATEST(required_qty − existing_qty, 0)` · procurement is cumulative (`actual_qty = actual_qty + new_qty`, no receipt table) · `fund_transfer` vs `actual_expense`, only `actual_expense` counts as material cost · late deduction tiers: <1h = 300, 1–3h = 400, >3h = 1,500, absent = 1,500 (uc 5S), advance-notice leave = no deduction · `net_wage = base_wage − total_deduction` · 5S idempotent · `net_profit = revenue − (labor + material)` · worker is linked to a location only through WORK_SCHEDULE · do not change thresholds, formulas, paid-status meaning or negative-pay treatment.

## 4. Open decisions — never guess (blocks listed in `docs/04` §2; decided ones in §2.1)
Open: D1b automatic 5S start from 4S has no payer — not built · D4 approval evidence (3S/6A/7A) · D5 funded extra-request path (2S→7A) · D6 6S project labor allocation · D7 6S date range vs billing_month · D8 5S deduction period + re-run guard · D9 where `is_paid` becomes true · D10 negative pay/refunds/reversals · D12 area/project authorization from MANAGES · D15 attendance row creation before 4S · D16 who sets contract 'active' · D17 9A summary storage · D18 shared `pending_supervisor` (2S/3S) and unreachable 7A funded branch · D-file contract file type/upload.
Decided (do not re-ask): D1 payer = `managed_by_id` = the supervisor who runs 5S manually (uc 5S steps 1–2) · D2 evidence by `assignment_id` · D3 `payroll_id` FK, NULL until applied · D11 attendance UNIQUE (schedule_id), replacement = new schedule row · D13 Supabase Auth · D14 `user_id` + UNIQUE (user_id, leave_date).

## 5. Architecture (proposal — adapt to the existing repo; details in `docs/06`)
`handler (Gin) → service (rules + transactions) → repository (parameterised SQL) → PostgreSQL`, with a notifier interface (real LINE + clearly-marked no-op that is never reported as "delivered"), a storage interface, and jobs (4S attendance, 5S payroll) that call services and are safe to re-run. SD controllers/repositories are design responsibilities, not one-file-per-class; keep similarly named participants separate unless mapped.

## 6. Conventions (details in `docs/06`)
Money = `NUMERIC` / decimal, never float · UTC `timestamptz`, Asia/Bangkok for day logic · follow the schema's ID type · transactions for multi-write flows (say whether source-stated or technical) · re-read inside the transaction for "recheck latest / prevent duplicate" · typed errors → consistent JSON, name failing fields · no secrets/PII in logs · reuse existing auth; Worker endpoints verify LINE identity server-side (check **current official LINE/LIFF docs**, not memory) · role **and** area checks where required · uploads via Supabase Storage with server-side validation · table-driven unit tests; integration tests only on a disposable DB · small functions, no speculative abstractions, no unrelated refactors.

## 7. Before any new work — verify (full checklist: `docs/07` §A)
`git status/log`; Go version/`go.mod`/Gin; tree and entry point; migrations vs `docs/02` for all 16 entities (types, PK/FK/unique); identity (one key or two); presence of `payroll.total_deduction/net_wage/is_paid`, `deduction_transaction.attendance_id/penalty_amount`, `expense_claim.requisition_id/expense_type`, requisition review columns; status values in use; existing auth; baseline `go build ./...`, `go vet ./...`, `go test ./...`. Report a table: expected / found / match?

## 8. Definition of done and reporting
Plan approved; matches the docs (differences reported); validations/errors match the use case; build/vet/tests pass (state exactly what ran); no secrets, no unrelated changes, **nothing committed**. Update `docs/TASKS.md` first. End every task with: **Inspected · Changed · Ran · Decisions touched (D#) · Blocked/open · Uncommitted changes present (yes)**. Code files named after diagram classes are not completion.

## 9. Suggested safety net (check current Claude Code docs for exact syntax)
`.claude/settings.json`: `{"permissions":{"deny":["Bash(git commit:*)","Bash(git push:*)","Bash(git add:*)","Bash(git reset:*)"]}}`
