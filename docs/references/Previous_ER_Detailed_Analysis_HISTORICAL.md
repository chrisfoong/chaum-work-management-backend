> **HISTORICAL — analysis of the PREVIOUS ER (not production). Use for relationship meanings and reasoning; names marked outdated in docs/02 §5.**

# Detailed analysis of the latest ER reference

Source: `C:/Users/PC01/Downloads/SA1-ER latest(ปรับทะแยง).drawio.svg`.

The SVG contains one embedded draw.io page, `ER latest(ปรับทะแยง)`. I inspected the rendered page and its editable XML. For cross-checking, I used the 23 team-edited sequence pages extracted from `SA1-SD-1A.drawio.svg`, as well as the use-case text supplied in this conversation where SQL is abbreviated in those pages.

No source diagram has been changed. This report distinguishes visible ER facts, connector problems, and inconsistencies with the sequence/use-case references. It does not assume the actual database schema is identical to either reference.

## 1. Overall structure

There are **16 entity rectangles, 23 named relationship diamonds, and a separate USER–WORKER connection labelled `(p, e)`**. The two relationships named GENERATES are separate diamonds with different endpoints. Small unnamed circles on lines are connector points, not additional entities or business attributes.

The main chains are:

- Project and location: CONTRACT_TOR → TOR_LOCATION_ASSIGNMENT ← LOCATION.
- Equipment: assignment → EQUIPMENT_REQUISITION → REQUISITION_ITEM → EQUIPMENT, with EXPENSE_CLAIM linked to the request.
- Work: assignment → WORK_SCHEDULE → ATTENDANCE, with WORKER linked to schedules, attendance, leave and evidence.
- Finance: attendance → DEDUCTION_TRANSACTION → PAYROLL; CONTRACT_TOR → COMPANY_INVOICE.

TOR_LOCATION_ASSIGNMENT is the central project/location context. USER is the shared account context. The layout largely supports the workflows already described, but several identity and association choices remain unsettled.

## 2. Entity and attribute inventory

Underlined identifiers below are shown as keys in the drawing. That visual notation does not prove that a physical database has the corresponding PRIMARY KEY constraint or a particular data type. Field spelling is preserved exactly.

| Entity | Shown key | Other shown attributes | Purpose |
|---|---|---|---|
| `USER` | `user_id` | `first_name`, `phone_number`, `last_name`, `created_at`, `updated_at`, `is_active`, `bank_account_no`, `daily_wage`, `bank_name`, `line_id`, `role` | Account, role, LINE identity, wage rate and banking details. |
| `WORKER` | None shown | `is_available` | Worker specialization/profile and availability. |
| `PAYROLL` | `payroll_id` | `period_start`, `period_end`, `base_wage`, `deduction`, `is_paid`, `created_at`, `net_pay` (derived), `payroll_slip_no` | A payment period, wage totals, deduction and payment state. |
| `CONTRACT_TOR` | `tor_id` | `start_date`, `project_name`, `end_date`, `contract_value`, `created_at`, `updated_at`, `partner_agency`, `contract_no`, `contract_file_url`, `status` | Contract and project requirements. |
| `ATTENDANCE` | `attendance_id` | `work_date`, `check_in`, `check_out`, `status` | Recorded attendance and check-in/check-out. |
| `WORK_SCHEDULE` | `schedule_id` | `work_date`, `shift_status`, `shift_start_time` | A dated shift assigned to a worker. |
| `LEAVE_REQUEST` | `request_id` | `user_id`, `leave_date`, `is_advance_notice`, `reason`, `created_at`, `status`, `leave_no` | Leave date, reason and approval state. |
| `TOR_LOCATION_ASSIGNMENT` | `assignment_id` | `required_workers` | A contract-to-location assignment and minimum staffing. |
| `LOCATION` | `location_id` | `address`, `latitude`, `location_name`, `longitude` | Work location and coordinates. |
| `WORK_EVIDENCE` | `evidence_id` | `photo_url`, `submitted_at`, `description` | Evidence photo, description and submission time. |
| `DEDUCTION_TRANSACTION` | `deduction_id` | `penalty_amount`, `reason`, `created_at` | An individual penalty/deduction. |
| `EQUIPMENT_REQUISITION` | `requisition_id` | `requisition_type`, `created_at`, `reason`, `requisition_no`, `status`, `reviewed_by`, `reviewed_at` | Equipment request, lifecycle and review metadata. |
| `REQUISITION_ITEM` | `item_id` | `required_qty`, `remark`, `actual_qty`, `actual_price`, `existing_qty`, `to_buy_qty` | A requested equipment line and quantities/pricing. |
| `EQUIPMENT` | `equipment_id` | `equipment_name`, `is_active` | Active equipment master records. |
| `EXPENSE_CLAIM` | `expense_id` | `receipt_photo_url`, `total_amount`, `transfer_ref_no`, `created_at`, `expense_no`, `expense_type` | Actual expenses or fund transfers and their evidence. |
| `COMPANY_INVOICE` | `invoice_id` | `billing_month`, `net_received`, `exat_deduction_amount`, `deduction_reason`, `expected_amount`, `status`, `invoice_no` | Project billing and amounts received. |

