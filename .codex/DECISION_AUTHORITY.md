# Codex Decision Authority

Codex may make private implementation choices only when those choices do **not** alter documented public APIs, persisted semantics, security boundaries, ownership boundaries, lifecycle behavior, compatibility promises, failure/recovery semantics, or externally observable behavior.

Examples of delegated implementation choices include exact private helper types, ordinary index details that preserve documented query semantics, internal registry representation, private cache implementation, exact CLI exit-code numbering, and library selection where the architecture explicitly delegates it and the selected library satisfies the contract.

Codex must not silently decide unresolved architecture. When a required choice crosses the boundary above, add an entry to `STATE/open-decisions.md` with the affected task, source conflict/gap, alternatives, and why the choice is architecture-significant.
