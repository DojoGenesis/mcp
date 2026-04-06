# Changelog

All notable changes to this project will be documented in this file.

## [3.0.0] - 2026-04-05

### Changed
- **Breaking:** 14 tools reduced to 7 focused methodology tools
- Handler now accepts `skillsPath` and `adrPath` parameters
- Skills loaded from filesystem (CoworkPlugins SKILL.md files) instead of hardcoded Go
- Resources URI scheme updated (`dojo://resources/`, `dojo://seeds/`, `dojo://skills/`)
- Version bumped to 3.0.0

### Added
- `dojo.scout` -- 4-step strategic analysis scaffold
- `dojo.invoke_skill` -- load full skill workflow by name
- `dojo.log_decision` -- write Architecture Decision Records to disk (anchor demo)
- `internal/skills/` package -- filesystem SKILL.md loader with YAML frontmatter parsing
- `internal/decisions/` package -- ADR file writer
- `internal/dojo/scaffolds.go` -- methodology templates for scout and reflect
- 15 bundled SKILL.md files via `go:embed` (works without `$DOJO_SKILLS_PATH`)
- `gopkg.in/yaml.v3` dependency for frontmatter parsing
- Environment variables: `DOJO_SKILLS_PATH`, `DOJO_ADR_PATH`
- MCP resources for all loaded skills (`dojo://skills/{plugin}/{name}`)
- Keyword + trigger-based skill search with relevance ranking
- 90+ tests across all packages

### Removed
- `dojo.get_seed` (use `dojo.apply_seed`)
- `dojo.list_seeds` (seeds accessible via `dojo.apply_seed` and `dojo.reflect`)
- `dojo.get_principles` (available as MCP resource)
- `dojo.create_thinking_room` (was static template)
- `dojo.trace_lineage` (was static template)
- `dojo.practice_inter_acceptance` (was static template)
- `dojo.explore_radical_freedom` (was static template)
- `dojo.check_pace` (folded into `dojo.reflect`)
- `dojo.search_wisdom` (replaced by `dojo.search_skills` + `dojo.reflect`)
- `dojo.get_skill` (use `dojo.invoke_skill`)
- `internal/dojo/new_handlers.go` (consolidated into handler.go)
- `internal/wisdom/skills.go` (replaced by `internal/skills/` package)
- MCP prompts (seeds now accessed via `dojo.apply_seed` tool)
- `RegisterPrompts` method

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
