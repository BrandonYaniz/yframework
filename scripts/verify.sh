#!/usr/bin/env bash
# Verify the repository after a task.
#
# Usage:
#   ./scripts/verify.sh          # normal run
#   RACE=1 ./scripts/verify.sh   # also run the race detector
#
# Exit status is 0 only when every applicable check passes.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1

fail=0
note() { printf '\n== %s ==\n' "$1"; }
pass() { printf 'ok: %s\n' "$1"; }
skip() { printf 'skip: %s\n' "$1"; }
bad() { printf 'FAIL: %s\n' "$1"; fail=1; }

note "task graph"
if python3 .codex/validate_task_graph.py; then
  pass "task graph"
else
  bad "task graph validation"
fi

note "gofmt"
unformatted="$(gofmt -l . 2>/dev/null)"
if [ -n "$unformatted" ]; then
  printf '%s\n' "$unformatted" | sed 's/^/  /'
  bad "gofmt reports unformatted files"
else
  pass "gofmt"
fi

note "go packages"
if ! packages="$(go list ./... 2>/dev/null)" || [ -z "$packages" ]; then
  skip "no Go packages yet (repository bootstrap)"
else
  printf '%s\n' "$packages" | sed 's/^/  /'

  note "go build"
  if go build ./...; then pass "go build"; else bad "go build"; fi

  note "go vet"
  if go vet ./...; then pass "go vet"; else bad "go vet"; fi

  note "go test"
  if [ "${RACE:-0}" = "1" ]; then
    if go test -race -count=1 ./...; then pass "go test -race"; else bad "go test -race"; fi
  else
    if go test -count=1 ./...; then pass "go test"; else bad "go test"; fi
  fi

  note "staticcheck"
  if command -v staticcheck >/dev/null 2>&1; then
    if staticcheck ./...; then pass "staticcheck"; else bad "staticcheck"; fi
  else
    skip "staticcheck not installed"
  fi
fi

if [ "$fail" -ne 0 ]; then
  printf '\nVERIFY FAILED\n'
  exit 1
fi

printf '\nVERIFY PASSED\n'
