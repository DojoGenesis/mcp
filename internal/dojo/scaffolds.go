package dojo

import (
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp-server/internal/skills"
	"github.com/DojoGenesis/mcp-server/internal/wisdom"
)

// scoutScaffold returns the 4-step strategic analysis scaffold.
// This is the methodology frame -- Claude provides the judgment.
func scoutScaffold(situation string, matchedSkills []skills.Skill) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, `# Strategic Scout: %s

## Step 1: Identify the Tension

Frame this as a tension between two competing ideas. What are the forces pulling in different directions?

**The tension:** [Articulate the core tension here]

## Step 2: Scout Routes

Generate 3-5 distinct approaches. For each:

| Route | Description | Risk | Time | Tradeoff |
|-------|-------------|------|------|----------|
| A | | | | |
| B | | | | |
| C | | | | |

## Step 3: Compare and Synthesize

Which routes can be combined? What's the hybrid that takes the best of each?

**Synthesis:** [Your recommendation here]

## Step 4: Decision

**Selected route:** [Which route and why]
**First action:** [The smallest concrete next step]
**Review point:** [When to check if this is working]`, situation)

	if len(matchedSkills) > 0 {
		sb.WriteString("\n\n---\n\n**Relevant methodology:**\n")
		for _, s := range matchedSkills {
			fmt.Fprintf(&sb, "- **%s** (%s): %s\n", s.Name, s.Plugin, firstSentence(s.Description))
		}
		sb.WriteString("\nUse `dojo.invoke_skill` to load the full workflow for any of these skills.")
	}

	return sb.String()
}

// reflectScaffold returns a structured reflection grounded in matched skills/seeds.
func reflectScaffold(sessionDescription string, matchedSkills []skills.Skill, matchedSeeds []wisdom.Seed) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "# Reflection: %s\n\n", sessionDescription)

	if len(matchedSkills) > 0 {
		sb.WriteString("## Relevant Skills from Your Methodology Library\n\n")
		for i, s := range matchedSkills {
			if i >= 3 {
				break
			}
			fmt.Fprintf(&sb, "### %d. %s (%s)\n", i+1, s.Name, s.Plugin)
			fmt.Fprintf(&sb, "**When to use:** %s\n\n", firstSentence(s.Description))
			// Extract first step from content if possible
			step := extractFirstStep(s.Content)
			if step != "" {
				fmt.Fprintf(&sb, "**First step:** %s\n\n", step)
			}
		}
	}

	if len(matchedSeeds) > 0 {
		sb.WriteString("## Relevant Seeds (Reusable Patterns)\n\n")
		for i, seed := range matchedSeeds {
			if i >= 3 {
				break
			}
			fmt.Fprintf(&sb, "### %d. %s\n", i+1, seed.Name)
			fmt.Fprintf(&sb, "**Core insight:** %s\n\n", firstSentence(seed.Content))
		}
	}

	if len(matchedSkills) == 0 && len(matchedSeeds) == 0 {
		sb.WriteString("No direct matches found in the methodology library for this session description.\n\n")
		sb.WriteString("Try being more specific, or use `dojo.search_skills` to explore available skills.\n\n")
	}

	sb.WriteString(`## Reflection Questions

1. What's the most important thing you learned or decided in this session?
2. What methodology or pattern would have helped you work more effectively?
3. What should you do differently next time?

Use ` + "`dojo.log_decision`" + ` to capture any key decisions made during this session.`)

	return sb.String()
}

// firstSentence returns the first sentence from a string.
func firstSentence(text string) string {
	text = strings.TrimSpace(text)
	// Skip markdown headers
	for strings.HasPrefix(text, "#") {
		idx := strings.Index(text, "\n")
		if idx == -1 {
			return text
		}
		text = strings.TrimSpace(text[idx+1:])
	}
	// Skip bold prefix like "**Core Insight:**"
	if strings.HasPrefix(text, "**") {
		idx := strings.Index(text[2:], "**")
		if idx != -1 {
			after := strings.TrimSpace(text[idx+4:])
			if after != "" {
				text = after
			}
		}
	}

	idx := strings.Index(text, ". ")
	if idx != -1 && idx < 200 {
		return text[:idx+1]
	}
	// Check for period at end of string
	if strings.HasSuffix(text, ".") && len(text) < 200 {
		return text
	}
	if len(text) > 200 {
		return text[:200] + "..."
	}
	return text
}

// extractFirstStep tries to find the first actionable step in skill content.
func extractFirstStep(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Look for numbered steps
		if len(trimmed) > 3 && (strings.HasPrefix(trimmed, "1.") || strings.HasPrefix(trimmed, "1 ")) {
			step := strings.TrimSpace(trimmed[2:])
			step = strings.TrimLeft(step, ". ")
			if step != "" {
				return firstSentence(step)
			}
		}
		// Look for "Step 1" patterns
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "### step 1") || strings.HasPrefix(lower, "## step 1") {
			// Return the header text
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
