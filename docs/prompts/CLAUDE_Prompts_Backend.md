# Prompts for Claude in VS Code — Backend repo

How to use: put `CLAUDE.md` in the backend repo root and the whole `docs/` folder next to it, then paste these **in order**. Every reference lives in `CLAUDE.md` and `docs/` — there is nothing else to look up.
Each prompt is plan-first. Do not move to the next until you have read the reply and
approved it. Prompts 0–2 are read-only.

Optional safety net: block commits at the tool level in `.claude/settings.json`
(check the current Claude Code docs for exact syntax):

```json
{
  "permissions": {
    "deny": ["Bash(git commit:*)", "Bash(git push:*)", "Bash(git add:*)", "Bash(git reset:*)"]
  }
}
```

---

## Shared preamble (paste at the top of every prompt, or keep it in CLAUDE.md)

```
Follow CLAUDE.md strictly. Read docs/TASKS.md first to learn the current state, and update it
at the end (only tick what you verified). All references are in CLAUDE.md and docs/ (start at docs/00_INDEX.md);
if something you need is not there, say so and ask — never invent it or claim to have read
a file that is not in the repo. Plan first and wait for my approval before editing.
Verify what actually exists in the repo before proceeding; do not assume.
Do NOT commit, push, add, reset or branch. Do NOT apply migrations or write to any
real database. Do not guess open decisions (D1–D14); stub and report them.
End with the reporting format in docs/07 §E.
```

---

## Prompt 0 — Verify the existing repo (READ-ONLY)

```
I initialised something in this repo before you started. Do not change any file.

1. Read CLAUDE.md and docs/00_INDEX.md through docs/07_verification_checklist.md (skim docs/references/ only when a doc points to it). Restate, in under 15 lines, the hard rules and the adopted
   decisions (especially user_id = worker_id) so I can confirm you understood.
2. Run the verification checklist in docs/07_verification_checklist.md §A using read-only commands only
   (git status/log, ls/tree, go env, go build ./..., go vet ./..., go test ./...).
   If a command could write or connect to a real database, do not run it; tell me instead.
3. Report a table: item / expected / found / match? Include:
   - what I already initialised (structure, packages, config, routes, migrations, tests)
   - anything that contradicts CLAUDE.md
   - anything you could not verify and why
4. List which docs/ files you read in full, which you skimmed, and which expected
   material is absent (see 'Material that is NOT available' in docs/00_INDEX.md, plus any
   schema in the repo). Do not claim to have read a file you did not open.
5. Do not propose code yet. End by asking only the questions that block Prompt 1.
```

---

## Prompt 1 — Schema reconciliation report (READ-ONLY)

```
Compare the physical schema in this repo (migration files / SQL schema; do not connect
to a live database unless I say so) against docs/02_domain_model.md and docs/04_decisions_and_corrections.md.

Produce docs/decisions/schema-mapping.md (create the file only after showing me the
outline and getting approval) with:
1. Entity map: each of the 16 ER entities → table, PK, FKs, unique constraints, types.
2. Canonical-name check for every field in docs/02 §1 and the corrections in docs/04 §1.3. Mark each as
   present / missing / named differently.
3. Identity check for user_id = worker_id: how WORKER relates to USER in the schema.
   Report whether both columns exist. Do not propose a migration yet.
4. Status value sets actually used per entity.
5. Constraint gaps that matter for the use cases: one schedule per worker per day,
   attendance uniqueness / UPSERT key, unique transfer_ref_no, unique contract_no,
   idempotency for payroll and deductions.
6. For each open decision D1–D14, what the schema evidence says (if anything).
7. A list of differences between schema and ER. For each: impact, and the smallest
   possible fix to *consider*. No migrations are written in this step.

Stop after the report and wait.
```

---

## Prompt 2 — Architecture and skeleton plan (PLAN ONLY)

