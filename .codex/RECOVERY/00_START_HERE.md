# yFramework Recovery Mode: Agent Operating Contract

**State: ACTIVE until Milestones 1, 2, and 3 have independently passed their behavioral acceptance gates.**

## Goal

Recover a truthful implementation baseline for the existing yFramework repository; repair incomplete and unsafe foundational code; prove the specified behavior through tests; then resume the existing 881-task roadmap to operational v1. Recovery is a **temporary, ordered overlay**, not a replacement architecture or competing task graph.

## Instruction order

1. `docs/architecture/` consolidated authoritative source documents. Read relevant sections, including newer wrapper clarifications.
2. `.codex/DECISION_AUTHORITY.md`, `.codex/ACCEPTANCE.md`, `.codex/INSTRUCTIONS.md`, and milestone/task specifications.
3. `AGENTS.md` for local execution and Git constraints.
4. This recovery overlay for corrective scheduling and evidence, only where consistent with the above.

If there is a conflict or unresolved public, persisted, security, lifecycle, or compatibility decision, record it in `.codex/STATE/open-decisions.md` with options and recommendation. Do not invent an answer. A small internal implementation choice that preserves contracts may be made and documented.

## Nonnegotiable architecture

- Module `github.com/BrandonYaniz/yframework`; applications own `main`. No executable framework template, global service locator, or untyped dependency bag.
- MariaDB only, with FreeBSD-reference MariaDB release policy; separate foreground/background pools, explicit startup connection, forward-only migrations, and version/capability validation.
- Explicit typed component dependencies and registries, deterministic dependency-first composition, freeze before startup mutation, and component ownership of domain data.
- Go rendering, route, authentication, job, workflow, scheduler, LLM, localization, and messaging architecture remain governed by their existing planned milestones.
- Passkey-only human authentication v1; no passwords, fallback MFA, or account recovery. Never weaken existing security boundaries to make tests easier.
- Framework owns coordination and outside-system integrations. Components do not independently contact external transports or run worker pools.
- No executable configuration stored in MariaDB. Stable registry names are compatibility sensitive.
- FreeBSD, Linux, and macOS support; Unix socket HTTP routing. Don't hardcode platform paths/service managers.

## Recovery iteration, one task at a time

1. Run `git status --short`, `git rev-parse HEAD`, and inspect the exact current source. Protect uncommitted work.
2. Read `01_REMEDIATION_TASKS.md` and select the next **R** task whose prerequisites are satisfied. Ignore stale progress numbers from `.codex/STATE/completed.md` until R01 has reconciled them.
3. Read exact original `.codex/TASKS/NNN_*.md` and relevant architecture sections (use `rg -n` then bounded `sed`, do not flood a 64K context window).
4. Write failing-before tests or equivalent verifiable reproduction. **Never add a production no-op to make tests pass.**
5. Implement the smallest correct fix, including actual error paths, crash/recovery, idempotence, observability, and security impacts.
6. Run relevant unit tests and `./scripts/verify.sh`; when database-sensitive, run real MariaDB acceptance tests on an isolated disposable DB. Check **assertions**, not just exit status.
7. Update `.codex/STATE/recovery-verification.md` (create if missing) with original task IDs, initial behavior, exact commands, tests run/not run, source assertions, commit hash, remaining risks, and status.
8. Commit an appropriately scoped change as `fix(task NNN): ...` or `recovery RNN: ...` per existing Git policy. Do not push or alter history automatically.
9. Stop after one R task, report evidence and next eligible R task. The next iteration resumes from files on disk.

## Reject false completion

A `t.Skip`, a `t.Log` without assertions, `return nil` in an unimplemented method, a test that treats unexpected errors as success, or an integration test that never opens a real DB **does not satisfy** an acceptance requirement. Ordinary fast tests may skip *optional* integration cases, but a milestone gate must fail closed if the necessary environment or required tests are unavailable. Missing dependency or known incomplete task means **blocked/incomplete**, never verified.

Do not retroactively turn historical `completed.md` entries into invented proof. Retain history, create separate verified-status records, and repair the selector after recording factual evidence.

## Safety

- `devframework` is not disposable. Do not DROP, TRUNCATE, reset, or implicitly migrate it during tests. Require an explicitly named isolated test DB or schema with narrow permissions; refuse to run a destructive fixture against a protected name.
- Don't print `.env.local`, passwords, DSNs, passkey material, API tokens, or secret test fixtures. Review tracked `.codex/STATE/local-environment.md` for public exposure.
- MariaDB DDL can implicitly commit. Never assume a transaction can roll back arbitrary migration DDL. Document partial-state recovery.
- Do not edit the authoritative architecture files or silently choose new public APIs.

## Resume criteria

Recovery ends only when Milestone 1, 2, and 3 acceptance requirements are demonstrated and documented, the task selector no longer treats stale headings as proof, and the next original task's prerequisites are independently verified. Then continue original milestones in dependency order. `yframework-admin` is a separate future module; do not start it during recovery.
