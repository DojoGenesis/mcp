package wisdom

import (
	"strings"
	"testing"
)

func TestNewBase(t *testing.T) {
	b := NewBase()
	if b == nil {
		t.Fatal("NewBase returned nil")
	}
	if len(b.seeds) == 0 {
		t.Error("NewBase created base with no seeds")
	}
	if len(b.resources) == 0 {
		t.Error("NewBase created base with no resources")
	}
	if b.principles == "" {
		t.Error("NewBase created base with empty principles")
	}
}

func TestSearch_ExactMatch(t *testing.T) {
	b := NewBase()
	results := b.Search("three_tiered_governance")
	if len(results) == 0 {
		t.Fatal("Search for known seed name returned no results")
	}
	found := false
	for _, r := range results {
		if r.Name == "three_tiered_governance" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Search results did not contain the expected seed 'three_tiered_governance'")
	}
}

func TestSearch_NoMatch(t *testing.T) {
	b := NewBase()
	results := b.Search("xyzzyplughqwerty12345")
	if len(results) != 0 {
		t.Errorf("Search for gibberish returned %d results, expected 0", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	b := NewBase()
	lower := b.Search("governance")
	upper := b.Search("GOVERNANCE")
	mixed := b.Search("Governance")

	if len(lower) == 0 {
		t.Fatal("Search for 'governance' returned no results")
	}
	if len(lower) != len(upper) {
		t.Errorf("Case sensitivity: 'governance' returned %d results, 'GOVERNANCE' returned %d", len(lower), len(upper))
	}
	if len(lower) != len(mixed) {
		t.Errorf("Case sensitivity: 'governance' returned %d results, 'Governance' returned %d", len(lower), len(mixed))
	}
}

func TestGetSeed_Found(t *testing.T) {
	b := NewBase()
	seed, err := b.GetSeed("three_tiered_governance")
	if err != nil {
		t.Fatalf("GetSeed returned error for existing seed: %v", err)
	}
	if seed == nil {
		t.Fatal("GetSeed returned nil seed for existing name")
	}
	if seed.Name != "three_tiered_governance" {
		t.Errorf("GetSeed returned seed with name %q, want %q", seed.Name, "three_tiered_governance")
	}
	if seed.Content == "" {
		t.Error("GetSeed returned seed with empty content")
	}
	if seed.Description == "" {
		t.Error("GetSeed returned seed with empty description")
	}
}

func TestGetSeed_NotFound(t *testing.T) {
	b := NewBase()
	seed, err := b.GetSeed("nonexistent_seed_name")
	if err == nil {
		t.Error("GetSeed did not return error for nonexistent seed")
	}
	if seed != nil {
		t.Error("GetSeed returned non-nil seed for nonexistent name")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error message %q does not contain 'not found'", err.Error())
	}
}

func TestListSeeds(t *testing.T) {
	b := NewBase()
	seeds := b.ListSeeds()
	if len(seeds) == 0 {
		t.Fatal("ListSeeds returned empty slice")
	}
	// Verify each seed has required fields
	for i, seed := range seeds {
		if seed.Name == "" {
			t.Errorf("Seed at index %d has empty name", i)
		}
		if seed.Description == "" {
			t.Errorf("Seed %q has empty description", seed.Name)
		}
		if seed.Content == "" {
			t.Errorf("Seed %q has empty content", seed.Name)
		}
	}
}

func TestGetPrinciples(t *testing.T) {
	b := NewBase()
	principles := b.GetPrinciples()
	if principles == "" {
		t.Fatal("GetPrinciples returned empty string")
	}
	if !strings.Contains(principles, "Beginner") {
		t.Error("Principles do not mention 'Beginner' (expected Beginner's Mind)")
	}
	if !strings.Contains(principles, "Self-Definition") {
		t.Error("Principles do not mention 'Self-Definition'")
	}
	if !strings.Contains(principles, "Understanding is Love") {
		t.Error("Principles do not mention 'Understanding is Love'")
	}
}

func TestGetResource_Found(t *testing.T) {
	b := NewBase()
	resources := b.ListResources()
	if len(resources) == 0 {
		t.Fatal("ListResources returned empty; cannot test GetResource")
	}
	name := resources[0].Name
	content, err := b.GetResource(name)
	if err != nil {
		t.Fatalf("GetResource(%q) returned error: %v", name, err)
	}
	if content == "" {
		t.Errorf("GetResource(%q) returned empty content", name)
	}
}

func TestGetResource_NotFound(t *testing.T) {
	b := NewBase()
	content, err := b.GetResource("nonexistent_resource_xyz")
	if err == nil {
		t.Error("GetResource did not return error for nonexistent resource")
	}
	if content != "" {
		t.Error("GetResource returned non-empty content for nonexistent resource")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error message %q does not contain 'not found'", err.Error())
	}
}

func TestListResources(t *testing.T) {
	b := NewBase()
	resources := b.ListResources()
	if len(resources) == 0 {
		t.Fatal("ListResources returned empty slice")
	}
	for i, r := range resources {
		if r.Name == "" {
			t.Errorf("Resource at index %d has empty name", i)
		}
		if r.Content == "" {
			t.Errorf("Resource %q has empty content", r.Name)
		}
	}
}

func TestCalculateRelevance_NameMatch(t *testing.T) {
	score := calculateRelevance("governance", "three_tiered_governance", "A governance framework", "Content about governance")
	if score <= 0 {
		t.Error("calculateRelevance returned zero or negative score for name match")
	}
	// Name match should give the highest base score (1.0)
	nameOnly := calculateRelevance("governance", "governance", "", "")
	descOnly := calculateRelevance("governance", "", "governance", "")
	if nameOnly <= descOnly {
		t.Errorf("Name match score (%.2f) should be higher than description-only match (%.2f)", nameOnly, descOnly)
	}
}

func TestCalculateRelevance_NoMatch(t *testing.T) {
	score := calculateRelevance("xyzzyplugh", "governance", "A framework", "Content about patterns")
	if score != 0 {
		t.Errorf("calculateRelevance returned %.2f for no match, expected 0", score)
	}
}

func TestGetSnippet_Found(t *testing.T) {
	content := "This is a long piece of content that talks about governance frameworks and how they apply to various systems in a meaningful way."
	snippet := getSnippet(content, "governance")
	if !strings.Contains(snippet, "governance") {
		t.Error("getSnippet did not return snippet containing the query term")
	}
}

func TestGetSnippet_NotFound(t *testing.T) {
	content := "This is content about architecture and design patterns for building reliable systems."
	snippet := getSnippet(content, "xyzzynonexistent")
	// When not found, should return beginning of content
	lowerContent := strings.ToLower(content)
	if !strings.HasPrefix(snippet, lowerContent[:20]) {
		// The function lowercases everything, so compare lowercase
		t.Errorf("getSnippet for no-match should start with beginning of content, got %q", snippet[:40])
	}
}

func TestGetSnippet_ShortContent(t *testing.T) {
	content := "Short text."
	snippet := getSnippet(content, "nonexistent")
	// Short content should be returned as-is (lowercased by the function)
	if snippet != strings.ToLower(content) {
		t.Errorf("getSnippet for short content returned %q, expected %q", snippet, strings.ToLower(content))
	}
}
