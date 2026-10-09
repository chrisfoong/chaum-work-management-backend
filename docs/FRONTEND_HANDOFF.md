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
| 9A operations dashboard | `GET /assignments/:id/operations-summary?period_start=YYYY-MM-DD&period_end=YYYY-MM-DD`; existing dashboard/assignments for area selection | Assistant / web |
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

Detailed audit update: see USECASE_AUDIT.md and API.md. Cards include address and stored future statuses; shortages require reason <=500. 7A first approved purchase does not require prior funding; partial rounds do. 2A all-zero acquisition can escalate. Profit JSON includes labor_details/material_details. For selected area use Assistant /assignments/:id/continuation. Read decision/delivery business_saved and notification.status separately; use Assistant /notifications/:id/retry to resend failed notices without resubmitting a purchase or delivery. Do not display historical summary notices as implemented; schema has no notice/read ledger.

Assistant POST /api/web/requisitions/:id/purchase/preview uses PurchaseInput and returns items(new_actual_qty/remaining_qty/new_actual_price/round_cost), round_total, next_status and persisted=false. It runs confirmation validation without quantity/expense/status writes. Use for 7A review; POST /purchase rechecks current data when confirming.

## Confirmed 9A integration

Frontend remains in its separate repository. Assistant chooses assignment and dates, queries operations-summary, displays only attendance/procurement, and refreshes current data. Use has_data/summary_notice for empty data; use can_continue/continuation_notice to disable new scheduling on ended/inactive/expired contracts. Never display unavailable notification/read history or an acknowledgement button.

Supervisor button “แจ้งผลสรุปให้ผู้ดูแลงาน” calls the manual notify endpoint only after processing/reviewing 5S/6S, with supervisor_reviewed=true. This assertion is not a stored approval. Show returned notification status separately from financial processing success. Configure ASSISTANT_DASHBOARD_URL privately to the real HTTPS Assistant page; .env.example contains only a placeholder. Dashboard handles tor_id/period_start/period_end and authenticates before fetching current operational data. No automatic summary from payroll/profit confirmation.
