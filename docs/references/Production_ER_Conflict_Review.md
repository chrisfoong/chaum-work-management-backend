> **REFERENCE COPY — audit of the production ER vs previous ER and the team SDs (before the corrections). Findings marked as corrected are listed in docs/04 §1.3; the rest are open decisions D1–D14.**

# Production ER compared with the previous ER and team-edited sequence diagrams

Sources compared:
- New: `SA1-ER PRODUCTION.drawio.svg`, embedded page `ER PRODUCTION`.
- Previous: `SA1-ER latest(ปรับทะแยง).drawio.svg`.
- Sequence diagrams: all 23 extracted team-edited SD pages from `SA1-SD-1A.drawio.svg`, including SD-3S and SD-3.1S. These are the team's versions, rather than earlier generated SDs.

The embedded models were inspected, the new SVG was rendered, and every team SD query label was reviewed. Some SQL labels contain `...`; omitted columns or join conditions cannot be verified from those labels. Original use-case SQL is identified separately below. No diagrams were modified.

## 1. Changes from the old ER

The same 16 entities remain. No entity was added or removed.

| Change | Consequence |
|---|---|
| PAYROLL.deduction → PAYROLL.total_deduction | Old ER-based class diagram and original 5S SQL use the previous name. New name agrees with 7W. |
| PAYROLL.net_pay → PAYROLL.net_wage | Old ER-based class diagram and 6S query use the previous name. New name agrees with 7W. It remains a derived attribute. |
| LEAVE_REQUEST.user_id oval removed | Worker-to-leave SUBMITS relationship remains. Need to establish whether the physical leave table retains user_id, uses worker_id, or uses an inherited user identifier. Removal alone does not prove SQL is invalid. |
| USER — PAYS — PAYROLL added, shown 1:M | RECEIVES remains, also USER-to-PAYROLL. They appear to distinguish payer and recipient, but role constraints and physical FK names are unspecified. Do not merge them into one relationship. |
| HOSTS connector reattached to TOR_LOCATION_ASSIGNMENT | Fixes the old embedded connector that ended on an M label. Intended LOCATION-to-assignment relationship is unchanged. |

Other changes reposition/redraw USER.created_at and PAYROLL.is_paid and adjust routing/cardinality label placement. Their fields remain. PAYROLL.base_wage, is_paid, created_at and payroll_slip_no remain. CONTRACT_TOR.status is visibly present; its embedded connector still attaches to an M label, so automatic extraction must not interpret it as removed.

The earlier class diagram and Claude handoff are now outdated specifically for payroll field names, the explicit leave user_id attribute, and the added PAYS relationship. They were not updated in this review.

## 2. Definite field/query conflicts, highest priority

| Priority | SD | What the SD uses | New ER | Practical consequence |
|---|---|---|---|---|
| 1 | 6S, labor-cost SELECT | SUM(net_pay) FROM payroll | net_wage | Query references the previous field name; fails against a schema using the new ER name without compatibility support. |
| 2 | 7W, payslip SELECT | p.total_wage | base_wage | No total_wage attribute. Needs the canonical field or an explicitly declared projection/alias. |
| 3 | 7W, payslip SELECT | p.status | is_paid | No payroll status attribute. A payment boolean is not automatically equivalent to an unrestricted status field. |
| 4 | 7W, deduction SELECT | dt.amount | penalty_amount | Query references a different deduction amount field. |
| 5 | 2A, equipment detail SELECT | equipment_name FROM REQUISITION_ITEM | equipment_name belongs to EQUIPMENT | Shown query has no equipment join; field is on a different entity. |
| 6 | 6S, material-cost JOIN | ec.expense_id = er.expense_id | One requisition GENERATES many claims | er.expense_id is neither an ER attribute nor the expected FK direction for this relationship. The model supports claims linked to requisitions. A single expense pointer would need an extra stated purpose. |

Original 5S use-case SQL explicitly inserts `base_wage, deduction, net_pay, is_paid`. Its `deduction` and `net_pay` column names conflict with the new ER. The team-edited SD-5S image abbreviates the insert as `INSERT INTO payroll (..., is_paid)`, so the image alone does not expose those two insert-column names. Its `createPayroll(..., totalDeduction, netPay)` parameters are operation parameter names, not proof of database column names.

What improved: 7W's `p.total_deduction` and `p.net_wage` now agree with the new ER. Its total_wage/status and deduction amount mismatches still remain. 5S's `AS total_deduction` aggregate alias is compatible; an SQL alias is not an extra database attribute.

## 3. Mapping differences that may be valid

These need a physical schema or a stated mapping; they should not all be reported as missing-column errors.

