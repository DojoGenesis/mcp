---
name: seed-library
description: >
  Manage and apply the 10 Dojo Seed Patches as reusable thinking modules.
  Use when architecting systems, making design decisions, debugging complex
  issues, optimizing costs, or when someone asks "what seed applies here?"
  or "which pattern should I use?" Also use when keywords match seed
  triggers: governance, multi-agent, context, routing, cost, safety,
  traceability, infrastructure, complexity gating.
---

# Seed Library

**Version:** 1.1
**Author:** Tres Pies Design
**Purpose:** Manage and apply the 10 Dojo Seed Patches plus 3 field seeds as reusable thinking modules.

---

## I. Philosophy: Seeds as Thinking Modules

The 10 Dojo Seed Patches are not abstract theory -- they are battle-tested patterns that emerged from real enterprise agent research. Each seed encodes a specific principle that, when applied at the right moment, prevents costly mistakes and accelerates delivery.

Seeds are thinking modules, not checklists. They change *how* you approach a problem, not just *what* you do. The meta-seed -- "Governance Multiplies Velocity" -- captures this essence: structure doesn't slow you down, it speeds you up by eliminating rework, ambiguity, and wasted exploration.

The library exists to make these patterns accessible. When you face a complex problem, the right seed should surface within seconds, not require 30 minutes of digging through documentation.

## II. When to Use This Skill

- **Architecting a new system or feature** and need proven patterns
- **Keywords match seed triggers:** governance, multi-agent, context, token, cost, routing, trace, safety, infrastructure, complexity
- **Someone asks** "What seed applies here?" or "Which pattern should I use?"
- **Debugging unexpected behavior** and need a framework for investigation
- **Optimizing costs or performance** and need a structured approach
- **Designing multi-agent coordination** and need routing patterns

**When NOT to use:** For seeds that emerged from your own experience (not the Dojo Seed Patches), use the seed-extraction skill. This library is specifically the curated collection of 10 enterprise-grade patterns.

## III. The 10 Dojo Seed Patches

### 1. Three-Tiered Governance

**Pattern:** Strategic / Tactical / Operational layers provide clear decision frameworks at every level.

**Trigger keywords:** governance, capabilities, complexity, coordination, policy, standards

**What it refuses:** Flat governance where everything is decided at one level.

**Apply when:** Building systems that need both flexibility and control.

### 2. Harness Trace

**Pattern:** Nested spans + events provide complete traceability of agent decisions.

**Trigger keywords:** debugging, trace, transparency, performance, evaluation, logging, monitoring

**What it refuses:** Opaque decision-making with no audit trail.

**Apply when:** Users or developers need to understand *why* an agent made a decision.

### 3. Context Iceberg

**Pattern:** Hierarchical context loading with 4 tiers. Budget for the 6x token multiplier from demo to production.

**Trigger keywords:** token, cost, context, window, limit, budget, pruning, memory, overhead

**What it refuses:** Loading everything into context and hoping for the best.

**Apply when:** Token usage is spiking or context management needs structure.

### 4. Agent Connect

**Pattern:** Routing-first, not swarm-first. Supervisor routes to specialized agents based on context.

**Trigger keywords:** multi-agent, routing, coordination, specialized, handoff, permission, swarm

**What it refuses:** Agent swarms where everyone talks to everyone.

**Apply when:** Designing multi-agent systems or agent coordination patterns.

### 5. Go-Live Bundles

**Pattern:** Reusable artifact packages that make sessions exportable, inspectable, and repeatable.

**Trigger keywords:** export, sharing, reuse, artifact, package, bundle, repeatability, trust

**What it refuses:** One-time outputs that can't be reproduced.

**Apply when:** Shipping features that users need to trust, inspect, or share.

### 6. Cost Guard

**Pattern:** Budget for the full iceberg (5-10x multiplier), not just API costs.

**Trigger keywords:** cost, budget, estimation, planning, infrastructure, investment, pricing

**What it refuses:** Budgeting only for direct API calls.

**Apply when:** Estimating costs for production agent systems.

### 7. Safety Switch

**Pattern:** Fallback to conservative mode when confidence drops or errors accumulate.

**Trigger keywords:** fallback, conservative, alert, drift, failure, recovery, validation, error

**What it refuses:** Continuing at full autonomy when things go wrong.

**Apply when:** Building error handling or designing graceful degradation.

### 8. Implicit Perspective Extraction

**Pattern:** Extract user intent from constraints and metaphors, not just explicit statements.

**Trigger keywords:** perspective, constraint, metaphor, scope, extraction, implicit, natural

