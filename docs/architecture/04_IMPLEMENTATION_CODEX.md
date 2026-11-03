# yFramework Consolidated Source: 04_IMPLEMENTATION_CODEX

This consolidated project source preserves the original documents listed below verbatim. The wrapper headings and checksum manifest are consolidation metadata only.

## Included source documents

- `19_PUBLIC_GO_API_SKELETON.md` — SHA-256 `8e5c19040d7d50dc56f423056cb0b499dd7edfd27ad65e0b9e5420f1b1a1bdea`
- `20_TESTING_STRATEGY.md` — SHA-256 `43a731cbdbe677586ac29fbdab6cd6f01608d081f38a749153d9d3ef025bf562`
- `21_IMPLEMENTATION_MILESTONES.md` — SHA-256 `b05d0692a64eeaa6967fcbd566a6d62a0992284a2388088fc87b92b011446a5d`
- `22_ACCEPTANCE_CRITERIA.md` — SHA-256 `f931e3158faa817fde7d127428d4eefd283da836f6d3d5f123ff51b004d01932`
- `23_DOCUMENTATION_OUTLINE.md` — SHA-256 `45ac763c02f70e72a1ae5aa56bec9006045a8664cfb5289d7bc5348cc7d65fab`
- `24_CODEX_OPERATING_SPEC.md` — SHA-256 `82c3d34a06eedff41eed88d7f00121b4b95172d42f6fc869ff0296db9a0affc8`
- `25_INITIAL_CODEX_TASK_GRAPH.md` — SHA-256 `c668b34d1c44a96abb46c08042a23d745b06ea6ef0bc16796b7f37f7e924684e`
- `26_FULL_CODEX_TASK_GRAPH_M04_M15.md` — SHA-256 `cbd79990a62cf1ef6b23d334ae4c6ccf48d625bd88278c34da0e0219128f250a`

---

<!-- ORIGINAL_SOURCE_BEGIN: 19_PUBLIC_GO_API_SKELETON.md -->

# Original Source File: `19_PUBLIC_GO_API_SKELETON.md`

# yFramework Public Go API Skeleton

## Goal

Freeze dependency direction and public contract categories without prematurely locking every method signature.

## Core Rule

Prefer narrow typed interfaces owned by consumers. Do not create a global service container, service locator, string-based dependency lookup, generic `Get(name)` API, or use `context.Context` as a service bag.

## Module

Initial public module:

`github.com/BrandonYaniz/yframework`

Applications live in separate repositories and consume tagged yFramework module releases explicitly.

## application

Owns application composition/building and runtime handoff. It does not own `main`.

Applications own process entry and call yFramework composition/lifecycle APIs.

`Application` should expose only lifecycle operations needed by the CLI/process boundary, not arbitrary subsystem retrieval.

## component

Reusable component contract includes a stable component ID, explicit dependency declaration, and declarative registration.

Component registration must be DB-free and deterministic.

## Registrar

The component registrar is composed from typed subsystem registrars, for example actions, permissions, routes, widgets, jobs, workflows, schedules, hooks/events, migrations, LLM purposes, localization, and notifications.

No untyped `Register(any)` or service lookup.

## Cross-Component Dependencies

Components define narrow typed ports. The application supplies implementations during composition.

Components do not call other components by registry name.

## action

Actions have strongly typed implementation boundaries. Dynamic route/input configuration maps into known typed inputs before execution.

Exact use of Go generics should be prototyped before signature freeze.

## Context

`context.Context` is used for cancellation, deadlines, and request/trace metadata. Dependencies and domain state are not hidden in context values.

## auth

Public auth APIs expose stable identities and safe metadata, not internal numeric IDs, token hashes, passkey public keys, or persistence rows.

Human and client/application identities remain distinguishable.

## authorization

Public APIs expose stable permission types, actors, and permission checking. Roles are administrative grouping only. Resource authorization stays typed and component-owned.

## mariadb / transaction

The framework owns pool construction and injection. Components do not independently open relational DB connections.

Use standard `database/sql` types/interfaces where practical rather than unnecessary abstraction layers.

No universal CRUD repository interface. Components define domain-specific repository ports.

## cache

Framework cache facilities are typed/narrow. Durable correctness never depends on cache contents.

## job

Typed job definitions/handlers and enqueue APIs are registered under stable names. Jobs remain durable, FIFO among eligible work, and framework-owned for leases/retries.

Exact heterogeneous typed registry API should be prototyped before signature freeze.

## workflow

Workflow public APIs encode prepare/commit/cleanup semantics and registered executable steps. Components own transitional data and final domain commit behavior.

Exact typed workflow API should be prototyped.

## schedule

Public schedule definitions cover once, cron, and interval plus missed/overlap policy. Schedules target registered jobs/workflows, not arbitrary Go functions/scripts.

## messaging

Components express normalized message intent. They never receive SMTP clients or transport credentials.

## notification

Components emit notification intent/events. Notification resolves subscriptions and creates messages; messaging owns delivery.

## external

The framework exposes normalized outbound transport/resilience primitives. Components define typed provider ports; application-selected adapters use the framework transport runtime.

Components never directly create external HTTP/SMTP/etc. clients outside the defined integration boundary.

## llm

Components use purpose-based LLM generation APIs. Framework owns durable request state and yllmd transport. Raw model-routing internals remain below the public component API unless a registered purpose explicitly permits options.

## render / route

Renderers and route definitions expose stable contracts while matcher/render internals stay private. Routes reference registered actions/renderers/permissions and typed declarative input mappings.

## events / hooks

Events/hooks are synchronous and typed. Durable or retryable follow-up work is enqueued rather than implemented through component goroutines.

## logging / health / config

Expose narrow structured logger, health reporter, and typed config/setting contracts. Subsystems report capability health; central lifecycle aggregation owns global readiness.

## testing

A public `testing` package provides component-author fakes and harnesses for public contracts without exposing private runtime internals.

## API Prototypes Required Before Exact Signature Freeze

- heterogeneous typed action registration,
- typed job registration/enqueue,
- workflow definitions and commit semantics,
- typed runtime settings,
- composed component registrar,
- action/render bridging.

Compile-only prototypes should prove ergonomics and avoid reflection-heavy or `any`-heavy APIs before public signatures become compatibility commitments.

<!-- ORIGINAL_SOURCE_END: 19_PUBLIC_GO_API_SKELETON.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 20_TESTING_STRATEGY.md -->

# Original Source File: `20_TESTING_STRATEGY.md`

# yFramework Testing Strategy

## Layers

1. Unit tests.
2. Package integration tests.
3. Real MariaDB integration tests.
4. Full framework runtime/lifecycle tests using a minimal internal consumer.
5. Public consumer compatibility tests.

## Real MariaDB Requirement

Database-sensitive behavior must be tested against real MariaDB, not mocks alone. This includes SKIP LOCKED, advisory locks, migrations, FKs, transaction behavior, lease races, and MariaDB temporal semantics.

Release validation should include the supported FreeBSD-ports MariaDB baseline. CI may use an equivalent MariaDB service/container where appropriate.

## Unit Coverage Areas

Config, validation, route matching/normalization, schedule calculations, permission resolution, parsing, and deterministic subsystem behavior should have fast table-driven unit tests.

## Migrations

Required tests include empty-to-current, prior-supported-to-current, checksum drift, failure recording, concurrent migration lock, crash lock release, unsupported MariaDB, dependency ordering, and partial-DDL recovery where applicable.

## Jobs

Required behavior includes FIFO eligibility, future availability, short claim transactions, concurrent claim safety, lease heartbeat/fencing, stale-owner rejection, retry timing, max attempts, cancellation, crash/restart recovery, and payload separation from claim scans.

## Workflows

Test prepare/commit/cleanup, production-data isolation before commit, record rejection as result data, crash recovery before/during commit, cleanup semantics, and correct job correlation.

## Scheduling

Use deterministic clocks. Test once, anchored interval, cron, restart behavior, missed runs, coalesce/replay, max start delay, all overlap policies, idempotent crash recovery, `scheduled_for`, and explicit DST spring-forward/fall-back behavior.

## Authentication and Authorization

Test multiple passkeys, final-passkey protection, disabled user/credential behavior, sessions/tokens expiry/revocation, token rotation, recent auth, application/installation capability narrowing, permission union/intersection, deny by default, and alternate entry-path authorization.

## CSRF / Origin

Use real HTTP requests to test unsafe cross-origin rejection, same-origin acceptance, configured additional origins, Fetch Metadata/Origin semantics, bearer-token API behavior, WebAuthn RP ID and exact approved origin.

## Routes

Test literal precedence, ambiguous structural parameter routes, 405, trailing slash normalization, repeated-slash rejection, input conversions, body/query/path/header/fixed/context mappings, and security-header restrictions.

## Cache

Test expiration, invalidation, revision invalidation, concurrent safety, bounded capacity, and correctness with cache disabled.

## LLM

Use a deterministic fake yllmd peer. Test persistence-before-submit, commit-before-wake, ambiguous submission, stable logical request IDs, attempts, crash quarantine, non-crash errors, retention/content purge, review holds, restart recovery, and nonresumable streaming semantics.

