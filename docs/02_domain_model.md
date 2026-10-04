# 02 — Domain model (production ER)

Authority for names and associations. Physical types, defaults and constraints are NOT specified here — they come from the real schema in this repo (verify, do not assume).

## 1. Entities and attributes
`{key}` = identifier shown underlined in the ER. Spelling is exact.

| Entity | Attributes |
|---|---|
| USER | `{user_id}`, first_name, last_name, phone_number, role, line_id, daily_wage, bank_name, bank_account_no, is_active, created_at, updated_at |
| WORKER | `is_available` only. Shares identity with USER (user_id = worker_id) |
| PAYROLL | `{payroll_id}`, period_start, period_end, base_wage, **total_deduction**, **net_wage** (derived; may be stored), is_paid, created_at, payroll_slip_no |
| CONTRACT_TOR | `{tor_id}`, project_name, contract_no, partner_agency, contract_value, start_date, end_date, contract_file_url, status, created_at, updated_at |
| TOR_LOCATION_ASSIGNMENT | `{assignment_id}`, required_workers |
| LOCATION | `{location_id}`, location_name, address, latitude, longitude |
| WORK_SCHEDULE | `{schedule_id}`, work_date, shift_start_time, shift_status |
| ATTENDANCE | `{attendance_id}`, work_date, check_in, check_out, status |
| LEAVE_REQUEST | `{request_id}`, leave_no, leave_date, reason, is_advance_notice, status, created_at (ER no longer draws user_id; the physical table may still have it) |
| WORK_EVIDENCE | `{evidence_id}`, photo_url, description, submitted_at |
| DEDUCTION_TRANSACTION | `{deduction_id}`, penalty_amount, reason, created_at (+ link to attendance; link to payroll) |
| EQUIPMENT_REQUISITION | `{requisition_id}`, requisition_no, requisition_type, status, reason, created_at, reviewed_by, reviewed_at (physical table also has `reject_reason` per the team's DB notes) |
| REQUISITION_ITEM | `{item_id}`, required_qty, existing_qty, to_buy_qty, actual_qty, actual_price, remark |
| EQUIPMENT | `{equipment_id}`, equipment_name, is_active |
| EXPENSE_CLAIM | `{expense_id}`, expense_no, expense_type, total_amount, receipt_photo_url, transfer_ref_no, created_at |
| COMPANY_INVOICE | `{invoice_id}`, invoice_no, billing_month, expected_amount, net_received, exat_deduction_amount (ER spelling; do not silently rename), deduction_reason, status |

Notes
- 16 entities total. TOR_LOCATION_ASSIGNMENT is the central project/location context.
- `USER.role` holds Supervisor/Assistant/Worker. Exact stored role strings: verify in schema.
- `CONTRACT_TOR.status` is visibly an attribute in the diagram even though its drawing connector is mis-bound.
- Money fields: base_wage, total_deduction, net_wage, penalty_amount, total_amount, actual_price, contract_value, expected_amount, net_received, exat_deduction_amount, daily_wage. Use exact decimals.

## 2. Relationships (name, endpoints, maximum cardinality)
Only maximum cardinality is known. Minimum participation (optional vs mandatory) is NOT specified — do not assume it.

| ID | Relationship | Max |
|---|---|---|
| R01 | USER RECEIVES PAYROLL | 1:M |
| R24 | **USER PAYS PAYROLL** (new in production ER; payer role) | 1:M |
| R02 | USER MANAGES CONTRACT_TOR | 1:M |
| R03 | USER REQUESTS EQUIPMENT_REQUISITION | 1:M |
| R04 | USER REVIEWS EQUIPMENT_REQUISITION | 1:M |
| R05 | USER CLAIMS EXPENSE_CLAIM | 1:M |
| R06 | CONTRACT_TOR BILLS COMPANY_INVOICE | 1:M |
| R07 | CONTRACT_TOR INCLUDES TOR_LOCATION_ASSIGNMENT | 1:M |
| R08 | LOCATION HOSTS TOR_LOCATION_ASSIGNMENT | 1:M |
| R09 | TOR_LOCATION_ASSIGNMENT REQUIRES EQUIPMENT_REQUISITION | 1:M |
| R10 | TOR_LOCATION_ASSIGNMENT GENERATES WORK_SCHEDULE | 1:M |
| R11 | TOR_LOCATION_ASSIGNMENT COLLECTS WORK_EVIDENCE | 1:M |
| R12 | EQUIPMENT_REQUISITION CONTAINS REQUISITION_ITEM | 1:M |
| R13 | EQUIPMENT IS_LISTED_IN REQUISITION_ITEM | 1:M |
| R14 | EQUIPMENT_REQUISITION GENERATES EXPENSE_CLAIM | 1:M |
| R15 | WORKER IS_ASSIGNED_TO WORK_SCHEDULE | 1:M |
| R16 | WORK_SCHEDULE RECORDS ATTENDANCE | 1:1 |
| R17 | WORKER LOGS ATTENDANCE | 1:M |
| R18 | WORKER REPLACES ATTENDANCE | 1:M |
| R19 | WORKER SUBMITS LEAVE_REQUEST | 1:M |
| R20 | WORKER CAPTURES WORK_EVIDENCE | 1:M |
| R21 | WORKER INCURS DEDUCTION_TRANSACTION | 1:M |
| R22 | ATTENDANCE TRIGGERS DEDUCTION_TRANSACTION | 1:M |
| R23 | PAYROLL APPLIES DEDUCTION_TRANSACTION | 1:M |
| U1 | USER — WORKER, marked `(p, e)` in the ER | specialization meaning: shared identity (decision adopted: user_id = worker_id) |

Meaning in practice
- A contract has many assignments; a location hosts many assignments. The assignment is a contract+location pairing with a required worker count. A uniqueness rule on the pairing is not stated.
- A requisition belongs to one assignment, has many items, and can generate many expense claims (so claims reference the requisition; the requisition has NO `expense_id`).
- Equipment master data (name, active) is on EQUIPMENT; the requisition item holds request-specific quantities and price.
- One WORK_SCHEDULE row = one worker on one date at one assignment. Scheduling several workers creates several rows.
- Schedule–attendance is shown 1:1 (maximum). Whether attendance exists before check-in is not stated.
- WORKER has both LOGS and REPLACES to ATTENDANCE: the original/replacement worker roles are not fully defined in the ER. Do not merge them or invent FK names.
- A deduction relates to a worker (INCURS), optionally an attendance (TRIGGERS) and a payroll (APPLIES).
- Payroll relates to USER via RECEIVES (recipient) and PAYS (payer). Payer semantics are an open decision (D1).
- Work evidence relates to an assignment (R11) and a worker (R20) in the ER, but some SDs also store `schedule_id` (open decision D2).

## 3. Observed value sets (from the references; physical sets are unverified)
| Field | Values seen |
|---|---|
| EQUIPMENT_REQUISITION.requisition_type | `tor_base`, `additional` |
| EQUIPMENT_REQUISITION.status | `pending_survey`, `pending_supervisor`, `pending_procurement`, `approved`, `completed`, (`rejected` implied by reject flows); `pending` appeared in an old 5W and is replaced by `pending_survey` |
| EXPENSE_CLAIM.expense_type | `fund_transfer` (supervisor sends money), `actual_expense` (assistant's receipts) |
| LEAVE_REQUEST.status | `pending`, `approved`, `rejected` |
| ATTENDANCE.status | on time, late, leave, absent (exact stored strings unverified) |
| PAYROLL.is_paid | boolean; payroll is created with `false` |
| CONTRACT_TOR.status / COMPANY_INVOICE.status | value sets not specified |

Earlier business-flow documents also mention record states such as CheckedIn → CheckedOut → Calculated → Confirmed (locked) and payroll Blocked/Calculated/Approved/Paid/Reconciled. These are NOT confirmed as physical columns and conflict in shape with the single `is_paid` boolean. Do not implement them without a decision.

## 4. Things the ER does NOT tell you
Data types and precision, ID type (UUID vs int), nullability, defaults, unique constraints, FK names, delete rules, status value sets, minimum multiplicities, the exact USER/WORKER table layout. Read the real schema. If it differs from this file, report before changing anything.

## 5. Overrides of the historical class diagram
The class-diagram package (included under `docs/references/` as HISTORICAL) was built from the previous ER. Where it disagrees with this file, THIS FILE wins:

| Historical | Current |
|---|---|
| PAYROLL.deduction | PAYROLL.total_deduction |
| PAYROLL.net_pay | PAYROLL.net_wage |
| LEAVE_REQUEST.user_id attribute | removed from ER (SUBMITS remains; physical key per schema) |
| USER RECEIVES PAYROLL only | RECEIVES and PAYS |
| 7W: total_wage, total_deduction, net_wage, status | base_wage, total_deduction, net_wage, is_paid |
| 7W: DEDUCTION_TRANSACTION.amount | penalty_amount |
| 2A: equipment_name from REQUISITION_ITEM | via EQUIPMENT join on equipment_id |
| 6S: er.expense_id join | ec.requisition_id = er.requisition_id |
| 5W: status pending | pending_survey |
| createDeduction(workerId, penaltyAmount, reason) | createDeduction(workerId, attendanceId, penaltyAmount, penaltyReason) |
| createPayroll(…, totalDeduction, netPay) / calculateNetPay | writes total_deduction, net_wage; calculateNetWage / netWage |
