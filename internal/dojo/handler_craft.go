package dojo

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) handleConverge(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Run git status --short
	statusOut, statusErr := runGit("status", "--short")
	// Run git log --oneline --since="7 days ago"
	logOut, logErr := runGit("log", "--oneline", "--since=7 days ago")

	var sb strings.Builder

	// Count dirty files
	dirtyCount := 0
	if statusErr == nil && strings.TrimSpace(statusOut) != "" {
		lines := strings.Split(strings.TrimSpace(statusOut), "\n")
		dirtyCount = len(lines)
	}

	// Count recent commits
	commitCount := 0
	if logErr == nil && strings.TrimSpace(logOut) != "" {
		lines := strings.Split(strings.TrimSpace(logOut), "\n")
		commitCount = len(lines)
	}

	// Determine signal
	signal := "GREEN"
	if dirtyCount >= 25 {
		signal = "RED"
	} else if dirtyCount >= 10 {
		signal = "YELLOW"
	}

	fmt.Fprintf(&sb, "# Convergence Check: %s\n\n", signal)
	fmt.Fprintf(&sb, "| Metric | Value |\n")
	fmt.Fprintf(&sb, "|--------|-------|\n")
	fmt.Fprintf(&sb, "| Dirty files | %d |\n", dirtyCount)
	fmt.Fprintf(&sb, "| Commits (last 7 days) | %d |\n", commitCount)

	// Optionally enrich with Gateway metrics
	if h.gw != nil {
		if memories, err := h.gw.Memories(ctx); err == nil {
			fmt.Fprintf(&sb, "| Stored memories | %d |\n", len(memories))
		}
		if seeds, err := h.gw.Seeds(ctx); err == nil {
			fmt.Fprintf(&sb, "| Gateway seeds | %d |\n", len(seeds))
		}
	}

	fmt.Fprintf(&sb, "\n")

	switch signal {
	case "RED":
		fmt.Fprintf(&sb, "**RED** — %d dirty files. Commit or stash before proceeding. Run `/converge` after cleanup.\n", dirtyCount)
	case "YELLOW":
		fmt.Fprintf(&sb, "**YELLOW** — %d dirty files. Consider committing work-in-progress before the next dispatch.\n", dirtyCount)
	default:
		fmt.Fprintf(&sb, "**GREEN** — workspace is clean. Ready for new work.\n")
	}

	if statusErr == nil && strings.TrimSpace(statusOut) != "" {
		fmt.Fprintf(&sb, "\n## Dirty Files\n\n```\n%s\n```\n", strings.TrimSpace(statusOut))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleHealth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured"), nil
	}

	status, err := h.gw.Health(ctx)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf(
			"# Gateway Health: UNREACHABLE\n\nError: %v", err,
		)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Gateway Health: %s\n\n", strings.ToUpper(status.Status))
	fmt.Fprintf(&sb, "| Field | Value |\n")
	fmt.Fprintf(&sb, "|-------|-------|\n")
	fmt.Fprintf(&sb, "| Status | %s |\n", status.Status)
	if status.Version != "" {
		fmt.Fprintf(&sb, "| Version | %s |\n", status.Version)
	}
	if status.Timestamp != "" {
		fmt.Fprintf(&sb, "| Timestamp | %s |\n", status.Timestamp)
	}
	if status.UptimeSeconds > 0 {
		hours := status.UptimeSeconds / 3600
		mins := (status.UptimeSeconds % 3600) / 60
		fmt.Fprintf(&sb, "| Uptime | %dh %dm |\n", hours, mins)
	}
	if len(status.Providers) > 0 {
		providerNames := make([]string, 0, len(status.Providers))
		for name := range status.Providers {
			providerNames = append(providerNames, name)
		}
		fmt.Fprintf(&sb, "| Providers | %s |\n", strings.Join(providerNames, ", "))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleDispositionList(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type disposition struct {
		Name        string
		Description string
		BestFor     string
	}

	dispositions := []disposition{
		{
			Name:        "focused",
			Description: "High-intensity single-track execution. All attention on one deliverable.",
			BestFor:     "Deadline-driven sprints, critical bug fixes, shipping a specific feature.",
		},
		{
			Name:        "balanced",
			Description: "Steady parallel progress across 2-3 tracks. Sustainable pace.",
			BestFor:     "Normal development cycles, maintaining multiple active projects.",
		},
		{
			Name:        "exploratory",
			Description: "Wide research and experimentation. Low commitment, high surface area.",
			BestFor:     "Architecture decisions, technology evaluation, problem diagnosis.",
		},
		{
			Name:        "deliberate",
			Description: "Slow, reflective, documentation-heavy. Quality over velocity.",
			BestFor:     "Post-ship stabilization, documentation sprints, knowledge capture.",
		},
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Dispositions\n\n")
	fmt.Fprintf(&sb, "| Name | Description | Best For |\n")
	fmt.Fprintf(&sb, "|------|-------------|----------|\n")
	for _, d := range dispositions {
		fmt.Fprintf(&sb, "| **%s** | %s | %s |\n", d.Name, d.Description, d.BestFor)
	}
	fmt.Fprintf(&sb, "\nTo apply a disposition, update `settings.json` with `\"disposition\": \"<name>\"` or use `/settings profile`.\n")

	return mcp.NewToolResultText(sb.String()), nil
}

// runGit runs a git command and returns stdout and any error.
func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(errBuf.String()))
	}
	return out.String(), nil
}
