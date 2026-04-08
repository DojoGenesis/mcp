---
name: debugging-troubleshooting
description: Systematic debugging and troubleshooting for code, systems, and workflows. Follows a 7-step methodology — reproduce, isolate, hypothesize, test, fix, verify, learn. Intuition is useful but methodology prevents rabbit holes.
---

# Debugging & Troubleshooting Skill

**Version:** 1.2
**Author:** Tres Pies Design
**Purpose:** Systematic debugging and troubleshooting that finds root causes, not just patches symptoms.

---

## I. Philosophy: Systematic Diagnosis

Debugging is systematic investigation, not random guessing. One change at a time, test, observe, learn. Intuition is useful but methodology prevents rabbit holes.

**Core principle:** You can't fix what you can't reproduce.

---

## II. When to Use This Skill

- Code behaves unexpectedly (errors, crashes, wrong output)
- Performance degradation (slow queries, long build times)
- Data inconsistency (state out of sync)
- Integration failure (API returns errors, sync breaks)
- Feature not working as expected (UI bug, logic error)

---

## III. The 7-Step Debugging Workflow

Always follow these steps in order. Never jump to fixing without reproducing.

### Step 1: Reproduce

Can the user reproduce the problem consistently?

- Identify exact steps to trigger the issue
- Run reproduction 3+ times
- Note if it happens every time or intermittently
- If intermittent, identify patterns (time of day, data size, user count, load)

**If you can't reproduce it, you can't fix it.** Invest time here.

**Information Gathering Template:**

```markdown
## Debug: [Issue Title]

**Gathered:** [Date/Time]

**Error Message:**
[Exact error text]

**Stack Trace:**
[Copy relevant portion]

**Logs:**
[Key log lines with timestamps]

**Context:**
- What triggered the issue:
- Recent changes:
- Environment:
  - Node: [version]
  - Database: [version]
  - OS: [platform]

**Expected Behavior:**
[What should happen]

**Actual Behavior:**
[What actually happened]
```

### Step 2: Isolate

Narrow the scope. Which component, file, function, or line?

**Bisection strategy:** When many variables exist, use binary search to narrow down:

1. Split the system in half
2. Test each half independently
3. The half that fails contains the bug
4. Repeat until you've isolated the specific component

**Variable isolation table:**

| Factor | Tested? | Result | Notes |
|--------|---------|--------|-------|
| [Component 1] | Yes/No | Pass/Fail | [Details] |
| [Component 2] | Yes/No | Pass/Fail | [Details] |

**Git bisect** for regression bugs:
```bash
git bisect start
git bisect bad          # current (broken) commit
git bisect good <hash>  # last known working commit
# Git tests commits between, finds exact breaking change
```

### Step 3: Hypothesize

Generate 3-5 possible causes ranked by likelihood:

| # | Hypothesis | Likelihood | Evidence For | Evidence Against |
|---|-----------|-----------|-------------|-----------------|
| 1 | [Most likely cause] | 70% | [What supports this] | [What contradicts] |
| 2 | [Second cause] | 20% | [What supports this] | [What contradicts] |
| 3 | [Edge case] | 10% | [What supports this] | [What contradicts] |

For each hypothesis: what evidence would **confirm** or **reject** it?

### Step 4: Test

For the most likely hypothesis, design a targeted test:

1. **Verify directly** — Add logging to confirm the cause
2. **Isolate variables** — Remove dependencies, simplify data, change environment
3. **Confirm causation** — Apply targeted fix, observe if problem resolves

**Document results:**
```markdown
Hypothesis: [Description]
Test Method: [How tested]
Result: Confirmed / Rejected
Evidence: [What proved/disproved it]
```

If rejected, move to next hypothesis. If all rejected, re-gather information.

### Step 5: Fix

Once cause is confirmed:

