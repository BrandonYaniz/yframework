# yFramework Consolidated Source: 01_CORE_RUNTIME

This consolidated project source preserves the original documents listed below verbatim. The wrapper headings and checksum manifest are consolidation metadata only.

## Included source documents

- `02_COMPONENTS_ACTIONS_ROUTES.md` — SHA-256 `a0db450142d575492258c28f20027df7eb6a7e5c3e7667d53cba3cc3195b0dcb`
- `03_RENDERING_WIDGETS_TEMPLATES.md` — SHA-256 `bb8e8a7d1e4630616b7ee12607abf99eb2f544371a89f1e2eb0af02088cd3b1d`
- `04_AUTH_TOKENS_SESSIONS.md` — SHA-256 `225bfeaef4f89c8a83585ef729199b94ca2fc32428114145c79d08eb2c17aace`
- `07_LOGGING_OBSERVABILITY.md` — SHA-256 `edc8d4fa05f43c9661b4049a81e74ff2325082fec561ff9978ecffc5206ca591`
- `08_HOOKS_EVENTS_COMPONENT_LIFECYCLE.md` — SHA-256 `2162893bbdb61609b75bad64612cf59d9362c609633b594d35ac006af43bd527`

---

<!-- ORIGINAL_SOURCE_BEGIN: 02_COMPONENTS_ACTIONS_ROUTES.md -->

# Original Source File: `02_COMPONENTS_ACTIONS_ROUTES.md`

# Components, Actions, and Routes

## Component Model

A component is a reusable Go package containing related application capabilities.

Example:

```text
components/
  metadata/
    component.go
    actions/
    repositories/
    jobs/
    workflows/
    widgets/
    events/
    hooks/
    templates/
    migrations/
```

A component registers its capabilities with framework registries during application startup.

## Naming Convention

Use noun first, then action:

```text
ArticleCreate
ArticleList
ArticleSearch
ArticleView
MetadataImportSubmit
MetadataImportProcess
```

`ArticleList` represents a collection. `ArticleListItem` represents one entry in that collection.

Avoid ambiguous names such as `Summary` when a more specific name is available.

Stable dynamic names should be namespaced:

```text
article.ArticleList
metadata.ISBNLookup
metadata.MetadataImportSubmit
```

## Typed Action Contract

Actions use concrete Go input and output types.

```go
type Action[I any, O any] interface {
    Execute(context.Context, I) (O, error)
}
```

Struct-based actions are the default. Typed function adapters may also be supported.

Internally, the runtime registry uses a non-generic wrapper so route records can dynamically reference registered actions by stable name.

## Query and Command Metadata

Queries and commands share the same execution contract but have different metadata.

- Query: reads or computes data without authoritative state mutation.
- Command: changes authoritative or user-owned state.
- Ingest: narrowly accepts untrusted or candidate data for staged processing without directly changing authoritative records.

Application-token access may invoke approved query and ingest actions. Direct commands require user authorization.

## Action Output

Action output must be presentation-neutral.

Do not return:

- HTTP response types
- HTML fragments
- JSON envelopes
- MCP protocol response objects

Return explicit public result structs. Renderers translate those results into protocol-specific output.

## Dependencies

Actions receive dependencies explicitly through constructors.

Good:

```go
func NewArticleList(repo ArticleRepository, auth Authorizer) *ArticleListAction
```

Avoid:

- Package globals
- Runtime service locators
- Hidden database access
- Global registry lookups from domain code

Components may call another component's domain functionality only through an explicitly injected typed Go interface. They do not invoke other components through the runtime action registry.

Framework infrastructure dependencies, including the LLM service, are also supplied explicitly through typed constructors or factories. A component that uses LLM functionality depends on a yFramework LLM interface, not on yllmd transport code, a socket path, or a concrete model implementation.

The component that owns data also owns the operations that protect and modify that data. Other components must not update its tables directly.

Registered actions are application entry points exposed through routes, JSON APIs, MCP, CLI commands, or other application-owned adapters. Dynamic action registration exists for dispatch at those boundaries, not as an internal component communication mechanism.

## Routes

Routes are application-owned configuration stored in MariaDB.

A route has stable identity independent of its current path.

Suggested structure:

```text
routes
  id
  parent_route_id
  name
  path
  enabled
```

```text
route_handlers
  id
  route_id
  method
  action_name
  action_kind
  renderer
  view
  layout
  access_policy
  success_behavior
  enabled
```

A logical route may have multiple method-specific handlers.

Example:

```text
/login
  GET  -> user.LoginView
  POST -> user.LoginSubmit
```

An unsupported method returns HTTP 405 rather than 404.

## Route Hierarchy

Hierarchy is explicit through `parent_route_id`.

Do not infer all hierarchy solely from path prefixes. Stable hierarchy is also used for widget-placement inheritance and future administrative relationships.

## Input Mapping

A route handler maps request data to a typed action input.

Supported sources should include:

