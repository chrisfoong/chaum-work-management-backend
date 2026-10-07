# API notes

api/openapi.yaml lists actual routes/request models (JSON is valid YAML1.2). Regenerate: go run ./cmd/apidoc. Documentation tooling is never an auth bypass in the server.

Web /api/web uses verified LINE Web tokens and Supervisor/Assistant; Mini App /api/liff uses Worker channel tokens and DB Worker role. Send Bearer authorization. Client IDs/roles never establish identity.

| Workflow | Permission |
| --- | --- |
| TOR preview/confirm, user/catalog writes, reviews, funding, payroll, invoices, profit | Supervisor |
| Schedules, replacement, survey/purchase | Assistant |
| Staff lists/dashboard, assignment QR, evidence | Supervisor/Assistant |
| Own schedules, leave, attendance, requisitions, payslips | Worker |

Lists: limit1-100(default50), offset0-100000(default0); response data/limit/offset. Dates YYYY-MM-DD; months YYYY-MM; money decimal strings; IDs UUID. Payroll periods must match1-15 or16-month end.

Check-in requires schedule_id, qr_token, latitude, longitude and accuracy. Staff assignment QR expires60s and is reusable across workers; one check-in/schedule prevents repeat check-in. All GPS fields required. Check-in must occur during eight-hour shift; no early window/grace. Checkout needs owned evidence and description.

POST /files accepts raw body with Content-Type, not multipart, max5MiB. Worker JPEG/PNG; staff also PDF. Use returned object path, not an external URL. GET /files?path=... enforces access and proxies private storage.

Errors have request IDs and sanitized internals. Business success does not promise notification delivery; missing QR/Storage explicitly blocks affected requests.
