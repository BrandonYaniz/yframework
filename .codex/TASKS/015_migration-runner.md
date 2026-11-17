# Task 015: Migration runner

**Milestone:** 2, MariaDB and Migrations  
**Primary predecessor:** 014  
**Expected package area:** mariadb, transaction, migrations

## Purpose

Deliver the smallest reviewable implementation increment for **Migration runner** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/02_DATA_EXECUTION.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `docs/architecture/03_SERVICES_CONFIGURATION.md`
- `.codex/MILESTONES/02.md`
- `.codex/DECISION_AUTHORITY.md`

## Original settled task detail

Purpose:
- Implement deterministic migration execution.

Requirements:
- framework -> application -> components
- component dependency order supplied later by composition layer
- apply only pending migrations
- checksum drift is fatal
- mark applied only after success
- no automatic down/rollback

Tests:
- empty DB
- pending migration
- checksum mismatch
- failed migration not marked applied

Depends on: 011-014.

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
