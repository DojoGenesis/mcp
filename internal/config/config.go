package config

import (
	"os"
)

// Config holds MCP server configuration.
type Config struct {
	GatewayURL   string
	GatewayToken string
	SkillsPath   string
	ADRPath      string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		GatewayURL:   envOr("DOJO_GATEWAY_URL", "http://localhost:7340"),
		GatewayToken: os.Getenv("DOJO_GATEWAY_TOKEN"),
		SkillsPath:   os.Getenv("DOJO_SKILLS_PATH"),
		ADRPath:      envOr("DOJO_ADR_PATH", "./decisions"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