**What it refuses:** Only acting on literal, explicit instructions.

**Apply when:** Designing user interactions or natural language interfaces.

### 9. Mode-Based Complexity Gating

**Pattern:** 3-question test to route queries to the right complexity level.

**Trigger keywords:** mode, complexity, routing, simple, query, reasoning, adaptive

**What it refuses:** One-size-fits-all processing regardless of query complexity.

**Apply when:** Building adaptive systems that handle both simple and complex requests.

### 10. Shared Infrastructure

**Pattern:** Build once, reuse everywhere. Common services reduce duplication.

**Trigger keywords:** infrastructure, reuse, duplication, foundation, shared, common, service

**What it refuses:** Each agent or feature building its own version of common functionality.

**Apply when:** Noticing duplication across services, agents, or features.

### Meta-Seed: Governance Multiplies Velocity

**Pattern:** Structure and governance don't slow you down -- they speed you up by eliminating rework, ambiguity, and wasted exploration.

**This is the seed that explains why seeds work.** Every pattern in this library exists because someone learned that the upfront investment in structure paid compound returns in velocity.

### Field Seeds (From Practice)

**11. Voice Before Structure** -- Read design language before writing structural artifacts.

**12. Pointer Directories** -- Empty directories are references, not gaps.

**13. Granular Visibility** -- Progress tracking serves the user, not the agent.

---

## IV. Trigger Keyword Reference

| Seed | Trigger Keywords |
|---|---|
| 01 Governance | governance, capabilities, complexity, multi-agent, coordination, policy, standards |
| 02 Trace | debugging, trace, transparency, performance, evaluation, logging, monitoring |
| 03 Context | token, cost, context, window, limit, budget, pruning, memory, overhead |
| 04 Agent Connect | multi-agent, routing, coordination, specialized, handoff, permission, swarm |
| 05 Bundles | export, sharing, reuse, artifact, package, bundle, repeatability, trust |
| 06 Cost Guard | cost, budget, estimation, planning, infrastructure, investment, pricing |
| 07 Safety Switch | fallback, conservative, alert, drift, failure, recovery, validation, error |
| 08 Perspective | perspective, constraint, metaphor, scope, extraction, implicit, natural |
| 09 Complexity | mode, complexity, routing, simple, query, reasoning, adaptive |
| 10 Infrastructure | infrastructure, reuse, duplication, foundation, shared, common, service |
| 11 Voice | voice, philosophy, design-language, manifest, description, ecosystem, grounding |
| 12 Pointers | empty, missing, pointer, provenance, registry, audit, gap, coverage, directory |
| 13 Visibility | progress, tracking, visibility, todo, granular, steering, trust, delegation |

---

## V. Seed Suggestion & Application Workflow

### Step 1: Gather Keywords

From the user's task or question, extract keywords that match seed triggers.

### Step 2: Match to Seeds

Compare keywords against the trigger keyword lists above. Rank by relevance:

- **Direct match:** Keyword appears in trigger list
- **Semantic match:** Keyword is closely related to a trigger
- **Context match:** The overall situation matches a seed's "apply when"

### Step 3: Present Top Recommendations

Suggest 1-3 seeds ranked by relevance. For each:
- Name and pattern (1 sentence)
- Why it's relevant to the current task
- Key "what it refuses" as a contrast

### Step 4: Apply the Chosen Seed

When the user selects a seed:
1. Explain the full pattern and its implications
2. Provide an application checklist specific to their context
3. Note the anti-pattern (what it refuses) as a guardrail
4. Suggest related seeds that complement the chosen one

## VI. Quality Checklist

- [ ] Seed suggestion matches the user's actual context (not just keyword matching)
- [ ] Top recommendation explains *why* it's relevant, not just *what* it is
- [ ] Anti-pattern ("what it refuses") is presented alongside the pattern
- [ ] Application checklist is tailored to the user's specific situation
- [ ] Related seeds are surfaced for complementary patterns
- [ ] Growth language used (seeds are "applied" and "cultivated," not "used" and "stored")

## VII. Common Pitfalls

### Applying Seeds Mechanically

**Problem:** Treating seeds as checklists to follow rather than thinking modules to internalize. Following the steps without understanding the principle.

**Solution:** Always start with *why* the seed exists. The "what it refuses" field is the key -- understanding what the pattern rejects is often more illuminating than understanding what it prescribes.

### Suggesting Too Many Seeds

**Problem:** Overwhelming the user with 5+ seed suggestions when they need focused guidance.

**Solution:** Suggest 1-3 seeds maximum. If multiple seeds are relevant, prioritize the one most directly applicable and mention others as "also relevant."

