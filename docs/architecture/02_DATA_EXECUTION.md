# yFramework Consolidated Source: 02_DATA_EXECUTION

This consolidated project source preserves the original documents listed below verbatim. The wrapper headings and checksum manifest are consolidation metadata only.

## Included source documents

- `05_DATA_MIGRATIONS_CACHE_TRANSACTIONS.md` — SHA-256 `9a82b8156192c739c718ea0b24759265f6befcc134c234d740e4ee6ce8e87b57`
- `06_JOBS_WORKFLOWS_SCHEDULER.md` — SHA-256 `ccc18f9678f34d4a74dbc977dd3aefe195b9489d8338edd0823707aa4959d815`
- `16_DATABASE_SCHEMA_DESIGN.md` — SHA-256 `ab73e7b7200b195f3629ec617b36a998c2e1ea5f309397e67725a3e89e74a6d6`
- `18_MIGRATION_IMPLEMENTATION_CONTRACT.md` — SHA-256 `71c0262a78e71376585ad76a87e4989379a09138ac282e3d296e9b8631866aa1`

## Canonical consolidation clarifications

The embedded source documents below preserve their historical text, but these newer settled decisions are authoritative where older wording differs:

- V1 has one global ordinary executable FIFO lane and **no named ordinary execution queues**. There is no ordinary `queue_name` scheduling dimension and no execution-priority field. Older wording that queue names may remain for administration, pause controls, or diagnostics is superseded. Administrative pause/resume applies to the ordinary background lane as a whole. The dedicated LLM FIFO remains a separate, explicit subsystem queue.
- Deployment configuration is primarily one strict TOML file. CLI configuration is intentionally narrow, primarily selecting the config path/bootstrap inputs. Environment variables do **not** automatically override arbitrary TOML fields. Any environment-backed value must be an explicitly approved/configured bootstrap or secret mechanism rather than a universal precedence layer. For registered dynamic settings, the settled precedence remains compiled default -> permitted TOML override -> MariaDB runtime override. Deployment-only settings have no MariaDB override.

These clarifications supersede only the stale wording identified above and do not otherwise rewrite the embedded originals.

---

<!-- ORIGINAL_SOURCE_BEGIN: 05_DATA_MIGRATIONS_CACHE_TRANSACTIONS.md -->

# Original Source File: `05_DATA_MIGRATIONS_CACHE_TRANSACTIONS.md`

# Data, Migrations, Cache, Transactions, and Errors

## MariaDB

MariaDB is the only supported relational database.

The framework should not introduce a generalized multi-database abstraction.

Repositories may use an executor interface compatible with both:

- `*sql.DB`
- `*sql.Tx`

## Component-Owned Schema

Each component owns:

- Tables
- Queries
- Repositories
- Migration files
- Indexes
- Component seed data

The framework owns migration ordering, execution, checksums, and tracking.

## Migration Tracking

Use stable component-qualified migration IDs:

```text
article:001_create_articles
metadata:004_add_edition_identifier
```

Suggested tracking fields:

```text
component
migration
checksum
applied_at
```

Applied migrations must not be edited. New changes require new migrations.

Component dependencies determine migration ordering, then component migrations run sequentially.

## Transaction Boundaries

Actions and component jobs own transaction boundaries.

Required consistency work happens within a transaction.

Hooks that may block the operation run before commit.

After successful commit:

- Invalidate affected caches
- Publish domain events
- Run after-commit callbacks

Do not publish success events when commit fails.

Nested transactions are not required initially. Called actions or services should share an existing transaction context when necessary.

## Shared Error Model

The framework error should include:

- Stable code
- Public message
- HTTP status or protocol mapping
- Internal cause
- Field errors
- Retryable indicator

Internal causes are logged, not exposed.

Validation errors should support per-field details.

## Repository-Owned Caching

Repositories decide:

- What to cache
- Cache keys
- Cache values
- Loading behavior
- Invalidation behavior

The framework provides a central cache manager that manages multiple isolated typed caches.

Avoid one giant untyped cache.

## Cache Timers

Canonical terms:

- `idle_ttl`: reset after a successful cache hit
- `absolute_ttl`: maximum age from load, never reset
- `negative_ttl`: optional caching of missing results
- `cleanup_interval`
- `max_entries`

Entries track:

- Loaded time
- Last access
- Value state

## Cache Behavior

- Lazy expiration on access
- Central periodic cleanup
- Capacity-based LRU eviction
- Per-key concurrent-load suppression
- Typed cache keys
- Direct invalidation after successful writes
- Broad list/search invalidation is acceptable initially

Critical correctness must not depend on optional events.

Avoid authorization leaks by caching neutral data or including visibility scope in cache keys.

Widget placement caches store resolved placement decisions, not rendered HTML.

## Cache Purge Modes

The cache manager supports three application-wide purge modes:

```text
Normal -> Reduced -> Hold
```

The modes change when an ordinary purge request becomes eligible for execution:

- `Normal`: purge immediately.
- `Reduced`: add the configured reduced-mode purge TTL to the request time.
- `Hold`: add the configured hold-mode purge TTL to the request time. This TTL is substantially longer than the reduced-mode TTL.

A delayed purge is tracked cache-manager state for the life of the process. Each request records its target and `purge_after` time. The periodic cache manager executes eligible purges even if the application has since returned to a less restrictive mode. A process restart clears an in-memory cache, so pending purge records do not require MariaDB persistence in v1.

Mode changes affect new purge requests. They do not extend, shorten, or cancel purge requests already recorded. Repeated requests for the same target retain the earliest `purge_after` time so continued writes cannot postpone invalidation indefinitely.

Purge delay is distinct from entry expiration. `idle_ttl`, `absolute_ttl`, and `negative_ttl` still govern normal entry lifetime. A pending purge may remove an entry earlier than those timers.

Delayed purge is permitted only for caches explicitly declared safe for bounded staleness. A hard purge bypasses the current mode and executes immediately. Hard purge is required for:

- Authorization and permission changes
- Token or session revocation
- Ownership or visibility changes
- Security-sensitive state
- Writes whose contract requires immediate read-after-write consistency
- Any cache not explicitly declared safe for delayed purge

The application configures:

- The initial purge mode
- The reduced-mode purge TTL
- The hold-mode purge TTL
- Which caches permit delayed purge

The framework validates that the reduced-mode TTL is non-negative and the hold-mode TTL is greater than the reduced-mode TTL. Administrative controls may change the active mode at runtime and must produce audit and structured log events.

## Configuration Model

Configuration is separated into four categories:

1. Bootstrap configuration needed before the full application can be built
2. Deployment configuration owned by the application
3. Sensitive values, which remain sensitive even when stored inline
4. MariaDB-managed runtime settings

Settings required to connect to MariaDB are never stored only in MariaDB.

Use Go for:

- Components
- Actions
- Hooks
- Events
- Repositories
- Widgets
- Business logic
- Workflow definitions
- Job definitions

Use TOML for:

- Service name
- Socket path
- Database connection
- Log settings
- Global cache defaults
- Initial cache purge mode
- Reduced-mode and hold-mode purge TTLs
- Worker and poller limits
- Runtime deployment options
- yllmd enablement, Unix-socket path, and connection-level timeouts
- MariaDB credentials, including the password in v1

Use MariaDB for:

- Routes
- Route mappings
- Access policies
- Widget instances and placements
- Tokens
- Schedules
- Simple registered application-wide runtime-setting overrides

Routes, route mappings, widget placement, tokens, and schedules retain dedicated schemas. They are not stored in the generic runtime-settings table.

Complex component configuration, provider accounts, domain rules, related records, per-user settings, and operational or transitional state belong in component-owned tables.


## LLM Configuration Boundary

yllmd is the supported LLM broker for yFramework. Connection settings required before runtime configuration is available belong in deployment TOML, including whether the LLM service is enabled, the yllmd Unix-socket path, and connection-level timeout settings.

Components may declare stable component-qualified LLM purposes such as:

```text
metadata.classify
article.summarize
```

An application maps those purposes to yllmd routing targets. A mapping may specify a yllmd group and optional profile. Components therefore express the semantic purpose of a request without dictating the application's yllmd route names or concrete model configuration. A generic request may intentionally omit a specialized purpose and use the application or yllmd default route.

Purpose-to-target mappings are non-executable configuration and may be eligible for MariaDB-managed dynamic settings after the application registers and validates the allowed purpose definitions. The yllmd socket path remains deployment configuration. Concrete model definitions, provider credentials, model lifecycle, and yllmd routing fallback configuration remain owned by yllmd and are not duplicated in yFramework MariaDB tables.

Prompt text and model output are sensitive application data. They are not written to logs by default. Durable LLM queue records retain only as much request content as is necessary for reliable delivery and recovery, and completed request content follows a bounded retention or cleanup policy. Large or domain-sensitive inputs should be represented by references to component-owned durable staging data when practical.

## Deployment Configuration

TOML is the recommended deployment format.

For v1, the general TOML configuration may contain the MariaDB password in clear text. The installer creates the file with restrictive filesystem permissions, redacts the password from all normal output, and warns that backups or copies of the configuration also contain the credential.

Separate secret files, environment-backed secrets, and external secret stores are optional future hardening. Their absence does not block v1.

Recommended precedence, from highest to lowest:

1. Explicit command-line bootstrap override
2. Approved environment override
3. Deployment TOML
4. Registered default

