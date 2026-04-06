package wisdom

import (
	"strings"
	"testing"
)

func TestGetSkills_Count(t *testing.T) {
	skills := getSkills()
	if len(skills) != 32 {
		t.Errorf("Expected 32 skills, got %d", len(skills))
	}
}

func TestGetSkills_AllHaveRequiredFields(t *testing.T) {
	skills := getSkills()
	for i, skill := range skills {
		if skill.Name == "" {
			t.Errorf("Skill at index %d has empty Name", i)
		}
		if skill.Description == "" {
			t.Errorf("Skill %q has empty Description", skill.Name)
		}
		if skill.Category == "" {
			t.Errorf("Skill %q has empty Category", skill.Name)
		}
		if skill.Content == "" {
			t.Errorf("Skill %q has empty Content", skill.Name)
		}
	}
}

func TestGetSkills_UniqueNames(t *testing.T) {
	skills := getSkills()
	seen := make(map[string]bool)
	for _, skill := range skills {
		if seen[skill.Name] {
			t.Errorf("Duplicate skill name: %q", skill.Name)
		}
		seen[skill.Name] = true
	}
}

func TestGetSkills_KnownSkillsPresent(t *testing.T) {
	skills := getSkills()
	expected := []string{
		"agent-to-agent-teaching",
		"patient-learning-protocol",
		"skill-creator",
		"strategic-scout",
		"pre-implementation-checklist",
		"skill-maintenance-ritual",
		"strategic-to-tactical-workflow",
		"transform-spec-to-implementation-prompt",
		"seed-reflector",
		"memory-garden-writer",
		"parallel-tracks-pattern",
		"iterative-scouting-pattern",
		"write-frontend-spec-from-backend",
		"product-positioning-scout",
		"retrospective",
		"multi-surface-product-strategy",
		"context-compression-ritual",
		"agent-handoff-protocol",
		"research-modes",
		"debugging-troubleshooting",
		"process-to-skill-workflow",
		"seed-to-skill-converter",
		"repo-context-sync",
		"project-exploration",
		"agent-workspace-navigator",
		"write-release-specification",
		"health-supervisor",
		"web-research",
		"status-writer",
		"decision-propagation-protocol",
		"era-architecture",
		"spec-constellation-to-prompt-suite",
	}

	nameSet := make(map[string]bool)
	for _, s := range skills {
		nameSet[s.Name] = true
	}

	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("Expected skill %q not found", name)
		}
	}
}

func TestListSkills(t *testing.T) {
	b := NewBase()
	skills := b.ListSkills()
	if len(skills) == 0 {
		t.Fatal("ListSkills returned empty slice")
	}
	if len(skills) != 32 {
		t.Errorf("ListSkills returned %d skills, expected 32", len(skills))
	}
}

func TestGetSkill_Found(t *testing.T) {
	b := NewBase()
	skill, err := b.GetSkill("strategic-scout")
	if err != nil {
		t.Fatalf("GetSkill returned error for existing skill: %v", err)
	}
	if skill == nil {
		t.Fatal("GetSkill returned nil for existing skill")
	}
	if skill.Name != "strategic-scout" {
		t.Errorf("GetSkill returned skill with name %q, want %q", skill.Name, "strategic-scout")
	}
	if skill.Content == "" {
		t.Error("GetSkill returned skill with empty content")
	}
}

func TestGetSkill_NotFound(t *testing.T) {
	b := NewBase()
	skill, err := b.GetSkill("nonexistent_skill_xyz")
	if err == nil {
		t.Error("GetSkill did not return error for nonexistent skill")
	}
	if skill != nil {
		t.Error("GetSkill returned non-nil skill for nonexistent name")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error message %q does not contain 'not found'", err.Error())
	}
}

func TestGetSkill_AllAccessible(t *testing.T) {
	b := NewBase()
	skills := b.ListSkills()

	for _, skill := range skills {
		retrieved, err := b.GetSkill(skill.Name)
		if err != nil {
			t.Errorf("GetSkill(%q) returned error: %v", skill.Name, err)
			continue
		}
		if retrieved == nil {
			t.Errorf("GetSkill(%q) returned nil", skill.Name)
			continue
		}
		if retrieved.Name != skill.Name {
			t.Errorf("GetSkill(%q) returned skill with name %q", skill.Name, retrieved.Name)
		}
	}
}

func TestSearchSkills_Found(t *testing.T) {
	b := NewBase()
	results := b.SearchSkills("strategic")
	if len(results) == 0 {
		t.Fatal("SearchSkills for 'strategic' returned no results")
	}
	// Should find at least strategic-scout (name contains "strategic")
	found := false
	for _, r := range results {
		if r.Name == "strategic-scout" {
			found = true
			break
		}
	}
	if !found {
		t.Error("SearchSkills for 'strategic' did not find 'strategic-scout'")
	}
}

func TestSearchSkills_NoMatch(t *testing.T) {
	b := NewBase()
	results := b.SearchSkills("xyzzynonexistent12345")
	if len(results) != 0 {
		t.Errorf("SearchSkills for gibberish returned %d results, expected 0", len(results))
	}
}

func TestSearchSkills_CaseInsensitive(t *testing.T) {
	b := NewBase()
	lower := b.SearchSkills("debugging")
	upper := b.SearchSkills("DEBUGGING")

	if len(lower) == 0 {
		t.Fatal("SearchSkills for 'debugging' returned no results")
	}
	if len(lower) != len(upper) {
		t.Errorf("Case sensitivity: 'debugging' returned %d results, 'DEBUGGING' returned %d", len(lower), len(upper))
	}
}
