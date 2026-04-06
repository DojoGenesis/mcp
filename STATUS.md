# Dojo Genesis MCP Server -- Status

**Version:** 2.1.0
**Module:** `github.com/DojoGenesis/mcp-server`
**Last Updated:** 2026-04-09

## Health

| Check | Status |
|-------|--------|
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go test ./...` | Pass |
| CI pipeline | Configured (.github/workflows/ci.yml) |
| Release config | Configured (.goreleaser.yml) |

## Inventory

| Category | Count | Details |
|----------|-------|---------|
| Tools | 14 | 6 core + 5 AROMA/Serenity Valley + 3 skill tools |
| Seed Patches | 20 | 10 Dojo Genesis + 10 AROMA/Serenity Valley |
| Skills | 32 | 8 categories (learning, strategy, process, workflow, meta, reflection, memory, development) |
| Resources | 8 | Philosophy, principles, norms, design patterns, synthesis, protocol, modes, planning |
| Prompts | 20 | One per seed patch |

## Dependencies

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.23+ | Required |
| `mcp-go` | v0.47.0 | Latest as of 2026-04-09 |
| `google/uuid` | v1.6.0 | Indirect (via mcp-go) |

## Architecture

- **Transport:** stdio (standard MCP transport)
- **Binary:** Single self-contained executable, no config files needed
- **All content embedded:** Seeds, skills, resources, and principles compiled into the binary
- **Search:** Keyword-based relevance scoring across all wisdom content

## Known Limitations

- Search uses keyword matching, not semantic/vector search
- Skills content is summary-length (full content would be loaded from external files in a future version)
- No hot-reload of content (requires rebuild)