Unknown keys and invalid values prevent installation or startup. MariaDB runtime settings do not silently override deployment settings unless a registered definition explicitly declares the key dynamic and database-managed.

The operating system timezone is authoritative.

## Component Configuration Definitions

Components register definitions, not deployment-specific or secret values. A definition declares:

- Qualified setting name
- Go value type
- Default
- Validation
- Whether the setting is dynamic
- Whether the setting is sensitive
- Reload behavior

Keys are qualified by component instance name. Component implementation selection occurs before registration and remains application-owned.

Dynamic configuration may reference only immutable registered capabilities. MariaDB configuration never contains executable scripts or arbitrary expressions.

## Runtime Settings

Simple dynamic settings use one framework-owned table:

```text
framework_runtime_settings
  setting_name primary key
  value_json
  revision
  updated_at
  updated_by_type
  updated_by_id
```

The registered Go definition remains authoritative. MariaDB stores only an override. A missing row uses the registered default, and deleting a row resets the setting to that default.

Values use canonical JSON so registered types survive storage. Revision-based optimistic concurrency is the recommended default. V1 runtime settings are application-wide only.

A runtime-setting change:

1. Authenticates and authorizes the actor.
2. Resolves the registered definition.
3. Parses and validates the proposed value.
4. Checks the current revision.
5. Persists the canonical value transactionally.
6. Records an audit event.
7. Publishes an internal configuration-change event.
8. Rebuilds the affected immutable runtime snapshot.
9. Atomically replaces the old snapshot.

If rebuilding fails, the current working snapshot remains active. Existing requests may finish with the old snapshot while new requests use the replacement.

The exact Go definition types, MariaDB column types and indexes, and snapshot APIs remain implementation details.

## Localization and Notification Configuration

Applications define a primary locale and the locales they support. A supported locale may also be manually preferred or adaptively preferred. Preferred locales are proactively localized during idle capacity. Supported non-preferred locales are localized on demand.

Adaptive preference uses aggregate meaningful localization demand rather than raw page hits. Promotion thresholds and the observation policy are configurable. V1 may automatically promote an eligible supported locale but does not automatically demote it. Removing preferred status is administrative. When a locale becomes preferred, existing eligible content is queued for idle-time backfill.

Localization coordination records may be framework-owned, but localized domain content remains in component-owned schema. Each localized unit needs enough provenance to determine at least:

```text
source_fingerprint
source_revision
status
human_modified
locked
review_required
generated_at
```

Only changed units become stale and eligible for regeneration. Human-modified or locked stale units require administrator review and are not automatically overwritten.

LLM request and result payloads are retained by default for 90 days. The retention period is configurable. Records flagged for administrator review or quarantine are exempt from normal automatic purge until the review condition is cleared.

Anonymous notification contacts and notification subscriptions are framework-owned operational data. Unverified contacts are retained for 30 days by default. They cannot receive requested completion notifications until verified and are purged after the verification window expires if verification never occurs. Verification tokens are stored only as hashes. Repeated verification sends are rate-limited or cooldown-controlled.

A verified contact remains purpose-limited. Verification proves control of the destination and does not create a user account, newsletter subscription, marketing consent, or permission for unrelated messaging.

## External Integration Configuration

External integration configuration is application-owned deployment configuration or registered runtime configuration, depending on whether a value is safe and appropriate to change dynamically.

Provider endpoints, transport credentials, TLS settings, and bootstrap connection values remain outside component-owned tables. Components register definitions and typed requirements, not provider secrets.

The external-integration runtime should support named provider or adapter configuration with validated settings for:

- Endpoint or destination policy
- Credential reference
- Connection and request timeout
- Retry policy
- Retry budget
- Rate limit
- Circuit-breaker thresholds and recovery interval
- Required or optional readiness classification

Secrets must be redacted from configuration displays, logs, diagnostics, and audit records. The adapter should receive only the credential material it needs through a framework-managed credential boundary rather than spreading secret strings through component code.

## Framework Table Ownership Convention

Every table owned and migrated by yFramework uses the `framework_` prefix. Component-owned tables do not use that prefix and remain under component migrations and repositories.

Components do not directly mutate framework-owned tables except through framework APIs. Framework infrastructure does not directly mutate component-owned production tables except through component-provided typed operations or explicit component commit callbacks.

## MariaDB Compatibility Policy

FreeBSD is the reference platform for MariaDB support. For each yFramework release, the current MariaDB release available in the FreeBSD ports tree when that yFramework release is cut becomes the minimum supported MariaDB version for that yFramework release.

Startup verifies the connected MariaDB version and required capabilities before migrations. Unsupported versions or missing required capabilities prevent startup before migration execution.

Normal startup automatically applies pending registered migrations after static application validation. Migration execution uses a bounded MariaDB-backed advisory lock, deterministic framework-then-component dependency ordering, immutable migration IDs/checksums, and success recording only after migration completion. Successfully applied migrations are not automatically reversed when a later migration fails.

Production migrations should be forward-compatible with the immediately preceding supported application version where practical and should use expand/migrate/contract patterns for breaking structural changes.

## V1 Database Design Contract

The conceptual schema and operational database design are defined in `16_DATABASE_SCHEMA_DESIGN.md`.

Framework-owned tables use the `framework_` prefix. High-write framework entities may use compact internal numeric primary keys plus stable UUID identifiers. Internal numeric IDs are database implementation details; stable UUIDs cross framework subsystem, component, cache, log, administrative, and API boundaries.

Caching reduces repeated reads but does not replace efficient InnoDB schema/index design or MariaDB authority for correctness-sensitive coordination. Live leases, job ownership, workflow state, token revocation, LLM request state, and message delivery remain database-authoritative.

Operational fields that are filtered, joined, sorted, constrained, or indexed remain typed columns. JSON/LONGTEXT is reserved for opaque payloads, type-specific configuration, and content not normally queried by individual field.

High-write coordination rows remain narrow. Large payload/content and append-only attempt history move to companion tables when doing so materially reduces hot-row width or write amplification.

Indexes must correspond to documented production query paths. Avoid speculative indexes and single-column low-cardinality boolean/status indexes without demonstrated query-plan benefit. Critical claim, recovery, authentication, and cleanup queries must be checked with `EXPLAIN` and realistic data volume before v1 release.

Retention is table/content specific. Cleanup uses bounded batches. Partitioning is not a v1 default and is added only after measured scale demonstrates a need.

<!-- ORIGINAL_SOURCE_END: 05_DATA_MIGRATIONS_CACHE_TRANSACTIONS.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 06_JOBS_WORKFLOWS_SCHEDULER.md -->

# Original Source File: `06_JOBS_WORKFLOWS_SCHEDULER.md`

# Jobs, Workflows, and Internal Scheduling

## Core Principle

The framework provides reusable durable job execution and simplified workflow coordination.

Components define the work, own transitional data, isolate domain-record failures, and implement the final production commit.

## Jobs

A job is one durable unit of component-owned background work.

The framework provides:

- Typed registration
- MariaDB queue
- Durable FIFO ordering
- Delayed availability
- Safe claiming
- Worker leases
- Heartbeats
- Retries
- Timeouts
- Cancellation
- Progress reporting
- Idempotency keys
- Durable results
- Request and trace relationships
- Recovery after restart
- Retention and cleanup
- Structured job logging

Possible contract:

```go
type Job[I any, O any] interface {
    Execute(context.Context, I, JobContext) (O, error)
}
```

## Component Ownership

Components own job definitions and handlers.

Example:

```text
components/metadata/jobs/
  CSVParse.go
  CSVValidate.go
  CSVMatch.go
  CSVCommit.go
```

The framework must not contain metadata-specific logic.

## Workflows

A workflow coordinates a durable sequence of component-defined jobs.

First-version features:

- Sequential steps
- One active step at a time
- Durable workflow state
- Typed step input and results
- Retry through the job subsystem
- Resume after restart
- Cancellation
- Progress
- Named failed step
- Commit-phase designation
- Optional cleanup phase

Do not initially build:

- Visual designer
- Arbitrary graph engine
- Scripting
- User-authored expressions
- Complex branching
- General fan-out/fan-in
- Compensation transactions

## Workflow Phases

```text
prepare
commit
cleanup
```

A workflow may contain:

- Zero or more prepare steps
- Zero or one commit step
- Zero or more cleanup steps

All prepare steps must succeed before commit.

Cleanup may follow a successful commit. Failure cleanup may be modeled explicitly later.

Not every workflow needs a commit phase.

## Transitional Data

Intermediate data belongs to the component and may be stored in:

- Durable staging tables
- Files
- A combination of both

Do not use connection-scoped SQL temporary tables for durable workflows.

Workflow payloads should pass references such as import IDs rather than large datasets.

The framework may provide safe workflow-scoped filesystem primitives, but it does not define file contents.

## Final Commit

Production database tables should be changed only in the final commit phase after all preparation steps succeed.

Where practical, the final commit uses one MariaDB transaction.

For very large operations, the component may use idempotent tracked batches while preserving the rule that production mutation occurs only during the final phase.

## Record-Level Failure Isolation

A bad record must not automatically fail an entire batch workflow.

Distinguish:

```text
record rejection -> result data
job execution failure -> error
```

A validation job may complete successfully with:

```text
processed: 10000
accepted: 9821
rejected: 179
```

Detailed per-record errors stay in component staging tables or files.

Only eligible records advance and commit.

The component decides thresholds that convert excessive record failures into workflow failure.

## Generic Statuses

Job statuses:

