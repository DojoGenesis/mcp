package dojo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// fetchHit is one merged cross-store search result.
type fetchHit struct {
	store   string // memory | skill | adr | seed
	key     string // fetchable key within the store
	title   string
	detail  string // one-line context (type, plugin, date…)
	snippet string
	score   float64
}

// adrContentScanCap bounds how many ADR files get a head-content check per
// query (newest first); filename matching covers the rest.
const adrContentScanCap = 50

func (h *Handler) handleFetch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query  string   `json:"query"`
		ID     string   `json:"id"`
		Stores []string `json:"stores"`
		Limit  int      `json:"limit"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if strings.TrimSpace(args.ID) != "" {
		return h.fetchByID(ctx, strings.TrimSpace(args.ID))
	}
	if strings.TrimSpace(args.Query) == "" {
		return mcp.NewToolResultError("Provide 'query' to search or 'id' (store:key) to fetch."), nil
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 8
	}
	if limit > 25 {
		limit = 25
	}

	enabled := map[string]bool{"memory": true, "skill": true, "adr": true, "seed": true}
	if len(args.Stores) > 0 {
		enabled = map[string]bool{}
		for _, s := range args.Stores {
			enabled[strings.ToLower(strings.TrimSpace(s))] = true
		}
	}

	var hits []fetchHit
	var notes []string

	if enabled["memory"] {
		if h.hub == nil {
			notes = append(notes, "memory store skipped (hub not configured)")
		} else if mh, err := h.searchMemoryLeg(ctx, args.Query, limit); err != nil {
			notes = append(notes, fmt.Sprintf("memory store error: %v", err))
		} else {
			hits = append(hits, mh...)
		}
	}
	if enabled["skill"] {
		hits = append(hits, h.searchSkillLeg(args.Query, limit)...)
	}
	if enabled["adr"] {
		hits = append(hits, h.searchADRLeg(args.Query, limit)...)
	}
	if enabled["seed"] {
		hits = append(hits, h.searchSeedLeg(args.Query, limit)...)
	}

	if len(hits) == 0 {
		msg := fmt.Sprintf("Nothing in the dojo store matches: \"%s\"", args.Query)
		if len(notes) > 0 {
			msg += "\n\n(" + strings.Join(notes, "; ") + ")"
		}
		return mcp.NewToolResultText(msg), nil
	}

	sort.SliceStable(hits, func(i, j int) bool { return hits[i].score > hits[j].score })
	if len(hits) > limit {
		hits = hits[:limit]
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# dojo store: \"%s\"\n\n", args.Query)
	for i, ht := range hits {
		fmt.Fprintf(&sb, "## %d. [%s] %s — `%s:%s`\n", i+1, ht.store, ht.title, ht.store, ht.key)
		if ht.detail != "" {
			fmt.Fprintf(&sb, "%s\n", ht.detail)
		}
		if ht.snippet != "" {
			fmt.Fprintf(&sb, "> %s\n", strings.ReplaceAll(ht.snippet, "\n", " "))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("Fetch full content: `dojo_fetch` with id `store:key`.\n")
	if len(notes) > 0 {
		fmt.Fprintf(&sb, "(%s)\n", strings.Join(notes, "; "))
	}
	return mcp.NewToolResultText(sb.String()), nil
}

// ─── id mode ──────────────────────────────────────────────────────────────────

func (h *Handler) fetchByID(ctx context.Context, id string) (*mcp.CallToolResult, error) {
	store, key, ok := strings.Cut(id, ":")
	if !ok || strings.TrimSpace(key) == "" {
		return mcp.NewToolResultError("id must be 'store:key', e.g. 'memory:night-shift-family', 'skill:debugging', 'adr:2026-07-05_title.md', 'seed:name'"), nil
	}
	key = strings.TrimSpace(key)

	switch strings.ToLower(strings.TrimSpace(store)) {
	case "memory":
		if h.hub == nil {
			return mcp.NewToolResultError(hubNotConfigured), nil
		}
		e, err := h.hub.GetMemory(ctx, key)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("memory %q: %v", key, err)), nil
		}
		return mcp.NewToolResultText(renderHubEntry(e)), nil

	case "skill":
		skill, err := h.skillsLoader.GetByName(key)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("skill %q not found — search with dojo_fetch or dojo_search_skills", key)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("# Skill: %s\n**Plugin:** %s\n\n%s", skill.Name, skill.Plugin, skill.Content)), nil

	case "adr":
		base := filepath.Base(key)
		if base != key || !strings.HasSuffix(base, ".md") || strings.Contains(base, "..") {
			return mcp.NewToolResultError("adr key must be a bare markdown filename, e.g. 'adr:2026-07-05_title.md'"), nil
		}
		data, err := os.ReadFile(filepath.Join(h.decisionWriter.BasePath(), base))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ADR %q not found in %s", base, h.decisionWriter.BasePath())), nil
		}
		return mcp.NewToolResultText(string(data)), nil

	case "seed":
		seed, err := h.wisdomBase.GetSeed(key)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("seed %q not found", key)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("# Seed: %s\n\n%s", seed.Name, seed.Content)), nil
	}

	return mcp.NewToolResultError(fmt.Sprintf("unknown store %q — use memory, skill, adr, or seed", store)), nil
}

// ─── query legs (each returns positionally scored hits) ──────────────────────

func positionScore(i int) float64 {
	s := 1.0 - 0.1*float64(i)
	if s < 0.1 {
		return 0.1
	}
	return s
}

func (h *Handler) searchMemoryLeg(ctx context.Context, query string, limit int) ([]fetchHit, error) {
	entries, err := h.hub.SearchMemories(ctx, query, "", limit)
	if err != nil {
		return nil, err
	}
	hits := make([]fetchHit, 0, len(entries))
	for i, e := range entries {
		snippet := e.Snippet
		if snippet == "" {
			snippet = e.Description
		}
		detail := "type: " + e.Type
		if e.Updated != "" {
			detail += " · updated: " + e.Updated
		}
		hits = append(hits, fetchHit{
			store: "memory", key: e.Slug, title: e.Name,
			detail: detail, snippet: snippet, score: positionScore(i),
		})
	}
	return hits, nil
}

func (h *Handler) searchSkillLeg(query string, limit int) []fetchHit {
	results := h.skillsLoader.Search(query, limit)
	hits := make([]fetchHit, 0, len(results))
	for i, s := range results {
		hits = append(hits, fetchHit{
			store: "skill", key: s.Name, title: s.Name,
			detail:  "plugin: " + s.Plugin,
			snippet: firstSentence(s.Description),
			score:   positionScore(i),
		})
	}
	return hits
}

func (h *Handler) searchSeedLeg(query string, limit int) []fetchHit {
	var hits []fetchHit
	for _, sr := range h.wisdomBase.Search(query) {
		if sr.Type != "seed" {
			continue
		}
		seed, err := h.wisdomBase.GetSeed(sr.Name)
		if err != nil {
			continue
		}
		hits = append(hits, fetchHit{
			store: "seed", key: seed.Name, title: seed.Name,
			snippet: firstSentence(seed.Content),
			score:   positionScore(len(hits)),
		})
		if len(hits) >= limit {
			break
		}
	}
	return hits
}

// searchADRLeg matches query terms against ADR filenames, with a bounded
// head-content check over the newest files.
func (h *Handler) searchADRLeg(query string, limit int) []fetchHit {
	base := h.decisionWriter.BasePath()
	dirents, err := os.ReadDir(base)
	if err != nil {
		return nil
	}

	var names []string
	for _, d := range dirents {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			names = append(names, d.Name())
		}
	}
	// Date-prefixed filenames: lexicographic desc ≈ newest first.
	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	terms := queryTerms(query)
	if len(terms) == 0 {
		return nil
	}

	type scored struct {
		name  string
		score float64
	}
	var matches []scored
	for idx, name := range names {
		lower := strings.ToLower(name)
		var hitCount int
		for _, t := range terms {
			if strings.Contains(lower, t) {
				hitCount++
			}
		}
		score := float64(hitCount) / float64(len(terms))

		if hitCount == 0 && idx < adrContentScanCap {
			if head := readHead(filepath.Join(base, name), 4096); head != "" {
				lowerHead := strings.ToLower(head)
				var contentHits int
				for _, t := range terms {
					if strings.Contains(lowerHead, t) {
						contentHits++
					}
				}
				score = 0.6 * float64(contentHits) / float64(len(terms))
			}
		}
		if score > 0 {
			matches = append(matches, scored{name: name, score: score})
		}
	}

	sort.SliceStable(matches, func(i, j int) bool { return matches[i].score > matches[j].score })
	if len(matches) > limit {
		matches = matches[:limit]
	}

	hits := make([]fetchHit, 0, len(matches))
	for _, m := range matches {
		hits = append(hits, fetchHit{
			store: "adr", key: m.name,
			title:  strings.TrimSuffix(m.name, ".md"),
			detail: "decision record",
			score:  m.score,
		})
	}
	return hits
}

func queryTerms(query string) []string {
	var terms []string
	for _, f := range strings.Fields(strings.ToLower(query)) {
		f = strings.Trim(f, `"'.,;:!?`)
		if len(f) >= 3 {
			terms = append(terms, f)
		}
	}
	return terms
}

func readHead(path string, n int) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close() //nolint:errcheck // read-only handle
	buf := make([]byte, n)
	read, _ := f.Read(buf)
	return string(buf[:read])
}
