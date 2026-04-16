package config

import (
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear all relevant env vars so defaults apply.
	t.Setenv("DOJO_GATEWAY_URL", "")
	t.Setenv("DOJO_GATEWAY_TOKEN", "")
	t.Setenv("DOJO_SKILLS_PATH", "")
	t.Setenv("DOJO_ADR_PATH", "")

	cfg := Load()

	if cfg.GatewayURL != "http://localhost:7340" {
		t.Errorf("GatewayURL: got %q, want %q", cfg.GatewayURL, "http://localhost:7340")
	}
	if cfg.ADRPath != "./decisions" {
		t.Errorf("ADRPath: got %q, want %q", cfg.ADRPath, "./decisions")
	}
	if cfg.GatewayToken != "" {
		t.Errorf("GatewayToken: got %q, want empty string", cfg.GatewayToken)
	}
	if cfg.SkillsPath != "" {
		t.Errorf("SkillsPath: got %q, want empty string", cfg.SkillsPath)
	}
}

func TestLoad_EnvVarOverrides(t *testing.T) {
	t.Setenv("DOJO_GATEWAY_URL", "http://gateway.example.com:9000")
	t.Setenv("DOJO_GATEWAY_TOKEN", "secret-token-abc")
	t.Setenv("DOJO_SKILLS_PATH", "/opt/dojo/skills")
	t.Setenv("DOJO_ADR_PATH", "/opt/dojo/adr")

	cfg := Load()

	if cfg.GatewayURL != "http://gateway.example.com:9000" {
		t.Errorf("GatewayURL: got %q, want %q", cfg.GatewayURL, "http://gateway.example.com:9000")
	}
	if cfg.GatewayToken != "secret-token-abc" {
		t.Errorf("GatewayToken: got %q, want %q", cfg.GatewayToken, "secret-token-abc")
	}
	if cfg.SkillsPath != "/opt/dojo/skills" {
		t.Errorf("SkillsPath: got %q, want %q", cfg.SkillsPath, "/opt/dojo/skills")
	}
	if cfg.ADRPath != "/opt/dojo/adr" {
		t.Errorf("ADRPath: got %q, want %q", cfg.ADRPath, "/opt/dojo/adr")
	}
}

func TestLoad_EmptyEnvVarsFallBackToDefaults(t *testing.T) {
	// Explicitly setting to empty string should fall back to defaults for
	// fields that use envOr (GatewayURL and ADRPath).
	t.Setenv("DOJO_GATEWAY_URL", "")
	t.Setenv("DOJO_ADR_PATH", "")
	t.Setenv("DOJO_GATEWAY_TOKEN", "")
	t.Setenv("DOJO_SKILLS_PATH", "")

	cfg := Load()

	if cfg.GatewayURL != "http://localhost:7340" {
		t.Errorf("GatewayURL: got %q, want default %q", cfg.GatewayURL, "http://localhost:7340")
	}
	if cfg.ADRPath != "./decisions" {
		t.Errorf("ADRPath: got %q, want default %q", cfg.ADRPath, "./decisions")
	}
	// Fields with no fallback should remain empty.
	if cfg.GatewayToken != "" {
		t.Errorf("GatewayToken: got %q, want empty", cfg.GatewayToken)
	}
	if cfg.SkillsPath != "" {
		t.Errorf("SkillsPath: got %q, want empty", cfg.SkillsPath)
	}
}

func TestLoad_PartialOverride(t *testing.T) {
	// Override only one of the two defaulted fields.
	t.Setenv("DOJO_GATEWAY_URL", "http://custom:7340")
	t.Setenv("DOJO_ADR_PATH", "")
	t.Setenv("DOJO_GATEWAY_TOKEN", "tok")
	t.Setenv("DOJO_SKILLS_PATH", "")

	cfg := Load()

	if cfg.GatewayURL != "http://custom:7340" {
		t.Errorf("GatewayURL: got %q, want %q", cfg.GatewayURL, "http://custom:7340")
	}
	// ADR_PATH empty → default
	if cfg.ADRPath != "./decisions" {
		t.Errorf("ADRPath: got %q, want default %q", cfg.ADRPath, "./decisions")
	}
	if cfg.GatewayToken != "tok" {
		t.Errorf("GatewayToken: got %q, want %q", cfg.GatewayToken, "tok")
	}
}