```
Using the Prompt 0 and Prompt 1 findings, propose how to structure the backend.
Reuse the existing layout wherever it already exists; do not restructure working code.

Deliver a plan (no code) containing:
1. Package/directory layout (handler → service → repository, notifier, storage, jobs).
2. A workflow matrix: use case → endpoint(s) → service → repository functions →
   entities touched → role/area permission → transaction boundary.
   Cover 1S–6S, 1A–9A, 1W–7W from docs/06_architecture_and_conventions.md (workflow map) and docs/03_workflows_operations.md.
3. Cross-cutting design: config, DB pool, transaction helper, error model, request
   validation, auth middleware (reuse existing), role + area permission, LINE identity
   resolution for Worker endpoints, Supabase Storage interface, notifier interface
   with a clearly marked no-op, job runner for 4S and 5S.
4. Which slices can start now vs which are blocked by which D#.
5. A proposed slice order with acceptance checks drawn from the use cases.
6. Risks and anything where the repo, schema and CLAUDE.md disagree.

Mark every new dependency or technology choice as a PROPOSAL with the reason.
Wait for my approval.
```

---

## Prompt 3 — Foundation slice (after approval)

```
Implement only the foundation from the approved plan:
config loading, DB pool, transaction helper, error model, validation helper,
logging, health endpoint, auth/role/area-permission middleware (reusing existing auth),
notifier and storage interfaces with test doubles.

Rules:
- Show the file list you will touch first; wait for my OK if it differs from the plan.
- No business workflows yet. No migrations unless the approved plan says so.
- Money type and time handling per docs/06 conventions.
- Add unit tests for the transaction helper, error mapping and permission checks.
- Run go build ./..., go vet ./..., go test ./... and report exact results.
- Leave everything uncommitted.
```

---

## Prompt 4 — Slice 1: contracts and requisition basics (1S, 1A)

```
Plan first, then implement after I approve: 1S Create TOR + scope and 1A Site survey.

Plan must cite the docs/ sections and SD operations you are following (the original use-case tables are not in the repo; ask if a rule is missing) and list every validation.
Behaviour to preserve:
- 1S: duplicate contract_no check, location validation/creation, TOR_LOCATION_ASSIGNMENT,
  initial requisition (requisition_type 'tor_base'), equipment find-or-create by name,
  requisition items. Single transaction (state whether source-stated or technical).
  Notification: sendNewContractNotification via the notifier interface.
- 1A: list pending_survey requisitions, validate existing quantities,
  to_buy_qty = GREATEST(required_qty - existing_qty, 0), update status per use case.
Tests: table-driven for to_buy_qty and validations; integration tests only on a
disposable test DB. Verify before proceeding, run build/vet/test, report, do not commit.
```

---

## Prompt 5 — Slice 2: procurement and funding (2S, 2A, 3S, 3.1S)

```
Plan first, then implement after approval: 2S, 2A, 3S, 3.1S.

Use the corrected references:
- 2A reads equipment_name through REQUISITION_ITEM → EQUIPMENT (equipment_id).
- Expense claims link by requisition_id; multiple claims per requisition are allowed.
- Procurement is cumulative (actual_qty = actual_qty + new_qty); expense_type is
  'actual_expense' for 2A and 'fund_transfer' for 2S; duplicate transfer_ref_no rejected.
- Receipt photo via the storage interface.
D4 and D5 are open: implement only what the SDs state for 3S/2S and mark the
7A-dependent behaviour TODO(decision-4)/TODO(decision-5). Do not invent approval stages.
Report which statuses each step writes. Run build/vet/test. Do not commit.
```

---

## Prompt 6 — Slice 3: scheduling and leave (3A, 4A, 1W, 2W)

```
Plan first, then implement after approval: 3A, 4A and Worker endpoints 1W, 2W.

Preserve: schedule validation, worker availability, no duplicate schedule per
worker/day, area-permission check (D12: implement what the use cases state and the
existing model supports; flag the rest), recheck latest data inside the transaction on
confirm, notifyAssignedWorkers; leave: duplicate-leave check, is_advance_notice
calculation, pending leave, replacement candidates, replacement schedule, notifications
to assistant / leaving worker / replacement worker.
Identity: user_id = worker_id (docs/04 §1.1). Worker endpoints resolve the LINE
identity to a linked USER server-side and only return that user's own data.
Verify against the schema first, run build/vet/test, do not commit.
```

---

## Prompt 7 — Slice 4: attendance and evidence (3W, 4W, 4S)

