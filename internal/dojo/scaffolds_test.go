package dojo

import (
	"strings"
	"testing"

	"github.com/DojoGenesis/mcp/internal/skills"
	"github.com/DojoGenesis/mcp/internal/wisdom"
)

// TestReflectScaffold locks the observable rendering contract of
// reflectScaffold: header, conditional sections, the cap-at-3 truncation,
// the "First step" conditional, and the both-empty isolation invariant.
func TestReflectScaffold(t *testing.T) {
	skillWithStep := skills.Skill{
		Name:        "Debugging Protocol",
		Description: "A method for killing hypotheses instead of trying fixes. It also covers config-vs-code triage.",
		Plugin:      "continuous-learning",
		Content:     "Some intro text.\n1. Form a hypothesis about the mechanism.\n2. Find the cheapest disproving observation.",
	}
	skillNoStep := skills.Skill{
		Name:        "No Step Skill",
		Description: "A skill with no numbered steps in its content.",
		Plugin:      "some-plugin",
		Content:     "Just prose with no numbered list anywhere in it.",
	}
	seedA := wisdom.Seed{
		Name:     "Seed A",
		Content:  "Config errors present as code errors. Always check settings first.",
		Category: "engineering",
	}

	t.Run("header always contains session description", func(t *testing.T) {
		got := reflectScaffold("my weird session", nil, nil)
		if !strings.Contains(got, "# Reflection: my weird session") {
			t.Errorf("expected header with session description, got:\n%s", got)
		}
	})

	t.Run("skills section absent when no skills matched", func(t *testing.T) {
		got := reflectScaffold("desc", nil, []wisdom.Seed{seedA})
		if strings.Contains(got, "## Relevant Skills from Your Methodology Library") {
			t.Errorf("did not expect skills heading with zero matched skills, got:\n%s", got)
		}
	})

	t.Run("skills section renders name, plugin, and When to use", func(t *testing.T) {
		got := reflectScaffold("desc", []skills.Skill{skillWithStep}, nil)
		if !strings.Contains(got, "## Relevant Skills from Your Methodology Library") {
			t.Errorf("expected skills heading, got:\n%s", got)
		}
		if !strings.Contains(got, "Debugging Protocol") {
			t.Errorf("expected skill Name to render, got:\n%s", got)
		}
		if !strings.Contains(got, "continuous-learning") {
			t.Errorf("expected skill Plugin to render, got:\n%s", got)
		}
		if !strings.Contains(got, "**When to use:** A method for killing hypotheses instead of trying fixes.") {
			t.Errorf("expected 'When to use' with first sentence of Description, got:\n%s", got)
		}
	})

	t.Run("skills cap at 3 -- 4th matched skill is absent", func(t *testing.T) {
		fourSkills := []skills.Skill{
			{Name: "Skill One", Description: "First.", Plugin: "p1"},
			{Name: "Skill Two", Description: "Second.", Plugin: "p2"},
			{Name: "Skill Three", Description: "Third.", Plugin: "p3"},
			{Name: "Skill Four", Description: "Fourth.", Plugin: "p4"},
		}
		got := reflectScaffold("desc", fourSkills, nil)
		for _, name := range []string{"Skill One", "Skill Two", "Skill Three"} {
			if !strings.Contains(got, name) {
				t.Errorf("expected %q to render (within cap of 3), got:\n%s", name, got)
			}
		}
		if strings.Contains(got, "Skill Four") {
			t.Errorf("expected 4th matched skill to be absent due to cap-at-3, got:\n%s", got)
		}
	})

	t.Run("seeds cap at 3 -- 4th matched seed is absent", func(t *testing.T) {
		fourSeeds := []wisdom.Seed{
			{Name: "Seed One", Content: "First insight."},
			{Name: "Seed Two", Content: "Second insight."},
			{Name: "Seed Three", Content: "Third insight."},
			{Name: "Seed Four", Content: "Fourth insight."},
		}
		got := reflectScaffold("desc", nil, fourSeeds)
		for _, name := range []string{"Seed One", "Seed Two", "Seed Three"} {
			if !strings.Contains(got, name) {
				t.Errorf("expected %q to render (within cap of 3), got:\n%s", name, got)
			}
		}
		if strings.Contains(got, "Seed Four") {
			t.Errorf("expected 4th matched seed to be absent due to cap-at-3, got:\n%s", got)
		}
	})

	t.Run("First step renders when extractFirstStep finds a numbered step", func(t *testing.T) {
		got := reflectScaffold("desc", []skills.Skill{skillWithStep}, nil)
		if !strings.Contains(got, "**First step:**") {
			t.Errorf("expected 'First step' line when content has a numbered step, got:\n%s", got)
		}
		if !strings.Contains(got, "Form a hypothesis about the mechanism.") {
			t.Errorf("expected extracted first step text to render, got:\n%s", got)
		}
	})

	t.Run("First step absent when extractFirstStep finds nothing", func(t *testing.T) {
		got := reflectScaffold("desc", []skills.Skill{skillNoStep}, nil)
		if strings.Contains(got, "**First step:**") {
			t.Errorf("did not expect 'First step' line when content has no numbered step, got:\n%s", got)
		}
	})

	t.Run("seeds section absent when no seeds matched", func(t *testing.T) {
		got := reflectScaffold("desc", []skills.Skill{skillWithStep}, nil)
		if strings.Contains(got, "## Relevant Seeds (Reusable Patterns)") {
			t.Errorf("did not expect seeds heading with zero matched seeds, got:\n%s", got)
		}
	})

	t.Run("seeds section renders name and Core insight", func(t *testing.T) {
		got := reflectScaffold("desc", nil, []wisdom.Seed{seedA})
		if !strings.Contains(got, "## Relevant Seeds (Reusable Patterns)") {
			t.Errorf("expected seeds heading, got:\n%s", got)
		}
		if !strings.Contains(got, "Seed A") {
			t.Errorf("expected seed Name to render, got:\n%s", got)
		}
		if !strings.Contains(got, "**Core insight:** Config errors present as code errors.") {
			t.Errorf("expected 'Core insight' with first sentence of Content, got:\n%s", got)
		}
	})

	t.Run("no-matches text renders only when both slices are empty", func(t *testing.T) {
		bothEmpty := reflectScaffold("desc", nil, nil)
		if !strings.Contains(bothEmpty, "No direct matches found in the methodology library") {
			t.Errorf("expected no-matches text when both skills and seeds are empty, got:\n%s", bothEmpty)
		}

		skillsOnly := reflectScaffold("desc", []skills.Skill{skillWithStep}, nil)
		if strings.Contains(skillsOnly, "No direct matches found in the methodology library") {
			t.Errorf("did not expect no-matches text when skills matched but seeds empty, got:\n%s", skillsOnly)
		}

		seedsOnly := reflectScaffold("desc", nil, []wisdom.Seed{seedA})
		if strings.Contains(seedsOnly, "No direct matches found in the methodology library") {
			t.Errorf("did not expect no-matches text when seeds matched but skills empty, got:\n%s", seedsOnly)
		}
	})

	t.Run("Reflection Questions block always present", func(t *testing.T) {
		cases := []struct {
			name   string
			skills []skills.Skill
			seeds  []wisdom.Seed
		}{
			{"both empty", nil, nil},
			{"skills only", []skills.Skill{skillWithStep}, nil},
			{"seeds only", nil, []wisdom.Seed{seedA}},
			{"both present", []skills.Skill{skillWithStep}, []wisdom.Seed{seedA}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := reflectScaffold("desc", tc.skills, tc.seeds)
				if !strings.Contains(got, "## Reflection Questions") {
					t.Errorf("expected Reflection Questions heading, got:\n%s", got)
				}
				if !strings.Contains(got, "dojo.log_decision") {
					t.Errorf("expected dojo.log_decision closing line, got:\n%s", got)
				}
			})
		}
	})
}
