# yFramework Consolidated Source: 00_FOUNDATION_ARCHITECTURE

This consolidated project source preserves the original documents listed below verbatim. The wrapper headings and checksum manifest are consolidation metadata only.

## Included source documents

- `README_ARCHITECTURE_FREEZE.md` — SHA-256 `ed5eea29ad9760c5171e9d2e117f57dba288757fc5507038c3df6c86c2cc77c1`
- `00_PROJECT_CHARTER.md` — SHA-256 `98a0b1e2906f0a20a0f846e1f34f4c08e33f93b6cb26e96f02dd731bdad3e9c8`
- `01_ARCHITECTURE_OVERVIEW.md` — SHA-256 `e249f07e910db8d4864702e9f42c35244c72ff1fe690229b2d7e11cc81c7244a`
- `14_OPEN_DECISIONS.md` — SHA-256 `4dd39b3162e988398e5a3242ccb93448a71ab71ba0d76fd577b925c35daa06b6`
- `15_REPOSITORY_PACKAGES_DEPENDENCIES.md` — SHA-256 `187edf4d9636c5bf42947bdbaabeddb412ce609c5d66539fad37ce9e9fdfe816`

## Canonical consolidation clarifications

The embedded source documents below preserve their historical text, but the following newer project-level clarifications are authoritative where older wording differs:

- The development-planning source set is no longer prospective. Files `15` through `26` are part of the current authoritative planning set. The Milestone 4 through 15 task-ID/sequencing extension is embedded in `04_IMPLEMENTATION_CODEX.md`.
- The charter phrase "Example and reference applications" does not authorize an `examples/` tree or an application template inside yFramework. Reference applications, when created, live as separate consuming repositories/modules.
- Consuming applications continue to own `main`; yFramework remains an importable framework.

These clarifications supersede only the stale wording identified above and do not otherwise rewrite the embedded originals.

---

<!-- ORIGINAL_SOURCE_BEGIN: README_ARCHITECTURE_FREEZE.md -->

# Original Source File: `README_ARCHITECTURE_FREEZE.md`

# yFramework Architecture v1

The architecture-v1 source set is frozen for implementation planning. Files `00` through `14` contain the architecture and open-decision record.

Development planning source files follow the frozen architecture and are equally authoritative for implementation constraints:

- `15_REPOSITORY_PACKAGES_DEPENDENCIES.md`
- `16_DATABASE_SCHEMA_DESIGN.md`

The initial public repository and Go module are `github.com/BrandonYaniz/yframework`. Consuming applications are separate repositories/modules. yFramework contains no `examples/` directory and does not own application `main` packages.

Further development-planning files will be added as configuration schemas, migration implementation, testing, milestones, acceptance criteria, documentation, and `.codex` instructions are settled.

## Development planning additions

The development source set also includes:

- `23_DOCUMENTATION_OUTLINE.md`
- `24_CODEX_OPERATING_SPEC.md`
- `25_INITIAL_CODEX_TASK_GRAPH.md`

These documents define documentation ownership, Codex operating constraints, and the initial Milestone 1-3 implementation task graph.

<!-- ORIGINAL_SOURCE_END: README_ARCHITECTURE_FREEZE.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 00_PROJECT_CHARTER.md -->

# Original Source File: `00_PROJECT_CHARTER.md`

# yFramework Planning Project Charter

## Purpose

This ChatGPT Project exists to plan and refine **yFramework**, an importable, modular Go framework for building web, API, MCP, background-processing, and component-oriented applications.

The Project is not initially intended to generate the completed framework. Its immediate purpose is to produce a complete, internally consistent project skeleton that Codex can use to implement the framework.

The eventual output of this planning Project will include:

- A proposed repository directory structure
- Framework package skeletons
- Component-development skeletons
- Application-development skeletons
- User and administrator documentation structure
- LLM integration through yllmd
- Persistent localization infrastructure
- Notification infrastructure
- Framework messaging infrastructure with SMTP as the initial delivery channel
- External-integration runtime and adapter architecture
- Developer documentation
- Architectural decision records
- Database schema and migration specifications
- Testing requirements
- Example and reference applications
- A `.codex` directory containing Codex-specific instructions, plans, prompts, task files, and development state

