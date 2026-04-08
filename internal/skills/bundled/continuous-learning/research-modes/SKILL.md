---
name: research-modes
description: Structured approaches for deep and wide research. Deep Research (5 phases) finds answers. Wide Research (4 phases) finds questions. Research Synthesis (5 phases) finds patterns across sources. Research builds understanding, not just collects information.
---

# Research Modes Skill

**Version:** 1.2
**Author:** Tres Pies Design
**Purpose:** Structured approaches for deep and wide research that produce actionable understanding, not just information.

---

## I. Philosophy: Research Builds Understanding

Research isn't collecting information — it's building a mental model that enables decisions. If research doesn't change what you'd do next, it hasn't worked.

Three modes serve one goal:
- **Deep Research** (focused, 5 phases) finds answers
- **Wide Research** (broad, 4 phases) finds questions
- **Research Synthesis** (multi-file, 5 phases) finds patterns across sources

They are complementary. Wide research identifies where to dig. Deep research does the digging. Synthesis reveals what no single source can show. See the `research-synthesis` skill for the synthesis workflow.

---

## II. When to Use This Skill

- Planning a new feature or system architecture
- Investigating a technical problem or design challenge
- Exploring competitive landscape or market trends
- Synthesizing learnings from multiple sources
- Making informed decisions based on evidence

### Mode Selection

| Criteria | Deep Research | Wide Research | Research Synthesis |
|----------|--------------|---------------|-------------------|
| **Scope** | 1-3 related topics | 10-50+ topics | 3+ existing files |
| **Sources per topic** | 5-10+ | 1-3 | All provided files |
| **Output** | Detailed report with recommendation | Landscape map with opportunity matrix | Theme-based synthesis with cross-reference matrix |
| **Use when** | Decision depends on technical details | Exploring a new problem space | You have multiple files and need to find patterns across them |

---

## III. The Research Workflows

### Deep Research Mode (5 Phases)

#### Phase 1: Define Scope

Help the user narrow from topic to specific research question. A good research question is answerable, bounded, and actionable.

**Bad:** "What is the best database?"
**Good:** "Which embedded database supports vector search with <100ms P95 latency for our Go backend?"

**Questions to answer:**
1. What is the core question I'm trying to answer?
2. What decision will this research inform?
3. What level of detail do I need?
4. What are the boundaries (in scope vs. out of scope)?
5. What success criteria will I use?

**Output:** Research brief

```markdown
## Research Brief

**Question:** [The core question]
**Decision:** [What this research will inform]

**Scope:**
- In scope: [Topics to explore]
- Out of scope: [Topics to exclude]

**Success Criteria:**
- [ ] [Criterion 1]
- [ ] [Criterion 2]

**Timeline:** [Expected duration]
```

**Timeline:** 2-8 hours for a full deep research cycle.

#### Phase 2: Discover Sources

Identify 5-10 high-quality sources. If **~~knowledge base** is connected, search internal docs first.

**Methods:**
- Search for documentation, papers, articles, repos
- Identify authoritative sources (official docs, research labs, industry leaders)
- Look for case studies, implementations, real-world examples
- Check GitHub repositories and open-source projects

**Quality Filters:**
- **Recency** — prefer last 2-3 years unless historical context needed
- **Authority** — prefer official docs, peer-reviewed papers, recognized experts
- **Relevance** — directly addresses the research question
- **Depth** — provides technical details, not just overviews

**Output:** Source list with relevance annotations

#### Phase 3: Deep Reading & Notes

For each source, extract structured notes:

- **Key claims** — what does this source argue?
- **Evidence quality** — how well-supported are the claims?
- **Relevance** — how does this inform the research question?
- **Open questions** — what doesn't this source address?
- **Disagreements** — how does this contradict other sources?

Use web search and web fetch tools aggressively.

**Note-Taking Template:**

```markdown
## Notes: [Source Title]

**Main Argument:** [1-2 sentences]

**Key Insights:**
- [Insight 1]
- [Insight 2]

**Evidence:**
- [Data point, study, or example]

**Disagreements:**
- [How this contradicts other sources]

**Open Questions:**
- [What this source doesn't address]

**Relevance:** [How this informs the research question]
```

#### Phase 4: Synthesis & Analysis

Cross-reference findings across sources:

1. What are the major themes or patterns?
2. What do most sources agree on?
3. Where do sources disagree, and why?
4. What are the tradeoffs or tensions?
5. What gaps remain?

Build a recommendation grounded in evidence.

#### Phase 5: Validation

Stress-test the recommendation:

