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

Check-in requires schedule_id, qr_token, latitude, longitude and accuracy. Staff assignment QR expires60s and is reusable across workers; one check-in/schedule prevents repeat check-in. All GPS fields required. Check-in must occur during eight-hour shift; no early window/grace. Checkout needs owned evidence and description.

POST /files accepts raw body with Content-Type, not multipart, max5MiB. Worker JPEG/PNG; staff also PDF. Use returned object path, not an external URL. GET /files?path=... enforces access and proxies private storage.

Errors have request IDs and sanitized internals. Business success does not promise notification delivery; missing QR/Storage explicitly blocks affected requests.

## Use Case Description reconciliation

- Assistant POST /leave-requests/:id/review: {status: approved|rejected, replacement_worker_id?: UUID}. Review and optional replacement are atomic; Supervisor is denied. Future-date leave is advance, same-day is emergency. 4S exempts only approved advance notice from absence.
- Scheduling requires active contract and completed initial TOR equipment; required_workers is a minimum. Worker GET /schedules lists own upcoming scheduled work, ascending.
- Assistant GET /requisitions/:id/inspection returns TOR requirements and other pending requests. POST /requisitions/:id/decision: {decision: purchase|no_purchase, reason: string}. Additional requests use this explicit decision; /survey is base only. Worker requests require active attendance.
- Purchase items include item_id, actual_qty (this round), actual_price and required expected_actual_qty (current cumulative quantity). At least one positive quantity; zero prices allowed. Positive cost requires receipt_path; zero cost returns empty expense_id. Partial completion -> pending_fund; funding -> pending_procurement. 2A latest price; 7A weighted average rounded to two decimals. Stale/repeated confirmations return 409 without duplicating quantity/expense. Receipt objects cannot be reused for another purchase. Additional receipts are photos.
- Assistant POST /requisitions/:id/delivery: {schedule_id, description, photo_paths: [owned JPEG/PNG paths]}. Request must be completed/additional and recipient must be original requester in matching assignment. Stores all photos atomically; matching replay succeeds, changed replay conflicts. GET /requisitions/:id/deliveries lists these evidence rows. Description including generated summary <=2000 chars; 1-20 photos. No persistent delivery status/FK is claimed.
- Supervisor POST /payroll/batch: {period_start, period_end}. All workers with scheduled work in period; atomic batch, repeat reuses existing exact-period slips. Worker GET /payroll?period_month=YYYY-MM shows own paid slips. Default worker list also paid only.
- Profit JSON/CSV/PDF accepts month=YYYY-MM or period_start/period_end, optional tor_id. Uses paid payroll periods fully inside dates; paid invoice net_received by included billing months; expenses by Bangkok purchase date. Allocation is estimated. PDF uses English financial fields/TOR UUIDs; JSON/CSV preserves project names.
- Supervisor POST /reports/profit/confirm: {tor_id, period_start, period_end} returns snapshot, confirmed_at, persisted=false. Does not change contract status. Assistant GET /contracts/:id/continuation returns assignments for active unexpired TOR.

Notifications are best effort. No Worker delivery notification is emitted during purchase. Persistent closing history/outbox/item-round history require a separately approved persistence design.
