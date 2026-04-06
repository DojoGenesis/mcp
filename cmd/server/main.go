package main

import (
	"log"

	"github.com/DojoGenesis/mcp-server/internal/dojo"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"2.1.0",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
	)

	// Initialize Dojo handler
	dojoHandler := dojo.NewHandler()

	// Register tools
	dojoHandler.RegisterTools(s)

	// Register prompts
	dojoHandler.RegisterPrompts(s)

	// Register resources
	dojoHandler.RegisterResources(s)

	// Start server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