`CONTRACT_TOR.status` is visibly placed as a contract attribute, but its connector attaches to an M label instead of the contract rectangle. It is included above as a visually intended field, with that attachment defect recorded below.

### Important field observations

- Fifteen entities have a shown underlined identifier. WORKER shows only `is_available`, with no separate `worker_id` or underlined key.
- USER contains `role`; separate Supervisor and Assistant entities are not drawn. These roles must not automatically become subclasses in a class draft.
- PAYROLL has `base_wage`, `deduction`, `net_pay`, and `is_paid`. `net_pay` is the only named attribute drawn with a dashed oval. This identifies it as derived in the ER. A derived value may also be stored; the drawing alone does not resolve that implementation choice.
- `to_buy_qty` and the other requisition quantities have ordinary solid ovals. Although calculations in the use cases produce these values, the ER represents them as ordinary attributes.
- `COMPANY_INVOICE.exat_deduction_amount` is the exact spelling shown. It looks like a spelling issue, but the intended replacement cannot be selected from this diagram alone.
- Identifier types, decimal precision, column nullability, defaults, unique constraints and status value sets are not specified. UUID versus integer must therefore come from a schema or explicit use-case requirement, not a guess based on the name.

## 3. Relationships and cardinalities

The following table records the **visible intended endpoints and 1/M labels**, including links whose XML attachments need repair. It does not add minimum participation. A label M indicates a many side, but does not establish whether zero related records are allowed.

| Entity A | A side | Relationship | B side | Entity B |
|---|---:|---|---:|---|
| `USER` | 1 | RECEIVES | M | `PAYROLL` |
| `USER` | 1 | MANAGES | M | `CONTRACT_TOR` |
| `USER` | 1 | REQUESTS | M | `EQUIPMENT_REQUISITION` |
| `USER` | 1 | REVIEWS | M | `EQUIPMENT_REQUISITION` |
| `USER` | 1 | CLAIMS | M | `EXPENSE_CLAIM` |
| `CONTRACT_TOR` | 1 | BILLS | M | `COMPANY_INVOICE` |
| `CONTRACT_TOR` | 1 | INCLUDES | M | `TOR_LOCATION_ASSIGNMENT` |
| `LOCATION` | 1 | HOSTS | M | `TOR_LOCATION_ASSIGNMENT` |
| `TOR_LOCATION_ASSIGNMENT` | 1 | REQUIRES | M | `EQUIPMENT_REQUISITION` |
| `TOR_LOCATION_ASSIGNMENT` | 1 | GENERATES | M | `WORK_SCHEDULE` |
| `TOR_LOCATION_ASSIGNMENT` | 1 | COLLECTS | M | `WORK_EVIDENCE` |
| `EQUIPMENT_REQUISITION` | 1 | CONTAINS | M | `REQUISITION_ITEM` |
| `EQUIPMENT` | 1 | IS_LISTED_IN | M | `REQUISITION_ITEM` |
| `EQUIPMENT_REQUISITION` | 1 | GENERATES | M | `EXPENSE_CLAIM` |
| `WORKER` | 1 | IS_ASSIGNED_TO | M | `WORK_SCHEDULE` |
| `WORK_SCHEDULE` | 1 | RECORDS | 1 | `ATTENDANCE` |
| `WORKER` | 1 | LOGS | M | `ATTENDANCE` |
| `WORKER` | 1 | REPLACES | M | `ATTENDANCE` |
| `WORKER` | 1 | SUBMITS | M | `LEAVE_REQUEST` |
| `WORKER` | 1 | CAPTURES | M | `WORK_EVIDENCE` |
| `WORKER` | 1 | INCURS | M | `DEDUCTION_TRANSACTION` |
| `ATTENDANCE` | 1 | TRIGGERS | M | `DEDUCTION_TRANSACTION` |
| `PAYROLL` | 1 | APPLIES | M | `DEDUCTION_TRANSACTION` |

