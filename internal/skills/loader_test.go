package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatter_Valid(t *testing.T) {
	content := `---
name: test-skill
description: A test skill for unit testing.
---

# Test Skill

This is the body content.`

	fm, body, err := parseFrontmatter(content)
	if err != nil {
		t.Fatalf("parseFrontmatter returned error: %v", err)
	}
	if fm.Name != "test-skill" {
		t.Errorf("name = %q, want %q", fm.Name, "test-skill")
	}
	if fm.Description != "A test skill for unit testing." {
		t.Errorf("description = %q, want %q", fm.Description, "A test skill for unit testing.")
	}
	if !strings.Contains(body, "# Test Skill") {
		t.Error("body does not contain expected heading")
	}
	if !strings.Contains(body, "body content") {
		t.Error("body does not contain expected content")
	}
}

func TestParseFrontmatter_NoDelimiter(t *testing.T) {
	content := "# Just a regular markdown file\n\nNo frontmatter here."
	_, _, err := parseFrontmatter(content)
	if err == nil {
		t.Error("expected error for content without frontmatter, got nil")
	}
}

func TestParseFrontmatter_NoClosingDelimiter(t *testing.T) {
	content := "---\nname: broken\n# No closing delimiter"
	_, _, err := parseFrontmatter(content)
	if err == nil {
		t.Error("expected error for content without closing delimiter, got nil")
	}
}

func TestParseFrontmatter_EmptyYAML(t *testing.T) {
	content := "---\n---\n\nBody here."
	fm, body, err := parseFrontmatter(content)
	if err != nil {
		t.Fatalf("parseFrontmatter returned error: %v", err)
	}
	if fm.Name != "" {
		t.Errorf("expected empty name, got %q", fm.Name)
	}
	if !strings.Contains(body, "Body here") {
		t.Error("body does not contain expected content")
	}
}

func TestParseSkillFile_Valid(t *testing.T) {
	content := `---
name: strategic-scout
description: Explore strategic tensions. Use when facing product decisions.
---

# Strategic Scout

Step 1: Identify the tension.`

	skill, err := parseSkillFile(content, "strategic-thinking", "/test/path/SKILL.md")
	if err != nil {
		t.Fatalf("parseSkillFile returned error: %v", err)
	}
	if skill.Name != "strategic-scout" {
		t.Errorf("name = %q, want %q", skill.Name, "strategic-scout")
	}
	if skill.Plugin != "strategic-thinking" {
		t.Errorf("plugin = %q, want %q", skill.Plugin, "strategic-thinking")
	}
	if skill.FilePath != "/test/path/SKILL.md" {
		t.Errorf("filepath = %q, want %q", skill.FilePath, "/test/path/SKILL.md")
	}
	if !strings.Contains(skill.Content, "Strategic Scout") {
		t.Error("content does not contain expected heading")
	}
}

func TestParseSkillFile_MissingName(t *testing.T) {
	content := `---
description: No name field here
---

Body.`

	_, err := parseSkillFile(content, "test", "/test")
	if err == nil {
		t.Error("expected error for SKILL.md with missing name")
	}
}

func TestExtractTriggers_UseWhen(t *testing.T) {
	desc := "Explore strategic tensions. Use when facing product decisions, strategic tensions, or exploring multiple possible directions."
	triggers := extractTriggers(desc)
	if len(triggers) == 0 {
		t.Fatal("extractTriggers returned no triggers")
	}
	found := false
	for _, tr := range triggers {
		if strings.Contains(tr, "product decisions") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("triggers %v do not contain expected 'product decisions'", triggers)
	}
}

func TestExtractTriggers_TriggerPhrases(t *testing.T) {
	desc := "Strategic scouting. Trigger phrases: 'scout this tension', 'explore multiple routes', 'what are our options'."
	triggers := extractTriggers(desc)
	if len(triggers) == 0 {
		t.Fatal("extractTriggers returned no triggers")
	}
	found := false
	for _, tr := range triggers {
		if strings.Contains(tr, "scout this tension") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("triggers %v do not contain expected 'scout this tension'", triggers)
	}
}

func TestExtractTriggers_NoTriggers(t *testing.T) {
	desc := "A simple skill description without trigger phrases."
	triggers := extractTriggers(desc)
	// Should return empty (no trigger keywords)
	if len(triggers) != 0 {
		t.Errorf("expected no triggers, got %v", triggers)
	}
}

func TestNewLoader_BundledFallback(t *testing.T) {
	// Empty path should use bundled skills
	loader, err := NewLoader("")
	if err != nil {
		t.Fatalf("NewLoader with empty path returned error: %v", err)
	}
	if loader.Count() == 0 {
		t.Fatal("NewLoader with empty path loaded 0 bundled skills")
	}
	if loader.Count() != 35 {
		t.Errorf("expected 35 bundled skills, got %d", loader.Count())
	}
}

func TestNewLoader_NonexistentPath(t *testing.T) {
	// Nonexistent path should fall back to bundled
	loader, err := NewLoader("/nonexistent/path/that/does/not/exist")
	if err != nil {
		t.Fatalf("NewLoader with nonexistent path returned error: %v", err)
	}
	if loader.Count() == 0 {
		t.Fatal("NewLoader with nonexistent path loaded 0 skills (should fall back to bundled)")
	}
}