- Path parameters
- Query parameters
- Request body
- Headers
- Fixed configured values
- Trusted request context
- Authenticated user identity
- Application identity
- Installation identity

Authenticated identity must never be accepted from ordinary client-controlled input.

Explicit mappings override simple inference.

Do not support executable expressions or scripts in route configuration.

Type conversion and validation happen before action execution.

## Startup Validation

At startup, route loading should reject:

- Unknown action names
- Unknown renderers
- Unknown permissions
- Missing required templates
- Invalid input mappings
- Duplicate or ambiguous route handlers
- Hierarchy cycles
- Invalid methods
- Invalid access-policy combinations

Routes are loaded and validated at startup. Live route reload is not required for the first implementation.

## Route Links

Templates and components should resolve links through stable route names or IDs rather than hardcoding public paths.

## Localizable Component Resources

Components that expose localizable domain content register explicit localization definitions. The framework must not inspect arbitrary tables or infer translatable fields.

A localization definition identifies:

- Stable component-qualified resource type
- Primary source locale or source-locale resolver
- Stable localization units within the resource
- Source extraction and contextual data
- Validation of generated localized units
- Component-owned persistence of localized results
- The point at which a source revision becomes eligible for localization

Localization units must have stable identity so only changed units need regeneration. Examples may include title, summary, SEO fields, and stable body blocks or paragraphs. Components may supply neighboring or broader context to the LLM while replacing only the units whose source fingerprints changed.

Generated localized content is derived domain data owned by the component. yFramework owns coordination and status, not one universal translation-content table.

Human editing of generated localized content is permitted. A localized unit may be marked human-modified and may be locked against automatic regeneration. If its source changes, a human-modified or locked unit is marked stale and review-required rather than overwritten automatically.

## External Service Ports

A component must not open network connections, speak provider protocols, load external-provider credentials, invoke SMTP, call yllmd directly, or use raw outbound HTTP to communicate with systems outside the application boundary.

When domain behavior requires an external capability, the component depends on a typed port that describes the domain need. For example:

```go
type BookMetadataLookup interface {
    LookupISBN(context.Context, string) (BookMetadata, error)
}
```

The application binds that port to an adapter. The adapter may use a yFramework transport primitive, such as the outbound API transport or safe web fetcher, but the component does not depend on provider-specific wire behavior.

External capability requirements participate in startup validation. Missing required bindings prevent startup. Optional bindings must be represented explicitly rather than discovered dynamically through a service locator.

## V1 Route Grammar and Matching

Route paths use only literal segments and named parameters using `{name}` syntax. Regex segments, optional segments, splats, wildcards, executable predicates, and arbitrary route expressions are outside v1.

Literal segments take precedence over parameter segments. Structurally identical parameterized route shapes are ambiguous regardless of parameter name and are rejected during startup. Paths are normalized to one canonical trailing-slash form; duplicate normalized routes are rejected. Route hierarchy remains explicit through `parent_route_id` and is not inferred from path structure.

Path matching occurs before method selection. A matched path with no enabled handler for the requested method returns HTTP 405. V1 permits exactly one enabled handler for a route/method pair and does not use handler priority or conditional competing handlers.

## Route Authorization

`access_policy` describes required identity presence only: `public`, `application`, `user`, or `application_and_user`.

Handlers may separately require zero or more stable registered permissions. All listed permissions are required. Routes reference permissions, never application role names. Unknown permissions and invalid identity/permission/action-kind combinations prevent startup.

Route authorization is a coarse gate. Sensitive protected operations enforce authorization again at the action/component boundary, and resource ownership, visibility, relationships, state, and other domain-specific rules remain component-owned typed Go logic.

## Route Input Mapping Contract

Route input mapping is declarative. Supported sources are path parameters, query parameters, request bodies, approved headers, fixed literal values, trusted request context, and trusted user/application/installation identity.

Mappings contain no scripts, SQL, database lookups, arbitrary transformations, or executable expressions. Explicit mappings take precedence over limited safe inference. Duplicate or ambiguous destination mappings are rejected at startup.

Mappings are validated against the frozen typed action input definition before readiness. Type conversion and structural validation occur before action execution. Authentication identity is never established from client-controlled input. Security-sensitive protocol headers such as Authorization, Cookie, Forwarded, and proxy identity headers are consumed by framework infrastructure rather than exposed as ordinary mapped input by default.

Framework-owned persisted mapping records should be relational and use the `framework_` table namespace.

<!-- ORIGINAL_SOURCE_END: 02_COMPONENTS_ACTIONS_ROUTES.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 03_RENDERING_WIDGETS_TEMPLATES.md -->

# Original Source File: `03_RENDERING_WIDGETS_TEMPLATES.md`

# Rendering, Widgets, and Templates

## Renderer Architecture

The common pipeline is:

```text
typed action result -> selected renderer -> protocol response
```

Likely built-in renderers:

- Web/HTML
- JSON
- Text
- None
- MCP adapter, either foundational or added as an extension