## Primary Goal

Create enough precise documentation that Codex can implement yFramework without having to infer major architectural decisions.

The implementation instructions should explain:

- What the framework is responsible for
- What applications are responsible for
- What components are responsible for
- Which decisions are fixed
- Which decisions remain open
- How the framework subsystems interact
- What must be implemented first
- How completion should be tested

## Current Reference Projects

yCMS and the **Book Dragon Metadata API** are separate examples of what will be built on top of yFramework.

yCMS is a reusable CMS composition layer built from yFramework components. It should exercise:

- HTML rendering
- Templates and application overrides
- Dynamic routes
- Widgets and widget placement
- Hooks and events
- Replaceable component providers
- Typed cross-component APIs
- User authentication and authorization
- Content-oriented jobs and workflows
- Persistent LLM-backed localization
- Preferred and on-demand localized content
- Verified ready-notification enrollment

The Book Dragon Metadata API is a concrete API service. It should exercise:

- MariaDB migrations and repositories
- Dynamic routes
- Typed actions
- JSON rendering
- Application tokens
- Installation tokens
- User tokens
- Read-only application access
- Anonymous metadata-ingestion submissions
- Jobs and workflows
- CSV metadata imports
- Repository-owned caching
- UNIX-socket service operation
- Structured logging

These are reference projects, not framework packages or required application features. Additional web, API, MCP, and background-processing applications must be able to use yFramework without depending on either yCMS or Book Dragon domain logic.

A small web consumer for the Book Dragon service should follow so the same actions can be tested through both JSON and HTML renderers.

## Fixed Project Constraints

- Go is the implementation language.
- The framework is an importable Go module.
- Applications own `main`.
- MariaDB is the only supported relational database.
- Application services expose HTTP over UNIX domain sockets; applications own the socket path.
- Components own domain capabilities.
- Public routes are application configuration stored in MariaDB.
- The framework must not contain yCMS-specific or Book Dragon-specific domain logic.
- Framework behavior must be machine-observable through structured logs.
- Jobs and workflows run within the monolith service.
- The operating system local timezone is authoritative for schedules and rollover.
- The design should favor explicit dependencies and stable names over reflection-heavy or convention-heavy behavior.
- Components remain passive registrants in v1.
- Cross-component functionality uses injected typed contracts, not action-registry dispatch or direct table mutation.
- Executable background work uses one global FIFO lane with at most one active item in v1.
- Configuration contains no executable scripts or arbitrary expressions.
- yllmd is the framework-supported LLM broker. Components do not communicate with yllmd directly.
- yFramework owns a durable internal FIFO queue for outbound LLM requests and their application correlation.
- The framework LLM queue is separate from both the global background-work FIFO and yllmd's inference queue.
- yFramework provides persistent LLM-backed localization infrastructure.
- Primary-language domain content remains authoritative; localized variants are component-owned derived data.
- yFramework provides a generic notification layer with verified anonymous contact subscriptions.
- Components never own external transport or communicate directly with external systems. They depend on typed ports bound by the application to framework-managed or application-provided adapters.
- External adapters reuse framework-managed policies for credentials, TLS, timeouts, retries, retry budgets, rate limits, circuit breakers, health, tracing, logging, redaction, and metrics.
- Remote SMTP is the required v1 outbound messaging transport. SMS and other channels are future extensions behind the same messaging boundary.
- FreeBSD, Linux, and macOS are first-class deployment targets. Core runtime behavior remains portable Unix behavior; service-manager integration belongs to packaging.
- Human authentication is passkey-only in v1. Passwords, MFA layers, and account recovery are not supported. Multiple passkeys per account are supported.
- API authentication uses explicit opaque bearer-style yFramework tokens in v1; OAuth/OIDC protocol implementation is outside v1.
- The current MariaDB release available in the FreeBSD ports tree when a yFramework release is cut defines that yFramework release's minimum supported MariaDB version.

## How to Use These Documents

These documents capture the current architectural state. Future planning should:

