---
name: seed-extraction
description: >
  Extract and document reusable patterns (seeds) from experiences. Use when
  reflecting on learnings, capturing a pattern you noticed, documenting a
  mistake that taught something, or formalizing an insight into a structured
  artifact. Also use when someone asks how to write a seed, what makes a
  good seed, or how to generalize from a specific experience.
triggers:
  - "extract a seed"
  - "plant a seed"
  - "document this insight as a seed"
  - "capture this pattern"
  - "write a seed from this experience"
  - "how do I write a seed"
---

# Seed Extraction

**Version:** 1.1
**Author:** Tres Pies Design
**Purpose:** Extract, document, and apply reusable patterns (seeds) from experiences.

---

## I. Philosophy: Every Experience Contains Seeds

A seed is a reusable pattern that emerged from experience -- not theory, not abstraction, but something that actually happened and taught something transferable. The practice of seed extraction is learning to see these patterns, name them precisely, and plant them where they'll grow.

Seeds compound. A seed planted today becomes a reference point next week, a skill component next month, and part of your institutional DNA next quarter. The compounding cycle -- experience > seed > skill > ecosystem -- is the engine of institutional learning.

But not every observation is a seed. The test is reusability: can someone else, facing a similar situation, apply this pattern and get value? If yes, it's a seed. If it's too specific to one context, it's a note. If it's too vague to act on, it's a platitude.

## II. When to Use This Skill

- **After completing a project or milestone** that produced learnings
- **When you notice a pattern repeating** across multiple situations
- **During memory compression** (Tier A > Tier B promotion)
- **When a mistake teaches something reusable** -- failures are often the richest soil
- **When a technique or approach works surprisingly well** and could apply elsewhere
- **When a reframe changes how you think** about a class of problems

**When NOT to use:** If the knowledge is a simple fact ("use port 3000 for dev") it's documentation, not a seed. If it's a multi-step process, it's a skill candidate, not a seed. Seeds encode *principles*, not *procedures*.

## III. The Seed Extraction Workflow

### Step 1: Identify the Candidate

Look for raw material in recent experiences. The richest sources are:

| Source | What to Look For |
|---|---|
| Decisions that worked well | What principle guided the decision? |
| Decisions that failed | What was the wrong assumption? |
| Patterns across multiple instances | What keeps repeating? |
| Insights that changed thinking | What do you see differently now? |
| Tensions or tradeoffs navigated | What framework helped you choose? |

**Questions to surface candidates:**
- What did I learn that I didn't know before?
- What pattern did I notice repeating?
- What would I tell someone facing this situation?
- What would I do differently next time?

### Step 2: Test for Reusability

Apply four tests before investing in documentation:

| Test | Pass | Fail |
|---|---|---|
| **General enough** | Applies across multiple contexts | "Use Mermaid.js for architecture diagrams" (too tool-specific) |
| **Specific enough** | Actionable when triggered | "Be thoughtful about design" (too vague) |
| **Grounded in experience** | Emerged from something real | "I think this might work" (unproven) |
| **Has a clear trigger** | You know when to apply it | "Use this... sometime?" (no trigger) |

**Decision point:** If a candidate fails any test, it's not ready. Either refine it or let it go. Not every observation deserves to be a seed.

### Step 3: Generalize the Pattern

This is where most value is created. Move from the specific experience to the reusable principle:

1. **Start with what happened** (the concrete experience)
2. **Ask "why did this work?"** (the mechanism)
3. **Abstract the mechanism** into a principle that applies beyond this one case
4. **Test the abstraction** by imagining 2-3 other contexts where it applies

**Example of generalization:**
- Experience: "We started the migration with the auth service because it had the least coupling"
- Mechanism: Starting with the lowest-coupling component reduces blast radius
- Pattern: "Decouple the most independent component first to build confidence and reduce risk"
- Name: "decouple-lowest-coupling-first"

### Step 4: Document the Seed

Write the seed with all 6 required fields:

