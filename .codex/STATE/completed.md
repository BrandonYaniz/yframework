## Task 013 — Migration history schema correction

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/history.go`
- Tests: `go test ./migrations` (N tests)
- Decisions: Corrected migration history table name from `migration_history` to `framework_migrations`
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 014 — Advisory locking with dedicated connections

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/lock.go`
- Tests: `go test ./migrations` (N tests)
- Decisions: Implemented advisory locking with dedicated connections instead of DB pool to prevent connection pool issues
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 015 — Critical migration runner fixes

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `migrations/runner.go`, `migrations/integration_test.go`
- Tests: `go test ./migrations` (N tests)
- Decisions: Fixed migration runner to execute DDL operations before marking migrations as applied, and rewrote integration tests 
- Follow-ups: none
- Architecture: no unresolved decision silently made.

## Task 021 — Dependency graph topological sort fix

- Commit: 976adf910b1b48a4c11b850b9e3475e7410b4c0a
- Files: `internal/registry/dependency_graph.go`
- Tests: `go test ./internal/registry` (N tests)
- Decisions: Fixed topological sort algorithm to correctly order components with dependencies
- Follow-ups: none
- Architecture: no unresolved decision silently made.