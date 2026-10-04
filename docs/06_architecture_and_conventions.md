# 06 — Architecture and conventions (PROPOSAL; adapt to what already exists in the repo)

Reuse existing structure. Do not restructure working code.

## Layers
```
handler (Gin) → service (rules, transactions) → repository (SQL) → PostgreSQL
                     ↘ notifier (LINE interface)   ↘ storage (Supabase Storage interface)
jobs (scheduler) → service
```
- Handlers: bind/validate input, call one service method, map errors to HTTP. No SQL and no business rules.
- Services: business rules (penalty tiers, to_buy_qty, status transitions, permission checks) and transaction boundaries.
- Repositories: one named function per SD operation (see `03_workflows_operations.md`), parameterized SQL only. These replace `DB Connector.executeQuery`.
- Notifier: interface with a real LINE implementation and a clearly marked no-op for tests. A no-op is never reported as delivered.
- Jobs: 4S daily attendance and 5S payroll are service functions called by a scheduler; re-runnable and idempotent.

The SD participants (controllers/repositories) are design responsibilities, not one-file-per-class. UI click messages become endpoints or frontend events. Similar names stay separate unless mapped (PayrollSummary Controller 5S vs Payroll Controller; Contract Repository vs TOR Repository).

## Suggested modules
`contract`, `funding`, `requisition`, `schedule`, `leave`, `attendance`, `payroll`, `report`, `evidence`, `worker`, plus `auth`, `notify`, `storage`, `jobs`.

## Workflow → module → interface
| Use case | Module | Interface | Notes |
|---|---|---|---|
| 1S | contract | Supervisor web | transaction across contract, location, assignment, requisition, items |
| 2S | funding | Supervisor web | fund_transfer claim, unique transfer ref |
| 3S, 3.1S | requisition | Supervisor web | approve/reject + notification |
| 4S | attendance | job | batch classification and UPSERT |
| 5S | payroll | job | deductions + payroll, idempotent |
| 6S | report | Supervisor web | revenue − labor − material; PDF export |
| 1A | requisition | Assistant web | to_buy_qty |
| 2A | requisition/procurement | Assistant web | cumulative qty, actual_expense, receipt |
| 3A | schedule | Assistant web | no duplicates, recheck |
| 4A | leave | Assistant web | decision + replacement |
| 5A, 6A, 7A, 8A | requisition (additional) | Assistant web | queue, forward/reject, purchase, delivery evidence |
| 9A | contract | Assistant web | continuation check |
| 1W | schedule | Worker LIFF | own schedule |
| 2W | leave | Worker LIFF | |
| 3W | attendance | Worker LIFF | substitute check-in |
| 4W | evidence/attendance | Worker LIFF | evidence + check-out |
| 5W | requisition (additional) | Worker LIFF | creates `pending_survey` |
| 6W | requisition | Worker LIFF | read result |
| 7W | payroll | Worker LIFF | own payslip |

## Coding conventions
- Money: `NUMERIC` in Postgres, `shopspring/decimal` (or the repo's existing choice) in Go. Never float for money.
- Time: store `timestamptz` UTC; day logic uses Asia/Bangkok explicitly.
- IDs: follow the physical schema.
- Transactions: wrap multi-write workflows; state whether the boundary is source-stated or added for technical atomicity. Re-read inside the transaction (`FOR UPDATE` where needed) for "recheck latest / prevent duplicate" steps.
- Errors: typed domain errors → consistent JSON error shape; list failing fields on validation errors; no stack traces to clients.
- Logging: structured; no secrets, tokens or personal data dumps.
- Auth: reuse the repo's existing web auth (D13). Worker endpoints verify LINE identity server-side.
- Authorization: role check plus assignment/area check where the use case requires it (D12).
- Storage: upload through Supabase Storage with server-side type/size validation; store only the URL/path.
- Supabase: be careful with Row Level Security and service-role keys; never expose the service key to clients.
- Tests: table-driven unit tests for pure rules; integration tests only against a disposable test database.
- Keep functions small; no speculative abstractions; no unrelated refactors; no dead code.
