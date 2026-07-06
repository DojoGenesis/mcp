-- 0002_dojo_mcp_ro_role.sql — SELECT-only role for the public MCP endpoint
--
-- The Lane B public endpoint (mcp.trespiesdesign.com, OCH-132) reads the
-- Memory Hub through internal/memhub. It must NEVER hold write on the hub:
-- the mirror pipeline (sync.sh) is the only writer.
--
-- DEPLOYMENT (operator or staging session):
--
--   ssh dojo-bridge \
--     docker exec -i dojo-memory-db-1 \
--       psql -U dojo_memory -d dojo_memory \
--         -v ro_password="'<GENERATED-PASSWORD>'" \
--         < /opt/dojo/memory-mirror/migrations/0002_dojo_mcp_ro_role.sql
--
--   Generate the password fresh (32+ random bytes); it lands in
--   /opt/dojo/.env as part of DOJO_MEMORY_DB_URL for the mcp-http service.
--   Never commit it.
--
-- Idempotent: re-running rotates the password (ALTER path). psql variables
-- do not expand inside DO $$ blocks, hence the format(...) + \gexec pattern.

SELECT format('CREATE ROLE dojo_mcp_ro LOGIN PASSWORD %L', :'ro_password')
WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'dojo_mcp_ro')
\gexec

SELECT format('ALTER ROLE dojo_mcp_ro WITH LOGIN PASSWORD %L', :'ro_password')
WHERE EXISTS (SELECT FROM pg_roles WHERE rolname = 'dojo_mcp_ro')
\gexec

GRANT CONNECT ON DATABASE dojo_memory TO dojo_mcp_ro;
GRANT USAGE ON SCHEMA public TO dojo_mcp_ro;
GRANT SELECT ON public.memories TO dojo_mcp_ro;

-- Deliberately NOT granted: mirror_runs (write bookkeeping), any sequences,
-- any future tables (no ALTER DEFAULT PRIVILEGES — new tables stay private
-- until explicitly granted here, reviewed against internal/memhub).
