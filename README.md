# Dojo Genesis MCP Server

The methodology layer for Claude Code -- the first MCP server that makes Claude measurably better at software development decisions by encoding 60 battle-tested thinking frameworks as active cognitive scaffolds.

Every other MCP server gives Claude more data. This one gives Claude better *methods*.

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
      "env": {
        "DOJO_SKILLS_PATH": "/path/to/CoworkPluginsByDojoGenesis",
        "DOJO_ADR_PATH": "./decisions"
      }
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
      "env": {
        "DOJO_SKILLS_PATH": "/path/to/CoworkPluginsByDojoGenesis",
        "DOJO_ADR_PATH": "./decisions"
      }
    }
  }
}
```

---

## What's Included

### Tools (7)

| Tool | Description |
|------|-------------|
| `dojo.scout` | 4-step strategic analysis scaffold. Frame tension, scout routes, synthesize, decide. |
| `dojo.invoke_skill` | Load a specific methodology skill by name. Returns the full workflow as actionable steps. |
| `dojo.search_skills` | Search the methodology library for skills matching a query. |
| `dojo.apply_seed` | Apply a reusable thinking pattern (seed) to a specific situation with checklist. |
| `dojo.log_decision` | Write an Architecture Decision Record (ADR) to disk. The only write-capable tool. |
| `dojo.reflect` | Structured reflection grounded in matched skills and seeds from the methodology library. |
| `dojo.list_skills` | List all available skills grouped by plugin category. |

### Skills (60 from CoworkPlugins, 15 bundled)

Skills are complete, battle-tested methodology workflows loaded from SKILL.md files at startup. When `DOJO_SKILLS_PATH` is set, all 60 CoworkPlugins skills are available. Without it, 15 key skills are bundled in the binary.

**Bundled skills:** `strategic-scout`, `release-specification`, `implementation-prompt`, `debugging`, `retrospective`, `pre-implementation-checklist`, `parallel-tracks`, `health-audit`, `seed-extraction`, `memory-garden`, `context-ingestion`, `research-modes`, `skill-creation`, `handoff-protocol`, `status-writing`

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

**AROMA & Serenity Valley**
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
| `inter_acceptance` | Accept yourself through the compassionate eyes of another. |
| `radical_freedom` | Agency and the power to choose your response within constraints. |

### Resources

Documentation resources accessible via MCP resource URIs:

- `dojo://resources/{name}` -- 8 documentation resources (AROMA philosophy, EIT principles, etc.)
- `dojo://seeds/{name}` -- 20 seed patches
- `dojo://skills/{plugin}/{name}` -- All loaded skills

---

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| `DOJO_SKILLS_PATH` | (bundled fallback) | Path to CoworkPlugins root directory containing `plugins/` |
| `DOJO_ADR_PATH` | `./decisions` | Directory where `dojo.log_decision` writes ADR files |

The server works out of the box with zero configuration (bundled skills, default ADR path).

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

# Run with skills path mounted
docker run -i --rm -v /path/to/CoworkPlugins:/skills -e DOJO_SKILLS_PATH=/skills dojo-mcp-server
```

## Project Structure

```
mcp-server/
  cmd/server/main.go           Server entry point, env vars
  internal/
    dojo/
      handler.go                7 MCP tool handlers
      scaffolds.go              Scout and reflect methodology templates
    skills/
      loader.go                 Filesystem SKILL.md loader
      search.go                 Keyword + trigger search
      bundled.go                go:embed fallback (15 key skills)
      bundled/                  Embedded SKILL.md files
    decisions/
      writer.go                 ADR file writer
    wisdom/
      base.go                   Wisdom base, search, helpers
      seeds.go                  20 seed patches
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

## License

MIT License -- see [LICENSE](LICENSE) for details.

Copyright (c) 2026 Dojo Genesis

## Contributing

Contributions are welcome. Please see [github.com/DojoGenesis](https://github.com/DojoGenesis) for organization-level guidelines.