1. Preserve settled decisions unless the user explicitly changes them.
2. Record new decisions in the appropriate specification.
3. Add unresolved architectural matters to `14_OPEN_DECISIONS.md`. Implementation details that Codex may decide should be recorded as implementation decisions rather than architectural open questions.
4. Avoid reintroducing previously rejected abstractions.
5. Keep examples clearly separated from framework requirements.
6. Treat yCMS and Book Dragon as reference projects, not as framework-owned functionality.

## Codex Decision Authority

Codex may make implementation decisions that do not alter documented public APIs, persisted-data semantics, security boundaries, ownership boundaries, lifecycle behavior, failure/recovery semantics, compatibility requirements, or externally observable framework behavior.

If implementation requires choosing among materially different architectural behaviors not resolved by these specifications, Codex records the matter in `14_OPEN_DECISIONS.md` and does not silently invent an architectural decision.

Codex may choose private helper types, package-private interfaces, internal algorithms, ordinary index details, test organization, unexported names, and other replaceable implementation details when those choices preserve the documented contracts.

## Development Planning Source Files

Development planning continues in numbered project source documents after the architecture specifications.

Current additions:

- `15_REPOSITORY_PACKAGES_DEPENDENCIES.md` records the public repository/module identity, package skeleton, dependency direction, import boundaries, versioning, and the decision to omit an examples directory.
- `16_DATABASE_SCHEMA_DESIGN.md` records the conceptual framework schema, performance and indexing rules, cache/database boundary, dual-ID strategy, retention model, high-write table design, and table-family ownership.

These documents are part of the project source of truth and constrain Codex implementation alongside files `00` through `14`.

<!-- ORIGINAL_SOURCE_END: 00_PROJECT_CHARTER.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 01_ARCHITECTURE_OVERVIEW.md -->

# Original Source File: `01_ARCHITECTURE_OVERVIEW.md`

# Architecture Overview

## Framework Form

yFramework is an importable Go module. It is not a forkable application template and does not own the application executable.

An application:

- Owns `main`
- Imports yFramework
- Loads deployment configuration
- Registers components
- Runs migrations
- Starts the UNIX-socket server
- Starts internal jobs, workflows, and scheduler services
- Chooses enabled renderers and adapters

The framework provides shared infrastructure. Components provide domain capabilities.

## Reference Projects

yCMS and the Book Dragon Metadata API are independent examples built on top of yFramework.

- yCMS is a reusable CMS composition layer that selects and wires content-oriented components while leaving final application composition and `main` under application control.
- The Book Dragon Metadata API is a concrete API service that uses framework capabilities for metadata lookup, staged ingestion, authentication, jobs, workflows, and JSON delivery.

Neither project defines the framework boundary. yFramework must remain usable by applications that import neither yCMS nor Book Dragon packages.

## Responsibility Boundaries

### Framework Owns

- Application lifecycle primitives
- Configuration loading and validation
- MariaDB connection support
- Migration execution and tracking
- Dynamic route loading and validation
- Typed action registration and invocation
- Rendering adapters
- Authentication and authorization infrastructure
- Token validation and session management
- Hooks and events infrastructure
- Cache-manager infrastructure
- Widget infrastructure
- Job queues
- Workflow coordination
- Internal scheduling
- Framework LLM service and yllmd transport integration
- Durable outbound LLM request queue, correlation, cancellation, and result tracking
- Structured logging
- Request and trace contexts
- Validation and shared errors
- Testing helpers
- CLI command infrastructure
- Localization coordination
- Notification infrastructure
- Channel-independent outbound messaging infrastructure
- Shared external-integration runtime for transport policy, credentials, resilience, health, and observability
- Safe outbound API transport and public-web fetch primitives
- LLM throughput estimation and request-time forecasting

### Application Owns

- `main`
- Enabled components
- Deployment configuration
- Public route configuration
- Global authorization policy
- UNIX-socket path
- Database credentials
- Logging paths and retention
- Worker and queue resource allocation
- Renderer availability
- Component wiring
- Component implementation selection
- Named component bindings and component-set composition
- Application-specific templates and overrides
- Whether LLM functionality is enabled
- yllmd socket and deployment connection settings
- Application mapping from component LLM purposes to yllmd groups and profiles
- External adapter selection and binding
- External provider endpoints, credentials, and provider-specific deployment configuration
- Required-versus-optional integration policy

