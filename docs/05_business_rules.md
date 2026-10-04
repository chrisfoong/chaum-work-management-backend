# 05 — Business rules and workflow behavior

Source: use-case descriptions as stated by the team. Where a rule is marked **[unverified]** it comes from design notes rather than a document you can open here — confirm before relying on it. Do not change formulas, thresholds, paid-status meaning or negative-pay treatment.

## Quantities and procurement
- Site survey (1A): `to_buy_qty = GREATEST(required_qty - existing_qty, 0)`. `existing_qty` is entered per item and validated.
- Procurement is cumulative per item (2A, 7A): `actual_qty = actual_qty + new_qty`. There is **no receipt table**; the receipt photo is stored on EXPENSE_CLAIM.
- Remaining equipment per requisition is computed from required/actual quantities. A request moves toward `completed` when quantities are met; partial procurement keeps it open (`determineProcurementStatus`, `updateRequestStatusFromSavedQuantities`).
- 7A must recheck the latest data and prevent duplicate submission; it uses row locks on the requisition and its items and a commit step (`lockOriginalRequest`, `lockRequestItems`, `commitTransaction`).

## Expense claims
- `fund_transfer`: created by the Supervisor when sending money (2S). Requires a transfer reference; duplicate `transfer_ref_no` is rejected.
- `actual_expense`: created by the Assistant on purchase (2A, 7A) with receipt photo.
- Claims reference the requisition (`requisition_id`); one requisition can have many claims.
- **Only `actual_expense` counts as material cost.** Fund transfers are never added to material cost.
- 2S loops the requisition status back to `pending_procurement` after funding.

## Requisition lifecycle (observed statuses)
- 1S creates the initial requisition with type `tor_base`.
- 1A (site survey) works on `pending_survey` requisitions.
- 3S/3.1S: supervisor approves or rejects a pending requisition (`checkPendingRequisition`, `approveRequisition`, `rejectRequisition`, rejection reason required on reject). Approval/rejection triggers a LINE notification to the creator through a background job.
- 5W (worker) creates `additional` requisitions as `pending_survey`; 5A/6A list those; 6A forwards (to the supervisor path) or rejects; 7A performs the additional purchase; 8A records delivery evidence.
- Which step writes which status, and the funded additional path, are partly unresolved (D4, D5). Implement only what each SD states.

## Scheduling and leave
- 3A: schedule rows are created per worker per date. Check worker is available and no duplicate schedule for the same worker and date. Check the area exists and the assistant has permission (D12). On confirm, recheck the latest data, then notify assigned workers.
- 2W: worker picks a leave date from their scheduled dates; duplicate leave for that date is rejected; `is_advance_notice` is computed from current timestamp vs leave date; leave is created `pending`; assistant notified.
- 4A: assistant sees pending leave for their area, picks approve/reject, optionally a replacement from candidates (available workers for that date). Recheck selected replacement and latest data before confirm. Approve creates a replacement schedule; notify leaving worker, replacement worker, and supervisor if no replacement was found.

## Attendance
- 3W: substitute check-in — validate LINE user and schedule, today's shift must match the worker, duplicate check-in rejected, attendance UPSERT (D11).
- 4W: work report with description + photo; exactly one open attendance must exist; saving the evidence also sets check-out.
- 4S: end-of-day batch classifies each scheduled worker: on time, late (compare check-in with shift start), leave (count approved advance-notice leave), absent. UPSERT attendance rows. Must be safe to re-run.

## Deductions and payroll
- Late deduction **[tiers as stated by the team]**: under 1 hour = 300; 1 to 3 hours = 400; over 3 hours = between 1,000 and 2,000 (the exact amount within the range is not specified — ask). Leave with advance notice = no deduction. Absent handling per the use case.
- Each deduction carries `worker_id`, `attendance_id`, `penalty_amount`, `reason`.
- Payroll period closes every two weeks (as-is payments on the 5th and 20th). **[unverified]**
- `base_wage = work days × daily_wage` (`calculateBaseWage(workDays, dailyWage)`); `net_wage = base_wage − total_deduction` (`calculateNetWage`). Negative pay: undefined (D10).
- 5S must be idempotent (no duplicate payroll or deductions on re-run) **[team requirement; guard design is D8]**.
- 5S creates payroll with `is_paid = false`.
- 7W payslip: worker chooses a month; shows base_wage, total_deduction, net_wage, is_paid and the list of deductions with their work date (via attendance_id), penalty_amount, reason. A worker only sees their own data.

## Financial report (6S)
- `net_profit = total_revenue − (total_labor_cost + total_material_cost)`.
- Revenue: COMPANY_INVOICE (`billing_month`, `net_received` / `expected_amount`; D7).
- Labor: sum of `net_wage` from paid payroll in the period (D6, D9).
- Material: `actual_expense` claims for the project's requisitions in statuses approved / pending_supervisor / pending_procurement / completed.
- Output is shown and exported as a PDF (`exportFinancialReportPDF`, `generateReportPDF`).

## Contracts (1S, 9A)
- 1S: validate contract format, reject duplicate `contract_no`, validate or create the location (reject duplicate location name), create the contract, the TOR_LOCATION_ASSIGNMENT with `required_workers`, the initial `tor_base` requisition, and its items (find-or-create equipment by name). Notify via LINE (new contract).
- 9A: given an assignment, check access permission and evaluate the contract's status and end date to decide whether work can continue.

## Worker (LINE) access
- Every Worker endpoint verifies the linked LINE identity (LINE user id / token) server-side and resolves it to a USER. A worker may only read or write their own records. Check current official LINE/LIFF documentation for verification details; do not rely on memory.

## General validation behavior
Each use case specifies required fields, formats and DB checks; return validation errors that name the failing fields. Repeated "login" steps in SDs are preconditions, not new endpoints.