### What these relationships mean

- A contract can include multiple assignments; a location can host multiple assignments. TOR_LOCATION_ASSIGNMENT represents the specific contract/location pairing and its staffing requirement. The ER does not explicitly show a uniqueness constraint on that pairing.
- A requisition belongs to one assignment, has multiple items, and can generate multiple expense claims. The separate claim entity therefore accommodates more than one payment/expense event for the same request.
- One equipment master record can occur on many request items. Equipment identity and descriptive data belong to EQUIPMENT; the item holds request-specific quantities and price.
- Each WORK_SCHEDULE is shown with one worker and one assignment. Assigning several workers in 3A can therefore mean creating several schedule records, as that use case explicitly does.
- WORK_SCHEDULE–ATTENDANCE is visibly 1:1. This is a maximum-cardinality statement in the drawing; it does not settle whether attendance is optional before check-in. The 3W compound UPSERT key `(schedule_id, worker_id)` also does not by itself prove that schedule_id is unique across attendance rows. The physical constraints would need to establish the intended maximum.
- WORKER has both LOGS and REPLACES relationships to ATTENDANCE. The second name suggests a replacement role, but the diagram does not define how the original worker, replacement worker, leave and replacement schedule are identified together. These two associations should not be merged or given invented foreign-key names.
- A payroll can have multiple deductions, and an attendance record can trigger multiple deductions. This permits more than one deduction per attendance in the ER, although 5S describes one calculated penalty for each late/absent item.
- PAYROLL is related directly to USER through RECEIVES. The payroll use cases instead refer to `worker_id`. These can represent compatible business ownership, but the identity mapping is not explicitly resolved here.
- WORK_EVIDENCE belongs to an assignment and is captured by a worker. There is no explicit WORK_SCHEDULE–WORK_EVIDENCE relationship on this ER page.

## 4. Four drawing attachment problems

These are confirmed from the embedded XML. They affect moving/editing the diagram and extracting relationships; they do not, by themselves, prove a flaw in the underlying database.

| Visible intention | Embedded attachment | Consequence |
|---|---|---|
| LOCATION — HOSTS — TOR_LOCATION_ASSIGNMENT | The assignment-side connector targets the M text cell `8tqF58g8p9KXpsFFBl3G-198`. | It stops at the cardinality label rather than being attached to the assignment entity. |
| ATTENDANCE — TRIGGERS — DEDUCTION_TRANSACTION | The deduction-side path terminates at M text cell `8tqF58g8p9KXpsFFBl3G-183` through a connector circle. | The visible line does not form a complete entity-to-entity attachment in the editable model. |
| EQUIPMENT_REQUISITION — REQUESTS — USER | Edge `8tqF58g8p9KXpsFFBl3G-203` has a fixed targetPoint and no target cell. | Its end is visually near USER but is not bound to USER. |
| CONTRACT_TOR — status | Edge `rlVyrjYtN52kc_-kk-_f-3` starts at M text cell `XjxD-yyzfwrAlXFz3RyX-4`. | Moving the contract independently can leave its status attribute incorrectly attached. |

