---
name: release-specification
description: Write production-ready, A+ quality specifications for software releases. A spec is a contract, not a wishlist. Encodes the disciplined engineering practice of creating comprehensive, precise, actionable, and testable release specifications.
triggers:
  - "write a release specification"
  - "write a spec for this release"
  - "create a production-ready spec"
  - "spec out this feature"
  - "write the release spec"
  - "document this release with a specification"
---

# Release Specification Skill

**Version:** 1.2
**Author:** Tres Pies Design
**Purpose:** Write production-ready, A+ quality specifications for software releases that autonomous agents can implement without asking questions.

---

## I. The Philosophy: Specification as Contract

A specification is not documentation — it is a **contract**. It is a formal agreement between the architect and the builder about what will be created, how it will work, and what success looks like. A vague specification invites confusion, rework, and failure. A rigorous specification is an act of respect for the builder's time and an investment in quality.

Specifications produced by this skill are:

- **Comprehensive:** Every question the builder might have is answered
- **Precise:** Technical details are specific, not hand-wavy
- **Actionable:** The path from specification to implementation is clear
- **Testable:** Success criteria are binary and measurable

**The standard:** 111/100 (A+). Good enough is not good enough.

---

## II. When to Use This Skill

Use this skill when:

- Planning a new software version or release with multiple features or components
- Designing a complex system architecture that requires detailed documentation
- Commissioning work to an autonomous agent that needs complete context
- Coordinating parallel development tracks where specifications serve as contracts between teams
- Communicating technical vision to stakeholders, developers, or future maintainers

Do NOT use this skill for:

- Small bug fixes or minor tweaks (use a simple task description)
- Exploratory prototypes (use scouting or rapid iteration)
- Features still being actively designed (finish discovery first)

---

## III. The Workflow

### Decision Point: Full Template or Lean Format?

Before starting, determine the right spec format:

**Use the Full Template (Section IV) when:**
- The system is new or the architecture is unfamiliar
- The audience includes stakeholders who need context
- Multiple teams or agents will implement from this spec
- The risk profile is high (production-critical, user-facing)

**Use the Lean Format when:**
- The architecture is established and the audience is the implementing agent
- The scope is well-defined (single feature, clear boundaries)
- The codebase patterns are well-documented
- The implementing agent is familiar with the codebase

**Lean Format structure:** Route layouts, component tables, behavior lists. No preamble. "Sonnet level chunks" — direct, precise, implementable.

**Rule:** The receiving agent should not default to the full template. Match format to scope.

---

### Step 1: Gather Context and Inspiration

Before writing, immerse yourself in the problem space:

1. **Read previous specifications** — Study 2-3 recent specs to understand the pattern and quality bar
2. **Review the codebase** — Understand the current architecture, types, and patterns
3. **Identify the problem** — What pain point, user need, or strategic goal is this release addressing?
4. **Scout alternatives** — If choosing between approaches, explore options before committing

**Output:** A clear understanding of the problem, the current state, and the desired future state.

---

### Step 1.5: Run Current State Audit

Before writing the spec, measure the codebase. Specs describe the delta from measured reality, not from assumptions.

**Run these audits against the target codebase:**

#### Testing
- Test file count: `find . -name "*.test.*" | wc -l`
- Test framework: [check package.json or equivalent]
- Coverage tool: [check scripts/config]

#### Accessibility
- Aria/role instances: `grep -r "aria-\|role=" --include="*.tsx" | wc -l`
- Error boundaries: `grep -r "ErrorBoundary" --include="*.tsx" | wc -l`

#### Performance
- Memoization usage: `grep -r "React.memo\|useMemo\|useCallback" --include="*.tsx" | wc -l`
- Code splitting: `grep -r "React.lazy" --include="*.tsx" | wc -l`

#### Dependencies
- UI framework: [check package.json]
- State management: [check package.json]
- Animation library: [check package.json]

#### File Structure
- Total source files: `find src -name "*.ts" -o -name "*.tsx" | wc -l`
- Route count: [check router config]

