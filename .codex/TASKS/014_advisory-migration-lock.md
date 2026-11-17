# Task 014: Advisory migration lock

**Milestone:** 2, MariaDB and Migrations  
**Primary predecessor:** 013  
**Expected package area:** mariadb, transaction, migrations

## Purpose

Deliver the smallest reviewable implementation increment for **Advisory migration lock** while preserving all settled yFramework architecture boundaries.

## Source documents

- `docs/architecture/02_DATA_EXECUTION.md`
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`
- `docs/architecture/03_SERVICES_CONFIGURATION.md`
- `.codex/MILESTONES/02.md`
- `.codex/DECISION_AUTHORITY.md`

## Original settled task detail

Purpose:
- Implement bounded MariaDB advisory lock around migration execution.

Requirements:
- Dedicated connection owns lock.
- Crash/connection loss releases lock.
- No persistent migration lock table.

Tests:
- two concurrent runners
- bounded acquisition failure
- connection-loss release

Depends on: 009.

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