```text
queued
scheduled
claimed
running
retry_wait
completed
failed
canceled
```

Workflow statuses:

```text
queued
running
commit_pending
committing
completed
completed_with_warnings
failed
canceled
```

## Suggested Tables

```text
framework_jobs
framework_workflows
framework_workflow_steps
framework_schedules
framework_concurrency_locks
```

Jobs should track:

- Queue and name
- Payload and result
- Attempts
- Availability
- Claim and lease
- Worker
- Progress
- Workflow and step
- Trace and request IDs
- User/application/installation identity
- Completion and failure timestamps
- Stable error code

## Monolith Runtime

Jobs, workflows, scheduler, and web/API/MCP handling run inside the monolith service.

Internal goroutines may include:

- Job poller
- Single background executor
- Workflow coordinator
- Schedule poller
- Log rollover
- Cache cleanup

The application decides deployment limits and background-pressure policy.

Components may not start independent background loops or worker pools in v1. Recurring or asynchronous component work must enter framework coordination through a registered job, workflow, schedule, or event subscriber. The framework owns worker, scheduler, lease, retry, resource-pressure, readiness, and shutdown behavior.

## FIFO Dispatch

Version one uses one global durable FIFO lane for executable background work.

- At most one job or workflow step runs at a time.
- The oldest eligible, claimable item runs first.
- Scheduled time controls eligibility, not precedence.
- Delayed, paused, canceled, or concurrency-blocked work does not prevent later eligible work from running.
- A scheduled item may run late under backlog. This is an accepted v1 tradeoff.
- One item is claimed per MariaDB transaction.

V1 has no background priority classes, queue weights, aging, reserved workers, per-component worker pools, or parallel workflow-step execution.

Queue names may remain for administration, pause controls, diagnostics, and future compatibility. They do not affect FIFO ordering.

Foreground requests remain concurrent and do not run through the background FIFO lane.

Lightweight queue control-plane work remains independent because it enables the lane to operate. This includes schedule evaluation, polling, lease heartbeats, lease recovery, wake signaling, and shutdown coordination.


## LLM Outbound Queue

yFramework maintains a durable FIFO queue for outbound LLM requests. This queue is framework infrastructure but is separate from the global durable background-work FIFO described above. It exists so the framework can identify who owns each inference request, preserve trace and request relationships, propagate cancellation, recover unsent work after process failure, and deliver results to the correct caller.

The LLM queue may receive requests from:

- Foreground actions
- Jobs
- Workflow steps
- Event-driven or system-originated framework work

Foreground LLM calls do not enter the global background-work FIFO merely because they use an LLM. A foreground action may synchronously await its LLM result while the framework LLM subsystem manages submission underneath. Background jobs and workflows retain their normal durable ownership and cancellation semantics while using the same LLM service.

The queue uses framework-generated request identity and tracks correlation fields where applicable, including component, LLM purpose, request ID, trace ID, job ID, workflow ID, workflow step, user identity, application identity, and installation identity. The framework request ID should also be used as the yllmd request ID when protocol constraints permit, otherwise an explicit mapping is stored.

Recommended lifecycle states are:

```text
queued
submitted
accepted
running
completed
failed
canceled
```

`queued` means yFramework has durably accepted the request but has not sent it to yllmd. `submitted` means the request was written to the yllmd connection. `accepted` and `running` reflect the corresponding yllmd lifecycle events. Terminal yllmd completion, cancellation, or error updates the framework request record and wakes or notifies the owning caller.

A queued request canceled before submission is removed from execution without contacting yllmd. Once submitted, framework cancellation is propagated to yllmd and the framework waits for a terminal outcome rather than pretending the inference stopped immediately. Go context cancellation is the normal propagation mechanism for foreground calls and job handlers.

The yFramework LLM queue does not replace the yllmd FIFO queue. The framework queue owns application correlation and durable outbound coordination. yllmd owns inference resource scheduling, model routing, and model execution.

Suggested framework table:

```text
framework_llm_requests
```

The exact schema remains an implementation detail, but it must preserve request identity, ownership/correlation fields, resolved yllmd target information, lifecycle timestamps, terminal status, normalized error information, finish reason, and token-usage counters where available. Request and response bodies require an explicit retention policy and should not be retained indefinitely merely for diagnostics.

## Polling

Use centralized adaptive polling rather than one poller per component.

Suggested pattern:

```text
work found -> immediate or very short repoll
idle -> gradual backoff
maximum idle delay -> a few seconds
```

A nonblocking in-process wake signal should notify the poller after local enqueue.

MariaDB remains the source of truth, so lost wake signals do not lose work.

## Execution Semantics

Job execution provides at-least-once semantics.

Exactly-once execution cannot be guaranteed across process crashes, network failures, or uncertain commit outcomes. Job handlers must therefore be idempotent, retry-safe, and responsive to context cancellation.

The framework prevents stale workers from updating framework-owned state through lease fencing. Components remain responsible for making domain mutations idempotent.

## Claiming

Workers claim jobs in short MariaDB transactions. They do not hold row locks while handlers execute.

The claim transaction:

1. Selects one eligible job using `FOR UPDATE SKIP LOCKED`.
2. Acquires any applicable concurrency key.
3. Assigns a process-specific worker identifier.
4. Assigns a unique cryptographically random lease token.
5. Sets the claim and lease timestamps using MariaDB time.
6. Commits before handler execution begins.

Claim eligibility must account for:

- Queue state
- Availability time
- Cancellation state
- Logical concurrency keys
- Whether the FIFO executor is already active
- Automatic or administrative pause state

The minimum supported MariaDB version must support the selected locking strategy.

## Claimed and Running States

`claimed` and `running` are distinct states:

```text
eligible -> claimed -> running
```

`claimed` means a worker reserved the durable record. `running` means the handler began execution.

The attempt count increments when the job enters `running`. A process failure during the short claim-to-start handoff does not consume an execution attempt.

## Lease Ownership and Heartbeats

Each claim records:

- Worker identifier
- Lease token
- Claim timestamp
- Heartbeat timestamp
- Lease expiration timestamp

MariaDB time is authoritative for durable timestamps and expiration comparisons.

Heartbeat, progress, completion, failure, retry, cancellation, workflow advancement, and concurrency-lock release operations must match the current lease token. If a lease-protected update affects no row, the worker no longer owns the job and must cancel the handler context.

Recommended deployment defaults:

```text
lease duration: 90 seconds
heartbeat interval: 20 seconds
```

The heartbeat interval should not exceed one-third of the lease duration. These values are configurable deployment defaults rather than framework constants.

## Lease Expiration and Recovery

An expired `claimed` or `running` job is eligible for recovery.

Recovery:

1. Locks the expired job.
2. Confirms expiration using MariaDB time.
3. Clears the obsolete lease identity.
4. Releases or replaces the associated concurrency lease.
5. Applies the registered retry policy.
6. Moves the job to `retry_wait`, or to `failed` when no attempts remain.
7. Records a structured recovery event.

An expired job that reached `running` has consumed an attempt. A job that expired while only `claimed` has not.

Job and concurrency leases are authoritative. A runtime-instance table may support diagnostics, but worker-presence records do not determine reclaim eligibility.

## Lease Fencing

A stale worker cannot:

- Mark a job completed, failed, retried, or canceled
- Store an authoritative framework-owned result
- Advance a workflow
- Release another worker's concurrency lock
- Update authoritative progress

Every authoritative framework update requires the current lease token.

For a workflow commit step, lease ownership must be verified in the same MariaDB transaction that changes production data. This prevents an expired worker from committing after another worker has reclaimed the step.

## Timeouts and Cooperative Cancellation

Maximum execution time is cooperative. The framework cancels the handler context when the execution deadline is reached and continues renewing the lease during a short cancellation grace period.

Go cannot safely terminate a non-cooperative goroutine. Therefore:

- Job handlers must observe context cancellation.
- Framework database and filesystem helpers must accept the job context.
- Cancellation behavior must be tested.
- Lease expiration must not be described as forcibly stopping a handler.

If a handler refuses to stop, terminating the process is the only reliable way to end it. Its durable lease remains in place until it expires.

## Retry Policy

Each job definition registers a bounded retry policy containing:

- Maximum attempts
- Initial delay
- Backoff multiplier
- Maximum delay
- Jitter
- Retryable error classifications

Recommended defaults:

- Unknown execution errors are retryable until attempts are exhausted.
- Explicit permanent errors fail immediately.
- Panics are recovered, logged as execution failures, and retried when attempts remain.
- Record-level rejections remain successful result data.
- Backoff includes jitter.
- Deployment policy may reduce component retry limits.
- Unlimited retries are not permitted.

Lease loss prevents the stale worker from selecting the durable final state.

## Cancellation

Cancellation is a durable request.

A queued job may move directly to `canceled`. For active work, the framework records the cancellation request, cancels the handler context, and waits for the handler to exit before recording `canceled`.

A job must not be marked canceled while its handler may still be running.

Workflow cancellation prevents creation of new steps and requests cancellation of the active job. Commit-phase cancellation behavior must be defined with workflow transition rules.

## Graceful Shutdown

Background shutdown proceeds in this order:

1. Stop new claims.
2. Continue heartbeats for active work.
3. Allow active handlers to drain.
4. Cancel remaining handler contexts when the background shutdown deadline arrives.
5. Allow a short cancellation grace period.
6. Release a job only after its handler has definitely stopped.
7. Leave leases intact for handlers that have not stopped.
8. Exit and allow remaining leases to expire naturally.

