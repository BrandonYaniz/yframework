---
description: Execute one verified yFramework recovery task (R00–R12)
---

You are executing **one recovery iteration**, not advancing the normal 881-task graph.

1. Read `.codex/RECOVERY/00_START_HERE.md` and the next eligible R-task in `.codex/RECOVERY/01_REMEDIATION_TASKS.md`.
2. Read the exact referenced original `.codex/TASKS/NNN_*.md` requirements, relevant `docs/architecture/` subsections, `.codex/DECISION_AUTHORITY.md`, and `.codex/RECOVERY/02_VERIFICATION_GATES.md`.
3. Inspect the *current* source and write a verified reproduction/failing-before test. Make the smallest safe correction, with real assertions and non-destructive database isolation.
4. Run formatting, tests, vet, and the applicable **behavioral** integration gate. A skipped/missing/log-only required test is a blocker, not success.
5. Update `.codex/STATE/recovery-verification.md` with exact evidence. Preserve historical completion records. Commit a scoped corrective change unless existing Git policy or user instruction prevents it.
6. Report changed files, tests actually executed, commit, remaining blockers, and next task. Stop after this single task.

If $ARGUMENTS specifies an R-task (for example `R06`), execute that task only if prerequisites are met. Otherwise select the first unfinished eligible R-task. Do not silently pick a later item.
