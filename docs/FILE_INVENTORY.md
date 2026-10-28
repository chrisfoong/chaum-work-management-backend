# Repository file inventory

Updated 2026-10-15 on branch `feature/repo-tidy`. 160 files, including this one and the local-only `.env` and `docs/references/`. Excluded: `.git/` internals. There is no `node_modules/` and no vendored dependencies.

## Root

| Path | Purpose |
|---|---|
| `.env` | Local secret runtime configuration (Supabase, LINE, database, storage, CORS, QR secret); never read, log or commit it. |
| `.env.example` | Committed placeholder template that documents every environment variable the server and tests expect. |
| `.gitignore` | Git ignore rules that keep secrets, build output and local artifacts out of version control. |
| `AGENTS.md` | Authoritative instructions for AI agents: business/schema authorities, the `.env` ban and the no-DDL-on-Supabase rule. |
| `CLAUDE.md` | Claude Code entry point that redirects to AGENTS.md, docs/00_INDEX.md and docs/TASKS.md. |
| `README.md` | Project overview with run instructions and links to requirements, task tracker, API notes and OpenAPI spec. |
| `go.mod` | Go module definition (`chrisfoong/chaum-work-management-backend`) and direct dependency list. |
| `go.sum` | Checksums that lock Go dependency versions. |
| `docker-compose.test.yml` | Throwaway tmpfs PostgreSQL container for local integration tests only; the password comes from `TEST_DB_PASSWORD`. |
| `Dockerfile` | Multi-stage Go 1.26 build of `./cmd/server` into a non-root Alpine runtime image with CA certificates and timezone data. |
| `.dockerignore` | Keeps `.git`, `.env` files, build output and `docs/references` out of the Docker build context. |
| `compose.yaml` | Docker Compose service that runs the backend image with `.env`, bound to localhost, with a health check. |
| `DEV_DOCKER.md` | Thai-language guide to running the backend with `go run` or Docker Compose. |

## `.github/`

| Path | Purpose |
|---|---|
| `.github/workflows/backend.yml` | GitHub Actions "Backend verification" workflow that runs format, vet, test and build checks. |

## `api/`

| Path | Purpose |
|---|---|
| `api/openapi.json` | Generated OpenAPI document (`go run ./cmd/apidoc > api/openapi.json`) listing the actual routes, request models and LINE ID-token bearer auth. |

## `cmd/`

| Path | Purpose |
|---|---|
| `cmd/apidoc/main.go` | CLI that builds the real route inventory and prints it as JSON for `api/openapi.json`. |
| `cmd/apidoc/models.go` | Request-model reflection helpers that the apidoc generator uses to describe request bodies. |
| `cmd/server/main.go` | The only server entry point (`go run ./cmd/server`); it calls `internal/application`. |

## `migrations/` (historical: not applied, never run against Supabase)

| Path | Purpose |
|---|---|
| `migrations/0001_init.up.sql` | Unapplied DDL draft for the 16 production-ER tables, derived from the data dictionary. |
| `migrations/0001_init.down.sql` | Unapplied, destructive rollback that drops every table 0001_init creates. |

## `internal/`

### `internal/application/`
| Path | Purpose |
|---|---|
| `internal/application/run.go` | Application bootstrap (`Main`/`Run`) that loads config, opens the DB pool, wires services and starts the HTTP server. |
| `internal/application/notify.go` | Adapter that connects the application's notification sender to the LINE messaging implementation. |

### `internal/apperr/`
| Path | Purpose |
|---|---|
| `internal/apperr/apperr.go` | Typed application errors (Validation, Unauthorized, Forbidden, NotFound, Conflict, Internal) and their HTTP status mapping. |
| `internal/apperr/apperr_test.go` | Tests for the HTTP status mapping and the wrapping of internal causes. |

