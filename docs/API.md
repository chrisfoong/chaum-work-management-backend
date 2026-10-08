# API notes

api/openapi.yaml lists actual routes/request models (JSON is valid YAML1.2). Regenerate: go run ./cmd/apidoc. Documentation tooling is never an auth bypass in the server.

Web /api/web uses verified LINE Web tokens and Supervisor/Assistant; Mini App /api/liff uses Worker channel tokens and DB Worker role. Send Bearer authorization. Client IDs/roles never establish identity.

| Workflow | Permission |
| --- | --- |
| TOR preview/confirm, user/catalog writes, procurement review, funding, payroll, invoices, profit | Supervisor |
| Schedules, leave review/replacement, survey, inspection/decision, purchase/delivery | Assistant |
| Staff lists/dashboard, assignment QR, evidence | Supervisor/Assistant |
| Own schedules, leave, attendance, requisitions, payslips | Worker |

Lists: limit1-100(default50), offset0-100000(default0); response data/limit/offset. Dates YYYY-MM-DD; months YYYY-MM; money decimal strings; IDs UUID. Payroll periods must match1-15 or16-month end.

Check-in requires schedule_id, qr_token, latitude, longitude and accuracy. Staff assignment QR expires60s and is reusable across workers; one check-in/schedule prevents repeat check-in. All GPS fields required. Check-in is allowed on the scheduled Bangkok calendar day (including before planned start) until the end of the eight-hour shift; overnight continuation remains supported. Early check-in is on_time. Checkout needs owned evidence and description.

POST /files accepts raw body with Content-Type, not multipart, max5MiB. Worker JPEG/PNG; staff also PDF. Use returned object path, not an external URL. GET /files?path=... enforces access and proxies private storage.

Errors have request IDs and sanitized internals. Business success does not promise notification delivery; missing QR/Storage explicitly blocks affected requests.

## Use Case Description reconciliation

- Assistant POST /leave-requests/:id/review: {status: approved|rejected, replacement_worker_id?: UUID}. Review and optional replacement are atomic; Supervisor is denied. Future-date leave is advance, same-day is emergency. 4S exempts only approved advance notice from absence.
- Scheduling requires active contract and completed initial TOR equipment; required_workers is a minimum. Worker GET /schedules lists own upcoming work with stored schedule statuses, ascending.
- Assistant GET /requisitions/:id/inspection returns TOR requirements and other pending requests. POST /requisitions/:id/decision: {decision: purchase|no_purchase, reason: string}. Additional requests use this explicit decision; /survey is base only. Worker requests require active attendance.
- Purchase items include item_id, actual_qty (this round), actual_price and required expected_actual_qty (current cumulative quantity). 7A additional purchases require at least one positive quantity; 2A base acquisition may report all-zero quantities and escalate missing supplies. Zero prices allowed. Positive cost requires receipt_path; zero cost returns empty expense_id. Partial completion -> pending_fund; funding -> pending_procurement. 2A latest price; 7A weighted average rounded to two decimals. Stale/repeated confirmations return 409 without duplicating quantity/expense. Receipt objects cannot be reused for another purchase. Additional receipts are photos.
- Assistant POST /requisitions/:id/delivery: {schedule_id, description, photo_paths: [owned JPEG/PNG paths]}. Request must be completed/additional and recipient must be original requester in matching assignment. Stores all photos atomically; matching replay succeeds, changed replay conflicts. GET /requisitions/:id/deliveries lists these evidence rows. Description including generated summary <=2000 chars; 1-20 photos. No persistent delivery status/FK is claimed.
- Supervisor POST /payroll/batch: {period_start, period_end}. All workers with scheduled work in period; atomic batch, repeat reuses existing exact-period slips. Worker GET /payroll?period_month=YYYY-MM shows own paid slips. Default worker list also paid only.
- Profit JSON/CSV/PDF accepts month=YYYY-MM or period_start/period_end, optional tor_id. Uses paid payroll periods fully inside dates; paid invoice net_received by included billing months; expenses by Bangkok purchase date. Allocation is estimated. PDF uses English financial fields/TOR UUIDs; JSON/CSV preserves project names.
- Supervisor POST /reports/profit/confirm: {tor_id, period_start, period_end} returns snapshot, confirmed_at, persisted=false. Does not change contract status. Assistant GET /contracts/:id/continuation returns assignments for active unexpired TOR.

Notifications are best effort. No Worker delivery notification is emitted during purchase. Persistent closing history/outbox/item-round history require a separately approved persistence design.

## Backend contract for the final frontend guide

