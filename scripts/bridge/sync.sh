#!/usr/bin/env bash
# /opt/dojo/memory-mirror/sync.sh - pull TresPies-AI-Orchestration + mirror markdown -> Postgres Memory Hub.
# Self-contained ongoing sync (cron). Read-only clone via deploy key; idempotent upsert.
set -euo pipefail
DIR=/opt/dojo/memory-mirror
REPO="$DIR/repo"
KEY=/root/.ssh/memory_repo_deploy
URL="git@github.com:TresPies-source/TresPies-AI-Orchestration.git"
export GIT_SSH_COMMAND="ssh -i $KEY -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new"
ts(){ date -u +%FT%TZ; }

if [ ! -d "$REPO/.git" ]; then
  echo "$(ts) clone $URL"; git clone --depth 1 "$URL" "$REPO"
else
  echo "$(ts) fetch+reset"; git -C "$REPO" fetch --depth 1 origin main && git -C "$REPO" reset --hard origin/main
fi
SHA="$(git -C "$REPO" rev-parse --short HEAD)"
python3 "$DIR/dojo_memory_mirror.py" "$REPO" "$DIR/memories.csv"
docker cp "$DIR/memories.csv" dojo-memory-db-1:/tmp/memories.csv >/dev/null
docker exec -i dojo-memory-db-1 psql -v ON_ERROR_STOP=1 -U dojo_memory -d dojo_memory <<SQL
CREATE TEMP TABLE stg (LIKE memories INCLUDING DEFAULTS);
\copy stg (id,name,slug,type,description,body,updated,origin_session,links) FROM '/tmp/memories.csv' WITH (FORMAT csv, HEADER true)
INSERT INTO memories AS m SELECT id,name,slug,type,description,body,updated,origin_session,links FROM stg
ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name,slug=EXCLUDED.slug,type=EXCLUDED.type,description=EXCLUDED.description,body=EXCLUDED.body,updated=EXCLUDED.updated,origin_session=EXCLUDED.origin_session,links=EXCLUDED.links;
INSERT INTO mirror_runs (files_seen,rows_upsert,source_sha) VALUES ((SELECT count(*) FROM stg),(SELECT count(*) FROM stg),'$SHA');
SQL
echo "$(ts) mirrored sha=$SHA rows=$(docker exec dojo-memory-db-1 psql -U dojo_memory -d dojo_memory -tAc 'SELECT count(*) FROM memories;')"
