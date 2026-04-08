---
name: pre-implementation-checklist
description: A checklist to ensure a smooth and successful handoff of a specification to an autonomous implementation agent. Use before committing a final specification.
---

# Pre-Implementation Checklist Skill

**Version:** 1.1
**Created:** 2026-02-04
**Last Updated:** 2026-04-06
**Author:** Manus AI  
**Purpose:** To provide a final quality gate before handing a specification to an autonomous development agent, ensuring all prerequisites are met for a successful implementation.

---

## I. The Philosophy of Preparing the Way

Handing a specification to an autonomous implementation agent is an act of trust. It is a request for a great work to be done. The Pre-Implementation Checklist is an **act of care** that precedes this request. It is the ritual of preparing the way, of ensuring that the path for the builder is clear, the materials are ready, and the destination is unambiguous.

By completing this checklist, we honor the builder's time and energy. We minimize the risk of confusion, rework, and failure. We transform the handoff from a hopeful toss over the wall into a deliberate and respectful commissioning.

---

## II. When to Use This Skill

-   **Always** use this skill as the final step before committing a `Final` specification document.
-   Use it as a final review gate before commissioning an implementation agent.
-   Use it to sanity-check a development plan before execution begins.

---

## III. The Checklist Workflow

### Step 1: Append the Checklist to the Specification

Once a specification is considered feature-complete and ready for final review, copy the checklist template from Section IV of this skill and append it to the very end of the specification document.

### Step 2: Collaboratively Review Each Item

Go through each item on the checklist. For each item, confirm that it has been addressed in the specification. Mark it with a `✅`.

-   If an item is not addressed, pause and update the specification accordingly.
-   If an item is not applicable, mark it as `N/A`.

### Step 3: Do Not Proceed Until All Items are Checked

The checklist is a gate. Do not proceed with the handoff to the implementation agent until every applicable item is marked with a `✅`. This discipline is crucial for maintaining a high standard of quality and ensuring successful autonomous implementation.

### Step 4: Commission the Implementation Agent

Once the checklist is complete, you can confidently commission the implementation agent, linking to the now-verified specification document.

---

## IV. Pre-Implementation Checklist Template

```markdown
---

## X. Pre-Implementation Checklist

**Instructions:** Before handing this specification to the implementation agent, ensure every item is checked. Do not proceed until the checklist is complete.

### 1. Vision & Goals

-   [ ] **Clarity of Purpose:** The "Vision" and "Core Insight" are clear, concise, and unambiguous.
-   [ ] **Measurable Goals:** The "Primary Goals" are specific and measurable.
-   [ ] **Testable Success Criteria:** The "Success Criteria" are concrete and can be objectively tested.
-   [ ] **Scope is Defined:** The "Non-Goals" clearly define what is out of scope for this release.

### 2. Technical Readiness

-   [ ] **Architecture is Sound:** The "Technical Architecture" is well-defined, and diagrams are clear.
-   [ ] **Code is Production-Ready:** All code examples are complete, correct, and follow existing patterns.
-   [ ] **APIs are Specified:** All API endpoints are fully specified (Method, Endpoint, Request, Response).
-   [ ] **Database Schema is Final:** The database schema is complete, including tables, columns, types, and indexes.
-   [ ] **Dependencies are Met:** All prerequisites and dependencies on other systems are identified and resolved.

### 3. Implementation Plan

-   [ ] **Plan is Actionable:** The week-by-week breakdown consists of specific, actionable tasks.
-   [ ] **Timeline is Realistic:** The timeline has been reviewed and is considered achievable.
-   [ ] **Testing Strategy is Comprehensive:** The testing strategy covers unit, integration, and E2E tests with clear targets.

### 4. Risk & Quality

-   [ ] **Risks are Mitigated:** Major risks have been identified, and clear mitigation strategies are in place.
-   [ ] **Rollback Plan is Clear:** A clear, step-by-step rollback procedure is documented.
-   [ ] **Feature Flags are Defined:** Necessary feature flags for a safe rollout are specified.

### 5. Handoff

-   [ ] **Final Review Complete:** This checklist has been fully completed and reviewed.
-   [ ] **Specification is Final:** The document status is marked as "Final".
-   [ ] **Implementation Ready:** You are now ready to commission the implementation agent.
```

---

## V. Best Practices

