# Confirmed requirements — 2026-10-08

See AGENTS.md for schema mapping, role rules and confirmed business policies. Supabase stays unchanged. Every role signs in with LINE; no Supabase Auth linkage is required. All assistant assignments are permitted. Supervisor reviews leave; assistant assigns replacements. Approved leave is unpaid without absent penalty. Work lasts 8 hours. Both GPS (200m, <=50m accuracy) and signed QR (60s) are mandatory. Payroll is half-monthly, penalty tiers 300/400/1500 THB with net capped at zero. Purchase each item once at the actual unit price. Revenue is paid net_received; payroll allocations use actual workdays per TOR and are estimates.

LINE channel IDs/provider, DB credentials, QR signing secret and Storage settings must come from the deployment environment. Enum labels need verification from metadata. Runtime schema preflight must reject incompatible labels/columns; never change the DB. No payment timestamps/reference, delivery tracking, permanent correction audit or durable outbox are claimed.

Worker record creation needs a per-user transaction lock because worker.user_id has no UNIQUE. Leave and payroll also need transactional guards shared by all backend writers. Direct DB writers can bypass those guards. Duplicate existing rows are conflicts. Generated to_buy_qty is read-only.