| Field | Purpose |
|---|---|
| **Name** | Short, memorable, verb-object or descriptive compound |
| **Context** | When and where this insight emerged (the specific story) |
| **Pattern** | The reusable principle, generalized beyond the original context |
| **Trigger** | Specific signals, situations, or keywords that should surface this seed |
| **Anti-pattern** | What the opposite looks like -- sharpens the pattern by contrast |
| **Example** | The concrete experience that birthed the seed (proof it's grounded) |

**Naming convention:** Use verb-object format (e.g., "decouple-lowest-coupling-first," "check-assumptions-before-code") or descriptive compounds (e.g., "the-reframe-is-the-prize," "parallel-tracks-pattern"). Names should be memorable and self-documenting.

**Full Seed Template (for comprehensive documentation):**

```markdown
## Seed: [Name]

**Pattern:** [One-sentence description of the reusable insight]

**Origin:** [Where this came from - project, experience, date]

**Why It Matters:** [The value or benefit of applying this seed]

**Trigger:** [When to apply this seed]
- [Context or situation 1]
- [Context or situation 2]

**How to Apply:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Anti-pattern:** [What the opposite looks like]

**Example (From Origin):**
[Concrete example from the experience where this seed emerged]

**Example (Applied):**
[Concrete example of applying this seed in a new context]

**Related Seeds:**
- [Seed that complements this one]
- [Seed that contrasts with this one]

**Cautions:**
- [When NOT to apply this seed]
- [Common misapplications]

**Evidence:**
- [Instance 1: date, context, outcome]
- [Instance 2: date, context, outcome]
- [Instance 3: date, context, outcome]

**Metadata:**
- **Created:** YYYY-MM-DD
- **Last Applied:** YYYY-MM-DD
- **Usage Count:** [Number]
- **Status:** Active | Experimental | Deprecated
```

### Step 5: Connect to Existing Seeds

Review the seed in context of the broader garden:

- **Complements:** What existing seeds work well alongside this one?
- **Contrasts:** What seeds represent an opposing perspective?
- **Compounds:** Does this seed + another seed = a larger pattern?

If a closely related seed already exists, consider whether to merge, differentiate, or let one supersede the other.

## IV. Best Practices

### Ground Every Seed in a Real Experience

Seeds without a concrete origin story are theories, not wisdom. The example field is not optional -- it's the proof that the pattern works. If you can't point to a specific experience, the seed isn't ready to plant.

### Name for Memory, Not for Completeness

Seed names should stick in your head. "backend-grounding-first" is memorable. "ensure-that-backend-implementation-precedes-frontend-work-to-prevent-disconnected-interfaces" is documentation. Choose names you'd use in conversation.

### Sharpen with Anti-patterns

The anti-pattern field is underrated. Describing what the opposite looks like makes the pattern click. "Error-codes-over-strings" becomes sharper when you add: "Anti-pattern: Every error is a freeform string message; client has to regex-parse errors to handle them."

### Prefer One Strong Seed Over Three Weak Ones

Ruthless curation applies to seeds too. One well-documented seed with a clear trigger and grounded example is worth more than three vague observations. Quality over quantity.

## V. Quality Checklist

- [ ] All 6 fields present (Name, Context, Pattern, Trigger, Anti-pattern, Example)
- [ ] Name follows verb-object or descriptive compound convention
- [ ] Pattern is generalized beyond the specific experience
- [ ] Trigger contains at least 2 concrete signals or situations
- [ ] Anti-pattern shows what the opposite looks like
- [ ] Example is grounded in a real experience (not hypothetical)
- [ ] Metadata complete (type: seed, tier: B, planted date, domain)
- [ ] Connected to related seeds where applicable

## VI. Common Pitfalls

### Over-Abstracting

**Problem:** The pattern is so general it's meaningless. "Design systems carefully" is true but provides zero guidance.

**Solution:** Include the mechanism. "Design systems carefully" becomes "decouple-lowest-coupling-first" which tells you *how* to be careful.

### Hoarding Ungrounded Seeds

**Problem:** Planting seeds from theory or intuition that haven't been tested in practice. The garden fills with speculative ideas.

**Solution:** Apply the "grounded in experience" test. If the seed doesn't have a real example, keep it as a Tier A note until experience validates it.

### Skipping the Generalization Step

**Problem:** The seed is just a specific lesson: "We should have started with the auth service." This is a note, not a pattern.

**Solution:** Always perform the generalization: "What principle would have told us to start there?" The answer is the seed.

## VII. Example

**Scenario:** After a sprint, the team discovered they'd built three separate utility libraries with overlapping functionality because three different developers solved similar problems independently.

**Step 1 (Identify):** A failure pattern -- duplication caused by lack of shared awareness.

**Step 2 (Test):** General enough? Yes -- applies to any team with parallel workstreams. Specific enough? Yes -- has a clear mechanism. Grounded? Yes -- just happened. Trigger? Yes -- when multiple people work on similar problems.

**Step 3 (Generalize):** Experience: Three duplicate utility libraries. Mechanism: Parallel work without shared visibility of in-progress solutions. Pattern: Make work-in-progress visible before starting implementation.

**Step 4 (Document):**
- **Name:** surface-work-in-progress-early
- **Context:** Sprint where 3 developers independently built overlapping utility libraries
- **Pattern:** Before building a new utility or shared component, surface your intent to the team. Duplicate work is a symptom of invisible work-in-progress.
- **Trigger:** Starting work on any shared component, utility, or cross-cutting concern; multiple people working in the same domain
- **Anti-pattern:** "I'll just build it quickly myself" -- leads to N implementations of the same thing
- **Example:** Three developers built date formatting utilities in the same sprint. A 5-minute standup mention would have prevented 2 of them.

**Step 5 (Connect):** Related to "parallel-tracks-pattern" (coordinate parallel work) and "shared-infrastructure" seed (build once, reuse everywhere).

## VIII. Seed Categories

### 1. Architectural Seeds
**Pattern:** Design decisions and system structures
**Examples:** Three-Tiered Governance, Harness Trace, Context Iceberg, Agent Connect

### 2. Process Seeds
**Pattern:** Workflows and methodologies
**Examples:** Planning with Files, Backend-First Development, Dual-Track Orchestration, Compression Cycle

### 3. Decision Seeds
**Pattern:** Frameworks for making choices
**Examples:** 3-Month Rule, Cost Guard, Safety Switch

### 4. Wisdom Seeds
**Pattern:** Principles and philosophies
**Examples:** Beginner's Mind, Understanding is Love, Knowing When to Shut Up

### 5. Technical Seeds
**Pattern:** Implementation patterns and best practices
**Examples:** Surgical Context, Graceful Degradation, Semantic Compression

---

## IX. Reflection Practice

### Daily Reflection (5-10 minutes)

1. What worked well today?
2. What didn't work as expected?
3. What pattern did I notice?
4. What would I do differently?
5. Is there a seed here?

**Output:** 1-2 candidate seeds for deeper reflection.

### Weekly Reflection (20-30 minutes)

1. What patterns emerged across this week?
2. Which candidate seeds are actually reusable?
3. Which seeds did I apply this week?
4. Which seeds need refinement?
5. Which seeds should be deprecated?

**Output:** Refined seed library, updated usage counts.

### Monthly Reflection (1-2 hours)

1. Which seeds have proven most valuable?
2. Which seeds have I stopped using?
3. What new categories of seeds are emerging?
4. How has my seed library evolved?
5. What seeds should I share with others?

**Output:** Curated seed collection, reflection document.

---

## X. Seed Application Workflow

### 1. Recognize the Trigger
Does this situation match a seed's trigger? Check context, keywords, and problem patterns.

### 2. Retrieve the Seed
Search seed library by keyword, browse by category, or recall from memory.

### 3. Apply the Seed
Read "How to Apply" steps, adapt to current context, check cautions to avoid misapplication.

### 4. Reflect on Outcome
Document whether the seed worked, what the outcome was, and whether the seed needs refinement. Increment usage count and add new examples if successful.

---

## XI. Related Skills

- **memory-garden** -- Seeds are Tier B artifacts in the three-tier hierarchy
- **compression-ritual** -- Seeds are often discovered during compression
- **seed-library** -- Reference collection of the 10 Dojo Seed Patches

---

**Last Updated:** 2026-04-06
**Status:** Active
