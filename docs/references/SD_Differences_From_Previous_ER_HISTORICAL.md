> **HISTORICAL — differences found against the previous ER; several are now corrected (docs/04 §1.3).**

# SD differences from the ER

ER is primary for domain names, attributes and associations. These differences remain visible; source ER and SD files have not been changed.

## Direct conflicts

| SD | Difference | Draft decision |
|---|---|---|
| 7W | PAYROLL uses total_wage, total_deduction, net_wage, status | Domain PAYROLL keeps base_wage, deduction, net_pay, is_paid |
| 7W | DEDUCTION_TRANSACTION.amount | Domain keeps penalty_amount |
| 2A | equipment_name is selected from REQUISITION_ITEM | Attribute stays on EQUIPMENT |
| 6S | EQUIPMENT_REQUISITION.expense_id join | Keep requisition-to-many-claims ER association; add no expense_id attribute |

## Differences that can be compatible but remain unresolved

- WORKER identity: SDs use worker_id/user_id; ER shows only is_available and its USER connection. No identifier fields or inheritance are invented.
- PAYROLL recipient: ER association is to USER, while SQL references WORKER. Keep the ER association.
- WORK_EVIDENCE: SDs use schedule_id; ER associates evidence with an assignment and a worker. No schedule association is silently added.
- Deduction attendance/payroll links: the ER expresses associations, but the 5S write flow does not demonstrate recording all links needed by 7W.
- SQL foreign keys can implement ER associations without needing separate ER ovals. Absence of an oval alone is not a contradiction.
- Derived net_pay may be computed and stored. Both can be compatible.

## Workflow differences separate from the ER

- 5W creates pending; 5A/6A select pending_survey. The references do not show the intervening transition.
- 6S selects approved equipment requests; 2A/7A move fulfilled requests to completed.
- 6S labor totals have no project condition. The ER does not establish a payroll-allocation rule across projects.
- 5S totals deductions by creation date, which is not necessarily the attendance work period, and its shown inserts do not establish a duplicate guard.

## Representation decisions

- SD operation names and parameters remain source-supported, including terminology that differs from ER fields. Such parameters do not add those names as entity attributes.
- Similar participants remain separate: for example, PayrollSummary Controller (5S) is not silently merged with Payroll Controller, and Contract Repository is not silently merged with TOR Repository.
- Automatic Job and EquipmentResult Controller have no received operations in the supplied messages. Their empty compartments are intentional.
- No software-to-domain arrows are invented from SQL table access alone. Software dependencies and ER business associations are distinct parts of the model.
- Exact ER spelling exat_deduction_amount is retained.
- Unknown minimum multiplicities, types and visibility are omitted. The draft does not claim a reconciled physical database schema.
