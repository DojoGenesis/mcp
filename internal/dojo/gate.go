package dojo

import (
	"context"
	"log"
	"time"

	"github.com/DojoGenesis/mcp/internal/authz"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// DispatchClassTools spend LLM provider budget through the gateway. In HTTP
// mode they require the caller's key label to be on the dispatch allowlist
// and are rate limited per label. (dojo_scout degrades to its offline
// scaffold instead of being blocked — see handleScout.)
var DispatchClassTools = map[string]bool{
	"dojo_dispatch":       true,
	"dojo_agent_dispatch": true,
	"dojo_agent_chat":     true,
}

// GateMiddleware enforces dispatch-class authorization + rate limits and
// logs every tool call: tool name, key label, duration, outcome. Payloads
// and key material are never logged.
func GateMiddleware(limiter *authz.Limiter) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tool := req.Params.Name
			label, hasIdentity := authz.Label(ctx)
			if !hasIdentity {
				label = "local"
			}

			if DispatchClassTools[tool] {
				if !authz.DispatchAllowed(ctx) {
					log.Printf("tool_call tool=%s key=%s outcome=dispatch_denied", tool, label)
					return mcp.NewToolResultError(
						"This API key is not authorized for dispatch-class tools (they spend LLM provider budget). " +
							"The operator can enable it by adding the key's label to DOJO_DISPATCH_ALLOWED_LABELS."), nil
				}
				if !limiter.Allow(label) {
					log.Printf("tool_call tool=%s key=%s outcome=rate_limited", tool, label)
					return mcp.NewToolResultError("Dispatch rate limit exceeded for this key. Try again in a minute."), nil
				}
			}

			start := time.Now()
			res, err := next(ctx, req)

			outcome := "ok"
			switch {
			case err != nil:
				outcome = "error"
			case res != nil && res.IsError:
				outcome = "tool_error"
			}
			log.Printf("tool_call tool=%s key=%s dur_ms=%d outcome=%s",
				tool, label, time.Since(start).Milliseconds(), outcome)

			return res, err
		}
	}
}