The framework never clears a lease merely because the shutdown deadline elapsed. Early release could allow another worker to start while the original handler is still executing.

## Timezone

Schedules use the operating system's local timezone.

Persist timestamps in UTC where practical, but calculate recurring local schedules with `time.Local` and real timezone rules.

Do not assume every day lasts exactly 24 hours.

## Missed Runs

Due work is selected with:

```text
next_run_at <= current time
```

If the service is paused, overloaded, or offline, due work remains due and runs when capacity returns.

Support:

- `coalesce`: one catch-up execution processes accumulated work
- `replay`: one execution per missed occurrence

Default to `coalesce`.

## Background Suspension

Background pressure control has three recommended states:

```text
normal -> throttled -> paused
```

In v1, the FIFO lane already limits execution to one item, so `throttled` may lengthen polling or defer new claims. `paused` claims nothing. Running work ordinarily continues, including heartbeats. The scheduler continues recording due work, so work is delayed rather than lost.

Recommended foreground-pressure signals are:

- Active requests relative to the request concurrency limit
- Foreground MariaDB pool utilization
- MariaDB connection wait duration
- Sustained request queueing or latency
- Explicit administrative pause

CPU and memory signals are optional future extensions.

Automatic transitions use sustained thresholds and hysteresis to prevent flapping. Exact thresholds and sampling intervals remain deployment and implementation details.

Use separate foreground and background connection pools against the same MariaDB database. This is resource isolation, not multi-database support. Background work must not consume every connection required by foreground requests.

Administrative controls include:

- Pause or resume all background claims
- Pause or resume a named queue
- Inspect why automatic throttling or pausing is active

Administrative pauses are durable. Automatic pressure state is process-local. A normal pause drains rather than cancels active work.

## Overlap Prevention

Prevent both:

1. Duplicate execution of the same job record
2. Concurrent execution of the same logical operation

Job rows use atomic claiming and leases.

Logical operations use concurrency keys such as:

```text
metadata-import:4821
yllmlog-error-analysis
provider-lookup:openlibrary
```

Recommended schedule overlap policies:

```text
allow
skip
queue_one
```

Default to `queue_one`.

`queue_one` allows one active execution and records at most one pending catch-up run.

Workflows may hold their concurrency key for their full lifetime.

## Concurrency-Lock Ownership

Concurrency locks use ownership tokens independent from job lease tokens:

```text
framework_concurrency_locks
  concurrency_key
  owner_type
  owner_id
  owner_token
  acquired_at
  heartbeat_at
  lease_expires_at
```

Job-scoped keys last for one job execution. Workflow-scoped keys may last across several steps and pauses.

Acquisition and job claiming occur in the same short MariaDB transaction. A lock may be inserted when absent or replaced only after its lease expires. Renewal, release, and replacement require the current concurrency ownership token, and MariaDB time governs expiration.

For workflow commit, the production transaction verifies both the current job lease token and the current workflow concurrency-lock token.

Concurrency keys remain useful in v1 for overlap prevention, coalescing, crash recovery, and future parallel execution even though only one executable background item runs at a time.

## Recovery

Workflow state is durable in MariaDB.

After restart, the coordinator resumes by inspecting:

- Current workflow step
- Active job status
- Expired leases
- Completed job result
- Need for next-step creation
- Commit readiness
- Cleanup readiness

## Idle-Only Work

The background subsystem supports work that is durable but eligible to execute only when the application is considered idle. This is an eligibility rule, not a second priority queue or worker pool.

Recommended execution policies:

```text
normal
idle_only
```

An `idle_only` item remains in normal durable background coordination but is skipped while foreground pressure or other configured signals indicate the service is busy. Later eligible normal work may continue to run. When idle conditions are satisfied, the oldest eligible idle-only item may run through the same single executable background lane.

Idle eligibility should reuse the framework's resource-pressure observations, including foreground request activity, MariaDB pool utilization and wait duration, request queueing or latency, and administrative pause state. Exact thresholds, sampling windows, and hysteresis are deployment details.

Preferred localization, adaptive-language backfill, and generation of framework operational localization messages are normal examples of idle-only work.

## Localization Scheduling

When a component reports a localization-eligible source revision:

- Every supported localized unit whose source changed becomes missing or stale as appropriate.
- Manually preferred and adaptive-preferred locales queue eligible regeneration as idle-only work.
- Supported non-preferred locales remain on-demand unless already requested.
- Duplicate work for the same resource, locale, source revision, and localization unit must be prevented through stable concurrency keys or equivalent durable uniqueness.
- Human-modified or locked stale units do not auto-regenerate and instead require administrator review.

When a supported locale is adaptively promoted, all existing eligible content missing or stale for that locale is queued for idle-time backfill. V1 does not automatically demote adaptive-preferred locales.

Notification delivery is also durable background work. Failure to deliver a notification does not reverse the already committed operation that triggered it.

## External Delivery and Retry Budgets

External integrations may use the background subsystem when the operation is a durable asynchronous side effect, but not every external call is a job. Request-path reads may remain synchronous when their semantics require it.

Retries for external providers are bounded at both the individual-operation level and, where appropriate, by an aggregate provider retry budget. A provider outage must not allow thousands of queued records to multiply into an uncontrolled retry storm.

When a retry budget is exhausted, further retryable work remains deferred until budget becomes available or an administrator intervenes. Circuit-breaker state may independently make a provider temporarily ineligible for new attempts.

Retry ownership must be explicit. The framework must avoid stacking independent retry loops across jobs, transport clients, provider SDKs, and adapters in a way that multiplies attempts unexpectedly.

## V1 Schedule Representation

V1 supports three schedule kinds:

```text
once
cron
interval
```

`once` identifies one absolute instant. `cron` uses standard five-field cron syntax with one-minute resolution and is interpreted in the operating system's local timezone. Per-schedule timezone overrides are outside v1. `interval` uses a fixed elapsed duration and a durable anchor so process restart does not reset cadence.

`next_run_at` is durable MariaDB state. A cron library may parse expressions and calculate occurrences but must never execute component work itself.

Schedules may define `max_start_delay`; an occurrence older than that limit is no longer eligible for catch-up. Omitting the value means no age limit. Existing `coalesce` and `replay` missed-run policies remain, with `coalesce` the default.

DST semantics are owned by yFramework rather than inherited accidentally from a cron library: nonexistent spring-forward wall-clock occurrences are skipped, and duplicated fall-back wall-clock occurrences represent one logical scheduled occurrence.

Work created from a schedule retains `schedule_id` and the original `scheduled_for` timestamp independently from enqueue/start/completion timestamps.

## Unified Background Execution Topology

V1 has one global executable background lane for all ordinary durable asynchronous work, including component jobs, workflow steps, scheduled executions, messaging delivery, notification-triggered message delivery, localization work, retries, and other framework/component asynchronous operations. At most one ordinary background item executes at a time.

Schedulers, pollers, workflow coordination, lease heartbeats, retry eligibility checks, wake signaling, and similar control-plane activity may run concurrently but do not independently execute component or delivery work.

The dedicated yFramework LLM outbound queue is the explicit v1 exception because it serves both foreground and background callers and has separate inference ownership/correlation semantics.

Components may not create independent goroutines, worker pools, retry workers, queue executors, or persistent background loops. Event subscribers that need durable asynchronous behavior enqueue framework-coordinated work rather than starting private execution loops.

## V1 Execution Schema Clarifications

The conceptual execution schema is defined in `16_DATABASE_SCHEMA_DESIGN.md`.

V1 ordinary background execution is one eligible-work FIFO lane. The job schema therefore has no named execution queue and no priority-scheduling field. Retries become eligible by durable `available_at` state and do not enter a second retry queue.

`framework_jobs` stores compact authoritative execution/lease state. Large input/results belong in `framework_job_payloads` where needed. Concise attempt history belongs in `framework_job_attempts`. Job attempt count is distinct from messaging-provider or LLM attempts owned by those subsystems.

Schedule targets are registered jobs or workflows only. `framework_schedules` stores the normalized once/cron/interval representation and durable `next_run_at`; optional arguments live in `framework_schedule_payloads` by recommended default. There is no separate v1 schedule-run table because generated jobs retain `schedule_id` and original `scheduled_for` correlation.

Logical resource concurrency continues to use `framework_concurrency_locks`, which is distinct from migration advisory locking and from a job's own execution lease.

<!-- ORIGINAL_SOURCE_END: 06_JOBS_WORKFLOWS_SCHEDULER.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 16_DATABASE_SCHEMA_DESIGN.md -->

# Original Source File: `16_DATABASE_SCHEMA_DESIGN.md`

# Database Schema Design and Operational Model

## Status

This document records the conceptual yFramework v1 database model and the performance, retention, indexing, caching, and ownership decisions that constrain concrete MariaDB DDL implementation.

Exact MariaDB column sizes, binary encodings, index names, and implementation-level CHECK constraints remain DDL engineering decisions unless explicitly settled below.

## Database Platform

### Settled requirement

MariaDB is the only relational database supported by yFramework.

FreeBSD is the reference platform for MariaDB support. Each yFramework release targets the current MariaDB version available in the FreeBSD ports tree when that yFramework release is cut. That MariaDB release becomes the minimum supported MariaDB release for that yFramework release.

Startup verifies the connected server version and required capabilities before migrations are attempted. Unsupported versions or missing required capabilities are fatal startup errors.

