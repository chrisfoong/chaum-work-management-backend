# 07 — Verification checklist and acceptance checks

## A. Before any new work (report as a table: expected / found / match?)
1. `git status`, `git log --oneline -n 10`, current branch. Confirm nothing will be overwritten.
2. Go version, module path, dependencies, Gin version.
3. Directory tree, entry point, config loading, `.env.example` (placeholders only).
4. Migrations/schema: for each of the 16 entities — table exists? columns match `02_domain_model.md`? types? PK/FK/unique?
5. Identity: one key (`user_id = worker_id`) or two columns?
6. Presence of: `payroll.total_deduction`, `payroll.net_wage`, `payroll.is_paid`, `deduction_transaction.attendance_id`, `deduction_transaction.penalty_amount`, `expense_claim.requisition_id`, `expense_claim.expense_type`, `equipment_requisition.requisition_type/status/reviewed_by/reviewed_at/reject_reason`, `leave_request` user key.
7. Status value sets actually used.
8. Existing auth/permission code. Reuse; do not add sign-in methods.
9. Baseline `go build ./...`, `go vet ./...`, `go test ./...` before changes.
10. Anything in the repo that contradicts the docs.

## B. Acceptance checks (derive concrete tests from these)
- Role and area access, preconditions, required-field/format/database validation and error behavior match each use case.
- 5W requests appear in the 5A/6A queue (`pending_survey`).
- 2A equipment names come from EQUIPMENT; quantities and procurement updates are correct and cumulative.
- Expense claims link to the original requisition; multiple claims allowed.
- Completed and partial actual purchases appear in material cost; fund transfers are not counted as purchase cost.
- Each 5S deduction carries the source attendance id and is retrievable with its work date in 7W.
- 5S writes and 7W reads base_wage, total_deduction, net_wage consistently; payment state uses is_paid.
- 6S labor uses net_wage; project/date scope validated once D6/D7 are decided.
- Worker, schedule, attendance, leave, evidence and payroll identities agree with user_id = worker_id and the real schema.
- No duplicate check-in; no conflicting schedule for the same worker and date.
- Notifications go to the correct recipient at the specified stage; no-op sends are not reported as delivered.
- Transactions and re-runs (4S, 5S) never leave partial or duplicate records.

## C. What the earlier verification did NOT prove
The class-diagram package verified diagram consistency only. It did not test any migration, application code, permission rule, integration or calculation. Do not treat it as application testing.

## D. Definition of done per slice
Plan approved; matches the docs (discrepancies reported); validation/errors match the use case; `go build ./...`, `go vet ./...` and relevant tests pass (state exactly what ran); no secrets, no unrelated changes, nothing committed; report what works, what was exercised, decisions touched, remaining gaps.

## E. Report format (end of every task)
**Inspected** (files actually read) · **Changed** (files) · **Ran** (commands + results) · **Decisions touched** (D#) · **Blocked / open** · **Uncommitted changes present** (yes).
