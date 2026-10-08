# Frontend handoff — Backend / Supabase

Business scope: final user guide, 2026-10-08. This repository implements Backend only. Existing Supabase exports remain the physical schema; no migrations or production fixtures are required.

All requests use the appropriate verified LINE ID token as Bearer authorization. `/api/web` accepts Supervisor/Assistant channels/roles; `/api/liff` accepts Worker. Roles and identity come from the database, never request-supplied user IDs.

| Use case | API / sequence | Permission |
|---|---|---|
| 1S TOR wizard | `POST /files` (raw PNG), `/contracts/info`, `/contracts/scope`, `/contracts/confirm`; submit `contract.contract_file_path` at confirm. `/locations?q=` picker; new locations in scope. | Supervisor / web |
| 3A schedule | `GET /assignments`, `/workers/available?work_date=YYYY-MM-DD`, `POST /schedules` | Assistant / web |
| 4A review / replacement | `GET /leave-requests?status=pending`, `/leave-requests/:id`, `/leave-requests/:id/candidates`, `POST /leave-requests/:id/review` with optional `replacement_worker_id` | Assistant / web |
| 1A base survey | `GET /requisitions?status=pending_survey&requisition_type=tor_base`, detail, `POST /requisitions/:id/survey` with every item_id/existing_qty | Assistant / web |
| 2A / 7A purchase | requisition detail, `POST /files`, `POST /requisitions/:id/purchase` with each item's current `expected_actual_qty`, round quantity and decimal unit price | Assistant / web |
| 2S funding | `POST /files`, `POST /requisitions/:id/fund-transfers` amount/reference/receipt_path | Supervisor / web |
| 5A / 6A inspect / decide | `GET /requisitions/:id/inspection`, `POST /requisitions/:id/decision` purchase/no_purchase + reason | Assistant / web |
| 3S procurement approval | `GET /requisitions?status=pending_approval`, `POST /requisitions/:id/review` approve/reject; reject reason required <=1000 | Supervisor / web |
| 8A handover | completed additional requisitions, `GET /requisitions/:id/delivery-schedules`, files, `POST /requisitions/:id/delivery`, `/deliveries` evidence | Assistant / web |
| 4S / 5S payroll | `POST /attendance/finalize`, `/payroll/batch`, `GET /payroll`, `/payroll/:id/mark-paid` after actual transfer | Supervisor / web |
| 6S profit | `GET /reports/profit` or `/reports/profit.pdf` with tor_id/period_start/period_end; month also supported | Supervisor / web |
| 9A work dashboard | `GET /dashboard`, `/contracts`, `/contracts/:id/continuation` for active contract | Assistant / web |
| 1W schedule | `GET /schedules`; own upcoming assigned shifts with location/project/start/end | Worker / liff |
| 2W leave | `POST /leave-requests` leave_date/reason <=500; date must have an assigned schedule | Worker / liff |
| 3W check-in | Obtain assignment QR, GPS coordinates/accuracy, then `POST /attendance/check-in`; server time only | Worker / liff |
| 4W checkout | `POST /files` photo, then `/attendance/check-out` with photo_path and description <=1000 | Worker / liff |
| 5W shortages | `GET /equipment`, `POST /requisitions` positive requested quantities while checked in | Worker / liff |
| 6W equipment result | LINE Chat notification; no extra browser page/inbox required | Worker |
| 7W payslip | `GET /payroll?period_month=YYYY-MM`; own paid slips with base/deductions/net | Worker / liff |

`/files` above uses the current platform prefix. QR issuance is staff `POST /api/web/assignments/:id/qr`, tied to the assignment; the Worker scan returns its signed qr_token. Check-in requires both QR AND GPS <=200m, accuracy <=50m. There is no accept/reject-tomorrow endpoint. Shift end is start +8 hours, not Figma mock 17:00; timestamps can cross midnight.

Object paths stay private. Files <=5MiB; TOR PNG only; purchase/checkout/handover JPEG/PNG; Supervisor funding also PDF. Upload and business confirmation are separate requests; business failure may leave an unreferenced object. File read proxy uses `GET /files?path=` and enforces role/ownership/reference checks.

Lists return `data`, `limit`, `offset`. Use SQL filters, not client filtering on the first page. Date range/assignment/TOR/status supported on schedules/attendance/leave/requisitions; requisition_type only on requisitions. `/contracts` is paginated; dashboard's latest 100 contracts is a compact view, not the full registry. Empty lists are arrays. Money is decimal strings; dates are Gregorian YYYY-MM-DD, month YYYY-MM; display Thai years/currency in frontend. UUID worker_id differs from user_id and LINE Channel ID.

Candidate lists are advisory; transaction validates active/availability/leave/one worker-day and payroll period again. Rejected leave cannot have replacements. Replacement approval and schedule creation commit together or roll back together. Approving without a replacement is supported and alerts Supervisor; no invented replacement-confirmation state.

Delivery is full handover of purchased items to original requester in matching area, only schedules <= server's Bangkok today. Photos commit atomically and retries are idempotent. Evidence uses existing WORK_EVIDENCE description convention; no new delivered enum/FK/table or item-round history is claimed. Notification reports actual request/item quantities and is best effort, not guaranteed delivery.

No paid_at field is available: do not substitute payroll created_at. `mark-paid` records the flag and never transfers money. Financial confirmation is a snapshot with persisted=false. PDF uses English financial fields/TOR UUID; Thai project names are available in JSON/CSV. Supabase/LINE/private Storage end-to-end delivery needs deployed configuration, distinct from isolated test success.
