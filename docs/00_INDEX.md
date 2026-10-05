# docs/ index — read this first

Everything an agent needs to work on this backend is inside this repo. **There are no other references.** If something you need is not in `CLAUDE.md` or `docs/`, say so and ask — do not invent it and do not claim to have read files that are not here.

**`docs/references/` is git-ignored and exists only on the maintainer's machine; it may be absent on a fresh clone.** `docs/00`–`07` (and `TASKS.md`) are the tracked summary of that material. If a reference file you need is missing, say so and ask the user; never guess its content.

| File | What it gives you |
|---|---|
| `TASKS.md` | **Project state**: phases, task checkboxes, blockers, decisions log, files changed (uncommitted). Read first, update last |
| `../CLAUDE.md` | Rules, adopted decisions, canonical names, quick map. Always loaded |
| `01_project_overview.md` | Business scope, roles, platforms, flow, notification types |
| `02_domain_model.md` | The 16 entities with attributes, all relationships, value sets, and overrides of the historical class diagram |
| `03_workflows_operations.md` | For each sequence diagram (1S–6S, 1A–9A, 1W–7W): participants and received operations (names from the historical snapshot — apply overrides) |
| `04_decisions_and_corrections.md` | Adopted decisions, adopted diagram corrections, decided and open decisions (D1–D18), source conflicts (C#), notes for the frontend repo |
| `05_business_rules.md` | Formulas, status behavior, validations, per-workflow behavior |
| `06_architecture_and_conventions.md` | Proposed layering, module map, coding conventions |
| `07_verification_checklist.md` | What to verify before work, acceptance checks, definition of done, report format |
| `prompts/CLAUDE_Prompts_Backend.md` | Ordered prompts for Claude in VS Code |
| `references/*` | **Local-only (git-ignored; may be absent).** Full copies of the source material (marked CURRENT or HISTORICAL) |
| `references/Chaum-Diagrams-PlantUML.txt` | CURRENT team diagrams (55, PlantUML): class, collaboration, component, data-flow, state. Authority for operation names, call order and who calls whom; not for column names/types (datadict) or rules (use cases) |
| `references/datadict.txt` | CURRENT team data dictionary (16 tables): column names, types, constraints, sample rows |
| `references/uc.txt` | CURRENT use-case descriptions with SQL for 1S–6S, 1A–9A, 1W–7W |
| `references/er-latest.png` | CURRENT ER drawing (image) |

## Precedence when files disagree
1. The user's current instruction.
2. For operation names, call order and caller/callee: `references/Chaum-Diagrams-PlantUML.txt` (over docs/03 and the historical mapping). Where it still uses names replaced by docs/04 §1.3, report it; do not choose.
3. For column names, types and constraints: `references/datadict.txt` plus the user decisions in docs/04; for workflow behavior and rules: `references/uc.txt`. Where these two disagree, list it in docs/04 §2.2; do not choose.
4. `CLAUDE.md` + `docs/02`, `docs/04` (production ER names, adopted corrections and decisions).
5. `docs/05` business rules.
6. The real schema/migrations in the repo for physical types and constraints (report differences with the ER; do not migrate silently).
7. `docs/03` operation names (historical snapshot).
8. `docs/references/*_HISTORICAL.md` (provenance only).

## Material that is NOT available (and must not be fabricated)
- (Now in `docs/references/`, local-only: the ER drawing `er-latest.png`, the use-case descriptions with SQL `uc.txt`, the data dictionary `datadict.txt`, the class/collaboration diagrams. If absent on this machine, ask.)
- The corrected SD files and their `Changes_and_SQL.md` (their effects are listed in `docs/04` §1.3).
- The team's real database DDL. `migrations/0001_init.*.sql` is Claude-designed from the data dictionary and has never been run.
- Use case 7S (Manage Location) and a normal (non-substitute) check-in use case (see D15).
- Sample rows for 11 of the 16 data-dictionary tables (empty in the file).
- Collaboration diagrams for alternate paths (only normal cases exist; e.g. 3S reject, 6A no_purchase).
- Sequence diagram images (the PlantUML file states they exist only as images). Class, collaboration, component, data-flow and state diagrams are available as PlantUML.

When a task depends on missing material, stop and ask the user (use the stop-and-ask template in the prompts file).