-   **Be Rigorous:** Do not check an item if it is only partially complete. The goal is 100% readiness.
-   **Treat it as a Conversation:** Use the checklist as a final opportunity to discuss any lingering uncertainties.
-   **Empower the Gatekeeper:** The agent responsible for the handoff (typically Manus) is empowered to halt the process if the checklist is not complete.
-   **Adapt as Needed:** If you find that items are consistently `N/A` or that new checks are needed, propose an update to this skill.

---

## VI. Example

**Problem:** The HTMLCraft Studio v1.0.0 specification was ready for handoff to the Zenflow implementation agent. The spec was 2,500 words covering the DIP scoring engine, Alpine.js template system, and KaTeX math rendering integration. Previous handoffs to Zenflow had experienced a ~30% rework rate due to ambiguous success criteria and missing dependency information.

**Process:**
1. Appended the Pre-Implementation Checklist template to the end of the HTMLCraft Studio v1.0.0 specification.
2. Reviewed each item collaboratively:
   - **Vision & Goals:** Checked. Purpose was clear ("self-contained HTML app builder"), goals were measurable ("408 tests passing"), success criteria were testable ("all templates render correctly in file:// protocol").
   - **Technical Readiness:** Found a gap -- the database schema section was missing the IndexedDB fallback strategy for `file://` protocol (IndexedDB is blocked on `file://`). Paused and added the localStorage fallback pattern.
   - **Implementation Plan:** Found another gap -- the testing strategy did not specify how to test the `@media print` stylesheet. Added a requirement for print preview screenshot comparison.
   - **Risk & Quality:** The rollback plan was missing. Added a simple rollback: "revert to the previous template directory if new templates fail validation."
   - **Handoff:** Marked all items as complete after addressing the 3 gaps.
3. Commissioned the Zenflow agent with the now-verified specification.

**Outcome:** HTMLCraft Studio v1.0.0 shipped with 408 tests passing and zero rework requests. The 3 gaps caught by the checklist (IndexedDB fallback, print testing, rollback plan) would each have caused a failed implementation attempt without the pre-implementation gate.

**Key Insight:** The checklist's value is not in confirming what is already there -- it is in forcing you to confront what is missing. The gaps found during checklist review are consistently the items that would have caused implementation failure.

---

## VII. Common Pitfalls

1. **Checking items without verifying them.** Marking "APIs are Specified" as complete without actually reading the API section to confirm it includes methods, endpoints, request bodies, and error codes.
   - *Solution:* For each checklist item, physically navigate to the relevant section of the spec and read it. The check mark means "I verified this," not "I believe this exists."

2. **Treating the checklist as a formality.** Running through the checklist quickly at the end, rubber-stamping everything, because the spec "feels done."
   - *Solution:* Schedule the checklist review as a distinct 15-30 minute activity, not a 2-minute addendum. The checklist review is where the last 10% of quality comes from.

3. **Not updating the spec when gaps are found.** Finding that the testing strategy is incomplete but checking the item anyway with a mental note to "fix it later."
   - *Solution:* The workflow is explicit: "If an item is not addressed, pause and update the specification accordingly." Do not proceed past a failed check.

4. **Using the same checklist for every type of spec.** Applying the full checklist (including database schema and rollback plan sections) to a documentation-only spec or a design system spec where those items are genuinely not applicable.
   - *Solution:* Mark items as N/A when they truly do not apply, but document why. If more than half the checklist is N/A, consider whether a simplified checklist variant is needed for that spec type.

5. **Skipping the checklist for "small" changes.** Deciding that a minor feature does not warrant the full pre-implementation checklist because "it is only a few files."
   - *Solution:* Scale the checklist, do not skip it. For small changes, a 5-minute abbreviated review (clarity of purpose, testable success criteria, dependencies met) still catches the most common handoff failures.

---

## VIII. Related Skills

- **write-release-specification** -- The upstream skill that produces the specifications this checklist validates. Use write-release-specification to create the spec, then this checklist to verify it.
- **parallel-tracks-pattern** -- When commissioning parallel tracks, apply this checklist to each track's specification independently before handoff.
- **health-supervisor** -- The health sprint commissions generated by health-supervisor should also pass this checklist before being sent to implementation agents.
- **write-frontend-spec-from-backend** -- Frontend specs are particularly prone to missing API contract details. This checklist catches those gaps before the frontend agent encounters them at build time.
- **pre-commission-alignment** -- A more comprehensive quality gate for multi-track commissions. Use pre-implementation-checklist for single specs, pre-commission-alignment for coordinated multi-spec handoffs.