### `internal/auth/`
| Path | Purpose |
|---|---|
| `internal/auth/auth.go` | Gin middleware that authenticates bearer tokens, stores the principal and enforces required roles. |
| `internal/auth/identity.go` | Onboarding-only middleware that verifies a LINE identity without granting an application principal. |
| `internal/auth/line.go` | LINE Login ID-token verifier that calls LINE's verification endpoint. |
| `internal/auth/supabase.go` | Supabase Auth access-token (JWT) verifier. |
| `internal/auth/users.go` | USER table lookups by user ID or LINE ID, used to resolve roles. |
| `internal/auth/auth_test.go` | Unit tests for authentication middleware, role checks, both verifiers and the user store. |
| `internal/auth/login_integration_test.go` | Integration test that runs the full pipeline: LINE verifier, USER lookup and role middleware. |

### `internal/config/`
| Path | Purpose |
|---|---|
| `internal/config/config.go` | Loads base runtime settings (port, DB URL and others) from environment variables. |
| `internal/config/runtime.go` | Loads and validates MVP deployment config (LINE, Storage, CORS, QR secret) and fails closed when it is invalid. |
| `internal/config/supabase_url.go` | Resolves the Supabase project URL from `SUPABASE_URL` or the `NEXT_PUBLIC_SUPABASE_URL` alias and rejects mismatches. |
| `internal/config/config_test.go` | Tests for loading environment configuration. |
| `internal/config/runtime_test.go` | Tests that confirm runtime config fails closed when required values are missing or bad. |
| `internal/config/supabase_url_test.go` | Tests that both config loaders accept the Supabase URL alias and reject conflicting values. |

### `internal/contract/` (1S: TOR contract and work scope)
| Path | Purpose |
|---|---|
| `internal/contract/model.go` | Request and domain models for 1S contract creation, with input normalisation. |
| `internal/contract/validate.go` | Format validation for contract info and scope data (uc 1S steps 4 and 7). |
| `internal/contract/service.go` | 1S service that submits info and scope, confirms contracts in a transaction, creates areas and sends notifications. |
| `internal/contract/repository.go` | SQL repository for contracts, locations, TOR assignments, requisition numbering, equipment and location search. |
| `internal/contract/handler.go` | HTTP handlers that expose the 1S operations and location search. |
| `internal/contract/handler_test.go` | Handler tests that use a recording fake store. |
| `internal/contract/validate_test.go` | Tests for contract and scope format validation. |
| `internal/contract/file_test.go` | Test confirming that a contract file must be present and verified before any write. |
| `internal/contract/testpool_test.go` | Test helper that opens the isolated test database pool. |
| `internal/contract/repository_integration_test.go` | Integration tests for unique and foreign-key violations, equipment find-or-create and requisition numbering. |
| `internal/contract/service_integration_test.go` | Integration tests for concurrent contract confirms and for lock release on cancel. |
| `internal/contract/location_coordinates_integration_test.go` | Integration test confirming that unassigned locations still expose their coordinates. |
| `internal/contract/new_location_coordinates_test.go` | Validation and integration tests for coordinates on newly created TOR locations. |

### `internal/db/`
| Path | Purpose |
|---|---|
| `internal/db/db.go` | PostgreSQL pgx pool creation, the `WithTx` transaction helper and unique/foreign-key violation detection. |
| `internal/db/db_test.go` | Tests for transaction commit, rollback, panic handling and violation detection. |

### `internal/health/`
| Path | Purpose |
|---|---|
| `internal/health/health.go` | Liveness and readiness endpoint that pings the database. |
| `internal/health/health_test.go` | Test for the health handler. |

### `internal/httpx/`
| Path | Purpose |
|---|---|
| `internal/httpx/httpx.go` | Shared Gin helpers: JSON binding with field-name errors, error responses and the request logger. |
| `internal/httpx/httpx_test.go` | Tests for error writing and JSON binding. |

### `internal/notify/`
| Path | Purpose |
|---|---|
| `internal/notify/notify.go` | LINE notification port, with no-op and disabled implementations. |
| `internal/notify/notify_test.go` | Tests confirming that the no-op and disabled notifiers never report a send. |

### `internal/schema/`
| Path | Purpose |
|---|---|
| `internal/schema/check.go` | Read-only startup check that compares the live DB metadata with the embedded Supabase export. |
| `internal/schema/columns.json` | Embedded copy of the Supabase column export (identical to docs/schema). |
| `internal/schema/constraints.csv` | Embedded copy of the Supabase constraint export (identical to docs/schema). |
| `internal/schema/sync_test.go` | Test that fails if the embedded exports drift from the `docs/schema/` authority. |
| `internal/schema/isolated-test-schema.sql` | DDL used only to build the isolated local test database; never run it on Supabase. |