**Include the results as a "Current State" section at the top of the spec.** The spec then describes what changes FROM this measured baseline.

---

### Step 2: Draft Vision and Goals

Start with the "why" before the "what":

1. **Write a compelling vision statement** — A single sentence that captures the essence of this release
2. **Explain the core insight** — 2-3 paragraphs on why this release matters
3. **Define specific, measurable goals** — What will be different after this release?
4. **List non-goals explicitly** — What is out of scope for this release?

**Output:** A clear, inspiring vision that motivates the work and sets boundaries.

---

### Step 3: Design Technical Architecture

This is the heart of the specification:

1. **Create a system overview** — How do the major components fit together?
2. **Design each component in detail:**
   - Purpose and responsibility
   - Backend implementation with production-ready code examples
   - Frontend implementation with production-ready code examples
   - API endpoints with request/response shapes
   - Database schema (if applicable)
   - Integration points with existing systems
   - Performance considerations
3. **Write production-ready code examples** — Not pseudocode. Real code that could be committed.

**Output:** A complete technical design that a skilled developer could implement without asking questions.

---

### Step 4: Plan Implementation Phases

Break the work into manageable phases:

1. **Define a phased approach** — 2-4 phases with clear focus areas
2. **Create a week-by-week breakdown** — Specific, actionable tasks for each week
3. **Identify dependencies** — What must be done before other work can start?
4. **Define the testing strategy** — Unit, integration, E2E, performance, and manual QA plans

**Output:** A realistic timeline with clear milestones and success criteria for each phase.

---

### Step 5: Assess Risks and Document

Anticipate what could go wrong:

1. **Identify major risks** — Technical, timeline, or integration risks
2. **Define mitigation strategies** — How will you reduce or eliminate each risk?
3. **Plan rollback procedures** — How will you safely undo this release if needed?
4. **Define monitoring and alerts** — How will you know if something goes wrong in production?
5. **Document user and developer documentation needs** — What needs to be written?

**Output:** A comprehensive risk assessment and contingency plan.

---

### Step 6: Review Against Checklist

Before finalizing:

1. **Run the quality checklist** (Section VI) — Ensure every item passes
2. **Get feedback** — Share with a peer or stakeholder
3. **Iterate** — Refine based on feedback
4. **Save** — As `[version]_specification_[feature].md` in the workspace

**Output:** A finalized, A+ quality specification ready for implementation.

---

## IV. The A+ Specification Template

