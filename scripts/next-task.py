#!/usr/bin/env python3
"""Report or select the next executable task from the .codex task graph.

Exit codes:
  0  a task is available (or --status/--task succeeded)
  1  every task is complete
  2  no task is eligible because prerequisites are missing/unknown

Usage:
  python3 scripts/next-task.py                 # describe the next task
  python3 scripts/next-task.py --path          # print only its path
  python3 scripts/next-task.py --status        # progress by milestone
  python3 scripts/next-task.py --task 019      # describe one task
  python3 scripts/next-task.py --json          # machine-readable
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
CODEX = ROOT / ".codex"
TASKS_DIR = CODEX / "TASKS"
MILESTONES_DIR = CODEX / "MILESTONES"
COMPLETED = CODEX / "STATE" / "completed.md"

TASK_FILE_RE = re.compile(r"^(\d{3})_(.+)\.md$")
TITLE_RE = re.compile(r"^# Task \d{3}:\s*(.+?)\s*$", re.MULTILINE)
MILESTONE_RE = re.compile(r"^\*\*Milestone:\*\*\s*(\d+)\s*,?\s*(.*?)\s*$", re.MULTILINE)
PREDECESSOR_RE = re.compile(r"^\*\*Primary predecessor:\*\*+\s*(\S+)", re.MULTILINE)
DEPENDS_RE = re.compile(r"^Depends on:\s*(.+?)\s*$", re.MULTILINE)
COMPLETED_RE = re.compile(r"^##\s+Task\s+(\d{3})\b", re.MULTILINE)


def milestone_titles() -> dict[int, str]:
    titles: dict[int, str] = {}
    if not MILESTONES_DIR.is_dir():
        return titles
    for path in sorted(MILESTONES_DIR.glob("*.md")):
        first = path.read_text(encoding="utf-8").splitlines()[0:1]
        if not first:
            continue
        match = re.match(r"#\s*Milestone\s+(\d+):\s*(.+)", first[0])
        if match:
            titles[int(match.group(1))] = match.group(2).strip()
    return titles


def expand_dependency(token: str) -> list[str]:
    token = token.strip().strip(".,;")
    if not token or token.lower() in {"none", "n/a", "-"}:
        return []
    range_match = re.fullmatch(r"(\d{1,3})\s*[-–]\s*(\d{1,3})", token)
    if range_match:
        start, end = int(range_match.group(1)), int(range_match.group(2))
        if start > end:
            start, end = end, start
        return [f"{n:03d}" for n in range(start, end + 1)]
    if re.fullmatch(r"\d{1,3}", token):
        return [f"{int(token):03d}"]
    return []


def parse_dependencies(text: str) -> list[str]:
    """Union of the formal 'Depends on' line and the primary predecessor."""
    deps: list[str] = []
    depends = DEPENDS_RE.search(text)
    if depends:
        for token in depends.group(1).split(","):
            deps.extend(expand_dependency(token))
    predecessor = PREDECESSOR_RE.search(text)
    if predecessor:
        deps.extend(expand_dependency(predecessor.group(1)))
    # De-duplicate while preserving numeric order.
    return sorted(set(deps))


def load_tasks() -> dict[str, dict]:
    tasks: dict[str, dict] = {}
    if not TASKS_DIR.is_dir():
        return tasks
    for path in sorted(TASKS_DIR.iterdir()):
        match = TASK_FILE_RE.match(path.name)
        if not match:
            continue
        number = match.group(1)
        text = path.read_text(encoding="utf-8")
        title_match = TITLE_RE.search(text)
        milestone_match = MILESTONE_RE.search(text)
        tasks[number] = {
            "number": number,
            "title": title_match.group(1) if title_match else match.group(2).replace("-", " "),
            "path": path.relative_to(ROOT).as_posix(),
            "milestone": int(milestone_match.group(1)) if milestone_match else 0,
            "dependencies": parse_dependencies(text),
        }
    return tasks


def load_completed() -> set[str]:
    if not COMPLETED.is_file():
        return set()
    return set(COMPLETED_RE.findall(COMPLETED.read_text(encoding="utf-8")))


def dependency_state(task: dict, completed: set[str]) -> tuple[list[str], list[str]]:
    done = [d for d in task["dependencies"] if d in completed]
    missing = [d for d in task["dependencies"] if d not in completed]
    return done, missing


def choose_next(tasks: dict[str, dict], completed: set[str]) -> tuple[dict | None, list[dict]]:
    incomplete = [t for number, t in sorted(tasks.items()) if number not in completed]
    eligible = [t for t in incomplete if not dependency_state(t, completed)[1]]
    if eligible:
        return eligible[0], incomplete
    return None, incomplete


def describe(task: dict, completed: set[str], titles: dict[int, str]) -> str:
    done, missing = dependency_state(task, completed)
    milestone = task["milestone"]
    milestone_label = f"{milestone} — {titles[milestone]}" if milestone in titles else str(milestone)
    if task["dependencies"]:
        deps = ", ".join(
            f"{d} ({'done' if d in done else 'missing'})" for d in task["dependencies"]
        )
    else:
        deps = "none"
    status = "complete" if task["number"] in completed else "incomplete"
    return (
        f"Task {task['number']} — {task['title']}\n"
        f"  milestone: {milestone_label}\n"
        f"  file:      {task['path']}\n"
        f"  depends:   {deps}\n"
        f"  status:    {status}"
    )


def status(tasks: dict[str, dict], completed: set[str], titles: dict[int, str]) -> str:
    lines = ["yFramework task graph", ""]
    by_milestone: dict[int, list[str]] = {}
    for number, task in tasks.items():
        by_milestone.setdefault(task["milestone"], []).append(number)
    for milestone in sorted(by_milestone):
        numbers = by_milestone[milestone]
        done = sum(1 for number in numbers if number in completed)
        label = titles.get(milestone, "Unknown milestone")
        lines.append(f"  {milestone:>2}  {label:<42} {done:>3}/{len(numbers):<3}")
    lines.append("")
    lines.append(f"  Total: {len(completed)}/{len(tasks)} tasks complete")
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Report or select the next .codex task.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    parser.add_argument("--path", "--print", dest="path_only", action="store_true",
                        help="print only the selected task file path")
    parser.add_argument("--status", action="store_true", help="show progress by milestone")
    parser.add_argument("--task", metavar="NNN", help="describe a specific task")
    parser.add_argument("--json", action="store_true", help="emit JSON")
    args = parser.parse_args()

    tasks = load_tasks()
    if not tasks:
        print(f"error: no task files found under {TASKS_DIR}", file=sys.stderr)
        return 2

    completed = load_completed()
    titles = milestone_titles()

    if args.task:
        number = f"{int(args.task):03d}"
        task = tasks.get(number)
        if task is None:
            print(f"error: task {number} not found", file=sys.stderr)
            return 2
        _, missing = dependency_state(task, completed)
        if args.json:
            print(json.dumps({**task, "missing_dependencies": missing,
                              "complete": number in completed}, indent=2))
        else:
            print(describe(task, completed, titles))
        return 2 if missing else 0

    if args.status:
        print(status(tasks, completed, titles))
        next_task, _ = choose_next(tasks, completed)
        if next_task:
            print("")
            print(f"Next: {next_task['number']} — {next_task['title']}")
        return 0

    next_task, incomplete = choose_next(tasks, completed)

    if next_task is None:
        if not incomplete:
            if args.json:
                print(json.dumps({"complete": True, "total": len(tasks)}))
            else:
                print(f"All {len(tasks)} tasks are complete.")
            return 1
        blocked = []
        for task in incomplete[:20]:
            _, missing = dependency_state(task, completed)
            blocked.append({"number": task["number"], "title": task["title"], "missing": missing})
        if args.json:
            print(json.dumps({"complete": False, "blocked": blocked}, indent=2))
        else:
            print("No eligible task: every incomplete task has an incomplete prerequisite.")
            for item in blocked:
                print(f"  {item['number']} needs {', '.join(item['missing'])}")
        return 2

    if args.path_only:
        print(next_task["path"])
        return 0

    if args.json:
        done, missing = dependency_state(next_task, completed)
        print(json.dumps({**next_task, "complete": False,
                          "missing_dependencies": missing,
                          "completed_count": len(completed),
                          "total": len(tasks)}, indent=2))
        return 0

    print(status(tasks, completed, titles))
    print("")
    print(describe(next_task, completed, titles))
    return 0


if __name__ == "__main__":
    sys.exit(main())
