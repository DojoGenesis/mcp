# scripts/bridge — Memory Hub mirror pipeline (repo-tracked copies)

Provenance: copied verbatim from dojo-bridge `/opt/dojo/memory-mirror/` on
2026-07-06 (OCH-131 follow-up, folded into OCH-132 per spec §6.5). The
canonical *running* copies live on the bridge; this directory exists so that
**schema changes are reviewed against their consumers** — `internal/memhub`
in this repo issues SELECTs against the `memories` table these scripts
maintain.

| File | Role on the bridge |
|---|---|
| `sync.sh` | Half-hourly cron: shallow-pull `TresPies-AI-Orchestration` (read-only deploy key), regenerate `memories.csv`, `\copy` into a staging table, upsert into `memories`, record a `mirror_runs` row. Column list in `\copy` is explicit — the generated `search_vector` column broke `LIKE`-clone staging once already (OCH-131). |
| `dojo_memory_mirror.py` | Markdown → CSV extraction (frontmatter + body → id, name, slug, type, description, body, updated, origin_session, links). |
| `dojo_memory_search.py` | Standalone CLI search against the hub (pre-dates the MCP tools; kept for shell use). |
| `migrations/0001_memory_fts.sql` | Generated `search_vector` tsvector (description+body) + GIN index. |
| `migrations/0002_dojo_mcp_ro_role.sql` | SELECT-only `dojo_mcp_ro` role for the public MCP endpoint (Lane B). The endpoint never writes the hub. |

## Schema contract with internal/memhub

`memories(id, name, slug, type, description, body, updated, origin_session,
links, search_vector)` — if you change this shape, update BOTH `sync.sh` +
`dojo_memory_mirror.py` AND the SQL in `internal/memhub/client.go`, then run
the guarded live test through the bridge tunnel:

```
DOJO_TEST_MEMORY_DB_URL="host=localhost port=5433 user=dojo_memory password=… dbname=dojo_memory sslmode=disable" \
  go test ./internal/memhub/ -run TestLiveHub -v
```

Notes `internal/memhub` relies on: `search_vector` covers description+body
only (client adds name/slug ILIKE legs); `updated` is read as `::text`;
`id` compared as `::text`.

## Deploying script changes

Edit here → review → copy to the bridge (`scp` to
`/opt/dojo/memory-mirror/`) → keep a dated `.bak` of the replaced file.
The bridge's cron and deploy key are untouched by repo changes.
