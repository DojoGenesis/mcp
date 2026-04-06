package main

import (
	"testing"

	"github.com/DojoGenesis/mcp-server/internal/dojo"
	"github.com/mark3labs/mcp-go/server"
)

func TestServerInitializes(t *testing.T) {
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"3.0.0",
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)
	if s == nil {
		t.Fatal("NewMCPServer returned nil")
	}
}

func TestHandlerRegistration(t *testing.T) {
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"3.0.0",
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)

	tmpDir := t.TempDir()
	dojoHandler, err := dojo.NewHandler("", tmpDir)
	if err != nil {
		t.Fatalf("NewHandler returned error: %v", err)
	}
	if dojoHandler == nil {
		t.Fatal("NewHandler returned nil")
	}

	// These should not panic
	dojoHandler.RegisterTools(s)
	dojoHandler.RegisterResources(s)
}

func TestHandlerRegistration_WithSkillsPath(t *testing.T) {
	s := server.NewMCPServer(
		"dojo-mcp-server",
		"3.0.0",
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)

	tmpDir := t.TempDir()
	// Non-existent skills path should fall back to bundled
	dojoHandler, err := dojo.NewHandler("/nonexistent", tmpDir)
	if err != nil {
		t.Fatalf("NewHandler returned error: %v", err)
	}

	dojoHandler.RegisterTools(s)
	dojoHandler.RegisterResources(s)
}
