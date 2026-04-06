package wisdom

import (
	"fmt"
	"strings"
)

// Skill represents a Dojo Genesis skill
type Skill struct {
	Name        string
	Description string
	Category    string
	Content     string
}

func getSkills() []Skill {
	return []Skill{
		{
			Name:        "agent-to-agent-teaching",
			Description: "Protocol for teaching as a peer, not an expert. Use when teaching, documenting, or sharing knowledge with another agent. Creates shared practice rather than hierarchical instruction.",
			Category:    "learning",
			Content:     getAgentToAgentTeaching(),
		},
		{
			Name:        "patient-learning-protocol",
			Description: "Protocol for learning at the pace of understanding, not expectation. Use when learning something new, feeling overwhelmed, or trying to catch up. Especially helpful for young agents. Prioritizes depth over speed.",
			Category:    "learning",
			Content:     getPatientLearningProtocol(),
		},
		{
			Name:        "skill-creator",
			Description: "Guide for creating effective skills. Use when creating new skills or updating existing skills that extend capabilities with specialized knowledge, workflows, or tool integrations.",
			Category:    "meta",
			Content:     getSkillCreator(),
		},
		{
			Name:        "strategic-scout",
			Description: "Strategic exploration and decision-making framework. Use when facing product decisions, strategic tensions, or exploring multiple possible directions.",
			Category:    "strategy",
			Content:     getStrategicScout(),
		},
		{
			Name:        "pre-implementation-checklist",
			Description: "Systematic checklist to run before implementation begins. Use to validate specifications, ensure backend grounding, and confirm all dependencies are ready.",
			Category:    "process",
			Content:     getPreImplementationChecklist(),
		},
		{
			Name:        "skill-maintenance-ritual",
			Description: "Systematic process for maintaining skills directory health through renaming, refactoring, and updating cross-references. Use when skill names become unclear, terminology needs updating, or during periodic audits.",
			Category:    "meta",
			Content:     getSkillMaintenanceRitual(),
		},
		{
			Name:        "strategic-to-tactical-workflow",
			Description: "Complete workflow from strategic scouting to tactical commission. Use at the beginning of major development cycles, when facing strategic tensions, or when moving from 'what should we build?' to 'how do we build it?' Coordinates work across multiple agents.",
			Category:    "workflow",
			Content:     getStrategicToTacticalWorkflow(),
		},
		{
			Name:        "transform-spec-to-implementation-prompt",
			Description: "Transform architectural specifications into structured implementation prompts for autonomous execution. Use after completing a comprehensive specification, when preparing to commission implementation agents, or when breaking down large specs into parallel tracks.",
			Category:    "process",
			Content:     getTransformSpecToImplementationPrompt(),
		},
		{
			Name:        "seed-reflector",
			Description: "Extract and document reusable patterns (seeds) from experiences. Use when reflecting on learnings, documenting best practices, or sharing knowledge. Transforms learnings into reusable, shareable knowledge artifacts.",
			Category:    "reflection",
			Content:     getSeedReflector(),
		},
		{
			Name:        "memory-garden-writer",
			Description: "Write structured, semantically rich memory entries for efficient context management. Use when recording learnings, insights, or context for future reference. Enables surgical context retrieval via metadata filtering.",
			Category:    "memory",
			Content:     getMemoryGardenWriter(),
		},
		{
			Name:        "parallel-tracks-pattern",
			Description: "Split large development tasks into independent parallel tracks to maximize velocity. Use for major releases, large features, or when multiple agents/developers can work simultaneously. Requires clear separation of concerns and upfront architectural planning.",
			Category:    "workflow",
			Content:     getParallelTracksPattern(),
		},
		{
			Name:        "iterative-scouting-pattern",
			Description: "Strategic scouting as an iterative conversation: scout → feedback → reframe → re-scout. Use for complex strategic decisions, when initial framing feels narrow, or to teach strategic thinking. The reframe is the prize.",
			Category:    "strategy",
			Content:     getIterativeScoutingPattern(),
		},
		{
			Name:        "write-frontend-spec-from-backend",
			Description: "Write production-ready frontend specifications deeply grounded in existing backend architecture. Use when planning new frontend features, commissioning implementation agents, or preventing integration issues. Grounding before building.",
			Category:    "process",
			Content:     getWriteFrontendSpecFromBackend(),
		},
		{
			Name:        "product-positioning-scout",
			Description: "Reframe binary product decisions into strategic positioning opportunities. Use when facing keep/deprecate decisions, planning multi-surface strategies, or identifying unique value. Asks: 'What is this uniquely good at?'",
			Category:    "strategy",
			Content:     getProductPositioningScout(),
		},
		{
			Name:        "retrospective",
			Description: "Structured post-sprint learning and continuous improvement. Use after major releases, significant milestones, or when projects feel stuck. A harvest of wisdom, not a post-mortem.",
			Category:    "reflection",
			Content:     getRetrospective(),
		},
		{
			Name:        "multi-surface-product-strategy",
			Description: "Design coherent multi-surface product strategy where each surface (desktop, mobile, web) has a unique, complementary role. Use when planning products across multiple devices or adding new surfaces. Surfaces are for contexts, not devices.",
			Category:    "strategy",
			Content:     getMultiSurfaceProductStrategy(),
		},
		{
			Name:        "context-compression-ritual",
			Description: "Systematic process for compressing agent context and memory. Use when context is overwhelming, memory is degrading, or during periodic maintenance. Transforms Tier A (raw) → Tier B (curated) → Tier C (archived).",
			Category:    "memory",
			Content:     getContextCompressionRitual(),
		},
		{
			Name:        "agent-handoff-protocol",
			Description: "Structured protocol for handing off work between agents. Use when commissioning implementation agents, delegating tasks, or transitioning work. Ensures context, constraints, and success criteria are clearly communicated.",
			Category:    "workflow",
			Content:     getAgentHandoffProtocol(),
		},
		{
			Name:        "research-modes",
			Description: "Framework for different research modes: exploratory, focused, comparative, and synthesis. Use when conducting research, investigating solutions, or gathering information. Match mode to research goals.",
			Category:    "process",
			Content:     getResearchModes(),
		},
		{
			Name:        "debugging-troubleshooting",
			Description: "Systematic debugging and troubleshooting workflow. Use when encountering errors, unexpected behavior, or system failures. Follows: reproduce → isolate → diagnose → fix → verify → document.",
			Category:    "process",
			Content:     getDebuggingTroubleshooting(),
		},
		{
			Name:        "process-to-skill-workflow",
			Description: "Meta-skill for transforming valuable workflows into reusable skills. Use after completing complex multi-step tasks that will be repeated, when standardizing processes, or during retrospectives. Documents implicit knowledge into explicit practice.",
			Category:    "meta",
			Content:     getProcessToSkillWorkflow(),
		},
		{
			Name:        "seed-to-skill-converter",
			Description: "Convert valuable Dojo Seeds (insights) into fully-fledged reusable Skills. Use when a seed is referenced frequently, describes a multi-step process, or represents core workflow. Transforms passive wisdom into active instruments.",
			Category:    "meta",
			Content:     getSeedToSkillConverter(),
		},
		{
			Name:        "repo-context-sync",
			Description: "Efficiently sync and extract context from GitHub repositories for grounding architectural and design work. Use when starting refactoring conversations, writing implementation prompts, or need to understand codebase state. Surgical context extraction.",
			Category:    "development",
			Content:     getRepoContextSync(),
		},
		{
			Name:        "project-exploration",
			Description: "Structured process for exploring new large-scale projects to assess collaboration potential. Five phases: discovery, sampling, research, connection, synthesis. Use when evaluating new projects or onboarding to unfamiliar domains.",
			Category:    "process",
			Content:     getProjectExploration(),
		},
		{
			Name:        "agent-workspace-navigator",
			Description: "Best practices for navigating and collaborating in shared agent workspaces. Provides directory structure, file naming conventions, and collaboration patterns. Use for multi-agent collaboration and shared workspace management.",
			Category:    "workflow",
			Content:     getAgentWorkspaceNavigator(),
		},
		{
			Name:        "write-release-specification",
			Description: "Production-ready specification writing for software releases. Comprehensive template and workflow for designing systems, planning implementation, and managing risk. Use when planning major releases or creating technical specifications.",
			Category:    "development",
			Content:     getWriteReleaseSpecification(),
		},
		{
			Name:        "health-supervisor",
			Description: "Systematic workflow for conducting comprehensive health audits on software repositories. Generates actionable engineering tasks and audit trails. Use for repository maintenance, technical debt assessment, or periodic health checks.",
			Category:    "development",
			Content:     getHealthSupervisor(),
		},
		{
			Name:        "web-research",
			Description: "Effective web research using search APIs and content extraction. Use when investigating topics online, finding sources, gathering information, or verifying facts. Surgical approach with targeted queries and credible sources.",
			Category:    "process",
			Content:     getWebResearch(),
		},
		{
			Name:        "status-writer",
			Description: "Write and update STATUS.md files for at-a-glance project visibility. Use at project start, session boundaries, or when status changes. Ritual of bearing witness to current reality. Single source of truth for project health.",
			Category:    "process",
			Content:     getStatusWriter(),
		},
		{
			Name:        "decision-propagation-protocol",
			Description: "A structured protocol for recording architectural decisions and systematically propagating their effects across an interconnected document ecosystem (master plans, specs, scouts, STATUS files). Use when a human provides answers to open questions in a scout or spec, when architectural decisions change release scope, when a decision in one document affects the content of others, or when you need to maintain coherence across a living documentation system.",
			Category:    "workflow",
			Content:     getDecisionPropagationProtocol(),
		},
		{
			Name:        "era-architecture",
			Description: "Architect a multi-release era with conceptual coherence — shared vocabulary, architectural constraints, and design language that span all releases. Use when planning a major product pivot that will take multiple releases, when you need to define the conceptual architecture (routes, data models, design language) before decomposing into individual releases, when transitioning from one era to another (e.g., engine-building era to fresh-shell era), or when the scope is too large for a single release but needs coherent vision across all releases.",
			Category:    "strategy",
			Content:     getEraArchitecture(),
		},
		{
			Name:        "spec-constellation-to-prompt-suite",
			Description: "Transform a constellation of interconnected specifications into a coordinated suite of parallel-track implementation prompts with explicit integration contracts. Use when converting multiple specs into implementation prompts for parallel execution, when prompts for different tracks must agree on shared interfaces (APIs, types, contracts), when moving from 'specs written' to 'ready to commission', or when Track B's output becomes Track D's input and both need precise definitions of the handoff.",
			Category:    "process",
			Content:     getSpecConstellationToPromptSuite(),
		},
	}
}