## Messaging / Notifications

Use a fake SMTP server. Test persistence before delivery, transient/permanent errors, retry ownership, send-after, attempt history, ambiguous delivery, notification-to-message flow, verification, and secret/content redaction.

## External Integration Runtime

Use controlled HTTP servers to test success, 429/5xx, timeouts, resets, slow/oversized responses, redirects, retry budgets, circuit breakers, rate limits, safe-fetch restrictions, and credential redaction.

## Logging / Lifecycle / Health

Test bootstrap logger before config, full logger activation, fatal startup stderr mirroring, secret redaction, startup phase failure boundaries, readiness transitions, quiesce/drain/DB close/termination, bounded shutdown with noncooperative work, socket cleanup, and bootstrap logger as final fallback.

## Security Audit

If security events are included, test append-only behavior, required auth/authz change events, actor/subject correlation, and absence of secrets.

## Retention

Test active-row preservation, bounded purge batches, content-vs-metadata retention, review holds, restartable cleanup, and prevention of unexpected FK cascades.

## Performance / Query Validation

Before v1, benchmark critical hot paths with realistic table sizes and review MariaDB `EXPLAIN` plans for job claims, schedules, session/token lookup, authorization resolution, LLM queue lookup, and message claims.

No arbitrary universal RPS target is required before baseline measurement. Regressions should be tracked against established baselines.

## Race and Fuzz Testing

Run applicable tests with `go test -race ./...`.

Fuzz parser and external-input boundaries such as routes, config, input mapping, URLs/safe-fetch, token parsing, permission names, and cron-wrapper parsing.

## Consumer Compatibility

`internal/testapp` acts as a minimal consumer and imports only public yFramework packages. It must never depend on private `internal/` implementation packages.

Real external applications such as yCMS may later add integration validation but are not required for yFramework's own test suite.

## CI Order

Recommended CI stages:

1. formatting/vet/static checks,
2. unit tests,
3. race tests,
4. MariaDB integration tests,
5. migration-from-zero/upgrade tests,
6. lifecycle/runtime tests,
7. consumer compile test.

Longer performance/recovery suites may run on release/nightly workflows.

## Definition of Testability

Subsystems accept explicit clocks, transports, persistence dependencies, and adapters where deterministic testing requires them. Tests must not depend on public internet access, arbitrary sleeps, or global mutable service state.

Behavioral acceptance criteria matter more than a single line-coverage percentage.

<!-- ORIGINAL_SOURCE_END: 20_TESTING_STRATEGY.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 21_IMPLEMENTATION_MILESTONES.md -->

# Original Source File: `21_IMPLEMENTATION_MILESTONES.md`

# yFramework Implementation Milestones

## Governing Rule

Every milestone must leave the default branch compiling and tested. Later milestones may not silently redesign settled earlier contracts. If implementation exposes an architectural conflict, record an open decision before proceeding.

Do not scaffold future milestones with fake production implementations merely to satisfy interfaces.

## Milestone 1 - Foundation and Repository Bootstrap

- initialize module `github.com/BrandonYaniz/yframework`,
- repository/package skeleton,
- errors, validation, logging, config,
- strict TOML,
- bootstrap logger,
- baseline `.codex` instructions.

Acceptance: builds; config validation/unknown-key rejection/secret redaction/bootstrap logging tests pass.

## Milestone 2 - MariaDB and Migrations

- mariadb, transaction, migrations,
- compatibility/capability check,
- advisory migration lock,
- `framework_migrations`,
- separate foreground/background pools.

Acceptance: empty migration, checksum drift, concurrent lock, unsupported-version and real-MariaDB tests pass.

## Milestone 3 - Composition and Registries

- application/component packages,
- private registry machinery,
- explicit dependency graph,
- deterministic registration/freeze,
- duplicate/missing/cycle validation.

Acceptance: internal minimal consumer composes; invalid graphs fail before DB mutation.

## Milestone 4 - Identity and Authorization

- users/passkeys/sessions/tokens,
- client/application identities/installations,
- roles/direct grants,
- authorization cache/revisions,
- security audit events.

Acceptance: passkey-only flow, final-passkey protection, revocation, permission semantics, secret-redaction tests pass.

## Milestone 5 - Actions, Routes, Server, Rendering

- typed action registration,
- database-backed routes/input mappings,
- authorization checks,
- UNIX-socket HTTP server,
- renderers,
- immutable route snapshot,
- browser cross-origin protection.

Acceptance: real HTTP route/auth/input/CSRF/render/socket tests pass.

## Milestone 6 - Widgets and Runtime Settings

- widget areas/instances/placements,
- validated immutable widget snapshots,
- registered runtime settings and revisioned snapshots.

Acceptance: invalid candidate config cannot replace active state; matching/order rules pass.

## Milestone 7 - Jobs

- jobs/payloads/attempts schema,
- FIFO claim,
- leases/heartbeats/fencing,
- retries/cancellation/recovery,
- one ordinary executable background lane.

Acceptance: real-MariaDB claim/recovery tests; stale workers cannot commit; FIFO eligibility proven.

## Milestone 8 - Workflows and Scheduling

- workflow coordination,
- prepare/commit/cleanup,
- once/cron/interval schedules,
- coalesce/replay,
- overlap policies,
- DST rules,
- durable `scheduled_for`.

Acceptance: pre-commit production isolation, crash safety, missed schedule/DST/idempotency tests pass.

## Milestone 9 - External Integration Runtime

- normalized HTTP/public-fetch transport,
- retry budgets,
- circuit breakers,
- rate limits,
- response bounds,
- health integration,
- adapter boundaries.

Acceptance: controlled server tests prove resilience/security; components cannot bypass framework integration runtime.

## Milestone 10 - Messaging and Notifications

- durable messages/content/attempts,
- SMTP adapter,
- notification contacts/subscriptions/verification,
- global background-lane execution,
- messaging-owned transport retries.

Acceptance: fake-SMTP tests; notification creates durable message; no direct component SMTP access.

## Milestone 11 - LLM and yllmd

- durable request/content/attempt schema,
- dedicated yFramework LLM FIFO,
- stable logical IDs,
- crash reconciliation/quarantine,
- commit-before-wake,
- retention/health.

Acceptance: submission-crash ambiguity, restart recovery, quarantine and content-retention tests pass.

## Milestone 12 - Localization

- derived-content coordination,
- component-owned localized persistence integration,
- background/idle processing,
- LLM integration where enabled.

Acceptance: no universal domain-content localization table; ownership/restart safety preserved.

## Milestone 13 - Lifecycle, CLI, Health Hardening

- serve/install/migrate/validate,
- full startup phases,
- readiness/degraded aggregation,
- quiesce/drain/DB-close/terminate,
- socket cleanup,
- logging shutdown ordering.

Acceptance: injected failure boundaries and bounded shutdown tests pass.

## Milestone 14 - Public API and Compatibility Freeze

- exported API review,
- minimize public surface,
- consumer compile fixture,
- API docs,
- semantic version/deprecation policy.

Acceptance: public-only consumer builds and no accidental dependency on private internals remains.

## Milestone 15 - Release Hardening

- full race/fuzz suites,
- migration zero/upgrade validation,
- performance baselines and query-plan review,
- retention/load testing,
- security review,
- documentation completion.

Acceptance: all named subsystem criteria pass on supported MariaDB baseline and repository is ready for a tagged release.

<!-- ORIGINAL_SOURCE_END: 21_IMPLEMENTATION_MILESTONES.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 22_ACCEPTANCE_CRITERIA.md -->

# Original Source File: `22_ACCEPTANCE_CRITERIA.md`

# yFramework Acceptance Criteria

## Purpose

This document defines implementation-completion criteria that Codex and human reviewers can use without substituting subjective "looks finished" judgments.

## Global Criteria for Every Milestone

A milestone is complete only when:

- repository builds cleanly,
- tests added for new behavior pass,
- relevant race/integration tests pass,
- no knowingly broken placeholder production implementation remains,
- new exported APIs are documented,
- architecture/source documents are updated when a settled design is refined,
- new unresolved architectural questions are recorded rather than silently decided,
- secrets are not emitted in logs/errors/tests,
- ownership boundaries remain intact,
- default branch is usable as the basis for the next milestone.

## Foundation

- Module path is `github.com/BrandonYaniz/yframework`.
- Public packages and `internal/` boundaries compile as designed.
- Strict TOML rejects unknown keys and invalid cross-field combinations.
- Bootstrap logger works before config and remains available as final fallback.
- Config diagnostics redact every registered secret field.

## MariaDB / Migrations

- Unsupported MariaDB fails before schema mutation.
- Migration lock serializes concurrent runners.
- Checksummed applied migrations cannot be silently modified.
- Failure does not mark a migration applied.
- Framework/application/component ordering is deterministic.
- Empty DB and prior-supported DB reach the expected current schema.
- Foreground/background pools are distinct.

## Composition

- Missing dependencies and cycles fail before DB mutation.
- Duplicate stable registry names fail deterministically.
- Registries freeze before normal runtime begins.
- Components cannot discover arbitrary services by string lookup.

## Identity / Authorization

