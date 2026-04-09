---
name: skill-audit
description: >
  Systematically audit all skills in a directory to assess completeness, grade
  each skill, identify ecosystem-wide patterns, and prioritize upgrades. Use for
  quarterly maintenance, after creating 5+ new skills, when skills are reported
  as confusing, before major releases, after major refactors, or when onboarding
  new contributors to assess current ecosystem health.
triggers:
  - "audit all skills"
  - "grade the skill ecosystem"
  - "skill ecosystem health"
  - "quarterly skill audit"
  - "which skills need upgrading"
  - "assess skill completeness"
---

# Skill Audit

## I. Philosophy: Ecosystem Health Through Systematic Assessment

A skills ecosystem, like any knowledge base, degrades over time. New skills are created with high standards, but older skills may lack structure, examples, or ecosystem connections. Without regular audits, quality becomes inconsistent, making the ecosystem harder to use and harder to maintain.

This skill is about *proactive* maintenance at the ecosystem level. While skill-maintenance operates on individual skills, skill-audit operates on the entire collection. By systematically grading every skill against clear criteria, we can identify patterns -- not just individual gaps, but systemic weaknesses that affect the whole ecosystem.

The core principle: **measure before you improve.** Grading and pattern analysis come before any upgrade work. Without data, you're guessing at what to fix.

## II. When to Use This Skill

- **Quarterly maintenance schedule** -- every 3 months as a routine health check
- **After creating 5+ new skills** -- to check that new additions are consistent with existing standards
- **When a skill is reported as confusing or incomplete** -- may indicate a broader pattern
- **Before a major release** -- to ensure all skills are production-ready
- **After a major refactor** -- to verify consistency after large-scale changes
- **When onboarding a new contributor** -- to assess and communicate current ecosystem state
- **When establishing standards for a new skills directory** -- to create the grading baseline

**When NOT to use:** For maintaining a single skill (use skill-maintenance). For creating new skills (use skill-creation). This skill is for *ecosystem-wide* assessment and prioritization.

## III. The Audit Workflow

### Step 1: Inventory

**Goal:** Build a complete picture of every skill in the directory without reading each one in full.

**Actions:**
1. List every skill directory
2. For each skill, record:
   - Name (directory name)
   - Line count (quick size indicator)
   - Has YAML frontmatter (Y/N)
   - Section headers present (extract with pattern matching)
3. Build an inventory table

**Output:** Skills inventory with basic metrics for every skill.

**Key Insight:** Use file listing and header extraction for the initial pass. Do not read every skill in full -- that comes in Step 2 only for borderline cases. Skills under 50 lines are likely incomplete.

### Step 2: Grade Each Skill

**Goal:** Assign an objective grade to every skill based on structural criteria.

**Grading Rubric:**

| Grade | Criteria |
|---|---|
| **A+** | All required sections + all recommended sections (philosophy, workflow, checklist, pitfalls, example, related skills) |
| **A** | All required sections (philosophy, workflow, checklist, pitfalls), some recommended |
| **B** | Has workflow and checklist, missing philosophy or pitfalls |
| **C** | Has workflow but missing checklist, pitfalls, or philosophy |
| **D** | Incomplete workflow |
| **F** | Stub or non-functional |

**Required sections:** Philosophy (I.), When to Use (II.), Workflow (III.), Best Practices (IV.), Quality Checklist (V.)

**Recommended sections:** Example (VI./VII.), Common Pitfalls (VI./VII.), Related Skills (VIII.)

**Actions:**
1. Check each skill against the rubric using the inventory data from Step 1
2. For borderline cases, read the skill in full before assigning a grade
3. Record the grade and any notable gaps for each skill

**Output:** Grade for every skill, with gap notes.

**Key Insight:** Grade objectively. "The skill works fine for me" is not an A+ -- structure matters for discoverability and maintainability.

### Step 3: Identify Patterns

**Goal:** Find systemic issues that affect multiple skills.

**Actions:**
1. Calculate grade distribution: how many A+, A, B, C, D, F?
2. Look for systematic gaps:
   - Do most skills lack philosophy sections? Systemic issue.
   - Are pitfalls consistently missing? Pattern.
   - Are newer skills better than older ones? Improving practices.
   - Are skills in one domain stronger than another? Uneven investment.
3. Compare against health thresholds:
   - **Healthy:** 80%+ at A or A+, no D/F skills
   - **Needs Attention:** 60-79% at A or A+, or any D/F skills
   - **Critical:** <60% at A or A+, or multiple D/F skills
4. Document patterns and their implications

**Output:** Pattern analysis with ecosystem health assessment.

**Key Insight:** Patterns are more valuable than individual grades. A pattern like "no skills have philosophy sections" tells you to update the creation process, not just fix individual skills.

### Step 4: Prioritize Upgrades

**Goal:** Create an actionable, ordered list of upgrades sorted by impact and effort.

**Actions:**
1. For each skill needing upgrade, assess:
   - **Impact:** How frequently is this skill used? Does it block other workflows?
   - **Effort:** How much work to upgrade? (Structural fixes: ~30-45 min. Adding content: ~15-30 min.)
2. Assign priority:
   - **High:** D/F skills, or C skills that are frequently used
   - **Medium:** B skills, or C skills that are occasionally used
   - **Low:** A skills that could become A+