I used the visible labels and layout to describe the intended relations above; an automatic XML-only conversion would miss or misread these four connections.

## 5. Confirmed inconsistencies with the latest sequences

### A. Worker identity is unresolved

**ER:** USER has `user_id`; WORKER has only `is_available`, and a direct USER–WORKER line labelled `(p, e)`.

**Sequences:** 3A, 1W, 3W, 4W and others use `WORKER.worker_id` and `WORKER.user_id`.

If WORKER is a subtype sharing USER's identifier, that is one design. If it has its own worker_id plus a reference to user_id, that is another design. The ER does not explicitly settle the difference. The `(p, e)` marker suggests a specialization constraint, but the file contains no legend defining it; I would not treat its exact meaning as confirmed. No worker_id should be invented during a faithful ER-only class conversion.

### B. Payroll names disagree in 7W

| ER / 5S terminology | 7W SELECT terminology |
|---|---|
| `base_wage` | `total_wage` |
| `deduction` | `total_deduction` |
| `net_pay` | `net_wage` |
| `is_paid` | `status` |

7W selects the latter names as actual PAYROLL columns, without aliases. This is not merely different display wording. The supplied 5S INSERT and the ER use the former names. The team-edited 6S also sums `net_pay` and filters `is_paid = true`. One canonical set must be selected before these references can be treated as a single consistent schema. `is_paid` and `status` should not be assumed equivalent without a defined mapping.

### C. Deduction fields and linkage disagree

**ER:** DEDUCTION_TRANSACTION has `penalty_amount`, with an ATTENDANCE association and a PAYROLL association.

**5S sequence:** inserts `(worker_id, penalty_amount, reason)` and totals by `created_at BETWEEN period_start AND period_end`.

**7W sequence:** selects `dt.amount` and joins using `dt.attendance_id`.

The amount column name differs. Also, the 5S INSERT does not explicitly supply an attendance link, although 7W needs that link to show `work_date`. The ER relationship expresses the business association, but the shown write query does not demonstrate how it is recorded.

Separately, filtering deductions by creation time does not necessarily select deductions for the work period: a deduction created in October for September work would not be included in a September creation-time range. Reprocessing 5S can also add duplicate deduction rows because its shown INSERT has no duplicate guard. The physical schema could have additional protection; none is established by this reference.

### D. Work evidence is linked at different levels

**ER:** TOR_LOCATION_ASSIGNMENT — COLLECTS — WORK_EVIDENCE, and WORKER — CAPTURES — WORK_EVIDENCE.

**4W and 8A sequences:** INSERT WORK_EVIDENCE includes `schedule_id` and `worker_id`.

The schedule relationship is not drawn in the ER. A schedule already provides an assignment context, but that does not make the ER's direct assignment association automatically equivalent to an explicit schedule association. A faithful class draft must record this mismatch rather than silently adding a schedule relationship.

### E. The additional-request status chain is discontinuous

**5W:** creates an additional request with `status = 'pending'`.

**5A and 6A:** select additional requests with `status = 'pending_survey'`.

No intermediate transition from pending to pending_survey is shown in these references. As written, a newly created 5W request is excluded from the 5A list. The ER's generic `status` field cannot resolve this workflow inconsistency by itself.

### F. 2A reads a field from a different entity

The team-edited 2A query selects `equipment_name` directly FROM REQUISITION_ITEM. In this ER, that field belongs to EQUIPMENT; REQUISITION_ITEM has no such attribute. Other workflows join EQUIPMENT to obtain the name. A hidden view or denormalized column could change the physical situation, but neither is shown in this ER.

### G. 6S expense linkage and filtering need reconciliation