- Human authentication is passkey-only.
- Multiple passkeys work; final usable passkey cannot be removed.
- Sessions/tokens store only credential verifiers/hashes, never raw secrets.
- Roles and direct grants resolve correctly.
- Application/user/installation permission intersections behave correctly.
- Resource authorization cannot be bypassed by alternate entry path.
- Security-relevant changes generate required audit records if the audit subsystem is enabled by the v1 implementation.

## Routing / Server / Rendering

- Route grammar only permits settled literal/named-param syntax.
- Structural ambiguity is rejected at startup.
- Literal precedence and 405 semantics are correct.
- Unknown permissions/actions/renderers/mappings prevent activation.
- Unsafe browser cross-origin requests are rejected according to framework policy.
- UNIX socket stale/live handling and removal semantics are correct.
- Runtime route snapshots swap atomically after successful validation.

## Jobs

- Eligible FIFO order is deterministic.
- Claim transactions are short and use correct MariaDB locking semantics.
- Lease token fencing prevents stale-worker completion/heartbeat.
- Retry, cancellation, crash recovery, and max-attempt behavior are deterministic.
- Large payload content is not required for claim scans.
- Only one ordinary executable background lane is active in v1.

## Workflows / Schedules

- Transitional component data does not mutate production state before final commit.
- Record rejection remains result data, not execution failure.
- Workflow crash recovery does not duplicate final commits.
- once/cron/interval behavior matches specification.
- missed/coalesce/replay/overlap/max-start-delay behavior is deterministic.
- DST spring-forward occurrence is skipped and repeated fall-back wall time runs once.
- Scheduler recovery does not duplicate a logical occurrence.

## Integrations

- Components cannot bypass framework-normalized external communications.
- Retry budgets and circuit breakers do not multiply retry loops.
- Size/time/rate limits are enforced.
- Provider credentials do not leak to logs or component-visible diagnostics.

## Messaging / Notifications

- Notification intent creates durable messaging work without direct SMTP calls.
- Messaging owns delivery retry semantics.
- SMTP transient/permanent failures are classified correctly.
- Message content and delivery state/history are separated as designed.
- Verification secrets are hash-only and recipient content is not logged.

## LLM

- Request is durable before submission.
- Completion is durable before waiting callers are released.
- Stable logical request identity survives retries.
- Ambiguous submission can be reconciled without silently duplicating application effects.
- Crash-attempt quarantine threshold works.
- 90-day default content retention and review/quarantine exemptions work.
- Streaming is not falsely reported as resumable across process crashes.

## Localization

- Framework owns coordination, not component domain localized content.
- Regeneration/backfill is durable and restart-safe.
- Component-owned source and localized records remain in component schema.

## Lifecycle / Health / CLI

- Startup phases execute in specified order.
- Static invalid configuration fails before migration.
- Readiness becomes true only after full runtime validation/start.
- Quiesce disables readiness and all new work sources before drain.
- DB remains open until correctness-sensitive active work is drained or deadline is reached.
- Shutdown has a hard bound and does not fake successful cancellation of still-running work.
- `install`, `migrate`, `validate`, and `serve` reuse the same composition/config/migration contracts.
- `validate` remains nonmutating.

## Release Criteria

Before a v1 release tag:

- public API review is complete,
- all named behavioral acceptance tests pass,
- race suite passes,
- relevant fuzz suites pass,
- migration from zero and supported upgrade path pass,
- critical MariaDB hot queries have acceptable indexed plans,
- performance baselines exist for job/message/LLM claims and auth lookups,
- security review has no unresolved release-blocking issue,
- user, administrator, and developer documentation are present,
- `.codex` instructions and task graph reflect the final architecture.

<!-- ORIGINAL_SOURCE_END: 22_ACCEPTANCE_CRITERIA.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 23_DOCUMENTATION_OUTLINE.md -->

# Original Source File: `23_DOCUMENTATION_OUTLINE.md`

# yFramework Documentation Outline

## Status

This document defines the documentation structure that accompanies implementation. It is a development planning source file, not end-user application documentation.

## Documentation root

The implementation repository should contain:

```text
docs/
├── architecture/
├── developer/
├── administrator/
├── user/
├── schemas/
└── adr/
```

The frozen planning/source documents belong under `docs/architecture/` when the implementation repository is created.

## User documentation

yFramework is an importable framework rather than an executable end-user application, so user documentation is intentionally small.

Recommended topics:

- What yFramework is and is not
- Supported operating systems
- Supported MariaDB policy
- Module/versioning policy
- How downstream applications consume yFramework
- Support and compatibility expectations
- Security reporting and project policy

Application-specific user documentation belongs in the consuming application repository.

## Administrator documentation

Administrator documentation covers deployment and operation of applications built on yFramework without assuming application-domain knowledge.

Recommended sections:

1. Deployment model
   - Application owns the executable and `main`
   - Foreground `serve` process
   - Service-manager supervision
   - FreeBSD, Linux, and macOS expectations

2. Configuration
   - TOML configuration
   - Required bootstrap settings
   - Secrets and filesystem permissions
   - Dynamic runtime settings
   - Configuration validation

3. MariaDB
   - Version/support policy
   - Required capabilities
   - Foreground/background connection pools
   - Backup expectations
   - Restore expectations

4. Lifecycle commands
   - `install`
   - `migrate`
   - `validate`
   - `serve`
   - `version`

5. Startup and readiness
   - Startup phases
   - Migration behavior
   - Readiness and degraded capability
   - Failure behavior

6. Shutdown
   - Quiesce
   - Drain
   - Database close
   - Socket removal
   - Process termination

7. UNIX socket and reverse proxy
   - Socket ownership and permissions
   - Stale socket handling
   - Trusted proxy rules
   - Public origin configuration

8. Logging and observability
   - Bootstrap logger
   - Structured logging
   - Log destinations
   - Correlation identifiers
   - Secret redaction

9. Health
   - Liveness
   - Readiness
   - Degraded capabilities
   - Public versus detailed diagnostics

10. Authentication administration
    - Passkey-only user access
    - Multiple passkeys
    - No password/MFA/recovery in v1
    - Session revocation
    - Personal API tokens

11. Application and installation identities
    - Application tokens
    - Installation tokens
    - Token rotation
    - Capability administration

12. Authorization administration
    - Permissions
    - Roles
    - Direct grants
    - Security audit history

13. Background work
    - Jobs
    - Workflows
    - Scheduling
    - Retry and recovery
    - Cancellation
    - Concurrency locks

14. Messaging and notifications
    - SMTP configuration
    - Delivery retries
    - Contact verification
    - Operational diagnostics

15. LLM/yllmd
    - Enabling LLM support
    - yllmd connection
    - Request retention
    - Quarantine/review
    - Health behavior

16. Data retention and cleanup
    - LLM content
    - Message content
    - Job/workflow history
    - Security events
    - Batched cleanup

17. Upgrades
    - Framework dependency update
    - Pre-running migrations
    - Validation
    - Rollback expectations
    - Expand/migrate/contract guidance

18. Troubleshooting
    - Startup failure
    - Migration lock failure
    - Unsupported MariaDB
    - Readiness failure
    - Stale socket
    - Background recovery
    - SMTP/yllmd/provider degradation

## Developer documentation

Developer documentation is the primary documentation set for framework consumers and component authors.

Recommended sections:

1. Getting started
   - Module import path
   - Repository model
   - Application owns `main`
   - No application template embedded in yFramework

2. Architecture overview
   - Framework/application/component ownership
   - Explicit dependency injection
   - Typed APIs
   - Stable registry names
   - Startup validation

3. Application composition
   - Builder/composition root
   - Component registration
   - Adapter selection
   - Dependency resolution

4. Components
   - Component contract
   - Component IDs
   - Dependencies
   - Registration
   - Schema ownership
   - Prohibited behaviors

5. Actions
   - Typed inputs/results
   - Validation
   - Error model
   - Authorization boundary

6. Routes
   - Route grammar
   - Methods
   - Access policy
   - Required permissions
   - Input mappings
   - Runtime snapshots

7. Authentication
   - Identity types
   - Passkey boundary
   - Sessions
   - Bearer tokens
   - Recent authentication

8. Authorization
   - Permission registration
   - Roles and direct grants
   - Application/installation capabilities
   - Resource authorization
   - Trusted/system actors

9. MariaDB and repositories
   - Explicit database dependency
   - Component-owned repositories
   - No universal CRUD repository
   - Transactions
   - Query/index guidance

10. Migrations
    - Registration
    - Ordering
    - Checksums
    - Crash recovery
    - Expand/migrate/contract
    - Large data backfills

11. Cache
    - Appropriate cache use
    - Revision/invalidation
    - Security-sensitive state
    - Database remains authoritative

12. Rendering and widgets
    - Renderers
    - Templates
    - Widget registration
    - Placement
    - Snapshot behavior

13. Jobs
    - Registration
    - Payload/result
    - Idempotency
    - Retry
    - Leases
    - Cancellation

14. Workflows
    - Prepare/commit/cleanup
    - Transitional data ownership
    - Production commit requirement
    - Record-level rejection

15. Scheduling
    - Once
    - Cron
    - Interval
    - Missed schedules
    - Overlap
    - DST semantics