The framework should support generic default rendering plus optional component-specific formatters.

Avoid a giant domain-type switch in the framework.

## Shared Error Rendering

The framework error model should be rendered differently by each adapter:

- Web: page or form feedback
- JSON: stable structured error envelope
- MCP: protocol-appropriate structured error
- Text: concise terminal-safe representation

Internal causes and stack data must not be exposed publicly.

## Web Page Context

The web renderer wraps action output in a standard page context containing:

- Site/application information
- Current route and handler
- Request context
- User, application, and installation identity
- Page metadata
- Assets
- Messages and validation errors
- Action output as `.Content`

## Template System

Use Go `html/template`.

Template levels:

```text
layouts
views
partials
```

Resolution order:

1. Application override
2. Application template
3. Component template
4. Framework fallback

Component templates should normally be embedded using `go:embed`. Applications may override them without forking the component.

View selection:

1. Route-handler-specific view
2. Action default view
3. Error if no compatible view exists

Layout selection:

1. Route-handler-specific layout
2. Application default layout

## Partials Versus Widgets

A partial renders data already supplied by the current rendering context.

A widget is a registered, configurable display capability that may prepare its own display data.

Widgets are display-only. They do not own POST processing or authoritative writes.

A form widget renders a form that submits to a route and action.

## Widget Model

Three concepts:

1. Widget definition: compiled Go functionality
2. Widget instance: configured use stored in MariaDB
3. Widget area: named placement region

A widget definition includes:

- Stable name
- Config schema
- Defaults
- Validation
- Display-data preparation
- Default template
- Assets

Template-selection order:

1. Definition default
2. Widget-instance override
3. Placement override
4. Inline override

Use stable template names, not arbitrary filesystem paths.

## Widget Storage

Suggested tables:

```text
widget_instances
widget_areas
widget_placements
widget_placement_routes
widget_placement_actions
widget_placement_handlers
```

A placement may include or exclude multiple routes. Exclusions win.

Placement can filter by:

- Route
- Descendants
- Exact route handler
- Action name
- Action kind

Parent-route inheritance is additive. Child routes may add placements or suppress inherited placements through exclusion rules.

## Widget Resolution

For a named area:

1. Load enabled candidate placements.
2. Require global scope or a matching include.
3. Apply route exclusions.
4. Apply route-handler rules.
5. Apply action rules.
6. Sort by position, then stable ID.
7. Render through the central widget renderer.

Cache resolved placements, not rendered HTML.

## Direct Widget Use

In Go templates:

```text
{{area "sidebar"}}
{{widgetInstance "homepage-search"}}
```

Stored content should use a narrow directive parser, not full Go-template execution:

```text
{{area name="article-after"}}
{{widget instance="homepage-search"}}
{{widget name="article.Recent" count="5"}}
```

The framework must detect cycles, cap nesting depth, and cap the number of widgets rendered per request.

## Asset Management

Widgets and templates register required CSS and JavaScript with a central asset manager.

The manager should:

- Deduplicate assets
- Preserve dependency order
- Separate header and footer assets
- Prevent arbitrary repeated script emission

## Localization Rendering State

Renderers may receive framework localization status in addition to component action output. The renderer does not perform localization itself.

When requested localized content is not yet current, the application may render framework-supplied state such as:

```text
missing
queued
generating
failed
```

The localization subsystem also supplies the already-localized operational message appropriate to the requested supported locale and, when available, a rounded estimated-ready value and estimate confidence. Applications decide whether to wait, show the status page/message, offer notification enrollment, or expose another application-approved presentation path.

Renderers must not trigger duplicate localization work merely because the same resource appears through multiple views or widgets. Localization requests are deduplicated by the localization subsystem.

A notification-enrollment form is presentation only. Submission goes through a normal route and typed action into the framework notification subsystem. Templates and widgets must not send email or SMS directly.


<!-- ORIGINAL_SOURCE_END: 03_RENDERING_WIDGETS_TEMPLATES.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 04_AUTH_TOKENS_SESSIONS.md -->

# Original Source File: `04_AUTH_TOKENS_SESSIONS.md`

# Authentication, Tokens, Sessions, and Authorization

## Identity Model

The framework should support a shared request identity that is independent of web, API, or MCP adapters.

```go
type RequestIdentity struct {
    User         *UserIdentity
    Application  *ApplicationIdentity
    Installation *InstallationIdentity
}
```

Authentication answers who or what is calling. Authorization answers what that identity may do.

Identity is stored in `context.Context`, not copied from untrusted action input.

## Trusted Internal Actors

Jobs, workflows, hooks, events, and typed cross-component calls must not manufacture a user
identity from ordinary payload data. When internal work acts on behalf of an initiating
request, the framework stores and reconstructs a limited trusted actor record containing the
relevant user, application, installation, authentication, and authorization context.

Authorization remains explicit at the operation that changes protected state. Framework
worker status does not automatically grant domain permission, and a component must not treat
all internal calls as administrative.