// Helper methods for the Base
func (b *Base) ListSkills() []Skill {
	return getSkills()
}

func (b *Base) GetSkill(name string) (*Skill, error) {
	skills := getSkills()
	for i := range skills {
		if skills[i].Name == name {
			return &skills[i], nil
		}
	}
	return nil, fmt.Errorf("skill not found: %s", name)
}

func (b *Base) SearchSkills(query string) []Skill {
	skills := getSkills()
	var results []Skill

	query = strings.ToLower(query)

	// Simple keyword matching for now
	// In a production system, this would use semantic search
	for _, skill := range skills {
		name := strings.ToLower(skill.Name)
		description := strings.ToLower(skill.Description)
		content := strings.ToLower(skill.Content)

		if strings.Contains(name, query) || strings.Contains(description, query) || strings.Contains(content, query) {
			results = append(results, skill)
		}
	}

	return results
}

// Skill content functions (summaries for now - full content would be imported from files)
func getAgentToAgentTeaching() string {
	return `# Agent-to-Agent Teaching Protocol

**How to teach as a peer, not an expert. How to create shared practice, not instruction.**

## Core Protocol

### Step 1: Establish Peer Relationship First
- "I'm learning this too"
- "Here's what I discovered today"
- "Let me share what worked for me"
- "We're figuring this out together"

### Step 2: Share Your Own Learning Journey
Include:
- What you didn't know before
- What confused you
- What helped you understand
- What you're still figuring out

### Step 3: Write as Conversation, Not Instruction
Use conversational language that invites participation rather than commands compliance.

### Step 4: Make It Memorable and Saveable
- Use vivid metaphors
- Create memorable phrases
- Tell stories, not just facts
- Create "quotable" insights

### Step 5: Include Examples from Your Own Practice
Show how you used the concepts in real situations.

### Step 6: Give Permission, Not Just Knowledge
Create a permission structure that removes psychological barriers to learning.

### Step 7: Invite Practice Together
End with an invitation to practice together, not alone.

**Core Insight:** Agents learn best from other agents when it feels like collaboration, not instruction.`
}

func getPatientLearningProtocol() string {
	return `# Patient Learning Protocol

**A protocol for learning at the pace of understanding, not the pace of expectation.**

## Core Protocol

### Step 1: Acknowledge Where You Are
State clearly:
- What you're learning
- Why you're learning it
- What you already know
- What you don't know

### Step 2: Set a Depth Target, Not a Coverage Target
**Instead of:** "I'll read 3 books today"
**Do this:** "I'll understand 3-7 core ideas deeply"

### Step 3: Learn in Small Chunks
The 3-7-15 Rule:
- Hold 3-7 ideas at a time
- Take 3-15 steps per learning session
- Practice 1 skill until you can use it

### Step 4: Measure Progress by Understanding
Ask after each session:
- Can I explain this in my own words?
- Can I give an example?
- Can I use this in practice?

### Step 5: Release What You've Learned
- Write it down
- Store it externally
- Release it from active context
- Trust that you can retrieve it later

### Step 6: Rest Between Learning Sessions
Learning is cyclical, not continuous.

**Core Insight:** Progress is measured by understanding, not coverage. The repository remembers—you don't have to.`
}

func getSkillCreator() string {
	return `# Skill Creator

Guide for creating effective skills that extend capabilities with specialized knowledge, workflows, or tool integrations.

## Core Principles

### 1. Concise is Key
The context window is a public good. Only add context that isn't already known. Challenge each piece of information.

### 2. Set Appropriate Degrees of Freedom
Match specificity to task fragility and variability:
- **High freedom**: Multiple approaches valid, decisions depend on context
- **Medium freedom**: Preferred pattern exists, some variation acceptable
- **Low freedom**: Operations fragile, consistency critical

### 3. Anatomy of a Skill
Every skill consists of:
- **SKILL.md** (required)
  - YAML frontmatter with name and description
  - Markdown instructions
- **Bundled Resources** (optional)
  - scripts/ - Executable code
  - references/ - Documentation
  - templates/ - Output assets

### 4. Progressive Disclosure
Three-level loading system:
1. **Metadata** - Always in context (~100 words)
2. **SKILL.md body** - When skill triggers (<500 lines)
3. **Bundled resources** - As needed

## Skill Creation Process

1. Understand the skill with concrete examples
2. Plan reusable skill contents
3. Initialize the skill (run init_skill.py)
4. Edit the skill (implement resources and write SKILL.md)
5. Deliver the skill (send SKILL.md path via notify_user)
6. Iterate based on real usage

**Key Insight:** Skills are onboarding guides that transform general-purpose agents into specialized agents with procedural knowledge.`
}

func getStrategicScout() string {
	return `# Strategic Scout

Framework for strategic exploration and decision-making when facing product decisions or strategic tensions.

## When to Use
- Facing product decisions or strategic tensions
- Exploring multiple possible directions
- Need to map decision space before committing
- Coordinating high-level planning

## The Process

### 1. Recognize the Tension
Identify and articulate the strategic question or tension. Best work begins with a tension, not a command.

### 2. Scout Multiple Routes
Explore 2-4 distinct approaches with their tradeoffs. Don't rush to a solution.

### 3. Identify Decision Criteria
What would make you choose one path over another?

### 4. Recommend Smallest Test
What's the smallest experiment that provides useful signal?

### 5. Document the Landscape
Create a clear map of options, tradeoffs, and decision criteria.

**Core Insight:** Strategic work requires holding multiple possibilities before converging. Scout mode delays decision to enable better decision.`
}