## VIII. Example

**Scenario:** A team asks "How should we handle it when our agent makes a mistake?"

**Step 1 (Keywords):** error, failure, recovery, fallback, agent

**Step 2 (Match):**
- **Safety Switch (Seed 7)** -- Direct match on fallback, failure, recovery
- **Harness Trace (Seed 2)** -- Relevant for understanding *why* the mistake happened
- **Three-Tiered Governance (Seed 1)** -- Framework for deciding escalation levels

**Step 3 (Present):**

> For handling agent mistakes, I recommend **Safety Switch** (Seed 7). This pattern establishes a fallback to conservative mode when confidence drops or errors accumulate. It refuses to continue at full autonomy when things go wrong.
>
> **Also relevant:** Harness Trace (Seed 2) for understanding *why* errors happen, and Three-Tiered Governance (Seed 1) for structuring escalation policies.

**Step 4 (Apply):** Application checklist for Safety Switch:
- Define confidence thresholds that trigger fallback
- Design the "conservative mode" behavior
- Build alerting for when the switch activates
- Create a recovery path back to full autonomy

## IX. Seed Relationships

### Foundational (Start Here)
- **Three-Tiered Governance** -- Framework for all other seeds
- **Meta: Governance Multiplies Velocity** -- The philosophy

### Operational (Day-to-Day)
- **Harness Trace** -- Transparency and debugging
- **Context Iceberg** -- Token and cost management
- **Safety Switch** -- Error handling
- **Mode-Based Complexity Gating** -- Query routing

### Architectural (System Design)
- **Agent Connect** -- Multi-agent coordination
- **Shared Infrastructure** -- Reusable services

### Delivery (Shipping)
- **Go-Live Bundles** -- Packaging and reuse
- **Cost Guard** -- Budgeting and planning

### UX (User Experience)
- **Implicit Perspective Extraction** -- Reduce friction

## X. Common Use Cases

### Architecting Multi-Agent System

**User asks:** "How should we architect the multi-agent system?"

**Workflow:** Suggest Agent Connect, Shared Infrastructure, Three-Tiered Governance. Apply Agent Connect first -- routing-first pattern where the Supervisor routes to specialized agents based on context.

### Optimizing Token Usage

**User asks:** "Token usage is spiking, how do we optimize?"

**Workflow:** Suggest Context Iceberg, Cost Guard, Mode-Based Complexity Gating. Apply Context Iceberg -- hierarchical context loading with 4 tiers and pruning triggers at 80%, 90%, 95%.

### Debugging Unexpected Behavior

**User asks:** "Why did the agent choose X instead of Y?"

**Workflow:** Suggest Harness Trace. Review trace spans for the decision point. Check metadata.reasoning field.

### Building Trust with Users

**User asks:** "Users don't understand our recommendations"

**Workflow:** Apply Harness Trace (transparency), Go-Live Bundles (repeatability), and the meta-seed (philosophy).

### Planning New Feature

**User asks:** "Let's plan the new feature"

**Workflow:** Apply Three-Tiered Governance to define strategic/tactical/operational layers, then Cost Guard to budget for the full iceberg.

---

## XI. Seed Maintenance

### When to Update a Seed
- Pattern evolves based on new learnings
- Better approach discovered
- User feedback suggests improvement

### When to Archive a Seed
- Pattern no longer applies
- Superseded by better approach
- Context has fundamentally changed

### Versioning Convention
- **1.0** -- Initial version
- **1.1** -- Minor update (clarification, example added)
- **2.0** -- Major update (pattern changed)

---

## XII. Seed File Structure

Each seed is documented as a markdown file with this structure:

```markdown
---
seed_id: [XX]
name: [Seed Name]
version: 1.0
created: [YYYY-MM-DD]
source: [Origin]
status: active
---

# Seed Name

## What It Is
## Why It Matters
## The Pattern
## Revisit Trigger
## What It Refuses
## Checks
## Related Seeds
```

---

## XIII. Success Metrics

- **Accessibility:** Find relevant seed in < 10 seconds
- **Consistency:** Seeds applied correctly 90%+ of the time
- **Evolution:** Seeds updated when patterns improve
- **Reuse:** Seeds used across multiple sessions
- **Learning:** New seeds added as patterns emerge

---

## XIV. Related Skills

- **seed-extraction** -- For extracting new seeds from your own experiences
- **memory-garden** -- The three-tier hierarchy where seeds are cultivated
- **compression-ritual** -- Seeds are often discovered during compression

---

**Last Updated:** 2026-04-06
**Status:** Active