For system-originated work, the application or component uses a distinct, auditable system
actor with narrowly defined capabilities.

## User Tokens

User tokens represent an authenticated user and may be used through:

- Web
- API
- MCP

User tokens may be:

- Generated in advance by a user or system
- Generated after successful login

Useful token categories:

```text
personal
session
```

Personal tokens support integrations, MCP clients, scripts, and explicit API access.

Session tokens support browser login sessions.

Suggested fields:

```text
user_tokens
  id
  user_id
  token_type
  name
  token_prefix
  token_hash
  application_id nullable
  issued_at
  expires_at
  idle_expires_at nullable
  last_used_at
  token_family_id nullable
  rotated_from_token_id nullable
  replaced_by_token_id nullable
  revoked_at nullable
  enabled
```

Store only token hashes. Show raw tokens once.

## Application Tokens

Application tokens are created by an authorized user and always expire.

They identify an application or supported release generation.

Application tokens are not durable secrets when embedded in distributed applications. They are useful for:

- Supported-client control
- Release-generation identification
- Application-level rate limiting
- Read access
- Narrow ingestion submission
- Manual revocation
- Expiration-based retirement

Suggested records:

```text
applications
application_tokens
```

One application may have overlapping active token generations so users have time to update.

## Read and Write Boundary

Application tokens are read-only by default.

They may invoke:

- Approved query actions
- Approved ingestion-submission actions

They may not directly:

- Create authoritative records
- Update authoritative records
- Delete authoritative records
- Publish content
- Change user-owned data

Direct writes require a user token.

## Anonymous Ingestion Exception

Book Dragon may upload a local metadata library without login.

The application token and installation token may authorize:

```text
metadata.import-submit
metadata.import-status-own
```

The upload is stored as staged or transitional data. It does not directly mutate canonical metadata.

Trusted internal jobs and workflows later validate, transform, and optionally commit approved records.

## Installation Identity

Each installation receives a random identity without using hardware IDs or device fingerprinting.

Recommended model:

- Application token identifies app/release generation.
- Installation token identifies one installation.
- Both are required for installation-scoped anonymous ingestion.

First launch:

1. App creates a random local installation identifier.
2. App registers using its application token.
3. Server returns a server-issued installation token.
4. Server stores only the token hash.

The installation token is bound to:

- Application
- Token generation or compatible generation
- Installation record
- Expiration
- Narrow capabilities

Reinstallation may create a new installation identity.

## Access Policies

Useful route policies:

```text
public
application
user
application_and_user
```

For requests requiring both identities, effective permission is the intersection of application and user permissions.

An application token never elevates the user, and a user token never extends an application's allowed client capabilities.

## Authorization Layers

- Route access check provides the first gate.
- Sensitive actions enforce capability permissions internally.
- Resource ownership belongs in component action logic or a component authorizer.
- Roles group permissions.
- Actions generally check permissions rather than role names.

## Web Sessions

Use an opaque session token in a secure cookie:

```text
Secure
HttpOnly
SameSite=Lax or Strict
Path=/
__Host- prefix where applicable
```

The same user identity is produced whether the token came from a cookie or an Authorization header.

## Compromised Cookie Detection

Do not rigidly bind a session to a browser fingerprint.

Instead, maintain a limited session profile and risk score using signals such as:

- Major User-Agent change
- Operating-system change
- IP-prefix or network change
- Geographic impossibility
- Accept-Language change
- Concurrent distant use
- Abnormal timing or velocity
- Reuse of a rotated or revoked token

Minor changes should update observations. Moderate risk may trigger rotation or recent reauthentication. Strong replay evidence should revoke the token family and require full login.

## Token Rotation

Rotate session tokens:

- After login
- After privilege elevation
- Periodically during long sessions
- After moderate environmental risk

If an old predecessor token is used after rotation, treat it as a strong replay signal and revoke the family.

## Recent Authentication

Sensitive actions may require recent authentication:

- Token creation and revocation
- Administrative changes
- Account deletion
- Sensitive exports

Identity context should record authentication and reauthentication times.

## Token Creation Rules

Creating an application token requires explicit permission.

The creator selects delegated capabilities, but cannot delegate permissions they do not possess.

Application-token creation must structurally reject ordinary write permissions.

## Passkey-Only User Authentication

V1 user authentication is passkey-only using WebAuthn. Passwords are not supported. A second authentication factor layered on a password is not supported because there is no password flow. Account recovery is not provided in v1.

One framework user identity may have multiple independent passkey credentials, including credentials synchronized or managed by different platform ecosystems. Applications own the policy that authorizes creation of the initial user identity; yFramework provides the WebAuthn registration/authentication machinery once an enrollment context is authorized.

A user may not remove or disable their final usable passkey in v1. Adding or removing passkeys requires a valid authenticated session, normal browser request protection, and recent passkey authentication.

Passkey ceremonies use cryptographically unpredictable server-generated one-time challenges and validate the expected RP ID and explicitly configured approved origin. Public origin/RP settings are never inferred solely from Unix-socket transport or untrusted forwarding headers.

