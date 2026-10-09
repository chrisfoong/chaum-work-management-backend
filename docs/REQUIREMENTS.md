# Backend requirements — 2026-10-08

Business source: [SA-Group5 Use Case Descriptions](https://docs.google.com/document/d/1Ovggbw2PwWJJdyzzkrCMP1CruIoPelmZ6f9whHYZryQ/edit?tab=t.0), explicitly adopted by the user. Physical mapping remains Supabase exports in schema/. Document SQL is illustrative; never execute it as a baseline. Backend only.

- LINE roles/audiences, unknown/inactive rejection and ownership remain unchanged. Assistant manages all assignments without Supervisor financial privileges.

- 3A: active/in-date contract, all initial TOR equipment completed, active available workers, minimum required_workers (not a maximum), one worker/day. Eight-hour shifts support midnight crossing.

- 2W/4A: future Bangkok date is advance leave, same day is emergency; past requests rejected. Assistant approves/rejects and optionally creates replacement atomically. Keep original schedule; approval without replacement alerts Supervisor. Existing replacement endpoint supports later assignment.

- 4S: only approved advance leave avoids absent status/penalty. Approved emergency leave without check-in remains absent. Retain shift-end +2h cutoff where the description has no precise cutoff.

- 3W permits assigned workers to confirm/check in on the scheduled calendar day before start; 4S classifies early check-in as on_time. Future-day and ended shifts are rejected; overnight continuation remains supported.

- Attendance retains both signed QR (60s) and GPS <=200m through LOCATION/TOR assignment, accuracy <=50m, server timestamps/Bangkok dates.

- 5W requires current active check-in/no checkout. 5A compares original additional request against base and pending additional requests across the TOR. Historical quantity is not proof of equipment condition.

- 6A explicitly decides purchase/no_purchase on original request. Purchase goes to Supervisor approval; no_purchase maps to rejected and sends explanation to requester without persisting it, per description.

- 1A surveys base requests only. 2A/7A support partial/repeated purchases with cumulative actual_qty. Clients send expected_actual_qty for stale/replay rejection. Positive-cost rounds require receipts and actual_expense; zero-cost rounds do not create expense rows because total_amount CHECK is >0.

- 2A stores latest actual unit price. 7A uses weighted cumulative unit price rounded to exported NUMERIC scale (two decimals). Actual expenses preserve exact per-round cost. Base and additional purchase receipts JPEG/PNG; funding receipts may include PDF.

- Enum mapping: approval pending_supervisor -> pending_approval; partial funding pending_supervisor -> pending_fund; approved -> pending_procurement with reviewed_by. No enum/schema alterations.

- 8A delivers fully purchased additional requests to requesting worker's matching assignment/schedule. One WORK_EVIDENCE row/photo with request, recipient, equipment summary and description. Atomic photos; matching retries succeed, mismatched retries conflict. Notify Worker after delivery, never after purchase. Explicit description prefix identifies records without inventing a FK/status.

- 5S batch all scheduled workers in half-month period, atomic slips/deductions, retry reuse. Base = actual on_time/late days * daily_wage, penalties 300/400/1500 THB, applied deduction <=base, net >=0, exact satang. 7W own paid slips with period_month filter.

- 6S profit filters TOR/date range; paid invoice net_received by billing month, paid payroll allocated by actual TOR workdays, actual_expense material only. Include payroll periods fully contained in selected dates; invoices cover selected calendar months. JSON/CSV/PDF. 9A active-contract continuation details, no automatic new schedules.

Limits: no persistent closing ledger, delivery relation/status, item-round history, replacement relation, wage snapshots or durable outbox. Closing confirmation returns a timestamped snapshot persisted=false, not contract termination. Delivery uses a description convention. Notification acceptance does not prove delivery. PDF uses TOR UUIDs and English financial fields; Thai project names remain in JSON/CSV. Direct DB writers can bypass application locks.

Never inspect .env or modify Supabase schema. Integration writes use an explicitly isolated loopback test database only.

## Frontend guide adopted on 2026-10-08

No advance replacement accept/reject flow. Worker confirms today's assigned shift through GPS AND QR check-in, server timestamp, exactly eight hours (08:00-16:00 when scheduled at 08:00), overnight supported. Photos are required at checkout, not check-in. Equipment notifications are LINE Chat only; no notification inbox. Payslip shows base/deductions/net without inventing paid_at. Assistant dashboard reports actual contract status and derived workflow_status/can_operate. Figma mock 09-hour shifts and alternative GPS/QR do not override these policies.

TOR confirmation requires a private uploaded PNG <=5MiB, verified server-side and stored in existing contract_file_url; no public URL or schema change. Wizard /contracts/info is a format preview; verification is enforced at /contracts/confirm. Base and additional purchases both use JPEG/PNG receipts when total >0; funding allows JPEG/PNG/PDF. Leave reasons <=500, Supervisor procurement rejection <=1000, checkout description <=1000.

Replacement candidates exclude original requester, inactive/unavailable, existing same-day schedules (including cancelled because DB uniqueness persists), approved leave and already calculated payroll periods. Eligibility is rechecked when scheduling; duplicate mappings conflict before candidate pagination. API reads filter in SQL before pagination and preserve Worker ownership.

Delivery rejects future schedules and incomplete procurement even if status was incorrectly marked completed. Guide 8A supplies schedule/description/photos, no new delivered-quantity field: all purchased quantities are handed over together, validated in the locked request transaction. Delivery LINE message includes request number/item quantities; messages split conservatively at LINE character bounds. Notifications support bounded keyed retries and a guarded manual resend API, without a durable queue.

## Corrections from detailed 22-use-case audit

2A permits all-zero acquisition with pending_fund and no actual expense; 7A requires a positive quantity. Initial Supervisor-approved 7A purchase can proceed directly; funding is needed to resume a partial request after pending_fund. Purchase approval metadata includes both reviewed_by and reviewed_at. Worker equipment reasons <=500. no_purchase records reviewer/time without replacing original reason and can decline inactive equipment. Worker schedule status/address are visible. Finance details reconcile to totals. Decision/delivery returns business_saved separately from notification acceptance; Assistant can retry notifications without business writes. 5S retains Worker payslip notices but never automatically sends the Assistant operational summary. 9A returns current assignment operations for selected dates; Supervisor sends a short dashboard link manually after reviewing 5S/6S. No snapshots, send/read history or acknowledgement are required in MVP. See USECASE_AUDIT.md.

## Confirmed 9A MVP — 2026-10-08

Keep exactly the existing 16 tables. No audit_events, outbox_events, summary snapshots, send timestamps, read receipts or acknowledgement. Reports query current source data on each request. Assistant operational fields exclude payroll, deductions, prices, costs, bank details and profit. Prior all-area Assistant authorization remains unchanged.

Supervisor explicitly reviews 5S/6S and manually triggers a short LINE project/period message with an HTTPS Assistant dashboard link. No Assistant summary is automatically sent from PayrollBatch or CloseSummary. Manual command requires a closed period, source data, payroll coverage for scheduled workers and supervisor_reviewed=true; 6S is recomputed, but prior human review is an attestation because no confirmation ledger exists.

No source records in the selected period returns has_data=false and “ยังไม่มีข้อมูลสรุปสำหรับรอบนี้”. An inactive/ended/expired contract returns can_continue=false and “กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ”. Frontend disables new schedules; Backend still rejects them. No database schema change is needed.
