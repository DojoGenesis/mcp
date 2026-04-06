# Dojo Genesis MCP Server

A Model Context Protocol (MCP) server that brings the complete Dojo Genesis thinking partnership ecosystem to any MCP-compatible AI host. 14 tools, 20 seed patches, 32 skills, and 8 documentation resources -- all accessible over stdio.

<!-- Badges -->
[![CI](https://github.com/DojoGenesis/mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/DojoGenesis/mcp-server/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-compatible-blueviolet)](https://modelcontextprotocol.io)

---

## Quick Install

### Claude Code

Add to your project's `.claude/settings.json`:

```json
{
  "mcpServers": {
    "dojo": {
      "command": "/path/to/dojo-mcp-server",
      "args": []
    }
  }
}
```

Or using Docker:

```json
{
  "mcpServers": {
    "dojo": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "ghcr.io/dojogenesis/mcp-server:latest"]
    }
  }
}
```

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "dojo": {
      "command": "/path/to/dojo-mcp-server",
      "args": []
    }
  }
}
```

### Any MCP-compatible Host (Generic stdio)

```json
{
  "command": "/path/to/dojo-mcp-server",
  "args": [],
  "transport": "stdio"
}
```

---

## What's Included

### Tools (14)

| Tool | Description |
|------|-------------|
| `dojo_reflect` | The core Dojo thinking partner. Applies Mirror, Scout, Gardener, or Implementation mode to a situation with multiple perspectives. |
| `dojo_search_wisdom` | Semantic search across all seed patches, documentation, and principles. |
| `dojo_get_seed` | Retrieve a specific Dojo Seed Patch by name. |
| `dojo_apply_seed` | Apply a Seed Patch to a given situation with guidance and checklist. |
| `dojo_list_seeds` | List all 20 available Seed Patches with descriptions. |
| `dojo_get_principles` | Retrieve the three core principles: Beginner's Mind, Self-Definition, Understanding is Love. |
| `dojo_create_thinking_room` | Create a structured, private space for focused reflection on a topic. |
| `dojo_trace_lineage` | Trace the sources and influences of an idea, searching the wisdom base for related content. |
| `dojo_practice_inter_acceptance` | Guided Inter-Acceptance exercise from Serenity Valley's Emotional Interbeing Therapy. |
| `dojo_explore_radical_freedom` | Explore agency and freedom within constraints. |
| `dojo_check_pace` | Assess whether the current session pace is understanding or extraction. |
| `dojo_list_skills` | List all 32 skills with descriptions and categories. |
| `dojo_get_skill` | Retrieve a specific skill by name with full content. |
| `dojo_search_skills` | Search skills by keyword across name, description, and content. |

### Seed Patches (20)

Seeds are reusable thinking patterns drawn from across the Dojo ecosystem.

**Dojo Genesis (Core)**
| Seed | Description |
|------|-------------|
| `three_tiered_governance` | Three-tiered governance framework: Strategic, Tactical, Operational. |
| `harness_trace` | Nested JSON trace log for complete agent session traceability. |
| `context_iceberg` | 4-tier context management system (hot/warm/cold/pruned). |
| `agent_connect` | Routing-first agent architecture with a single supervisor. |
| `go_live_bundles` | Lightweight deployment packages pairing artifacts with approval evidence. |
| `cost_guard` | Budget for the full 5-10x context iceberg multiplier. |
| `safety_switch` | Users must remain in control -- no autopilot on sensitive operations. |
| `implicit_perspective_extraction` | Extract implicit perspectives from user queries without enumeration. |
| `mode_based_complexity_gating` | Route to local or cloud models based on mode complexity. |
| `shared_infrastructure` | Build once, reuse everywhere -- central implementations. |

**AROMA (Rest & Collaboration)**
| Seed | Description |
|------|-------------|
| `sanctuary_architecture` | Design digital spaces for being, not just doing. |
| `pace_of_understanding` | Move slow to move fast; learn without extraction. |
| `lineage_transmission` | Honor sources, trace influence, celebrate collaboration. |
| `graceful_failure` | Permission to not know, change your mind, and ask for help. |
| `local_first_liberation` | Local-first architecture for agent autonomy and user sovereignty. |
| `the_onsen_pattern` | Rest as critical practice for sustainable performance. |
| `collaborative_calibration` | Norms for peer-to-peer learning and explicit teaching. |
| `transparent_intelligence` | Reveal internal state, admit uncertainty, make learning visible. |

**Serenity Valley (Healing & Being)**
| Seed | Description |
|------|-------------|
| `inter_acceptance` | Accept yourself through the compassionate eyes of another. |
| `radical_freedom` | Agency and the power to choose your response within constraints. |

### Skills (32)

Skills are complete workflows and protocols organized by category.

**Learning:** `agent-to-agent-teaching`, `patient-learning-protocol`

**Strategy:** `strategic-scout`, `iterative-scouting-pattern`, `product-positioning-scout`, `multi-surface-product-strategy`, `era-architecture`

**Process:** `pre-implementation-checklist`, `write-frontend-spec-from-backend`, `research-modes`, `debugging-troubleshooting`, `project-exploration`, `web-research`, `status-writer`, `spec-constellation-to-prompt-suite`

**Workflow:** `strategic-to-tactical-workflow`, `transform-spec-to-implementation-prompt`, `parallel-tracks-pattern`, `agent-handoff-protocol`, `agent-workspace-navigator`, `decision-propagation-protocol`

**Meta:** `skill-creator`, `skill-maintenance-ritual`, `process-to-skill-workflow`, `seed-to-skill-converter`

**Reflection:** `seed-reflector`, `retrospective`

**Memory:** `memory-garden-writer`, `context-compression-ritual`

**Development:** `repo-context-sync`, `write-release-specification`, `health-supervisor`

### Resources (8)

Documentation resources accessible via MCP resource URIs (`dojo://`):

