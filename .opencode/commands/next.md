---
description: Pick up the next .codex task and complete one iteration
---
Start exactly one iteration of the yFramework task loop.

1. Run `python3 scripts/next-task.py` and read the task it reports.
2. Read the task file and only the architecture sections it needs.
3. Implement the task with tests.
4. Run `./scripts/verify.sh` and fix every failure.
5. Append the completion entry to `.codex/STATE/completed.md`.
6. Commit with `task NNN: <title>`.

Do not start a second task in this iteration. Stop and report when the task is
done, or when an architecture decision is required.

$ARGUMENTS
