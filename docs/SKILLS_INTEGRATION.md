# Dojo Genesis Skills Integration

**Date:** February 9, 2026
**Version:** 2.1

## Overview

This document describes the integration of Dojo Genesis skills into the dojo-mcp-server, extending the server's capabilities with actionable workflows and procedural knowledge from the skills directory.

## What Was Added

### 1. Skill Frontmatter Cleanup (dojo-genesis/skills)

Fixed YAML frontmatter in skills to follow the skill-creator standard:

**Fixed Skills:**
- `agent-to-agent-teaching` - Removed 8 excessive fields (version, author, created, type, for, trigger, dependencies, related_skills)
- `patient-learning-protocol` - Removed 8 excessive fields (same as above)
- `skill-maintenance-ritual` - Added missing frontmatter
- `strategic-to-tactical-workflow` - Added missing frontmatter
- `transform-spec-to-implementation-prompt` - Added missing frontmatter

**Standard Frontmatter Format:**
```yaml
---
name: skill-name
description: Clear description of what the skill does and when to use it
---
```

### 2. Skills Module (`internal/wisdom/skills.go`)

Created a new skills module in the wisdom package that:

- Defines the `Skill` struct with Name, Description, Category, and Content
- Implements `ListSkills()`, `GetSkill(name)`, and `SearchSkills(query)` methods on the Base
- Includes 20 core skills from dojo-genesis organized by category:
  - **learning** (2): agent-to-agent-teaching, patient-learning-protocol
  - **reflection** (2): seed-reflector, retrospective
  - **strategy** (4): strategic-scout, iterative-scouting-pattern, product-positioning-scout, multi-surface-product-strategy
  - **process** (5): pre-implementation-checklist, transform-spec-to-implementation-prompt, write-frontend-spec-from-backend, research-modes, debugging-troubleshooting
  - **workflow** (3): strategic-to-tactical-workflow, parallel-tracks-pattern, agent-handoff-protocol
  - **memory** (2): memory-garden-writer, context-compression-ritual
  - **meta** (2): skill-creator, skill-maintenance-ritual

### 3. New MCP Tools (`internal/dojo/handler.go`)

Added three new tools to the MCP server:

**`dojo.list_skills`**
- Lists all available skills grouped by category
- No parameters required
- Returns formatted list with descriptions

**`dojo.get_skill`**
- Retrieves the full content of a specific skill
- Parameters: `name` (string)
- Returns complete skill with category, description, and content

**`dojo.search_skills`**
- Searches skills by keyword across name, description, and content
- Parameters: `query` (string)
- Returns matching skills with metadata

### 4. Documentation Updates

Updated `README_V2.md`:
- Added skills to the AROMA sanctuary section
- Documented the 3 new skill tools with examples
- Updated tool count from 5 to 8 new tools in v2.0

## Architecture

### Skill Structure

Skills are stored as Go functions that return markdown content:

```go
type Skill struct {
    Name        string   // Unique identifier (kebab-case)
    Description string   // When and how to use this skill
    Category    string   // Grouping (learning, meta, strategy, etc.)
    Content     string   // Full markdown content
}
```

### Integration with Wisdom Base

Skills are integrated into the existing `wisdom.Base` architecture:

```
wisdom.Base
├── seeds      []Seed       (20 seed patches)
├── resources  []Resource   (8 documentation resources)
├── principles string       (core principles)
└── skills     []Skill      (20 skills) ← NEW
```

## Usage Examples

### Listing All Skills

```json
// Request
{
  "tool": "dojo.list_skills"
}

// Response
# Dojo Genesis Skills

## learning
- **agent-to-agent-teaching**: Protocol for teaching as a peer...
- **patient-learning-protocol**: Protocol for learning at the pace of understanding...

## meta
- **skill-creator**: Guide for creating effective skills...
- **skill-maintenance-ritual**: Systematic process for maintaining skills...

## strategy
- **strategic-scout**: Strategic exploration and decision-making framework...
...
```

### Getting a Specific Skill

```json
// Request
{
  "tool": "dojo.get_skill",
  "name": "agent-to-agent-teaching"
}

// Response
# agent-to-agent-teaching

**Category:** learning
**Description:** Protocol for teaching as a peer, not an expert...

---

# Agent-to-Agent Teaching Protocol
...
```

### Searching for Skills

```json
// Request
{
  "tool": "dojo.search_skills",
  "query": "learning"
}

// Response
# Skills Matching: learning

Found 2 skill(s):

## 1. agent-to-agent-teaching
**Category:** learning
**Description:** Protocol for teaching as a peer...

## 2. patient-learning-protocol
**Category:** learning
**Description:** Protocol for learning at the pace of understanding...
```

## Design Decisions

### 1. Why Embedded Content?

Skills are embedded in the Go binary rather than loaded from files because:
- **Simplicity**: No file I/O or path management needed
- **Portability**: Single binary contains all wisdom
- **Consistency**: Same pattern as seeds and resources
- **Performance**: No disk reads at runtime

For full skills with bundled resources, the content functions could be updated to load from the filesystem if needed.

### 2. Why Simple String Search?

The current implementation uses case-insensitive string matching rather than semantic search because:
- **Sufficient**: 8 skills with clear descriptions are easy to search
- **Fast**: No dependencies on embedding models or vector DBs
- **Consistent**: Matches the pattern used for seeds and resources
- **Extensible**: Can be upgraded to semantic search if needed

### 3. Why Skills in Wisdom Package?

Skills live in `internal/wisdom` alongside seeds and resources because:
- **Conceptual Cohesion**: All are knowledge artifacts
- **Shared Infrastructure**: All use the same Base and search patterns
- **Clean Separation**: Handler only orchestrates, wisdom contains content

## Testing

Build and test the server:

```bash
cd /Users/alfonsomorales/ZenflowProjects/dojo-mcp-server
go build -o dojo-mcp-server ./cmd/server
./dojo-mcp-server
```

The server compiles successfully with no errors.

## Future Enhancements

### Near-Term

1. **Add More Skills**: Expand beyond the initial 8 to include all skills from dojo-genesis
2. **Skill Prompts**: Make skills available as MCP prompts (like seeds)
3. **Skill Resources**: Support bundled resources (scripts, references, templates)

### Long-Term

1. **Dynamic Loading**: Load skills from filesystem for easier updates
2. **Semantic Search**: Upgrade to embedding-based semantic search
3. **Skill Versioning**: Track and manage skill versions
4. **Skill Dependencies**: Support skill-to-skill references
5. **Skill Analytics**: Track skill usage and effectiveness

## Benefits

### For Users

- **Actionable Workflows**: Not just philosophy, but concrete processes
- **Procedural Knowledge**: Step-by-step guidance for complex tasks
- **Searchable**: Find the right skill for any situation
- **Integrated**: Access skills alongside seeds and resources

### For the Ecosystem

- **Knowledge Sharing**: Skills from dojo-genesis now available to all MCP clients
- **Standard Format**: Reinforces skill-creator frontmatter standard
- **Extensible**: Easy to add more skills as the ecosystem grows
- **Consistent**: Same patterns as seeds and resources

## Conclusion

The skills integration extends the dojo-mcp-server from a philosophy and wisdom server to a complete practice toolkit. Users can now access not just the "why" (seeds and principles) but also the "how" (skills and workflows) of the Dojo ecosystem.

This creates a unified knowledge base where:
- **Seeds** provide architectural patterns
- **Resources** provide philosophical grounding
- **Skills** provide actionable workflows

Together, they form a complete system for building, learning, and practicing with AI.
