# yFramework — Agent Operating Manual

This repository is **yFramework**, an importable Go framework. Applications own
`main`; this module never ships one. MariaDB is the only relational database.

The authoritative design is in `docs/architecture/` and is read-only. `.codex/`
is the execution control surface: task graph, milestone contracts, and state.

## The unit of work is one task

`.codex/TASKS/` contains 881 numbered tasks. Do **one task per iteration**.
Never batch tasks, and never start a task whose prerequisites are incomplete.

Start every iteration with:

```bash
python3 scripts/next-task.py
```

Then follow this loop:

1. **Read the task file** and the "Source documents" it lists.
2. **Read only the relevant architecture sections** — see the context budget.
3. **Plan** the smallest change that satisfies the task.
4. **Implement** code and tests together.
5. **Verify** with `./scripts/verify.sh`. Fix failures before continuing.
6. **Log** the task in `.codex/STATE/completed.md`.
7. **Commit** with `task NNN: <title>`.

## Context budget

This project runs on a local 64k-context model. The five architecture documents
total ~320 KB (~80k tokens), so reading them whole is impossible and unnecessary.

- Find a section with `rg -n "^#{1,3} " docs/architecture/<file>.md`.
- Read only that section with `sed -n 'START,ENDp'`, not the whole file.
- Prefer the task file, the milestone contract (`.codex/MILESTONES/NN.md`), and
  `docs/architecture/04_IMPLEMENTATION_CODEX.md`.
- Open `00_FOUNDATION_ARCHITECTURE.md` for module layout and boundaries.
- Open `01_CORE_RUNTIME.md`, `02_DATA_EXECUTION.md`, or
  `03_SERVICES_CONFIGURATION.md` only for the subsystem the task touches.
- Never dump an entire architecture file into the conversation.
- If a required fact is missing, say so; do not guess.

## Definition of done

A task is done only when **all** of these hold:

- The task's requirements and acceptance criteria are met.
- New behavior has tests; changed behavior updates existing tests.
- `./scripts/verify.sh` passes: task graph, gofmt, build, vet, tests, staticcheck.
- No forbidden architectural change was made.
- `.codex/STATE/completed.md` has a new entry.
- The work is committed.

## Completion log format

Append to `.codex/STATE/completed.md`, newest last:

```markdown
## Task 001 — Repository bootstrap

- Commit: <sha or "uncommitted">
- Files: `go.mod`, `errors/`, ...
- Tests: `go test ./...` (N tests)
- Decisions: <private implementation choices, or "none">
- Follow-ups: <discovered work, or "none">
- Architecture: no unresolved decision silently made.
```

Update `.codex/CURRENT_MILESTONE.md` only when a milestone gate task passes.

## Architecture boundaries — do not violate

- No `examples/` directory and no framework-owned `main` or executable template.
- MariaDB only; no multi-database abstraction.
- No service locator, generic dependency bag, or string-based service lookup.
- No independent component worker pool; no direct external transport from
  components.
- No executable scripts or arbitrary expressions stored in MariaDB.
- Do not move application/component domain ownership into framework packages.
- Public API, persisted semantics, security boundaries, lifecycle behavior,
  compatibility promises, and failure/recovery semantics are settled. Changing
  them requires an architecture decision, not an implementation choice.

## When to stop and ask

`DECISION_AUTHORITY.md` defines what may be decided alone. If a choice would
alter a public API, persisted data, security boundary, ownership boundary,
lifecycle, compatibility, or failure semantics:

1. Do not decide it silently.
2. Append an entry to `.codex/STATE/open-decisions.md` (task, conflict, options,
   recommendation).
3. Stop and report the blocker to the user.

Private choices — helper types, internal representation, library choice where
delegated — are yours to make and record in the completion log.

## Git rules

- One commit per task: `task NNN: <short title>`.
- Never `git push`, `git reset --hard`, `git clean`, or rewrite history.
- Never commit `.codex/` or `.localc/` (gitignored by user choice).
- Do not commit changes to `docs/architecture/`; those documents are
  authoritative and read-only.

## Local development database

MariaDB 13.0.2 is installed and running on `127.0.0.1:3306` (Homebrew).

- Database: `devframework`
- User: `devframework`
- Password: `YFRAMEWORK_DB_PASSWORD` in `.env.local` (gitignored)
- Full details: `.codex/STATE/local-environment.md`

Rules:

- Load `.env.local` when a task needs the database. Never hardcode credentials in
  source, tests, docs, commits, or log output, and never print the password.
- The `mariadb` client is on `PATH`. Check connectivity with
  `mariadb -h 127.0.0.1 -u devframework -p devframework -e "SELECT 1"`.
- Do not drop, truncate, or reset the `devframework` database. Schema-mutating
  tests create and drop their own objects.
- Use this server from task 009 onward (connection config) and task 018 onward
  (migrations, locks, transactions). A separate `devframework_test` database is
  recommended before task 018 and is not yet created.

## Known repository facts

- Module path is `github.com/BrandonYaniz/yframework` (set correctly by task 001).
- Go 1.27.3 is installed. `staticcheck` is at `~/Development/Go/bin/staticcheck`.
- This repo lives in iCloud Drive. If the file provider stalls a build, pause
  iCloud sync for the directory rather than moving the repo.
- `.codex/TASKS/008` and `018` contain stray `# Milestone N` headings. They are
  known defects; ignore those headings.
