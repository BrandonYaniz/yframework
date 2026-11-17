# Task 026: Milestone 3 public API review

**Milestone:** 3, Composition and Registries  
**Primary predecessor:** 025  
**Expected package area:** application, component, internal/registry, internal/testapp

## Purpose

Deliver the smallest reviewable implementation increment for **Milestone 3 public API review** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/00_FOUNDATION_ARCHITECTURE.md`
- `docs/architecture/01_CORE_RUNTIME.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `.codex/MILESTONES/03.md`
- `.codex/DECISION_AUTHORITY.md`

## Original settled task detail

Purpose:
- Review composition/component exported APIs before later subsystems build on them.

Acceptance:
- no service locator
- no circular package dependencies
- minimal exported symbols
- consumer fixture compiles
- `go test -race ./...` passes for implemented concurrent registry code

Depends on: 019-025.

## Forbidden architectural changes

- Do not add multi-database support. MariaDB remains the only relational database.
- Do not add a framework-owned application `main` or executable template.
- Do not introduce a service locator or generic dependency bag.
- Do not move application/component domain ownership into framework packages.
- Do not store executable scripts or arbitrary expressions in MariaDB.
- Do not change public, persisted, security, lifecycle, compatibility, or recovery semantics unless an authoritative source already settles the change.

## Completion notes

When completing this task, record:
- files changed,
- tests run,
- notable private implementation choices,
- follow-up tasks discovered,
- confirmation that no unresolved architectural decision was silently made.
