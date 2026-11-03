# yFramework Consolidated Source: 03_SERVICES_CONFIGURATION

This consolidated project source preserves the original documents listed below verbatim. The wrapper headings and checksum manifest are consolidation metadata only.

## Included source documents

- `09_LLM_YLLMD_INTEGRATION.md` — SHA-256 `a248cb5712d0cfab195dc5b9438460e879038cfcf3e477dd672053d5ba904475`
- `10_LOCALIZATION.md` — SHA-256 `4034bea704a4383c89c4c341c3eb5607c2b3288ebe971273d15737dbf8fb83e5`
- `11_NOTIFICATIONS.md` — SHA-256 `1b40a6ac512bdb81e4c85aa7af0e9e7d93d6127ea2da1ae9b1800df23c33d8b8`
- `12_MESSAGING.md` — SHA-256 `0f12c4fb69caa09673297a400fd1a01340792f363bc55c865860d8c1d8ab8012`
- `13_EXTERNAL_INTEGRATIONS.md` — SHA-256 `48ebcaba35d161089d3c5cd8a1a00f800357fd2804cf0290d307eee76412fe33`
- `17_CONFIGURATION_SCHEMA.md` — SHA-256 `334c37639d8d1a61f6455fa1f4fd549e16706a8bddc3ccf0f08ea28b72cccb3a`

---

<!-- ORIGINAL_SOURCE_BEGIN: 09_LLM_YLLMD_INTEGRATION.md -->

# Original Source File: `09_LLM_YLLMD_INTEGRATION.md`

# LLM and yllmd Integration

## Status and Scope

This document defines the yFramework boundary for LLM use. yllmd is the supported LLM broker. yFramework does not implement a competing provider abstraction, model router, model lifecycle manager, or direct provider integrations.

The first implementation should provide the framework foundation needed for components and applications to use yllmd without depending on yllmd transport details.

## Responsibility Boundary

### Components Own

- Why an LLM is needed
- Prompt and message construction
- Domain context supplied to the request
- Expected output or schema
- Validation and interpretation of the result
- Business behavior that follows from the result
- Stable component-qualified LLM purpose definitions when specialized routing is needed

Components do not open yllmd sockets, encode JSON Lines, manage yllmd queue state, issue yllmd cancellation messages directly, or select concrete model files.

### yFramework Owns

- Typed component-facing LLM APIs
- LLM purpose registration and validation
- Application purpose-to-yllmd-target resolution
- Durable outbound LLM FIFO queue
- Request identity and ownership correlation
- yllmd Unix-socket transport and protocol handling
- Context-driven cancellation propagation
- Normalized framework errors
- Foreground result delivery and asynchronous result lookup primitives as required
- Trace and structured logging integration
- Queue recovery and bounded request-content cleanup

### Application Owns

- Whether LLM functionality is enabled
- yllmd socket and deployment connection configuration
- Purpose-to-yllmd group/profile mappings
- Component composition and injection of the framework LLM capability
- Application-level defaults and administrative policy

### yllmd Owns

- Inference queueing
- Routing groups and profiles
- Concrete model resolution
- Operational model fallbacks
- Model loading and lifecycle
- Provider implementation and evolution
- Model inventory and administrative model operations

yFramework must not duplicate those responsibilities.

## Component API

Components consume a typed framework interface through explicit dependency injection. Exact Go names remain implementation details, but the normal complete-response capability should resemble:

```go
type Generator interface {
    Generate(context.Context, Request) (Response, error)
}
```

Streaming should be a separate capability so components that only require complete responses do not need to implement stream handling:

```go
type StreamingGenerator interface {
    GenerateStream(context.Context, Request) (Stream, error)
}
```

The framework API should expose application-meaningful request and response types instead of leaking raw yllmd wire messages.

## LLM Purposes and yllmd Targets

Reusable components should not hard-code yllmd group and profile names as architectural requirements. A component may register a stable purpose such as:

```text
metadata.classify
metadata.normalize-author
article.summarize
```

The application maps that purpose to a yllmd target:

```text
metadata.classify
  -> group: metadata
     profile: classification
```

A component therefore states what kind of LLM work it needs while the application decides which yllmd route fulfills that need. The yllmd route then resolves to a concrete model according to yllmd configuration.

A request may intentionally omit a specialized purpose. In that case the configured application/default behavior may allow yllmd to use its global default target.

Exact-model targeting should be available only when an application or specialized component has a concrete requirement for it. It should not be the normal reusable-component path.

## Prompt Ownership

Application-specific prompting remains outside yllmd. In yFramework, prompts belong to the component or application capability that understands the domain operation. The framework does not maintain a global prompt library and does not reinterpret component prompts.

Structured-output helpers may be provided by yFramework, but the component owns the expected schema and validates the returned model output. Content-quality or schema-validation failure is not automatically treated as an operational yllmd fallback condition.

## Durable Outbound Queue

yFramework maintains its own durable FIFO queue before sending inference requests to yllmd. The purpose is application ownership and recoverability, not model-resource scheduling.

```text
component/action/job/workflow
  -> yFramework LLM API
  -> durable yFramework LLM FIFO
  -> yllmd client
  -> yllmd FIFO
  -> model/provider
```

The LLM queue is separate from the global background-work FIFO. Foreground actions may submit LLM work and synchronously wait for the response without waiting behind unrelated jobs. Background jobs and workflow steps may use the same LLM subsystem while retaining their normal durable job ownership.

The queue tracks who owns every request. Correlation data should include where applicable:

```text
llm_request_id
component
llm_purpose
action
request_id
trace_id
job_id
workflow_id
workflow_step
user_id
application_id
installation_id
```