16. Hooks and events
    - Synchronous execution
    - Durable effects via jobs/messages
    - No component worker goroutines

17. Messaging and notifications
    - Message intent
    - SMTP abstraction
    - Notification subscriptions
    - Retry ownership

18. External integrations
    - Typed component ports
    - Application-selected adapters
    - Framework transport/runtime
    - Retry budgets and circuit breakers
    - Safe web fetch

19. LLM
    - Registered purposes
    - Durable requests
    - yllmd boundary
    - Crash handling
    - Retention

20. Localization
    - Component-owned localized content
    - Persistent derived data
    - Backfill/regeneration

21. Logging and observability
    - Structured fields
    - Correlation IDs
    - Sensitive data rules

22. Testing
    - Public testing helpers
    - MariaDB integration tests
    - Fake transports
    - Deterministic clocks
    - Consumer compile fixture

23. Public API compatibility
    - Semantic versioning
    - Exported symbol policy
    - Deprecation
    - Breaking-change discipline

## Schema documentation

`docs/schemas/` should contain implementation-aligned references for:

- Framework MariaDB tables
- Foreign-key and retention policies
- Critical indexes and the queries they support
- TOML configuration schema
- Dynamic runtime settings
- Stable permission names
- Stable route/action/job/workflow/widget/LLM-purpose names where useful
- Status/state definitions

Schema documentation must describe the implemented schema and should be updated in the same changeset as schema changes.

## Architecture documentation

`docs/architecture/` contains the frozen planning/source documents and later architectural clarifications.

Do not rewrite settled history to conceal architectural changes. Material future architectural changes should receive an ADR.

## ADR policy

`docs/adr/` is for decisions made after the initial architecture freeze when a future change materially affects public API, persistence, security, lifecycle, ownership, failure/recovery, compatibility, or external behavior.

Do not create ADRs merely for ordinary private implementation choices.

<!-- ORIGINAL_SOURCE_END: 23_DOCUMENTATION_OUTLINE.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 24_CODEX_OPERATING_SPEC.md -->

# Original Source File: `24_CODEX_OPERATING_SPEC.md`

# yFramework Codex Operating Specification

## Purpose

The `.codex` directory is the implementation control plane for Codex. It is not a second architecture specification. Codex must use the source documents as authoritative input and use `.codex` to track implementation scope, tasks, decisions, blockers, and completion state.

## Required structure

```text
.codex/
├── README.md
├── INSTRUCTIONS.md
├── DECISION_AUTHORITY.md
├── ARCHITECTURE_INDEX.md
├── CURRENT_MILESTONE.md
├── ACCEPTANCE.md
├── TASKS/
└── STATE/
    ├── completed.md
    ├── blockers.md
    └── open-decisions.md
```

## Core operating rule

Codex may make private implementation choices that do not alter:

- documented public APIs
- persisted-data semantics
- security boundaries
- ownership boundaries
- lifecycle behavior
- compatibility requirements
- failure/recovery semantics
- externally observable framework behavior

If an implementation choice would affect one of those categories, Codex must record an open decision and stop that architectural branch rather than silently choosing.

## Architecture preservation rules

Codex must preserve the following unless a source document is explicitly changed:

- MariaDB is the only relational database.
- The application owns `main` and the executable.
- yFramework is an importable Go module, not an application template.
- Components do not contain application-specific domain logic.
- Components do not directly communicate with external systems.
- Components do not create independent worker pools/background services in v1.
- Dependencies are explicit and typed.
- No global service locator or generic `Get(name)` dependency system.
- No executable scripts, arbitrary expressions, or policy languages stored in MariaDB.
- Production workflow changes occur only in final commit phase.
- Record rejection is result data, not execution failure.
- Ordinary async work uses one global FIFO execution lane in v1.
- The LLM queue is a separate explicit exception.
- Passkey-only user authentication remains intact.
- Authorization checks permissions/capabilities, not role names in business code.

## Task discipline

Every task file must contain:

- task ID
- title
- milestone
- purpose
- source documents
- prerequisites
- scope
- expected packages/files
- behavioral requirements
- required tests
- acceptance criteria
- forbidden architectural changes
- completion notes

Tasks should be small enough to review independently.

Bad task:

```text
Implement authentication.
```

Good task:

```text
Create framework passkey persistence schema and repository with credential-ID uniqueness, multi-passkey user support, and final-passkey-safe repository operations.
```

## Branch/state rule

Each completed task must leave the repository compiling and relevant tests passing.

Do not scaffold later milestones with fake production implementations simply to satisfy compilation.

Temporary test fakes belong in test-only code.

## Source hierarchy

When sources conflict, use this precedence:

1. Newer explicitly settled source document decision
2. Subsystem architecture document
3. Cross-cutting development-plan document
4. Task file
5. Private implementation assumption

If conflict cannot be resolved by this order, add it to `.codex/STATE/open-decisions.md`.

## Public API discipline

- Minimize exported symbols.
- Prefer narrow typed interfaces owned by the consumer.
- Do not expose internal registry/runtime implementation.
- Do not use `context.Context` as a service container.
- Do not introduce generic framework-wide CRUD repositories.
- Do not expose internal numeric database IDs across subsystem/component/API boundaries.

## Database discipline

- Use framework-owned `framework_` tables only for framework state.
- Components own component tables/migrations.
- Respect dual IDs where specified: internal numeric key plus stable UUID.
- Keep high-write operational rows narrow.
- Separate large content/history where specified.
- Do not add speculative indexes.
- Every critical index should support a documented query.
- Do not introduce schema auto-sync or inferred destructive migrations.

## Testing discipline

- Unit tests for deterministic behavior.
- Real MariaDB tests for database semantics.
- Race tests for concurrent in-memory/runtime code.
- Failure injection for lifecycle and crash recovery.
- No tests requiring public internet access.
- Avoid arbitrary sleeps; inject clocks and controlled transports where needed.

## Documentation discipline

A task that changes public behavior, configuration, schema, or administration is incomplete until corresponding documentation is updated.

## Completion notes

Task completion notes should record:

- files changed
- tests run
- notable private implementation decisions
- any follow-up tasks discovered
- confirmation that no unresolved architectural decision was silently made

<!-- ORIGINAL_SOURCE_END: 24_CODEX_OPERATING_SPEC.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 25_INITIAL_CODEX_TASK_GRAPH.md -->

# Original Source File: `25_INITIAL_CODEX_TASK_GRAPH.md`

# Initial Codex Task Graph

## Scope

This task graph covers Milestones 1 through 3 only. Later milestone task graphs should be generated after these foundations compile and their public/private boundaries have been validated.

# Milestone 1: Foundation and Repository Bootstrap

## 001 Repository bootstrap

Purpose:
- Create the public Go module `github.com/BrandonYaniz/yframework`.
- Establish repository directories and baseline project files.

Requirements:
- No `examples/` directory.
- No framework-owned executable application.
- No `cmd/yframework` runtime executable.
- Public packages follow the package plan.
- Private implementation uses `internal/`.

Tests/acceptance:
- `go test ./...` succeeds on initial repository.
- Module path is exact.

Depends on: none.

## 002 Error model foundation

Purpose:
- Implement shared framework error primitives used by higher packages.

Requirements:
- Stable code concept.
- Public-safe message separate from underlying cause.
- Structured validation/field-error support where needed.
- No HTTP-specific ownership in core error type.

Tests:
- wrapping/preservation
- public-message behavior
- stable code behavior

Depends on: 001.

## 003 Validation foundation

Purpose:
- Implement reusable validation primitives without coupling to routes or database code.

Tests:
- field errors
- typed validation aggregation
- deterministic output

Depends on: 002.

## 004 Bootstrap logger

Purpose:
- Implement the process-lifetime stderr bootstrap logger.

Requirements:
- Available before full config.
- Remains available as final fallback after full logger shutdown.
- No platform-specific syslog/service-manager assumptions.

Tests:
- stderr capture
- startup failure output
- post-full-logger fallback

Depends on: 001, 002.

## 005 Structured logging core

Purpose:
- Implement framework logging interfaces and structured full logger foundation.

Requirements:
- Structured fields.
- Secret-safe formatting hooks.
- Bootstrap logger remains separate fallback.
- Do not force a third-party logger type into public component APIs.

Tests:
- level filtering
- structured fields
- redaction
- output lifecycle

Depends on: 002, 004.

## 006 Typed TOML configuration parser

Purpose:
- Implement strict top-level typed deployment configuration.

Requirements:
- Reject unknown keys.
- Apply compiled defaults.
- Support secret metadata/redaction.
- No generic `map[string]any` as primary framework config representation.
- TOML is the v1 deployment format.

Tests:
- minimal valid config
- defaults
- unknown key rejection
- invalid type rejection
- secret redaction

Depends on: 002, 003, 005.

## 007 Configuration cross-field validation

Purpose:
- Add static validation for database pools, background timings, server/auth origin relationships, and other settled invariants.

Tests:
- heartbeat < lease
- max idle <= max open
- invalid/conflicting DB connection modes
- required public origin/RP ID validation

Depends on: 006.

## 008 Foundation public API review

Purpose:
- Review Milestone 1 exported symbols and move unnecessary implementation details under `internal/`.

