package skills

import (
	"strings"
	"testing"
)

func TestSearch_ExactNameMatch(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("strategic-scout", 5)
	if len(results) == 0 {
		t.Fatal("Search for exact name returned no results")
	}
	if results[0].Name != "strategic-scout" {
		t.Errorf("first result = %q, want %q", results[0].Name, "strategic-scout")
	}
}

func TestSearch_PartialNameMatch(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("debug", 5)
	if len(results) == 0 {
		t.Fatal("Search for 'debug' returned no results")
	}
	found := false
	for _, r := range results {
		if strings.Contains(r.Name, "debug") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Search for 'debug' did not return any skill with 'debug' in name")
	}
}

func TestSearch_DescriptionMatch(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("retrospective", 5)
	if len(results) == 0 {
		t.Fatal("Search for 'retrospective' returned no results")
	}
}

func TestSearch_NoMatch(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("xyzzynonexistent12345", 5)
	if len(results) != 0 {
		t.Errorf("Search for gibberish returned %d results, expected 0", len(results))
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("", 5)
	if len(results) == 0 {
		t.Fatal("Search with empty query should return skills")
	}
	if len(results) > 5 {
		t.Errorf("Search with maxResults=5 returned %d results", len(results))
	}
}

func TestSearch_EmptyQueryNoLimit(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("", 0)
	if len(results) != loader.Count() {
		t.Errorf("Search with empty query and no limit returned %d, expected %d", len(results), loader.Count())
	}
}

func TestSearch_MaxResultsRespected(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("skill", 2)
	if len(results) > 2 {
		t.Errorf("Search with maxResults=2 returned %d results", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	lower := loader.Search("debugging", 5)
	upper := loader.Search("DEBUGGING", 5)
	if len(lower) == 0 {
		t.Fatal("Search for 'debugging' returned no results")
	}
	if len(lower) != len(upper) {
		t.Errorf("case sensitivity: 'debugging' returned %d, 'DEBUGGING' returned %d", len(lower), len(upper))
	}
}

func TestSearch_RelevanceOrder(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("strategic-scout", 5)
	if len(results) < 2 {
		t.Skip("not enough results to test ordering")
	}
	// The exact name match should be first
	if results[0].Name != "strategic-scout" {
		t.Errorf("exact name match should be first, got %q", results[0].Name)
	}
}

func TestSearch_MultipleKeywords(t *testing.T) {
	loader, err := NewLoader("")
	if err != nil {
		t.Fatal(err)
	}
	results := loader.Search("write specification release", 5)
	if len(results) == 0 {
		t.Fatal("Search for multiple keywords returned no results")
	}
}

func TestCalculateSkillRelevance_HighForExactName(t *testing.T) {
	skill := Skill{
		Name:        "debugging",
		Description: "Systematic debugging methodology",
		Content:     "Step 1: Reproduce the issue.",
		Triggers:    []string{"debug", "troubleshoot"},
	}
	score := calculateSkillRelevance("debugging", []string{"debugging"}, skill)
	if score < 1.0 {
		t.Errorf("exact name match score = %.2f, want >= 1.0", score)
	}
}

func TestCalculateSkillRelevance_ZeroForNoMatch(t *testing.T) {
	skill := Skill{
		Name:        "debugging",
		Description: "Systematic debugging",
		Content:     "Step 1.",
	}
	score := calculateSkillRelevance("xyzzyplugh", []string{"xyzzyplugh"}, skill)
	if score != 0 {
		t.Errorf("no match score = %.2f, want 0", score)
	}
}