### `internal/server/`
| Path | Purpose |
|---|---|
| `internal/server/server.go` | Builds the Gin engine and the authenticated route groups (`/api/web`, `/api/liff`). |
| `internal/server/cors.go` | Strict CORS allowlist middleware that never reflects arbitrary origins. |
| `internal/server/server_test.go` | Test confirming that each router group applies its own auth. |
| `internal/server/cors_test.go` | Tests for CORS allowlist behaviour. |

### `internal/storage/`
| Path | Purpose |
|---|---|
| `internal/storage/storage.go` | File-upload port for Supabase Storage, plus an in-memory fake. |
| `internal/storage/storage_test.go` | Test for the fake upload implementation. |

### `internal/work/` (confirmed MVP workflows)

Services and handlers:
| Path | Purpose |
|---|---|
| `internal/work/service.go` | Core work service: transactions, listing and paging, locking, validation helpers and exact money/decimal types. |
| `internal/work/handler.go` | Registers every work-workflow HTTP route on the web and LIFF groups. |
| `internal/work/http_helpers.go` | Shared handler helpers for responses, actor extraction, input binding and pagination. |
| `internal/work/attendance.go` | Signed QR issuance and checks, GPS checks, check-in/check-out and attendance finalisation. |
| `internal/work/catalog.go` | Read endpoints for contracts, contract state, locations and equipment. |
| `internal/work/dashboard.go` | Dashboard aggregate read. |
| `internal/work/events.go` | Best-effort post-commit event notifications (the schema has no durable outbox). |
| `internal/work/export.go` | Profit export endpoint. |
| `internal/work/external.go` | Supabase Storage and LINE push clients, with owner checks, PNG/media verification and keyed sends. |
| `internal/work/frontend.go` | Frontend read models: filtered lists, available workers, leave candidates, assignments, delivery schedules and message splitting. |
| `internal/work/notification_retry.go` | Retries a notification from committed business rows only, without replaying writes. |
| `internal/work/operations_summary.go` | 9A Assistant operations summary and the Supervisor's manual notification, with financial fields excluded. |
| `internal/work/payroll.go` | Payroll preview and creation, atomic period batches, payment, monthly views, and invoice issue and receipt. |
| `internal/work/people.go` | User create/update, `/me`, schedules, leave requests, and atomic leave review with replacement (4A). |
| `internal/work/invitation.go` | LINE invitation onboarding: create, renew, preview and claim invitation links using the existing USER schema (token digest only). |
| `internal/work/procurement.go` | Requisitions, surveys, request inspection and decision (5A), funding, and purchase preview and execution. |
| `internal/work/report.go` | Profit and invoice report computation, date-range profit and the 6S closing snapshot. |
| `internal/work/delivery.go` | 8A delivery recording and listing via the WORK_EVIDENCE description convention, plus evidence-photo checks. |
| `internal/work/continuation.go` | Contract and 9A selected-area continuation read models (`can_continue`). |
| `internal/work/tor_equipment.go` | Read of active initial TOR equipment for an assignment, filtered by owner. |
| `internal/work/report_http.go` | HTTP report endpoints, including a dependency-free financial profit PDF generator. |

Repositories (SQL against the exported schema):
| Path | Purpose |
|---|---|
| `internal/work/repository_queries.go` | Shared Row/Rows/Exec query primitives; SQL lives in repositories and services choose the order. |
| `internal/work/repository_attendance.go` | SQL for attendance. |
| `internal/work/repository_catalog.go` | SQL for catalog reads. |
| `internal/work/repository_dashboard.go` | SQL for dashboard aggregates. |
| `internal/work/repository_events.go` | SQL for event/notification data. |
| `internal/work/repository_frontend.go` | SQL for frontend read models and filters. |
| `internal/work/repository_payroll.go` | SQL for payroll and invoices. |
| `internal/work/repository_people.go` | SQL for users, schedules and leave. |
| `internal/work/repository_procurement.go` | SQL for requisitions, funding and purchases. |
| `internal/work/repository_report.go` | SQL for profit and invoice reports. |

