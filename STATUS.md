# Dojo Genesis MCP Server -- Status

**Version:** 3.0.0
**Module:** `github.com/DojoGenesis/mcp-server`
**Last Updated:** 2026-04-05

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
| Tools | 7 | scout, invoke_skill, search_skills, apply_seed, log_decision, reflect, list_skills |
| Seed Patches | 20 | 10 Dojo Genesis + 10 AROMA/Serenity Valley |
| Skills | 15 bundled / 60 with CoworkPlugins | Loaded from SKILL.md files at startup |
| Resources | 8 + skills | Philosophy, principles, norms, design, synthesis, protocol, modes, planning + all loaded skills |

## Dependencies

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.23+ | Required |
| `mcp-go` | v0.47.0 | MCP protocol library |
| `gopkg.in/yaml.v3` | latest | YAML frontmatter parsing |

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| `DOJO_SKILLS_PATH` | (bundled) | Path to CoworkPlugins directory |
| `DOJO_ADR_PATH` | `./decisions` | ADR output directory |

## Architecture

- **Transport:** stdio (standard MCP transport)
- **Binary:** Single self-contained executable
- **Skills:** Loaded from filesystem at startup, with bundled fallback
- **Write capability:** `dojo.log_decision` writes ADR markdown files to disk
- **Search:** Keyword + trigger-based relevance scoring

## What Changed in v3.0.0

- 14 tools -> 7 tools (fewer, each does something real)
- Hardcoded skills -> file-backed from CoworkPlugins (60 SKILL.md files)
- Static templates -> structured methodology scaffolds
- Read-only -> one write tool (`log_decision`)
- No configuration -> `DOJO_SKILLS_PATH` and `DOJO_ADR_PATH` env vars
