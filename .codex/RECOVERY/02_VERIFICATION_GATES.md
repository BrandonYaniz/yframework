# Verification Gates and Evidence Requirements

## Core principle

An exit code of zero from `go test` means the tests that ran did not signal a failure. It does **not** mean the documented requirements are met. A passed `t.Log("placeholder")` is not evidence; a skipped required test is a blocker; a method that silently returns `nil` is not an implementation.

For each requirement record **source citation by path/section**, **code under test**, **test name**, **assertion**, **failure injection/negative case**, **command**, **observed result**, and **commit**. Link the result in `.codex/STATE/recovery-verification.md`. No invented results or commit hashes.

## Universal verification

- `gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test -count=1 ./...` pass; `staticcheck` runs when available; `RACE=1 ./scripts/verify.sh` where applicable.
- `.codex/validate_task_graph.py` passes and completed/verified dependencies are consistent.
- No placeholder methods in *accepted scope*; deliberate future interfaces may be explicitly deferred but cannot satisfy present acceptance.
- No required `t.Skip`, `t.Log`-only tests, swallowed unexpected errors, or asynchronous tests that complete without joining goroutines.
- No new security vulnerabilities, leaked credentials, unauthorized access, non-idempotent retries, or silent recovery assumptions.
- No change to public/persisted contracts without settled source authority or documented human decision.

## Gate M1, Foundation (Tasks 001–008)

| Requirement | Mandatory evidence |
|---|---|
| Strict TOML | Parse valid file into typed config, reject unknown key/section and invalid type |
| Defaults | Assert exact documented default values |
| Cross-field checks | Test invalid pool bounds, heartbeat/lease, conflicting DB settings, and public-origin/RP ID constraints |
| Secrets | Assert passwords/tokens absent from config/log/error output |
| Structured logger | Decode actual JSON Lines, assert event/level/fields, filtering, and shutdown behavior |
| Bootstrap logging | Assert stderr works before configuration and after full logger shutdown |
| API consistency | Review exported APIs for source contract, no service locator or generic config maps |

**Fail closed** if a required case cannot be executed or its assertion is absent.

## Gate M2, MariaDB (Tasks 009–018)

Require an explicitly isolated, disposable MariaDB test database. The primary development database must remain untouched by destructive schema fixtures. Credentials are loaded from local protected config and never logged.

| Requirement | Mandatory evidence |
|---|---|
| Supported MariaDB | Verified version/capability query and expected acceptance |
| Unsupported MariaDB | Prove failure occurs **before** schema mutation |
| Migration from empty | Run real DDL and assert created schema/constraints via DB metadata |
| Repeat execution | Rerun, assert no duplicate effects and stable history |
| History identity | Two owners can register same migration name without collision, stable checksum stored |
| Checksum drift | Alter source registration checksum after initial success, assert fatal mismatch |
| Migration failure | Execute a deliberate failing operation, assert no success record and prior records intact |
| Advisory lock | Two independent DB connections cannot simultaneously hold lock/execute migration |
| Connection loss | Kill/close lock-owning connection, assert lock becomes available to another connection |
| Partial DDL | Test explicitly documented recoverable partial state; never claim full transactional rollback |
| Transaction | Real commit, rollback, canceled context, and error paths |
| Runtime pools | Separate foreground/background settings and cleanup on partial initialization failure |
| CLI validate | Assert checks/reads only, **no** schema/history mutation |

Run a gate command that **fails** when required tests are missing, skipped, or cannot connect. `go test -json` plus `scripts/check_required_go_tests.py` may be used for name/skip detection, but a human or a second independent audit must confirm tests assert real behavior. Don't make tests that only invoke `LoadConfig()` and print success.

## Gate M3, Composition (Tasks 019–026)

| Requirement | Mandatory evidence |
|---|---|
| Stable component ID | Valid/invalid names, duplicates, explicit dependencies |
| Graph | Exact dependency-first chain/diamond order, repeated deterministic output |
| Missing/cycles | Specific startup failure and no mutation |
| Registry | Duplicate rejection, freeze prevents writes, concurrent reads |
| Typed interface | No generic global service locator or public `GetService` bag |
| Composition | Successful minimal public-API consumer, deterministic registration |
| Failure atomicity | Invalid composition fails before DB mutation or runtime start |

## Gate M4 onward: operational completeness, original spec remains authoritative

After recovery, follow the exact milestone specifications in `.codex/MILESTONES/04.md` through `15.md` and the behavioral acceptance criteria in `docs/architecture/04_IMPLEMENTATION_CODEX.md`. In particular:

- **M4:** passkeys, account protections, sessions, API tokens, revocation, role/grant authorization, audit events.
- **M5:** typed actions, routes, UDS HTTP, reverse proxy/origin boundaries, authorization, CSRF, secure server-rendered templates, assets.
- **M6:** widgets and runtime settings, snapshot validation, rollback-on-invalid-candidate behavior.
- **Remaining milestones:** run their domain-specific integration, concurrency, recovery, idempotency, observability, and security acceptance gates. Do not claim full v1 completion on task counts alone.

## Gate result reporting template

```markdown
### Gate M2 at <commit> on <date>
- Environment: isolated MariaDB test database <nonsecret identifier>
- Required cases: <list>
- Executed: <list with test names and commands>
- Skipped/missing: <list, or none>
- Assertion review: <what database changes and failures were actually tested>
- Race/vet/staticcheck: <results>
- Defects: <links>
- Outcome: PASS or FAIL
- Release/next-task eligibility: <decision>
```
