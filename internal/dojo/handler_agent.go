package dojo

import (
	"context"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp-server/internal/gateway"
	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) handleAgentList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — agent operations require a connected Gateway"), nil
	}

	agents, err := h.gw.Agents(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list agents: %v", err)), nil
	}

	if len(agents) == 0 {
		return mcp.NewToolResultText("No agents found.\n\nUse `dojo_agent_dispatch` to create and dispatch an agent."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Agents (%d)\n\n", len(agents))
	fmt.Fprintf(&sb, "| ID | Name | Mode | Status |\n")
	fmt.Fprintf(&sb, "|----|------|------|--------|\n")
	for _, a := range agents {
		name := a.Name
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(&sb, "| %s | %s | %s | %s |\n", a.ID, name, a.Mode, a.Status)
	}
	fmt.Fprintf(&sb, "\nUse `dojo_agent_chat` with an agent ID to send a message.\n")

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleAgentDispatch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — agent operations require a connected Gateway"), nil
	}

	var args struct {
		Name    string `json:"name"`
		Mode    string `json:"mode"`
		Message string `json:"message"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Message == "" {
		return mcp.NewToolResultError("'message' is required"), nil
	}
	if args.Mode == "" {
		args.Mode = "chat"
	}

	createReq := gateway.CreateAgentRequest{
		Name: args.Name,
		Mode: args.Mode,
	}
	created, err := h.gw.CreateAgent(ctx, createReq)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create agent: %v", err)), nil
	}

	response, err := h.gw.AgentChatSync(ctx, created.ID, args.Message)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Agent created (ID: %s) but chat failed: %v", created.ID, err)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Agent Dispatched\n\n")
	fmt.Fprintf(&sb, "**Agent ID:** %s\n", created.ID)
	if args.Name != "" {
		fmt.Fprintf(&sb, "**Name:** %s\n", args.Name)
	}
	fmt.Fprintf(&sb, "**Mode:** %s\n\n", args.Mode)
	fmt.Fprintf(&sb, "---\n\n")
	fmt.Fprintf(&sb, "**Response:**\n\n%s\n", response)

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleAgentChat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — agent operations require a connected Gateway"), nil
	}

	var args struct {
		AgentID string `json:"agent_id"`
		Message string `json:"message"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.AgentID == "" {
		return mcp.NewToolResultError("'agent_id' is required"), nil
	}
	if args.Message == "" {
		return mcp.NewToolResultError("'message' is required"), nil
	}

	response, err := h.gw.AgentChatSync(ctx, args.AgentID, args.Message)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Chat failed for agent %s: %v", args.AgentID, err)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "**Agent:** %s\n\n", args.AgentID)
	fmt.Fprintf(&sb, "---\n\n")
	fmt.Fprintf(&sb, "%s", strings.TrimSpace(response))

	return mcp.NewToolResultText(sb.String()), nil
}
