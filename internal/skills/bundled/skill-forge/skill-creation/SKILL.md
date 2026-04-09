---
name: skill-creation
version: "1.1"
description: >
  Guide for creating effective, reusable skills that extend agent capabilities
  with specialized knowledge and workflows. Use when creating a new skill from
  scratch, updating an existing skill, designing skill architecture, or when a
  user says "make a skill for X." Also use when formalizing a workflow, capturing
  domain knowledge, or building institutional memory.
triggers:
  - "create a skill"
  - "build a new skill"
  - "make a skill for"
  - "write a skill"
  - "formalize a workflow into a skill"
  - "add a skill"
---

# Skill Creation

## I. Philosophy: Skills Encode Thinking, Not Tool Usage

A skill is not a script. It is not a checklist. It is a *pattern of thinking* -- a way of approaching a class of problems that any competent agent can execute, regardless of what tools are available.

The agent is already smart. It can write code, parse files, search the web, and reason about complex problems. What it lacks is *your* domain knowledge: the institutional conventions, the hard-won patterns, the "we tried X and it failed because of Y" lessons that took months to learn.

A good skill provides exactly this: the knowledge gap between general intelligence and domain expertise. Nothing more. The context window is a shared resource -- every token in a skill must justify its cost.

### Progressive Disclosure

Skills load in three tiers:

1. **Metadata** (~100 words) -- Always in context. The YAML `description` field determines whether the skill activates.
2. **SKILL.md body** (<500 lines) -- Loaded when the skill triggers. Contains the core workflow.
3. **Bundled resources** -- Loaded on demand. Templates, references, scripts, examples.

This architecture means: **the description field is the most important part of the skill.** If the description doesn't capture when to use the skill, the body never loads.

## II. When to Use This Skill

- **Creating a new skill from scratch** for a domain, tool, or workflow
- **Updating an existing skill** that needs restructuring or enrichment
- **Designing skill architecture** for a new skills ecosystem
- **Formalizing a repeatable workflow** that agents should follow consistently
- **Capturing domain expertise** that would otherwise live only in someone's head
- **Reviewing a skill's quality** against structural standards
- **Onboarding someone** to the skill creation process

**When NOT to use:** If the knowledge is a simple fact or reminder (not a process), it should be a seed, not a skill. If the knowledge is tool-specific and non-transferable, it should be documentation, not a skill.

## III. The Skill Creation Workflow

### Step 1: Understand the Skill

**Goal:** Develop a concrete understanding of what the skill provides and when it triggers.

**Actions:**
1. Ask what capability the skill provides that the agent doesn't already have
2. Identify at least 5 trigger scenarios (situations that should activate this skill)
3. Identify inputs (what the skill starts with) and outputs (what it produces)
4. Clarify the boundary: what is in scope vs. out of scope

**Output:** Clear skill specification with triggers, inputs, outputs, and scope.

**Key Insight:** The most common failure is creating skills that teach the agent what it already knows. Challenge every piece of content: "Does the agent really need this?"

### Step 2: Plan Contents

**Goal:** Determine what the skill needs beyond SKILL.md.

**Actions:**
1. Evaluate whether supporting files are needed:
   - `templates/` -- For structured outputs the skill produces
   - `references/` -- For external knowledge the skill relies on
   - `scripts/` -- For executable, deterministic steps
   - `examples/` -- For real-world process walkthroughs
2. Apply the test: if the content would exceed 500 lines in SKILL.md, split it

**Output:** Content plan listing SKILL.md sections and any supporting files.

**Key Insight:** Most skills are SKILL.md-only. Only add supporting files when there's a clear need. Avoid duplication -- information lives in SKILL.md OR references, not both.

### Step 3: Initialize

**Goal:** Create the skill directory with proper structure and frontmatter.

**Actions:**
1. Create directory: `skills/[skill-name]/`
2. Create SKILL.md with YAML frontmatter:
   - `name`: The skill name (verb-object pattern preferred)
   - `description`: Comprehensive trigger description including when to use AND when not to use
3. Create supporting directories only if planned in Step 2

**Output:** Initialized skill directory ready for content.

**Key Insight:** The YAML `description` is the most important field. It determines whether the skill ever gets loaded. Be comprehensive: include trigger words, use cases, and anti-use-cases.

### Step 4: Write the SKILL.md Body

**Goal:** Create the core knowledge document.

**Before writing, consult proven design patterns:**
- **Multi-step processes:** Study existing workflow skills for sequential steps and conditional logic
- **Output formats or quality standards:** Study existing template skills for output patterns
- **Progressive disclosure:** Study how existing skills split content between SKILL.md and references/

**Required sections:**

