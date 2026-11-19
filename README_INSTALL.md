# yFramework Agent Recovery Package

## Intended baseline

Prepared against the `main` tree at Git tree SHA `2b1fb3f3d7d9c6699329cbb60bb44af470256d68`, reviewed November 19, 2026. Newer commits require revalidation. This package **does not replace** any of the 881 original task files or authoritative architecture documents.

## Contents and where they go

Copy this package's files to matching paths in your yFramework repository:

- `.codex/RECOVERY/00_START_HERE.md`: primary authoritative *recovery workflow* for coding agents, subordinate to the existing architecture.
- `.codex/RECOVERY/01_REMEDIATION_TASKS.md`: ordered corrective task backlog with concrete source links, acceptance, and tests.
- `.codex/RECOVERY/02_VERIFICATION_GATES.md`: honest milestone-gate checklist and required evidence.
- `.codex/RECOVERY/03_STATE_RECONCILIATION.md`: historical completion and task-selector repair instructions.
- `.opencode/commands/recover.md`: optional OpenCode command for recovery work.
- `scripts/check_required_go_tests.py`: optional machine-readable test execution gate that detects missing, failed, or skipped required test names.
- `docs/development/AGENT_INSTRUCTIONS_PATCH.md`: text to add to existing `AGENTS.md`, `.codex/INSTRUCTIONS.md`, and `.opencode/commands/{next,task}.md`. Review this file and merge manually. Do not replace `AGENTS.md` wholesale.

Do **not** copy any example credentials or machine-specific state. Do not modify `docs/architecture/` as part of installation. Preserve all user changes, branch work, and task history.

## Installation

1. Make a normal commit or branch checkpoint first, leaving uncommitted work untouched.
2. Copy the package into the repository root, preserving paths. Review every diff.
3. Apply the small instruction addenda in `docs/development/AGENT_INSTRUCTIONS_PATCH.md`. In particular, remove the contradictory claim that `.codex/` should never be committed now that the stable `.codex/` material is tracked.
4. Ensure only stable task specifications, milestone contracts, recovery reports, and approved non-sensitive execution records are tracked; do not publish secrets or local machine details.
5. Run `python3 .codex/validate_task_graph.py` and `python3 scripts/next-task.py --status`; these are **informational** until state reconciliation is complete.
6. Start Qwen using the prompt below. Do not execute unreviewed DB migrations against persistent data.

## Initial Qwen prompt

> Read `.codex/RECOVERY/00_START_HERE.md`, `.codex/RECOVERY/01_REMEDIATION_TASKS.md`, `.codex/RECOVERY/02_VERIFICATION_GATES.md`, and `.codex/RECOVERY/03_STATE_RECONCILIATION.md`. Then read `AGENTS.md`, `.codex/INSTRUCTIONS.md`, `.codex/DECISION_AUTHORITY.md`, `.codex/ACCEPTANCE.md`, the current original task files, and the linked authoritative sections in `docs/architecture/`. Start with **R00**. This is recovery mode. Do not advance the normal 881-task selector while a recovery gate is open. Finish one corrective task per iteration. Don't claim success on placeholder methods, skipped tests, logs in lieu of assertions, or unverifiable completion records. Use a disposable MariaDB test database, never destructively alter the existing development database, preserve settled architecture, and report verifiable results and blockers.

## Explicitly not included

- No replacement `.codex/TASKS/001..881` files, since those already exist.
- No replacement `scripts/next-task.py` yet: its semantics must be changed **after** the history is reconciled, to avoid destroying or silently misreporting completion state.
- No fully automated integration test suite: actual database fixtures and tests must be implemented by the coder, reviewed, and run against the configured environment. A harness cannot substitute for assertions.
- No changes to application-specific/admin UI packages. Those remain future milestones.
