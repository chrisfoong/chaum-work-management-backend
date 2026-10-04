# docs/ index — read this first

Everything an agent needs to work on this backend is inside this repo. **There are no other references.** If something you need is not in `CLAUDE.md` or `docs/`, say so and ask — do not invent it and do not claim to have read files that are not here.

| File | What it gives you |
|---|---|
| `TASKS.md` | **Project state**: phases, task checkboxes, blockers, decisions log, files changed (uncommitted). Read first, update last |
| `../CLAUDE.md` | Rules, adopted decisions, canonical names, quick map. Always loaded |
| `01_project_overview.md` | Business scope, roles, platforms, flow, notification types |
| `02_domain_model.md` | The 16 entities with attributes, all relationships, value sets, and overrides of the historical class diagram |
| `03_workflows_operations.md` | For each sequence diagram (1S–6S, 1A–9A, 1W–7W): participants and received operations (names from the historical snapshot — apply overrides) |
| `04_decisions_and_corrections.md` | Adopted decisions, adopted diagram corrections, and open decisions D1–D14 |
| `05_business_rules.md` | Formulas, status behavior, validations, per-workflow behavior |
| `06_architecture_and_conventions.md` | Proposed layering, module map, coding conventions |
| `07_verification_checklist.md` | What to verify before work, acceptance checks, definition of done, report format |
| `prompts/CLAUDE_Prompts_Backend.md` | Ordered prompts for Claude in VS Code |
| `references/*` | Full copies of the source material (marked CURRENT or HISTORICAL) |

## Precedence when files disagree
1. The user's current instruction.
2. `CLAUDE.md` + `docs/02`, `docs/04` (production ER names, adopted corrections and decisions).
3. `docs/05` business rules.
4. The real schema/migrations in the repo for physical types and constraints (report differences with the ER; do not migrate silently).
5. `docs/03` operation names (historical snapshot).
6. `docs/references/*_HISTORICAL.md` (provenance only).

## Material that is NOT available (and must not be fabricated)
- The production ER drawing itself (its content is transcribed in `docs/02`).
- The use-case description tables and the exact SQL of each sequence diagram (SQL was abbreviated with `...` in the SDs). `docs/05` states the rules the team described; anything beyond it must be confirmed.
- The corrected SD files and their `Changes_and_SQL.md` (their effects are listed in `docs/04` §1.3).
- The database schema, unless it is already in the repo.
- Sequence/class diagram images.

When a task depends on missing material, stop and ask the user (use the stop-and-ask template in the prompts file).
