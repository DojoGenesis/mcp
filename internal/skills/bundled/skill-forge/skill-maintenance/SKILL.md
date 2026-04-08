---
name: skill-maintenance
description: >
  Systematic process for maintaining skill health through accuracy checks,
  completeness reviews, and terminology updates. Use when skill names become
  unclear, terminology needs updating, periodic audits are due, new skills need
  cross-reference validation, or skills are reported as confusing or outdated.
---

# Skill Maintenance

## I. Philosophy: Maintenance Is Keeping What Works Working Well

A skills directory is a living knowledge base. As the project evolves, skill names may become unclear, terminology may need updating, and cross-references may drift out of sync. Without systematic maintenance, the knowledge base degrades: skills become hard to find, references break, and terminology becomes inconsistent.

This skill transforms maintenance from a reactive chore into a proactive ritual. By following a structured process -- read, check, propose, execute, verify, document -- we ensure the skills directory remains clear, consistent, and valuable.

Good maintenance prevents future problems. A well-maintained knowledge base is easier to search, easier to understand, and easier to extend.

## II. When to Use This Skill

- **Skill names become unclear or outdated** -- names no longer describe what the skill does
- **Terminology needs updating** -- tool-specific language should be generalized, or industry terms have shifted
- **Adding new skills** -- new skills reference existing skills and need cross-reference validation
- **Deprecating or merging skills** -- old skills need to be removed or consolidated
- **Conducting periodic audits** -- regular health checks (quarterly recommended)
- **After a major refactor** -- ensuring consistency after large-scale changes
- **Onboarding new contributors** -- ensuring the knowledge base is accessible

**When NOT to use:** For creating new skills (use skill-creation). For ecosystem-wide assessment (use skill-audit). This skill is for maintaining *individual* skills or small groups.

## III. The Maintenance Workflow

### Step 1: Recognize the Need

**Goal:** Identify when maintenance is needed and clarify scope.

**Actions:**
1. Identify the trigger: user request, noticed inconsistency, scheduled audit, or new skill addition
2. Clarify scope with the user: What specifically needs to change? What should stay the same?
3. Document the maintenance goal

**Output:** Clear understanding of maintenance scope.

**Key Insight:** Always pause and clarify scope before large refactors. Over-refactoring is worse than under-refactoring.

### Step 2: Read Before Proposing

**Goal:** Understand what the skills actually do before suggesting changes.

**Actions:**
1. Read all skills that will be affected by the maintenance, in full
2. Understand the actual purpose and workflow of each skill
3. Note cross-references between skills
4. Identify patterns in naming or terminology

**Output:** Deep understanding of affected skills.

**Key Insight:** Never propose renames or refactors without reading the actual content first. Names should reflect reality, not assumptions.

### Step 3: Check Accuracy

**Goal:** Verify everything in the skill is still true.

**Actions:**
1. Check: Have conventions changed since the skill was written?
2. Check: Are referenced tools, APIs, or patterns still current?
3. Check: Do examples still represent real, valid scenarios?
4. Check: Has the domain itself evolved?

**Output:** List of accuracy issues (if any).

### Step 4: Check Completeness

**Goal:** Identify anything missing.

**Actions:**
1. Check: Are there new trigger scenarios the skill should cover?
2. Check: Are there new pitfalls discovered through real usage?
3. Check: Are there new best practices that have emerged?
4. Check: Should the skill reference new related skills?

**Output:** List of completeness gaps (if any).

### Step 5: Check Triggers and Naming

**Goal:** Ensure the skill is discoverable and clearly named.

**Actions:**
1. Check: Would a new agent recognize when to use this skill from the YAML description alone?
2. Check: Does the description include enough trigger words?
3. Check: Are 5+ concrete scenarios listed in "When to Use"?
4. Check: Does the skill name follow verb-object convention?
5. If renaming is needed: propose names following verb-object pattern, get user confirmation

**Output:** Trigger and naming assessment.

### Step 6: Check Workflow and Checklist

**Goal:** Verify the workflow steps and quality criteria are current.

**Actions:**
1. Check: Are steps in the right order? Should any be added, removed, or reordered?
2. Check: Are decision points clearly marked?
3. Check: Do quality checklist items still apply?
4. Check: Are there new criteria that should be added?

**Output:** Workflow and checklist assessment.

### Step 7: Check Supporting Files

**Goal:** Verify templates, references, and examples are current.

**Actions:**
1. Check: Do templates match the current workflow?
2. Check: Are references still valid and accessible?
3. Check: Do examples still demonstrate the skill accurately?

**Output:** Supporting files assessment.

### Step 8: Execute Changes

**Goal:** Apply all identified changes systematically.

**Actions:**
1. If renaming: rename directories using `mv` AND update internal metadata (name field, title heading)
2. If refactoring terminology: search and catalog all references with `grep -r -i`, determine strategy (what changes, what stays), execute with batch edits, verify no broken references
3. If updating content: apply changes to SKILL.md and supporting files
4. If changes are non-trivial: increment version (e.g., v1.0 to v1.1)
5. Update "Last Updated" date
6. Verify renames with `ls` command

**Example Batch Edit:**
```json
{
  "edits": [
    {"all": true, "find": "old-term", "replace": "new-term"},
    {"all": true, "find": "Old Term", "replace": "New Term"}
  ]
}
```

**Output:** Updated skill files.

**Key Insight:** When refactoring terminology, always catalog before editing. Use systematic search to find all references -- do not rely on memory. Batch edits with `all: true` are more efficient than one-by-one replacements.

### Step 8.5: Verify and Commit

**Goal:** Ensure all changes are correct and commit with comprehensive documentation.