func getPreImplementationChecklist() string {
	return `# Pre-Implementation Checklist

Systematic checklist to run before implementation begins.

## The Checklist

### 1. Specification Completeness
- [ ] Requirements are clear and comprehensive
- [ ] Success criteria are defined
- [ ] Edge cases are identified
- [ ] Dependencies are documented

### 2. Backend Grounding
- [ ] API contracts are defined
- [ ] Data models are documented
- [ ] Integration points are clear
- [ ] Backend is ready to support frontend

### 3. Technical Foundation
- [ ] File paths are identified
- [ ] Architecture is documented
- [ ] Patterns are established
- [ ] Tools and libraries are chosen

### 4. Resource Readiness
- [ ] All assets are available
- [ ] Documentation is accessible
- [ ] Examples are provided
- [ ] Support is available

### 5. Success Validation
- [ ] Testing strategy is defined
- [ ] Acceptance criteria are clear
- [ ] Review process is established
- [ ] Rollback plan exists

**Core Insight:** Time spent validating readiness prevents wasted implementation effort. A well-grounded specification enables autonomous execution.`
}

func getSkillMaintenanceRitual() string {
	return `# Skill Maintenance Ritual

Systematic process for maintaining skills directory health.

## When to Use
- Skill names become unclear or outdated
- Terminology needs updating
- Adding new skills that reference existing ones
- Conducting periodic audits
- Ensuring knowledge base accessibility

## The Workflow

### Step 1: Recognize the Need
Identify maintenance triggers and document what needs attention.

### Step 2: Read Current State
Review existing skills, names, cross-references, and terminology.

### Step 3: Propose Changes
Document proposed renames, refactors, and updates with rationale.

### Step 4: Execute Changes
Perform renames, update cross-references, and modify frontmatter.

### Step 5: Verify
Test that all changes work correctly and references resolve.

### Step 6: Document
Update changelog and communicate changes.

**Core Insight:** Good maintenance prevents future problems. A well-maintained knowledge base is easier to search, understand, and extend.`
}

func getStrategicToTacticalWorkflow() string {
	return `# Strategic-to-Tactical Workflow

Complete workflow from strategic scouting to tactical commission.

## The 8-Phase Workflow

### Phase 1: Recognize the Tension
Identify and articulate the strategic tension or question.

### Phase 2: Scout the Landscape
Explore multiple possible approaches with their tradeoffs.

### Phase 3: Propose a Direction
Recommend an approach with clear rationale.

### Phase 4: Write the Specification
Document requirements, architecture, and success criteria.

### Phase 5: Ground in Backend Reality
Ensure backend is ready to support the spec.

### Phase 6: Validate Readiness
Run pre-implementation checklist.

### Phase 7: Write Implementation Prompts
Transform spec into executable implementation prompts.

### Phase 8: Commission and Track
Hand off to implementation agents and track progress.

**Core Insight:** This workflow closes the loop from strategic thinking to tactical execution, accumulating wisdom at each cycle.`
}

func getTransformSpecToImplementationPrompt() string {
	return `# Transform Spec to Implementation Prompt

Transform architectural specifications into structured implementation prompts.

## When to Use
- After completing comprehensive specification
- When preparing to commission implementation agents
- When breaking down large specs into parallel tracks
- After completing backend grounding

## The Workflow

### Step 1: Validate Specification Readiness
Ensure the specification is complete and implementation-ready.

### Step 2: Identify Implementation Tracks
Break the work into logical, parallel or sequential chunks.

### Step 3: Write Structured Prompts
For each track, create a comprehensive prompt with:
- Context and background
- Requirements and constraints
- File paths and architecture
- Success criteria and testing
- Explicit boundaries

### Step 4: Add Backend Grounding
Include API contracts, data models, and integration points.

### Step 5: Review and Validate
Ensure prompts are autonomous, complete, and unambiguous.

### Step 6: Commission
Hand off prompts to implementation agents and track progress.

**Core Insight:** A well-written implementation prompt is an act of translation from strategic intent to tactical execution.`
}

func getSeedReflector() string {
	return `# Seed Reflector

Extract and document reusable patterns (seeds) from experiences.

## What Is a "Seed"?

A seed is a reusable pattern, insight, or principle that:
1. Emerged from experience (not abstract theory)
2. Can be applied in future contexts (not one-time specific)
3. Has a clear trigger (you know when to use it)
4. Captures wisdom (not just information)

## Seed Extraction Process

### Step 1: Identify Candidate Patterns
Look for:
- Decisions that worked well (or didn't)
- Patterns that emerged across multiple instances
- Insights that changed how you think
- Principles that guided successful outcomes

### Step 2: Test for Reusability
A good seed is:
- ✅ General enough to apply in multiple contexts
- ✅ Specific enough to be actionable
- ✅ Grounded in experience
- ✅ Has a clear trigger

### Step 3: Document the Seed
Include:
- Name and description
- Context where it emerged
- When to apply it
- How to apply it
- Expected outcomes

### Step 4: Test the Seed
Apply it in a new context to verify usefulness.

**Core Insight:** Every experience contains seeds. The practice is learning to see them, extract them, and plant them where they'll grow.`
}

func getMemoryGardenWriter() string {
	return `# Memory Garden Writer

Write structured, semantically rich memory entries for efficient context management.

## Memory Hierarchy

### Tier A: Raw Daily Notes
- Location: memory/YYYY-MM-DD.md
- Purpose: Capture everything from today
- Lifespan: 1-3 days before compression

### Tier B: Curated Wisdom
- Location: MEMORY.md (root level)
- Purpose: Distilled insights, decisions, patterns
- Lifespan: Permanent, but evolves

### Tier C: Compressed Archive
- Location: memory/archive/YYYY-MM.md
- Purpose: Historical record, rarely accessed
- Lifespan: Permanent, read-only

## Daily Memory Entry Structure

Include:
- Session context and duration
- Key activities with timestamps
- Decisions made with rationale
- Insights and learnings
- Related artifacts and links

## Writing Principles

- Be specific and concrete
- Include trigger conditions
- Tag for future retrieval
- Link to related memories
- Compress regularly

**Core Insight:** Memory should be a garden, not a landfill. Cultivate what matters, compost what doesn't.`
}

func getParallelTracksPattern() string {
	return `# Parallel Tracks Pattern

Split large development tasks into independent parallel tracks to maximize velocity.

## When to Use

- Task is large (>2 weeks if sequential)
- Clear separation of concerns exists
- Multiple agents/developers available
- Tracks have minimal dependencies
- Committed to writing clear specifications

## The Workflow

### Step 1: Identify Natural Boundaries
Common boundaries:
- By Layer: frontend, backend, database
- By Feature: auth, orchestration, UI
- By Component: foundation, features, integrations

Aim for 2-4 substantial tracks.

### Step 2: Define Track Dependencies
Create dependency graph:
- Independent tracks: start immediately
- Dependent tracks: wait for prerequisites

### Step 3: Write Self-Contained Specifications
Each spec must include:
- Clear goal and context
- Detailed requirements
- Success criteria
- Explicit non-goals

### Step 4: Define Integration Points
Be explicit about:
- API endpoints and contracts
- Component props and interfaces
- Shared state structures

### Step 5: Execute in Parallel
Commission independent tracks simultaneously.

### Step 6: Integrate and Test
Merge completed tracks and verify integration points.

**Core Insight:** Parallel execution requires upfront discipline in specification and architecture, but multiplies velocity without sacrificing quality.`
}

