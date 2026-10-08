# Backend acceptance and live configuration

Backend only. Frontend belongs to a separate repository. Implemented use-case routes are mapped in [FRONTEND_HANDOFF.md](FRONTEND_HANDOFF.md). Do not execute historical migrations or change the shared Supabase schema. Agents must not inspect or modify .env.

## Evidence boundaries

- Local tests use an explicitly isolated loopback PostgreSQL database. LINE and Storage HTTP responses are simulated where credentials would otherwise be needed.
- The [shared verification conversation](https://chatgpt.com/s/cx_6ac7abfc16a881918c18b7b7b8418d98) reports Supervisor LINE login and `/api/web/me` returning 200 with an active Supervisor from Supabase. This is historical reported live evidence, not a live recheck by this task.
- The same conversation reports Worker Endpoint/openid/mapping and QR/Storage/Messaging configuration gaps. Those configuration values have not been inspected in this task, so their current state remains unverified.
- A successful HTTP 200 is insufficient unless the body is the expected backend JSON identity. A tunnel warning page does not prove login.

## Configuration owned by the operator

| Feature | Required configuration / data | Acceptance |
| --- | --- | --- |
| Database | DATABASE_URL for the existing Supabase PostgreSQL database | `/health/ready` returns 200; startup read-only schema check passes |
| Web login | LINE_WEB_CHANNEL_ID, LINE_PROVIDER_ID; Web LIFF openid and correct HTTPS Endpoint | `/api/web/me` returns active Supervisor/Assistant matching the intended test account |
| Worker login | LINE_WORKER_CHANNEL_ID, same verified Provider; Worker MINI App openid/HTTPS Endpoint; active USER and exactly one WORKER linked by user_id | `/api/liff/me` returns the intended Worker |
| QR + GPS | Operator-generated QR_SIGNING_SECRET of at least 32 bytes; actual assignment/location coordinates | Check-in accepts both valid QR and GPS <=200m; rejects either missing/invalid, expired/wrong-area QR and inaccurate GPS |
| Private files | SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY, STORAGE_BUCKET for an existing private bucket | Owned PNG/JPEG uploads and authorized reads succeed; wrong-owner reads fail; TOR requires PNG <=5MiB |
| LINE notifications | LINE_MESSAGING_TOKEN for Messaging API under the same Provider, recipient account able to receive bot messages | Actual recipient receives procurement decision/delivery message; purchase alone does not notify delivery |
| Browser access | WEB_ALLOWED_ORIGINS exact frontend HTTPS origins when calling backend across origins | Authorized preflight succeeds, unexpected origin denied |

No secret belongs in API requests, source code, screenshots or test evidence. URL/publishable key alone do not configure this backend's PostgreSQL connection. The operator manages configuration privately; ordinary DB row changes are read on the next request, while process configuration changes require a server restart.

## Login diagnosis

1. No bearer -> 401 `missing bearer token`.
2. Invalid/expired token, wrong channel or verification failure -> 401 `invalid or expired token`; USER lookup has not succeeded.
3. Verified subject with no active matching USER -> 401 `user is not registered or is inactive`.
4. Valid identity with wrong DB role -> 403.
5. Worker identity missing WORKER mapping -> 404; more than one mapping -> 409.
6. Success -> 200 backend JSON with correct user_id, role and is_active.

Changing the test-page role selector only selects an API. It does not change the token audience or DB role. USER.line_id stores the verified LINE subject, not a LIFF ID, Channel ID, searchable LINE username or access token. Worker/User primary keys remain distinct.

## Live acceptance still pending

- [ ] Assistant login and Supervisor-only route rejection with a real LINE token.
- [ ] Worker login using Worker channel, correct mapping and ownership restrictions.
- [ ] Read-only dashboard/contracts/schedules against the intended Supabase project.
- [ ] Storage byte/MIME/ownership behavior against the configured private bucket.
- [ ] Real LINE delivery to the intended recipient.
- [ ] QR+GPS check-in, checkout evidence, leave/replacement, purchasing/delivery, payroll and reports across a complete controlled workflow.

Run workflows that write data only on an approved isolated environment; local automated tests must never point at the shared Supabase project. No production schema or data edits are authorized by this checklist.

## Continuous integration

[backend.yml](../.github/workflows/backend.yml) runs formatting, tests, vet and build on push/PR. PostgreSQL17 is a disposable service reachable through loopback; only that service receives isolated-test-schema.sql. No repository environment secrets, Supabase credentials, deployments or migrations of the live database are used. Skipped tests fail CI so integration cannot silently be counted as passing. Coverage/results are retained for seven days. The workflow has to run on GitHub before remote CI is considered verified.

Schema limitations remain documented in REQUIREMENTS.md: durable outbox, historical wage snapshots, dedicated delivery/closing ledger and item-round history are not implemented persistence facilities. Passing tests does not remove these limitations.

## 9A acceptance

Confirmed scope excludes persisted summary/send/read history and acknowledgement. Isolated tests cover empty dates, current attendance, retrospective changes, no financial fields, contract continuation and schedule rejection, role/review/payroll/URL guards, absence of automatic Assistant summaries, and the manual link message using simulated LINE. Operator must configure the real HTTPS Assistant Dashboard URL and LINE transport privately. Successful real manual delivery and opening the authenticated Dashboard still require controlled live acceptance; automated tests do not prove them.