### Component Owns

- Domain actions
- Domain types
- Business rules
- MariaDB repositories and queries
- Component migrations
- Cache usage and invalidation decisions
- Hooks and event definitions
- Event subscribers
- Widgets
- Job definitions
- Workflow definitions
- Transitional data
- Final workflow commit behavior
- Optional templates and renderer formatters
- Component-level authorization rules
- LLM prompt construction, domain context, expected output, result validation, and business interpretation
- Stable component-qualified LLM purpose definitions when specialized routing is required
- Domain-shaped external-service ports describing what the component needs, without owning provider transport
- Domain interpretation of external-service results

Components do not own the application lifecycle, global router, global database pool, public paths, global logger, service sockets, application-wide authentication policy, direct yllmd transport, SMTP transport, provider credentials, raw outbound HTTP transport, or direct communication with any system outside the application boundary.

Components are passive registrants in v1. They expose a permanent identifier and one registration method. The application constructs them with explicit typed dependencies, and the framework owns validation, migrations, readiness, background coordination, and shutdown.

## Major Request Path

```text
Request
  -> protocol adapter
  -> request context and identity
  -> route match
  -> route-handler selection
  -> input mapping and validation
  -> typed action
  -> typed result or shared error
  -> renderer
  -> response
```

## Background Work Path

```text
Trigger or schedule
  -> durable job or workflow record
  -> internal scheduler/poller
  -> global FIFO eligibility
  -> worker claim and lease
  -> component-owned job
  -> durable result
  -> workflow transition
  -> final component commit step
  -> cleanup
```


## LLM Integration Path

yFramework integrates with the separately operated `yllmd` project. Components do not open yllmd sockets, encode its protocol, manage its queue, or select concrete models directly. They consume typed framework LLM capabilities supplied through normal dependency injection.

```text
component or action/job
  -> yFramework LLM API
  -> purpose-to-target resolution
  -> durable yFramework outbound LLM FIFO
  -> yllmd Unix-socket client
  -> yllmd FIFO inference queue
  -> resolved model/provider
```

The two framework queues have different responsibilities. The global background FIFO coordinates durable application work. The LLM outbound FIFO coordinates inference requests from foreground requests, jobs, workflows, and system work without forcing foreground LLM use through the background-work lane. yllmd retains its own inference queue for model-resource scheduling.

The framework LLM layer owns request identity, ownership and trace correlation, durable pre-submission state, cancellation propagation, normalized errors, response delivery, and structured observability. yllmd owns model routing, groups and profiles, concrete model selection, operational fallbacks, model lifecycle, and provider evolution.

Components own their prompts, domain context, expected result shape, result validation, and the business decision made from the result. Application-specific prompting does not belong in yllmd or in a framework-global prompt library.

## Localization Path

yFramework treats localization as persistent derived content, not as repeated request-time translation. Components define localizable resources and stable localization units, own localized domain storage, and identify when source content becomes localization-eligible.

```text
component source content
  -> localization eligibility/change detection
  -> preferred or on-demand localization request
  -> yFramework LLM outbound queue
  -> yllmd
  -> validated localized units
  -> component-owned localized data
```

Applications define a primary locale, supported locales, manually preferred locales, adaptive-preferred policy, and presentation behavior. Preferred locales are proactively generated during idle capacity. Supported non-preferred locales are generated on first demand. Adaptive promotion may move a supported locale into preferred behavior when configured usage thresholds are reached; v1 does not automatically demote it. Promotion queues existing eligible content for idle-time backfill.

Framework operational localization messages, including localization-in-progress and content-unavailable messages, are generated and stored for every supported locale during idle capacity so user-facing status does not itself require an on-demand LLM call.

## Notification Path

Components use a framework notification service rather than sending email or SMS directly. Notification delivery is durable and separate from the transaction that produced the underlying durable fact.

```text
durable completion event
  -> notification subscription match
  -> durable notification delivery
  -> configured channel provider
```