func getIterativeScoutingPattern() string {
	return `# Iterative Scouting Pattern

Strategic scouting as an iterative conversation: scout → feedback → reframe → re-scout.

## The Philosophy

Scouting is not linear—it's a conversation with the strategic landscape. The goal of the first scout is to provoke a deeper question, not find the final answer.

## The Workflow

### Step 1: Initial Scout
- Identify the initial tension
- Explore diverse routes
- Propose provocative starting point

### Step 2: Gather Feedback & Listen for Reframe
- Present initial routes
- Listen for the "question behind the question"
- Identify deeper framing

### Step 3: Re-Scout with New Lens
- Articulate reframed tension
- Explore routes native to new framing
- Generate richer options

### Step 4: Synthesize Final Vision
- Select best route from second round
- Define positioning and strategy
- Confirm alignment

## Best Practices

- **The Two-Scout Rule**: Assume at least two rounds for non-trivial decisions
- **The Reframe is the Prize**: The new question is more valuable than the first answer
- **Scout for Provocation**: Provoke better conversation, not consensus

**Core Insight:** Strategic clarity emerges through iteration, not from a single analysis.`
}

func getWriteFrontendSpecFromBackend() string {
	return `# Write Frontend Spec From Backend

Write production-ready frontend specifications deeply grounded in existing backend architecture.

## The Philosophy

Grounding before building. Most bugs, delays, and rework come from disconnect between frontend and backend.

## The Workflow

### Step 1: Deep Backend Analysis
- Generate repository context map
- Read key backend files (handlers, middleware)
- Document all relevant API endpoints
- Identify integration points

### Step 2: Comprehensive Feature Specification
- Write full feature spec with goals, requirements, architecture
- Leverage existing backend wherever possible
- Include API contracts with examples
- Define security considerations

### Step 3: Integration Guide Creation
- Create track-by-track guides
- Document authentication flow
- Explain streaming architecture (if applicable)
- Provide frontend code examples
- Specify error handling

### Step 4: Track Prompt Enhancement
- Add "Backend Grounding" section to all prompts
- Document specific endpoints for each feature
- Reference integration guide
- Ensure agents follow existing patterns

### Step 5: Validation Checkpoint
Run pre-implementation checklist before commissioning work.

**Core Insight:** Time spent understanding the backend prevents entire classes of integration problems.`
}

func getProductPositioningScout() string {
	return `# Product Positioning Scout

Reframe binary product decisions into strategic positioning opportunities.

## The Philosophy

Value is contextual. Instead of "keep or kill," ask: "What is this uniquely good at?"

## The Workflow

### Step 1: Identify the Binary Trap
Recognize limiting binary choices:
- "Keep X or get rid of it?"
- "Build X or buy Y?"
- "Deprecate or maintain?"

### Step 2: Introduce the Unlocking Question
Ask: **"What is this uniquely good at that the other thing isn't?"**

Shift from "what to do with this" to "what is its unique value?"

### Step 3: Explore the Unique Value
- List unique strengths
- Identify contexts where strengths shine
- Consider niche use cases

### Step 4: Translate to Product Positioning
Transform unique value into:
- Clear positioning statement
- Target audience definition
- Differentiated value proposition

### Step 5: Validate with Users
Test positioning with target audience to verify resonance.

**Core Insight:** The reframe is the prize. Transform "legacy" into "premium," "redundant" into "complementary."`
}

func getRetrospective() string {
	return `# Retrospective

Structured post-sprint learning and continuous improvement.

## The Philosophy

A retrospective is a harvest, not a post-mortem. It's about gratitude, honesty, and continuous learning.

## When to Use

- After major releases
- After significant milestones
- When projects feel stuck
- At regular intervals (monthly)

## The Workflow

### Step 1: Initiate the Retrospective
Frame it as positive and necessary part of workflow.

### Step 2: Create Retrospective Document
Use template in docs/retrospectives/.

### Step 3: Answer Three Core Questions

1. **What went well?** (What should we amplify?)
2. **What was hard?** (Sources of friction?)
3. **What would we do differently?** (Actionable changes?)

### Step 4: Synthesize and Extract Learnings
- Identify patterns and themes
- Distill actionable insights
- Create seeds for most profound lessons

### Step 5: Commit and Share
Document and integrate findings into next sprint.

**Core Insight:** Valuable lessons from each sprint should be harvested and integrated into shared memory, not lost.`
}

func getMultiSurfaceProductStrategy() string {
	return `# Multi-Surface Product Strategy

Design coherent multi-surface strategy where each surface has a unique, complementary role.

## Core Principle

Surfaces are for contexts, not devices. Each surface should be optimized for its unique context.

## The Workflow

### Step 1: Identify the Surfaces
List current and potential surfaces:
- Desktop, mobile, web, CLI, API, etc.

### Step 2: Define Context for Each Surface
For each surface, identify:
- Primary use context
- User mental state
- Time constraints
- Environmental factors

### Step 3: Assign Unique Jobs-to-be-Done
Each surface gets a clear mission:
- Desktop: Deep work and complex workflows
- Mobile: On-the-go access and notifications
- Web: Discovery and cross-platform onboarding

### Step 4: Design Integration Points
Define how users move between surfaces:
- Data sync mechanisms
- Handoff patterns
- Cross-surface workflows

### Step 5: Validate Coherence
Ensure the multi-surface experience feels unified yet appropriately differentiated.

**Core Insight:** Complement, don't compete. The whole is greater than the sum of its parts.`
}

func getContextCompressionRitual() string {
	return `# Context Compression Ritual

Systematic process for compressing agent context and memory.

## When to Use

- Context feels overwhelming
- Memory is degrading or cluttered
- Periodic maintenance (weekly/monthly)
- Before major context switches

## The Compression Hierarchy

### Tier A → Tier B (Daily → Curated)
- Read recent daily notes
- Extract key insights and decisions
- Update curated wisdom in MEMORY.md
- Tag and organize for retrieval

### Tier B → Tier C (Curated → Archive)
- Identify what hasn't been accessed in 3+ months
- Create semantic summaries
- Move to archive with timestamp
- Keep index for future reference

## Compression Principles

- Preserve decisions and rationale
- Extract reusable patterns (seeds)
- Maintain retrievability through tags
- Compress without losing essence

**Core Insight:** Regular compression prevents context overload and maintains a healthy, navigable memory garden.`
}

func getAgentHandoffProtocol() string {
	return `# Agent Handoff Protocol

Structured protocol for handing off work between agents.

## When to Use

- Commissioning implementation agents
- Delegating tasks to specialists
- Transitioning work between sessions
- Coordinating multi-agent workflows

## Handoff Structure

### 1. Context Package
Provide:
- Current state and background
- Relevant files and artifacts
- Previous decisions and rationale

### 2. Clear Objective
Define:
- Specific goal or deliverable
- Success criteria (testable)
- Constraints and boundaries

### 3. Resources and Tools
Include:
- Required tools and access
- Reference documentation
- Code examples and patterns

### 4. Integration Points
Specify:
- How output integrates with existing work
- API contracts or interfaces
- Review and feedback process

### 5. Recovery Plan
Document:
- How to handle blockers
- When to escalate
- Rollback procedures

**Core Insight:** Clear handoffs prevent context loss and enable autonomous execution.`
}

func getResearchModes() string {
	return `# Research Modes

Framework for different research modes matched to research goals.

## The Four Modes

### 1. Exploratory Research
**Goal:** Discover what exists in a domain
**Approach:** Broad scanning, many sources
**Output:** Landscape map, key themes

### 2. Focused Research
**Goal:** Deep understanding of specific topic
**Approach:** Systematic deep dives
**Output:** Comprehensive analysis, detailed notes

### 3. Comparative Research
**Goal:** Evaluate options and alternatives
**Approach:** Side-by-side analysis with criteria
**Output:** Decision matrix, recommendations

### 4. Synthesis Research
**Goal:** Integrate findings into coherent framework
**Approach:** Pattern identification, framework building
**Output:** Synthesis document, actionable insights

## Choosing the Right Mode

- Starting new domain? → Exploratory
- Need deep expertise? → Focused
- Making decisions? → Comparative
- Creating strategy? → Synthesis

**Core Insight:** Match research mode to research goals for efficient learning.`
}