func TestNewLoader_FromFilesystem(t *testing.T) {
	// Create a temporary directory with skill files
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins", "test-plugin", "skills", "test-skill")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: test-skill
description: A test skill.
---

# Test Skill

Body content here.`

	if err := os.WriteFile(filepath.Join(pluginDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader, err := NewLoader(tmpDir)
	if err != nil {
		t.Fatalf("NewLoader returned error: %v", err)
	}
	if loader.Count() != 1 {
		t.Errorf("expected 1 skill, got %d", loader.Count())
	}

	skill, err := loader.GetByName("test-skill")
	if err != nil {
		t.Fatalf("GetByName returned error: %v", err)
	}
	if skill.Plugin != "test-plugin" {
		t.Errorf("plugin = %q, want %q", skill.Plugin, "test-plugin")
	}
}

func TestNewLoader_MultiplePlugins(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two plugins with skills
	for _, plugin := range []string{"plugin-a", "plugin-b"} {
		for _, skill := range []string{"skill-1", "skill-2"} {
			dir := filepath.Join(tmpDir, "plugins", plugin, "skills", skill)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			content := "---\nname: " + plugin + "-" + skill + "\ndescription: test\n---\n\nBody."
			if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	loader, err := NewLoader(tmpDir)
	if err != nil {
		t.Fatalf("NewLoader returned error: %v", err)
	}
	if loader.Count() != 4 {
		t.Errorf("expected 4 skills, got %d", loader.Count())
	}
	plugins := loader.PluginNames()
	if len(plugins) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(plugins))
	}
}

func TestGetByName_Found(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	skill, err := loader.GetByName("strategic-scout")
	if err != nil {
		t.Fatalf("GetByName returned error: %v", err)
	}
	if skill.Name != "strategic-scout" {
		t.Errorf("name = %q, want %q", skill.Name, "strategic-scout")
	}
}

func TestGetByName_NotFound(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	_, err = loader.GetByName("nonexistent-skill-xyz")
	if err == nil {
		t.Error("expected error for nonexistent skill, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error %q does not contain 'not found'", err.Error())
	}
}

func TestListByPlugin(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	byPlugin := loader.ListByPlugin()
	if len(byPlugin) == 0 {
		t.Fatal("ListByPlugin returned empty map")
	}
	// Bundled skills should have multiple plugins
	total := 0
	for _, skills := range byPlugin {
		total += len(skills)
	}
	if total != loader.Count() {
		t.Errorf("total skills in ListByPlugin = %d, Count = %d", total, loader.Count())
	}
}

func TestAllSkills(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	all := loader.AllSkills()
	if len(all) != loader.Count() {
		t.Errorf("AllSkills returned %d skills, Count = %d", len(all), loader.Count())
	}
	for i, s := range all {
		if s.Name == "" {
			t.Errorf("skill at index %d has empty name", i)
		}
		if s.Plugin == "" {
			t.Errorf("skill %q has empty plugin", s.Name)
		}
	}
}

func TestPluginNames_Sorted(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	names := loader.PluginNames()
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("PluginNames not sorted: %q comes after %q", names[i], names[i-1])
		}
	}
}

func TestBundledSkills_AllHaveRequiredFields(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range loader.AllSkills() {
		if s.Name == "" {
			t.Error("bundled skill has empty name")
		}
		if s.Description == "" {
			t.Errorf("bundled skill %q has empty description", s.Name)
		}
		if s.Plugin == "" {
			t.Errorf("bundled skill %q has empty plugin", s.Name)
		}
		if s.Content == "" {
			t.Errorf("bundled skill %q has empty content", s.Name)
		}
	}
}

func TestBundledSkills_KnownSkillsPresent(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"strategic-scout",
		"release-specification",
		"implementation-prompt",
		"debugging-troubleshooting",
		"retrospective",
		"pre-implementation-checklist",
		"parallel-tracks",
		"health-audit",
		"seed-extraction",
		"memory-garden",
		"context-ingestion",
		"research-modes",
		"skill-creation",
		"handoff-protocol",
		"status-writing",
	}

	for _, name := range expected {
		_, err := loader.GetByName(name)
		if err != nil {
			t.Errorf("expected bundled skill %q not found", name)
		}
	}
}

func TestSkipInvalidFiles(t *testing.T) {
	tmpDir := t.TempDir()
	dir := filepath.Join(tmpDir, "plugins", "test", "skills", "bad-skill")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write a SKILL.md without valid frontmatter
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# No frontmatter\n\nJust content."), 0644); err != nil {
		t.Fatal(err)
	}

	// Also write a valid one
	dir2 := filepath.Join(tmpDir, "plugins", "test", "skills", "good-skill")
	if err := os.MkdirAll(dir2, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\nname: good-skill\ndescription: works\n---\n\nBody."
	if err := os.WriteFile(filepath.Join(dir2, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader, err := NewLoader(tmpDir)
	if err != nil {
		t.Fatalf("NewLoader returned error: %v", err)
	}
	// Should only load the valid one
	if loader.Count() != 1 {
		t.Errorf("expected 1 valid skill (skipping invalid), got %d", loader.Count())
	}
}