Anonymous contacts must be verified before receiving completion notifications. Email is the initial channel; SMS is a later channel. Verification proves control of a destination but does not grant consent for unrelated communication.

## Module Organization

The first implementation should remain one Go module with focused packages, such as:

```text
action
auth
cache
config
errors
events
hooks
job
llm
localization
logging
mariadb
migrations
notification
render
route
server
transaction
validation
widget
workflow
```

The framework should avoid one giant root package and avoid splitting into many separately versioned modules before a demonstrated need exists.

## Design Principles

- Stable names at dynamic boundaries
- Strong typing inside component code
- Type erasure only at runtime registries
- Explicit dependencies
- No service locator
- No global mutable component state
- MariaDB as durable source of truth
- In-memory registries and routers built from validated startup state
- Framework mechanics separated from domain meaning
- Safe defaults that can be overridden intentionally
- Application-owned provider selection before registration
- Immutable registries and runtime snapshots after validation

## External Integration Boundary

Components never own external transport. A component expresses an external need through an explicitly injected, domain-shaped typed port. Application composition selects and binds an adapter. The adapter performs external communication using framework-managed integration infrastructure.

```text
component
  -> typed domain port
  -> application-selected adapter
  -> yFramework external-integration runtime
  -> external system
```

The framework must not collapse all external communication into a generic service locator such as `external.Call`. Purpose-specific interfaces remain preferred, such as LLM generation, messaging, metadata lookup, or public-web fetching.

The shared integration runtime provides reusable mechanics where applicable:

- Credential resolution and redaction
- TLS and connection policy
- Timeouts and cancellation
- Bounded retries with backoff and jitter
- Aggregate retry budgets
- Rate limiting
- Circuit breakers
- Required/optional health and readiness reporting
- Request and trace correlation
- Structured logging and metrics
- Provider-error normalization
- Test doubles and adapter contract-test support

Provider-specific or domain-specific adapters do not automatically belong in yFramework. Application or integration packages may implement adapters using framework transport primitives. Components see only the typed port.

Outbound work is not universally queued. Synchronous reads may execute synchronously when required by the request path. Durable side effects such as messages use durable delivery queues. LLM work uses the dedicated LLM queue. The operation semantics determine the coordination model.

## Canonical Process Lifecycle

### Startup

The canonical startup sequence is:

```text
bootstrap stderr logger
  -> load and statically validate deployment configuration
  -> initialize full structured logging
  -> compose application and construct selected components/adapters
  -> register components and validate/freeze static registries
  -> open a startup MariaDB connection
  -> verify MariaDB version/capabilities
  -> acquire the migration lock
  -> verify and apply framework/component migrations
  -> release the migration lock
  -> load and validate database-managed runtime state
  -> create runtime foreground/background MariaDB pools
  -> construct framework runtime services
  -> start queue/scheduler/control-plane coordination
  -> bind the UNIX socket and establish readiness
  -> accept normal requests
```

Registration is declarative and database-free. Static application validation completes before database mutation.

The bootstrap logger writes to stderr and remains available for the entire process lifetime. The configured structured logger starts as early as practical and remains active through nearly all shutdown work. FreeBSD, Linux, and macOS service packaging is responsible for capturing or redirecting stderr; the framework does not depend on rc.d, systemd, launchd, syslog, journald, or platform-specific log paths.

Normal startup automatically applies pending registered migrations under a MariaDB-backed migration lock. Applied migrations are immutable and checksum-validated. Successfully applied migrations are not automatically rolled back if a later migration fails. Production migrations should use forward-compatible expand/migrate/contract patterns where practical.

### Readiness and Health

Liveness, readiness, and degraded capability are distinct. MariaDB is readiness-critical. Optional integration failure produces capability-specific degradation rather than process failure. Required integration policy is application-owned. Public health output is minimal; detailed subsystem diagnostics require administrative or local operational access.

### Shutdown

Shutdown uses four top-level phases:

```text
QUIESCE
  -> DRAIN
  -> CLOSE DATABASE
  -> TERMINATE
```

Quiesce disables every listener, scheduler, trigger, queue claimant, retry activation path, workflow advancement mechanism, event-driven dispatcher, and other source capable of initiating new work. Durable queued records remain intact.