```markdown
# [Project Name] v[X.X.X]: [Memorable Tagline]

**Author:** [Name]
**Status:** [Draft | Final | Approved]
**Created:** [Date]
**Grounded In:** [What this builds on — previous versions, research, feedback]

---

## 1. Vision

> A single, compelling sentence that captures the essence of this release.

**The Core Insight:**

[2-3 paragraphs explaining WHY this release matters]

**What Makes This Different:**

[2-3 paragraphs on what makes this approach unique or better than alternatives]

---

## 1.5 Current State (Audit Results)

> Include measured codebase metrics here. This grounds the spec in reality.

**Testing:** [X] test files, [framework], [coverage tool]
**Accessibility:** [X] aria/role instances, [X] error boundaries
**Performance:** [X] memoization instances, [X] code splitting instances
**Dependencies:** [list key deps from package.json]
**File Structure:** [X] source files, [X] routes, [X] shared components

---

## 2. Goals & Success Criteria

**Primary Goals:**
1. [Specific, measurable goal]
2. [Specific, measurable goal]
3. [Specific, measurable goal]

**Success Criteria:**
- [ ] [Binary, testable criterion]
- [ ] [Binary, testable criterion]
- [ ] [Binary, testable criterion]

**Non-Goals (Out of Scope):**
- [What this release explicitly does NOT include]
- [What is deferred to future versions]

---

## 3. Technical Architecture

### 3.1 System Overview

[How components fit together]

**Key Components:**
1. **[Component Name]** — [Purpose and responsibility]
2. **[Component Name]** — [Purpose and responsibility]

### 3.2 [Component] — Detailed Design

**Purpose:** [What and why]

**Backend Implementation:**

```[language]
// Production-ready code with types, error handling, and imports
```

**Frontend Implementation:**

```typescript
// Production-ready code with interfaces, props, and state
```

**API Endpoints:**

| Method | Endpoint | Request Body | Response Body | Error Cases |
|--------|----------|-------------|---------------|-------------|
| POST | `/api/v1/resource` | `{ field: type }` | `{ id: string }` | 400, 401, 500 |

**Database Schema:**

```sql
CREATE TABLE table_name (
    id TEXT PRIMARY KEY,
    field TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Integration Points:**
- Integrates with [component] via [method]

**Performance Considerations:**
- [Specific constraint or optimization]

---

## 4. Implementation Plan

### 4.1 Phased Approach

| Phase | Duration | Focus | Deliverables |
|-------|----------|-------|--------------|
| 1 | Week 1-2 | [Focus] | [Deliverables] |
| 2 | Week 3-4 | [Focus] | [Deliverables] |

### 4.2 Week-by-Week Breakdown

**Week 1: [Focus]**
- [ ] [Specific task]
- [ ] [Specific task]

**Success Criteria:** [What "done" looks like for this week]

**File Manifest:**
- Create: `path/to/new_file.ts`
- Modify: `path/to/existing_file.ts`

### 4.3 Dependencies & Prerequisites

**Blocking:** [What must complete first]
**Parallel:** [What can run simultaneously]

### 4.4 Testing Strategy

**Unit Tests:** [Components], Target: [X]%
**Integration Tests:** [Scenarios]
**E2E Tests:** [User flows]
**Performance Tests:** [Metrics and targets]

---

## 5. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| [Description] | High/Med/Low | High/Med/Low | [Strategy] |

---

## 6. Rollback & Contingency

**Feature Flags:** `flag_name` — Controls [feature], default: false
**Rollback Procedure:** [Ordered steps]
**Monitoring:** [Metrics and alert thresholds]

---

## 7. Documentation

**User-Facing:** [What to update]
**Developer:** [API docs, schema docs, code examples]
**Release Notes:** [Changelog, breaking changes, migration guide]

---

## 8. Appendices

### 8.1 Related Work & Inspiration

- [Project/Paper]: [What we learned from it]
- [Tool/System]: [How it influenced this design]

### 8.2 Future Considerations

**v[X+1] Candidates:**
- [Feature that didn't make this release but is planned]
- [Enhancement that can build on this foundation]

### 8.3 Open Questions

- [ ] [Question that needs resolution before or during implementation]
- [ ] [Decision that can be made during development]

### 8.4 References

1. [Link to related spec]
2. [Link to research paper]
3. [Link to GitHub issue or discussion]
```

---

## V. Best Practices

### 1. Start with Vision, Not Features

**Why:** Features without vision are just a list of tasks. Vision provides meaning and direction.

**How:** Write the vision statement first. If you can't articulate why this release matters in one sentence, you're not ready to write the spec.

### 2. Write Production-Ready Code Examples

**Why:** Pseudocode leaves too much room for interpretation. Real code is a contract.

**How:** Write code examples that could be committed to the repository. Include imports, error handling, and types.

### 3. Use Realistic Timelines Based on Complexity

**Why:** Underestimating timelines leads to rushed work and technical debt.

**How:** Use past releases as benchmarks. A 1,000-line feature typically takes 1-2 weeks, not 2 days.

### 4. Document Integration Points Explicitly

**Why:** Most bugs happen at the boundaries between systems.

**How:** For every new component, explicitly document how it connects to existing systems (APIs, props, state, events).

### 5. Include Risk Mitigation from the Start

**Why:** Identifying risks after implementation is too late.

**How:** During the architecture phase, ask "What could go wrong?" and document mitigation strategies.

### 6. Make Success Criteria Binary and Testable

**Why:** Ambiguous success criteria lead to scope creep and endless iteration.

**How:** Every success criterion should be a yes/no question. "The user can create a new project" is testable. "The UI is intuitive" is not.

### 7. Reference Existing Patterns

**Why:** Consistency reduces cognitive load and makes the codebase easier to maintain.

**How:** When designing a new component, reference an existing component that follows the same pattern.

---

## VI. Quality Checklist

Before finalizing a specification, verify ALL items:

### Vision & Goals
- [ ] Vision statement is a single, compelling sentence
- [ ] Goals are specific, measurable, and achievable
- [ ] Non-goals are explicitly stated

### Technical Architecture
- [ ] Every major component has detailed design with code examples
- [ ] All API endpoints are fully specified (method, path, request, response, errors)
- [ ] Integration points with existing systems are documented
- [ ] Performance considerations are addressed
- [ ] All code examples are production-ready (real types, not pseudocode)

### Implementation Plan
- [ ] Timeline is realistic
- [ ] Week-by-week breakdown includes specific, actionable tasks
- [ ] File manifests list every file to create or modify
- [ ] Testing strategy covers unit, integration, E2E, and performance

### Risk & Documentation
- [ ] Major risks identified with mitigation strategies
- [ ] Rollback procedure documented
- [ ] User and developer documentation needs listed

**If you cannot answer "yes" to all items, the specification is not ready.**

---

## VII. Example: Dojo Genesis v0.0.23

**The Problem:** Users needed a way to calibrate the agent's behavior and communication style to match their preferences.

**The Vision:** "The Collaborative Calibration" — A release that transforms the agent from a fixed personality into an adaptive partner.

**What Made It A+:**

1. **Clear Vision:** The tagline "Collaborative Calibration" immediately communicated the essence
2. **Comprehensive Architecture:** Detailed design of the calibration UI, backend storage, and agent integration
3. **Production-Ready Code:** Complete Go and TypeScript examples that could be implemented directly
4. **Realistic Timeline:** 3-week phased approach with weekly milestones
5. **Risk Mitigation:** Identified the risk of "calibration drift" and defined a validation mechanism

**Key Decisions:**
- Store calibration preferences in SQLite for local-first architecture
- Use a multi-dimensional calibration model (tone, verbosity, formality)
- Implement real-time preview of calibration changes

**Outcome:** The specification was commissioned to an agent and implemented in 2.5 weeks with minimal rework.

---

## VIII. Common Pitfalls

| Pitfall | Fix |
|---------|-----|
| Vague Goals: "Improve user experience" | "Reduce context loading time by 50%" |
| Missing Code Examples: High-level description only | Complete, runnable code |
| Unrealistic Timelines: "2 days for full backend" | "2 weeks with phased approach" |
| No Risk Assessment: Assumes everything will work | Identifies risks and mitigations |
| Incomplete Testing: "We'll test it" | Specific test cases and coverage targets |
| No Integration Points: Treats feature as isolated | Documents how it connects to existing system |

---

## IX. Related Skills

- **`strategic-to-tactical-workflow`** — The complete workflow from scouting to implementation (this skill is Phase 6)
- **`strategic-scout`** — For exploring strategic tensions before committing to an approach
- **`frontend-from-backend`** — For frontend specs that need deep backend grounding
- **`implementation-prompt`** — For converting this spec into implementation prompts
- **`parallel-tracks`** — For splitting large specs into parallel execution tracks
- **`repo-context-sync`** — For gathering codebase context before writing specs
- **`memory-garden`** — For documenting learnings from implementation
- **`seed-extraction`** — For extracting reusable patterns from specs
- **`pre-implementation-checklist`** — For verifying specs are ready before commissioning

---

## X. Skill Metadata

**Token Savings:** ~10,000-15,000 tokens per specification (no need to re-read old specs for patterns)
**Quality Impact:** Ensures consistency across all specifications
**Maintenance:** Update when new patterns emerge from successful releases

**When to Update This Skill:**
- After completing 3-5 new specifications (to incorporate new patterns)
- When a specification leads to significant rework (to identify what was missing)
- When commissioning to a new type of agent (to adapt the template)

---

**Last Updated:** 2026-04-06
**Status:** Active
