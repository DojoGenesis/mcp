---
name: implementation-prompt
description: Transform specifications into structured, self-contained prompts for autonomous implementation agents. A prompt is a commission — a formal request for craftsmanship.
---

# Implementation Prompt Skill

**Version:** 1.1
**Author:** Tres Pies Design
**Purpose:** Write clear, comprehensive, and self-contained prompts that enable autonomous implementation agents to execute high-quality work without asking questions.

---

## I. The Philosophy: The Art of Commissioning

A prompt to an autonomous agent is not a command; it is a **commission**. It is a formal request for a work of craftsmanship. The quality of the commission directly determines the quality of the work. A vague, incomplete, or ambiguous prompt invites confusion, rework, and failure. A clear, comprehensive, and well-grounded prompt is an act of respect for the builder's time and capability.

This skill transforms prompt writing from a hopeful guess into a deliberate and rigorous engineering discipline. By following this structure, we provide the agent with everything it needs to succeed — enabling it to work with precision, autonomy, and a deep understanding of the existing codebase.

---

## II. When to Use This Skill

- After a specification has been finalized and has passed the pre-implementation checklist
- When commissioning a new development task to an autonomous implementation agent
- When breaking down a large specification into smaller, manageable implementation chunks
- When preparing work for parallel track execution

---

## III. The Prompt Writing Workflow

### Step 1: Validate Specification Readiness

**Goal:** Ensure the specification is complete and implementation-ready before writing prompts.

**Actions:**
1. Verify the specification has passed the pre-implementation-checklist (if applicable)
2. Confirm backend grounding is complete (API contracts, data models, integration points)
3. Check that all architectural decisions are documented
4. Verify there are no missing dependencies or unclear requirements
5. Identify if the work should be split into multiple tracks or sequential chunks

**Key Insight:** Never write an implementation prompt from an incomplete specification. The prompt quality is directly determined by the specification quality.

### Step 2: Identify Track Boundaries (If Applicable)

**Goal:** If the specification is large, identify which track this prompt addresses and what its boundaries are.

**Actions:**
1. Review the parallel tracks decomposition (if applicable)
2. Identify the scope of THIS track: what does it include? What does it explicitly NOT include?
3. Confirm track dependencies: what must be complete before this track? What will others depend on?
4. Verify track integration points: what props, APIs, state shapes will this track produce or consume?

**Key Insight:** When working with parallel tracks, be crystal clear about what THIS track does and doesn't do. Explicit non-goals prevent scope creep.

### Step 3: Ground in the Codebase

Before writing, gather all necessary context. The agent has full access to the repository, so leverage this. Your primary job is to be an excellent librarian, pointing the agent to the right information.

**Actions:**
1. **Identify pattern files:** Find 2-3 existing files that demonstrate the desired structure, style, and patterns
2. **List files to read/modify:** Enumerate all files that will be touched by this track
3. **Identify files to create:** List new files that will be created
4. **Map integration points:** Define the props, APIs, state shapes that this track will consume or produce
5. **Reference backend grounding:** Link to the backend integration guide for API contracts and data models

**Key Insight:** Pattern files are force multipliers. Implementation agents work best when you point them to existing examples of the desired structure and style.

### Step 4: Write the Prompt Using the Template

Create a new markdown file (e.g., `prompt_[feature]_[track].md`) and fill out the template from Section IV. Be precise and thorough.

**Actions:**
1. Fill out each section systematically
2. Ensure requirements are step-by-step and specific
3. Make success criteria binary (yes/no, no ambiguity)
4. Explicitly state what NOT to do in constraints
5. Keep the objective to a single, clear sentence
6. Reference backend grounding even if not calling APIs yet

### Step 5: Review Against the Checklist

Before sending the prompt to the agent, review it against the quality checklist in Section V. Every item must pass. This is the final quality gate.

**Key Insight:** The quality checklist catches gaps before commissioning. A prompt that passes the checklist has a much higher success rate.