## Framework Table Ownership

### Settled requirement

Every framework-owned table uses the `framework_` prefix.

Components do not directly modify framework-owned tables. They use framework APIs.

The framework does not directly modify component-owned production tables except through component-defined repositories, typed callbacks, workflow commit operations, or equivalent component-owned APIs.

Framework-defined foreign keys must not reach into component-owned tables.

## Storage Engine and Relational Model

### Recommended default

Use InnoDB for framework-owned durable tables.

Use relational typed columns for operational fields that are regularly filtered, sorted, joined, constrained, or indexed. JSON/LONGTEXT is appropriate for opaque payloads, type-specific configuration, provider metadata, or content not normally queried by individual field.

Do not hide lifecycle, scheduling, authorization, queue, identity, retry, lease, or correlation fields inside JSON for schema convenience.

## Internal and Stable Identifiers

### Settled requirement

For durable framework entities that cross subsystem, component, cache, logging, admin, or API boundaries, use two identities when beneficial:

```text
id
    compact internal database identity, normally BIGINT UNSIGNED

<entity>_id
    stable UUID identity used outside the owning database implementation
```

Internal framework foreign keys normally use the numeric `*_internal_id` convention. Stable cross-boundary identifiers use `*_id` UUIDs.

Examples:

```text
framework_users.id                  internal
framework_users.user_id             stable
framework_passkeys.user_internal_id internal FK
component author_user_id             stable user UUID
```

Pure association rows may use composite keys and do not require synthetic UUIDs.

### Recommended default

Where supported by the selected MariaDB release and Go UUID implementation, use time-ordered UUIDs such as UUIDv7 for new stable identifiers. Concrete storage representation is a DDL implementation decision after benchmarking.

## Caching and Identifier Performance

### Settled requirement

Caching reduces repeated MariaDB reads but does not replace efficient durable schema/index design. InnoDB still maintains clustered and secondary indexes on writes even when an object is cached.

Cache keys use stable public identifiers where appropriate. Internal numeric IDs remain database implementation details and must not leak across component boundaries.

Read-mostly identity/configuration data may be cached aggressively with correct invalidation. Correctness-sensitive volatile coordination state such as leases, job ownership, workflow transitions, token revocation, and message delivery remains MariaDB-authoritative.

## Hot Rows, Content, and History

### Settled requirement

High-write coordination tables contain the minimum authoritative operational state needed for claiming, scheduling, recovery, and status inspection.

Large content/payloads and append-only attempt history are stored separately where doing so materially reduces hot-row width or write amplification.

The standard pattern is:

```text
current-state table
content/payload table when needed
attempt/history table when useful
```

Examples:

```text
framework_jobs
framework_job_payloads
framework_job_attempts

framework_llm_requests
framework_llm_request_content
framework_llm_request_attempts

framework_messages
framework_message_content
framework_message_attempts
```

## Index Design

### Settled database guideline

Indexes are designed from documented production query paths. Do not add one index per column or speculative indexes for every administrative search.

Each material schema index should have a documented query or invariant it supports.

Avoid single-column boolean or low-cardinality status indexes unless `EXPLAIN`/benchmarking demonstrates value. Compound indexes may begin with status/enablement when that structure supports a real selective or ordered hot-path query.

Critical claim, recovery, lookup, expiration, and cleanup queries must be validated with `EXPLAIN` and realistic row counts before v1 release.

## Time

### Recommended default

Use UTC durable timestamps with microsecond resolution, normally `DATETIME(6)`, for operational state such as queued, available, heartbeat, lease expiry, started, completed, and retention timestamps.

MariaDB time is authoritative for database coordination such as lease and claim expiry.

Calendar schedule calculation uses the configured OS local timezone according to the scheduling specification, then stores resulting absolute instants durably.

## Retention and Cleanup

### Settled requirement

Retention is table/content specific. Active records remain durable while active. Large/sensitive content may expire sooner than lightweight operational metadata.

Cleanup operates in bounded batches and commits between batches. Do not run unbounded mass deletes against hot operational tables.

Tables subject to cleanup expose an indexed timestamp such as `completed_at`, `finished_at`, or `retain_until` so cleanup does not inspect JSON or infer age through expensive joins.

### Recommended default

Do not partition framework tables in v1. Add time partitioning later only when measured table size/retention cost justifies it and the foreign-key/lifecycle implications are understood.

## Foreign Keys and Cascades

### Settled requirement

Use real foreign keys within framework-owned schema families where lifecycle and retention align.

Prefer stable UUID correlation rather than hard cross-subsystem foreign keys when independent retention is desirable.

Examples of appropriate hard FKs:

```text
passkey -> user
user_role -> user
user_role -> role
message_content -> message
message_attempt -> message
llm_content -> llm_request
llm_attempt -> llm_request
```

Examples where stable UUID correlation is preferred:

```text
LLM request -> job
message -> workflow
job -> schedule
notification/message correlation
```

Avoid broad cascading deletes through large operational/audit histories. Disable/tombstone durable principals where history must remain. Cascades are appropriate for tightly bound small configuration children when deletion semantics are clear.

## Foreground and Background Pools

### Settled requirement

Runtime uses separate foreground and background MariaDB pools against the same database. Startup/migration work uses a dedicated startup/bootstrap connection before normal runtime pools are established.

Background capacity is deliberately bounded so jobs, queue polling, cleanup, messaging, workflows, and LLM persistence cannot consume all request-serving database connections.

# Conceptual Table Inventory

## Migration and Runtime Configuration

### `framework_migrations`

Purpose: authoritative append-only record of applied framework, application, and component migrations.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
owner_type          bounded machine name
owner_name          bounded stable name
migration_name      bounded stable name
checksum            checksum
applied_at          DATETIME(6)
duration_ms         nullable
```

Critical uniqueness:

```text
UNIQUE(owner_type, owner_name, migration_name)
```

Applied migration content is immutable. A checksum mismatch for an already-applied identity is a startup failure.

Retention: permanent.
Cache: none required.

### `framework_runtime_settings`

Purpose: registered dynamic application-wide settings only. It is not a generic dumping-ground key/value table.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
setting_name        bounded stable name UNIQUE
value_json          JSON/LONGTEXT
revision            BIGINT UNSIGNED
updated_at          DATETIME(6)
updated_by_type     bounded actor type
updated_by_id       UUID nullable
```

The revision increments when the value changes and supports immutable runtime snapshots, optimistic administration, and cache invalidation.

Routes, tokens, schedules, widgets, provider/domain account records, component domain rules, and complex component configuration retain dedicated schemas.

# Identity and Authentication

## `framework_users`

Purpose: minimal framework security identity, not application profile data.

Conceptual fields:

```text
id                      BIGINT UNSIGNED PK
user_id                 UUID UNIQUE
status                  active/disabled
authorization_revision  BIGINT UNSIGNED
created_at              DATETIME(6)
disabled_at             DATETIME(6) nullable
```

Names, addresses, avatars, phone numbers, biographies, preferences, and domain profile data are component/application-owned.

User identities are normally disabled/tombstoned rather than immediately hard-deleted while referenced by operational or audit history.

## `framework_passkeys`

Purpose: durable WebAuthn credentials. Multiple usable passkeys per user are supported.

Conceptual fields include:

```text
id                      BIGINT UNSIGNED PK
passkey_id              UUID UNIQUE
user_internal_id        BIGINT UNSIGNED FK
credential_id           unique binary credential identifier
public_key/credential data required by selected WebAuthn implementation
label                   nullable
status                  active/disabled
created_at              DATETIME(6)
last_used_at            DATETIME(6) nullable
disabled_at             DATETIME(6) nullable
sign_count              nullable when required
```

Exact WebAuthn storage fields are chosen with the WebAuthn library and must preserve all information required by that implementation.

Disabling/removing the final usable passkey is prohibited and enforced transactionally, not only in the UI.

Usage timestamps should be throttled rather than written on every authentication unless correctness of a specific WebAuthn field requires otherwise.

## `framework_user_sessions`

Purpose: browser/session authentication state after passkey login. Browser sessions are stored separately from personal bearer API tokens because their security and expiry semantics differ.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
session_id          UUID UNIQUE
user_internal_id    BIGINT UNSIGNED FK
token_hash          unique cryptographic verifier
token_prefix        optional non-secret display prefix
created_at          DATETIME(6)
last_seen_at        DATETIME(6)
expires_at          DATETIME(6)
idle_expires_at     DATETIME(6) nullable
recent_auth_at      DATETIME(6)
revoked_at          DATETIME(6) nullable
```

Raw session secrets are never stored. `recent_auth_at` supports sensitive operations requiring recent passkey authentication.

`last_seen_at` updates are throttled so read-heavy session validation does not become a write-per-request workload.

## `framework_user_tokens`

Purpose: personal bearer API credentials.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
token_id            UUID UNIQUE
user_internal_id    BIGINT UNSIGNED FK
name                bounded human label
token_prefix        non-secret display prefix
token_hash          unique cryptographic verifier
created_at          DATETIME(6)
last_used_at        DATETIME(6) nullable
expires_at          DATETIME(6) nullable
revoked_at          DATETIME(6) nullable
```

The UUID identifies the credential administratively but is not itself the bearer secret. Raw bearer values are returned only when created and are never stored.

Framework-generated high-entropy bearer secrets use an efficient cryptographic equality verifier/hash rather than slow password hashing intended for low-entropy human passwords.

`last_used_at` updates are throttled.