The framework-generated LLM request ID should also be used as the yllmd protocol request ID when safe and supported. Otherwise the queue stores an explicit mapping.

## Queue Lifecycle

Recommended framework states:

```text
queued
submitted
accepted
running
completed
failed
canceled
quarantined
```

Definitions:

- `queued`: durably accepted by yFramework but not yet sent to yllmd.
- `submitted`: written to the yllmd transport.
- `accepted`: yllmd acknowledged and queued the request.
- `running`: yllmd reported generation start.
- `completed`: yllmd returned a successful terminal response.
- `failed`: submission, protocol, transport, or terminal generation failure.
- `canceled`: cancellation reached a terminal state.
- `quarantined`: repeated crash-associated execution failures exceeded the configured threshold and automatic submission is stopped pending administrator review.

Suggested framework table:

```text
framework_llm_requests
```

The exact MariaDB schema remains an implementation detail, but records must retain enough information for ownership, recovery, cancellation, result delivery, observability, and cleanup. They should include resolved target information, lifecycle timestamps, normalized errors, finish reason, and token usage when available.

## Request Content and Retention

Durable recovery may require the framework to retain enough input to submit a queued request after restart. Prompt content, source documents, and model output can contain private or proprietary information and must not be retained indefinitely merely for diagnostics.

Settled rules:

- Persist request and result payloads for operational recovery and review.
- Default retention is 90 days and is configurable.
- Records flagged for administrator review or quarantine are not automatically purged until the review condition is cleared.
- Prefer references to component-owned staging data for very large or domain-sensitive inputs when practical.
- Do not log prompts or generated content by default.
- Do not store provider credentials in LLM queue records.

## Cancellation

Go context cancellation is the normal framework cancellation signal.

If a request is still `queued`, yFramework cancels it locally and never sends it to yllmd. If it has been submitted, the framework propagates cancellation through the yllmd protocol and waits for a terminal outcome. The framework must not claim that an active inference stopped before yllmd reports or the transport establishes a terminal condition.

Job and workflow cancellation naturally propagate through the handler context into the LLM service. Foreground request cancellation uses the same mechanism.

## Recovery, Crash Handling, and Quarantine

The framework distinguishes ordinary generation failure, timeout, transport loss, yllmd crash, runner crash, and framework-process interruption where evidence permits. A lost connection must not automatically be labeled an LLM crash.

Crash-associated requests are retried with the same stable framework request identity when it is safe to do so. Every crash-associated attempt is logged and retained in request history. The default maximum is three crash-associated attempts and is configurable. When the threshold is reached, the request becomes `quarantined`, automatic resubmission stops, and administrator review is required.

Recommended failure classes include:

```text
yllmd_crash
runner_crash
transport_lost
timeout
generation_error
invalid_output
```

Only configured crash-associated classes count toward the crash-quarantine threshold.

The framework must also distinguish a request that was never submitted from one whose submission outcome is uncertain because yFramework failed after sending it. The framework should use one stable LLM request ID as the yllmd protocol request ID when yllmd permits it. If yllmd later exposes durable lookup or idempotent resubmission by request ID, yFramework should use that facility for reconciliation.

Until durable broker-side reconciliation exists, authoritative domain writes must remain idempotent. Localization, for example, identifies work by resource, locale, source revision/fingerprint, and localization unit. If the intended current localization already exists, a recovered retry must not create a conflicting duplicate authoritative result.

Administrator tooling must expose quarantined requests and their failure history and permit an authorized manual retry after the underlying prompt, source data, model route, or operational condition has been reviewed.

## Observability

Stable events should include:

```text
llm.request.queued
llm.request.submitted
llm.request.accepted
llm.request.started
llm.request.completed
llm.request.failed
llm.request.cancel_requested
llm.request.canceled
llm.request.quarantined
llm.request.retry_scheduled
```

Useful fields include:

```text
llm_request_id
yllmd_request_id
llm_purpose
yllmd_group
yllmd_profile
yllmd_model
yllmd_fallback
yllmd_fallback_from
llm_finish_reason
llm_input_tokens
llm_output_tokens
llm_total_tokens
llm_queue_wait_us
yllmd_queue_wait_us
llm_generation_duration_us
llm_first_output_at
llm_time_to_first_output_us
llm_estimated_ready_at
llm_estimate_confidence
```

Existing request, trace, component, action, job, workflow, and identity fields remain part of the same log context.

## Configuration

Deployment TOML owns bootstrap connection settings such as:

```toml
[llm]
enabled = true
socket = "/var/run/yllmd/yllmd.sock"
retention_days = 90
crash_quarantine_attempts = 3
throughput_sample_count = 15
```

Connection and request timeout defaults may also be deployment settings. Purpose-to-target mappings are non-executable application configuration and may later be MariaDB-managed dynamic settings if their definitions are registered and validated.

The yllmd configuration remains authoritative for concrete models, model aliases, provider credentials, route fallback order, and model lifecycle.

## Queue Separation

Three distinct coordination layers may exist and must not be conflated:

1. The yFramework background FIFO coordinates durable application jobs and workflow steps.
2. The yFramework LLM FIFO coordinates outbound inference ownership, correlation, and delivery.
3. The yllmd FIFO coordinates inference-resource scheduling and model execution.

They may all be FIFO in v1, but they protect different resources and have different durability and ownership semantics.

## Security Boundary

LLM output is untrusted data. A model response does not gain authority merely because it was produced through yllmd. Any later feature that allows model-directed tools or actions must pass through explicit typed capabilities, normal authorization, validation, auditing, and transaction boundaries.

yFramework must not allow model output to become executable Go, SQL, shell commands, route configuration, authorization changes, or arbitrary stored expressions merely by returning text that requests those operations.