3. Sort by priority, then by impact within each priority level

**Output:** Prioritized upgrade list with impact and effort estimates.

**Key Insight:** Prioritize by impact times effort, not alphabetically. A frequently-used C skill matters more than a rarely-used B skill.

### Step 5: Produce the Audit Report

**Goal:** Create a comprehensive, actionable report.

**Actions:**
1. Compile the skills inventory table (name, grade, gap summary)
2. Include grade distribution with percentages
3. Document identified patterns
4. List top 5 upgrade recommendations with: current grade, target grade, gap, impact, and effort
5. Assign ecosystem health verdict
6. Save as `docs/audits/[YYYY-MM-DD]_skills_audit.md`

**Output:** Complete audit report.

**Report Structure:**
- Skills Inventory (table)
- Grade Distribution (counts and percentages)
- Patterns Identified (bullet list)
- Top 5 Upgrade Recommendations (detailed)
- Ecosystem Health Verdict (Healthy / Needs Attention / Critical)

## IV. Best Practices

### Define Criteria Before Grading

Create or reference the grading rubric before reading any skills. This prevents bias -- you grade against the criteria, not against your expectations.

### Use Automation for the Initial Pass

Manually reading 10+ skills is time-consuming and error-prone. Use file listing, line counts, and header extraction to build the inventory quickly. Save full reads for borderline cases.

### Focus on Structural Issues First

Structure is easier to fix than content, and has higher impact on usability. A skill with proper section numbering, headers, and flow is more usable even if the content is imperfect.

### Track Trends Across Audits

If you conduct audits regularly, compare results over time. Is the ecosystem getting healthier? Are the same gaps recurring? Trend data is more valuable than point-in-time snapshots.

### Celebrate the Work

Maintenance is often invisible. Make audit results visible: share the report, highlight improvements, acknowledge the effort.

## V. Quality Checklist

- [ ] Every skill in the directory was inventoried
- [ ] Every skill received a grade based on the rubric
- [ ] Grade distribution was calculated with percentages
- [ ] Patterns were identified across the ecosystem
- [ ] Ecosystem health verdict was assigned using defined thresholds
- [ ] Upgrades were prioritized by impact times effort, not alphabetically
- [ ] Top 5 recommendations include current grade, target grade, gap, impact, and effort
- [ ] Audit report was saved to `docs/audits/`
- [ ] No skills were skipped or forgotten in the inventory

## VI. Common Pitfalls

### No Clear Rubric

**Problem:** Without a defined rubric, grading becomes subjective and inconsistent. "This feels like a B" is not reproducible.

**Solution:** Always reference or define the grading rubric before starting. The rubric in Section III, Step 2 is the default.

### Reading Every Skill in Full

**Problem:** Attempting to read 10+ skills in full is time-consuming and unnecessary for the initial assessment.

**Solution:** Use the two-pass approach: inventory first (headers and line counts), then full reads only for borderline cases.

### Alphabetical Prioritization

**Problem:** Upgrading skills in alphabetical order instead of by impact. The most important skill might start with "Z."

**Solution:** Always sort by impact times effort. A skill that blocks three other workflows is more urgent than one that's rarely used.

### Upgrading Everything at Once

**Problem:** Trying to bring all skills to A+ in a single session is overwhelming and leads to burnout or incomplete work.

**Solution:** Focus on high-priority upgrades first. Set a realistic scope (e.g., "upgrade 2-3 skills per session"). A skills ecosystem is never "done."

### Skipping the Report

**Problem:** Conducting the audit mentally but not producing a written report. The audit work is invisible and not repeatable.

**Solution:** Always produce the report. It takes 15-20 minutes and serves as the historical record, the communication artifact, and the baseline for the next audit.

## VII. Example: Skills Ecosystem Audit

**Context:** After a productive week of creating 5 new skills, a team audited their 11-skill ecosystem to check for consistency.

**Inventory Results:**
- 11 skills total
- Line counts ranged from 93 to 412 lines
- All had YAML frontmatter

**Grade Distribution (Before):**
- A+: 4/11 (36%)
- A: 5/11 (45%)
- B+: 2/11 (18%)

**Patterns Identified:**
- 2 skills lacked section numbering (structural issue)
- 2 skills missing Example, Pitfalls, and Related Skills sections
- All recently created skills were A+ (improving practices)

**Upgrades Performed:**
- 2 B+ skills upgraded to A+ (added section numbering, examples, pitfalls, related skills)
- ~200 lines added per skill
- ~1.5 hours total investment

**Grade Distribution (After):**
- A+: 6/11 (55%)
- A: 5/11 (45%)
- B+ or lower: 0/11 (0%)

**Ecosystem Health:** Healthy (100% at A or A+)

## VIII. Related Skills

- **skill-maintenance** -- For maintaining individual skills (this skill audits the whole ecosystem)
- **skill-creation** -- For creating new skills that meet audit standards from the start
- **process-extraction** -- For capturing the audit process itself as institutional knowledge

## IX. Audit Schedule

| Frequency | Trigger |
|---|---|
| **Quarterly** | Routine ecosystem health check |
| **After 5+ new skills** | Consistency check for new additions |
| **Before major releases** | Ensure production-readiness |
| **After major refactors** | Verify consistency post-changes |
| **On report of confusion** | Investigate whether it's a systemic issue |