- Implement **minimal fix** (don't refactor everything)
- Add error handling or validation where appropriate
- Add monitoring or logging to catch future instances
- If the fix is risky, define a **rollback plan**

If **~~repository** is connected, suggest specific code changes.

### Step 6: Verify

Confirm the fix works without side effects:

- [ ] Original reproduction steps no longer fail
- [ ] Edge cases also work
- [ ] No regressions introduced (other features still work)
- [ ] Fix addresses root cause, not just symptom

### Step 7: Learn

Ask: **Is this a pattern that could recur?**

- If yes, offer to `/plant` a seed (wisdom-garden cross-reference)
- Document what happened for future reference
- Identify systemic improvements that would prevent similar bugs

---

## IV. Best Practices

### Read Logs Strategically

Don't read entire log files into context. Instead:
- Search for error codes or keywords
- Extract log lines around error timestamp (plus/minus 5 minutes)
- Look for recurring patterns

**Useful search commands:**
```bash
# Search for specific error
grep -i "error" app.log | tail -20

# Search around timestamp
grep "2026-02-02T14:00" app.log -A 5 -B 5

# Extract stack traces
grep -A 10 "Error:" app.log
```

### Use Logging to Verify

When in doubt, add a log. You can always remove later. Add logs at critical points:

| Point | Log What | Why |
|-------|---------|-----|
| Before database query | Query string, parameters | See what's executed |
| After database query | Rows returned, time taken | Performance check |
| Before API call | Request payload | Verify what's sent |
| After API call | Response status, body | Verify what's received |
| Before file write | File path, content | Verify what's written |

### Binary Search for Root Cause

When many variables exist, narrow down systematically instead of testing everything.

**Example:**
```
Database query slow (4s)

Test 1: Run query with 10 rows
Result: 0.1s (fast) -> Table size not cause

Test 2: Run EXPLAIN on query
Result: "seq scan" -> Missing index confirmed

Fix: Add index on queried columns
Result: 0.05s -> Confirmed
```

### Reproduce in Isolated Environment

When an issue only happens in production:
- Create minimal reproduction in dev
- Simplify data (use synthetic data)
- Remove external dependencies (mock APIs)
- Match environment versions

**Goal:** Isolate whether the issue is data, code, or environment.

### Version Control Time Travel

Use git history to find when something broke:
```bash
# What changed recently?
git log --oneline -10

# When did it start failing?
git log --oneline --since="2026-02-02T13:00"

# Bisect to find breaking commit
git bisect start
git bisect bad
git bisect good <last-working-commit>
```

---

## V. Quality Checklist

- [ ] Follows 7-step methodology in order — never jumps to fixing without reproducing
- [ ] Hypotheses ranked by likelihood with evidence
- [ ] Always ends with a learning step (is this a pattern that could recur?)
- [ ] Debug report includes prevention seed
- [ ] Uses bisection strategy for isolation when cause is unclear
- [ ] Minimal fix applied — no scope creep into refactoring

---

## VI. Example: rAF Batching Fix for SSE Reactive Hang (March 2026)

**The Problem:** The Dojo Gateway's SvelteKit dashboard froze when receiving Server-Sent Events (SSE) with skill execution traces. The UI became unresponsive during long-running skill executions.

**The Process:**

1. **Reproduce:** Triggered a skill execution with 50+ trace spans. Dashboard froze consistently after ~20 spans. Intermittent with fewer spans.
2. **Isolate:** Used bisection — disabled SSE processing (UI worked), disabled store updates (UI worked), re-enabled both (freeze). Isolated to: Svelte store updates triggered per SSE chunk.
3. **Hypothesize:** (1) Svelte reactivity triggering full re-render per store update (70%), (2) SSE backpressure overwhelming EventSource (20%), (3) Memory leak in trace accumulation (10%).
4. **Test:** Added `console.time` around store updates. Confirmed: each SSE chunk triggered a synchronous Svelte re-render. 50 chunks = 50 synchronous re-renders in rapid succession.
5. **Fix:** Implemented `requestAnimationFrame` batching — accumulate store updates and flush at frame rate (16ms intervals) instead of per-chunk.
6. **Verify:** Tested with 200+ spans. Dashboard remained responsive. No regressions in other SSE consumers.
7. **Learn:** Extracted seed: "rAF Batching for SSE" — any reactive framework (Svelte, React, Vue) will hang if you update stores per SSE chunk. Batch at frame rate.

**The Outcome:** Dashboard stayed responsive under 10x the original load. The seed was applied preventatively to 2 other SSE consumers in the codebase.

**Key Insight:** The hypothesis table was critical. Without it, we would have chased SSE backpressure (hypothesis 2) first, which would have been a rabbit hole.

---

## VII. Common Pitfalls

### Pitfall 1: Shotgun Debugging

**Problem:** Changing everything at once makes it impossible to know which change fixed the issue (or introduced new ones).

**Solution:** One change at a time, test, observe. The 7-step methodology enforces this discipline.

### Pitfall 2: Ignoring Error Messages

**Problem:** Assuming "probably network" or "probably a library bug" without reading the exact error message and stack trace.

**Solution:** Read the exact error. Search for it. Stack traces tell you exactly where the problem is.

### Pitfall 3: Fixing Without Verifying

**Problem:** Assuming the change works without re-running the reproduction steps.

**Solution:** Step 6 (Verify) is mandatory. Run the original reproduction steps after every fix.

### Pitfall 4: Assuming Environment Parity

**Problem:** "Works on my machine" — the dev environment doesn't match production.

**Solution:** Compare environment variables, dependency versions, OS, and configuration between working and broken environments.

### Pitfall 5: Refactoring During Debugging

**Problem:** Discovering messy code during debugging and starting a refactor instead of fixing the bug.

**Solution:** Step 5 says "minimal fix." File a separate ticket for the refactor. Debugging and refactoring are different activities with different goals.

---

## VIII. Related Skills

- **`seed-extraction`** — Step 7 (Learn) feeds into seed extraction when a bug reveals a reusable pattern
- **`health-audit`** — Proactive health audits can catch bug-prone code before debugging is needed
- **`research-modes`** — Deep research mode for investigating unfamiliar error conditions or library behaviors
- **`status-writing`** — Document ongoing debugging investigations in the project status
- **`retrospective`** — Significant bugs should be reviewed in retrospectives for systemic improvements

---

## IX. Troubleshooting Categories Reference

### Code Errors
**Common causes:** Null/undefined access, async timing, type mismatches, module not found
**Approach:** Read error message carefully, check stack trace, search error code, isolate function

### Performance Issues
**Common causes:** N+1 queries, unnecessary re-renders, memory leaks, network latency
**Approach:** Profile the slow operation, identify bottleneck, optimize specifically, measure before/after

### Data Inconsistency
**Common causes:** Database out of sync, cache invalidation, race conditions, unapplied migration
**Approach:** Compare expected vs actual state, check sync timestamps, verify data integrity

### Integration Failures
**Common causes:** API contract mismatch, auth failures, network connectivity, service downtime
**Approach:** Verify request format, check credentials, test API directly, check service status

### Environment-Specific Issues
**Common causes:** Missing env vars, path differences, dependency version conflicts, file permissions
**Approach:** Compare working vs broken environment, check env vars, verify dependency versions

---

## X. Output Format: Debug Report

```markdown
## Debug Report: [Issue Title]

**Date:** [Date]
**Reported By:** [Name]

### Symptom

[What was observed — exact error message, unexpected behavior]

### Reproduction Steps

1. [Step 1]
2. [Step 2]
Expected: [What should happen]
Actual: [What actually happens]

### Hypotheses Tested

| # | Hypothesis | Result | Evidence |
|---|-----------|--------|----------|
| 1 | [Hypothesis] | Confirmed/Rejected | [Evidence] |
| 2 | [Hypothesis] | Confirmed/Rejected | [Evidence] |

### Root Cause

[What actually caused the issue and why]

### Fix Applied

[What was changed to resolve it]

### Verification

- [ ] Original issue resolved
- [ ] Edge cases tested
- [ ] No regressions

### Prevention Seed

**Pattern:** [Is this a pattern that could recur?]
**Seed:** [Lesson to preserve for the wisdom-garden]
**Prevention:** [What systemic change would prevent this class of bug?]
```

---

## XI. Skill Metadata

**Token Efficiency:** ~8,000-15,000 tokens per debugging session (systematic approach, targeted testing)
**Quality Impact:** Ensures bugs are fixed at root cause, not symptoms; prevents regressions
**Maintenance:** Update when new debugging patterns emerge

---

**Last Updated:** 2026-04-06
**Status:** Active