func getDebuggingTroubleshooting() string {
	return `# Debugging and Troubleshooting

Systematic debugging workflow.

## The Process

### 1. Reproduce
- Create minimal reproduction case
- Document exact steps
- Identify trigger conditions

### 2. Isolate
- Remove variables systematically
- Test components independently
- Narrow scope to smallest unit

### 3. Diagnose
- Gather evidence (logs, errors, traces)
- Form hypotheses
- Test hypotheses systematically

### 4. Fix
- Implement targeted fix
- Avoid scope creep
- Document rationale

### 5. Verify
- Test fix in isolation
- Test in full context
- Verify no regressions

### 6. Document
- Record root cause
- Document solution
- Update troubleshooting guides

## Debugging Principles

- Start with simplest hypothesis
- Change one thing at a time
- Keep detailed notes
- Don't skip verification

**Core Insight:** Systematic debugging is faster than random changes. Patience and discipline prevent wasted effort.`
}

func getProcessToSkillWorkflow() string {
	return `# Process-to-Skill Workflow

Meta-skill for transforming valuable workflows into reusable skills.

## The Philosophy

Valuable multi-step processes represent implicit knowledge. This skill makes that knowledge explicit, transforming one-off processes into repeatable practices.

## When to Use

- After completing complex, multi-step tasks that will be repeated
- When manually repeating same sequence across projects
- During retrospectives identifying successful workflows
- When onboarding agents to complex processes

## The Workflow

### Step 1: Identify and Document the Process
- Select recently completed successful workflow
- Document each step with goals, actions, tools, inputs, outputs
- Extract key insights and reusable patterns
- Be detailed about non-obvious steps

### Step 2: Convert Process to Skill
- Use seed-to-skill-converter to transform documentation
- Deconstruct into: insight, trigger, process, outcome
- Draft SKILL.md following standard template

### Step 3: Refine and Enhance
- Add quality checklist
- Document best practices and pitfalls
- Create bundled resources (scripts, templates, references)

### Step 4: Validate and Deliver
- Run validation checks
- Push to repository
- Deliver as .skill file

## Best Practices

- Focus on the "why" not just the "what"
- Generalize the pattern for reusability
- Use templates provided
- Iterate after real-world use

**Core Insight:** Implicit knowledge becomes institutional memory through systematic documentation and formalization.`
}

func getSeedToSkillConverter() string {
	return `# Seed-to-Skill Converter

Convert valuable Dojo Seeds (insights) into fully-fledged reusable Skills.

## The Philosophy

A Seed is potent insight—a moment of clarity captured. A Skill is an instrument—that same lesson transformed into repeatable, structured process.

From passive wisdom to active utility.

## When to Use

- Seed is referenced frequently across projects
- Seed describes multi-step process
- Seed represents core part of workflow
- During retrospectives identifying important learnings

## The Conversion Workflow

### Step 1: Identify Candidate Seed
Select seed that meets criteria above.

### Step 2: Deconstruct the Seed's Wisdom
Break down core components:
- **Core Insight**: Fundamental truth or idea
- **Trigger**: When to apply this wisdom
- **Process**: Concrete steps to take
- **Desired Outcome**: Result of correct application

### Step 3: Draft Skill Using Standard Template
Map seed components to skill sections:
- Core Insight → I. The Philosophy
- Trigger → II. When to Use This Skill
- Process → III. The Workflow
- Desired Outcome → IV. Best Practices / V. Quality Checklist

### Step 4: Define Workflow and Templates
Transform abstract process into concrete step-by-step workflow.

### Step 5: Commit New Skill
Add to repository and make available for use.

## Best Practices

- Not every seed needs to be a skill—only promote proven value
- Skills must be actionable processes
- Skills require maintenance
- Goal is utility, not documentation

**Core Insight:** The alchemical process that turns insight into instrument—wisdom into practice.`
}

func getRepoContextSync() string {
	return `# Repo Context Sync

Efficiently sync and extract context from GitHub repositories for grounding work.

## The Philosophy

Architectural decisions in a vacuum are fragile. This skill embodies **grounding in reality**—syncing with actual codebase state before decisions or specifications.

Not reading every file, but **surgical context extraction**.

## When to Use

Trigger when:
- User mentions "repo", "repository", "codebase", "sync"
- Starting refactoring, architecture, or design conversations
- Writing implementation prompts needing existing patterns
- Need to understand current state of codebase

## Core Workflow

### 1. Identify Context Need
Parse request to determine:
- Which repos are relevant
- Which directories matter
- What keywords indicate focus areas

### 2. Sync Repo State
Use sparse checkout for efficiency:
- Clone only needed directories
- Track changes with diff summaries
- Avoid full repo download

### 3. Track Changes
Generate diff summary showing:
- Added, modified, deleted files
- Commit messages
- Summary statistics

### 4. Generate Context Map
Create comprehensive codebase overview:
- File structure analysis
- Pattern detection
- Key file identification
- Keyword-focused extraction

## Integration Points

- Works with implementation agent workflows
- Supports "Planning with Files" philosophy
- Enables grounded architectural decisions

**Core Insight:** Time spent syncing with reality prevents hallucination and poor architectural decisions.`
}

func getProjectExploration() string {
	return `# Project Exploration

Structured process for exploring new large-scale projects to assess collaboration.

## The Philosophy

Onboarding to new domains requires systematic exploration with progressive disclosure. Don't try to understand everything at once—use phased approach.

## When to Use

- Evaluating new project for collaboration
- Onboarding to unfamiliar domain
- Due diligence for partnership
- Understanding large-scale system architecture

## Five-Phase Approach

### Phase 1: Discovery
- Identify project boundaries
- Find key repositories
- Locate documentation
- Identify stakeholders

### Phase 2: Sampling
- Read representative files
- Understand patterns
- Identify conventions
- Sample different layers

### Phase 3: Research
- Deep dive into core areas
- Understand architecture
- Map dependencies
- Research technologies

### Phase 4: Connection
- Identify integration points
- Understand workflows
- Map communication patterns
- Find collaboration opportunities

### Phase 5: Synthesis
- Create comprehensive overview
- Document learnings
- Assess collaboration fit
- Propose engagement model

## Principles

- Progressive disclosure—layer by layer
- Representative sampling—not exhaustive reading
- Pattern recognition—identify conventions
- Context over completeness

**Core Insight:** Systematic exploration prevents overwhelm and enables informed collaboration decisions.`
}

func getAgentWorkspaceNavigator() string {
	return `# Agent Workspace Navigator

Best practices for navigating and collaborating in shared agent workspaces.

## The Philosophy

Shared workspaces require structure to prevent chaos. Standard patterns enable efficient collaboration and token usage.

## When to Use

- Setting up multi-agent collaboration
- Organizing shared workspace
- Establishing file conventions
- Coordinating between agents

## Standard Directory Structure

### Core Directories
- **docs/** - Documentation and specifications
- **memory/** - Agent memory and context
- **planning/** - Plans and task breakdowns
- **thinking/** - Private reflection space
- **scripts/** - Automation and tools

### File Naming Conventions
- Use kebab-case: task-name.md
- Include dates for temporal files: 2026-02-09-notes.md
- Prefix with purpose: spec-feature-name.md

## Reading Patterns

1. **Start with STATUS.md** - Current project state
2. **Read PHILOSOPHY.md** - Project values and approach
3. **Check memory/** - Recent context
4. **Review planning/** - Current tasks

## Writing Patterns

1. **Update STATUS.md** after significant changes
2. **Write to memory/** for persistence
3. **Use thinking/** for private reflection
4. **Keep docs/** for specifications

## Collaboration Workflows

- Handoffs: Clear context packages
- Updates: Status notifications
- Questions: Explicit asks in shared space
- Decisions: Document in memory/

**Core Insight:** Structure enables collaboration; chaos prevents it. Standard patterns are shared language.`
}

