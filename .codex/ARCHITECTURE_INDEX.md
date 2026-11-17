# Architecture Index

Treat these five files collectively as authoritative. Each file contains labeled embedded original source documents. Read relevant embedded sections rather than assuming topic solely from the consolidated filename.

- `docs/architecture/00_FOUNDATION_ARCHITECTURE.md`, charter, architecture overview, open decisions, repository/package/dependency plan, lifecycle, UDS, CLI, and cross-cutting ownership.
- `docs/architecture/01_CORE_RUNTIME.md`, components, actions, routes, rendering, widgets, authentication/authorization, logging, hooks/events, lifecycle behavior.
- `docs/architecture/02_DATA_EXECUTION.md`, MariaDB, migrations, transactions/cache, jobs, workflows, scheduling, conceptual database schema, migration contract.
- `docs/architecture/03_SERVICES_CONFIGURATION.md`, LLM/yllmd, localization, notifications, messaging, external integrations, and configuration schema.
- `docs/architecture/04_IMPLEMENTATION_CODEX.md`, public API skeleton, testing strategy, milestones, acceptance criteria, documentation outline, Codex operating spec, and task graphs.

When sources conflict, follow the source-precedence rule in `INSTRUCTIONS.md`. Wrapper-level superseding clarifications in the current consolidated files are newer than stale embedded wording they explicitly supersede.