Tests:
| Path | Purpose |
|---|---|
| `internal/work/policy_test.go` | Unit tests for exact money, penalty boundaries, periods and overnight shifts, GPS+QR and file ownership. |
| `internal/work/authorization_test.go` | Tests for mutation role boundaries and rejection of requests without an identity. |
| `internal/work/contract_file_test.go` | Test that contract PNG verification checks the bytes, owner and size. |
| `internal/work/external_test.go` | Tests that Storage rejects other owners and redirects, and that LINE push authenticates and surfaces failures. |
| `internal/work/frontend_test.go` | Tests for frontend filter validation, Unicode message splitting and read-model delivery guards. |
| `internal/work/usecases_test.go` | Tests for profit PDF pagination and fields, and for rejecting invalid use-case input before the DB. |
| `internal/work/completion_extra_test.go` | Tests for bounded LINE retries, continuation and notification HTTP, and unfunded additional purchases. |
| `internal/work/usecase_completion_test.go` | Tests for the worker equipment reason limit and integration of use-case completion. |
| `internal/work/integration_test.go` | End-to-end MVP integration test on the isolated database pool. |
| `internal/work/usecases_integration_test.go` | Integration tests for emergency leave with atomic replacement, and for purchase rounds, delivery and the paid report. |
| `internal/work/concurrent_integration_test.go` | Integration test for atomic concurrent attendance and rollback. |
| `internal/work/login_http_integration_test.go` | Integration test for `/me` handlers, DB role resolution and the Worker mapping gate. |
| `internal/work/operations_summary_test.go` | Integration tests for the 9A operations summary and its manual notification. |
| `internal/work/overnight_read_test.go` | Integration test confirming that workers can read ongoing overnight shifts. |
| `internal/work/checkout_end_test.go` | Integration test confirming checkout waits for the scheduled end of an overnight shift. |
| `internal/work/invitation_test.go` | Unit and integration tests for invitation tokens, Supervisor-only creation, forged identity rejection, cancel and reissue. |
| `internal/work/onboarding_platform_test.go` | Integration test confirming onboarding routes users by their verified database role. |
| `internal/work/payroll_role_history_test.go` | Integration test confirming payroll keeps scheduled work after a user role change. |
| `internal/work/tor_equipment_test.go` | Integration test for assignment equipment filtering by TOR and owner. |

## `docs/`

### Current docs
| Path | Purpose |
|---|---|
| `docs/00_INDEX.md` | Documentation index that marks older docs and migrations as historical. |
| `docs/03_workflows_operations.md` | Current workflow-operation summary under the adopted SA-Group5 use cases. |
| `docs/07_verification_checklist.md` | Current verification commands (gofmt, test, vet, build) and integration-test prerequisites. |
| `docs/9A_DESIGN_AND_LIVE_TEST_PLAN.md` | Confirmed 9A MVP design (16-table schema kept) and the live acceptance test plan. |
| `docs/API.md` | API notes covering the web and LIFF route groups and auth, and how to regenerate the OpenAPI spec. |
| `docs/BACKEND_ACCEPTANCE.md` | Backend acceptance evidence boundaries and live-configuration steps. |
| `docs/FINAL_USECASE_REVIEW.md` | Final 2026-10-09 review of the backend against the SA-Group5 use cases. |
| `docs/ONBOARDING.md` | Thai-language guide to adding users with LINE invitation links without changing the ER. |
| `docs/FRONTEND_HANDOFF.md` | Frontend handoff that maps use cases to backend routes and auth requirements. |
| `docs/REQUIREMENTS.md` | Authoritative backend requirements confirmed by the user on 2026-10-08. |
| `docs/SERVER_AND_TEST_STATUS.md` | Thai-language guide to starting the server, with the current test status. |
| `docs/TASKS.md` | Canonical task tracker with implemented, verified and evidence sections. |
| `docs/UI_ACCEPTANCE_20261010.md` | Thai-language 2026-10-10 use-case and screen acceptance results, with what remains unverified. |
| `docs/USECASE_AUDIT.md` | 2026-10-08 audit of use-case function coverage. |
| `docs/VERIFICATION_AND_PUBLISH.md` | Final verification results (177 tests, coverage figures) and the publication attempt log. |
| `docs/FILE_INVENTORY.md` | This file: an inventory of every repository file with a one-line purpose. |

