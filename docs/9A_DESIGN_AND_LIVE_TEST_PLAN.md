# Confirmed 9A MVP and live acceptance plan

User confirmed on 2026-10-08: preserve the 16-table schema. The earlier two-table persistence proposal is superseded. No audit/outbox/summary/read tables, DDL, permanent snapshots, send/read history or acknowledgement are required.

Assistant reads current operational reports by area/date range: attendance, absence, leave and procurement progress. No wages/deductions/profit/prices/bank data. Current-data edits change reports. Existing all-area Assistant permissions remain in force.

Supervisor manually calls contracts/:id/operations-summary/notify after processing/reviewing 5S/6S. It checks payroll coverage and recomputes profit; prior human review is an explicit attestation, not a stored confirmation. It sends a short project/date LINE notice and HTTPS Dashboard link only. PayrollBatch and financial confirmation do not automatically send an Assistant summary. No write to business tables is needed for 9A.

Assistant assignments/:id/operations-summary returns has_data/summary_notice and can_continue/continuation_notice. No source records means “ยังไม่มีข้อมูลสรุปสำหรับรอบนี้”; ended/inactive/expired contract means “กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ”. UI scheduling is disabled and Backend guards remain enforced.

Report period basis: attendance by scheduled work_date; procurement uses current cumulative quantities for requests created in the selected Bangkok date range. This is not purchase-round history. LINE acceptance is distinct from delivery/read; no durable retry queue is claimed. Configure the real HTTPS Assistant Dashboard URL privately; never inspect .env.

## Live checks performed on 2026-10-08

The existing local server at 127.0.0.1:8080 was reachable. Five read-only requests returned expected backend JSON/status:

| Request | Expected/actual |
|---|---|
| GET /health/ready | 200 / 200 |
| GET /api/web/me, no bearer | 401 / 401 |
| GET /api/liff/me, no bearer | 401 / 401 |
| GET /api/web/me, deliberate invalid bearer | 401 / 401 |
| GET /api/liff/me, deliberate invalid bearer | 401 / 401 |

This confirms readiness and negative authentication behavior of the running process. It does not establish the intended Supabase project, current code revision, successful real LINE login or complete live integration. No identity/bank/secret responses were printed, no .env was inspected, no server was started/restarted and no business mutation was requested.

## Remaining live test sequence

1. Obtain Web login and Worker LIFF/MINI App URLs, plus the intended test environment. Human performs LINE login; password/2FA and tokens must not be sent into chat or logs.
2. Real Web Supervisor and Assistant /me: correct active role/account. Assistant cannot access Supervisor financial routes. Real Worker /me: exactly one mapping; own schedules and records only. Negative role/ownership requests must fail.
3. An explicitly isolated test Supabase project is required for mutations. Verify the operator-approved project/account/assignment before uploads or check-ins. Existing rules do not authorize integration writes to the shared database. No baseline migration, schema alteration or data cleanup on the shared project.
4. Upload a small owned PNG/JPEG and read it privately; reject wrong-owner access. Confirm actual file type/size. Use dedicated test paths; cleanup only specifically approved test objects, never a bucket purge.
5. At a real test assignment, use actual browser/device GPS and a fresh signed area QR. Missing/expired/wrong-area QR, invalid GPS/accuracy and duplicate check-in must fail; valid GPS AND QR must succeed. Never fabricate coordinates and call that live GPS verification. Checkout requires a real test photo.
6. Exercise test procurement handover/no_purchase and have the designated test LINE recipient confirm the actual message. This step sends a real message and needs that test recipient/environment; API acceptance alone is insufficient.
7. Record each observed status/result and whether it was simulated or live. Full acceptance is pending until all intended roles/providers/services have been exercised.

The current normal server runs attendance finalization every minute. Do not launch it against the shared project for a read-only diagnostic: a controlled test server must use isolated data or explicitly disable jobs/mutations. AGENTS.md allows normal configuration loading only with separate authorization; that never grants permission to inspect .env contents.
