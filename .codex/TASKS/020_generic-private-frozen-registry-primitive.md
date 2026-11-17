# Task 020: Generic private frozen-registry primitive

**Milestone:** 3, Composition and Registries  
**Primary predecessor:** 019  
**Expected package area:** application, component, internal/registry, internal/testapp

## Purpose

Deliver the smallest reviewable implementation increment for **Generic private frozen-registry primitive** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/00_FOUNDATION_ARCHITECTURE.md`
- `docs/architecture/01_CORE_RUNTIME.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `.codex/MILESTONES/03.md`
- `.codex/DECISION_AUTHORITY.md`

## Original settled task detail

Purpose:
- Implement reusable private registry machinery under `internal/registry`.

Requirements:
- deterministic registration
- duplicate rejection
- freeze prevents mutation
- private implementation only

Tests:
- duplicate registration
- freeze
- concurrent reads after freeze

Depends on: 008.

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