If v1 supports token-specific permission narrowing, use a normalized token-permission relationship rather than a JSON permission array. Do not add it unless the authorization API requires it.

# Authorization

## `framework_roles`

Purpose: application-defined administrative groupings of registered permissions. Components do not check role names.

Conceptual fields:

```text
id              BIGINT UNSIGNED PK
role_id         UUID UNIQUE
name            bounded stable name UNIQUE
label           human label
description     nullable
enabled         boolean
created_at      DATETIME(6)
updated_at      DATETIME(6)
```

## `framework_role_permissions`

Pure relationship:

```text
role_internal_id    BIGINT UNSIGNED FK
permission_name     bounded registered name
created_at          DATETIME(6)
PRIMARY KEY(role_internal_id, permission_name)
```

Permission names are validated against the frozen compiled registry.

## `framework_user_roles`

```text
user_internal_id    BIGINT UNSIGNED FK
role_internal_id    BIGINT UNSIGNED FK
created_at          DATETIME(6)
granted_by_type     bounded actor type
granted_by_id       UUID nullable
PRIMARY KEY(user_internal_id, role_internal_id)
```

A reverse `(role_internal_id, user_internal_id)` index may support administration.

## `framework_user_permissions`

Direct grants only; no persisted deny grants in v1.

```text
user_internal_id    BIGINT UNSIGNED FK
permission_name     bounded registered name
created_at          DATETIME(6)
granted_by_type     bounded actor type
granted_by_id       UUID nullable
PRIMARY KEY(user_internal_id, permission_name)
```

Effective global user permissions are the union of role-derived and direct grants.

Changing role membership, role permissions, or direct grants increments the affected user's authorization revision and invalidates the resolved authorization cache immediately.

# Registered Application and Installation Identity

The authenticated non-human principal is distinct from the host `application.Application` composition object. Go APIs should prefer a name such as `ClientIdentity` or `RegisteredApplication` for the principal to avoid ambiguity.

## `framework_applications`

```text
id                      BIGINT UNSIGNED PK
application_id          UUID UNIQUE
name                    bounded stable name UNIQUE
label                   human label
status                  active/disabled
authorization_revision  BIGINT UNSIGNED
created_at              DATETIME(6)
disabled_at             DATETIME(6) nullable
```

## `framework_application_permissions`

```text
application_internal_id BIGINT UNSIGNED FK
permission_name         bounded registered name
created_at              DATETIME(6)
granted_by_type         bounded actor type
granted_by_id           UUID nullable
PRIMARY KEY(application_internal_id, permission_name)
```

Authority belongs to the application principal, not to the token credential itself unless a future token-specific narrowing feature is explicitly added.

## `framework_application_tokens`

```text
id                          BIGINT UNSIGNED PK
token_id                    UUID UNIQUE
application_internal_id     BIGINT UNSIGNED FK
name                        bounded label
token_prefix                non-secret display prefix
token_hash                  unique verifier
created_at                  DATETIME(6)
last_used_at                DATETIME(6) nullable
expires_at                  DATETIME(6) nullable
revoked_at                  DATETIME(6) nullable
generation                  nullable implementation field for overlapping rotation generations
```

Raw token values are never stored. Usage timestamps are throttled.

## `framework_installations`

```text
id                          BIGINT UNSIGNED PK
installation_id             UUID UNIQUE
application_internal_id     BIGINT UNSIGNED FK
name                        nullable label
status                      active/disabled
authorization_revision      BIGINT UNSIGNED
created_at                  DATETIME(6)
last_seen_at                DATETIME(6) nullable
disabled_at                 DATETIME(6) nullable
```

No hardware fingerprint, serial number, or other device identifier is required by the framework.

## `framework_installation_permissions`

```text
installation_internal_id    BIGINT UNSIGNED FK
permission_name             bounded registered name
created_at                  DATETIME(6)
granted_by_type             bounded actor type
granted_by_id               UUID nullable
PRIMARY KEY(installation_internal_id, permission_name)
```

Effective installation authority is constrained by the registered application. For a request requiring both user and application/install identity, effective permission is the intersection of the user permission set and the delegated non-human capability sets.

## `framework_installation_tokens`

Same credential principles as application tokens:

```text
id                              BIGINT UNSIGNED PK
token_id                        UUID UNIQUE
installation_internal_id        BIGINT UNSIGNED FK
token_prefix                    non-secret display prefix
token_hash                      unique verifier
created_at                      DATETIME(6)
last_used_at                    DATETIME(6) nullable
expires_at                      DATETIME(6) nullable
revoked_at                      DATETIME(6) nullable
generation                      nullable
```

# Security Audit

## `framework_security_events`

### Recommended v1 requirement

Provide a durable append-only security and authorization audit stream distinct from ordinary logs.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
security_event_id   UUID UNIQUE
event_name          bounded stable name
actor_type          bounded actor type
actor_id            UUID nullable
subject_type        bounded subject type
subject_id          UUID nullable
result              bounded result value
request_id          UUID nullable
trace_id            UUID nullable
metadata_json       bounded JSON/LONGTEXT nullable
created_at          DATETIME(6)
```

Events may include passkey changes, token/session revocation, role assignment/removal, direct permission changes, application/installation disablement, and contact verification.

The audit table never stores bearer secrets, passkey credential material, message bodies, LLM prompts/responses, or other unrestricted sensitive content. Retention is longer than ordinary operational history and administrator-configurable.

# Routes

## `framework_routes`

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
route_id            UUID UNIQUE
name                bounded stable name UNIQUE
path                normalized canonical path
parent_internal_id  BIGINT UNSIGNED nullable self-FK
enabled             boolean
created_at          DATETIME(6)
updated_at          DATETIME(6)
```

Structural path ambiguity such as `/articles/{id}` versus `/articles/{slug}` is rejected during startup/snapshot validation even if SQL uniqueness cannot express the normalized shape rule.

## `framework_route_handlers`

```text
id                      BIGINT UNSIGNED PK
handler_id              UUID UNIQUE
route_internal_id       BIGINT UNSIGNED FK
method                  bounded HTTP method
action_name             registered action name
action_kind             bounded action kind
renderer_name           nullable registered renderer
view_name               nullable
layout_name             nullable
access_policy           public/user/application/application_and_user
recent_auth_required    boolean
enabled                 boolean
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

V1 has one enabled handler per route/method. There is no handler priority, competing conditional handler, executable predicate, or stored expression.

## `framework_route_handler_permissions`

```text
handler_internal_id BIGINT UNSIGNED FK
permission_name     bounded registered name
PRIMARY KEY(handler_internal_id, permission_name)
```

All listed permissions are required. Routes reference permissions, never roles.

## `framework_route_input_mappings`

This is a subordinate configuration row and does not require a public UUID.

Conceptual fields:

```text
id                  BIGINT UNSIGNED PK
handler_internal_id BIGINT UNSIGNED FK
destination         bounded typed action-input field
source_type         path/query/body/header/fixed/context/user_identity/application_identity/installation_identity
source_name         nullable
fixed_value_json    nullable
required            boolean
created_at          DATETIME(6)
updated_at          DATETIME(6)
```

Ambiguous duplicate destination mappings fail validation. Security-sensitive framework-owned headers are not available as ordinary mapped input.

Routes, handlers, permissions, and input mappings are loaded into a validated immutable in-memory route snapshot. Normal request routing does not query MariaDB.

# Widgets

## `framework_widget_areas`

```text
id          BIGINT UNSIGNED PK
area_id     UUID UNIQUE
name        stable name UNIQUE
label       nullable
enabled     boolean
created_at  DATETIME(6)
updated_at  DATETIME(6)
```

## `framework_widget_instances`

```text
id                  BIGINT UNSIGNED PK
widget_id           UUID UNIQUE
widget_name         registered widget type
instance_name       stable/human identifier
configuration_json  type-specific configuration nullable
enabled             boolean
created_at          DATETIME(6)
updated_at          DATETIME(6)
```

The registered widget decodes and validates its typed configuration. JSON is appropriate because type-specific widget configuration is not normally queried by individual field.

## `framework_widget_placements`

```text
id                  BIGINT UNSIGNED PK
placement_id        UUID UNIQUE
widget_internal_id  BIGINT UNSIGNED FK
area_internal_id    BIGINT UNSIGNED FK
position            integer
enabled             boolean
created_at          DATETIME(6)
updated_at          DATETIME(6)
```

Placement scope is normalized through:

```text
framework_widget_placement_routes
framework_widget_placement_actions
framework_widget_placement_handlers
```

Each scope relation records include/exclude behavior. Explicit exclusions win. Includes, when present, require a match. Remaining placements sort deterministically by position and stable tie-breaker.

No executable placement expressions are stored in MariaDB.

Widget configuration is loaded into validated immutable in-memory snapshots and pre-indexed for normal rendering.

# Jobs

## `framework_jobs`

Purpose: authoritative current state for all ordinary durable executable background work.

### Settled requirement

V1 has one global ordinary executable FIFO lane. There are no named execution queues, priority scheduling, queue weights, per-component worker pools, or hidden independent worker systems.

Therefore the v1 job row has no `queue_name` and no execution `priority` field.

Conceptual fields:

```text
id                      BIGINT UNSIGNED PK
job_id                  UUID UNIQUE
job_name                registered job name
component_name          nullable stable owner
status                  queued/running/retry_wait/completed/failed/canceled
available_at            DATETIME(6)
queued_at               DATETIME(6)
started_at              DATETIME(6) nullable
completed_at            DATETIME(6) nullable
attempt_count           integer
max_attempts            integer
lease_owner             nullable
lease_token             secret binary token nullable
lease_expires_at        DATETIME(6) nullable
heartbeat_at            DATETIME(6) nullable
cancel_requested_at     DATETIME(6) nullable
workflow_id             UUID nullable
workflow_step_id        UUID nullable
schedule_id             UUID nullable
scheduled_for           DATETIME(6) nullable
actor_type              nullable
actor_id                UUID nullable
trace_id                UUID nullable
request_id              UUID nullable
error_code              nullable
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