| Section | Purpose | Minimum Bar |
|---|---|---|
| I. Philosophy | Why this skill exists | 2-3 paragraphs explaining the principle |
| II. When to Use | Trigger scenarios | 5+ concrete scenarios |
| III. Workflow | Step-by-step process | 3+ steps with decision points |
| IV. Best Practices | Actionable wisdom | Each practice explains "why" |
| V. Quality Checklist | Verification criteria | 5+ binary yes/no items |
| VI. Common Pitfalls | What goes wrong | 2+ problem/solution pairs |

**Recommended sections:**

| Section | Purpose |
|---|---|
| VII. Example | Concrete walkthrough on a real scenario |
| VIII. Related Skills | Links to complementary skills |

**Writing guidelines:**
- Use imperative/infinitive form
- Prefer concise examples over verbose explanations
- Match specificity to fragility (narrow bridge = specific guardrails; open field = high freedom)
- Keep SKILL.md under 500 lines

**Output:** Complete SKILL.md body.

### Step 5: Create Supporting Files

**Goal:** Build any templates, references, scripts, or examples planned in Step 2.

**Actions:**
1. Create each supporting file
2. Test any scripts to verify they work
3. Ensure SKILL.md references supporting files with clear "when to read" instructions
4. Remove any unused placeholder files

**Output:** Supporting files ready for use.

### Step 6: Deliver and Iterate

**Goal:** Present the skill and refine based on real usage.

**Actions:**
1. Present the completed skill to the user
2. Ask: "Does this capture the pattern accurately? What's missing?"
3. Highlight assumptions that should be verified
4. Iterate based on feedback

**Iteration cycle (post-delivery):**
1. Use the skill on real tasks
2. Notice struggles or inefficiencies during execution
3. Identify how SKILL.md or bundled resources should be updated
4. Implement changes and test again

**Output:** Finalized skill approved by the user, with iteration plan for post-usage refinement.

**Key Insight:** The first version is rarely perfect. Real usage reveals gaps that planning cannot. The most valuable feedback comes right after using a skill, with fresh context of how it performed.

## IV. Best Practices

### Write Descriptions Like Headlines

The YAML `description` field is the skill's advertisement. It must communicate in ~50 words what the skill does and when to use it. Include trigger keywords generously.

### Match Freedom to Fragility

High-freedom instructions ("choose an appropriate approach") work when multiple approaches are valid. Low-freedom instructions ("follow these exact steps") work when operations are fragile. Most skills need a mix.

### Keep the Context Window Clean

Every line in SKILL.md costs tokens. Challenge each paragraph: "Does this justify its token cost?" Move variant-specific details to reference files. Keep the core workflow in SKILL.md.

### Name with Verb-Object Convention

Good: `write-release-specification`, `audit-skill-ecosystem`, `transform-spec-to-prompt`
Bad: `spec-writer`, `audit`, `prompt-tool`

Verb-object names are self-documenting. They tell you what the skill *does*, not what it *is*.

## V. Quality Checklist

- [ ] YAML frontmatter has `name` and `description` fields
- [ ] Description includes 3+ trigger keywords or phrases
- [ ] Philosophy section explains "why," not just "what"
- [ ] When to Use lists 5+ concrete trigger scenarios
- [ ] Workflow has 3+ numbered steps with goals and outputs
- [ ] Each workflow step has a clear decision point or output
- [ ] Quality checklist has 5+ binary items
- [ ] At least 2 common pitfalls with problem/solution pairs
- [ ] SKILL.md is under 500 lines
- [ ] No references to specific agents or proprietary tools
- [ ] Skill name follows verb-object or descriptive compound convention

## VI. Common Pitfalls

### Teaching the Agent What It Already Knows

**Problem:** Filling skills with instructions the agent can infer (e.g., "write clean code," "handle errors gracefully"). This wastes context tokens.

**Solution:** Before adding content, ask: "Would the agent do this wrong without this instruction?" If no, remove it.

### Vague Trigger Descriptions

**Problem:** Description says "Use for skill-related tasks" -- too broad to trigger correctly.

**Solution:** Include specific trigger scenarios in the description: "Use when creating a new skill, updating an existing skill, or reviewing skill quality."

### Monolithic Skills

**Problem:** A single SKILL.md trying to cover too much ground, exceeding 500 lines.

**Solution:** Split into core workflow (SKILL.md) and variant details (references/). The reader should find the main path in SKILL.md and branch to references only when needed.

### Missing Decision Points

**Problem:** Workflow steps that are purely sequential with no branching. Real processes have decision points.

**Solution:** At each step, ask: "What could go differently here?" Add conditional guidance for the common branches.

### Skills That Are Really Documentation

**Problem:** A "skill" that describes a tool's API rather than encoding a pattern of thinking.

**Solution:** Skills describe *patterns*, not *tools*. If the content would become obsolete when the tool changes, it's documentation. If the reasoning would survive a tool change, it's a skill.
