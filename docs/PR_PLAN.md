# Draft PR delivery plan

Local commits are complete. Direct push was rejected by GitHub (403): Sxthxwit has no Write permission on chrisfoong/chaum-work-management-backend. No remote feature branches or Draft PRs have been created. No merge/deploy/migrations occurred.

Branches are stacked. Each builds on the preceding branch; service features become publicly routed in the final runtime/API PR. Review and merge in this order after publication.

| Branch | PR base | Local commit |
| --- | --- | --- |
| feature/schema-line-auth | develop | e23087a |
| feature/catalog-contracts | feature/schema-line-auth | 2fffcaf |
| feature/scheduling-leave | feature/catalog-contracts | 5bd404f |
| feature/attendance | feature/scheduling-leave | cfc76ec |
| feature/procurement | feature/attendance | 1c1dd3b |
| feature/payroll | feature/procurement | e4e9a29 |
| feature/reports | feature/payroll | c875fd1 |
| feature/runtime-api | feature/reports | 4cadd01 |

## Reconcile Supabase schema and verify LINE identities

Replace migration assumptions with exported Supabase metadata, quoted public."USER" and correct Worker/Payroll references. Harden LINE issuer/audience/expiry and DB identity checks; repair the preexisting contract integration build failures. Final runtime wiring follows in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base develop. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Add catalog and contract views

Add contract views/status updates, location/equipment management and dashboards using existing columns. Names are guarded in application transactions where live schema has no matching uniqueness constraint. HTTP registration follows in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/schema-line-auth. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Add scheduling, leave review and replacement services

Add staffing capacity/availability guards, 24-hour advance leave, Supervisor review and Assistant replacement. Preserve the original schedule, exclude approved unpaid leave from absence and guard finalized payroll periods. HTTP/end-to-end wiring follows in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/catalog-contracts. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Require QR and GPS attendance with private evidence

Require assignment-bound signed QR plus GPS within 200m and accuracy <=50m. Add atomic check-in, evidence checkout, rerunnable absent finalization and private Storage access; reject replay per schedule and roll back evidence failures. Live QR/Storage deployment configuration remains required.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/scheduling-leave. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Add requisition survey, funding and purchasing

Add requisition survey, Supervisor review/funding and Assistant purchases. Preserve generated quantities, store actual unit prices and reject repeat item purchases; matching funding-reference retries are idempotent. End-to-end verification is included in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/attendance. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Calculate half-month payroll and record invoice receipts

Calculate half-month payroll with integer satang, exact penalty boundaries, capped deductions and net floor zero. Resolve recipient through worker.user_id, guard concurrent/overlapping payroll and record invoice receipts. End-to-end financial verification follows in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/procurement. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Allocate TOR costs and export profit reports

Allocate net Payroll by actual workdays per TOR, mark estimates and separate actual material expenses from fund transfers. Add profit JSON/CSV using paid invoice receipts; reject duplicate Worker mappings before aggregation. End-to-end accounting verification follows in feature/runtime-api.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/payroll. Do not merge out of dependency order. No schema changes, migrations or deployment.

## Wire all MVP routes and document the backend

Wire all implemented MVP services into the existing Gin Foundation and use verified LINE tokens on both Web and Mini App. Add graceful lifecycle/absence jobs, safe CORS/errors, request schemas and current/historical docs. Full isolated integration covers concurrency, rollback, leave/replacement, procurement, payroll, invoice/profit and ownership; real LINE/Supabase/Storage remain unverified.

Validation: go test -count=1 ./..., go vet ./..., go build ./... passed on this branch. Database tests used a disposable loopback PostgreSQL17, never the existing Supabase. Formatting and diff whitespace checked. Real service credentials/configuration were not used.

Stack: base feature/reports. Do not merge out of dependency order. No schema changes, migrations or deployment.