Acceptance:
- no service locator
- no generic global config maps
- `go test ./...`
- `go vet ./...`

Depends on: 002-007.

# Milestone 2: MariaDB and Migrations

## 009 MariaDB connection/config integration

Purpose:
- Implement MariaDB connection construction from typed config.

Requirements:
- MariaDB only.
- Dedicated startup/migration connection support.
- Runtime foreground/background pools are separately configurable.

Tests:
- valid/invalid config
- connection failure behavior

Depends on: 007.

## 010 MariaDB version/capability validation

Purpose:
- Verify supported MariaDB version and required capabilities before migration mutation.

Requirements:
- Release baseline derives from the documented FreeBSD ports policy.
- Unsupported database fails before schema mutation.

Tests:
- supported server accepted
- unsupported version rejected
- capability failure rejected

Depends on: 009.

## 011 Transaction helper foundation

Purpose:
- Implement explicit transaction helper primitives compatible with `database/sql`.

Requirements:
- No global DB handle.
- Keep public abstractions narrow.

Tests:
- commit
- rollback on error
- context cancellation

Depends on: 009.

## 012 Migration registry

Purpose:
- Implement registration of framework/application/component migrations.

Requirements:
- Stable owner type/name/migration identity.
- Stable checksums.
- Immutable applied migrations.
- Registration is DB-free.

Tests:
- duplicate registration rejection
- deterministic ordering within owner
- checksum identity

Depends on: 002.

## 013 Migration history schema

Purpose:
- Add initial `framework_migrations` schema and repository.

Requirements:
- permanent append-only applied history
- unique owner/migration identity
- checksum stored

Tests:
- schema creation
- duplicate prevention

Depends on: 010, 012.

## 014 Advisory migration lock

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

## 015 Migration runner

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

## 016 Runtime MariaDB pools

Purpose:
- Construct foreground and background pools only after successful migration preparation.

Tests:
- separate limits
- cleanup on construction failure

Depends on: 009, 015.

## 017 Migration CLI/service API primitives

Purpose:
- Expose reusable non-executable migration/validation operations for later CLI composition.

Requirements:
- no framework `main`
- migrate mutates
- validate reports without mutation

Depends on: 015, 016.

## 018 Milestone 2 MariaDB integration suite

Purpose:
- Run real MariaDB tests for migrations, locks, transactions, and supported-version behavior.

Acceptance:
- migration-from-zero passes
- concurrent migration test passes
- checksum drift test passes
- unsupported DB performs no schema mutation

Depends on: 009-017.

# Milestone 3: Composition and Registries

## 019 Stable component identity types

Purpose:
- Define component IDs and dependency declaration types.

Requirements:
- explicit dependencies
- no runtime string-based service lookup

Depends on: 008.

## 020 Generic private frozen-registry primitive

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

## 021 Component dependency graph

Purpose:
- Validate and sort required component dependencies.

Tests:
- deterministic topological order
- missing dependency failure
- cycle detection

Depends on: 019.

## 022 Component registrar skeleton

Purpose:
- Implement the composed typed registrar boundary.

Requirements:
- no `GetService`
- no generic service container
- registrar provides typed subsystem registration interfaces
- initial subsystem registrars may be minimal contracts until later milestones

Depends on: 019, 020.

## 023 Application builder skeleton

Purpose:
- Implement application composition root without lifecycle/server implementation.

Requirements:
- application owns selected components and adapters
- yFramework does not own `main`
- application builder validates composition before DB mutation

Depends on: 021, 022.

## 024 Registration/freeze pipeline

Purpose:
- Compose components, validate dependency graph, register definitions in deterministic order, and freeze registries.

Tests:
- duplicate stable names fail
- component registration order deterministic
- mutation after freeze fails
- failure occurs before DB mutation

Depends on: 020-023.

## 025 Internal consumer test application

Purpose:
- Create `internal/testapp` fixture that consumes only public framework APIs as a downstream application would.

Requirements:
- not an example application
- not a production executable template
- must not import private framework implementation except where test package rules explicitly require internal test harnessing

Acceptance:
- fixture compiles
- basic application composition succeeds

Depends on: 023, 024.

## 026 Milestone 3 public API review

Purpose:
- Review composition/component exported APIs before later subsystems build on them.

Acceptance:
- no service locator
- no circular package dependencies
- minimal exported symbols
- consumer fixture compiles
- `go test -race ./...` passes for implemented concurrent registry code

Depends on: 019-025.

# Gate before Milestone 4

Milestones 1-3 are complete only when:

- repository is cleanly importable as `github.com/BrandonYaniz/yframework`
- configuration is strict and typed
- MariaDB/migrations work against a real server
- component dependency resolution is deterministic
- registries freeze correctly
- application composition exists without owning `main`
- no service locator has been introduced
- internal consumer fixture compiles
- all relevant tests pass

Only then should the auth/authorization implementation task graph be generated.

<!-- ORIGINAL_SOURCE_END: 25_INITIAL_CODEX_TASK_GRAPH.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 26_FULL_CODEX_TASK_GRAPH_M04_M15.md -->

# Original Source File: `26_FULL_CODEX_TASK_GRAPH_M04_M15.md`

# Full Codex Task Graph, Milestones 4 through 15

## Status and use

This document extends `25_INITIAL_CODEX_TASK_GRAPH.md`. Tasks 001 through 026 remain the detailed Milestone 1 through 3 foundation graph. Tasks 027 through 881 below are the settled implementation-planning task-ID and sequencing index for Milestones 4 through 15.

Before Codex executes any task below, `.codex/TASKS/` must contain the corresponding detailed task file. That file must use the authoritative subsystem specifications and satisfy the task-file contract in `24_CODEX_OPERATING_SPEC.md`: sources, prerequisites, scope, expected files/packages, behavioral requirements, tests, acceptance criteria, forbidden architectural changes, and completion notes. The task index does not grant Codex authority to invent public, persisted, security, lifecycle, compatibility, or failure-recovery semantics.

# Milestone 4: Identity and Authorization

- **027** Identity schema migrations
- **028** User repository
- **029** Passkey repository
- **030** Final-passkey invariant
- **031** WebAuthn library integration prototype
- **032** Passkey registration service
- **033** Passkey authentication service
- **034** Session persistence
- **035** Session issuance/rotation
- **036** Personal bearer-token persistence
- **037** Roles/user auth schema
- **038** Permission registry
- **039** User authorization resolver
- **040** Authorization revision/cache
- **041** Application-principal schema
- **042** Installation schema/persistence
- **043** Client token services
- **044** Client capability resolution
- **045** Framework actor/identity model
- **046** Authorization checker/service
- **047** Recent-auth service
- **048** Security-event schema
- **049** Security-event writer
- **050** Integrate security events
- **051** Identity secret-redaction pass
- **052** Authentication concurrency tests real MariaDB
- **053** Authentication restart tests
- **054** Milestone 4 public API review
- **055** Milestone 4 gate

**Milestone 4 gate:** Passkey-only identity, token/session safety, permission semantics, capability intersections, authorization revision invalidation, and security-event behavior are proven.

# Milestone 5: Actions, Routes, Server, and Rendering

- **056** Action registry contract
- **057** Action execution adapter
- **058** Route schema migrations
- **059** Route definition types
- **060** Route grammar parser
- **061** Structural ambiguity validator
- **062** Route matcher
- **063** Path normalization
- **064** Route hierarchy validation
- **065** Input mapping representation
- **066** Action input introspection
- **067** Mapping validator
- **068** Scalar conversion
- **069** Request body mapping
- **070** Header mapping boundary
- **071** Route access-policy evaluator
- **072** Route permission gate
- **073** Recent-auth route requirement
- **074** Action-level auth enforcement pattern
- **075** Route repository
- **076** Route snapshot builder
- **077** Immutable route snapshot manager
- **078** Server config types
- **079** Safe UDS preparation
- **080** HTTP listener adapter
- **081** Request ID/trace init
- **082** Request identity middleware
- **083** Trusted proxy boundary
- **084** Cookie policy
- **085** Cross-origin middleware
- **086** WebAuthn HTTP adapters
- **087** Renderer registry
- **088** Framework response abstraction
- **089** JSON renderer
- **090** Text/none renderers
- **091** Template registry/loading
- **092** Web page context
- **093** HTML renderer
- **094** Template link resolver
- **095** Asset manager foundation
- **096** HTTP dispatch pipeline
- **097** Request cancellation propagation
- **098** HTTP error mapping
- **099** Access logging
- **100** Full route activation validation
- **101** DB-route integration real MariaDB
- **102** Browser-security integration
- **103** UDS lifecycle suite
- **104** Alternate-entry auth suite
- **105** Snapshot atomicity
- **106** M5 performance baseline
- **107** M5 public API review
- **108** Milestone 5 gate

**Milestone 5 gate:** Typed actions, validated DB-backed routes, secure UDS HTTP serving, browser-origin protections, rendering, and immutable route snapshots are proven.

# Milestone 6: Widgets and Runtime Settings