## CSRF and Browser-Origin Protection

Passkeys protect authentication ceremonies but do not eliminate CSRF after a browser session cookie exists.

Unsafe browser requests authenticated by a session cookie are protected at the HTTP boundary using Go's cross-origin protection facilities or an equivalent framework wrapper based on Fetch Metadata and strict origin validation. Safe methods such as GET, HEAD, and OPTIONS must never perform authoritative state mutation.

Session cookies default to `Secure`, `HttpOnly`, `SameSite=Lax`, `Path=/`, and the `__Host-` prefix where applicable. SameSite is defense in depth rather than the sole CSRF defense. Same-origin is trusted by default; additional trusted origins require explicit application configuration.

Explicit bearer-token API requests are not required to present CSRF tokens merely because they mutate state. The protection decision follows the authentication mechanism rather than renderer type.

## API Authentication Protocol Boundary

V1 API authentication uses explicit opaque yFramework bearer-style tokens through the `Authorization: Bearer` header. Personal, application, installation, and other defined token categories remain yFramework credentials. V1 does not implement OAuth 2.x or OpenID Connect protocol flows. OAuth/OIDC may be added later when delegated third-party authorization is a demonstrated requirement.

## Authorization Model

Authorization is deny-by-default.

Components and applications register stable permissions. Application-owned roles group permissions for administration. Reusable component business logic checks permissions rather than role names. Users may receive permissions through roles and through direct user-specific grants; effective global user permissions are the union of both sources.

V1 does not support explicit deny grants or role inheritance. Application identities use delegated application capabilities rather than human roles. For `application_and_user`, effective permission remains the intersection of application and user authority.

Global permission is only the coarse layer. Component-owned resource authorization evaluates ownership, relationships, attributes, visibility, resource state, and other domain facts through typed Go logic. yFramework does not impose one universal ACL table and does not store executable authorization expressions in MariaDB.

Protected mutations must not rely exclusively on route-layer checks because the same domain capability may be reached through web, API, MCP, CLI, jobs, workflows, or typed internal calls.

## Identity Data Ownership

yFramework owns the minimum durable records necessary to establish, authenticate, authorize, revoke, and audit an identity. This includes framework user IDs, passkey credentials, sessions/tokens, applications/installations, roles, role assignments, and direct permission grants.

Applications and components own profile and domain data such as display names, avatars, biographies, preferences, organization-specific attributes, and business records. Such data references the stable framework user ID. A passkey proves control of a credential bound to the framework identity; it does not define a user's domain profile.

## Development Schema Clarifications

The conceptual v1 schema is defined in `16_DATABASE_SCHEMA_DESIGN.md`.

Browser sessions and personal bearer API tokens use separate framework tables because their expiry, recent-authentication, and security semantics differ. Framework-owned relational joins use compact internal numeric IDs where appropriate, while stable UUID identity crosses subsystem/component/API/cache boundaries.

Resolved authorization sets are cacheable, but authority changes require immediate invalidation. Users, registered application principals, and installations carry an authorization revision that changes when their effective role/permission/capability inputs change. TTL alone is not sufficient for security-sensitive invalidation.

Usage timestamps such as session `last_seen_at` and token `last_used_at` are throttled rather than updated on every request, except where a correctness-sensitive WebAuthn field requires immediate persistence.

A durable append-only framework security audit stream records passkey, session/token, role/permission, application/installation, and similar security-administration events without logging secrets or unrestricted sensitive content.

<!-- ORIGINAL_SOURCE_END: 04_AUTH_TOKENS_SESSIONS.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 07_LOGGING_OBSERVABILITY.md -->

# Original Source File: `07_LOGGING_OBSERVABILITY.md`

# Logging and Observability

## Machine-First Logging

Logs are intended for yLLMLog and other automated processors.

Use JSON Lines with stable field names and stable event names.

Human-readable messages may be included, but machine processing should rely on structured fields.

## Log Separation

Recommended files:

```text
application.log
access.log
database.log
security.log
audit.log
events.log
jobs.log
```

All files use the same common schema.

## Common Envelope

Common fields should include:

```text
timestamp
level
log
event
service
instance
request_id
trace_id
component
message
```

Frequently queried fields should remain top-level rather than being buried in a generic data object.

## Canonical Context Fields

Identity:

```text
user_id
application_id
application_token_id
installation_id
token_type
auth_method
```

Routing:

```text
route_id
route_handler_id
route_name
action
action_kind
renderer
view
layout
```

Request:

```text
method
path
status
duration_us
response_bytes
remote_ip
user_agent
```

Component and repository:

```text
component
repository
operation
record_type
record_id
cache_name
cache_result
cache_purge_mode
cache_purge_kind
purge_after
```

Database:

```text
query_name
transaction_id
rows_read
rows_written
duration_us
attempt
```

Errors:

```text
error_code
error_type
error_message
error_chain
retryable
stack
```

