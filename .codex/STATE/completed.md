## Task 001 — Repository bootstrap

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `go.mod`, `errors/`, `validation/`, `logging/`
- Tests: `go test ./...` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 002 — Error model foundation

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `errors/errors.go`, `errors/errors_test.go`
- Tests: `go test ./errors` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 003 — Validation foundation

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `validation/validation.go`, `validation/validation_test.go`
- Tests: `go test ./validation` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 004 — Bootstrap logger

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `logging/bootstrap.go`, `logging/bootstrap_test.go`
- Tests: `go test ./logging` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 005 — Structured logging core

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `logging/logging.go`, `logging/logging_test.go`
- Tests: `go test ./logging` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 006 — Typed TOML configuration parser

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `config/config.go`, `config/config_test.go`
- Tests: `go test ./config` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 007 — Configuration cross-field validation

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `config/config.go`, `config/config_test.go`
- Tests: `go test ./config` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 008 — Foundation public API review

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `errors/`, `validation/`, `logging/`, `config/`
- Tests: `go test ./...` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 009 — MariaDB connection/config integration

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/config.go`, `mariadb/config_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 010 — MariaDB version capability validation

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/config.go`, `mariadb/config_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 011 — Transaction helper foundation

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/transaction.go`, `mariadb/transaction_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 012 — Migration registry

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/registry.go`, `migrations/registry_test.go`
- Tests: `go test ./migrations` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 013 — Migration history schema

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/history.go`, `mariadb/history_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 014 — Advisory migration lock

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/lock.go`, `mariadb/lock_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 015 — Migration runner execution order

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/runner.go`, `migrations/runner_test.go`
- Tests: `go test ./migrations` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 016 — Runtime MariaDB pools

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/pools.go`, `mariadb/pools_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 017 — Migration CLI service API primitives

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/cli.go`, `migrations/cli_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 018 — Milestone 2 MariaDB integration suite

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `mariadb/integration_test.go`
- Tests: `go test ./mariadb` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 019 — Stable component identity types

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `application/component.go`
- Tests: `go test ./application` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 020 — Generic private frozen-registry primitive

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `internal/registry/frozen_registry.go`, `internal/registry/frozen_registry_test.go`
- Tests: `go test ./internal/registry` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 021 — Component dependency graph topological sort

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `internal/registry/dependency_graph.go`, `internal/registry/dependency_graph_test.go`
- Tests: `go test ./internal/registry` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 022 — Component registrar skeleton

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `internal/registry/registrar.go`, `internal/registry/registrar_test.go`
- Tests: `go test ./internal/registry` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 023 — Application builder skeleton

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `application/builder.go`, `application/builder_test.go`
- Tests: `go test ./application` (N tests)
- Decisions: none
- Follow-ups: none
- Architecture: no unresolved decision silently made.