- **109** Widget definition registry
- **110** Widget area schema
- **111** Widget instance schema
- **112** Widget placement schema
- **113** Widget placement scope model
- **114** Widget placement hierarchy semantics
- **115** Widget placement resolver
- **116** Widget instance configuration decoder
- **117** Widget display preparation API
- **118** Widget rendering pipeline
- **119** Widget template resolution
- **120** Template widget functions
- **121** Stored-content widget directive parser
- **122** Widget recursion and resource limits
- **123** Widget asset integration
- **124** Widget repository
- **125** Widget candidate validator
- **126** Immutable widget snapshot
- **127** Widget snapshot integration tests
- **128** Runtime-setting definition registry
- **129** Typed runtime-setting API prototype
- **130** Runtime settings schema
- **131** Configuration precedence implementation
- **132** Runtime-setting decoder
- **133** Runtime-setting immutable snapshot
- **134** Runtime-setting mutation service
- **135** Runtime-setting optimistic concurrency
- **136** Reset runtime override
- **137** Runtime-setting change event
- **138** Dynamic-setting activation policies
- **139** Secret runtime-setting handling
- **140** Unknown runtime-setting validation
- **141** Runtime settings startup loader
- **142** Snapshot generation model
- **143** Cross-snapshot consistency
- **144** Administrative mutation boundary
- **145** Widget rule matrix tests
- **146** Widget recursion security tests
- **147** Runtime-setting persistence tests
- **148** Snapshot failure tests
- **149** Cache interaction tests
- **150** Concurrency/race tests
- **151** Milestone 6 public API review
- **152** Milestone 6 gate

**Milestone 6 gate:** Widget placement/configuration and typed dynamic runtime-setting snapshots activate atomically and reject invalid candidates.

# Milestone 7: Jobs

- **153** Job definition registry
- **154** Typed job API prototype
- **155** Job schema migration
- **156** Job payload/result schema
- **157** Job attempt schema
- **158** Job repository
- **159** Job enqueue service
- **160** Idempotency-key behavior
- **161** Local wake signal
- **162** Eligibility query
- **163** Single-lane execution coordinator
- **164** Adaptive poller
- **165** Claim transaction
- **166** Concurrent claimer tests
- **167** Worker identity
- **168** Start transition
- **169** Handler context
- **170** Lease generation
- **171** Heartbeat loop
- **172** Heartbeat fault handling
- **173** Lease-fenced progress updates
- **174** Successful completion transaction
- **175** Failed execution transition
- **176** Panic recovery
- **177** Retry-policy implementation
- **178** Retry-delay calculation
- **179** Max-attempt behavior
- **180** Queued cancellation
- **181** Active cancellation
- **182** Remote cancellation observation
- **183** Cancellation grace behavior
- **184** Execution deadline
- **185** Non-cooperative handler test
- **186** Expired-claim recovery
- **187** Expired-running recovery
- **188** Recovery scanner
- **189** Startup recovery pass
- **190** Completion fencing test
- **191** Failure fencing test
- **192** Concurrency-lock schema
- **193** Job-scoped concurrency-key acquisition
- **194** Concurrency-key contention
- **195** Concurrency-lock expiry recovery
- **196** Execution-policy model
- **197** Idle-state provider boundary
- **198** Global background pause
- **199** Automatic throttled/paused state boundary
- **200** Job retry ownership contract
- **201** Job structured logging
- **202** Metrics/operational counters
- **203** Completed-job cleanup
- **204** FIFO integration tests
- **205** Claim transaction tests
- **206** Crash boundary matrix
- **207** Restart/recovery suite
- **208** Lease timing tests
- **209** Retry tests
- **210** Cancellation tests
- **211** Payload-size/query-plan test
- **212** Million-row claim benchmark
- **213** Job API review
- **214** Internal consumer job fixture
- **215** Milestone 7 documentation
- **216** Milestone 7 gate

**Milestone 7 gate:** One ordinary durable FIFO lane, MariaDB claims, leases, fencing, retries, cancellation, recovery, and concurrency locks are proven.

# Milestone 8: Workflows and Scheduling

- **217** Workflow definition registry
- **218** Workflow public API prototype
- **219** Workflow schema migration
- **220** Workflow-step schema
- **221** Workflow creation service
- **222** Prepare phase coordinator
- **223** Commit phase coordinator
- **224** Commit fencing contract
- **225** Commit idempotency
- **226** Cleanup phase coordinator
- **227** Transitional-data ownership contract
- **228** Record-rejection semantics
- **229** Workflow progress aggregation
- **230** Job completion to workflow transition
- **231** Workflow advancement idempotency
- **232** Workflow coordinator
- **233** Prepare-step failure
- **234** Commit-step failure
- **235** Cleanup failure
- **236** Queued/pre-commit cancellation
- **237** Commit-phase cancellation
- **238** Post-commit cancellation
- **239** Restart state evaluator
- **240** Workflow crash matrix
- **241** Schedule schema migration
- **242** Schedule definition/runtime types
- **243** Cron parser wrapper
- **244** Interval schedule calculation
- **245** Once schedule
- **246** Spring-forward semantics
- **247** Fall-back semantics
- **248** DST deterministic tests
- **249** Durable next occurrence
- **250** Schedule poller
- **251** Logical occurrence key
- **252** scheduled_for propagation
- **253** Coalesce policy
- **254** Replay policy
- **255** max_start_delay
- **256** allow
- **257** skip
- **258** queue_one
- **259** Overlap concurrency-key integration
- **260** Policy matrix tests
- **261** Registered target validation
- **262** Typed schedule payload validation
- **263** Schedule repository
- **264** Schedule configuration validator
- **265** Resource pause semantics
- **266** Scheduler downtime recovery
- **267** Concurrent scheduler test
- **268** Scheduler crash matrix
- **269** Workflow events
- **270** Schedule events
- **271** Workflow happy-path integration
- **272** Record rejection test
- **273** Failed-prepare test
- **274** Commit crash/retry test
- **275** Cleanup failure test
- **276** Workflow cancellation matrix
- **277** Schedule once test
- **278** Anchored interval test
- **279** Cron restart test
- **280** Missed schedule tests
- **281** Overlap tests
- **282** DST suite
- **283** Schedule query-plan benchmark
- **284** Workflow public API review
- **285** Schedule public API review
- **286** Internal testapp workflow fixture
- **287** Administrator documentation
- **288** Developer documentation
- **289** Milestone 8 gate

**Milestone 8 gate:** Sequential prepare/commit/cleanup workflows and once/cron/interval schedules satisfy missed-run, overlap, DST, idempotency, and crash-recovery semantics.

# Milestone 9: External Integration Runtime

- **290** Integration definition registry
- **291** Typed domain-port contract enforcement
- **292** External transport API prototype
- **293** Credential-provider boundary
- **294** Credential redaction tests
- **295** HTTP client construction
- **296** Request timeout hierarchy
- **297** Response-size enforcement
- **298** Request-body bounds
- **299** Redirect policy
- **300** Safe-fetch API
- **301** SSRF protection
- **302** DNS/rebinding protection tests
- **303** Integration error model
- **304** HTTP status classification
- **305** Retry-policy primitive
- **306** Retry-attempt accounting
- **307** Provider SDK retry suppression
- **308** Aggregate retry-budget model
- **309** Retry-budget accounting
- **310** Retry-budget exhaustion behavior
- **311** Outbound rate limiter
- **312** Server-provided rate-limit hints
- **313** Circuit-breaker state machine
- **314** Circuit breaker and retry-budget interaction
- **315** Integration health reporter
- **316** Readiness aggregation integration
- **317** Integration health probing
- **318** Integration structured events
- **319** Integration log fields
- **320** Metrics
- **321** Fake transport
- **322** Controlled HTTP test server
- **323** Deterministic retry testing
- **324** Adapter startup validation
- **325** Required binding validation
- **326** Adapter contract-test harness
- **327** Synchronous integration path
- **328** Durable external side effects
- **329** Retry ownership matrix
- **330** Header sanitization
- **331** URL validation
- **332** TLS defaults
- **333** Response-data trust boundary
- **334** Connection-pool reuse test
- **335** Concurrency-limit test
- **336** Outage storm simulation
- **337** External package review
- **338** Internal testapp adapter fixture
- **339** Developer integration guide
- **340** Administrator integration guide
- **341** Milestone 9 gate

**Milestone 9 gate:** Typed adapter boundaries, bounded HTTP/safe-fetch transport, retries, budgets, rate limits, breakers, health, SSRF protection, and redaction are proven.

# Milestone 10: Messaging and Notifications