## Event Naming

Prefer stable names such as:

```text
request.completed
route.matched
action.failed
cache.hit
cache.purge.requested
cache.purge.completed
cache.purge_mode.changed
database.query.slow
authentication.failed
authorization.denied
session.rotated
session.replay_detected
job.started
job.claimed
job.heartbeat_failed
job.lease_lost
job.lease_expired
job.recovered
job.retry_scheduled
job.cancel_requested
job.canceled
workflow.step.completed
event.subscriber.failed
logging.archive.created
background.pressure.changed
background.claims.paused
background.claims.resumed
llm.request.queued
llm.request.submitted
llm.request.accepted
llm.request.started
llm.request.completed
llm.request.failed
llm.request.cancel_requested
llm.request.canceled
concurrency_lock.acquired
concurrency_lock.recovered
```


## LLM Context Fields

LLM activity should use stable structured fields where applicable:

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
```

Normal request, trace, component, action, job, workflow, user, application, and installation fields remain the authoritative ownership context and should accompany LLM events when available. This makes an LLM request traceable back to the foreground request, background job, workflow step, or system operation that created it.

Prompt text, message bodies, source documents, and generated output are not logged by default. Logging metadata about model routing and usage must not accidentally expose prompt content or secrets.

## Request and Trace IDs

Every incoming request gets a request ID.

A trace ID follows related work into:

- Actions
- Repositories
- Events
- Jobs
- Workflows
- Follow-up internal operations

Background work receives its own job or workflow ID while retaining the originating trace ID.

## Job Execution Context

Job logs should include stable fields where applicable:

```text
job_id
job_name
queue
workflow_id
workflow_step
worker_id
lease_expires_at
attempt
max_attempts
concurrency_key
cancel_requested
```

Never log a lease token. It is an ownership credential and must be redacted like other secrets.

Concurrency-lock ownership tokens and deployment credentials, including an inline MariaDB password, are also secrets and must never appear in normal logs, diagnostics, audit records, or configuration displays.

Lease recovery logs must distinguish:

- Expiration while `claimed`
- Expiration while `running`
- Retry scheduled
- Attempts exhausted
- Stale-worker update rejected
- Shutdown drain or cancellation

## Redaction

Never log raw:

- Authorization headers
- Cookies
- Set-Cookie headers
- Passwords
- Tokens
- Token hashes
- Secrets
- Private keys
- Full request bodies by default

The logger should centrally sanitize dangerous field names.

## Audit Logging

Audit logs should answer:

- Who acted
- What changed
- Which target changed
- When
- Through which request
- Whether the operation succeeded

Audit logging remains subject to secret and sensitive-value redaction.

## Daily Rollover

Each active log rolls over at the start of a new server-local day.

Example archives:

```text
access-2026-07-23.jsonl.zst
security-2026-07-23.jsonl.zst
```

Rollover sequence:

1. Flush and close the active file.
2. Rename to an uncompressed dated `.jsonl` archive.
3. Open the new active file immediately.
4. Compress the completed archive.
5. Verify compression.
6. Atomically rename the temporary archive.
7. Remove the uncompressed archive.

Use temporary names such as:

```text
access-2026-07-23.jsonl.zst.tmp
```

## Compression

Default to Zstandard (`.zst`) at a moderate level, approximately level 3.

Support gzip as an alternate if required, but Zstandard is the preferred format for fast machine processing.

## Timezone

Log rollover uses the system timezone.

Timestamps include an explicit UTC offset.

The implementation calculates the next local midnight using timezone rules rather than adding 24 hours.

## Recovery and Safety

At startup, inspect for:

- Completed uncompressed archives
- Interrupted `.tmp` archives
- Existing verified archives
- Date-mismatched active files

Recover safely without overwriting valid archives.

Do not archive empty logs.

Support emergency size-based rollover in addition to daily rollover.

## Retention

Retention is configurable by log type.

Security and audit logs may be retained longer than access or debug-oriented logs.

Deletion should itself produce a structured logging event.

## Writers

Use bounded asynchronous writers for ordinary logs.

Security and audit records may require stronger flush or durability behavior.

If a log file cannot be written, fall back to stderr and emit a clear application failure signal.

## Localization and Notification Observability

Stable localization events should include:

```text
localization.queued
localization.started
localization.completed
localization.failed
localization.review_required
localization.locked
localization.unlocked
localization.locale_promoted
localization.backfill_queued
```

Useful fields include resource type and ID, localization unit, source locale, target locale, source revision or fingerprint, preferred-kind, human-modified flag, locked flag, review-required flag, LLM request ID, and estimated-ready metadata when available.

Stable notification events should include:

```text
notification.contact.pending
notification.verification.sent
notification.contact.verified
notification.contact.expired
notification.subscription.created
notification.delivery.queued
notification.delivery.completed
notification.delivery.failed
```

Do not place raw email addresses or phone numbers in ordinary structured logs. Log notification contact/subscription IDs and masked destination data only where operationally necessary.

LLM throughput observations should record enough information to reconstruct estimates without logging prompts or model output. Useful fields include concrete model, accepted/start/first-output/completion timestamps, input tokens, output tokens, total tokens, derived generation rate, sample count, estimate confidence, and predicted wait where applicable.

## External Integration Observability

External integration logs should use stable fields where applicable:

```text
integration
adapter
provider
operation
external_request_id
attempt
retry_budget_remaining
circuit_state
rate_limit_state
endpoint_name
```

Do not log provider credentials, Authorization values, complete sensitive destinations, or arbitrary response bodies by default.

Useful stable events include:

```text
integration.request.started
integration.request.completed
integration.request.failed
integration.retry_scheduled
integration.retry_budget_exhausted
integration.circuit.opened
integration.circuit.half_open
integration.circuit.closed
integration.rate_limited
integration.health.changed
integration.destination_blocked
```

Readiness and health logs must distinguish required integrations from optional integrations so an optional provider outage does not appear equivalent to loss of a mandatory dependency.

## Bootstrap Logging Lifecycle

A minimal bootstrap logger writes to stderr and exists for the full process lifetime. It is the first logging facility available during startup and the final fallback after the full structured logger is closed.

The full structured logger initializes immediately after bootstrap configuration and logging settings are statically valid. Fatal startup failures are mirrored to stderr. During shutdown, the structured logger remains active as long as practical so database, socket, integration, and shutdown failures can be recorded. After it flushes and closes, the bootstrap logger records final logger-shutdown and process-termination diagnostics.

The framework does not hardcode syslog, journald, macOS Unified Logging, or operating-system-specific log paths. FreeBSD rc.d/daemon, Linux service managers, macOS launchd, or other deployment packaging are responsible for retaining or redirecting stderr.

<!-- ORIGINAL_SOURCE_END: 07_LOGGING_OBSERVABILITY.md -->

---

<!-- ORIGINAL_SOURCE_BEGIN: 08_HOOKS_EVENTS_COMPONENT_LIFECYCLE.md -->

# Original Source File: `08_HOOKS_EVENTS_COMPONENT_LIFECYCLE.md`

# Hooks, Events, Component Lifecycle, and Dependencies

## Hooks

Hooks are synchronous Joomla-like extension points.

They may:

- Inspect active processing
- Mutate allowed context
- Add validation
- Block an operation
- Return an error

Hooks run in priority order.

A hook failure may stop the active operation.

Hooks are appropriate before commit or during rendering when extension logic must participate in the current outcome.

## Events

Events are immutable, past-tense notifications of completed durable facts.

Examples:

```text
ArticlePublished
MetadataImported
ApplicationTokenRevoked
```

Events occur after successful commit.

Subscribers may fail without reversing the already committed original action. Failures are logged and collected.

Initial delivery may be synchronous and in-process, but event semantics remain post-commit.

Do not make core correctness depend only on an event subscriber.

Events should carry IDs and facts, not pointers to mutable domain objects.

## Hooks Versus Events

Use hooks when an extension must influence or block the active operation.

Use events when something has already completed and other capabilities may react.

Suggested services:

```text
hooks.Dispatch(...)
events.Publish(...)
```

## Plugins

A separate plugin runtime is not required initially.

A plugin can be treated as a component focused on:

- Hooks
- Event subscribers
- Widgets
- Render adapters
- Extensions to other components

## Component Lifecycle

Lifecycle states and coordination belong to the application runtime:

```text
Constructed
Registered
Validated
Starting
Ready
Stopping
Stopped
Failed
```

Components remain passive and expose only:

```go
type Component interface {
    Identifier() ComponentID
    Register(Registrar) error
}
```

Registration declares framework-facing definitions such as actions, widgets, jobs, workflows, hooks, event subscribers, migrations, configuration definitions, and structural dependencies. It performs no database access, migration execution, business operations, configuration loading, or background startup.

Repositories remain private component implementation details and are not registered with the framework.

Components should not start independent queue pollers for their jobs. They register handlers with the framework runtime.

For v1, components may not start independent background services. Recurring or asynchronous component work must enter framework coordination through a job, workflow, schedule, or event subscriber. A narrowly defined managed-service facility may be considered later if a concrete application requires one.

## Dependencies

Dependencies are explicit and serve composition validation and ordering. They do not provide runtime service lookup.

Reusable components normally depend on application-defined contracts:

```go
r.Dependencies().RequireContract(identity.Contract)
```

The application binds a selected component as the provider and injects the provider's typed Go API through the dependent component's constructor. Contract identifiers belong to the application or composition package that defines the contract, not to yFramework.

An exact component dependency may be declared when a component genuinely requires one specific implementation. This should be rare.

Avoid a generic service locator.

Cross-component interaction uses:

- Explicitly injected typed interfaces for direct functionality
- Events for post-commit reactions
- Hooks for intentional extension points
- Jobs for asynchronous durable work
- Workflows for durable multi-step coordination

Components do not invoke one another through the runtime action registry and do not directly modify tables owned by another component. The component owning data also owns the operations that protect that data.

Registered actions are application entry points exposed through routes, APIs, MCP, CLI commands, or other application-owned adapters. The action registry is a dispatch boundary, not a component communication bus.

Optional integrations are selected explicitly during application composition. When selected, the application supplies the typed dependency, and the component registers the corresponding structural dependency edge.

The framework LLM service is a framework infrastructure dependency rather than a cross-component call. Components that need LLM functionality receive a typed yFramework LLM capability during construction. They do not receive a yllmd socket path and do not import or call yllmd protocol code directly. The application controls whether LLM support is enabled and supplies the validated purpose-to-target configuration used by the framework service.

Dependency requirements must resolve to selected providers. Singular contracts resolve to exactly one provider. Registration order has no dependency meaning, dependency cycles prevent installation and startup, and migration ordering follows the resolved dependency graph.

Because v1 applications are compiled as complete binaries, compatibility is enforced through Go modules, interface satisfaction, compilation, tests, and migration validation. Runtime contract-version negotiation is outside v1.

## Application Composition

Applications compose named component bindings. Three identities remain distinct:

- Permanent component identifier
- Application-owned component instance name
- Application or composition-package contract or role

Conceptually:

```go
type ComponentBinding struct {
    Name      InstanceName
    Roles     []Role
    Component Component
}
```

Reusable components never hard-code application instance names.

## Component Sets

A component set is an immutable composition description containing bindings and required roles. The application builder may accept several sets, allowing yCMS, another composition package, and application-specific components to remain independent.

Returned binding and role slices are copies. A set cannot be altered after it is handed to the builder.

yCMS is one example of a composition layer. It returns a component set and does not create an application, own `main`, retain a global builder, register components, or start the runtime. The Book Dragon Metadata API may compose components directly. Future applications and composition layers remain equally valid.

Replacement of a default provider occurs inside the composition package before its set is returned. A typed factory receives the selected typed dependencies. Only the chosen provider is constructed and registered.

## Application Builder

The recommended v1 builder is small and single-use:

```go
builder := framework.NewApplicationBuilder(definition)
if err := builder.AddSet(componentSet); err != nil {
    return nil, err
}
app, err := builder.Build()
```

`Build()`:

1. Rejects reuse of a completed builder.
2. Collects bindings and required roles.
3. Validates identifiers and instance names.
4. Rejects duplicate instance names.
5. Rejects duplicate implementations where multiple instances are unsupported.
6. Resolves singular required roles to exactly one provider.
7. Calls `Register()` on every selected component.
8. Rejects duplicate qualified capability names.
9. Validates dependencies and cycles.
10. Validates embedded resource references.
11. Freezes registries.
12. Returns the constructed application.

Addition order is preserved only for deterministic diagnostics. It does not provide dependency injection or migration ordering. Independent validation failures should be collected into one structured report where practical.

Both install and serve use the same application-composition function with fully populated typed configuration. Installation collects proposed deployment configuration and constructs infrastructure before composing components. Static validation completes before database mutation.

Exact interface names and whether factories use interfaces or typed function adapters remain implementation details.

## Localization and Notification Events

Localization completion is a durable fact only after the component-owned localized data has committed successfully. A completion event may then trigger notification matching and other post-commit reactions. Notification delivery failures do not roll back localization or any other completed domain operation.

Components that need notification functionality receive a typed framework notification capability or publish durable completion facts that the notification subsystem can match. They never send SMTP, SMS, or other provider requests directly. Outbound notification delivery always passes through the framework messaging subsystem and its application-selected channel adapter.

Localization and notification administrative operations remain subject to normal authorization, auditing, and request identity rules.

## External Integration Dependencies

External integrations follow the same explicit dependency principles as cross-component contracts.

A reusable component declares a typed external-capability port. Application composition chooses and constructs the adapter, then injects it through the component constructor. The component never discovers adapters through a runtime service locator and never obtains provider credentials directly.

Adapters may be framework-provided when they are general infrastructure, such as the yllmd transport, SMTP messaging transport, outbound API transport, or safe public-web fetcher. Domain/provider-specific adapters may live in application or reusable integration packages when they are not general framework responsibilities.

Required external capabilities must resolve during startup validation. Optional capabilities must be explicitly modeled as optional. Required integration health may affect readiness, while optional integration failure reports degraded capability without necessarily making the application unready.

## Application Lifecycle Coordination

The application runtime owns formal startup/shutdown phases and the root cancellation context. Successfully acquired startup resources register cleanup so startup failure can unwind safely. Components remain passive and never implement independent service startup or shutdown loops.

Static application composition and component registration complete before MariaDB mutation. Runtime services and background coordination start only after migrations and database-backed runtime validation succeed.

<!-- ORIGINAL_SOURCE_END: 08_HOOKS_EVENTS_COMPONENT_LIFECYCLE.md -->