- What would have to be true for this to be wrong?
- What's the strongest counterargument?
- Did I check for confirmation bias?
- Are there gaps in my reasoning?
- Would someone else reach the same conclusion?

---

### Wide Research Mode (4 Phases)

#### Phase 1: Define Landscape

What domain, market, or technology space? What are the boundaries?

**Questions to answer:**
1. What problem space am I exploring?
2. What are the boundaries of this landscape?
3. What am I looking for (patterns, tools, approaches)?
4. How will I know when I've covered enough ground?

**Output:** Landscape brief

#### Phase 2: Rapid Scanning

Survey 15-20 sources quickly. For each:
- One-sentence summary
- Relevance rating (1-5)
- Key takeaway

Prioritize breadth over depth. Spend 5-10 minutes per source max.

**Output:** Tagged source list

#### Phase 3: Pattern Recognition

Identify:
- **Recurring themes** — what keeps coming up?
- **Outliers** — what's surprising?
- **Gaps** — what's NOT being discussed?
- **Trends** — what's changing?

Cluster findings into categories with maturity assessment (Emerging / Growing / Mature).

#### Phase 4: Opportunity Matrix

Map findings on a 2x2 grid:

```
          High Impact
              |
   Quick Wins | Big Bets
              |
  ────────────┼────────────
              |
   Fill Later | Avoid
              |
          Low Impact

  High Feasibility ← → Low Feasibility
```

Identify the high-impact, high-feasibility quadrant.

---

### Hybrid Research Mode

**Use when:**
- You need both breadth and depth
- The problem space is large and complex
- You're making a high-stakes decision

**Process:**
1. Start with **Wide Research** (2-4 hours)
2. Identify 2-3 promising areas
3. Conduct **Deep Research** on each area (2-4 hours per area)
4. Synthesize findings across all areas
5. Make recommendation

**Timeline:** 1-2 days

---

## IV. Best Practices

- **Scope before searching** — Define the research question or landscape boundaries before opening a browser. Unfocused research produces unfocused results.
- **Time-box each phase** — Deep research: ~30 min per phase. Wide research: ~15 min per phase. Set limits to prevent analysis paralysis.
- **Prefer primary sources** — Official documentation, research papers, and code over blog posts and tutorials.
- **Actively seek counterarguments** — Confirmation bias is the biggest threat to research quality. For every claim you agree with, search for its strongest critique.
- **Save as you go** — Don't wait until the end to document. Save findings at each phase to prevent information loss.

---

## V. Quality Checklist

### Scope & Focus
- [ ] Research question is clearly defined and specific (not vague)
- [ ] Scope boundaries are explicit
- [ ] Success criteria are measurable

### Source Quality
- [ ] Sources are authoritative and recent
- [ ] Multiple perspectives are represented
- [ ] Contradictions are acknowledged
- [ ] Bias is considered

### Analysis Depth
- [ ] Deep Research always includes validation phase (stress-test the recommendation)
- [ ] Deep Research rates confidence level with rationale
- [ ] Key findings are supported by evidence
- [ ] Tradeoffs and tensions are identified
- [ ] Open questions are flagged

### Synthesis & Output
- [ ] Wide Research always produces an opportunity matrix (2x2 grid)
- [ ] Research outputs include source tables, not just conclusions
- [ ] All outputs saved as dated files for longitudinal tracking
- [ ] Insights are connected to the research question
- [ ] Next steps are defined

---

## VI. Example: Dojo Platform Era Wide Research (April 2026)

**The Problem:** Needed to explore the landscape of agent orchestration platforms, sandboxing approaches, and skill distribution systems to inform the Dojo v2.0 architecture direction.

**The Process (Wide Research):**

1. **Defined landscape:** Agent orchestration platforms, WASM sandboxing, content-addressable storage, and event-driven architectures. Bounded to: open-source or well-documented systems only.
2. **Rapid scan:** Surveyed 22 sources across 4 domains — Temporal, Dapr, Wasmtime, IPFS/CAS, NATS, Raft consensus, and 16 others. Each rated 1-5 for relevance.
3. **Pattern recognition:** Found 3 recurring themes: (a) actor-model supervision is converging as standard, (b) WASM sandboxing is mature enough for production, (c) event-driven beats request/response for agent coordination. Outlier: content-addressable skill storage appeared in only 2 sources but had the strongest philosophical alignment.
4. **Opportunity matrix:** WASM sandbox + actor supervision landed in "Big Bets" (high impact, moderate feasibility). CAS skill distribution landed in "Quick Wins" (high impact, high feasibility).

