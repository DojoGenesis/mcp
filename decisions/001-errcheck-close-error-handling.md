# ADR-001: errcheck Handling Strategy for Non-Actionable Close() Errors

**Date:** 2026-05-27
**Status:** Accepted
**Affects:** `internal/gateway/client.go`, `internal/fsutil/atomic.go`, future HTTP and file I/O code

## Context

golangci-lint runs with `errcheck` enabled (default config — no `.golangci.yml` exists). `errcheck`
flags every call whose return value is discarded, including `resp.Body.Close()` in `defer` statements
and `tmp.Close()` calls in error-return branches.

Two categories of unchecked `Close()` arose:

1. **`defer resp.Body.Close()`** — The HTTP response body is fully consumed by `json.Decode` or
   `io.ReadAll` before the deferred close fires. For SSE streams, the stream is fully read by
   `parseSSE` before the function returns. In both cases the close error is non-actionable: it
   surfaces no information the caller can use, and there is no recovery path.

2. **`tmp.Close()` in error-return branches** — In `fsutil.AtomicWriteFile`, two error paths call
   `Close()` solely to release the OS file descriptor before returning. The original `Write` or `Sync`
   error is what the caller receives; a secondary close error on an already-failed handle is
   irrelevant.

The question is how to suppress these warnings consistently without hiding real errors.

## Routes Considered

| Route | Description | Risk | Effort | Key Tradeoff |
|-------|-------------|------|--------|-------------|
| A — Full error check | Log or propagate `Close()` error | Low | Medium | Noisy: defer requires a closure; non-defer mixes concerns in error paths |
| B — `_ = fn()` discard | Assign to blank identifier at call site | None | Minimal | Explicit at call site; doesn't work cleanly inside `defer` without a closure |
| C — `//nolint:errcheck` per-site | Annotate only the specific suppressed call | None | Minimal | Locates suppression context at the exact call; requires justification comment |
| D — `.golangci.yml` global exclusion | Exclude `(net/http.Response).Body.Close` etc. from errcheck globally | None | Small | Silences the entire function class; no per-site context; requires adding a new config file |

## Decision

**Use `_ = fn()` for non-defer cleanup paths; use `//nolint:errcheck // reason` for `defer` patterns.**

Route B (`_ =`) is chosen for non-defer error-path cleanup because it is idiomatic Go, does not
require a lint suppression directive, and signals intent clearly at the call site.

Route C (`//nolint:errcheck`) is chosen for `defer` patterns because `defer` cannot directly assign
to a blank identifier without a wrapper closure. A wrapper closure would add allocation overhead on
every HTTP call, is visually noisy, and provides no correctness benefit over the annotation.

Route D (global `.golangci.yml` exclusion) was not chosen because it removes the warning for all
callsites of those functions — including future callsites where the error genuinely matters (e.g.,
closing a file you've written, where `Close()` flushes buffered data and its error is real).

Route A (full error check) was not chosen for these specific sites because there is no recovery path
for a response body close error. Logging it would produce noise with no actionable signal; propagating
it would incorrectly shadow the actual HTTP or SSE error.

## What Was Not Chosen

- **Full error check (Route A):** The `defer resp.Body.Close()` pattern is idiomatic Go for HTTP
  handlers. Net/http guarantees the body is non-nil after a successful `Do`; the close error after
  full body consumption is structurally always nil on success paths. Adding a closure to surface it
  would be defensive theatre.
- **Global `.golangci.yml` exclusion (Route D):** Adds a new config file and sets a precedent of
  silencing checks globally rather than at the call site. Excluded for now; if the per-site annotation
  count grows past ~10, this tradeoff should be revisited.

## Consequences

**Easier:**
- Future contributors can follow the pattern by example: `_ =` for non-defer discards, `//nolint:errcheck // reason` for defer.
- golangci-lint CI annotations remain clean without suppressing the check globally.

**Harder:**
- Per-site annotations require a justification comment — reviewers should reject bare `//nolint:errcheck` without explanation.
- The absence of a `.golangci.yml` means lint configuration is implicit. If the default linter set changes in a future golangci-lint release, new violations may appear.

**Technical debt accepted:**
- No `.golangci.yml` exists. Running with defaults means the enabled linter set can drift between golangci-lint versions. A follow-up should pin the linter list explicitly.

## Propagation

- [x] `internal/fsutil/atomic.go` — lines 33, 37 converted to `_ = tmp.Close()`
- [x] `internal/gateway/client.go` — lines 245, 302, 496, 526 annotated with `//nolint:errcheck // reason`
- [ ] Add `.golangci.yml` to pin linter set and optionally configure `errcheck` exclusions for `http.Response.Body.Close` if per-site annotation count grows — tracked as future work
- [ ] Future HTTP call sites in `gateway/client.go`: follow `defer resp.Body.Close() //nolint:errcheck // body fully consumed before Close fires`
- [ ] Future file I/O error paths: follow `_ = f.Close()` pattern for error-return branches
