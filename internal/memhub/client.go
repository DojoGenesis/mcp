// Package memhub is a read-only client for the Postgres Memory Hub on
// dojo-bridge: the queryable mirror of the TresPies/DojoGenesis markdown
// institutional memory (table `memories`, FTS via the generated
// `search_vector` column — see scripts/bridge/ for the mirror pipeline).
//
// The public endpoint connects with a SELECT-only role (`dojo_mcp_ro`);
// nothing in this package writes.
package memhub

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Entry is one memory row. Body is populated only by Get.
type Entry struct {
	Slug        string
	Name        string
	Type        string
	Description string
	Snippet     string
	Body        string
	Updated     string
	Rank        float32
}

// ErrNotFound is returned by Get when no row matches.
var ErrNotFound = errors.New("memory not found")

// Client is a pooled read-only hub connection.
type Client struct {
	pool *pgxpool.Pool
}

// New builds a lazy connection pool for dbURL. The pool dials on first use,
// so a temporarily unreachable hub does not prevent server startup; Ping
// reports actual reachability.
func New(ctx context.Context, dbURL string) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parse memory db url: %w", err)
	}
	cfg.MaxConns = 4
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("build memory db pool: %w", err)
	}
	return &Client{pool: pool}, nil
}

// Ping verifies connectivity.
func (c *Client) Ping(ctx context.Context) error { return c.pool.Ping(ctx) }

// Close releases the pool.
func (c *Client) Close() { c.pool.Close() }

const searchSQL = `
SELECT slug, name, type, coalesce(description, ''),
       coalesce(ts_headline('english', coalesce(body, ''),
                websearch_to_tsquery('english', $1),
                'StartSel=**, StopSel=**, MaxWords=35, MinWords=10, MaxFragments=2, FragmentDelimiter= … '), '') AS snippet,
       coalesce(ts_rank(search_vector, websearch_to_tsquery('english', $1)), 0) AS rank,
       coalesce(updated::text, '')
FROM memories
WHERE (search_vector @@ websearch_to_tsquery('english', $1)
       OR name ILIKE '%' || $2 || '%'
       OR slug ILIKE '%' || $2 || '%')
  AND ($3 = '' OR type = $3)
ORDER BY rank DESC, updated DESC NULLS LAST
LIMIT $4`

// SearchMemories runs ranked full-text search (websearch syntax: quoted
// phrases, OR, -exclusions) with name/slug substring fallback, optionally
// filtered by type.
func (c *Client) SearchMemories(ctx context.Context, query, typ string, limit int) ([]Entry, error) {
	limit = clampLimit(limit, 8, 25)
	rows, err := c.pool.Query(ctx, searchSQL, query, escapeLike(query), typ, limit)
	if err != nil {
		return nil, fmt.Errorf("hub search: %w", err)
	}
	defer rows.Close()

	var out []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Slug, &e.Name, &e.Type, &e.Description, &e.Snippet, &e.Rank, &e.Updated); err != nil {
			return nil, fmt.Errorf("hub search scan: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

const getSQL = `
SELECT slug, name, type, coalesce(description, ''), coalesce(body, ''), coalesce(updated::text, '')
FROM memories
WHERE slug = $1 OR id::text = $1
ORDER BY updated DESC NULLS LAST
LIMIT 1`

// GetMemory returns the full body of one memory by slug (or raw id).
func (c *Client) GetMemory(ctx context.Context, slug string) (*Entry, error) {
	var e Entry
	err := c.pool.QueryRow(ctx, getSQL, slug).
		Scan(&e.Slug, &e.Name, &e.Type, &e.Description, &e.Body, &e.Updated)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("hub get: %w", err)
	}
	return &e, nil
}

const recentSQL = `
SELECT slug, name, type, coalesce(description, ''), left(coalesce(body, ''), 240),
       coalesce(updated::text, '')
FROM memories
WHERE ($1 = '' OR type = $1)
ORDER BY updated DESC NULLS LAST
LIMIT $2`

// RecentMemories lists the most recently updated memories, optionally
// filtered by type — orientation on connect.
func (c *Client) RecentMemories(ctx context.Context, typ string, limit int) ([]Entry, error) {
	limit = clampLimit(limit, 10, 50)
	rows, err := c.pool.Query(ctx, recentSQL, typ, limit)
	if err != nil {
		return nil, fmt.Errorf("hub recent: %w", err)
	}
	defer rows.Close()

	var out []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Slug, &e.Name, &e.Type, &e.Description, &e.Snippet, &e.Updated); err != nil {
			return nil, fmt.Errorf("hub recent scan: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// clampLimit applies a default and an upper bound.
func clampLimit(n, def, max int) int {
	if n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

// escapeLike neutralizes LIKE metacharacters in user input used for the
// name/slug substring legs.
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}