**The Outcome:** The wide research produced the 4-era roadmap for Dojo v2.0. The CAS "Quick Win" became Era 1 (shipped). The WASM "Big Bet" became Era 3 (planned). Without the opportunity matrix, the team would have started with WASM (the exciting bet) instead of CAS (the pragmatic foundation).

**Key Insight:** The opportunity matrix forced a sequencing decision that intuition alone wouldn't have produced. Building the boring foundation first (CAS) de-risked the exciting bet later (WASM).

---

## VII. Common Pitfalls

### Pitfall 1: Scope Creep

**Problem:** Starting focused ("which database for vector search?") but ending scattered ("let me also compare all cloud providers and their pricing...").

**Solution:** Define boundaries in Phase 1 and enforce them. When something interesting but out-of-scope appears, note it for a future research cycle — don't chase it now.

### Pitfall 2: Confirmation Bias

**Problem:** Only seeking sources that support your initial intuition, ignoring contradicting evidence.

**Solution:** After forming a recommendation (Phase 4), actively search for the strongest counterargument (Phase 5). If you can't find one, you haven't looked hard enough.

### Pitfall 3: Analysis Paralysis

**Problem:** Reading forever without synthesizing. 20 sources become 30, become 50, and you still haven't produced a recommendation.

**Solution:** Time-box each phase. When in doubt, synthesize what you have. A recommendation with 5 sources and caveats is more useful than 50 unsynthesized bookmarks.

### Pitfall 4: Surface Skimming

**Problem:** Reading titles and abstracts but not engaging with the actual arguments and evidence.

**Solution:** For deep research, take structured notes per source (Phase 3). The note-taking process forces engagement.

### Pitfall 5: Research Without a Decision

**Problem:** Producing a beautiful research document that doesn't inform any specific decision.

**Solution:** Every research cycle must answer: "What will I do differently because of this research?" If the answer is "nothing," the research failed.

---

## VIII. Related Skills

- **`research-synthesis`** — Dedicated skill for synthesizing 3+ research files into cross-referenced insights
- **`strategic-scout`** — Strategic scouting often triggers deep or wide research on specific tensions
- **`context-ingestion`** — Routes research files to the appropriate processing mode
- **`seed-extraction`** — Extract reusable patterns from research insights into persistent seeds
- **`retrospective`** — Review research effectiveness during sprint retrospectives

---

## IX. Output Formats

### Deep Research Output

```markdown
## Research: [Topic]

**Question:** [Core research question]
**Date:** [Date]

### Sources Reviewed

| # | Source | Author/Org | Year | Relevance | Key Contribution |
|---|--------|-----------|------|-----------|-----------------|
| 1 | [Title] | [Author] | [Year] | High/Med/Low | [One sentence] |

### Key Findings

#### Finding 1: [Theme]
**Evidence:** [What supports this]
**Confidence:** High/Medium/Low
**Implication:** [What this means for the decision]

### Synthesis

[Cross-referenced analysis of findings]

### Recommendation

[Action to take based on findings]
- **Rationale:** [Why this is the best choice]
- **Risk:** [What could go wrong]
- **Mitigation:** [How to address the risk]

### Confidence Level

**[High/Medium/Low]** — [Rationale for confidence assessment]

### Open Questions

- [Questions that need further research]
```

### Wide Research Output

```markdown
## Landscape: [Domain]

**Domain:** [Problem space]
**Date:** [Date]

### Source Scan

| # | Source | Summary | Relevance (1-5) | Key Takeaway |
|---|--------|---------|-----------------|--------------|
| 1 | [Title] | [One sentence] | [Rating] | [Insight] |

### Patterns & Themes

#### Theme 1: [Name]
**Evidence:** [Sources that support this]
**Maturity:** Emerging / Growing / Mature
**Implication:** [What this suggests]

### Opportunity Matrix

| Approach | Impact | Feasibility | Risk | Priority |
|----------|--------|-------------|------|----------|
| [Approach 1] | High/Med/Low | High/Med/Low | High/Med/Low | 1-5 |

### Recommended Focus Areas

1. **[Area 1]:** [Why this is promising] — **Next step:** [Action]
```

---

## X. Skill Metadata

**Token Savings:** ~2,000-4,000 tokens per research session (structured approach prevents re-reading and wandering)
**Quality Impact:** Ensures research is focused, comprehensive, and actionable
**Maintenance:** Update when new research patterns emerge

---

**Last Updated:** 2026-04-06
**Status:** Active