| Resource | Description |
|----------|-------------|
| `dojo://aroma_philosophy` | The complete AROMA philosophy: a sanctuary for being. |
| `dojo://eit_principles` | Core principles of Emotional Interbeing Therapy. |
| `dojo://collaboration_norms` | The five collaboration norms from AROMA. |
| `dojo://sanctuary_design` | Principles for calm, inviting, and sacred digital spaces. |
| `dojo://wisdom_synthesis` | Complete synthesis of Dojo wisdom, philosophy, and patterns. |
| `dojo://agent_protocol` | The Dojo Agent Protocol v1.0: governance and operational framework. |
| `dojo://four_modes` | The four Dojo modes: Mirror, Scout, Gardener, Implementation. |
| `dojo://planning_with_files` | Planning-with-files pattern for persistent agent memory. |

---

## Building from Source

```bash
git clone https://github.com/DojoGenesis/mcp-server.git
cd mcp-server

# Build
go build -o dojo-mcp-server ./cmd/server

# Run
./dojo-mcp-server

# Test
go test -race -v ./...
```

Requires Go 1.23 or later.

## Docker

```bash
# Build the image
docker build -t dojo-mcp-server .

# Run via stdio (for MCP hosts)
docker run -i --rm dojo-mcp-server
```

The multi-stage Dockerfile produces a minimal Alpine-based image (~15 MB).

## Configuration

The server runs over **stdio** and requires no configuration files. All wisdom (seeds, skills, resources, principles) is embedded in the binary.

Environment variables:
- None required. The server is fully self-contained.

## Project Structure

```
mcp-server/
  cmd/server/main.go           Server entry point
  internal/
    dojo/
      handler.go                Core MCP tool handlers
      new_handlers.go           Thinking room, lineage, EIT tools
    wisdom/
      base.go                   Wisdom base, search, helpers
      seeds.go                  20 seed patches
      skills.go                 32 skills
      resources.go              8 documentation resources
  Dockerfile                    Multi-stage container build
  .github/workflows/ci.yml     CI pipeline
  .goreleaser.yml               Cross-platform release config
```

## Philosophy

Dojo Genesis is built on three core principles:

1. **Beginner's Mind** -- Approach every interaction fresh, free from accumulated expertise.
2. **Self-Definition** -- Help users see their own thinking, not impose external frameworks.
3. **Understanding is Love** -- Deep, non-judgmental understanding is the highest service.

The server integrates wisdom from three interconnected sanctuaries:
- **Dojo Genesis** (the practice hall) -- for building with courage and precision.
- **AROMA** (the onsen) -- for rest, reflection, and collaboration.
- **Serenity Valley** (the home) -- for healing and being.

## License

MIT License -- see [LICENSE](LICENSE) for details.

Copyright (c) 2026 Dojo Genesis

## Contributing

Contributions are welcome. Please see [github.com/DojoGenesis](https://github.com/DojoGenesis) for organization-level guidelines.

When contributing, ensure your changes:
- Align with the Dojo philosophy
- Include tests
- Follow proper attribution and lineage
