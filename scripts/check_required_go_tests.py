#!/usr/bin/env python3
"""Fail if required named Go tests were absent, failed, or skipped in `go test -json`.

Usage: go test -count=1 -json ./migrations > /tmp/migrations.json
       python3 scripts/check_required_go_tests.py /tmp/migrations.json \
           TestMigrationFromZero TestConcurrentMigration TestChecksumDrift

This is a small gate helper, NOT proof of test assertion quality or real DB use.
"""
from __future__ import annotations
import json
import pathlib
import sys


def main(argv: list[str]) -> int:
    if len(argv) < 3:
        print("usage: check_required_go_tests.py GO_TEST_JSON TEST_NAME [TEST_NAME ...]", file=sys.stderr)
        return 2
    path = pathlib.Path(argv[1])
    expected = set(argv[2:])
    seen: dict[str, set[str]] = {name: set() for name in expected}
    pkg_failed = False
    try:
        with path.open(encoding="utf-8") as f:
            for lineno, line in enumerate(f, 1):
                if not line.strip():
                    continue
                obj = json.loads(line)
                action = obj.get("Action")
                name = obj.get("Test")
                if name in seen and action in {"run", "pass", "fail", "skip"}:
                    seen[name].add(action)
                if obj.get("Action") == "fail" and not name:
                    pkg_failed = True
    except (OSError, json.JSONDecodeError) as exc:
        print(f"invalid/missing Go test report: {exc}", file=sys.stderr)
        return 2
    problems = []
    for name, actions in sorted(seen.items()):
        if "run" not in actions:
            problems.append(f"MISSING: {name} did not run")
        if "skip" in actions:
            problems.append(f"SKIPPED: {name}")
        if "fail" in actions:
            problems.append(f"FAILED: {name}")
        if "pass" not in actions:
            problems.append(f"NO PASS: {name}")
    if pkg_failed:
        problems.append("PACKAGE FAILURE")
    if problems:
        for problem in problems:
            print(problem, file=sys.stderr)
        print("GATE FAILED: named-test status is not evidence of behavior", file=sys.stderr)
        return 1
    print("Named tests executed and passed without skips. Review real assertions and DB effects separately.")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