- **342** Message-purpose registry
- **343** Typed messaging API prototype
- **344** Messaging schema migration
- **345** Message-content schema
- **346** Message-attempt schema
- **347** Message repository
- **348** Durable message creation service
- **349** Message idempotency
- **350** SMTP adapter
- **351** SMTP connection policy
- **352** SMTP sender-identity policy
- **353** SMTP envelope/header separation
- **354** Message eligibility
- **355** Message execution through global lane
- **356** Message claim semantics
- **357** SMTP ambiguous-delivery handling
- **358** Stable message identity across retry
- **359** Messaging retry policy
- **360** Temporary failure retry
- **361** Permanent failure
- **362** Authentication/configuration failure
- **363** SMTP retry-budget integration
- **364** SMTP circuit breaker
- **365** Notification-contact schema
- **366** Destination protection
- **367** Destination normalization
- **368** Contact deduplication
- **369** Verification-token generation
- **370** Verification message flow
- **371** Verification cooldown
- **372** Verification completion
- **373** Unverified-contact retention
- **374** Notification subscription schema
- **375** Typed notification API prototype
- **376** Anonymous subscription enrollment
- **377** User-associated contacts
- **378** Durable completion fact boundary
- **379** Subscription matching
- **380** Verified-only delivery
- **381** Notification-to-message creation
- **382** Fulfillment semantics
- **383** Localized operational messages
- **384** Localization fallback policy
- **385** Enrollment rate limiting
- **386** Verification-send rate limiting
- **387** Purpose limitation
- **388** Destination log redaction
- **389** Message-content redaction
- **390** Message-content retention
- **391** Message attempt retention
- **392** Contact/subscription cleanup
- **393** Messaging capability health
- **394** Readiness classification
- **395** Message inspection API
- **396** Retry/review operation
- **397** Contact administration
- **398** Fake SMTP server
- **399** Persistence-before-delivery test
- **400** SMTP retry test
- **401** Permanent rejection test
- **402** Ambiguous-delivery test
- **403** Retry-budget outage test
- **404** Notification verification test
- **405** Notification-before-verification test
- **406** Subscription deduplication test
- **407** Purpose-limitation test
- **408** Secret/redaction test
- **409** Cleanup test
- **410** Concurrency test
- **411** Messaging API review
- **412** Notification API review
- **413** Internal testapp messaging fixture
- **414** Developer messaging guide
- **415** Developer notification guide
- **416** Administrator messaging guide
- **417** Privacy/retention guide
- **418** Milestone 10 gate

**Milestone 10 gate:** Durable messaging and notification verification/subscription flows work through the global background lane with messaging-owned delivery retries.

# Milestone 11: LLM and yllmd

- **419** LLM purpose registry
- **420** Purpose-to-yllmd target mapping
- **421** Typed LLM API prototype
- **422** LLM schema migration
- **423** Authoritative LLM request row
- **424** LLM content table
- **425** LLM attempt table
- **426** Durable request creation service
- **427** Stable yllmd request identity
- **428** LLM FIFO coordinator
- **429** LLM poller/wake behavior
- **430** LLM request eligibility
- **431** yllmd Unix-socket client
- **432** yllmd protocol adapter
- **433** Submission transition
- **434** Accepted transition
- **435** Running transition
- **436** Successful terminal completion
- **437** Caller wake mechanism
- **438** Completion-after-process-restart behavior
- **439** LLM failure classification
- **440** Generation failure
- **441** Invalid output
- **442** Uncertain submission state
- **443** Reconciliation strategy
- **444** Uncertain-submission recovery test
- **445** Crash-attempt counter
- **446** Quarantine transition
- **447** Manual retry from quarantine
- **448** Queued cancellation
- **449** Submitted/running cancellation
- **450** Cancellation propagation from foreground
- **451** Cancellation propagation from jobs/workflows
- **452** Outermost deadline rule
- **453** Purpose/application timeout defaults
- **454** Queue-expiration behavior
- **455** Streaming API
- **456** Streaming persistence boundary
- **457** Streaming restart semantics
- **458** LLM throughput sample collection
- **459** Rolling sample window
- **460** Per-model partitioning
- **461** Output-token generation rate
- **462** Time-to-first-output
- **463** Estimate API
- **464** Queue-ahead estimation
- **465** Request-size estimation
- **466** Output-size estimation
- **467** Presentation rounding
- **468** LLM disabled state
- **469** Optional LLM unavailable
- **470** Required LLM unavailable
- **471** yllmd health monitoring
- **472** LLM output trust boundary
- **473** Prompt/output logging exclusion
- **474** Administrative content permissions
- **475** LLM content retention
- **476** Metadata retention
- **477** Review/quarantine retention hold
- **478** Bounded LLM cleanup
- **479** Stable LLM events
- **480** LLM structured fields
- **481** Deterministic fake yllmd server
- **482** Persistence-before-submit test
- **483** Commit-before-wake test
- **484** Crash quarantine test
- **485** Restart recovery test
- **486** Streaming interruption test
- **487** Retention test
- **488** Throughput estimator test
- **489** LLM claim/query-plan benchmark
- **490** LLM public API review
- **491** Internal testapp LLM fixture
- **492** Developer LLM guide
- **493** Administrator yllmd guide
- **494** Milestone 11 gate

**Milestone 11 gate:** Durable LLM FIFO, yllmd transport, stable identities, ambiguous-submission recovery, quarantine, cancellation, retention, and forecasting are proven.

# Milestone 12: Localization

- **495** Localization resource registry
- **496** Stable localization-unit contract
- **497** Source fingerprinting
- **498** Localization eligibility callback
- **499** Application locale configuration
- **500** Locale identifier handling
- **501** Locale categories
- **502** Unsupported-locale behavior
- **503** Framework localization status model
- **504** Localization coordination record
- **505** Component-owned localized persistence
- **506** Eligible-source comparison
- **507** New localization-unit behavior
- **508** Removed localization-unit behavior
- **509** Human-modified state
- **510** Locked localization state
- **511** Source change against human-modified/locked unit
- **512** Unlock/review operation
- **513** Preferred localization enqueue
- **514** Idle-only localization job
- **515** Preferred backfill
- **516** Backfill batching
- **517** On-demand request API
- **518** On-demand deduplication key
- **519** Join-existing operation
- **520** Wait versus return policy
- **521** Localization LLM purpose
- **522** Localization request construction
- **523** Structured localization output
- **524** LLM result validation
- **525** Component localization commit contract
- **526** Commit-before-current
- **527** Localization idempotency
- **528** Source-changed-during-generation behavior
- **529** Locked-during-generation behavior
- **530** Localization failure transition
- **531** Public failed-localization behavior
- **532** Crash-quarantined LLM request
- **533** Operational-message registry
- **534** Operational-message localization storage
- **535** Operational-message generation
- **536** Bootstrap availability problem
- **537** Meaningful localization-demand counter
- **538** Aggregate demand
- **539** Adaptive promotion threshold
- **540** No automatic demotion
- **541** Manual preferred override
- **542** Localization availability API
- **543** Current localized content path
- **544** Pending rendering state
- **545** Ready-time estimate integration
- **546** Localization completion fact
- **547** Ready-notification subscription
- **548** No notification on failed commit
- **549** Duplicate completion handling
- **550** Localize-resource operation
- **551** Localize-missing-locale operation
- **552** Regenerate stale content
- **553** Localization status query
- **554** Review-required operation
- **555** Localization events
- **556** Localization log fields
- **557** Stable-unit change test
- **558** Preferred-locale proactive test
- **559** On-demand locale test
- **560** Human-modified protection test
- **561** Locked-unit protection test
- **562** Source-race test
- **563** Lock-race test
- **564** Failure test
- **565** Restart test
- **566** Adaptive-promotion test
- **567** Operational-message test
- **568** Notification integration test
- **569** Component-storage ownership test
- **570** Localization scalability test
- **571** Localization public API review
- **572** Internal testapp localization fixture
- **573** Developer localization guide
- **574** Administrator localization guide
- **575** Application integration guide
- **576** Milestone 12 gate

**Milestone 12 gate:** Localization remains component-owned derived content with stable units, fingerprinted change detection, preferred/on-demand coordination, review protection, and notification integration.

# Milestone 13: Lifecycle, CLI, and Health Hardening

