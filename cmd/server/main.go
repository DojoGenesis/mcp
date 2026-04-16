package main

import (
	"context"
	"log"

	"github.com/DojoGenesis/mcp-server/internal/config"
	"github.com/DojoGenesis/mcp-server/internal/dojo"
	"github.com/DojoGenesis/mcp-server/internal/gateway"
	"github.com/mark3labs/mcp-go/server"
)

// version is set via ldflags at build time: -X main.version=vX.Y.Z
var version = "3.1.0"

func main() {
	cfg := config.Load()

	s := server.NewMCPServer(
		"dojo-mcp-server",
		version,
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)

	gw := gateway.New(cfg.GatewayURL, cfg.GatewayToken)

	ctx := context.Background()
	if gw.IsOnline(ctx) {
		log.Printf("dojo-mcp-server: gateway online at %s", cfg.GatewayURL)
	} else {
		log.Printf("dojo-mcp-server: gateway offline at %s (offline mode active)", cfg.GatewayURL)
	}

	handler, err := dojo.NewHandler(cfg.SkillsPath, cfg.ADRPath, gw)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	handler.RegisterTools(s)
	handler.RegisterResources(s)

	log.Printf("dojo-mcp-server v%s starting (skills: %s, adr: %s)",
		version, cfg.SkillsPath, cfg.ADRPath)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
