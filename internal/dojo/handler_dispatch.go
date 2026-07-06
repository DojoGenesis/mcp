package dojo

import (
	"context"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp/internal/gateway"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleDispatch sends a prompt through the gateway's /v1/chat and returns
// the collected reply. Authorization and rate limiting happen in
// GateMiddleware before this runs (dispatch-class tool).
func (h *Handler) handleDispatch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.gw == nil {
		return mcp.NewToolResultError("Gateway is not configured — dispatch requires a connected Gateway"), nil
	}

	var args struct {
		Prompt    string `json:"prompt"`
		SessionID string `json:"session_id"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if strings.TrimSpace(args.Prompt) == "" {
		return mcp.NewToolResultError("'prompt' is required and cannot be empty"), nil
	}

	reply, err := h.gw.ChatSync(ctx, gateway.ChatRequest{
		Message:   args.Prompt,
		SessionID: args.SessionID,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Gateway dispatch failed: %v", err)), nil
	}
	if strings.TrimSpace(reply) == "" {
		return mcp.NewToolResultError("Gateway returned an empty reply"), nil
	}
	return mcp.NewToolResultText(reply), nil
}