```
Plan first. Before coding, state how you will handle D2 (evidence link) and D11
(attendance UPSERT key / replacement flow). If the schema does not settle them, stop
and ask me.

When cleared, implement:
- 3W substitute check-in: validate LINE user and schedule, today's shift matches worker,
  duplicate check-in rejected, UPSERT attendance per the agreed key.
- 4W work report: exactly one open attendance, description + photo validation, save
  evidence, set check-out, in one transaction.
- 4S daily attendance job: classify on-time / late / leave / absent, count approved
  advance leave, safe to re-run.
Tests for classification boundaries and duplicate/re-run behaviour. Build/vet/test.
Do not commit.
```

---

## Prompt 8 — Slice 5: payroll and payslip (5S, 7W)

```
Plan first. Blocked in part by D1, D3, D8, D9, D10. Implement only what is settled and
flag the rest.

Use canonical, corrected references:
- createDeduction(workerId, attendanceId, penaltyAmount, penaltyReason) inserts
  worker_id, attendance_id, penalty_amount, reason.
- Payroll insert: worker_id, period_start, period_end, base_wage, total_deduction,
  net_wage, is_paid=false. net_wage = base_wage - total_deduction (do not change the
  formula or negative-pay treatment; D10 is open).
- Penalty tiers per docs/05_business_rules.md. Re-running 5S must not create duplicate payroll or
  deduction rows; propose the guard in the plan and wait for approval (D8).
- 7W payslip: base_wage, total_deduction, net_wage, is_paid; deductions with their work
  date via attendance_id; penalty_amount. Own data only.
Money uses decimal, never float. Build/vet/test. Do not commit.
```

---

## Prompt 9 — Slice 6: extra equipment (5W, 5A, 6A, 7A, 8A)

```
Plan first. D2, D4 and D5 affect this slice; settle or stub them explicitly.

Implement what is specified:
- 5W creates additional requests with status 'pending_survey' so they appear in 5A/6A.
- 5A/6A list, review, forward or reject; notify per use case.
- 7A cumulative purchase with row locks and duplicate-submission prevention,
  'actual_expense' claim linked by requisition_id, status from saved quantities.
- 8A delivery evidence and notifications.
Acceptance: a 5W request appears in the 5A queue; no double purchase on concurrent
submit; claims link to the original requisition. Build/vet/test. Do not commit.
```

---

## Prompt 10 — Slice 7: financial report (6S)

```
Plan first. Blocked by D6 and D7; do not guess either. Propose the options with their
effect on totals and wait for my decision.

When cleared, implement 6S:
- Revenue from COMPANY_INVOICE per the agreed scope.
- Labor = SUM(net_wage) of paid payroll within the agreed project/period scope.
- Material = actual_expense claims joined by requisition_id for requisitions in
  approved, pending_supervisor, pending_procurement, completed. Fund transfers excluded.
- Net profit = revenue - (labor + material).
- Export per the use case.
Verify with a hand-calculated sample fixture, including a case with a fund_transfer and
an actual_expense on the same requisition (must not double count). Build/vet/test.
Do not commit.
```

---

## Prompt 11 — Notifications and jobs wiring

```
Plan first. Wire the notifier interface to the real LINE implementation and the job
runner for 4S and 5S.

- Use CURRENT official LINE documentation for message and token APIs; do not rely on memory.
- Use only supplied configuration. Never print or commit secrets.
- Keep the no-op notifier for tests; ensure results distinguish "sent" from "skipped/failed".
- Recipient must match the use case (e.g. original requester, leaving worker, replacement).
Report anything you could not verify against a real LINE channel.
```

---

## Prompt 12 — Review pass (READ-ONLY)

```
Review all uncommitted changes against CLAUDE.md and the use cases. Do not edit.
Check: canonical names, no float money, parameterised SQL, permissions, transaction
boundaries, idempotency, unresolved decisions stubbed (not guessed), no secrets,
no unrelated changes, tests meaningful. Report findings by severity with file:line,
and list what remains unverified. Confirm nothing was committed.
```

---

## Stop-and-ask template (use anytime)

```
Pause. Do not edit anything. Tell me: (1) what you inspected, (2) what conflicts with
CLAUDE.md or the schema, (3) which D# this touches, (4) the options and their impact.
Wait for my decision.
```
