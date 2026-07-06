-- 0001_memory_fts.sql — add full-text search column + GIN index to memories table
--
-- DEPLOYMENT (integrator):
--   Apply once on bridge after the memories table already exists:
--
--     ssh dojo-bridge \
--       docker exec -i dojo-memory-db-1 \
--         psql -U dojo_memory -d dojo_memory \
--           < /opt/dojo/memory-mirror/migrations/0001_memory_fts.sql
--
--   Idempotent: the column and index are created with IF NOT EXISTS guards, so
--   re-running is safe.  The column is STORED so Postgres maintains it
--   automatically on INSERT/UPDATE — no application changes needed.
--
-- USAGE (after applying):
--   Replace the ILIKE WHERE clause in dojo_memory_search.py with:
--
--     WHERE search_vector @@ plainto_tsquery('english', :query)
--
--   and ORDER BY can be enhanced with ts_rank:
--
--     ORDER BY ts_rank(search_vector, plainto_tsquery('english', :query)) DESC,
--              updated DESC NULLS LAST
--
-- PERFORMANCE:
--   GIN index on a tsvector column reduces full-text search latency from O(n)
--   full-scan (ILIKE) to O(log n + k) — typically single-digit milliseconds even
--   at tens of thousands of rows.

-- 1. Add the generated tsvector column (Postgres 12+).
--    coalesce() normalises NULL description/body to '' so the concatenation
--    never produces NULL and every row always participates in the index.
ALTER TABLE memories
    ADD COLUMN IF NOT EXISTS search_vector tsvector
        GENERATED ALWAYS AS (
            to_tsvector(
                'english',
                coalesce(description, '') || ' ' || coalesce(body, '')
            )
        ) STORED;

-- 2. GIN index on the generated column.
--    The index name includes the migration number so it's unambiguous in
--    pg_indexes and won't collide with future indexes on the same table.
CREATE INDEX IF NOT EXISTS memories_search_vector_gin_0001
    ON memories
    USING GIN (search_vector);