Drain permits already-running foreground and durable background work to finish within bounded deadlines while preserving leases, heartbeats, and state needed to record correct completion. The framework never clears a lease while the old handler may still execute.

Once durable work can no longer start or finish, MariaDB is closed. The UNIX listener is closed and its filesystem socket entry is removed. Noncritical runtime machinery does not require elaborate graceful teardown and may terminate with the process.

The structured logger remains available as long as practical, then flushes and closes near the end. The lifetime bootstrap logger reports final logger shutdown, socket-cleanup failures, and final process termination to stderr.

## UNIX Socket Server Contract

V1 serves standard HTTP over an application-configured Unix domain socket. The framework does not define a custom framing protocol. Startup never blindly unlinks an existing socket: a socket confirmed stale may be removed, while a live or ambiguous existing endpoint prevents startup.

Socket ownership, mode, and parent-directory safety are validated and remain deployment configuration. Forwarded headers are not trusted implicitly; applications explicitly configure trusted proxy behavior and public hosts/origins. WebAuthn RP ID and approved origin configuration are explicit and are never inferred solely from the backend socket or untrusted forwarded headers.

The server applies bounded header/body/time policies and propagates request cancellation through `context.Context`. Binding the socket does not by itself make the application ready; normal request handling is gated until final readiness is established.

## CLI Boundary

Applications own `main` and the executable. yFramework provides reusable CLI infrastructure and lifecycle capabilities. V1 standard lifecycle commands are `serve`, `install`, `migrate`, and `validate`; `version` is recommended.

All lifecycle commands reuse the same application composition, configuration, validation, and migration infrastructure. `serve` remains a foreground process and never daemonizes itself. `validate` is non-mutating. Platform service supervision belongs to FreeBSD, Linux, or macOS packaging.

<!-- ORIGINAL_SOURCE_END: 01_ARCHITECTURE_OVERVIEW.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 14_OPEN_DECISIONS.md -->

# Original Source File: `14_OPEN_DECISIONS.md`

# Open Decisions

## Purpose

This file contains architectural decisions that must not be silently invented during implementation. When an item is settled, update the owning specification and remove it from this file.

## Architecture v1 Status

There are currently no known architecture-blocking decisions remaining for the initial development plan.

The following are intentionally delegated implementation decisions rather than open architecture questions:

- Exact exported Go names and helper structs where the public semantic contract is already specified.
- Exact MariaDB column types and ordinary index details that do not change persisted semantics.
- Choice of WebAuthn, cron parsing, and HTTP-routing libraries, provided they satisfy the documented framework behavior.
- Internal registry, snapshot, cache, and lifecycle helper implementations.
- Detailed cleanup cadences and administrative method names.
- Exact CLI exit-code numbering beyond stable success/non-success behavior.

If implementation reveals a choice that would alter public APIs, persisted-data semantics, security boundaries, component/application/framework ownership, lifecycle behavior, failure/recovery behavior, compatibility, or other externally observable framework behavior, add it here before proceeding.

## Development Planning Status

The repository/package structure and conceptual framework database model are no longer open architecture decisions. They are recorded in `15_REPOSITORY_PACKAGES_DEPENDENCIES.md` and `16_DATABASE_SCHEMA_DESIGN.md`.

Concrete DDL details such as exact bounded string sizes, UUID encoding, binary hash sizes, index names/order verified through `EXPLAIN`, and non-LLM retention defaults remain implementation engineering within the settled schema behavior. They become open architectural decisions only if a proposed implementation would change persistence semantics, security, subsystem ownership, compatibility, or failure/recovery behavior.

<!-- ORIGINAL_SOURCE_END: 14_OPEN_DECISIONS.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 15_REPOSITORY_PACKAGES_DEPENDENCIES.md -->

# Original Source File: `15_REPOSITORY_PACKAGES_DEPENDENCIES.md`

# Repository, Packages, and Dependency Boundaries

## Status

This document records the settled development structure for the initial public yFramework repository and the package dependency rules Codex must preserve.

## Public Repository and Module

### Settled requirement