- POST /api/web/contracts/confirm requires contract.contract_file_path from POST /api/web/files (PNG bytes <=5MiB). Persists the verified private object path in contract_file_url. Never submit an arbitrary external URL. /contracts/info previews textual formats; final confirmation verifies the uploaded PNG.
- Purchase receipts for both 2A/7A: JPEG/PNG when positive total; funding receipts retain JPEG/PNG/PDF. Raw upload stays <=5MiB; no multipart.
- GET /api/web/assignments?tor_id=UUID&limit=50&offset=0 supplies assignment, contract, location, required_workers and can_schedule.
- GET /api/web/workers/available?work_date=YYYY-MM-DD lists valid free workers for scheduling. Assistant only.
- GET /api/web/leave-requests/:id gives affected schedule/location/project. GET /api/web/leave-requests/:id/candidates supplies free candidates excluding the leaver. Assistant only. Lists remain data/limit/offset; candidates are not a reservation and are rechecked during submission.
- GET /schedules, /leave-requests, /attendance, /requisitions in their authorized web/liff groups accept status, assignment_id, tor_id, period_start, period_end, limit and offset; requisitions also accepts requisition_type. Filters run before pagination. Dates on requisitions mean creation date in Bangkok. Worker upcoming-date restriction/ownership still applies; stored completed/cancelled statuses are visible.
- Schedule cards include project_name, contract_no, worker names, shift_start_at and shift_end_at as timestamps; shift_end_at=start+8h, possibly next date. Worker check-in still requires schedule_id, latitude, longitude, accuracy_m, qr_token; there is no advance accept/reject endpoint.
- GET /api/web/requisitions/:id/delivery-schedules lists only original requester's matching area schedules dated today or earlier. Assistant only. POST delivery also rechecks this rule and full procurement inside transaction. Guide payload remains schedule_id/description/photo_paths; all items are delivered together.
- Staff GET /contracts adds workflow_status (registered/active/ended) and can_operate, without changing stored enum values. Staff dashboard includes latest 100 contracts; use paginated /contracts for the complete list.
- /requisitions/:id/decision purchase reason is appended to stored reason for Supervisor review; no_purchase explanation is sent to Worker without persisting, per Use Case.
- 6W uses LINE Chat. Delivery notifications include requisition number/equipment/quantity, no Worker notification at purchase. Missing config/network failure remains delivery-unverified; no inbox/replacement-acceptance/paid_at is invented.

Frontend process payroll by calling Supervisor /attendance/finalize followed by /payroll/batch after period/cutoff closure. Retrying batch reuses existing slips; mark-paid records confirmation only and does not execute a bank transfer.

## Use Case audit completion

See USECASE_AUDIT.md for all 22 descriptions, function names and honest limits.

- Worker schedule cards include nullable LOCATION.address and actual stored future schedule statuses.
- Worker shortage reasons <=500. Inspection accepts only pending_survey additional requests. Supervisor approval accepts only additional requests.
- no_purchase writes Assistant reviewed_by/reviewed_at, preserves original reason, and may decline inactive equipment. The explanation is not persisted; resupply it for manual notification retry.
- Initial approved 7A purchases need no fund_transfer record. A partial purchase moves to pending_fund; only Supervisor funding moves it back to pending_procurement. 7A requires a positive quantity. 2A permits an all-zero acquisition to escalate without expense or quantity writes.
- Profit JSON adds labor_details (payroll_id/workdays/allocated_net_wage) and material_details (expense_id/requisition_id/amount/created_at); sums use the same filters and exact cent allocation as totals.
- Assistant GET /assignments/:id/continuation accepts an assignment UUID, returns contract/location/minimum/dates/status/can_continue/notice. Ended contracts are visible with a contact-Supervisor notice. summary_notification_history_available=false means history is unavailable, not no summary has been sent.
- Decision and delivery POST success returns {business_saved:true,notification:{status,accepted_recipients,failed_recipients,delivery_verified:false,persisted:false}}. Status is accepted/partial/failed/disabled. LINE failure does not roll back saved business data. Matching delivery replay reuses evidence and stable notification key.
- Assistant POST /notifications/:id/retry: {kind,reason?}. Kinds: no_purchase, delivery, approval_needed, purchase_funding, leave_review, replacement_needed. id is the matching requisition UUID or leave-request UUID. Guards require the committed state; no_purchase additionally requires the original reviewing Assistant and a reason <=1000. No arbitrary recipient or free message API. Resend reads only and does not create quantity/expense/evidence/review rows.
- LINE transport retries transient failures at most three times under a bounded timeout, retaining the same retry key/body. Same event/recipient/chunk/content uses a stable key within LINE's deduplication window. Acceptance is not recipient delivery; changed content is a new message. No durable outbox.

Assistant POST /api/web/requisitions/:id/purchase/preview uses PurchaseInput and returns items(new_actual_qty/remaining_qty/new_actual_price/round_cost), round_total, next_status and persisted=false. It runs confirmation validation without quantity/expense/status writes. Use for 7A review; POST /purchase rechecks current data when confirming.
