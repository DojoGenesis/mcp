package wisdom

import (
	"testing"
)

func TestGetSeeds_Count(t *testing.T) {
	seeds := getSeeds()
	if len(seeds) != 20 {
		t.Errorf("Expected 20 seeds, got %d", len(seeds))
	}
}

func TestGetSeeds_AllHaveRequiredFields(t *testing.T) {
	seeds := getSeeds()
	for i, seed := range seeds {
		if seed.Name == "" {
			t.Errorf("Seed at index %d has empty Name", i)
		}
		if seed.Description == "" {
			t.Errorf("Seed %q has empty Description", seed.Name)
		}
		if seed.Content == "" {
			t.Errorf("Seed %q has empty Content", seed.Name)
		}
		if seed.Category == "" {
			t.Errorf("Seed %q has empty Category", seed.Name)
		}
		if seed.Triggers == "" {
			t.Errorf("Seed %q has empty Triggers", seed.Name)
		}
	}
}

func TestGetSeeds_UniqueNames(t *testing.T) {
	seeds := getSeeds()
	seen := make(map[string]bool)
	for _, seed := range seeds {
		if seen[seed.Name] {
			t.Errorf("Duplicate seed name: %q", seed.Name)
		}
		seen[seed.Name] = true
	}
}

func TestGetSeeds_KnownSeedsPresent(t *testing.T) {
	seeds := getSeeds()
	expected := []string{
		"three_tiered_governance",
		"harness_trace",
		"context_iceberg",
		"agent_connect",
		"go_live_bundles",
		"cost_guard",
		"safety_switch",
		"implicit_perspective_extraction",
		"mode_based_complexity_gating",
		"shared_infrastructure",
		"sanctuary_architecture",
		"pace_of_understanding",
		"lineage_transmission",
		"graceful_failure",
		"local_first_liberation",
		"the_onsen_pattern",
		"collaborative_calibration",
		"transparent_intelligence",
		"inter_acceptance",
		"radical_freedom",
	}

	nameSet := make(map[string]bool)
	for _, s := range seeds {
		nameSet[s.Name] = true
	}

	for _, name := range expected {
		if !nameSet[name] {
			t.Errorf("Expected seed %q not found", name)
		}
	}
}

func TestGetSeeds_CategoriesValid(t *testing.T) {
	seeds := getSeeds()
	validCategories := map[string]bool{
		"dojo_genesis":    true,
		"aroma":           true,
		"aroma_serenity":  true,
		"serenity_valley": true,
	}

	for _, seed := range seeds {
		if !validCategories[seed.Category] {
			t.Errorf("Seed %q has unexpected category %q", seed.Name, seed.Category)
		}
	}
}

func TestGetSeed_ByName_AllAccessible(t *testing.T) {
	b := NewBase()
	seeds := b.ListSeeds()

	for _, seed := range seeds {
		retrieved, err := b.GetSeed(seed.Name)
		if err != nil {
			t.Errorf("GetSeed(%q) returned error: %v", seed.Name, err)
			continue
		}
		if retrieved == nil {
			t.Errorf("GetSeed(%q) returned nil", seed.Name)
			continue
		}
		if retrieved.Name != seed.Name {
			t.Errorf("GetSeed(%q) returned seed with name %q", seed.Name, retrieved.Name)
		}
	}
}