**Actions:**
1. Verify no unintended references remain: `grep -i "<old term>" <directory>`
2. Check git status: `git status`
3. Stage changes: `git add skills/`
4. Write comprehensive commit message with file-by-file breakdown
5. Commit and push

**Commit Message Template:**
```
<Action> in skills directory

<Summary paragraph>

Changes:
- <file1> (<N> replacements)
  - <change 1>
  - <change 2>

- <file2> (<N> replacements)
  - <change 1>

Kept <term> only when:
- <context 1>
- <context 2>

<Rationale paragraph>
```

**Key Insight:** Comprehensive commit messages are documentation. Future maintainers need to understand *why* changes were made, not just *what* changed.

### Step 9: Document and Report

**Goal:** Create a maintenance report for future reference.

**Actions:**
1. Create maintenance report documenting:
   - What was checked (all 7 dimensions)
   - What changed (with rationale)
   - Current grade (A+ through F)
2. Save as `docs/maintenance/[YYYY-MM-DD]_[skill-name]_maintenance.md`

**Output:** Maintenance report.

**Grading Rubric:**

| Grade | Criteria |
|---|---|
| **A+** | Complete philosophy, workflow, checklist, pitfalls, supporting files, process example |
| **A** | Complete philosophy, workflow, checklist, pitfalls |
| **B** | Has workflow and checklist, missing philosophy or pitfalls |
| **C** | Has workflow but missing checklist, pitfalls, or philosophy |
| **D** | Incomplete workflow |
| **F** | Stub or non-functional |

## IV. Best Practices

### Pause and Clarify Scope

Before any refactor, clarify: What exactly needs to change? What should stay the same? What's the scope? Avoid the temptation to "fix everything while we're at it."

### Read Before Proposing

Never suggest changes without reading the actual content. Understanding reality prevents proposing changes that don't make sense.

### Use Systematic Search

Don't rely on memory to find all references. Use search to systematically catalog all instances before making changes.

### Preserve Contextually Appropriate References

Not all references should be changed. When refactoring terminology, preserve tool-specific mentions when they're contextually appropriate (e.g., when listing multiple tools by name).

### Batch Edits for Efficiency

Use the file tool's multiple edits feature with `all: true` to replace all occurrences in a single operation. Batch edits are more efficient than one-by-one replacements.

### Comprehensive Commit Messages

Write commit messages that explain: **What** changed (file-by-file breakdown), **Why** it changed (rationale), **What** was preserved (and why).

### Document Immediately

Create maintenance reports right after completing the work. Details fade quickly, and future maintainers will need this context.

## V. Quality Checklist

- [ ] All 9 maintenance steps were followed, not a subset
- [ ] Every skill was read in full before changes were proposed
- [ ] Scope was clarified with the user before work began
- [ ] Accuracy, completeness, triggers, workflow, checklist, and supporting files were all checked
- [ ] If renaming: both directory and internal metadata were updated
- [ ] If refactoring terminology: all references were cataloged before changes began
- [ ] No broken cross-references remain after changes
- [ ] Version was bumped if changes were non-trivial
- [ ] Maintenance report was created with grade
- [ ] Changes were committed with a comprehensive message (what, why, what was preserved)

## VI. Common Pitfalls

### Over-Refactoring

**Problem:** Changing references that should stay -- historical documents, tool-specific mentions, examples that use real names for clarity.

**Solution:** Scope appropriately. Only change what needs changing. When in doubt, preserve and note it in the maintenance report.

### Proposing Without Reading

**Problem:** Suggesting renames or changes based on assumptions about what a skill does, rather than what it actually does.

**Solution:** Read every affected skill completely before proposing any changes. Names should reflect reality.

### Missing References

**Problem:** Relying on memory to find all references that need updating, leading to inconsistencies after the refactor.

**Solution:** Use systematic search (grep or equivalent) to catalog every reference before making any changes.

### Vague Commit Messages

**Problem:** Committing maintenance changes with messages like "Updated skills" that don't explain what changed or why.

**Solution:** Write comprehensive commit messages: what changed, why it changed, and what was intentionally preserved.

### Skipping the Report

**Problem:** Completing maintenance without recording the process. Future maintainers don't know what was checked, what was changed, or why.

**Solution:** Always produce a maintenance report. It takes 10 minutes and saves hours of future confusion.

## VII. Example: Skill Rename and Terminology Refactor

**Context:** After creating 5 new skills, 4 existing skills had names that didn't follow the verb-object convention, and 32 references used tool-specific terminology that needed generalization.

**Process:**
1. User requested renames for clarity
2. Read all 4 skills to understand their actual purpose
3. Proposed verb-object names, user refined (e.g., "spec-writer" became "write-release-specification")
4. Renamed directories and updated all internal metadata
5. Searched repository: found 530 total references to the old terminology
6. User scoped down: skills directory only (32 references)
7. Determined strategy: what to change, what to preserve
8. Executed batch edits across 4 files (32 replacements)
9. Verified, committed, and documented

**Outcome:** 4 skills renamed, 32 references updated, 6 tool-specific references intentionally preserved. All cross-references intact.

**Time:** ~2 hours total.

## VIII. Maintenance Schedule

| Frequency | Trigger |
|---|---|
| **Quarterly** | Full audit of skills directory |
| **After 5+ new skills** | Consistency check for new additions |
| **On request** | When naming or terminology issues are identified |
| **Before major releases** | Ensure consistency before public releases |

## IX. Related Skills

- **`compression-ritual`** -- For preserving insights before large refactors
- **`process-extraction`** -- For creating new skills from maintenance processes
- **`health-audit`** -- For comprehensive repository health audits
- **`documentation-audit`** -- For identifying documentation drift

---

**Last Updated:** 2026-04-06
**Status:** Active