## Implementation Decisions Delegated to Codex

The architecture is settled. Exact exported request/response helper names, internal queue helper types, ordinary index details, administrative method names, and cleanup scheduling may be selected during implementation provided they preserve the contracts in this specification.

## Throughput Sampling and Wait Estimates

The framework maintains rolling LLM throughput observations so callers can receive a rough estimate of queue and generation delay. The rolling window is based on a configurable number of qualifying completed requests, not a time period. The recommended default is the most recent 15 qualifying requests.

Samples should be partitioned by concrete model, and by materially different runtime configuration if later evidence shows that is necessary. Do not combine substantially different models into one throughput estimate.

Where available, record:

```text
accepted_at
started_at
first_output_at
completed_at
input_tokens
output_tokens
total_tokens
```

The current yllmd protocol reports token usage on completion and streaming provides output deltas, so yFramework can at minimum derive queue delay, time to first output for streamed generation, total generation duration, and approximate output-token generation rate. If yllmd or yllama-runner later exposes prompt-evaluation duration or more precise token timing, yFramework should consume those measurements rather than replacing them with local inference.

Before exact token counts are known, request size may be estimated from words or characters. Historical completed measurements should use actual token counts whenever available. Estimated output size may use application or localization-specific heuristics and must remain explicitly approximate.

The estimator may combine work ahead in the yFramework outbound LLM queue, known or estimated queued work, current model throughput observations, and the current request estimate. Returned estimates should include a confidence classification such as:

```text
unavailable
low
medium
high
```

Applications should present rounded language such as `about 1 minute` rather than false second-level precision. Estimates must be recalculated as queue state and completed measurements change.

## External Integration Rule

yllmd integration is one specialized implementation of the framework-wide external-integration rule. Components never speak the yllmd protocol or open its Unix socket directly. The LLM subsystem owns the transport boundary and reuses shared framework conventions for health, tracing, redaction, failure normalization, and bounded operational retry where applicable.

The yFramework LLM durable queue remains separate from generic external retry coordination because it carries LLM-specific ownership, token usage, forecasting, and result state.

## Authoritative Durable Request Record

Every yFramework LLM operation has one durable authoritative `framework_llm_requests` record. The record contains a stable request ID, ownership/correlation data, the outgoing request payload or durable references to large/component-owned source material, the current lifecycle status, resolved routing metadata, lifecycle timestamps, normalized error information, token usage, and the received response when available.

Prompt/request and response content are retained by default for the configured LLM retention period to support crash recovery, troubleshooting, and authorized audit. Large or particularly sensitive source data may instead be represented through durable component-owned references. Prompt and response content remain excluded from normal logs.

The received terminal response is durably committed before completion is reported to or used to wake the caller. The current row is authoritative for present state; an append-only `framework_llm_request_attempts` history may record retries, crash-associated attempts, and reconciliation events.

Administrative access to LLM metadata and access to prompt/response content should be separately permissionable.

## V1 Recovery and Timeout Closure

A submission whose outcome cannot be proven is treated as uncertain rather than assumed successful or failed. Stable request identity is reused for reconciliation/idempotent retry when supported. Component final writes remain idempotent so duplicate inference cannot create conflicting authoritative state.

Streaming responses are not resumable after yFramework process interruption in v1. Complete-response requests retain durable recovery semantics.

The caller or owning job/workflow context is the outermost deadline. Internal queue-wait, transport/connect, and yllmd execution timeouts must fit within that outer deadline and may not extend it. Purpose/application policy may select defaults subject to deployment caps.

## V1 LLM Database Shape

The conceptual schema is defined in `16_DATABASE_SCHEMA_DESIGN.md`.

`framework_llm_requests` is the authoritative compact current-state row for a logical request. The stable framework LLM request ID survives retries. Large request/response content lives in `framework_llm_request_content`, and attempt history lives in `framework_llm_request_attempts`.

Durable completion is committed, including response content and final metadata, before the caller is awakened or success is externally reported. Large component source data may be referenced instead of duplicated only when the durable reference identifies an immutable/versioned source sufficient to establish exactly what was submitted.

Default prompt/response content retention remains 90 days, configurable, with review/quarantine exemptions. Lightweight metadata may be retained longer. No streaming-chunk persistence is required for v1 resumability.

<!-- ORIGINAL_SOURCE_END: 09_LLM_YLLMD_INTEGRATION.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 10_LOCALIZATION.md -->

# Original Source File: `10_LOCALIZATION.md`

# Localization

## Core Model

yFramework provides persistent, LLM-backed localization infrastructure. Primary-language content remains authoritative. Localized content is derived data that is generated through the framework LLM service and yllmd, validated, persisted, and reused on later requests.

The framework does not assume English. Every application defines a primary locale, and individual components or records may support an explicit source locale where necessary.

Localization is broader than literal translation. A component may localize prose, titles, summaries, SEO text, image alt text, or other linguistic content while preserving domain-specific meaning and structure.

## Responsibility Boundary

### Components Own

- Which resource types are localizable
- Stable localization units within each resource
- Source extraction and contextual material
- Localization-specific prompt instructions appropriate to the domain
- Validation of generated localized output
- Component-owned localized domain tables or files
- The source revision or eligibility point at which content is ready for localization
- Final persistence and read behavior for localized domain data

### yFramework Owns

- Localization registration and coordination
- Locale status and availability APIs
- Preferred and adaptive-preferred policy mechanics
- Idle-only scheduling of proactive localization
- On-demand localization coordination
- Duplicate-generation prevention
- Source revision/fingerprint comparison
- Review-required state
- Framework operational localization messages
- LLM request submission through the yFramework LLM queue
- Wait estimates supplied by the LLM estimator
- Administrative coordination APIs

