> **REFERENCE COPY — current handoff written by the team for the desktop app. Authority order and corrections here match docs/04. It mentions files (Production_ER_SD_Corrections.zip, Whole_Class_Diagram_Five_Pages.zip, the production ER svg) that are NOT available in this repo; their relevant content is summarized in docs/02 and docs/04.**

# Desktop Application Handoff — Production ER and Corrected Sequence Diagrams

## Purpose

Read this document together with the supplied diagrams, use-case descriptions and repository. It is self-contained; no prior conversation is needed.

The application manages contracts, work locations, equipment procurement, staffing, leave, attendance, wages and financial reporting. The target is a desktop application for the applicable Supervisor and Assistant workflows. Worker workflows currently use a LINE Mini App; do not silently replace that interface with desktop screens.

**If the request is “Read and analyze,” analyze the actual supplied files and produce the outputs below. That phrase alone does not authorize generating or changing application code.** When implementation is explicitly requested, proceed with sufficiently specified work and identify decisions that block particular workflows.

## Current references and authority

| Reference | Role and current status |
|---|---|
| `SA1-ER PRODUCTION.drawio.svg` | Current authority for canonical domain names and intended business associations. Supersedes the older ER. |
| `Production_ER_SD_Corrections.zip` | Current corrected versions of SD-2A, SD-5S, SD-6S, SD-5W and SD-7W. Use the `Fixed/` versions for these cases. |
| Team sequence diagrams from `SA1-SD-1A.drawio.svg` or `Team_Edited_Sequence_Diagrams.zip` | Current team baseline for other SDs. Five affected baseline pages are superseded by the correction package. |
| Supplied use-case tables | Full behavior, validations, branches, formulas and stated transaction requirements. The listed corrections below intentionally supersede their conflicting fields/queries/statuses. |
| `Class_Diagram_ER_First.zip` | Historical class-design mapping built from the old ER and team SD baseline. Useful for responsibilities and provenance, but not fully synchronized with current references. |
| `Whole_Class_Diagram_Five_Pages.zip`, if supplied | Compact document views of the provided whole class diagram. These are presentation exports, not a subsequent domain-model reconciliation. |
| `Production_ER_Conflict_Review.md` | Audit performed before the correction package. Its seven targeted findings are now addressed as listed below; other findings remain open. |
| Actual schema/migrations and repository | Authority for existing physical types, constraints, storage mappings and stack. Report discrepancies with the intended ER; do not silently migrate a working database. |

Explicit current user requirements take priority. The correction package supersedes older SD/use-case wording only for its documented changes. Preserve other stated business behavior. Do not select an old field name just because the historical class diagram still shows it.

A conceptual ER does not establish every physical FK, database type, default, minimum participation or authorization rule. Missing FK ovals are not automatically contradictions. A collaboration/communication diagram is optional and does not supersede the ER.

## Reading order and package paths

1. Inspect the production ER and actual database schema, if supplied.
2. Extract `Production_ER_SD_Corrections.zip`; its root is `Production_ER_SD_Corrections/`.
3. Read `README.md`, `Changes_and_SQL.md` and `Verification.md` in that root.
4. Read the five corrected draw.io files under `Fixed/`, using the role/use-case folder structure. `Old/` is comparison material only. Each case also contains two Google Docs images.
5. Inspect remaining team SDs and the relevant full use-case tables. SQL containing `...` is abbreviated documentation, not executable SQL.
6. Use the class packages to map participants and responsibilities, applying the current-reference overrides below.
7. Read the remaining issues in this document and the conflict review.

The correction package contains five cases, not the complete SD set:

- `Fixed/Assistant(A)/2A/SD-2A.drawio`
- `Fixed/Supervisor(S)/5S/SD-5S.drawio`
- `Fixed/Supervisor(S)/6S/SD-6S.drawio`
- `Fixed/Worker(W)/5W/SD-5W.drawio`
- `Fixed/Worker(W)/7W/SD-7W.drawio`