| Area | SD evidence | Why unresolved |
|---|---|---|
| Leave identity | 2W selects/inserts LEAVE_REQUEST.user_id; 4S filters leave by user_id; 4A uses leave joins with omitted parts | New ER removes the explicit user_id oval but retains WORKER SUBMITS LEAVE_REQUEST. An FK can implement the relationship without an oval. Shared USER/WORKER identity could make user_id valid, but it is not specified. |
| Worker identity | 3A and many W diagrams use w.worker_id and w.user_id | WORKER still explicitly shows only is_available and the USER connection labelled (p,e). Identifier inheritance or separate worker_id mapping is not clearly defined. |
| Payroll ownership and payment | 7W joins p.worker_id to WORKER; original 5S inserts worker_id | New ER links payroll to USER through both RECEIVES and PAYS. Worker-as-user mapping could reconcile the recipient, but the payer role is not shown in the payroll write flow. |
| Work evidence scope | 8A and 4W insert schedule_id and worker_id | ER associates evidence with assignment and worker, not directly with WORK_SCHEDULE. The assignment may be inferred through the schedule, but that storage/mapping decision is not shown. |
| Deduction links | 5S inserts worker_id, penalty_amount, reason; 7W joins dt.attendance_id; ER also links deductions to PAYROLL | Shown 5S insert does not record attendance/payroll links needed by downstream views and ER relationships. Hidden defaults or later updates are not stated; this is a flow gap rather than a proven nonexistent column. |
| Derived net wage | ER uses a dashed net_wage oval; original 5S stores a calculated result | Storing a derived calculation can be valid. The derived notation alone does not prohibit persistence. |
| Relationship foreign keys | SDs use tor_id, assignment_id, requisition_id, equipment_id and related IDs | Many implement relationships already drawn in the ER. Missing FK ovals alone are not conflicts. |
| New PAYS association | USER pays PAYROLL in addition to receiving it | Extra relationship is not automatically a business conflict. Need distinct payer/recipient semantics if both are implemented. |

## 4. Existing workflow gaps still present

These are SD-to-SD or processing-scope gaps; changing ER field names does not fix them.

- **5W → 5A/6A:** 5W inserts additional requests with status `pending`; 5A and 6A select `pending_survey`. The transition is not shown, so the new request can disappear from the review queue.
- **3S → 7A approval evidence:** 7A refers to reviewed_by/reviewed_at and prior approval for the funded incomplete path. The shown 3S approve query updates status only. 6A writes review fields on rejection, not on the shown escalation to pending_supervisor. Required approval evidence is not established in the shown approval flow.
- **2S → 7A funding status:** 2S changes a funded request to `pending_procurement`; 7A's shown list admits `approved` or a specified `pending_supervisor` funding path. The funded additional request path needs reconciliation. 2A does admit pending_procurement, but that does not explain the intended 7A path.
- **6S material status:** 6S includes only approved requisitions, while 2A/7A set fulfilled requisitions to completed. Completed purchases can be omitted.
- **6S material expense types:** total_amount is summed without an expense_type restriction; both fund_transfer and actual_expense can exist. If they represent the same funds at different stages, counting both can overstate cost. Accounting treatment is unspecified.
- **6S project scope:** labor-cost query sums paid payroll for the date interval without a project condition. A project report could include other projects' wages. ER does not establish a project payroll allocation rule.
- **6S date input:** report describes a date range, while revenue is selected by billing_month. The conversion from range to billing months is unspecified.
- **5S deduction period:** deductions are selected by created_at, while penalties originate from attendance work dates. Generating deductions after the work period can omit them from that period's total. The shown flow also lacks a reprocessing duplicate guard.
- **5S → paid reports:** 5S inserts is_paid=false; 6S filters is_paid=true and 7W assumes payment occurred. A separate payment process may exist, but it is outside the shown 23 SDs.

## 5. Embedded ER drawing issues

- HOSTS endpoint issue from the previous ER is fixed.
- TRIGGERS still has an endpoint attached to the M label rather than DEDUCTION_TRANSACTION; the rendered line indicates the intended association.
- REQUESTS still uses an unattached fixed endpoint rather than binding directly to USER; visually intended USER-to-request connection remains.
- CONTRACT_TOR.status still attaches to a cardinality label rather than directly to the contract entity.

These are editable-diagram connection defects, not evidence that the intended business relationships were deleted.

## 6. Suggested reconciliation order

1. Treat production ER payroll names as canonical, then reconcile 5S/6S/7W columns together.
2. Resolve 2A equipment lookup and 6S requisition-to-expense join.
3. Define USER/WORKER identity and payer/recipient/leave FK mappings.
4. Define evidence and deduction relationship storage.
5. Reconcile additional-request statuses and recorded approval evidence across 5W, 5A, 6A, 3S, 2S and 7A.
6. Resolve financial period/project/expense-type rules.
7. Update the class diagram and handoff after these decisions; repair remaining ER connector bindings.

This is an audit only. No source ER, SD, class diagram or handoff was changed.