### Applications Own

- Primary locale
- Supported locales
- Manually preferred locales
- Adaptive-promotion thresholds and policy
- Whether a request waits, shows status, or follows another presentation path while localization is pending
- Administrative UI and authorization policy

### yllmd Owns

- Inference scheduling
- Route and model resolution
- Model execution
- Operational model fallback

## Locale Categories

Every public locale is either unsupported or supported. A supported locale may additionally be preferred.

Preferred status may come from:

```text
manual preferred
adaptive preferred
```

Effective preferred behavior is the union of those states. Preferred locales are proactively localized during idle capacity. Supported non-preferred locales are localized on first demand.

V1 may automatically promote a supported locale when aggregate meaningful localization demand reaches a configured threshold. It does not automatically demote the locale. An administrator removes preferred status manually.

When a locale is promoted to preferred, all existing eligible missing or stale content for that locale is queued for idle-time backfill.

## Status

The framework exposes at least:

```text
missing
queued
generating
current
stale
failed
```

Additional administrative state includes:

```text
human_modified
locked
review_required
```

A failed localization is logged and flagged for administrator review. A public request for that failed localized content receives the localized framework content-unavailable message by default. The application may expose additional fallback options, but yFramework does not silently substitute primary-language content.

## Framework Operational Messages

The framework owns a small set of operational messages needed before or during content localization, including at least:

```text
localization in process
localized content unavailable
notification verification
content ready notification
```

For every supported locale, these messages are generated through the LLM subsystem and persisted during idle capacity. They are versioned or fingerprinted so a changed canonical message causes its localized derivatives to become stale and eligible for idle regeneration.

This ensures a visitor can receive a localized status message without requiring another synchronous translation request.

## Stable Localization Units

Only changed content should be regenerated. Components therefore expose stable localization units instead of treating every resource as one indivisible translation blob.

Possible units include:

```text
article.title
article.summary
article.body.paragraph.001
article.body.paragraph.002
article.seo_title
article.seo_description
```

Each unit retains enough state to determine whether its source changed, including a source fingerprint and/or source revision. A component may provide surrounding context to the model for consistency, but only output belonging to stale units is replaced.

## Human Editing and Locks

Humans may edit generated localized content.

A localized unit may be:

```text
generated
human_modified
locked
```

`human_modified` records that a person changed the generated result. `locked` prevents automatic regeneration. A unit may be both human-modified and locked.

If source content changes:

- Generated, unlocked affected units become stale and may regenerate automatically according to locale policy.
- Human-modified or locked affected units become stale and `review_required`.
- yFramework must not silently overwrite a human-modified or locked unit.

Administrator tooling must make review-required units discoverable and provide enough source and localization provenance to review the change.

## Triggering Localization

Components identify when source content becomes localization-eligible. Framework code must not assume every database save is ready for localization.

Examples include publication, an explicit localization-ready revision, or another component-defined lifecycle point.

When an eligible source changes:

1. Compare stable unit fingerprints.
2. Mark only changed localized units stale.
3. Queue affected unlocked units for every preferred locale as idle-only work.
4. Leave supported non-preferred locales on demand unless they have already been requested.
5. Flag affected human-modified or locked units for administrator review instead of regeneration.

## On-Demand Localization

When a supported non-preferred locale is requested and the localization is missing or stale, yFramework creates or joins one durable localization operation for that exact resource, locale, source revision/fingerprint, and localization unit set.

Duplicate visitor requests must not create duplicate LLM work. Stable concurrency keys or equivalent durable uniqueness should identify the localization operation.

While localization is pending, yFramework exposes status, a pre-generated localized in-process message, and an approximate ready-time estimate when enough LLM throughput data exists.

Applications may choose whether a caller waits for completion or receives the status response immediately. The framework supplies the state and estimate; presentation remains application-owned.

## Failure

If localization fails:

- Record structured failure information.
- Mark the affected localization failed.
- Create or set administrator review state.
- Do not expose partial or unvalidated localized content as current.
- Serve the localized content-unavailable operational message for requests that require the failed locale.

Retry behavior follows the LLM subsystem rules. Repeated crash-associated failures may quarantine the underlying LLM request and require administrator review.

## Preferred and Idle Localization

Preferred localization is opportunistic background work. It is durable but `idle_only`, so foreground load and normal background operations can take precedence under configured resource-pressure rules.

Examples:

- New eligible article queues Spanish and French because they are preferred.
- German is supported but not preferred, so it is generated only when first requested.
- Once German crosses the configured adaptive demand threshold, it becomes adaptive-preferred.
- Future eligible content queues German proactively.
- Existing German-missing content is queued for idle-time backfill.

## Administrative Operations

The framework should expose authorization-aware operations that applications such as yCMS can use to build administration interfaces for:

- Localize one resource into one or more locales
- Localize all missing content for a locale
- Regenerate stale unlocked content
- Inspect missing, queued, generating, current, stale, failed, and review-required counts
- Lock or unlock localized units
- Review human-modified or locked stale units
- Retry failed or quarantined work where allowed
- Inspect adaptive locale promotion state

The framework provides APIs and durable state. yCMS or another application owns the administrative pages and workflows.

## Relationship to Notifications

A visitor waiting on localization may optionally subscribe for a one-time notification when the localized content becomes current. The localization subsystem publishes a durable completion fact only after the localized domain data commits successfully. The notification subsystem then handles verification, subscription matching, and delivery.

<!-- ORIGINAL_SOURCE_END: 10_LOCALIZATION.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 11_NOTIFICATIONS.md -->

# Original Source File: `11_NOTIFICATIONS.md`

