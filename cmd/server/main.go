package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/DojoGenesis/mcp/internal/authz"
	"github.com/DojoGenesis/mcp/internal/config"
	"github.com/DojoGenesis/mcp/internal/dojo"
	"github.com/DojoGenesis/mcp/internal/gateway"
	"github.com/DojoGenesis/mcp/internal/httpserver"
	"github.com/mark3labs/mcp-go/server"
)

// version is set via ldflags at build time: -X main.version=vX.Y.Z
var version = "3.2.0"

func main() {
	cfg := config.Load()

	gw := gateway.New(cfg.GatewayURL, cfg.GatewayToken)

	ctx := context.Background()
	if gw.IsOnline(ctx) {
		log.Printf("dojo-mcp-server: gateway online at %s", cfg.GatewayURL)
	} else {
		log.Printf("dojo-mcp-server: gateway offline at %s (offline mode active)", cfg.GatewayURL)
	}

	handler, err := dojo.NewHandler(cfg.SkillsPath, cfg.ADRPath, gw, cfg.WorkspaceRoot)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	if cfg.HTTPAddr != "" {
		runHTTP(cfg, handler)
		return
	}
	runStdio(cfg, handler)
}

// runStdio serves MCP over stdio — the local (Lane A) mode, unchanged.
func runStdio(cfg *config.Config, handler *dojo.Handler) {
	s := newMCPServer(handler)

	log.Printf("dojo-mcp-server v%s starting (skills: %s, adr: %s)",
		version, cfg.SkillsPath, cfg.ADRPath)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// runHTTP serves MCP over streamable HTTP behind Bearer-key auth — the
// public endpoint (Lane B) mode. Refuses to start without a valid key table.
func runHTTP(cfg *config.Config, handler *dojo.Handler) {
	keys, err := authz.ParseKeys(cfg.APIKeysRaw)
	if err != nil {
		log.Fatalf("DOJO_MCP_API_KEYS invalid: %v", err)
	}
	if keys.Len() == 0 {
		log.Fatalf("DOJO_HTTP_ADDR is set but DOJO_MCP_API_KEYS is empty — refusing to serve an unauthenticated public endpoint")
	}

	dispatchLabels := parseDispatchLabels(cfg.DispatchLabels, keys)
	limiter := authz.NewLimiter(cfg.DispatchRPM)

	s := newMCPServer(handler, server.WithToolHandlerMiddleware(dojo.GateMiddleware(limiter)))

	streamable := server.NewStreamableHTTPServer(s,
		server.WithStateLess(true),
		server.WithHeartbeatInterval(30*time.Second),
		server.WithHTTPContextFunc(httpserver.IdentityContextFunc),
	)

	h := httpserver.NewHandler(httpserver.Options{
		Version:        version,
		Keys:           keys,
		DispatchLabels: dispatchLabels,
		MCP:            streamable,
	})

	log.Printf("dojo-mcp-server v%s serving HTTP on %s (keys: %v; dispatch-enabled: %v; dispatch rate: %d/min; skills: %s, adr: %s)",
		version, cfg.HTTPAddr, keys.Labels(), sortedLabels(dispatchLabels), cfg.DispatchRPM, cfg.SkillsPath, cfg.ADRPath)

	if err := httpserver.Run(cfg.HTTPAddr, h); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

// newMCPServer builds the MCP server and registers all dojo capabilities.
func newMCPServer(handler *dojo.Handler, extra ...server.ServerOption) *server.MCPServer {
	opts := append([]server.ServerOption{
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	}, extra...)

	s := server.NewMCPServer("dojo-mcp-server", version, opts...)
	handler.RegisterTools(s)
	handler.RegisterResources(s)
	return s
}

// parseDispatchLabels resolves the dispatch allowlist against the key table,
// warning on labels that don't correspond to any configured key.
func parseDispatchLabels(raw string, keys *authz.KeySet) map[string]bool {
	out := make(map[string]bool)
	known := make(map[string]bool, keys.Len())
	for _, l := range keys.Labels() {
		known[l] = true
	}
	for _, l := range strings.Split(raw, ",") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if !known[l] {
			log.Printf("WARN: DOJO_DISPATCH_ALLOWED_LABELS contains %q which matches no configured key label", l)
			continue
		}
		out[l] = true
	}
	return out
}

func sortedLabels(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for l := range m {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}
