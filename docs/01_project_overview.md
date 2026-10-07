# 01 — Project overview

## What the system is
"Chaum Resource Management" (ระบบจัดการทรัพยากร ชะอุ่ม): a system for a tree-care and maintenance business that works under TOR contracts with a government agency (EXAT). It replaces an Excel/LINE/phone-based manual process.

It manages: contracts (TOR), work locations, equipment requirements and procurement, worker scheduling, leave and replacement, attendance, deductions, payroll, equipment shortages during work, and a project profit report.

## Repos and platforms
- **This repo: backend only.** Frontend is a separate repo (Next.js). Backend is built first.
- Backend: Go + Gin. Database: PostgreSQL on Supabase. Storage: Supabase Storage.
- Supervisor and Assistant use the **desktop web** UI.
- Workers use **only the LINE Mini App (LIFF)** on mobile. Worker workflows must not be moved to desktop screens.

## Roles (all are values of `USER.role`; they are NOT subclasses)
| Role | Thai | Interface | Typical actions |
|---|---|---|---|
| Supervisor | ผู้ควบคุมงาน | Web | Create TOR contract and scope, transfer funds, approve/reject requisitions, run daily attendance and payroll jobs, view financial report |
| Assistant | ผู้ดูแลงาน | Web | Site survey, procurement, create schedules, review leave and assign replacements, review/forward extra equipment requests, record extra purchases and delivery, check contract continuation |
| Worker | คนสวน | LINE Mini App | View schedule, request leave, substitute check-in, submit work report (also checks out), report missing equipment, view request result, view payslip |
| Replacement worker | — | LINE Mini App | A worker who covers an absent worker (substitute check-in, 3W) |

**Identity:** `user_id = worker_id`. A worker is a user. See `docs/04_decisions_and_corrections.md`.

## To-be flow in one paragraph
Supervisor signs a TOR and records it (contract, locations, required workers, required equipment). Assistant checks existing equipment against the TOR (site survey) and computes what to buy. Supervisor transfers funds; assistant buys with a receipt photo; procurement is cumulative until complete. Assistant schedules workers per location and date. Workers may request leave; the assistant approves/rejects and may assign a replacement who then checks in as a substitute. Workers check in/out via the LINE app (the work report with photo evidence performs check-out). A daily job classifies attendance (on time / late / leave / absent). A periodic job creates deductions (penalties) and payroll. If equipment runs short during work, a worker requests it; the assistant reviews, forwards or rejects, buys, and records delivery. The supervisor views a financial report: revenue minus labor minus material cost.

## Use-case codes
Sequence diagrams and use cases use codes: `1S–6S` (+`3.1S`), `1A–9A`, `1W–7W`.
A use case "7S Manage Location" (Supervisor adds/edits work sites tied to a TOR) exists in the business flow but is **not** in the sequence-diagram inventory, so no operations are specified for it here. Treat it as unspecified; ask before implementing beyond basic location CRUD that 1S already needs.

## Notifications (LINE)
Notifications go through a LINE notification service. Message types seen in the SDs: new contract, procurement funding, shortage, approval/rejection of a request, guidance to the original requester, delivery summary, leave notice to assistant, additional request notice, equipment received push, no-purchase guidance push, assigned-worker schedule notice, leave decision notices (leaving worker, replacement worker), missing replacement notice to supervisor. Exact operations are in `docs/03_workflows_operations.md`.

## Team/document context (informational)
The class and sequence diagrams are university systems-analysis deliverables for this project. They are the specification; there is no other source.
