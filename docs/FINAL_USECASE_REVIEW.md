# Backend final Use Case review — 2026-10-09

Reviewed branch: `feature/assistant-operations-summary` in `C:\SA\chaum-work-management-backend`. Existing uncommitted TASKS.md changes preserved. No runtime Go code changes were needed during this verification.

Source: [SA-Group5](https://docs.google.com/document/d/1Ovggbw2PwWJJdyzzkrCMP1CruIoPelmZ6f9whHYZryQ/edit?tab=t.0), freshly exported through the browser on 2026-10-09. Compared substantive Use Case Descriptions with REQUIREMENTS.md, USECASE_AUDIT.md, registered routes and implementation/tests. Illustrative SQL and old enum labels do not override Supabase exports. Latest user-confirmed 9A overrides automatic Assistant summary/history. Frontend remains a separate repository.

## Coverage by Use Case

All 22 have Backend flows. This is implementation and local verification coverage, not proof of every scenario or live external integration.

| Case | Backend behavior / functions | Verification and remaining limits |
|---|---|---|
| 1S | SubmitContractInfo, SubmitScopeData, ConfirmContract: TOR validation, private PNG verification, areas/equipment and transaction | Contract validation, role, file, concurrent confirmation and rollback tests pass. Private Supabase Storage live workflow pending. |
| 2S | Fund: funding amount, receipt/reference, replay guard and procurement transition | Isolated procurement flow passes; records funding, never executes bank transfer. |
| 3S | ReviewRequest: Supervisor approves/rejects additional procurement | Role/state/validation tests pass; actual notification delivery pending. |
| 4S | Finalize: on_time/late/absent/approved advance leave after shift end +2h | Attendance concurrency, leave rules, penalty boundaries and overnight tests pass. |
| 5S | PreviewPayroll, Payroll, PayrollBatch, Pay: half-month batch, exact money, capped deduction and paid flag | Batch/retry/rollback tests pass. Worker payslip notices remain; Assistant operational summary is manual 9A. No bank transfer occurs. |
| 6S | Invoice, Receive, ProfitRange, CloseSummary, ExportProfit/PDF: paid revenue/labor, actual expense, workday allocation | Paid-report, allocation/detail reconciliation and PDF tests pass. Confirmation is persisted=false; Thai PDF typography remains limited. |
| 1A | Survey: base request survey and computed shortage | Isolated survey/procurement tests pass; generated to_buy_qty is not written. |
| 2A | PreviewPurchase, Purchase: partial/repeated base acquisition and latest price | Zero-cost/all-zero and partial procurement tests pass; no durable item-round ledger. |
| 3A | AvailableWorkers, Schedule: active contract/readiness/minimum workers, one worker-day | Candidate filters, duplicates, transactional rechecks and concurrency tests pass. |
| 4A | LeaveDetail, LeaveCandidates, ReviewLeaveAndReplace, Replacement | Assistant approval, advance/emergency, candidate exclusion and atomic replacement tests pass; no advance accept/reject flow. |
| 5A | InspectRequest: existing additional request versus TOR/pending requests | State/type/inspection integration tests pass. |
| 6A | DecideRequest, RetryNotification: purchase/no_purchase on original request | Reviewer/state, original reason, notification-only retry tests pass; no_purchase explanation not persisted. |
| 7A | PreviewPurchase, Purchase: initial approved acquisition and weighted repeated purchases | Initial approval, partial funding, stale/replay and weighted-price tests pass. |
| 8A | DeliverySchedules, Deliver, Deliveries: full handover to original requester with photos | Recipient/area/date/quantity, atomic photos and matching retry tests pass. WORK_EVIDENCE convention, not a new relation. |
| 9A | OperationsSummary, NotifyOperationsSummary, AssignmentContinuation | Current operational data only, empty/ended messages, financial exclusion, Supervisor manual trigger and no-auto-summary tests pass. Actual LINE/dashboard delivery pending; no send/read/snapshot history required. |
| 1W | FilteredList/Schedules: own upcoming shifts with project/location and start/end | Ownership/read-model tests pass; live authenticated empty-list read 200. |
| 2W | Leave: assigned current/future date, reason and advance/emergency | Leave validation/integration tests pass; live write not exercised. |
| 3W | QR, CheckIn: both GPS <=200m and signed area QR, server time, eight-hour shift | QR tamper/expiry/area, GPS accuracy, ownership and duplicate/concurrent check-in tests pass. Physical GPS/QR end-to-end pending. |
| 4W | CheckOut, file upload/verification: owned photo, description and atomic checkout | File transport/ownership and attendance rollback tests pass. Live upload/checkout pending. |
| 5W | Requisition: additional equipment during active attendance | Active attendance, quantities/reason and request transaction tests pass; live write pending. |
| 6W | LinePush.SendKey, delivery/no_purchase notices and RetryNotification | HTTP transport, bounded retry keys and no duplicate evidence tests pass. Actual LINE recipient delivery pending. |
| 7W | PayrollMonth/Payrolls: own paid slips, month and deduction detail | Paid/ownership/period tests pass; real authenticated empty payroll read 200. |

## Fresh checks

- gofmt -l tracked Go files: no unformatted files.
- go test -count=1 -json ./... with coverage: **233 tests/subtests passed, 0 failed, 0 skipped**; 16 test packages passed.
- go vet ./...: passed.
- go build ./...: passed.
- git diff --check: passed (Git emitted a line-ending conversion advisory for TASKS.md).
- Statement coverage measured this run: **63.1%**. Passing tests do not imply exhaustive coverage.

Integration used newly created disposable `chaum_final_audit_20261009` on isolated PostgreSQL17 at 127.0.0.1:55439. Test schema was applied there only. Existing Supabase was not used for automated mutation tests. Evidence in the task workspace: backend-final-tests.jsonl and backend-final-coverage.out.

## Real LINE and shared-database cleanup

Earlier live read suites: Supervisor **16/16**, Assistant **16/16**, Worker **11/11** after user-authorized creation of one Worker mapping. Most lists were empty, so those suites establish authentication/routes/query availability, not populated business acceptance.

User subsequently requested cleanup. Deleted only test worker_id `fb3629c2-65cd-48dc-b72f-c802a0ef7bda` linked to target user_id `11111111-1111-1111-1111-111111111111`. Transaction first locked the target records and checked all Worker FK references; related business records would abort cleanup. Verified final target state: **role=supervisor, is_active=true, worker_rows=0**. No other account or business record was removed. Worker success is historical evidence; the restored Supervisor account should now be denied Worker APIs.

No .env access, Supabase DDL/migrations/seeds, server/job startup, new push, PR, merge or deployment in this round.

## Acceptance limits

No missing Backend flow was identified against the 22 reviewed descriptions plus confirmed user overrides. Live private Storage/upload, physical GPS+QR check-in/out, actual LINE equipment/9A delivery and frontend Dashboard navigation remain unverified. Previously observed missing QR signing/messaging configuration must be resolved by the operator; current secret values were not inspected. Durable outbox, send/read history, item-round ledger and closing ledger are not implemented; fixed-schema MVP limits remain explicit. Thai PDF text rendering remains a limitation.