Eligible work is FIFO using durable eligibility/enqueue order and stable row tie-breaking. Retries become eligible by changing `available_at`; they do not enter a separate retry queue.

Critical claim/recovery queries use short transactions and `FOR UPDATE SKIP LOCKED`, then execute outside the claim transaction.

Completion, heartbeat, and lease operations are fenced with the current lease token. The lease token is secret and never logged.

Authoritative job state is not cached.

## `framework_job_payloads`

One-to-one optional large content:

```text
job_internal_id BIGINT UNSIGNED PK/FK
input_payload   JSON/LONGTEXT nullable
result_payload  JSON/LONGTEXT nullable
created_at      DATETIME(6)
updated_at      DATETIME(6)
```

Queue scans do not touch large payload content.

## `framework_job_attempts`

Concise attempt history:

```text
id                  BIGINT UNSIGNED PK
attempt_id          UUID UNIQUE
job_internal_id     BIGINT UNSIGNED FK
attempt_number      integer
started_at          DATETIME(6)
finished_at         DATETIME(6) nullable
worker_id           nullable
outcome             nullable
error_code          nullable
error_detail        nullable bounded text
lease_lost          boolean
crash_recovered     boolean
duration_ms         nullable
```

Job attempt count means execution attempts of the background work item. It is distinct from transport attempts owned by messaging or LLM subsystems.

# Workflows

## `framework_workflows`

Stores coordination metadata only. Component transitional/domain data remains component-owned.

Conceptual fields:

```text
id              BIGINT UNSIGNED PK
workflow_id     UUID UNIQUE
workflow_name   registered workflow name
component_name  nullable
status          bounded lifecycle state
current_phase   nullable
created_at      DATETIME(6)
started_at      DATETIME(6) nullable
completed_at    DATETIME(6) nullable
actor_type      nullable
actor_id        UUID nullable
trace_id        UUID nullable
request_id      UUID nullable
error_code      nullable
```

## `framework_workflow_steps`

```text
id                      BIGINT UNSIGNED PK
workflow_step_id        UUID UNIQUE
workflow_internal_id    BIGINT UNSIGNED FK
step_name               registered/stable name
phase                   prepare/commit/cleanup
sequence_number         integer
status                  bounded lifecycle state
job_id                  UUID nullable
created_at              DATETIME(6)
started_at              DATETIME(6) nullable
completed_at            DATETIME(6) nullable
error_code              nullable
```

Workflow coordinates; jobs execute. Production domain changes occur only in the final commit phase. Record-level rejection is component result data rather than job/workflow execution failure.

# Schedules

## `framework_schedules`

```text
id                      BIGINT UNSIGNED PK
schedule_id             UUID UNIQUE
name                    stable name UNIQUE
target_type             job/workflow
target_name             registered target name
schedule_kind           once/cron/interval
run_at                  DATETIME(6) nullable
cron_expression         nullable standard 5-field cron
interval_seconds        BIGINT UNSIGNED nullable
interval_anchor_at      DATETIME(6) nullable
next_run_at             DATETIME(6) nullable
missed_policy           coalesce/replay
overlap_policy          allow/skip/queue_one
max_start_delay_seconds BIGINT UNSIGNED nullable
enabled                 boolean
last_scheduled_for      DATETIME(6) nullable
last_enqueued_at        DATETIME(6) nullable
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

Schedule targets resolve only to registered job/workflow definitions. MariaDB stores `next_run_at` so due scanning does not repeatedly parse every schedule.

The critical due index supports enabled schedules ordered by `next_run_at`.

`scheduled_for` remains distinct from actual enqueue/start time on generated jobs.

Scheduler behavior is idempotent across process crashes and must not accidentally enqueue the same logical occurrence twice. Exact uniqueness representation is an implementation choice constrained by coalescing semantics.

## `framework_schedule_payloads`

### Recommended default

One-to-one schedule argument payload:

```text
schedule_internal_id BIGINT UNSIGNED PK/FK
payload_json         JSON/LONGTEXT nullable
```

This keeps the scheduler's hot due row compact.

No separate `framework_schedule_runs` table is required in v1; generated jobs carry schedule correlation and original `scheduled_for`.

# Concurrency

## `framework_concurrency_locks`

Logical resource concurrency, distinct from migration advisory locking:

```text
id                  BIGINT UNSIGNED PK
concurrency_key     bounded stable key UNIQUE
owner_type          bounded owner type
owner_id            UUID
owner_token         secret binary fencing token
acquired_at         DATETIME(6)
heartbeat_at        DATETIME(6)
lease_expires_at    DATETIME(6)
```

Release/reclaim operations are fenced with `owner_token`. The token is never logged. Coordination state is never authoritative in cache.

# LLM

## `framework_llm_requests`

One authoritative current-state row per logical LLM request:

```text
id                      BIGINT UNSIGNED PK
llm_request_id          UUID UNIQUE
component_name          stable owner
purpose_name            registered purpose
status                  bounded lifecycle state
queued_at               DATETIME(6)
submitted_at            DATETIME(6) nullable
accepted_at             DATETIME(6) nullable
started_at              DATETIME(6) nullable
completed_at            DATETIME(6) nullable
attempt_count           integer
crash_attempt_count     integer
yllmd_request_id        nullable stable broker identifier
yllmd_group             nullable
yllmd_profile           nullable
yllmd_model             nullable
job_id                  UUID nullable
workflow_id             UUID nullable
workflow_step_id        UUID nullable
user_id                 UUID nullable
application_id          UUID nullable
installation_id         UUID nullable
trace_id                UUID nullable
request_id              UUID nullable
finish_reason           nullable
error_code              nullable
input_tokens            nullable
output_tokens           nullable
total_tokens            nullable
review_required         boolean
retain_until            DATETIME(6)
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

The stable logical request identity persists across attempts. The durable response and final state are committed before completion is reported/waiters are awakened.

## `framework_llm_request_content`

```text
llm_request_internal_id BIGINT UNSIGNED PK/FK
request_payload         LONGTEXT/JSON
response_payload        LONGTEXT/JSON nullable
request_content_type    nullable
response_content_type   nullable
request_size_bytes      nullable
response_size_bytes     nullable
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

A durable source reference may be used instead of duplicating large component-owned source data only when the reference identifies an immutable/versioned source sufficient to determine exactly what was submitted.

## `framework_llm_request_attempts`

```text
id                      BIGINT UNSIGNED PK
attempt_id              UUID UNIQUE
llm_request_internal_id BIGINT UNSIGNED FK
attempt_number          integer
started_at              DATETIME(6)
finished_at             DATETIME(6) nullable
yllmd_request_id        nullable
yllmd_model             nullable
outcome                 nullable
error_code              nullable
error_detail            nullable bounded text
crash_associated        boolean
submission_uncertain    boolean
input_tokens            nullable
output_tokens           nullable
total_tokens            nullable
duration_ms             nullable
```

Crash quarantine defaults to three crash-associated attempts, configurable. Ordinary invalid output or provider failure does not automatically increment the crash counter.

Default prompt/response content retention is 90 days, configurable, with review/quarantine records exempt from ordinary cleanup. Lightweight metadata may be retained longer than content.

Streaming chunks are not persisted for resumability in v1.

# Messaging

## `framework_messages`

Authoritative current delivery state:

```text
id                      BIGINT UNSIGNED PK
message_id              UUID UNIQUE
component_name          nullable
purpose_name            registered/stable purpose
channel                 bounded channel
status                  queued/sending/retry_wait/delivered/failed/canceled
recipient_type          bounded type
recipient_id            UUID nullable
sender_name             nullable
locale                  nullable
send_after              DATETIME(6)
queued_at               DATETIME(6)
started_at              DATETIME(6) nullable
completed_at            DATETIME(6) nullable
attempt_count           integer
max_attempts            integer
error_code              nullable
job_id                  UUID nullable
workflow_id             UUID nullable
notification_id         UUID nullable
trace_id                UUID nullable
request_id              UUID nullable
retain_until            DATETIME(6) nullable
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

Messaging owns delivery retry semantics. The global background lane executes eligible message-delivery work, but the job subsystem does not create a second independent SMTP retry loop.

## `framework_message_content`

```text
message_internal_id BIGINT UNSIGNED PK/FK
subject             nullable
body_text           LONGTEXT nullable
body_html           LONGTEXT nullable
content_type        nullable
template_name       nullable
template_version    nullable
created_at          DATETIME(6)
updated_at          DATETIME(6)
```

Large attachments are not stored inline by default; future support should use durable references or a dedicated attachment abstraction.

## `framework_message_attempts`

```text
id                  BIGINT UNSIGNED PK
attempt_id          UUID UNIQUE
message_internal_id BIGINT UNSIGNED FK
attempt_number      integer
provider_name       nullable
provider_message_id nullable
started_at          DATETIME(6)
finished_at         DATETIME(6) nullable
outcome             nullable
error_code          nullable
error_detail        nullable bounded text
duration_ms         nullable
```