# Notifications

## Purpose

yFramework provides a generic durable notification layer for subscriptions and event-driven notification intent. Actual outbound delivery is performed only through the framework messaging subsystem. Components and applications do not implement direct SMTP, email-provider, SMS-provider, or other delivery transport.

Localization-ready messages are one use case. Other future uses may include completed exports, reports, imports, publication events, or other durable application facts.

## Responsibility Boundary

### Components Own

- The durable fact or domain event that can trigger a notification
- Domain-specific notification meaning and any component-specific template content
- Whether a particular operation offers notification enrollment

### yFramework Owns

- Notification contact records
- Verification state
- Notification subscriptions
- Durable delivery records
- Channel-independent notification APIs
- Verification-token generation and validation
- Resend throttling/cooldowns
- Expiration and purge of unverified contacts
- Creation of channel-independent messaging requests
- Notification fulfillment state
- Structured logging and audit integration

### Applications Own

- Which notification channels are enabled
- Channel provider configuration
- Public forms or administration UI
- Authorization and enrollment presentation
- Application-specific templates and branding

## Channels

Email is the initial notification channel. SMS/text messaging is a later extension. Notifications target messaging channels, but provider delivery belongs to the messaging subsystem.

Components must not assume SMTP, a particular email API, or a future SMS vendor, and notification code never opens provider connections directly.

## Contacts and Subscriptions

A contact is distinct from a subscription.

Conceptually:

```text
notification_contacts
  id
  channel
  destination
  verification_status
  verification_token_hash
  created_at
  verification_sent_at
  verified_at
  expires_at

notification_subscriptions
  id
  contact_id
  event/resource relationship
  locale
  status
  created_at
  fulfilled_at
```

Exact table names and columns remain implementation details.

This separation lets one verified contact subscribe to more than one eligible notification without re-verifying for every subscription, subject to retention and application policy.

## Anonymous Verification

Anonymous completion notifications require contact verification.

Flow:

```text
contact submitted
  -> pending contact created or reused
  -> verification message sent subject to resend limits
  -> user verifies control of destination
  -> contact becomes verified
  -> requested notification subscriptions may deliver
```

An unverified contact cannot receive the requested content-ready or other completion notification.

Unverified contacts are retained for 30 days by default. During that window, repeated attempts to enroll the same destination should reuse or recognize the pending contact rather than continuously creating new records or verification sends. Verification resends are subject to a configurable cooldown or rate limit.

If the contact remains unverified for 30 days, the pending contact and associated unusable subscriptions are purged according to the retention implementation.

Verification tokens are cryptographically random, shown only through the verification delivery, and stored only as hashes.

## Abuse Prevention

Mandatory verification prevents the notification facility from being used as an unrestricted message relay to third parties.

The framework should additionally support:

- Verification resend cooldowns
- Per-destination and per-source enrollment rate limits
- Generic public responses that do not unnecessarily disclose whether an address is already registered
- Normal security and audit logging without exposing the full destination

Exact thresholds are deployment settings.

## Purpose Limitation

Verification proves control of a destination only.

A contact collected for `notify me when this content is ready` must not automatically become:

- A user account
- A newsletter subscriber
- A marketing contact
- Consent for unrelated notifications

Any broader use requires separate application consent and data handling.

## Delivery

Notification delivery occurs only after the triggering durable fact commits successfully.

```text
domain/localization commit
  -> post-commit event or durable notification creation
  -> framework message created
  -> durable messaging delivery queued
  -> messaging channel adapter sends
  -> notification fulfillment records the delivery outcome
```

A messaging-delivery failure does not reverse the underlying completed operation. Messaging owns bounded retries, retry budgets, circuit-breaker behavior, and provider failure handling. Notification state may be flagged for administrator review when delivery ultimately fails.

## Localization Integration

When localized content becomes current, matching verified subscriptions may receive a one-time ready notification in the requested locale.

If the content is still queued or generating, yFramework may expose both an approximate ready time and an option to subscribe. The notification itself uses pre-generated or otherwise validated localized notification content so sending the message does not require an additional synchronous localization operation.

## Privacy and Logging

Raw email addresses and phone numbers are sensitive operational data.

- Do not place full destinations in general logs.
- Use contact and subscription IDs for correlation.
- Mask destinations only when human diagnostics genuinely require it.
- Do not include verification tokens or token hashes in normal logs.
- Apply explicit retention policies to contacts, subscriptions, and delivery records.

Verified-contact retention beyond active subscriptions remains an application/data-retention policy and should be documented before implementation. The settled 30-day period applies specifically to unverified contacts.

## V1 Notification Database Shape

The conceptual schema is defined in `16_DATABASE_SCHEMA_DESIGN.md`.

`framework_notification_contacts` stores verified contact endpoints using protected destination storage and hash-only verification tokens. Clear destinations must not enter ordinary logs. `framework_notification_subscriptions` stores normalized subscription relationships and contains no executable filter expression language.

The notification subsystem does not maintain a separate transport-attempt table. Once notification intent becomes a durable message, messaging owns provider/channel delivery history and retry behavior.

<!-- ORIGINAL_SOURCE_END: 11_NOTIFICATIONS.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 12_MESSAGING.md -->

# Original Source File: `12_MESSAGING.md`

# Messaging

## Purpose

yFramework provides a channel-independent outbound messaging subsystem. Components create typed message intent through yFramework and never communicate directly with SMTP servers, SMS providers, push providers, or similar delivery systems.

Remote SMTP email is the required v1 delivery channel. SMS and other channels are future extensions that must fit behind the same messaging boundary without changing component-facing semantics.

## Responsibility Boundary

### Components Own