The initial public repository is:

```text
github.com/BrandonYaniz/yframework
```

The initial Go module path is:

```go
module github.com/BrandonYaniz/yframework
```

The repository may later be transferred to a company organization. A transfer that changes the Go module path is a deliberate module-path migration for consuming applications and is not assumed to be transparent.

Applications such as yCMS and Book Dragon are separate repositories and separate Go modules. They consume tagged yFramework module releases through ordinary Go module dependencies.

A consuming application does not automatically receive newer yFramework code merely because it is recompiled. Dependency upgrades are explicit so builds remain reproducible. During local development, Go workspaces or local replacement directives may be used without changing the published dependency model.

## No Framework Executable Application

### Settled requirement

Applications own `main`. yFramework is an importable framework and does not provide an application template executable or `cmd/yframework` runtime.

A separate developer utility may be created in the future only through an explicit design decision.

## No Examples Directory

### Settled requirement

The yFramework repository does not contain an `examples/` tree. Real applications and separate reference repositories are used to exercise and demonstrate the framework.

Application-specific code, including yCMS and Book Dragon domain logic, must not enter framework packages.

## Repository Skeleton

### Settled development structure

```text
yframework/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── SECURITY.md
├── CHANGELOG.md
│
├── action/
├── application/
├── auth/
├── authorization/
├── cache/
├── cli/
├── component/
├── config/
├── errors/
├── events/
├── external/
├── health/
├── hooks/
├── job/
├── lifecycle/
├── llm/
├── localization/
├── logging/
├── mariadb/
├── messaging/
├── migrations/
├── notification/
├── render/
├── route/
├── schedule/
├── server/
├── transaction/
├── validation/
├── widget/
├── workflow/
│
├── internal/
│   ├── registry/
│   ├── runtime/
│   └── testapp/
│
├── testing/
│
├── docs/
│   ├── architecture/
│   ├── developer/
│   ├── administrator/
│   ├── user/
│   ├── schemas/
│   └── adr/
│
└── .codex/
```

Empty directories are not created merely to satisfy the skeleton. Packages appear when implementation requires them.

## Package Purpose

### Settled requirement

The package names describe focused responsibilities. Do not introduce broad dumping-ground packages such as `common`, `utils`, `helpers`, `services`, `manager`, or a generalized multi-database `database` package.

`mariadb` is intentionally specific because MariaDB is the only relational database supported by yFramework.

The initial package responsibilities are:

- `action`: typed action contracts and execution metadata.
- `application`: application composition and validated construction; does not own `main`.
- `auth`: passkeys, browser sessions, bearer-token authentication, and identity establishment.
- `authorization`: permissions, roles, direct grants, delegated client/install capabilities, recent-auth policy, and authorization decisions.
- `cache`: central typed cache manager and cache primitives.
- `cli`: reusable CLI infrastructure used by application executables.
- `component`: component contracts, registration, and declared dependencies.
- `config`: typed deployment/bootstrap configuration and static validation.
- `errors`: stable framework error model.
- `events`: synchronous in-process event contracts.
- `external`: normalized outbound integration transport, resilience, safe fetch, and credential boundaries.
- `health`: liveness, readiness, degradation, and health-reporting primitives.
- `hooks`: synchronous lifecycle/operation hook contracts.
- `job`: durable ordinary background work and execution coordination.
- `lifecycle`: bootstrap, startup, readiness transition, quiesce, drain, and termination orchestration.
- `llm`: yFramework LLM request service, durable request state, and yllmd transport integration.
- `localization`: framework localization coordination and generated-localization workflow support.
- `logging`: bootstrap/full structured logging and logging contracts.
- `mariadb`: MariaDB connection pools, capability/version verification, and database-specific infrastructure.
- `messaging`: durable outbound message delivery and channel/provider attempt state.
- `migrations`: migration registry, checksums, locking, ordering, execution, and tracking.
- `notification`: verified contacts, subscriptions, and normalized notification intent.
- `render`: renderer contracts, response rendering, templates/layout integration.
- `route`: route definitions, matching, access gates, permissions, and declarative input mapping.
- `schedule`: once/cron/interval schedule definitions and due-occurrence calculation.
- `server`: HTTP-over-UNIX-socket server adapter and request lifecycle integration.
- `transaction`: database transaction/executor helpers.
- `validation`: reusable structural/input validation.
- `widget`: widget registration, placement, and resolved placement snapshots.
- `workflow`: durable workflow coordination and phase/step state.

