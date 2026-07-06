package dojo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp/internal/memhub"
	"github.com/mark3labs/mcp-go/mcp"
)

const hubNotConfigured = "Memory Hub is not configured on this server (DOJO_MEMORY_DB_URL unset). " +
	"Hub-backed memory tools are available on deployments wired to the Postgres Memory Hub."

func (h *Handler) handleSearchMemoryHub(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.hub == nil {
		return mcp.NewToolResultError(hubNotConfigured), nil
	}

	var args struct {
		Query string `json:"query"`
		Type  string `json:"type"`
		Limit int    `json:"limit"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if strings.TrimSpace(args.Query) == "" {
		return mcp.NewToolResultError("'query' is required and cannot be empty"), nil
	}

	entries, err := h.hub.SearchMemories(ctx, args.Query, args.Type, args.Limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Memory Hub search failed: %v", err)), nil
	}
	if len(entries) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No hub memories match: \"%s\"\n\nTry broader terms, websearch syntax (\"exact phrase\", OR), or dojo_recent_memories for orientation.", args.Query)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Memory Hub: \"%s\"", args.Query)
	if args.Type != "" {
		fmt.Fprintf(&sb, " (type: %s)", args.Type)
	}
	fmt.Fprintf(&sb, "\n\n%d match(es):\n\n", len(entries))
	for i, e := range entries {
		fmt.Fprintf(&sb, "## %d. %s `%s`\n", i+1, e.Name, e.Slug)
		fmt.Fprintf(&sb, "**Type:** %s", e.Type)
		if e.Updated != "" {
			fmt.Fprintf(&sb, " · **Updated:** %s", e.Updated)
		}
		sb.WriteString("\n")
		if e.Description != "" {
			fmt.Fprintf(&sb, "%s\n", e.Description)
		}
		if e.Snippet != "" {
			fmt.Fprintf(&sb, "> %s\n", strings.ReplaceAll(e.Snippet, "\n", " "))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("Full body: `dojo_get_memory` with the slug.\n")
	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleGetMemoryHub(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.hub == nil {
		return mcp.NewToolResultError(hubNotConfigured), nil
	}

	var args struct {
		Slug string `json:"slug"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if strings.TrimSpace(args.Slug) == "" {
		return mcp.NewToolResultError("'slug' is required and cannot be empty"), nil
	}

	e, err := h.hub.GetMemory(ctx, args.Slug)
	if errors.Is(err, memhub.ErrNotFound) {
		return mcp.NewToolResultError(fmt.Sprintf("No hub memory with slug %q. Find slugs via dojo_search_memory or dojo_recent_memories.", args.Slug)), nil
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Memory Hub get failed: %v", err)), nil
	}

	return mcp.NewToolResultText(renderHubEntry(e)), nil
}

func (h *Handler) handleRecentMemoriesHub(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.hub == nil {
		return mcp.NewToolResultError(hubNotConfigured), nil
	}

	var args struct {
		Type  string `json:"type"`
		Limit int    `json:"limit"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	entries, err := h.hub.RecentMemories(ctx, args.Type, args.Limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Memory Hub recent failed: %v", err)), nil
	}
	if len(entries) == 0 {
		return mcp.NewToolResultText("The Memory Hub returned no entries."), nil
	}

	var sb strings.Builder
	sb.WriteString("# Memory Hub — most recently updated")
	if args.Type != "" {
		fmt.Fprintf(&sb, " (type: %s)", args.Type)
	}
	sb.WriteString("\n\n")
	for i, e := range entries {
		fmt.Fprintf(&sb, "%d. **%s** `%s` (%s", i+1, e.Name, e.Slug, e.Type)
		if e.Updated != "" {
			fmt.Fprintf(&sb, ", %s", e.Updated)
		}
		sb.WriteString(")")
		if e.Description != "" {
			fmt.Fprintf(&sb, " — %s", firstSentence(e.Description))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\nFull body: `dojo_get_memory` with the slug.\n")
	return mcp.NewToolResultText(sb.String()), nil
}

func renderHubEntry(e *memhub.Entry) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", e.Name)
	fmt.Fprintf(&sb, "**Slug:** `%s` · **Type:** %s", e.Slug, e.Type)
	if e.Updated != "" {
		fmt.Fprintf(&sb, " · **Updated:** %s", e.Updated)
	}
	sb.WriteString("\n\n")
	if e.Description != "" {
		fmt.Fprintf(&sb, "_%s_\n\n---\n\n", e.Description)
	}
	sb.WriteString(e.Body)
	return sb.String()
}