- The domain reason for sending a message
- Stable message purpose
- Domain data needed to render message content
- Whether an operation should request a message

Components do not own SMTP configuration, credentials, provider SDKs, transport retries, or delivery state machines.

### yFramework Owns

- Typed messaging API
- Message-purpose registration and startup validation
- Durable outbound message records
- Channel routing
- Email delivery through the configured SMTP adapter
- Future channel-adapter registration
- Retry classification, backoff, retry budgets, and circuit breakers
- Rate limiting
- Delivery status and attempt history
- Sender-identity policy
- Tracing, structured logging, redaction, metrics, and administration

### Applications Own

- Enabled channels
- SMTP endpoint and credentials
- Allowed sender identities and branding
- Provider selection for future channels
- Message-purpose policy overrides where supported

## Message Model

A message uses a channel-independent envelope. Exact Go types and schemas remain implementation details, but the model should represent:

```text
message_id
purpose
component
channel
recipient/contact reference
locale
sender_identity
subject nullable
content/template data
created_at
send_after
status
attempt_count
request_id nullable
trace_id
job_id nullable
workflow_id nullable
```

Suggested statuses:

```text
queued
sending
retry_wait
sent
failed
canceled
review_required
```

Channel-specific details live behind channel adapters rather than contaminating the base message model.

## SMTP Channel

V1 email delivery uses a remote SMTP server. Application deployment configuration supplies validated settings such as host, port, TLS mode, authentication, credentials, connection timeout, send timeout, and allowed sender identities.

Components request a stable sender identity when needed. They do not invent arbitrary From addresses or receive SMTP credentials.

Verification email, localization-ready email, and all other application email use the same messaging subsystem.

## Durable Delivery

Outbound messages are durable before delivery. Delivery is a post-commit side effect and cannot roll back an already committed domain operation.

```text
domain fact commits
  -> message created
  -> durable delivery becomes eligible
  -> channel adapter attempts delivery
  -> sent, retry_wait, failed, or review_required
```

## Retry, Retry Budget, and Circuit Breaker

Temporary transport or provider failures may retry with bounded exponential backoff and jitter. Permanent rejection fails without pointless retries.

Messaging also uses an aggregate retry budget per provider or configured SMTP integration. A large queue must not create a retry storm during an outage.

Repeated provider-level failures may open a circuit breaker. While open, new delivery attempts remain deferred rather than repeatedly contacting a dependency expected to fail. After the recovery interval, a controlled half-open probe determines whether the circuit can close.

Retry ownership must be singular and explicit. Provider-library retries, transport retries, and messaging retries must not multiply unintentionally.

## Failure Classification

At minimum, delivery should distinguish:

```text
connection_failure
authentication_failure
timeout
temporary_provider_rejection
permanent_provider_rejection
invalid_destination
rate_limited
provider_unavailable
```

Authentication or configuration failures should be surfaced operationally and may suspend attempts until configuration is corrected.

## Localization

Message purposes may resolve localized subject and body content through the localization subsystem. A message requesting `de-DE` should use a current German message localization when available rather than initiating ad hoc provider-specific translation inside the channel adapter.

## Security

The messaging service is not an unrestricted relay. Components register stable purposes and the framework may enforce allowed channels, sender identities, recipient rules, rate limits, and authorization by purpose.

Full recipient destinations, credentials, verification tokens, and sensitive message bodies are excluded from normal structured logs.

## Future Channels

SMS, push notification, and other delivery mechanisms are optional future extensions. They implement the channel adapter contract and reuse the same durable message envelope, resilience primitives, observability, and administrative model.

## Execution Through the Global Background Lane

Messaging owns durable message state, delivery attempt history, retry classification, retry budgets, circuit breakers, rate limits, and provider outcome handling. It does not own an independent worker pool in v1.

When a message becomes eligible for delivery, the delivery attempt executes through the single global yFramework background execution lane. Notification fulfillment creates messaging intent; it does not create another delivery executor. Retry eligibility changes durable message state, while the global background lane performs the next eligible attempt.

## V1 Messaging Database Shape

The conceptual schema is defined in `16_DATABASE_SCHEMA_DESIGN.md`.

`framework_messages` stores compact authoritative delivery state. Large subject/body content lives in `framework_message_content`, and provider/channel attempt history lives in `framework_message_attempts`.

Messaging owns delivery retries even though actual eligible delivery work executes through the single ordinary background lane. A job execution attempt and a provider delivery attempt are distinct scopes and must not be conflated or allowed to create multiplying retry loops.

The message claim path uses short MariaDB transactions and performs external SMTP/provider communication outside the transaction. Message content and contact destinations are sensitive and are not written to ordinary logs.

<!-- ORIGINAL_SOURCE_END: 12_MESSAGING.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 13_EXTERNAL_INTEGRATIONS.md -->

# Original Source File: `13_EXTERNAL_INTEGRATIONS.md`

# External Integrations

## Core Rule

Components never own transport to systems outside the application boundary. They express external needs through typed, domain-shaped ports. Application composition selects and binds adapters. External adapters use yFramework integration infrastructure for common transport, resilience, security, health, and observability behavior.

```text
component
  -> typed port
  -> application-selected adapter
  -> framework integration runtime
  -> external system
```

This rule applies to LLM access, messaging, remote APIs, public-web fetching/scraping, and future external providers.

## Ports and Adapters

A component port describes domain capability rather than provider mechanics. For example, a metadata component should prefer a `BookMetadataLookup` contract over direct knowledge of Google Books or Open Library HTTP endpoints.

Provider-specific adapters may live outside yFramework when they are application- or domain-specific. yFramework supplies general infrastructure and selected general-purpose adapters, not every third-party integration.

