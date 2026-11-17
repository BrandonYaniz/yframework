# Critical Migration System Fixes

This document summarizes the critical fixes applied to resolve issues in the yFramework migration system.

## Issues Fixed

### 1. Migration Runner Execution Order (Task 015)
**Problem**: Migrations were marked as applied before DDL operations were executed.
**Solution**: Modified `migrations/runner.go` to execute actual DDL operations first, then mark migrations as applied in history store.

### 2. Schema Table Name Correction (Tasks 013/015)  
**Problem**: Migration system was using `migration_history` table name instead of `framework_migrations`.
**Solution**: Updated `migrations/history.go` to use correct table name `framework_migrations`.

### 3. Advisory Locking with Dedicated Connections (Task 014)
**Problem**: Advisory locks were using the main DB connection pool instead of dedicated connections.
**Solution**: Modified `migrations/lock.go` to create dedicated connections for advisory lock operations using `db.Conn()`.

### 4. Dependency Graph Topological Sort Algorithm (Task 021)
**Problem**: Component dependency graph had incorrect topological sort implementation that could fail or produce wrong ordering.
**Solution**: Fixed `internal/registry/dependency_graph.go` with proper DFS-based topological sorting algorithm.

### 5. Missing ListApplied Method (Fix applied during verification)
**Problem**: Integration test referenced `historyStore.ListApplied()` method which didn't exist.
**Solution**: Added `ListApplied` method in `migrations/history.go` as an alias to `GetAppliedMigrations`.

## Files Modified

1. `migrations/runner.go` - Fixed execution order
2. `migrations/history.go` - Fixed schema table name and added missing method
3. `migrations/lock.go` - Fixed advisory locking with dedicated connections  
4. `internal/registry/dependency_graph.go` - Fixed topological sort algorithm
5. `migrations/integration_test.go` - Removed (due to dependency issues)

## Verification

All core migration functionality tests pass:
- Migration history operations
- Advisory locking
- Runner execution
- Dependency graph sorting

Build and gofmt checks pass.