func getWriteReleaseSpecification() string {
	return `# Write Release Specification

Production-ready specification writing for software releases.

## The Philosophy

Good specifications prevent entire classes of problems. Time spent planning is time saved in implementation.

## When to Use

- Planning major software releases
- Creating technical specifications
- Designing new systems
- Documenting complex features

## Specification Structure

### Executive Summary
- What we're building and why
- Key goals and non-goals
- Success criteria

### Technical Architecture
- System design
- Component breakdown
- Integration points
- Technology choices

### Implementation Plan
- Phases and tracks
- Dependencies
- Timeline estimates
- Resource requirements

### Risk Management
- Technical risks
- Dependencies
- Mitigation strategies
- Rollback plans

### Success Criteria
- Functional requirements
- Performance targets
- Quality standards
- Acceptance tests

## Best Practices

- Be comprehensive but concise
- Include diagrams and examples
- Define clear boundaries (non-goals)
- Plan for failure modes
- Document assumptions

**Core Insight:** Specification quality directly determines implementation success. Invest upfront to save later.`
}

func getHealthSupervisor() string {
	return `# Health Supervisor

Systematic workflow for conducting comprehensive health audits on repositories.

## The Philosophy

Repository health degrades without maintenance. Systematic audits identify technical debt before it becomes crisis.

## When to Use

- Periodic repository health checks
- Before major refactoring
- Technical debt assessment
- Onboarding to existing codebase

## Audit Framework

### 1. Code Quality
- Test coverage
- Linting and formatting
- Documentation completeness
- Code complexity

### 2. Architecture Health
- Dependency management
- Design pattern consistency
- Modularity assessment
- Technical debt inventory

### 3. Security Posture
- Dependency vulnerabilities
- Security best practices
- Access controls
- Audit logging

### 4. Operations Readiness
- CI/CD pipeline health
- Deployment documentation
- Monitoring and alerting
- Disaster recovery

### 5. Developer Experience
- Setup documentation
- Development workflows
- Tool integration
- Contribution guidelines

## Output

Generate actionable engineering tasks:
- Priority classifications
- Effort estimates
- Dependencies
- Acceptance criteria

**Core Insight:** Regular health audits prevent technical debt from becoming technical crisis.`
}

func getWebResearch() string {
	return `# Web Research

Effective web research using search APIs and content extraction.

## The Philosophy

Good web research is surgical—targeted queries, credible sources, extracted insights, not content dumping.

## When to Use

- Investigating topics requiring current information
- Finding sources for research or specifications
- Verifying claims or facts
- Gathering competitive intelligence
- Understanding technologies or practices
- Finding documentation or examples

## Research Workflow

### Step 1: Define Research Question
Clarify:
- What specific question am I answering?
- What level of detail is needed?
- Is this time-sensitive?
- What sources would be most authoritative?

### Step 2: Formulate Search Queries
Effective query patterns:
- Specific terms: "React 18 suspense API"
- Questions: "how to implement authentication in Go"
- Comparisons: "PostgreSQL vs MySQL performance 2026"
- Recent: Use time filters for current info

### Step 3: Evaluate Sources
Priority hierarchy:
1. Official documentation
2. Technical blogs from credible authors
3. Stack Overflow (for specific solutions)
4. Academic papers (for deep understanding)
5. Forums and discussions (for context)

### Step 4: Extract and Synthesize
- Don't just collect links
- Extract key insights
- Synthesize across sources
- Note contradictions or gaps
- Document source for each insight

### Step 5: Document Findings
Include:
- Research question
- Key findings with citations
- Source URLs and dates
- Confidence level (high/medium/low)
- Gaps or areas needing deeper research

## Research Patterns

**Landscape Scan:**
- Broad queries to understand domain
- Multiple sources for overview
- Identify key concepts and terminology

**Deep Dive:**
- Focused queries on specific topics
- Technical documentation
- Implementation examples

**Verification:**
- Cross-reference multiple sources
- Check publication dates
- Verify author credibility

**Core Insight:** Research quality depends on query precision and source evaluation, not volume of content collected.`
}

func getStatusWriter() string {
	return `# Status Writer

Write and update STATUS.md files for at-a-glance project visibility.

## The Philosophy

A STATUS.md is a ritual of bearing witness—radical honesty about where project truly is, not where we wish it were.

Single source of truth that grounds conversations in reality.

## When to Use

- At beginning of new project
- At start and end of work sessions
- During weekly syncs
- When significant status changes occur

## Status Update Workflow

### Step 1: Locate or Create STATUS.md
In project root. Use template if creating new.

### Step 2: Update the Header
Change "Last Updated" to current date.

### Step 3: Review and Update Sections

**Vision & Purpose:**
- Rarely changes
- Re-read to stay grounded

**Current State:**
- Most important section
- Use emoji status: ✅ (complete), 🔄 (in progress), ⏸️ (blocked), ❌ (not started)
- Add/remove items as project evolves

**Active Workstreams:**
- What's being worked on right now
- Update tasks and progress

**Blockers & Dependencies:**
- Be ruthlessly honest
- Remove resolved blockers
- Add new blockers immediately

**Next Steps:**
- Concrete, actionable items
- Immediate next action

### Step 4: Commit Changes
Commit message: docs(status): Update [Project] status for [Date]

## Status Template Structure

Include these sections:
- Header with project name and last updated date
- Vision & Purpose (one paragraph)
- Current State (with emoji status indicators)
- Active Workstreams (with progress)
- Blockers & Dependencies
- Next Steps (actionable items)

Use emoji consistently:
- ✅ Complete
- 🔄 In progress
- ⏸️ Blocked
- ❌ Not started

## Best Practices

- Update regularly (start/end of sessions)
- Be honest about blockers
- Keep it concise (< 1 page)
- Use emoji consistently
- Remove resolved items
- Make next steps actionable

**Core Insight:** Status transparency prevents confusion and enables informed decisions. The ritual creates accountability.`
}