Do not create a generic service locator such as `external.Call`. Strongly typed capability interfaces remain the default.

## Shared Integration Runtime

The framework provides reusable primitives for adapters where applicable:

- Credential lookup and redaction
- TLS and connection policy
- Request timeout and context cancellation
- Bounded retry with backoff and jitter
- Aggregate retry budgets
- Rate limiting
- Circuit breakers
- Required/optional health and readiness state
- Trace and request correlation
- Structured logging and metrics
- Error normalization
- Administrative enable/disable controls
- Testing helpers and fake transports

Adapters may opt into only the primitives appropriate for their protocol.

## Retry Ownership

Exactly one layer should own a given retry policy. The framework must not unknowingly stack job retries, adapter retries, transport retries, and provider-SDK retries.

Retries require:

- Bounded attempts
- Retryable error classification
- Backoff and jitter
- Idempotency or retry-safe semantics where state may change

Aggregate retry budgets limit total retry pressure against a failing provider. Once the configured budget is exhausted, otherwise retryable work is deferred until capacity returns.

## Circuit Breakers

A circuit breaker protects the application and failing dependency from repeated futile attempts.

Recommended states:

```text
closed
open
half_open
```

Repeated qualifying failures open the circuit. Open circuits fail fast or defer work without contacting the provider. After the recovery interval, a limited half-open probe determines whether normal operation resumes.

Exact thresholds are adapter/deployment configuration rather than universal constants.

## Synchronous Versus Durable Operations

Not every external operation goes through a durable queue.

- Request-path reads may be synchronous when the caller needs the result immediately.
- Durable side effects such as outbound messages use durable delivery records.
- LLM inference uses the dedicated durable LLM request queue.
- Long-running or multi-step domain work may use jobs/workflows.

Semantics determine the coordination mechanism.

## Outbound API Transport

The framework should provide a reusable outbound API transport for adapters. It owns mechanics such as:

- HTTP client lifecycle and connection pooling
- TLS policy
- Authentication injection through credential references
- Timeouts and cancellation
- Body-size limits
- Rate limiting
- Retry classification and Retry-After handling
- Trace propagation where appropriate
- Response/error normalization
- Safe logging and redaction

Components normally depend on a domain port implemented by an adapter using this transport rather than consuming the raw transport directly.

## Public Web Fetching and Scraping

Web-page consumption is a distinct integration category because public websites do not provide the same stability or semantics as APIs.

The safe web-fetch layer should support:

- HTTP/HTTPS scheme policy
- Redirect limits and redirect revalidation
- DNS and resolved-address checks
- Response-size and decompression limits
- Content-type validation
- Connection and total request timeouts
- User-Agent policy
- Rate limiting
- Character-encoding normalization where required
- HTML parsing helpers

Domain-specific extraction remains component or adapter logic.

### SSRF and Destination Policy

User- or data-influenced URLs must not create unrestricted access to internal resources. The web-fetch layer must enforce outbound destination policy across initial resolution and redirects.

At minimum, arbitrary-public-web mode should reject prohibited destinations such as loopback, link-local, private/internal network ranges, and other application-defined restricted ranges unless an explicit trusted integration permits them.

Known integrations should prefer configured host allowlists.

## Credentials

Components never receive raw external-provider credentials. Adapters refer to application-configured credential identities and obtain required secret material through the framework integration boundary.

The design should permit future environment-backed or external secret stores without changing component APIs.

## Health and Readiness

Each configured external integration is classified as required or optional.

- Failure of a required integration may make the application unready or degraded according to application policy.
- Failure of an optional integration reports the capability unavailable without necessarily making the whole application unready.

Health checks must avoid excessive provider traffic and may use circuit-breaker state and recent operational evidence where appropriate.

## Error Model

The framework normalizes common operational categories without discarding internal provider diagnostics:

```text
timeout
unavailable
rate_limited
authentication_failed
permission_denied
invalid_response
temporary_failure
permanent_failure
destination_blocked
```

Components receive stable typed errors. Sanitized provider-specific details remain available to logs and administration.

## Registration and Startup Validation

Components declare required typed capabilities during registration/composition. The application binds adapters before startup validation completes.

Startup rejects missing required bindings, invalid provider configuration, unsupported transport combinations, and malformed destination policies. Optional integrations must be represented explicitly as optional rather than silently missing.

## Testing

Component tests use fake port implementations and do not require real external systems.

Adapter tests verify provider/protocol translation separately. Integration tests may exercise controlled external or local test services. The framework should provide deterministic fakes for common transport and resilience behavior, including rate limits, circuit states, transient failures, and retry-budget exhaustion.

## Retry Ownership Matrix

Every logical operation has one retry owner. Durable job/workflow execution retries are owned by the job subsystem; messaging transport retries are owned by messaging; LLM request retries/crash handling are owned by the LLM subsystem; synchronous foreground integration retries are owned by the selected integration adapter/runtime.

Provider SDK and transport-library retries must be disabled, bounded, or explicitly incorporated into the owning policy so attempts do not multiply accidentally. Nested retries are permitted only when they represent clearly different logical scopes. Retry budgets and circuit breakers control aggregate pressure and attempt eligibility; they are not additional retry loops.

<!-- ORIGINAL_SOURCE_END: 13_EXTERNAL_INTEGRATIONS.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 17_CONFIGURATION_SCHEMA.md -->

# Original Source File: `17_CONFIGURATION_SCHEMA.md`

# yFramework Configuration Schema

## Status

This document defines the settled v1 configuration model.

## Configuration Classes