The ER shows one requisition generating many expense claims. The supplied write flows link claims using `EXPENSE_CLAIM.requisition_id`, including the explicit 7A INSERT.

The team-edited 6S cost query instead joins on `ec.expense_id = er.expense_id`. A separate `EQUIPMENT_REQUISITION.expense_id` pointer is not shown in this ER. Such a pointer would also select a particular claim rather than directly express the shown one-to-many request/claim association.

The same 6S query filters requests to `approved`, while 2A/7A change fulfilled requests to `completed`. Consequently, completed procurement is excluded by the shown approved-only filter. It also does not distinguish `fund_transfer` from `actual_expense`; if both are represented for the same funding/purchase cycle and both meet the other filters, summing both can count the same funding and expenditure twice. These are query/workflow questions, not permission to change the supplied diagram.

### H. Project-specific profit is not established for wages

6S describes a project report. Its labor query sums all paid payrolls within the period without a project condition. The ER links PAYROLL to USER, not directly to CONTRACT_TOR or TOR_LOCATION_ASSIGNMENT. Work records can provide project context, but allocating an aggregated worker payroll across projects requires a rule that is not supplied here.

The revenue query additionally uses `billing_month = $2`, while the material and labor queries use date-period conditions. A month and a date range are different scopes. The references do not establish their mapping.

## 6. Other limits and consistency questions

- MANAGES links USER to CONTRACT_TOR. Area-level permission checks are described in several Assistant use cases, but the ER does not specify how contract management translates to permissions for individual assignments. The role field alone does not define that mapping.
- LEAVE_REQUEST visibly stores user_id but is associated with WORKER through SUBMITS. This can be compatible with a worker's account, once worker identity is resolved; it is not automatically an error.
- `reviewed_by` and REVIEWS both identify a review context. The field may implement that relationship. Showing both is not inherently duplication, but neither establishes review-history storage; only the shown review metadata is explicit.
- `work_date` occurs on both WORK_SCHEDULE and ATTENDANCE. This can be intentional, but consistency between the two values is not specified by the ER.
- Unique transfer references, one schedule per worker/day, and the attendance UPSERT key are rules in the use cases. The ER does not explicitly prove their physical enforcement. An actual schema would be needed to verify that.
- 4S reads check_in using schedule_id and worker_id, then updates by attendance_id. Its supplied check-in SELECT does not return attendance_id. It also uses a worker parameter position in the leave query's user_id predicate. The necessary identity/attendance context is not explicitly demonstrated by those query labels.
- Status values and transitions are spread across the use cases. A generic status attribute does not validate that the entire transition chain is complete. The pending/pending_survey mismatch above is a concrete example.

## 7. Consequences for a class-diagram draft

The ER supplies a strong starting inventory for **16 candidate domain classes**, with shown attributes and associations. The latest sequence diagrams supply the operations and participating controller/repository/service objects for a design-level diagram.

These should be used under the existing constraints:

1. Keep the draft simple. Do not turn every SQL query, return value, loop or branch into a class.
2. Preserve explicit names. Conflicting names go in an unresolved-items list rather than being silently combined.
3. Treat USER–WORKER inheritance as provisional until its identity and specialization meaning are confirmed. Do not automatically add Supervisor/Assistant subclasses.
4. Do not add foreign-key attributes merely because a conceptual relationship exists. Associations can represent those references; physical key fields require schema evidence or explicit use-case SQL.
5. Do not infer composition, aggregation, ownership/deletion rules, visibility or parameter/return data types from ordinary ER relationship lines.
6. Do not convert M automatically to 1..* or 0..* without evidence for minimum participation.
7. Use only operations actually shown in the latest sequence diagrams if a design-level class diagram is requested. An ER alone does not establish public methods.
8. Keep unresolved identity, payroll, deduction, evidence and expense mappings visible. These affect actual class structure, not just formatting.

The reference is sufficient for an initial documented draft, but not for claiming a fully reconciled physical schema or final implementation-ready class diagram.