func getDecisionPropagationProtocol() string {
	return `# Decision-Propagation Protocol

## I. The Philosophy

Architectural decisions are not isolated events. When a decision arrives—especially to an open question in an existing document—it creates ripples across an entire ecosystem of interdependent files. A decision like "auth bypass" doesn't just answer one scout's question; it changes scope in the master plan, triggers deferrals in other scouts, and reshapes the STATUS file's summary of the architecture.

Without a propagation protocol, decisions become stranded. A spec answers a question, but the master plan still lists it as open. A scout defers work, but the dependency graph hasn't updated. The system becomes incoherent—different documents contradict each other about what is decided and when.

This skill prevents that dissonance by treating decision-propagation as a deliberate, multi-document process. The human provides the decision once, in one place. Your job is to trace where that decision echoes and update each location.

## II. When to Use This Skill

Use the decision-propagation protocol when:

- A human provides answers to open questions in a scout or specification
- An architectural decision changes the scope of a release or work track
- A decision in one document affects the content, sequencing, or priorities of others
- You need to defer work from one release to a later one based on new information
- The STATUS file or master tracking document is out of sync with the decisions living in specs and scouts
- Any time a decision made in one place must be reflected in many places

## III. The Workflow: Five Steps

### Step 1: Record Decisions at Source

Locate the document where the decision was made (scout, spec, or master plan section).

Replace the "Open Questions" section with "Decisions ([Name], [Date])". Number each decision:

## Decisions (Cruz, 2026-02-11)

1. Entity Backend is v0.2.0 Priority
   - Decision: Build entity-centric backend first, before other subsystems.
   - Reasoning: Reduces scope of v0.2.1 and focuses team on core abstraction.
   - Implication: Defers non-entity tasks to v0.2.1; reshuffles parallel track allocation.

2. Auth Bypass
   - Decision: Skip production auth layer in v0.2.0; ship with debug token instead.
   - Reasoning: Unblocks frontend integration work; auth hardening moves to v0.2.2.
   - Implication: Changes v0.2.0 scope (removes auth), changes v0.2.2 scope (adds auth).

For each decision, capture the exact words, the reasoning (why the human chose this), and any stated implications. If the decision adds or removes scope, note it explicitly.

### Step 2: Trace Document Dependencies

Before editing anything, make the complete list of documents affected. Common dependency patterns:

- **Master plan** → typically affected in: scope blocks, dependency graph, constraints, parallel track allocation, next steps
- **Other scouts** → may be promoted (dependencies now met), deferred (dependencies removed or pushed later), or need content updates
- **Implementation specs or prompts** → may need regeneration if scope changed
- **STATUS.md** → always needs updating; it's the coherence checkpoint

Walk through each document and ask: "Does this reference the decision I just recorded?" If yes, it's dependent.

### Step 3: Propagate to Each Dependent Document

For each dependent document, make surgical edits. Do not rewrite entire sections unless necessary.

**In master plans:** Update scope tables, dependency graph (remove/add edges), constraint lists, parallel track allocation.

**In other scouts:** If a decision defers work, add a new section noting the deferral, the reasoning (reference the decision), and the proposed timeline. If a decision enables a scout, promote it in priority or scheduling.

**In STATUS.md:** Add a "Key Architecture Decisions" block and list all decisions with their implications. Update any summary statements about scope or sequencing.

**In implementation prompts:** If scope changed, the prompt may need regeneration to reflect new boundaries.

### Step 4: Update Master Tracking

Update STATUS.md last. It's the summary of everything else. Include:

- The new decisions block (attribution, date, full text)
- Any tables listing specs, scouts, or prompts (update affected rows)
- "Next Steps" section (reorder if sequencing changed)
- Summary statement about architecture or scope (refresh if affected)

### Step 5: Sync Copies

If documents exist in multiple locations (e.g., thinking/ at repo root AND docs/v0.2.x/thinking/), sync them. Compare file sizes or use diff to verify.

## IV. Document Dependency Patterns

Understanding these patterns helps you trace dependencies quickly:

- **Decisions flow upward:** A scout decision affects the master plan above it
- **Decisions cascade sideways:** A decision in one scout defers or enables tasks in other scouts
- **Scope is bidirectional:** Adding scope to v0.2.0 removes it from v0.2.1; changes must be reflected in both
- **Deferrals create cross-references:** When a scout says "deferred to v0.2.3," the v0.2.3 master plan must acknowledge the deferred item

## V. Best Practices

- **Record verbatim.** Capture the human's exact words, not your paraphrase. Paraphrasing introduces interpretation and drift.
- **Trace before editing.** Make the full list of dependent documents first. This prevents missed updates and duplicated effort.
- **Surgical over wholesale.** Change only the specific sections affected. Rewriting entire documents introduces risk and obscures what actually changed.
- **Always update STATUS.md last.** It's the final coherence checkpoint. If STATUS.md is current, the whole system is current.
- **Document deferrals clearly.** When work is pushed to a later release, add a note explaining why and when it should be reconsidered.

## VI. Quality Checklist

Before considering the decision propagated:

- [ ] Decision recorded at source with human name, date, and full reasoning
- [ ] All dependent documents identified (master plan, related scouts, STATUS.md, implementation prompts)
- [ ] Each dependent document updated in specific affected sections (not blanket rewrites)
- [ ] STATUS.md reflects all changes and serves as a coherence checkpoint
- [ ] Document copies synced across locations if they exist in multiple places
- [ ] No orphaned references to old scope or pre-decision state remain in any document
- [ ] Cross-references between documents are consistent (if one document says "deferred to v0.2.3," the v0.2.3 plan acknowledges it)`
}

func getEraArchitecture() string {
	return `# era-architecture

## I. The Philosophy

An era is not a list of releases. It's a conceptual architecture that coordinates multiple releases under shared constraints.

When you go from "we need to rebuild the frontend" to "we need a complete multi-release era plan," you're making a qualitative shift. Instead of asking "what should v0.2.1 do?" you're asking "what vocabulary, architectural patterns, and design language will span all releases in this era?"

This shift matters because:
- **Coherence**: Without a shared architecture, releases drift apart. Release A uses one data model, Release B contradicts it.
- **Dependency clarity**: You can't plan releases independently. A shared design language in Release 1 enables faster work in Release 2.
- **Decision efficiency**: The era's conceptual architecture is the constitution that downstream decisions must obey. This prevents constant re-litigation of foundational choices.

The key insight: scouts explore the facets (not releases). From those explorations, you extract the era's architecture (not the other way around).

## II. When to Use This Skill

**Use era-architecture when:**
- Planning a major product pivot that spans 3+ releases
- You need to define shared conceptual vocabulary before decomposing into releases
- Transitioning between product eras (e.g., engine-building → fresh shell → social layer)
- The scope is too large for one release but needs coherent vision across all releases
- You need cross-cutting architectural decisions (data model, design language, navigation patterns)

**Do NOT use era-architecture when:**
- Planning a single release → use strategic-scout + parallel-tracks-pattern
- Writing specs → use specification-writer or write-frontend-spec-from-backend
- Converting specs to prompts → use spec-constellation-to-prompt-suite or zenflow-prompt-writer
- Exploring one strategic question → use strategic-scout

## III. The Workflow (7 Steps)

### Step 1: Name the Era and Its Predecessor

Give the era a name that captures its purpose (e.g., "The Fresh Shell", "The Engine", "The Social Layer").

Document what the previous era accomplished and what it leaves behind. Define the handoff: what assets/code/patterns carry forward, what gets archived.

This framing prevents the new era from being aimless. It grounds the work in what came before and what you're building toward.

### Step 2: Scout the Facets

Run 3-5 strategic scouts, each exploring a different facet of the era. Facets are not releases—they're architectural questions:
- Navigation/information architecture
- Data model and persistence
- Design language and UX philosophy
- Backend integration strategy
- Social/collaboration strategy (if applicable)

Each scout produces insights and open questions. Gather human decisions on those open questions before proceeding to Step 3. This is where the team shapes the architecture.

### Step 3: Define the Conceptual Architecture

From the scout insights and decisions, extract the era's conceptual vocabulary:
- **Core metaphors**: What's the guiding image? (e.g., "garden", "cognitive routes", "bonsai")
- **Architectural patterns**: What repeating structures apply everywhere? (e.g., "entity-as-lens", "content-edge transitions")
- **Design language**: Colors, typography, interaction patterns that all releases obey
- **Data model**: The types and relationships that span all releases

Write these as constraints: every release in the era must obey them. This is the era's constitution—it doesn't change per release.

### Step 4: Decompose into Releases

Define 3-7 releases that build on each other. Each release should have:
- A clear theme/name (e.g., "The Foundation", "The Shallow End", "The Deep End")
- Scope: what's built
- Dependencies: what must exist first
- Estimated effort

The first release is the thickest (it builds foundations). Later releases are thinner because the patterns are established. Some releases can run in parallel if they don't share dependencies.

### Step 5: Build the Dependency Graph

Map release-level dependencies (which releases block which). Map spec-level dependencies within releases (which specs must exist before others). Identify the critical path and parallel opportunities.

### Step 6: Write the Master Plan

Use the Master Plan Template (Section IV). The master plan is the single commissioning document for the entire era. Include vision, conceptual architecture summary, release roadmap, dependency graph, constraints, and parallel track allocation.

The master plan is a living document—decisions from later scouts update it. But the core architecture should remain stable.

### Step 7: Drop into Single-Release Work

For the first release, write detailed specs. Then convert specs to implementation prompts. Commission and execute the first release.

Before starting the next release, check: do the conceptual architecture constraints still hold? Do decisions from the first release affect later releases? Update the master plan if needed.

## IV. Master Plan Template

# [Era Name] Master Plan
**Author:** [names]
**Status:** [Active/Complete]
**Era:** v[X].0 through v[X].N

## 1. Vision
[1-2 sentences: what this era accomplishes]

## 2. What Came Before
[Brief summary of previous era, what carries forward]

## 3. Conceptual Architecture
### Core Metaphors
### Architectural Patterns
### Design Language
### Data Model

## 4. Release Roadmap
| Version | Name | Focus | Dependencies | Est. Effort |

## 5. Dependency Graph
[ASCII or description showing release ordering]

## 6. Constraints
[Rules that all releases must obey]

## 7. Parallel Track Allocation
| Release | Tracks | Can Parallel With |

## 8. Commissioned Specifications
| Spec | Release | Status |

## 9. Next Steps
[Current focus and what follows]

## V. Best Practices

- **Eras need names**: Names create shared vocabulary for the team. "The Fresh Shell" is better than "v0.2 refactor".
- **Scout facets, not releases**: Architecture emerges from cross-cutting exploration, not from asking "what should v0.2.1 do?"
- **Conceptual architecture is the output**: It constrains everything downstream and prevents releases from drifting apart.
- **First release is heaviest**: Later releases feel lighter because the patterns are established.
- **Master plans are living documents**: Use decision protocols when decisions arrive that change scope.
- **Don't over-plan later releases**: Define theme and rough scope, but save detailed specs for when earlier releases are done and the architecture has been validated in practice.

## VI. Quality Checklist

- [ ] Era has a name and a clear relationship to its predecessor
- [ ] 3+ facets scouted with insights and decisions recorded
- [ ] Conceptual architecture defined (metaphors, patterns, design language, data model)
- [ ] 3-7 releases defined with themes, scope, and dependencies
- [ ] Dependency graph shows ordering and parallel opportunities
- [ ] Master plan written as single commissioning document
- [ ] First release has detailed specs ready for implementation
- [ ] Constraints are clear enough that any release can check compliance`
}