- **577** Lifecycle state model
- **578** Root process context
- **579** Lifetime bootstrap logger
- **580** Full logger activation
- **581** Startup-failure logger fallback
- **582** Deployment configuration load
- **583** Static configuration validation gate
- **584** Application composition invocation
- **585** Registration/freeze lifecycle gate
- **586** Required integration binding validation
- **587** Startup resource cleanup stack
- **588** Startup failure phase attribution
- **589** Dedicated startup DB connection
- **590** MariaDB compatibility gate
- **591** Migration lock phase
- **592** Migration execution phase
- **593** Migration resource cleanup
- **594** Database-managed runtime-state phase
- **595** Runtime snapshot construction gate
- **596** Foreground pool construction
- **597** Background pool construction
- **598** Pool failure cleanup
- **599** Runtime service construction
- **600** Control-plane startup
- **601** Recovery-before-readiness
- **602** UDS preflight
- **603** Existing-socket classification
- **604** Socket bind
- **605** Pre-readiness request gate
- **606** Health status types
- **607** Narrow liveness
- **608** Central readiness aggregator
- **609** MariaDB readiness
- **610** Required integration readiness
- **611** Optional integration degradation
- **612** Disabled subsystem state
- **613** Background administrative pause health
- **614** Health reason model
- **615** Minimal public health response
- **616** Detailed administrative health
- **617** Health transition logging
- **618** Readiness activation
- **619** Ready event
- **620** Reusable CLI root
- **621** serve
- **622** install
- **623** migrate
- **624** validate
- **625** version
- **626** Common CLI configuration path
- **627** CLI exit codes
- **628** CLI secret redaction
- **629** Signal registration
- **630** Repeated shutdown signal behavior
- **631** Readiness-off first
- **632** Stop accepting normal requests
- **633** Stop new background claims
- **634** Stop schedule generation
- **635** Stop LLM new submissions
- **636** Stop new event-driven work
- **637** Quiesce completeness assertion
- **638** Foreground request drain
- **639** Background-handler drain
- **640** LLM active-request drain
- **641** Drain deadline
- **642** Cancellation phase
- **643** Non-cooperative work handling
- **644** Hard shutdown bound
- **645** Database-close gate
- **646** Close background pool
- **647** Close foreground pool
- **648** DB-close failure logging
- **649** Listener final close
- **650** Socket filesystem cleanup
- **651** Socket cleanup failure
- **652** Stop control-plane loops
- **653** Runtime cache teardown
- **654** Full logger flush
- **655** Bootstrap final diagnostics
- **656** Final lifecycle state
- **657** Failure before DB
- **658** Migration-phase failure
- **659** Runtime-state failure
- **660** Pool-construction failure
- **661** Control-plane startup failure
- **662** Socket-bind failure
- **663** Startup readiness test
- **664** DB-loss transition test
- **665** Optional-integration outage test
- **666** Required-integration outage test
- **667** Disabled-capability health test
- **668** Detailed-health authorization test
- **669** Idle shutdown
- **670** Active foreground shutdown
- **671** Active job shutdown
- **672** Job exceeding drain deadline
- **673** Active workflow commit shutdown
- **674** Active message delivery shutdown
- **675** Active LLM shutdown
- **676** Socket replacement race test
- **677** serve integration test
- **678** install integration test
- **679** migrate integration test
- **680** validate nonmutation test
- **681** Common-config consistency test
- **682** Lifecycle events
- **683** Startup timing fields
- **684** Shutdown timing fields
- **685** Lifecycle API review
- **686** Health API review
- **687** CLI API review
- **688** Administrator lifecycle guide
- **689** CLI documentation
- **690** Health operations guide
- **691** Service-manager integration guidance
- **692** Milestone 13 gate

**Milestone 13 gate:** Canonical startup, readiness, CLI commands, bounded quiesce/drain/database-close/termination, socket cleanup, and logging shutdown ordering are proven.

# Milestone 14: Public API and Compatibility Freeze

- **693** Exported API inventory
- **694** Public package inventory
- **695** Internal package boundary audit
- **696** Remove accidental exports
- **697** Avoid configuration structs as accidental API
- **698** Avoid persistence structs in public API
- **699** Internal numeric ID audit
- **700** Secret-bearing API audit
- **701** Service-locator audit
- **702** Context-value dependency audit
- **703** Cross-component API audit
- **704** Stable-name namespace review
- **705** Stable-name type safety
- **706** Action API freeze review
- **707** Action registration ergonomics
- **708** Auth public model review
- **709** Authorization API freeze
- **710** Resource-authorization boundary
- **711** MariaDB public API review
- **712** Transaction API review
- **713** Repository guidance freeze
- **714** Job API freeze
- **715** Workflow API freeze
- **716** Schedule API freeze
- **717** Rendering API freeze
- **718** Widget API freeze
- **719** External package freeze
- **720** Safe-fetch API review
- **721** Messaging API freeze
- **722** Notification API freeze
- **723** LLM API freeze
- **724** Localization API freeze
- **725** Error taxonomy freeze
- **726** Error wrapping contract
- **727** Public error redaction
- **728** Semantic-versioning policy
- **729** What counts as a compatibility promise
- **730** Explicit non-contract internals
- **731** Deprecation policy
- **732** Harden internal/testapp
- **733** Public-only import enforcement
- **734** Consumer build test
- **735** Representative external component fixture
- **736** Go documentation audit
- **737** Package documentation
- **738** Lifecycle examples in docs
- **739** Configuration schema freeze review
- **740** Unknown-key policy confirmation
- **741** Runtime-setting compatibility
- **742** Persisted identifier audit
- **743** Migration upgrade fixture
- **744** No schema reinitialization dependency
- **745** Dangerous-helper audit
- **746** Trusted actor construction audit
- **747** Administrative API permission review
- **748** Goroutine ownership audit
- **749** Background API ownership audit
- **750** Public testing package review
- **751** Test-helper compatibility
- **752** Exported-symbol baseline
- **753** Dependency-direction check
- **754** Forbidden-pattern scan
- **755** Package-by-package human review
- **756** Cross-package naming review
- **757** Constructor consistency review
- **758** Interface ownership review
- **759** Public API reference
- **760** Compatibility guide
- **761** Component-author guide review
- **762** Application-author guide review
- **763** API compile baseline test
- **764** Persistence-name compatibility test
- **765** Config compatibility test
- **766** Behavioral compatibility suite
- **767** Public API freeze report
- **768** Architecture conformance review
- **769** Open-decision review
- **770** Codex authority audit
- **771** Milestone 14 gate

**Milestone 14 gate:** The exported Go surface, stable names, persisted/config compatibility, documentation, and consumer fixtures are explicitly frozen and reviewed.

# Milestone 15: Release Hardening

- **772** Release-candidate branch gate
- **773** Full unit-test suite
- **774** Race suite
- **775** Repeated race stress
- **776** Route fuzz suite
- **777** Configuration fuzz suite
- **778** Input-mapping fuzz suite
- **779** Token/parser fuzz suite
- **780** URL/safe-fetch fuzz suite
- **781** Cron-wrapper fuzz suite
- **782** Stored-content directive fuzz suite
- **783** Supported MariaDB baseline validation
- **784** Empty-database migration test
- **785** Prior-supported-version upgrade test
- **786** Interrupted-migration recovery
- **787** Concurrent migrator release test
- **788** Migration checksum drift test
- **789** Critical query inventory
- **790** Realistic database fixtures
- **791** Job claim EXPLAIN review
- **792** Schedule lookup EXPLAIN review
- **793** Session/token lookup EXPLAIN review
- **794** Authorization resolution EXPLAIN review
- **795** Messaging claim EXPLAIN review
- **796** LLM claim EXPLAIN review
- **797** Retention cleanup EXPLAIN review
- **798** Baseline methodology
- **799** Authentication lookup benchmark
- **800** Route dispatch benchmark
- **801** Job claim benchmark
- **802** Message claim benchmark
- **803** LLM FIFO benchmark
- **804** Snapshot-read benchmark
- **805** Regression threshold policy
- **806** Full crash-point matrix
- **807** Process-kill tests
- **808** Repeated crash-recovery stress
- **809** Retention inventory
- **810** Active-record protection
- **811** Bounded purge test
- **812** Content-vs-metadata retention test
- **813** LLM 90-day default test
- **814** Unverified-contact retention test
- **815** Foreground load with background work
- **816** Background backlog test
- **817** Schedule backlog recovery
- **818** External-provider outage storm
- **819** SMTP outage backlog
- **820** yllmd outage backlog
- **821** Memory-growth test
- **822** File-descriptor/socket leak test
- **823** Authentication security review
- **824** Session security review
- **825** Bearer-token security review
- **826** Authorization bypass review
- **827** CSRF/origin review
- **828** UDS/server security review
- **829** SSRF review
- **830** Template/directive security review
- **831** Secret-leak review
- **832** LLM trust-boundary review
- **833** Privileged administrative API review
- **834** Security-event review
- **835** Dependency vulnerability review
- **836** Security blocker policy
- **837** MariaDB outage test
- **838** Disk/logging failure test
- **839** SMTP failure test
- **840** yllmd failure test
- **841** External integration failure test
- **842** Corrupt runtime configuration test
- **843** Stable-name upgrade test
- **844** Public consumer upgrade test
- **845** Configuration upgrade test
- **846** Data-preservation upgrade test
- **847** User documentation completion
- **848** Administrator documentation completion
- **849** Developer documentation completion
- **850** Configuration reference
- **851** Database schema reference
- **852** CLI reference
- **853** Operational recovery guide
- **854** Final .codex directory audit
- **855** Full task graph synchronization
- **856** Completed-task state audit
- **857** Open-decision state audit
- **858** Architecture index verification
- **859** Repository layout audit
- **860** Generated/temporary artifact audit
- **861** License/security files
- **862** Dependency cleanup
- **863** Clean-clone build
- **864** Clean MariaDB install test
- **865** Offline testability
- **866** Deterministic test audit
- **867** Correlation propagation test
- **868** Logging-volume review
- **869** Operational event coverage
- **870** Subsystem acceptance checklist
- **871** Forbidden-architecture review
- **872** Security final sign-off
- **873** Performance final sign-off
- **874** Documentation final sign-off
- **875** Compatibility final sign-off
- **876** Release manifest
- **877** CHANGELOG entry
- **878** Release candidate tag
- **879** RC soak validation
- **880** Final release gate
- **881** Milestone 15 gate

**Milestone 15 gate:** All subsystem acceptance criteria, race/fuzz, upgrade, performance/query-plan, security, retention/load, documentation, and release checks pass.

<!-- ORIGINAL_SOURCE_END: 26_FULL_CODEX_TASK_GRAPH_M04_M15.md -->

