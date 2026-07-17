# Dojo Genesis MCP Server

The methodology layer for Claude Code -- the first MCP server that makes Claude measurably better at software development decisions by encoding 99 first-party working methods across 10 behavioral plugins as active cognitive scaffolds.

Every other MCP server gives Claude more data. This one gives Claude better *methods*.

<!-- Badges -->
[![CI](https://github.com/DojoGenesis/mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/DojoGenesis/mcp/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-compatible-blueviolet)](https://modelcontextprotocol.io)
[![LobeHub](https://lobehub.com/badge/mcp/dojogenesis-mcp-server)](https://lobehub.com/mcp/dojogenesis-mcp-server)

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

### Tools (28)

| Group | Tools |
|-------|-------|
| Methodology | `dojo_scout`, `dojo_invoke_skill`, `dojo_search_skills`, `dojo_list_skills`, `dojo_apply_seed`, `dojo_reflect` |
| Decisions | `dojo_log_decision` (writes ADR markdown to `DOJO_ADR_PATH`) |
| Gateway memory | `dojo_memory_list`, `dojo_memory_store`, `dojo_memory_search` (gateway session memory) |
| **Memory Hub** | `dojo_search_memory`, `dojo_get_memory`, `dojo_recent_memories` — read-only Postgres mirror of the institutional memory (`DOJO_MEMORY_DB_URL`) |
| **Unified fetch** | `dojo_fetch` — one tool to search AND fetch across memory hub + skills + ADRs + seeds (typed ids: `memory:slug`, `skill:name`, `adr:file.md`, `seed:name`) |
| Seeds | `dojo_seed_list`, `dojo_seed_create`, `dojo_seed_search` |
| Agents² | `dojo_agent_list`, `dojo_agent_dispatch`², `dojo_agent_chat`² |
| **Dispatch²** | `dojo_dispatch`² — prompt → LLM through the gateway |
| Project | `dojo_project_status`, `dojo_project_track`, `dojo_project_decision` |
| Disposition | `dojo_disposition_list`, `dojo_disposition_set` |
| Craft | `dojo_converge`, `dojo_health` |

² **Dispatch-class** — spends LLM provider budget through the gateway. In HTTP
mode these require a dispatch-enabled API key and are rate limited per key;
`dojo_scout`'s LLM path degrades to its offline scaffold for non-dispatch keys.

### Skills (99 from CoworkPlugins, 35 bundled)

Skills are complete, battle-tested methodology workflows loaded from SKILL.md files at startup. When `DOJO_SKILLS_PATH` is set, all 99 first-party CoworkPlugins skills are available. Without it, 35 key skills across 7 plugins are embedded in the binary via `go:embed`.

**Bundled plugins and skills:**

| Plugin | Skills |
|--------|--------|
| `agent-orchestration` | agent-teaching, decision-propagation, handoff-protocol, workflow-router, workspace-navigation |
| `continuous-learning` | debugging, project-exploration, research-modes, research-synthesis, retrospective, web-research-external |
| `skill-forge` | mcp-cloudflare-builder, mcp-server-builder, process-extraction, skill-audit, skill-creation, skill-maintenance |
| `specification-driven-development` | context-ingestion, frontend-from-backend, implementation-prompt, parallel-tracks, pre-implementation-checklist, release-specification |
| `strategic-thinking` | iterative-scouting, multi-surface-strategy, product-positioning, strategic-scout |
| `system-health` | documentation-audit, health-audit, semantic-clusters, status-writing |
| `wisdom-garden` | compression-ritual, memory-garden, seed-extraction, seed-library |

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
| `DOJO_ADR_PATH` | `./decisions` | Directory where `dojo_log_decision` writes ADR files |
| `DOJO_GATEWAY_URL` | `http://localhost:7340` | Dojo AgenticGateway base URL |
| `DOJO_GATEWAY_TOKEN` | (unset) | Bearer token for the gateway, if it requires one |
| `DOJO_MEMORY_DB_URL` | (unset → hub tools disabled) | Postgres DSN for the Memory Hub (URL or keyword form; use the SELECT-only role) |
| `DOJO_HTTP_ADDR` | (unset → stdio) | Opt into HTTP mode, e.g. `:8091` |
| `DOJO_MCP_API_KEYS` | (unset) | HTTP mode only: comma-separated `label:key` pairs (individually revocable) |
| `DOJO_DISPATCH_ALLOWED_LABELS` | (unset → none) | Key labels allowed to run dispatch-class tools |
| `DOJO_DISPATCH_RATE_PER_MIN` | `6` | Per-label rate limit for dispatch-class tools |

The server works out of the box with zero configuration (bundled skills, default ADR path, stdio).

---

## HTTP Mode (public endpoint)

Setting `DOJO_HTTP_ADDR` serves MCP streamable-HTTP instead of stdio:

- `POST/GET /mcp` — the MCP endpoint, Bearer-key required
  (`Authorization: Bearer <key>`). Keyless or wrong-key requests get 401.
- `/mcp/k/<key>` — same endpoint for clients that cannot send custom
  headers; the key is redacted in logs.
- `GET /health` — unauthenticated liveness (status + version only).

The server refuses to start in HTTP mode without a valid, non-empty
`DOJO_MCP_API_KEYS`. Every tool call is logged as
`tool_call tool=… key=<label> dur_ms=… outcome=…` — labels only, never key
material, never payloads.

```bash
DOJO_HTTP_ADDR=:8091 \
DOJO_MCP_API_KEYS="win:$(openssl rand -hex 32)" \
./dojo-mcp-server

curl -s localhost:8091/health
curl -s -X POST localhost:8091/mcp \
  -H "Authorization: Bearer <key>" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"curl","version":"0"}}}'
```

Client wiring (Claude Code): `claude mcp add --transport http dojo-remote
https://<host>/mcp --header "Authorization: Bearer <key>"`.

Deployment notes: ingress is expected to be a Cloudflare tunnel (no host
ports); pin the container image **by digest** so auto-pull never
surprise-deploys a new public surface. Image publishing:
`.github/workflows/docker-publish.yml` → `ghcr.io/dojogenesis/mcp`
(supersedes the old manually-pushed `mcpbydojogenesis` image).

---

## Building from Source

```bash
git clone https://github.com/DojoGenesis/mcp.git
cd mcp

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