The historical class package root is `Class_Diagram_ER_First/`. Useful files include:

- `Drawio/Class_Diagram_Full.drawio`: complete historical editable model.
- `Plan_and_Check/Reference_Mapping.md`: original class/member provenance.
- `Plan_and_Check/SD_Differences_From_ER.md`: historical differences, some now corrected.
- `Plan_and_Check/Latest_ER_Detailed_Analysis.md`: analysis of the previous ER, not production ER.
- `Plan_and_Check/Verification.md`: diagram/export checks, not application tests.

The compact package root is `Whole_Class_Diagram_Five_Pages/`. `Reference/Full_Class_Signatures.md` and `Drawio/Whole_Class_Diagram_Full_Signatures.drawio` preserve the source signatures. Compact images show method names with omitted parameters. These full signatures still reflect the source model before the latest SD corrections.

If a file is absent, say so. Never claim inspection based only on this handoff or an image filename.

## Production ER changes to apply

The same 16 entities remain:

| Area | Entities |
|---|---|
| Accounts | USER, WORKER |
| Contracts and locations | CONTRACT_TOR, LOCATION, TOR_LOCATION_ASSIGNMENT |
| Work and leave | WORK_SCHEDULE, ATTENDANCE, LEAVE_REQUEST, WORK_EVIDENCE |
| Procurement | EQUIPMENT, EQUIPMENT_REQUISITION, REQUISITION_ITEM, EXPENSE_CLAIM |
| Finance | PAYROLL, DEDUCTION_TRANSACTION, COMPANY_INVOICE |

Changes from the previous ER:

- `PAYROLL.deduction` is now **`total_deduction`**.
- `PAYROLL.net_pay` is now **`net_wage`**, still shown as derived.
- PAYROLL retains **`base_wage`** and **`is_paid`**. Do not introduce `total_wage` or a separate payroll `status` column merely to match old SD wording.
- The explicit `LEAVE_REQUEST.user_id` oval was removed. The WORKER-to-LEAVE_REQUEST SUBMITS relationship remains. Physical identity/FK mapping is still a decision.
- Added **USER — PAYS — PAYROLL**, shown 1:M, alongside USER — RECEIVES — PAYROLL. Keep payer and recipient roles distinct; their exact identifiers and authorization rules are not supplied.
- Repaired the HOSTS connector binding to TOR_LOCATION_ASSIGNMENT without changing the intended LOCATION-to-assignment relationship.

CONTRACT_TOR.status is visibly present. Some embedded connectors still bind to cardinality labels or use fixed endpoints: CONTRACT_TOR.status, ATTENDANCE-to-deduction TRIGGERS and USER-to-requisition REQUESTS. Treat these as drawing defects to report, not proof that the intended fields/relationships are absent.

Preserve the source spelling `COMPANY_INVOICE.exat_deduction_amount` in mappings. Any intentional code/schema spelling correction needs an explicit old/new mapping and migration plan.

## Corrections already made — do not treat these as unanswered questions

These are the adopted diagram corrections, not a claim that application code or database migrations already exist.

| Case | Current correction |
|---|---|
| 7W payslip fields | Reads `base_wage`, `total_deduction`, `net_wage`, `is_paid`. Payslip display uses base_wage. |
| 7W deduction amount | Reads `DEDUCTION_TRANSACTION.penalty_amount`. |
| 2A equipment lookup | Joins REQUISITION_ITEM to EQUIPMENT through equipment_id to read equipment_name. |
| 6S expense association | Joins `ec.requisition_id = er.requisition_id`, supporting multiple claims per requisition. |
| 5W initial request status | Creates additional requests as **pending_survey**, matching the existing 5A/6A review entry. No extra status-transition step was added. This replaces the original 5W pending value. |
| 6S material costs | Includes requisitions with status **approved, pending_supervisor, pending_procurement or completed**, counting existing actual purchase claims even while procurement is incomplete. Uses **expense_type = 'actual_expense'**, excluding funding transfers from this purchase-cost total. |
| 5S deduction attendance link | Carries attendance_id from the existing penalty result into `createDeduction(workerId, attendanceId, penaltyAmount, penaltyReason)` and inserts worker_id, attendance_id, penalty_amount and reason. No additional query or participant was added. |
| 5S payroll write | Explicitly inserts worker_id, period_start, period_end, base_wage, total_deduction, net_wage and is_paid=false. The calculation/return names now use calculateNetWage/netWage. |
| 6S labor amount | Sums `net_wage`, replacing `net_pay`. |