### Step 6: Commission and Monitor

**Goal:** Commission the prompt to the appropriate agent and monitor progress.

**Actions:**
1. Save the prompt to the appropriate location
2. Commission to the appropriate agent (strategic implementation vs. tactical implementation)
3. Monitor progress through success criteria
4. If parallel tracks, verify completion before commissioning dependent tracks

**Key Insight:** A well-written prompt enables autonomous execution. The agent should complete the work without needing clarification or back-and-forth.

---

## IV. Implementation Prompt Template

```markdown
# Implementation Commission: [Brief, Descriptive Title]

**Objective:** [A single sentence describing the high-level goal of this task.]

---

## 1. Context & Grounding

**Primary Specification:**
- [Link to the final specification document]

**Pattern Files (Follow these examples):**
- `[path/to/existing_file_1.tsx]`: Use this for component structure and styling.
- `[path/to/existing_file_2.go]`: Use this for backend API structure and error handling.

**Files to Read:**
- [List files the agent should read for context before starting work]

---

## 2. Detailed Requirements

[Step-by-step, unambiguous list of implementation requirements. Be ruthlessly specific.]

**Backend:**
1. In `[path/to/file.go]`, create a new function `[FunctionName]` that...
2. Add a new API endpoint `GET /api/v1/[resource]` that...
3. The endpoint must return a JSON object matching this interface: `[interface]`

**Frontend:**
1. Create a new component at `[path/to/component.tsx]` named `[ComponentName]`.
2. The component must fetch data from `GET /api/v1/[resource]`.
3. Render data following the styling patterns in the pattern file.

---

## 3. File Manifest

[Complete list of all files to be created or modified.]

**Create:**
- `[path/to/new_file_1.ts]` — [Purpose]
- `[path/to/new_file_2.tsx]` — [Purpose]

**Modify:**
- `[path/to/existing_file_1.go]` — [What changes]
- `[path/to/existing_file_2.tsx]` — [What changes]

---

## 4. Success Criteria

[Binary, testable criteria. Each must be evaluable as TRUE or FALSE.]

- [ ] The new `[ComponentName]` renders correctly at the `/page` route.
- [ ] Clicking [element] triggers a call to `GET /api/v1/[resource]`.
- [ ] The backend returns `200 OK` with the correct JSON payload.
- [ ] All new code is covered by tests with at least [X]% coverage.
- [ ] The build passes with zero errors and zero new warnings.

---

## 5. Constraints & Boundaries

[What the agent must NOT do.]

- **DO NOT** modify any files outside the File Manifest.
- **DO NOT** introduce any new third-party dependencies without explicit approval.
- **DO NOT** refactor existing code outside the scope of this task.
- **DO NOT** address [related feature] — it is out of scope.
- **DO NOT** change the existing API contract for [existing endpoint].

---

## 6. Integration Points

[Where this work connects to other tracks or existing code.]

- This component receives data from `[ContextProvider/Store]` via `[hook/prop]`.
- This endpoint is consumed by `[frontend component]` at `[path]`.
- The shared interface between this track and Track [N] is: `[interface definition]`.

---

## 7. Testing Requirements

**Unit Tests:**
- Test `[function/component]` with [scenarios].
- Coverage target: [X]%.

**Integration Tests:**
- Test the flow from [start] to [end].
- Verify [specific behavior].

**Edge Cases:**
- [Empty state]
- [Error state]
- [Maximum data]
```

---

## V. Quality Checklist

- [ ] **Is the Objective a single, clear sentence?**
- [ ] **Is the specification linked?**
- [ ] **Are there 1-3 relevant Pattern Files listed?**
- [ ] **Are the Requirements specific, step-by-step, and unambiguous?**
- [ ] **Is the File Manifest complete?** (every file mentioned in requirements appears in the manifest)
- [ ] **Are the Success Criteria binary and testable?**
- [ ] **Are the Constraints clear about what NOT to do?**
- [ ] **Are Integration Points documented with both sides defined?**
- [ ] **Are Testing Requirements specific** (not just "write tests")?
- [ ] **Does the prompt respect existing codebase patterns?**
- [ ] **Is the prompt self-contained?** (agent needs no clarification to begin)

