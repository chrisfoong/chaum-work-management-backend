# Backend Use Case audit — 2026-10-08

Business source: [SA-Group5 Use Case Descriptions](https://docs.google.com/document/d/1Ovggbw2PwWJJdyzzkrCMP1CruIoPelmZ6f9whHYZryQ/edit?tab=t.0). The requested tab `t.5fr53uyz81w0` is **5.2 QUERY TABLE**; the substantive descriptions are in the SA tab. The separate Use Case template tab is empty. Read descriptions as requirements; illustrative SQL does not override exported UUIDs, enums or constraints.

Existing develop Foundation and Supabase schema are preserved. This audit adds `feature/usecase-completion` on top of `feature/backend-verification`. Backend only. No `.env` inspection, schema changes, production test writes, merge or deployment.

## Functions and coverage

“Implemented” means the Backend flow exists and is covered by local tests, not that all scenarios or deployed integrations have been verified. Shared prerequisites: server-verified LINE identity, database roles/active accounts, ownership, exact money, transaction guards and schema export mapping.

| Use case | Backend functions | API (platform prefix omitted) | Result and limits |
|---|---|---|---|
| 1S TOR | SubmitContractInfo, SubmitScopeData, ConfirmContract | POST contracts/info, scope, confirm; files | Implemented: 10 digits, dates/value, private PNG bytes/size/owner, areas/equipment, atomic confirmation. |
| 2S funding | Fund | POST requisitions/:id/fund-transfers | Implemented: positive exact amount, unique reference/replay, receipt, transactional state. |
| 3S approval | ReviewRequest | POST requisitions/:id/review | Implemented: Supervisor, pending additional request only, reject reason, review metadata, background LINE. |
| 4S attendance | Finalize | POST attendance/finalize; scheduled job | Implemented: end+2h, on_time/late/absent, approved advance leave exemption, eight-hour overnight shifts. |
| 5S wages | PreviewPayroll, Payroll, PayrollBatch, Pay | POST payroll/preview, payroll/batch; GET payroll; mark-paid | Implemented: periods, exact wages/penalty cap, duplicate/retry guards, post-commit project/area summary notification. mark-paid confirms payment; no bank transfer executed. |
| 6S finance | Invoice, Receive, Profit, ProfitRange, CloseSummary, ExportProfit, ExportProfitPDF | invoices; reports/profit, .csv, .pdf, /confirm | Implemented: paid revenue/labor, actual expenses only, workday/cent allocation and labor/material breakdown. Confirmation snapshot has persisted=false. PDF uses English/UUIDs; Thai typography remains limited. |
| 1A survey | Survey | POST requisitions/:id/survey | Implemented: all base items, existing quantities; generated to_buy_qty never written. |
| 2A acquisition | PreviewPurchase, Purchase | POST requisitions/:id/purchase | Implemented: cumulative quantity, latest unit price, receipts for positive cost, partial funding. All-zero acquisition can escalate without creating an expense. |
| 3A scheduling | AvailableWorkers, Schedule | workers/available; POST schedules | Implemented: active contract/readiness/minimum staffing, free active workers, worker-day uniqueness and transactional rechecks. |
| 4A leave/replacement | LeaveDetail, LeaveCandidates, ReviewLeaveAndReplace, Replacement | leave-requests/:id, candidates, review, replacement | Implemented: Assistant review, future advance/same-day emergency, optional atomic replacement/rollback, notification-only retry. No advance accept/reject flow. |
| 5A inspect | InspectRequest | GET requisitions/:id/inspection | Implemented: pending additional only, original items, TOR requirements and other pending requests. No writes. |
| 6A decision | DecideRequest, RetryNotification | POST requisitions/:id/decision; notifications/:id/retry | Implemented: purchase/no_purchase, inactive equipment rejected for purchase but may be declined, no_purchase review metadata and original reason retained. HTTP separates business_saved from LINE acceptance. Explanation must be supplied again for retry because description says not to persist it. |
| 7A additional purchase | PreviewPurchase, Purchase | POST requisitions/:id/purchase | Implemented: initial Supervisor approval permits first purchase; partial rounds require funding state, positive quantity, stale/replay guards, weighted unit price, exact per-round expenses, no premature delivered notice. |
| 8A delivery | DeliverySchedules, Deliver, Deliveries, RetryNotification | requisitions/:id/delivery-schedules, /delivery, /deliveries | Implemented: complete quantities/original recipient/area/date, atomic owned photos, matching replay, item summary. LINE failure/resend does not create evidence again. Delivery uses WORK_EVIDENCE convention; no invented relation/status. |
| 9A operations | OperationsSummary, NotifyOperationsSummary, AssignmentContinuation | assignments/:id/operations-summary; contracts/:id/operations-summary/notify; assignments/:id/continuation | Implemented confirmed MVP: current attendance/leave/absence and procurement progress only; manual Supervisor LINE link after 5S/6S review; no financials or persisted snapshot/send/read history. Real LINE delivery remains unverified. |
| 1W schedules | FilteredList, Schedules | GET schedules | Implemented: own future schedules, project/location/address/start/end and stored scheduled/completed/cancelled status. Check-in separately rejects invalid states. |
| 2W leave | Leave | POST leave-requests | Implemented: own assigned date, reason <=500, no past/duplicates, advance/emergency, notify Assistant. |
| 3W check-in | QR, CheckIn | POST attendance/check-in | Implemented: both GPS <=200m and signed/expiring area QR, server Bangkok time, schedule ownership/unique attendance. Deployed QR/GPS verification pending. |
| 4W checkout | CheckOut, Storage.Upload/Verify | files; POST attendance/check-out | Implemented: checked-in/not checked-out, description <=1000, owned actual JPEG/PNG, atomic evidence/time. |
| 5W shortage | Requisition | equipment; POST requisitions | Implemented: current active attendance/assigned area, active equipment, positive quantities, reason <=500, pending survey and notification. |
| 6W result | RetryNotification, LinePush.SendKey | LINE Chat; staff notification retry | Implemented transport/flow: delivery summary only after real handover or no_purchase explanation; no Worker acknowledgement page. Actual LINE delivery remains unverified. |
| 7W payslip | PayrollMonth, Payrolls | GET payroll?period_month=YYYY-MM | Implemented: own paid slips, base/deductions/date/reason/net, exact period mapping through payroll.user_id. No invented paid_at. |

## Gaps corrected in this feature

- Added actual LOCATION.address to schedule cards and allowed future stored statuses instead of silently hiding completed/cancelled cards.
- Enforced shortage reason <=500. Restricted inspection to pending additional requests and Supervisor approval to additional requests.
- Corrected no_purchase: inactive equipment can be declined; reviewed_by/reviewed_at are written while original Worker reason remains unchanged. Purchase still requires active equipment.
- Separated 2A all-zero acquisition from 7A positive-quantity requirement. Removed unintended funding prerequisite on the initial approved 7A purchase; partial funding transition remains enforced.
- Added labor/material detail arrays using the same exact allocation and filters as totals. No bank information included.
- Added assignment continuation read model. The initial post-payroll Assistant summary was superseded by the confirmed manual 9A command. No automatic schedule or contract mutation.
- Added notification-only retry with role/state/actor guards, synchronous decision/delivery status, bounded LINE retries and stable event/message retry keys. No business write is repeated.

LINE retries follow [official retry-key guidance](https://developers.line.biz/en/docs/messaging-api/retrying-api-request/). An accepted request (including keyed 409 with accepted-request header) is acceptance, not proof of recipient delivery. Keys are reused for the same event/recipient/chunk/content; LINE's deduplication window is finite. Manual resend with changed content is a new message. No persistent delivery log or restart recovery queue is claimed.

## Remaining work that cannot be called complete

- Real LINE read suites now passed: Supervisor 16/16, Assistant 16/16, Worker 11/11. Private Storage, actual LINE Chat delivery and QR/GPS end-to-end still need operator verification. See FINAL_USECASE_REVIEW.md for current evidence and test-data cleanup.
- 9A persisted summary/send/read history and acknowledgement are explicitly outside the confirmed MVP, not a completion blocker. Current-data reports may change after retrospective edits. No new tables are required. Legacy continuation history_available=false is not a claim that no notice was sent.
- Durable notification scheduling/recovery after restart, closing ledger, item-by-item purchase history, wage snapshots and replacement relations are absent schema facilities. Retry APIs cover immediate/manual resends without replaying business writes.
- Thai text in PDF remains limited. JSON/CSV preserve names; PDF identifies TOR by UUID.

Tests: final measured results are recorded in TASKS.md and VERIFICATION_AND_PUBLISH.md. Do not infer 100% test coverage from all tests passing. Further pushes are allowed only after this audit and final checks; no deploy/merge.

7A review step: Assistant POST requisitions/:id/purchase/preview uses the same validation/state/expected-quantity/receipt checks as confirmation, returns round_total/new quantities/remaining/weighted price/next_status with persisted=false, and makes no business writes. Confirmation always checks current state again; preview is not a reservation.
