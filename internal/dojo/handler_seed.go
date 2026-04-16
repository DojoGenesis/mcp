package dojo

import (
	"context"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp-server/internal/gateway"
	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) handleSeedList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Try Gateway first; fall back to local wisdom base if offline or gw not configured.
	if h.gw != nil {
		seeds, err := h.gw.Seeds(context.Background())
		if err == nil {
			if len(seeds) == 0 {
				return mcp.NewToolResultText("No seeds found in Gateway.\n\nUse `dojo_seed_create` to create a seed."), nil
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "# Seeds — Gateway (%d)\n\n", len(seeds))
			for i, s := range seeds {
				desc := s.Description
				if len(desc) > 100 {
					desc = desc[:100] + "..."
				}
				fmt.Fprintf(&sb, "## %d. %s\n", i+1, s.Name)
				if desc != "" {
					fmt.Fprintf(&sb, "%s\n", desc)
				}
				fmt.Fprintf(&sb, "\n")
			}
			return mcp.NewToolResultText(sb.String()), nil
		}
		// Gateway reachable but returned error — fall through to local
	}

	// Offline fallback: local wisdom base
	localSeeds := h.wisdomBase.ListSeeds()
	if len(localSeeds) == 0 {
		return mcp.NewToolResultText("No seeds available locally. Gateway is unreachable."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Seeds — Local (%d)\n\n", len(localSeeds))
	fmt.Fprintf(&sb, "_Gateway offline — showing local seeds only._\n\n")
	for i, s := range localSeeds {
		fmt.Fprintf(&sb, "## %d. %s\n", i+1, s.Name)
		if s.Description != "" {
			fmt.Fprintf(&sb, "%s\n", s.Description)
		}
		fmt.Fprintf(&sb, "\n")
	}
	fmt.Fprintf(&sb, "Use `dojo_apply_seed` to apply any seed to a situation.\n")

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleSeedCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — seed creation requires a connected Gateway"), nil
	}

	var args struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Name == "" {
		return mcp.NewToolResultError("'name' is required"), nil
	}
	if args.Content == "" {
		return mcp.NewToolResultError("'content' is required"), nil
	}

	req := gateway.CreateSeedRequest{
		Name:        args.Name,
		Description: args.Description,
		Content:     args.Content,
	}
	created, err := h.gw.CreateSeed(context.Background(), req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create seed: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Seed created.\n\n**ID:** %s\n**Name:** %s", created.ID, created.Name)), nil
}

func (h *Handler) handleSeedSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Query == "" {
		return mcp.NewToolResultError("'query' is required"), nil
	}

	queryLower := strings.ToLower(args.Query)

	// Try Gateway seeds; fall back to local.
	var matches []gateway.Seed
	if h.gw != nil {
		gwSeeds, err := h.gw.Seeds(context.Background())
		if err == nil {
			for _, s := range gwSeeds {
				if strings.Contains(strings.ToLower(s.Name), queryLower) ||
					strings.Contains(strings.ToLower(s.Description), queryLower) ||
					strings.Contains(strings.ToLower(s.Content), queryLower) {
					matches = append(matches, s)
				}
			}
			if len(matches) == 0 {
				return mcp.NewToolResultText(fmt.Sprintf("No seeds found matching: \"%s\"", args.Query)), nil
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "# Seed Search: \"%s\"\n\n", args.Query)
			fmt.Fprintf(&sb, "Found %d match(es) in Gateway:\n\n", len(matches))
			for i, s := range matches {
				content := s.Content
				if len(content) > 150 {
					content = content[:150] + "..."
				}
				fmt.Fprintf(&sb, "## %d. %s\n", i+1, s.Name)
				if s.Description != "" {
					fmt.Fprintf(&sb, "_%s_\n\n", s.Description)
				}
				fmt.Fprintf(&sb, "%s\n\n", content)
			}
			return mcp.NewToolResultText(sb.String()), nil
		}
	}

	// Offline fallback: search local wisdom base
	localSeeds := h.wisdomBase.ListSeeds()
	var localMatches []string
	for _, s := range localSeeds {
		if strings.Contains(strings.ToLower(s.Name), queryLower) ||
			strings.Contains(strings.ToLower(s.Description), queryLower) ||
			strings.Contains(strings.ToLower(s.Content), queryLower) {
			localMatches = append(localMatches, fmt.Sprintf("**%s** — %s", s.Name, firstSentence(s.Description)))
		}
	}

	if len(localMatches) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No seeds found matching: \"%s\" (Gateway offline, searched local seeds only)", args.Query)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Seed Search: \"%s\" (local, Gateway offline)\n\n", args.Query)
	fmt.Fprintf(&sb, "Found %d match(es):\n\n", len(localMatches))
	for _, line := range localMatches {
		fmt.Fprintf(&sb, "- %s\n", line)
	}
	fmt.Fprintf(&sb, "\nUse `dojo_apply_seed` to apply any seed to a situation.\n")

	return mcp.NewToolResultText(sb.String()), nil
}
