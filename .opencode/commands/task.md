---
description: Implement a specific .codex task by number (e.g. /task 009)
---
Implement task $ARGUMENTS from `.codex/TASKS/`.

1. Run `python3 scripts/next-task.py --task $ARGUMENTS` to confirm the file and
   its prerequisites.
2. Read the task file and the architecture sections it needs.
3. Implement it with tests, run `./scripts/verify.sh`, append the completion
   entry to `.codex/STATE/completed.md`, and commit as `task NNN: <title>`.

If a prerequisite is incomplete or an architecture decision is required, stop
and report instead of proceeding.