## Public and Internal Code

### Settled requirement

Public root packages expose only supported framework contracts. Replaceable implementation machinery belongs under `internal/` or a package-local `internal/` tree where practical.

`internal/registry` may provide private generic registry machinery. Public packages expose typed registrars and frozen views rather than the generic registry implementation.

`internal/runtime` may depend broadly because it is private composition machinery. Application and component code must not import it.

The public `testing` package provides reusable test helpers for component/application authors. Internal framework tests may additionally use `internal/testapp`.

## Dependency Layers

### Settled development structure

```text
Level 1: Foundation
  errors
  validation
  logging
  config

Level 2: Core infrastructure
  mariadb
  transaction
  migrations
  cache
  external
  health

Level 3: Framework contracts
  component
  action
  auth
  authorization
  events
  hooks
  render

Level 4: Coordinated subsystems
  route
  widget
  job
  workflow
  schedule
  llm
  localization
  notification
  messaging
  server

Level 5: Composition
  application
  lifecycle
  cli
```

Dependencies normally flow from higher levels toward lower levels. This is an allowed-direction model, not a requirement that a package import every lower layer.

## Import and Dependency Rules

### Settled requirement

A lower-level package must not import a higher-level coordinating package merely to gain access to a service.

Examples of prohibited dependencies include:

```text
mariadb -> application
cache -> component
logging -> route
external -> messaging
component implementation -> application runtime
component implementation -> lifecycle
component implementation -> server
```

The application composition layer supplies explicit typed dependencies instead.

### Settled requirement

Subsystem packages must be designed to avoid cycles. Registration contracts should be narrow and directional. For example, component registration may consume `action.Registrar`, `job.Registrar`, or similar interfaces while those subsystem packages must not import concrete component implementations.

### Settled requirement

Cross-component access never occurs through a service locator or by looking up another component's actions in a global registry. A component that requires another component capability receives a typed interface explicitly wired by the application.

## Shared Types

### Recommended default

Do not create generic `common`, `shared`, or `types` packages merely to break import cycles. Shared identifier/correlation types should live with the lowest package that semantically owns them. A narrowly named neutral package may be added only when a genuine cross-subsystem abstraction emerges.

## Component Repository/Layout Guidance

### Recommended default

A component may use a structure such as:

```text
component-name/
├── component.go
├── permissions.go
├── config.go
├── actions/
├── repositories/
├── jobs/
├── workflows/
├── events/
├── hooks/
├── widgets/
├── templates/
├── migrations/
└── internal/
```

Directories are created only when the component needs them. The layout is guidance, not a requirement that every component contain every directory.

## Application Layout Guidance

### Recommended default

A consuming application may use a structure such as:

```text
myapp/
├── go.mod
├── go.sum
├── cmd/
│   └── myapp/
│       └── main.go
├── app/
│   ├── application.go
│   ├── components.go
│   ├── adapters.go
│   └── config.go
├── components/
├── integrations/
├── templates/
├── migrations/
├── config/
└── docs/
```

`main.go` remains thin. Composition belongs in application code rather than becoming an unstructured `main` function.

## Versioning and Compatibility

### Settled requirement

Published applications consume tagged yFramework Go module versions. Public Go APIs become compatibility commitments once v1 is released. Internal implementation may change without compatibility guarantees where Go's `internal` rules prevent external dependence.

Semantic versioning is the expected release model. Breaking public API changes require the corresponding deliberate major-version/module-path treatment.

## Codex Enforcement

### Settled requirement

Codex must preserve these package and import boundaries. If implementation pressure requires a dependency that violates this document, Codex records the conflict as an open architectural decision rather than introducing the dependency silently.

<!-- ORIGINAL_SOURCE_END: 15_REPOSITORY_PACKAGES_DEPENDENCIES.md -->