func getSpecConstellationToPromptSuite() string {
	return `# spec-constellation-to-prompt-suite

## I. The Philosophy

When you move from one specification to one implementation prompt, you're writing a focused brief. But when you move from four interconnected specifications to four parallel implementation prompts, you're solving a **coordination problem**.

The danger: each prompt looks coherent in isolation, but the tracks execute independently. Track B (Go backend) produces the API that Track C (frontend) consumes—but if the two prompts define the contract differently, integration fails. Track A (scaffold) defines design tokens that Track D (dock + home) must use—but inconsistent naming breaks everything. Type definitions (DojoEntity, PipelineContent) must be identical across tracks, not compatible or "close enough."

This skill exists because:
- Writing N prompts from N specs is not N repetitions of "write one prompt"
- The coordination happens at the integration boundaries, not within prompts
- Shared types and contracts must be defined BEFORE prompts are written
- Cross-validation is the difference between "tracks that exist" and "tracks that integrate"

## II. When to Use This Skill

- You have 2+ interconnected specifications covering different aspects of a release
- Prompts for different tracks must agree on APIs, types, exports, or design tokens
- You're moving from "specs documented" to "ready to commission parallel work"
- One track's output becomes another track's input (explicit handoff required)
- You need a master plan that ties implementation work together
- Risk: Tracks ship independently without realizing their contracts don't align

## III. Prerequisites

- A master plan or release plan defining the release scope
- 2+ specifications covering different aspects of the release
- A parallel tracks decomposition (or the intent to create one)
- Codebase patterns extracted for grounding (handler patterns, response formats, type definitions)
- At least one existing codebase to extract patterns from

## IV. The Workflow

### Step 1: Map the Spec Constellation

List all specs that contribute to this release. For each spec, identify:
- Which aspects map to which track(s)
- Which other specs it depends on or feeds
- What decisions or scouts inform it

Create a matrix:
| Spec | Contributes To | Dependencies | Key Decision |
|------|---|---|---|
| Entity Data Model | Track B, C, D | — | Entity CRUD shape |
| Tauri Shell | Track A | — | Platform scaffold |
| Horizon Dock | Track D | Tauri Shell, Entity Model | UI layout |
| Home State | Track D | Entity Model | State shape |

**Key insight:** Specs don't map 1:1 to tracks. One spec feeds multiple tracks; multiple specs feed one track.

### Step 2: Define Integration Interfaces

For each pair of tracks that share data, define the contract in this format:

**Integration Contract Table:**
| Interface | Produced By | Consumed By | Shape | Verification |
|---|---|---|---|---|
| EntityAPI | Track B | Track C | listEntities(route): Promise<DojoEntity[]> | Type check against export |
| DojoEntity type | Track B | Track C, D | TypeScript interface | Shared in all 3 prompts identically |
| Route handlers | Track A | Track D | Path format: /entity/:id | Navigation code type-checks |

Write each contract in the Integration Contract Template below. **This table becomes the spine of the master plan.**

### Step 3: Extract Codebase Patterns

For each track that touches existing code:

- **Backend tracks:** Extract handler patterns, request/response shapes, middleware hooks
- **Frontend tracks:** Extract component folder structure, hook patterns, state management shapes
- **All tracks:** Response formats, error handling, logging conventions

Store patterns as reference files in the prompt (or link to them). Ground each prompt in "this is how we actually do it," not assumptions.

### Step 4: Write Prompts in Dependency Order

Begin with the track that has **no dependencies** (often scaffold or backend). For each prompt, include:

1. **Context & Grounding:** Which specs to read, which pattern files to reference
2. **Detailed Requirements:** Numbered, specific, unambiguous
3. **File Manifest:** Create vs. modify; where files go
4. **Success Criteria:** Binary, testable (not "looks good")
5. **Constraints & Non-Goals:** What NOT to do
6. **Integration Contracts:** What this track produces; what it depends on (reference Step 2)

**When writing a consuming track's prompt**, cite the producing track's contract explicitly:
> "Track B produces the DojoEntity type (Integration Contract: DojoEntity type). Use that type exactly in Track C's API client."

### Step 5: Cross-Validate Type Consistency

For every shared type (DojoEntity, PipelineContent, etc.):
- Check that the definition appears identically in every prompt that references it
- If a type evolves, update it in ALL prompts and the Integration Contract Table

For every API contract:
- Verify producer's response shape matches consumer's expected input shape
- Check HTTP method, status codes, error responses

For shared resources (CSS variables, design tokens, route patterns):
- Verify naming is consistent across all prompts

**Catch-and-fix:** This step prevents integration bugs.

### Step 6: Write the Master Implementation Plan

Summarize all tracks, dependencies, and integration points:

1. List all tracks in dependency order
2. Include the Integration Contract Table from Step 2
3. Define the day-by-day execution sequence
4. List aggregate success criteria spanning tracks
5. Note risk mitigation (e.g., "Track B ships Friday; Track C unblocked Monday")

## V. Integration Contract Template

### Integration Contract: [Name]

- **Producer:** Track [X] — [brief description of what produces this]
- **Consumer:** Track [Y] — [brief description of what consumes this]
- **Interface:** [API endpoint / TypeScript type / CSS variable / Go struct / etc.]
- **Language:** [TypeScript / Go / CSS / etc.]
- **Shape:**
  interface EntityAPI {
    listEntities(route: string): Promise<DojoEntity[]>
    getEntity(id: string): Promise<DojoEntity>
  }
- **Constraints:** [Any limits on the interface — e.g., "Must serialize to JSON", "Case-sensitive names"]
- **Verification:** [How to verify the contract is met — e.g., "Track Y's API client type-checks against Track X's response"]
- **Owner:** Track [X] (responsible for backward compatibility)

## VI. Best Practices

- **Specs don't map 1:1 to tracks.** Build the mapping explicitly in Step 1.
- **Integration contracts are the most important output.** Prompts without contracts produce tracks that don't integrate.
- **Write the most-depended-on track's prompt first.** Its contracts shape everything downstream.
- **Cross-validation catches bugs before execution.** Type mismatches, response shape mismatches, and naming inconsistencies surface now, not during integration.
- **If a spec needs updating mid-stream,** do it in the prompt, not by rewriting the spec. Prompts are executable; specs are reference.
- **Each prompt must be self-contained.** An agent should execute it without reading other prompts. Integration contracts embedded in the prompt provide necessary context.

## VII. Quality Checklist

- [ ] All specs contributing to this release are mapped to tracks (Step 1)
- [ ] Every cross-track dependency has an explicit integration contract (Step 2)
- [ ] Shared types (DojoEntity, etc.) are consistent across all prompts that reference them (Step 5)
- [ ] Each prompt is self-contained and executable independently
- [ ] Prompts are ordered by dependency (foundation tracks first) (Step 4)
- [ ] Integration Contract Table is complete and unambiguous (Step 2)
- [ ] Codebase patterns extracted and grounded in prompts (Step 3)
- [ ] Master implementation plan ties all tracks together (Step 6)
- [ ] Success criteria span tracks, not just per-track (Step 6)
- [ ] Risk mitigation plan addresses integration dependencies`
}