1. Bootstrap/deployment configuration: TOML, required before MariaDB is available.
2. Secrets: supplied through deployment configuration in v1 and treated as sensitive values.
3. Dynamic framework settings: registered typed settings persisted in MariaDB after startup/migrations.
4. Component/domain configuration: component-owned typed static configuration or component-owned durable tables.

## Settled Requirements

- One application deployment TOML file is the primary deployment configuration source.
- TOML is parsed strictly into strongly typed Go structures.
- Unknown TOML sections and keys are errors.
- Deployment topology and process-level infrastructure are not runtime-dynamic.
- Settings requiring reconstruction of infrastructure take effect only after restart.
- MariaDB runtime settings exist only for registered settings.
- There is no untyped generic component configuration map.
- Secrets never appear in normal logs, validation output, configuration dumps, health output, panic formatting, or CLI diagnostics.
- CLI overrides remain narrow and bootstrap-oriented, primarily the configuration-file path.
- Environment variables do not automatically override every TOML field.
- Framework security protections are not exposed as casual production disable switches.

## Top-Level TOML Sections

Recommended v1 top-level sections:

- `[application]`
- `[server]`
- `[database]`
- `[database.foreground]`
- `[database.background]`
- `[logging]`
- `[cache]`
- `[background]`
- `[health]`
- `[auth]`
- `[llm]`
- `[messaging]`
- `[messaging.smtp]`
- `[external]`

Component-specific sections may be registered only through typed component configuration contracts.

## Application

Typical fields:

- `name` required, stable process/application identifier.
- `environment` optional operational label such as development, test, staging, or production.

Environment labels must not silently weaken correctness or security behavior.

## Server

Deployment-only fields may include:

- UNIX socket path.
- socket mode/ownership where supported.
- explicit public origin.
- request limits/timeouts where exposed.
- shutdown drain deadline.

Public browser origin is explicit and is never inferred from an untrusted reverse-proxy header or backend UNIX socket.

## Database

Bootstrap fields include MariaDB host/socket, port where applicable, database name, user, password, TLS/connect settings as supported.

- TCP and UNIX-socket connection configuration may both be supported.
- Contradictory connection settings fail validation.
- Supported MariaDB version is defined by the yFramework release, not user configuration.

### Foreground and Background Pools

The foreground and background pools are configured separately. Defaults are provided.

Validation includes:

- nonnegative connection limits,
- idle limits not exceeding open limits,
- bounded background pool capacity.

The migration/bootstrap connection is separate from both runtime pools.

## Logging

The bootstrap logger always writes to stderr.

Full logging configuration may define:

- log level,
- structured format,
- destination,
- file/retention settings where yFramework owns them.

Core yFramework does not hardcode rc.d, systemd, launchd, journald, syslog, or platform-specific log paths.

## Cache

The cache section controls framework cache-manager defaults such as enablement, capacity, and general TTL policy.

Subsystems may impose stricter TTL/invalidation rules, especially for security-sensitive state.

Durable correctness must never depend on cache contents.

## Background Execution

Deployment-only mechanics may include:

- poll interval,
- lease duration,
- heartbeat interval,
- shutdown drain deadline,
- cleanup batch size.

Validation requires heartbeat intervals to remain safely inside lease durations.

There are no v1 configuration switches for multiple ordinary worker pools, per-subsystem workers, or job priority queues. V1 has one ordinary executable background lane.

## Health

Health configuration controls exposure and readiness policy, not the truth of subsystem health.

Subsystems report capability health. Central aggregation determines liveness/readiness/degraded state according to application policy.

## Authentication

Authentication deployment configuration may include:

- WebAuthn RP ID,
- exact approved public origin(s),
- session lifetime,
- session idle timeout,
- recent-authentication window,
- safe cookie naming options.

Secure, HttpOnly, SameSite, cross-origin protection, and other core security semantics are framework requirements rather than convenient production-off switches.

## LLM

Deployment-level configuration may include:

- enabled/disabled,
- yllmd socket/connection,
- queue/connect/execution timeout defaults,
- crash-quarantine threshold,
- content retention default,
- cleanup batch size.

Purpose-specific behavior remains registered in code.

## Messaging

Messaging configuration defines framework messaging behavior and channel adapters.

V1 requires remote SMTP support. SMTP configuration may include host, port, username, password, default sender identity, and constrained TLS mode.

Future channels may add typed sections without changing the component-facing messaging API.

## External Integration Runtime

Framework-wide outbound defaults may include:

- connect timeout,
- request timeout,
- maximum response size,
- retry-budget defaults,
- circuit-breaker/rate-limit defaults where configurable.

Provider-specific accounts and domain mappings remain application/component adapter configuration.

## Runtime Settings

A runtime setting must be registered in compiled code with:

- stable name,
- Go type,
- default,
- validation,
- dynamic/non-dynamic classification,
- secret classification where applicable,
- activation semantics.

Runtime precedence for an explicitly dynamic setting is:

1. compiled default,
2. permitted deployment TOML override,
3. MariaDB runtime override.

Deployment-only settings have no MariaDB override.

## Validation Phases

Static validation happens before database mutation and includes TOML structure, required values, durations, enum values, filesystem/socket settings, DB configuration, pool bounds, known adapter selections, and cross-field constraints.

Runtime validation occurs after registration and database loading and includes dynamic settings, routes, widgets, schedules, permission/action references, and other registry-backed configuration.

A static configuration failure prevents migrations.

## Configuration Reference Artifact

The implementation documentation must provide a structured reference for every setting with:

- path,
- type,
- required/default value,
- secret classification,
- deployment-only/runtime-dynamic classification,
- validation constraints,
- restart requirement,
- description.

<!-- ORIGINAL_SOURCE_END: 17_CONFIGURATION_SCHEMA.md -->