Use `Changes_and_SQL.md` for exact before/after wording and corrected SQL. Foreign keys used by these queries implement ER relationships; they do not add new business entities. The attendance FK mapping is now part of the correction proposal, but its physical constraint must be checked against the actual schema.

The correction package preserves original versions and changes 14 labels across five diagrams. Participants, message counts and nested-fragment counts are unchanged. The package does not edit the use-case tables, production ER, class diagrams or earlier handoff files. Reconcile those artifacts explicitly if implementation requires synchronized documentation.

## How to use the class design

The historical model contains 64 logical classes: 16 domain entities and 48 software participants. Its 354 supported operation signatures and relationship inventory describe the earlier source snapshot. They are provenance counts, not proof that the current model has been reconciled.

Repeated classes on document pages are context views, not extra classes. Apply the payroll-field changes, removal of the explicit leave user_id attribute, added PAYS association and corrected SD operation signatures before generating a current class/schema mapping.

Responsibilities:

- Boundary participants represent Web UI and LINE Mini App interactions.
- Controllers coordinate the stated actions and validations.
- Repositories implement supported retrieval/persistence operations.
- DB Connector, LINE Notification Service and explicitly shown jobs represent infrastructure/integration responsibilities.

These are design responsibilities, not a mandatory one-file-per-class architecture. Necessary framework adapters, view models and transaction utilities may be proposed as technical implementation choices. Do not present them as new business requirements from the diagrams.

Do not automatically merge similarly named participants. PayrollSummary Controller (5S) and Payroll Controller, for example, remain distinct in the source mapping. Any implementation consolidation needs a responsibility mapping.

Notation limits:

- 1/* marks maximum cardinality; minimum participation is unspecified.
- USER–WORKER `(p,e)` does not clearly settle inheritance or shared versus separate identifiers.
- Underlined/{key} identifiers do not establish types or actual constraints.
- `/net_wage` is derived; storing a computed value can still be valid.
- Dashed software dependencies reflect supported calls, not necessarily stored object-reference fields.
- Ellipses indicate omitted documentation content, not parameters or executable SQL.
- Types, visibility and empty operation compartments leave implementation details open.

## Remaining material decisions

Keep these separate from the seven corrected findings.

### Identity and persistence

- **USER/WORKER identity:** existing SDs use worker_id and user_id, while the ER leaves the worker identity strategy unclear. Establish shared or separate identity from schema evidence.
- **Payroll roles:** define recipient mapping and the new payer role. Corrected SQL retains the existing worker_id mapping; that does not settle the PAYS relationship.
- **Leave ownership:** establish the physical FK implementing WORKER SUBMITS LEAVE_REQUEST. Existing 2W/4S SQL still uses user_id; removing its ER oval alone does not make the query invalid.
- **Evidence context:** ER links evidence to assignment and worker; 4W/8A use schedule_id. Define whether assignment is inferred through schedule or stored separately.
- **Deduction-to-payroll link:** attendance linkage is now carried by 5S, but the ER’s PAYROLL APPLIES deductions relationship still needs a physical mapping/population rule.
- **Permissions:** role values alone do not establish project/area access. Use the supplied authorization mapping.
- **Replacement attendance:** schedule identity, replacement flow and the compound attendance UPSERT key must agree with physical uniqueness constraints.

### Workflow and financial scope

- **Approval evidence:** 7A relies on reviewed_by/reviewed_at/prior approval, while shown 3S approval only updates status. Establish how approval evidence is recorded.
- **Funded additional requests:** 2S sets pending_procurement, while 7A lists approved or its specific funded pending_supervisor path. Resolve the intended path without inventing an extra approval stage.
- **Project labor allocation:** 6S still sums paid payroll for a date interval without a project condition. Do not attribute all projects’ wages to the selected project.
- **Reporting dates:** revenue uses billing_month, while other totals use date/period ranges. Establish conversion and inclusion rules.
- **Deduction period:** 5S still sums by created_at; attendance dates can differ. Define attribution and safe repeat processing.
- **Payment completion:** 5S creates unpaid payroll; 6S selects paid payroll and 7W assumes money has been transferred. A payment process may exist elsewhere; do not claim it is shown here.
- **Financial edge cases:** negative net wages, refunds and cancellation/reversal behavior are not defined by the current corrections. Do not silently add rules.

These decisions should not stop unrelated specified workflows. Identify exactly which implementation work depends on each unresolved decision.

## Desktop scope and implementation rules

Map applicable Supervisor/Assistant actions to desktop views. Preserve worker LINE workflows and integration contracts unless explicitly asked to replace them. Diagram page layout is documentation layout, not a desktop UI specification.

Inspect the existing stack and schema before proposing a framework, dependencies or migration. If absent, identify operating-system target, database deployment, network assumptions and desktop framework as decisions. Do not assume offline synchronization or local-only storage.

Prefer simple implementations that preserve stated behavior. Do not add approval stages, repeated login steps, extra queries, nested workflow boxes or new business features merely to fill design gaps. Technical atomicity can be justified when required to preserve related writes; distinguish it from a source-stated transaction boundary.

Implement monetary precision according to the chosen schema/stack. Preserve documented formulas and deduction thresholds. Integrations must use supplied configuration and current official API documentation; simulated notifications are not proof of real delivery.

## Required response to “Read and analyze”

1. List the files actually inspected, their versions and any missing references.
2. Explain the business roles, interfaces, layers and principal relationships.
3. Produce an updated entity/schema and operation mapping, showing where historical classes are overridden by current references.
4. State which listed conflicts are resolved in the diagrams and which remaining decisions affect implementation.
5. Propose an application structure compatible with the existing stack; mark any new technology choice as a proposal.
6. Give an ordered implementation plan with acceptance checks.
7. Ask only for material missing decisions. Do not ask again whether to adopt the already documented corrections.

When implementation is authorized, inspect the repository, establish mappings, implement a small complete workflow, then extend in dependency order. Validate behavior rather than creating empty files for every diagram class.

## Acceptance checks

- Role/area access, preconditions, required-field/format/database validation and error behavior match their use cases.
- 5W requests are visible in the existing 5A/6A review queue.
- 2A equipment names come from EQUIPMENT; quantities and procurement updates remain correct.
- Expense claims link to the original requisition, permitting multiple claims.
- Completed and partial actual purchases appear in material costs; funding transfers are not counted again as purchase costs.
- Each 5S deduction carries the source attendance ID and is retrievable with its work date in 7W.
- 5S writes and 7W reads base_wage, total_deduction and net_wage consistently; payment state uses is_paid.
- 6S labor totals use net_wage, with project/date scope validated once those rules are established.
- Worker, schedule, attendance, leave, evidence and payroll identities agree with the adopted physical schema.
- Stated transaction boundaries, notifications and recipient selection are preserved.
- Reprocessing does not leave partial or duplicate records under the agreed constraints.

The correction package verified artifact consistency and selected SQL behaviors with SQLite-compatible fixtures. It did not validate a deployed PostgreSQL schema, application code, authorization or real integrations. Treat its verification as documentation/query evidence, not completed application testing.

When reporting implementation completion, state the working workflows, executed checks, adopted decisions and remaining concrete gaps. Having code files named after all classes is not sufficient evidence of completion.
