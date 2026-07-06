# Dojo Genesis MCP Server -- Status

**Version:** 3.2.0
**Module:** `github.com/DojoGenesis/mcp`
**Last Updated:** 2026-07-06

## Health

| Check | Status |
|-------|--------|
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go test -race ./...` | Pass (90+ tests) |
| CI pipeline | Configured (.github/workflows/ci.yml) |
| Release config | Configured (.goreleaser.yml) |

## Inventory

| Category | Count | Details |
|----------|-------|---------|
| Tools | 28 | 7 methodology + 3 gateway-memory + 3 Memory Hub (Postgres) + dojo_fetch + dojo_dispatch + 3 seeds + 3 agents + 3 project + 2 disposition + converge/health |
| Seed Patches | 20 | 10 Dojo Genesis + 10 AROMA/Serenity Valley |
| Skills | 43 bundled / 683 with CoworkPlugins | 683 community + 80 first-party; loaded from SKILL.md files at startup |
| Resources | 8 + skills | Philosophy, principles, norms, design, synthesis, protocol, modes, planning + all loaded skills |

## Dependencies

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.23+ | Required |
| `mcp-go` | v0.47.0 | MCP protocol library (stdio + streamable HTTP) |
| `gopkg.in/yaml.v3` | latest | YAML frontmatter parsing |
| `jackc/pgx/v5` | v5.10.0 | Memory Hub (read-only Postgres) |

## Configuration

See README Configuration — full env table (skills, ADR, gateway, memory hub,
HTTP mode: `DOJO_HTTP_ADDR`, `DOJO_MCP_API_KEYS`,
`DOJO_DISPATCH_ALLOWED_LABELS`, `DOJO_DISPATCH_RATE_PER_MIN`).

## Architecture

- **Transport:** stdio (default) or streamable HTTP behind Bearer-key auth (`DOJO_HTTP_ADDR`, OCH-132 Lane B)
- **Binary:** Single self-contained executable
- **Skills:** Loaded from filesystem at startup, with bundled fallback
- **Write capability:** `dojo.log_decision` writes ADR markdown files to disk
- **Search:** Keyword + trigger-based relevance scoring

## Bundled Plugins

| Plugin | Skills | Notes |
|--------|--------|-------|
| agent-orchestration | 5 | Dispatch, handoff, workspace navigation, teaching, decision propagation |
| continuous-learning | 5 | Research modes, project exploration, retrospective, research synthesis, debugging |
| skill-forge | 5 | Skill creation, maintenance, process extraction, file management, community scan |
| specification-driven-development | 5 | Spec writer, release spec, implementation prompt, frontend-from-backend, parallel tracks |
| strategic-thinking | 5 | Strategic scout, multi-surface strategy, product positioning, iterative scouting |
| system-health | 5 | Health audit, status writing, documentation audit, repo context sync, hooks reference |
| wisdom-garden | 5 | Memory garden, seed extraction, seed library, compression ritual, seed-to-skill converter |
| dojo-craft | 8 | adr-writer, claude-md-guardian, codebase-viewer, convergence-checker, memory-curator, project-scaffolder, scout-writer, seed-curator |

DojoCraft integration added 2026-04-15: 8 skills covering ADR authoring, CLAUDE.md auditing, codebase intelligence, convergence gating, memory management, project scaffolding, strategic scouting, and seed lifecycle.

## What Changed in v3.0.0

- 14 tools -> 7 tools (fewer, each does something real)
- Hardcoded skills -> file-backed from CoworkPlugins (671 SKILL.md files: 599 community + 72 first-party)
- Static templates -> structured methodology scaffolds
- Read-only -> one write tool (`log_decision`)
- No configuration -> `DOJO_SKILLS_PATH` and `DOJO_ADR_PATH` env vars
