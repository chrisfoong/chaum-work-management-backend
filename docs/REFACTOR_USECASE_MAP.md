# Use Case module organization — 2026-11-08

This refactor organizes the existing implementation by the adopted Use Case descriptions. It does not introduce new Use Cases, change the schema, rewrite queries, or change authorization or business behavior. Public service methods and request/response types retain their names and signatures.

| Use Case | Implementation after refactor |
|---|---|
| 1S — area QR | `internal/work/area_qr.go`: permanent QR claims, signing, verification, issuance |
| 4S — attendance finalization | `internal/work/attendance_finalize.go`: `Finalize`, independent of payroll |
| 3W / 4W — attendance | `internal/work/attendance.go`: GPS, check-in, checkout, attendance reads |
| 5W — equipment request | `internal/work/procurement.go`: request models, creation, ownership-filtered list/detail |
| 1A — initial survey | `internal/work/procurement_survey.go`: `Survey` |
| 3S — procurement approval | `internal/work/procurement_approval.go`: `ReviewRequest` |
| 2S — funding evidence | `internal/work/procurement_funding.go`: `Fund` |
| 2A / 7A — purchase rounds | `internal/work/procurement_purchase.go`: models, preview and shared transactional commit |
| 5A / 6A — inspection and decision | `internal/work/procurement_inspection.go`: `InspectRequest`, `DecideRequest` |
| 5S / 7W — payroll and paid slips | `internal/work/payroll.go`: calculation, batch, confirmation and reads |
| 6S — project invoicing | `internal/work/invoice.go`: invoice models, validation and receipt confirmation |

The shared `work` package, HTTP routes, repository transactions, SQL references, enums, exact-money arithmetic, notification timing, GPS/QR guards, eight-hour shifts and role/ownership checks are unchanged. Invoice deduction/budget helpers and administrative edit workflows remain in their existing modules.

Verification for this change: 317 test/subtest passes, zero failures or skips on an isolated loopback PostgreSQL database; `go vet ./...` and `go build ./...` pass. A source AST comparison found the same 271 production declarations in `internal/work`, including function bodies and signatures. This is not a new real-device/LINE end-to-end test or a production deployment.
