> **HISTORICAL — what the class-diagram package verified (diagram/export checks only, not application tests).**

# Verification

Passed:

- 16 domain classes exactly match named ER attributes, including visually intended CONTRACT_TOR.status.
- WORKER contains only the explicitly shown is_available attribute; no worker identifiers added.
- All 354 operation signatures have receiving-call evidence in the latest SDs.
- Actor labels and human receipt actions are excluded from software methods.
- No return values become methods; no SQL bodies are placed in class boxes.
- All multiline SQL labels with larger layout distances still map to DB Connector.
- All 64 logical classes and all their methods are covered by the split views.
- All 76 software dependencies and R01–R23 plus U1 are covered by the split views.
- Full draw.io has one complete definition per logical class.
- All editable connectors have valid source/target cells.
- No class boxes overlap; every connector route avoids the interior of class boxes in full and split layouts.
- Association IDs avoid class boxes and do not overlap one another.
- All 23 split PNGs decode correctly at 2000 × 2828 pixels.

Visual review: all split pages and the full overview inspected. Detailed domain and dense software/service pages inspected at larger size. Google Docs placement is checked against the intended width; no live Google Doc was created.

Source differences remain documented, especially worker identity, payroll terminology, deduction linkage, evidence context and expense linkage. These are unresolved reference decisions, not silently repaired source diagrams.

## Source fingerprints

- SA1-ER latest(ปรับทะแยง).drawio.svg: SHA-256 `dd486277235a5d2ddd6ccd88a0c3b017339d137943ad11e3db1c9c10f925101d`
- SA1-SD-1A.drawio.svg: SHA-256 `067519065676fb7ef2067c80978cfebde58f3881da3c2aa89c2b20c68a5365c9`
