#!/usr/bin/env python3
"""dojo_memory_search.py — query the Memory Hub Postgres instance via docker exec.

Runs ON bridge (dojo-bridge); shells out to docker exec so no network credentials
are needed — the container trusts the local Unix socket for user `dojo_memory`.

DEPLOYMENT (integrator):
  1. Copy this file into /opt/dojo/memory-mirror/
     scp dojo_memory_search.py dojo-bridge:/opt/dojo/memory-mirror/
  2. Make it executable:
     ssh dojo-bridge chmod +x /opt/dojo/memory-mirror/dojo_memory_search.py
  3. Invoke:
     ssh dojo-bridge python3 /opt/dojo/memory-mirror/dojo_memory_search.py "<query>"
     ssh dojo-bridge python3 /opt/dojo/memory-mirror/dojo_memory_search.py "nightshift" --type seed --limit 5

INDEXING NOTE (v1 vs future):
  This version uses ILIKE '%<term>%' over (description || ' ' || body) — correct
  but O(n) full-scan. When the table grows past ~50 K rows, apply the migration
  migrations/0001_memory_fts.sql to add a GIN-backed tsvector column; then swap
  the WHERE clause to:
      to_tsvector('english', coalesce(description,'') || ' ' || coalesce(body,''))
        @@ plainto_tsquery('english', $1)
  which reduces latency to single-digit milliseconds at any realistic corpus size.

STDLIB ONLY — no pip dependencies.
"""
from __future__ import annotations

import argparse
import subprocess
import sys
import textwrap

CONTAINER = "dojo-memory-db-1"
DB_USER   = "dojo_memory"
DB_NAME   = "dojo_memory"

# Maximum characters from description shown per result line
DESC_WIDTH = 80


def _pg_literal(value: str) -> str:
    """Wrap a value as a safe single-quoted SQL string literal (double any embedded
    single quotes). Injection-safe for string values under standard_conforming_strings
    (the default) — we inline the term directly because psql's `:var` interpolation
    does NOT fire inside a `-c` command string sent over `docker exec`."""
    return "'" + value.replace("'", "''") + "'"


def build_sql(query: str, type_filter: str | None, limit: int) -> str:
    """Return the SELECT with the search term + type safely inlined as escaped SQL
    string literals (see _pg_literal). `limit` is an int from argparse."""
    q = _pg_literal(query)
    where_clauses = [
        f"(description ILIKE '%' || {q} || '%' OR body ILIKE '%' || {q} || '%')"
    ]
    if type_filter:
        where_clauses.append(f"type = {_pg_literal(type_filter)}")

    where = " AND ".join(where_clauses)

    sql = textwrap.dedent(f"""\
        SELECT
            id,
            type,
            name,
            left(coalesce(description, ''), {DESC_WIDTH}) AS desc_snippet
        FROM memories
        WHERE {where}
        ORDER BY updated DESC NULLS LAST
        LIMIT {int(limit)};
    """)
    return sql


def run_query(query: str, type_filter: str | None, limit: int) -> int:
    """Shell out to docker exec psql and print results. Returns exit code."""
    sql = build_sql(query, type_filter, limit)

    cmd = [
        "docker", "exec", CONTAINER,   # no -i: query comes via -c, and inheriting
        "psql",                        # an interactive stdin lets psql consume the
                                       # caller's stdin (e.g. a heredoc driving it)
        "-U", DB_USER,
        "-d", DB_NAME,
        "--no-password",
        "-t",               # tuples only (no column headers or row counts)
        "-A",               # unaligned output — one row per line, fields split by |
        "-F", " | ",        # field separator
        "-c", sql,
    ]

    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            stdin=subprocess.DEVNULL,
        )
    except FileNotFoundError:
        print("ERROR: docker not found on PATH. Run this script on dojo-bridge.", file=sys.stderr)
        return 2

    if result.returncode != 0:
        print(f"ERROR: psql exited {result.returncode}", file=sys.stderr)
        if result.stderr:
            print(result.stderr.strip(), file=sys.stderr)
        return result.returncode

    output = result.stdout.strip()
    if not output:
        print("(no results)")
        return 0

    # Print header then rows
    print(f"{'id':<36}  {'type':<16}  {'name':<40}  description")
    print("-" * 120)
    for line in output.splitlines():
        line = line.strip()
        if not line:
            continue
        # Rows arrive as: id | type | name | desc_snippet
        # Split on first 3 occurrences of " | " to preserve " | " in the snippet.
        parts = line.split(" | ", 3)
        if len(parts) == 4:
            row_id, row_type, row_name, row_desc = parts
            print(f"{row_id:<36}  {row_type:<16}  {row_name:<40}  {row_desc}")
        else:
            # Unexpected format — print raw
            print(line)

    return 0


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Search the Memory Hub (memories table) on dojo-bridge.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=textwrap.dedent("""\
            examples:
              dojo_memory_search.py "nightshift"
              dojo_memory_search.py "kata acquisition" --type seed --limit 5
              dojo_memory_search.py "gateway port" --limit 20
        """),
    )
    parser.add_argument(
        "query",
        help="Search term (case-insensitive; matched against description and body).",
    )
    parser.add_argument(
        "--type",
        dest="type_filter",
        metavar="TYPE",
        default=None,
        help="Filter by memory type (e.g. seed, project, feedback, reference).",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=10,
        metavar="N",
        help="Maximum number of results to return (default: 10).",
    )
    args = parser.parse_args()

    if args.limit < 1:
        parser.error("--limit must be a positive integer")

    sys.exit(run_query(args.query, args.type_filter, args.limit))


if __name__ == "__main__":
    main()
