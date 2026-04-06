package main

import (
	"log"
	"os"

	"github.com/DojoGenesis/mcp-server/internal/dojo"
	"github.com/mark3labs/mcp-go/server"
)

const version = "3.0.0"

func main() {
	skillsPath := os.Getenv("DOJO_SKILLS_PATH")
	adrPath := os.Getenv("DOJO_ADR_PATH")

	s := server.NewMCPServer(
		"dojo-mcp-server",
		version,
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)

	handler, err := dojo.NewHandler(skillsPath, adrPath)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	handler.RegisterTools(s)
	handler.RegisterResources(s)

	log.Printf("dojo-mcp-server v%s starting (skills: %s, adr: %s)",
		version, skillsPath, adrPath)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
