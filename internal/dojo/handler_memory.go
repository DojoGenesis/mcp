package dojo

import (
	"context"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp-server/internal/gateway"
	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) handleMemoryList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — memory operations require a connected Gateway"), nil
	}

	memories, err := h.gw.Memories(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Gateway error: %v", err)), nil
	}

	if len(memories) == 0 {
		return mcp.NewToolResultText("No memories stored yet.\n\nUse `dojo_memory_store` to save a memory."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Stored Memories (%d)\n\n", len(memories))
	for i, m := range memories {
		content := m.Content
		if len(content) > 120 {
			content = content[:120] + "..."
		}
		fmt.Fprintf(&sb, "## %d. [%s] %s\n", i+1, m.ID, m.Type)
		fmt.Fprintf(&sb, "%s\n\n", content)
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleMemoryStore(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — memory operations require a connected Gateway"), nil
	}

	var args struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Content == "" {
		return mcp.NewToolResultError("'content' is required"), nil
	}
	if args.Type == "" {
		args.Type = "general"
	}

	req := gateway.StoreMemoryRequest{
		Content: args.Content,
		Type:    args.Type,
	}
	stored, err := h.gw.StoreMemory(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to store memory: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Memory stored.\n\n**ID:** %s\n**Type:** %s", stored.ID, stored.Type)), nil
}

func (h *Handler) handleMemorySearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — memory operations require a connected Gateway"), nil
	}

	var args struct {
		Query string `json:"query"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Query == "" {
		return mcp.NewToolResultError("'query' is required"), nil
	}

	matches, err := h.gw.SearchMemories(ctx, args.Query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	if len(matches) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No memories found matching: \"%s\"", args.Query)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Memory Search: \"%s\"\n\n", args.Query)
	fmt.Fprintf(&sb, "Found %d match(es):\n\n", len(matches))
	for i, m := range matches {
		content := m.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		fmt.Fprintf(&sb, "## %d. [%s] %s\n", i+1, m.ID, m.Type)
		fmt.Fprintf(&sb, "%s\n\n", content)
	}

	return mcp.NewToolResultText(sb.String()), nil
}