### `docs/schema/` (physical schema authority)
| Path | Purpose |
|---|---|
| `docs/schema/columns.json` | Supabase export of every table column (16 tables, 121 columns) and its type. |
| `docs/schema/constraints.csv` | Supabase export of every table constraint (PK, FK, UNIQUE, CHECK). |

### `docs/historical/` (archived 2026-10-05 design: do not apply)
| Path | Purpose |
|---|---|
| `docs/historical/01_project_overview.md` | Archived system overview: contracts, locations, procurement, scheduling, attendance and so on. |
| `docs/historical/02_domain_model.md` | Archived production-ER entity and attribute model. |
| `docs/historical/03_workflows_operations.md` | Archived per-sequence-diagram workflow operation inventory. |
| `docs/historical/04_decisions_and_corrections.md` | Archived adopted decisions and open items (for example `user_id = worker_id`). |
| `docs/historical/05_business_rules.md` | Archived business rules for quantities, procurement, scheduling and other areas. |
| `docs/historical/06_architecture_and_conventions.md` | Archived layering and conventions proposal. |
| `docs/historical/07_verification_checklist.md` | Archived verification checklist and acceptance checks. |
| `docs/historical/TASKS-2026-10-05.md` | Archived 2026-10-05 task tracker snapshot. |

### `docs/historical/pr/` (past PR and review working notes)
| Path | Purpose |
|---|---|
| `docs/historical/pr/COPILOT_REVIEW.md` | Prompt for an external Copilot review of the feature stack. |
| `docs/historical/pr/FEATURE_PRS.md` | Status of the stacked feature branches and draft-PR preparation. |
| `docs/historical/pr/PR_PLAN.md` | Draft PR delivery plan for the stacked branches after the GitHub 403 push rejection. |
| `docs/historical/pr/PR_USECASE_DESCRIPTION.md` | PR description text for the use-case-alignment changes and their validation. |

### `docs/prompts/`
| Path | Purpose |
|---|---|
| `docs/prompts/CLAUDE_Prompts_Backend.md` | Ordered, plan-first prompt sequence for driving Claude in VS Code on this repo. |

### `docs/references/`
| Path | Purpose |
|---|---|
| `docs/references/uc.txt` | Thai use-case descriptions, sequence-diagram text and collaboration-diagram text (from 3A onward). |
| `docs/references/datadict.txt` | Team data dictionary with physical names, types and constraints. |
| `docs/references/er-latest.png` | Image of the latest ER diagram (2274×1668). |
| `docs/references/Chaum-Diagrams-PlantUML.txt` | PlantUML sources for the class diagrams (sequence diagrams exist only as images). |
| `docs/references/Desktop_App_Handoff_CURRENT.md` | Current self-contained handoff on the production ER and the corrected sequence diagrams. |
| `docs/references/Production_ER_Conflict_Review.md` | Comparison of the production ER with the previous ER and the team-edited sequence diagrams. |
| `docs/references/Class_Package_Verification_HISTORICAL.md` | Historical visual verification notes for the class-diagram package. |
| `docs/references/Class_Reference_Mapping_HISTORICAL.md` | Historical mapping of 64 classes (16 ER and 48 SD participants) to their operations. |
| `docs/references/Previous_ER_Detailed_Analysis_HISTORICAL.md` | Historical detailed analysis of the previous draw.io ER. |
| `docs/references/SD_Differences_From_Previous_ER_HISTORICAL.md` | Historical list of conflicts between the sequence diagrams and the previous ER. |

## Notes

- `internal/schema/columns.json` and `internal/schema/constraints.csv` must stay byte-identical to the `docs/schema/` copies; `internal/schema/sync_test.go` enforces this.
- `.env` and `docs/references/` are gitignored and exist only locally.
- Keep this file current: CI (`.github/workflows/backend.yml`) fails when a tracked file is missing here or a listed, non-ignored path no longer exists.