A message attempt is a provider/channel delivery attempt and is distinct from a job execution attempt.

# Notifications

## `framework_notification_contacts`

Conceptual fields:

```text
id                      BIGINT UNSIGNED PK
contact_id              UUID UNIQUE
user_id                 UUID nullable
channel                 bounded channel
destination_hash        binary nullable
destination_encrypted   encrypted value
status                  pending/verified/disabled
verification_token_hash binary nullable
verification_expires_at DATETIME(6) nullable
verified_at             DATETIME(6) nullable
created_at              DATETIME(6)
updated_at              DATETIME(6)
```

The clear destination is required for delivery but must not appear in ordinary logs. A stable hash may support deduplication/lookup where product semantics permit. Exact encryption/key-management design belongs in the security implementation plan.

Verification tokens are hash-only at rest.

## `framework_notification_subscriptions`

```text
id                  BIGINT UNSIGNED PK
subscription_id     UUID UNIQUE
contact_internal_id BIGINT UNSIGNED FK
component_name      registered/stable component
event_name          registered/stable event/resource notification name
resource_type       nullable
resource_id         UUID nullable
locale              nullable
status              active/disabled
created_at          DATETIME(6)
updated_at          DATETIME(6)
disabled_at         DATETIME(6) nullable
```

Subscriptions do not store executable filter expressions. Components own more complex domain subscription logic/data and emit normalized notification intent.

No `framework_notification_attempts` table is required. Once notification intent becomes a message, messaging owns delivery history.

# Cache/Snapshot Strategy for Configuration and Authorization

### Settled requirement

Routes, widget placement, registered runtime settings, and resolved authorization sets are strong cache/snapshot candidates because they are read frequently and change relatively rarely.

Routes/widgets are loaded as complete validated immutable snapshots. Configuration changes build and validate a candidate snapshot, then atomically replace the active snapshot only after successful validation.

Process-local generations may track active snapshot versions. A generic persisted snapshot table is not required in v1.

Authorization caches use stable principal IDs and a principal authorization revision. Permission/role/capability changes increment the revision and cause immediate invalidation. TTL alone is insufficient for security-sensitive authority changes.

# Usage Timestamp Write Reduction

### Recommended default

Do not write usage timestamps on every authenticated request. Throttle updates to fields such as:

```text
session.last_seen_at
user token.last_used_at
application token.last_used_at
installation token.last_used_at
passkey.last_used_at
installation.last_seen_at
```

The exact interval is configurable/implementation-defined. Correctness-sensitive WebAuthn counters are exempt where the authenticator/library requires immediate persistence.

# Conceptual Inventory

The v1 conceptual framework schema currently consists of:

```text
framework_migrations
framework_runtime_settings

framework_users
framework_passkeys
framework_user_sessions
framework_user_tokens

framework_roles
framework_role_permissions
framework_user_roles
framework_user_permissions

framework_applications
framework_application_permissions
framework_application_tokens
framework_installations
framework_installation_permissions
framework_installation_tokens

framework_security_events

framework_routes
framework_route_handlers
framework_route_handler_permissions
framework_route_input_mappings

framework_widget_areas
framework_widget_instances
framework_widget_placements
framework_widget_placement_routes
framework_widget_placement_actions
framework_widget_placement_handlers

framework_jobs
framework_job_payloads
framework_job_attempts

framework_workflows
framework_workflow_steps

framework_schedules
framework_schedule_payloads

framework_concurrency_locks

framework_llm_requests
framework_llm_request_content
framework_llm_request_attempts

framework_messages
framework_message_content
framework_message_attempts

framework_notification_contacts
framework_notification_subscriptions
```

## Intentionally Absent in v1

Do not add these without a new requirement:

```text
framework_queues
framework_workers
framework_dead_letter_queue
framework_retry_queue
framework_schedule_runs
framework_event_log
framework_hook_state
framework_cache_entries
framework_integration_state
framework_llm_stream_chunks
framework_notification_attempts
```

Hooks/events are in-process in v1. Durable/retryable effects become jobs or messages rather than acquiring a second event-execution infrastructure.

# Remaining DDL Engineering

The following are implementation details to be completed in the concrete schema artifact and verified on the target MariaDB release:

- exact bounded `VARCHAR` sizes
- exact UUID binary/native representation
- exact token/hash binary sizes
- foreign-key names and delete/update actions
- index names and column direction/order
- generated/check constraints where useful
- exact batch cleanup defaults
- exact non-LLM operational retention defaults
- query plans and benchmark results for claim/recovery/authentication hot paths

Codex may make these choices only within the behavioral and ownership constraints in this document. A choice that changes a public API, persistence semantics, security behavior, compatibility, failure/recovery behavior, or subsystem ownership is an architectural decision and must be surfaced rather than invented silently.

<!-- ORIGINAL_SOURCE_END: 16_DATABASE_SCHEMA_DESIGN.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 18_MIGRATION_IMPLEMENTATION_CONTRACT.md -->

# Original Source File: `18_MIGRATION_IMPLEMENTATION_CONTRACT.md`

# yFramework Migration Implementation Contract

## Owners

Migration owners are framework, application, and component.

Durable migration identity is the tuple:

- owner type,
- owner name,
- migration name.

Applied migration identities are immutable.

## Registration

Migrations are compiled into Go code or embedded migration assets registered by compiled code.

- No writable runtime migration-directory scanning.
- No executable migration scripts stored in MariaDB.
- Component migration registration is declarative and DB-free.

## Ordering

Deterministic order:

1. framework migrations,
2. application migrations,
3. component migrations in resolved component dependency order.

Within each owner, migrations run in declared sequential order.

## Migration History

`framework_migrations` stores migration identity, checksum, applied time, and useful timing metadata.

Behavior:

- missing record: migration is pending,
- matching record/checksum: already applied,
- matching identity with differing checksum: fatal validation/startup error.

Applied migrations are never silently rewritten.

## Checksums

Every migration registration has a stable checksum derived from its migration source artifact. Runtime compiled-machine-code hashing is not used.

## MariaDB Locking

V1 uses a MariaDB advisory lock on the dedicated migration connection with bounded acquisition wait.

Failure to obtain the lock is fatal for startup/migrate. The lock is explicitly released and is also released by MariaDB if the connection dies.

The advisory migration lock is separate from framework logical concurrency locks.

## Dedicated Migration Connection

Migrations run before runtime foreground/background pools are constructed:

bootstrap config/logger -> compose/register -> dedicated DB connection -> MariaDB compatibility -> migration lock -> migrations -> release/close -> runtime construction.

## MariaDB Compatibility

Before mutation, yFramework verifies the connected MariaDB release and required capabilities.

For each yFramework release, the minimum supported MariaDB version is the current MariaDB version available in the FreeBSD ports tree when that yFramework release is cut.

Unsupported MariaDB versions fail before migrations.

## Transaction Semantics

Use transactions where MariaDB operations support meaningful transactional behavior, but yFramework does not promise rollback of all DDL.

Migration authors must account for MariaDB nontransactional/implicitly committing DDL.

## Recording Success

A migration is recorded as applied only after successful execution and required verification.

Never insert migration history before executing the migration.

## Crash Recovery and Idempotency

Strict expected-state migrations are preferred. `IF EXISTS` / `IF NOT EXISTS` is not used indiscriminately because it can hide drift.

Where known nontransactional DDL can leave a partial state, the migration must explicitly recognize and safely recover that partial state.

## No Automatic Rollback

If a later migration fails, previously successful migrations remain applied.

The framework does not automatically execute down migrations or attempt destructive rollback.

## Expand / Migrate / Contract

Production schema evolution should use forward-compatible expansion, data migration/backfill, and later contraction where practical.

Large data transformations should normally be moved to background job/workflow processing after a compatible schema expansion rather than blocking startup for long periods.

## Compatibility Goal

Schema changes should preserve compatibility with the immediately previous supported application/framework release where reasonably practical.

## Foreign Keys and Destructive Changes

Migration order must respect FK dependencies. Existing data is validated before adding new FK constraints.

Destructive changes require explicit migrations. yFramework does not auto-diff Go structs into destructive DDL.

## CLI Behavior

`install` and `migrate` reuse the same migration engine.

`install` reaches current schema by executing migrations, not a separate drifting schema dump.

`migrate` loads configuration/composition, validates DB compatibility, acquires the lock, applies migrations, and exits without starting request/background runtimes.

`validate` is nonmutating and reports migration consistency, checksums, pending migrations, DB compatibility, and composition problems.

## Observability

Structured migration logs include owner, migration name, checksum, start/end, duration, and result. Sensitive row values are not dumped.

## Required Tests

- empty DB -> current schema,
- previous supported schema -> current,
- checksum drift fails,
- failed migration not marked applied,
- concurrent migrators serialize,
- crashed migration connection releases advisory lock,
- unsupported MariaDB mutates nothing,
- component ordering is respected,
- missing dependency prevents migration,
- known partial DDL recovery where applicable,
- real MariaDB testing, not mocks alone.

## Authority Boundary

Components migrate only their schema. Framework migrations own framework schema. Application migrations own application schema. Cross-owner mutation is prohibited unless a deliberately documented compatibility migration requires it.

<!-- ORIGINAL_SOURCE_END: 18_MIGRATION_IMPLEMENTATION_CONTRACT.md -->
