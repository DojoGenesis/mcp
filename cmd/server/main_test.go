package main

import (
	"testing"

	"github.com/DojoGenesis/mcp-server/internal/dojo"
	"github.com/mark3labs/mcp-go/server"
)

func TestServerInitializes(t *testing.T) {
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"2.1.0",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
	)
	if s == nil {
		t.Fatal("NewMCPServer returned nil")
	}
}

func TestHandlerRegistration(t *testing.T) {
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"2.1.0",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
	)

	dojoHandler := dojo.NewHandler()
	if dojoHandler == nil {
		t.Fatal("NewHandler returned nil")
	}

	// These should not panic
	dojoHandler.RegisterTools(s)
	dojoHandler.RegisterPrompts(s)
	dojoHandler.RegisterResources(s)
}
