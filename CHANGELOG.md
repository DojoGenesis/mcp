# Changelog

All notable changes to this project will be documented in this file.

## [2.1.0] - 2026-04-09

### Changed
- Migrated to DojoGenesis org (`github.com/DojoGenesis/mcp-server`)
- Upgraded `mcp-go` from v0.8.0 to v0.47.0
- Fixed all `json.MarshalIndent` error handling (previously ignored errors)
- Updated server name to `dojo-mcp-server` and version to `2.1.0`

### Added
- Comprehensive test suite covering all 14 tool handlers
- CI pipeline (GitHub Actions: vet, build, test, lint)
- GoReleaser config for multi-arch releases (darwin/linux, amd64/arm64)
- Product-quality README.md with install snippets, full inventory, and build instructions
- CHANGELOG.md

### Removed
- Stale `dojo-mcp-server-v2-fixed.tar.gz` archive
- Internal-only `COMPLETION_REPORT.md`

## [2.0.0] - 2026-02-09

### Added
- Thinking room creation (`dojo.create_thinking_room`)
- Lineage tracing (`dojo.trace_lineage`)
- Inter-Acceptance practice (`dojo.practice_inter_acceptance`)
- Radical Freedom exploration (`dojo.explore_radical_freedom`)
- Pace check assessment (`dojo.check_pace`)
- 10 new AROMA and Serenity Valley seed patches (11-20)
- 32 skills across 8 categories
- 4 new documentation resources (AROMA philosophy, EIT principles, collaboration norms, sanctuary design)
- Skill tools: `dojo.list_skills`, `dojo.get_skill`, `dojo.search_skills`

## [1.0.0] - 2026-02-06

### Added
- Initial release with 6 core tools: reflect, search_wisdom, get_seed, apply_seed, list_seeds, get_principles
- 10 Dojo Genesis seed patches
- 4 documentation resources (wisdom synthesis, agent protocol, four modes, planning with files)
- MCP prompts for all seed patches
- Docker support