---

## VI. Best Practices

### From Specification to Prompt: The Translation

The specification describes the **what** and **why**. The prompt describes the **how** and **where**. Your job is to translate strategic intent into tactical execution.

**Example Transformation:**

**Specification says:** "The desktop shell should have a three-column layout with two collapsible sidebars."

**Prompt says:** "Create a new component at `frontend/src/app/desktop/components/DesktopShell.tsx`. Use a CSS Grid layout with three columns: 240px (sidebar 1), 200px (sidebar 2), and 1fr (main content). Add local state for `sidebar1Collapsed` and `sidebar2Collapsed` using `useState`. Follow the layout pattern in `frontend/src/app/layout.tsx`."

### Pattern Files are Force Multipliers

Pointing the implementation agent to 2-3 existing files that demonstrate the desired structure and style is more effective than writing detailed style guides. Implementation agents learn by example.

**Good:** "Follow the component structure in `src/components/Header.tsx`"
**Bad:** "Use functional components with TypeScript, export as default, use Tailwind for styling..."

### Chunk Your Prompts

Break large features into smaller implementation chunks. A single prompt should ideally represent 1-2 days of focused work.

### The File Manifest Prevents Surprises

Explicitly listing every file to be created or modified sets clear expectations and makes it easy to verify completion. It also helps identify potential merge conflicts in parallel tracks.

### Success Criteria Must Be Binary

Success criteria must be answerable with yes or no, with no subjective judgment.

**Good:** "The component renders at the /desktop route"
**Bad:** "The component should look good"

### Constraints Prevent Scope Creep

Explicitly stating what NOT to do is as important as stating what to do. This prevents the implementation agent from "helpfully" implementing adjacent features that belong in other tracks.

### Backend Grounding Even When Not Calling APIs

Even if a track doesn't call backend APIs yet, reference the backend grounding document. This ensures the implementation agent doesn't make decisions that would conflict with the backend architecture.

---

## VII. Example: v0.0.31 Track 1

**Specification:** `docs/v0.0.31/v0_0_31_desktop_architecture.md`

**Key Elements:**
- **Objective:** "Replace Next.js router with React Router and create the main DesktopShell component"
- **Pattern Files:** `frontend/src/app/layout.tsx` (for provider setup), `frontend/src/app/desktop/page.tsx` (for entry point)
- **Requirements:** 5 numbered, specific steps (install React Router, create router config, wrap app, create shell, update page)
- **File Manifest:** 2 files to create, 3 files to modify
- **Success Criteria:** 5 binary checkboxes
- **Constraints:** "DO NOT implement sidebar or main content area content"
- **Backend Grounding:** "No direct backend integration. Refer to backend integration guide for future tracks."

**Outcome:** The agent completed Track 1 autonomously without clarification, enabling Tracks 2-4 to start in parallel.

---

## VIII. Related Skills

- **`release-specification`** -- For writing the specification that this prompt translates
- **`frontend-from-backend`** -- Ensures backend architecture is documented before frontend specs
- **`parallel-tracks`** -- Guides the decomposition of large specs into parallel tracks
- **`pre-implementation-checklist`** -- Validates specification readiness before prompt writing

---

## IX. The Vision

This skill is designed to create **autonomous execution**. When an implementation prompt is written following this process, the agent should be able to:

1. Understand exactly what to build
2. Know where to build it (file paths)
3. Follow existing patterns (pattern files)
4. Verify completion (success criteria)
5. Stay within boundaries (constraints)
6. Integrate cleanly (backend grounding)

The result is high-quality implementation with minimal back-and-forth, enabling parallel work and compounding velocity.

---

**Last Updated:** 2026-04-06
**Status:** Active
