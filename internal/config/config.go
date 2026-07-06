package config

import (
	"os"
	"strconv"
)

// Config holds MCP server configuration.
type Config struct {
	GatewayURL    string
	GatewayToken  string
	SkillsPath    string
	ADRPath       string
	WorkspaceRoot string

	// HTTP mode (public endpoint). Unset HTTPAddr = stdio mode, exactly as before.
	HTTPAddr       string // DOJO_HTTP_ADDR, e.g. ":8091" — opt-in streamable-HTTP serving
	APIKeysRaw     string // DOJO_MCP_API_KEYS — comma-separated label:key pairs
	MemoryDBURL    string // DOJO_MEMORY_DB_URL — postgres:// DSN for the Memory Hub (read-only role)
	DispatchLabels string // DOJO_DISPATCH_ALLOWED_LABELS — key labels allowed to use dispatch-class tools
	DispatchRPM    int    // DOJO_DISPATCH_RATE_PER_MIN — per-label dispatch rate limit (default 6)
}

// Load reads configuration from environment variables.
func Load() *Config {
	workspaceRoot := os.Getenv("DOJO_WORKSPACE_ROOT")
	if workspaceRoot == "" {
		if wd, err := os.Getwd(); err == nil {
			workspaceRoot = wd
		}
	}

	return &Config{
		GatewayURL:    envOr("DOJO_GATEWAY_URL", "http://localhost:7340"),
		GatewayToken:  os.Getenv("DOJO_GATEWAY_TOKEN"),
		SkillsPath:    os.Getenv("DOJO_SKILLS_PATH"),
		ADRPath:       envOr("DOJO_ADR_PATH", "./decisions"),
		WorkspaceRoot: workspaceRoot,

		HTTPAddr:       os.Getenv("DOJO_HTTP_ADDR"),
		APIKeysRaw:     os.Getenv("DOJO_MCP_API_KEYS"),
		MemoryDBURL:    os.Getenv("DOJO_MEMORY_DB_URL"),
		DispatchLabels: os.Getenv("DOJO_DISPATCH_ALLOWED_LABELS"),
		DispatchRPM:    envIntOr("DOJO_DISPATCH_RATE_PER_MIN", 6),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
