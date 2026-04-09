package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skill represents a loaded CoworkPlugins skill.
type Skill struct {
	Name        string   // from YAML frontmatter "name"
	Description string   // from YAML frontmatter "description"
	Plugin      string   // derived from directory path (e.g., "strategic-thinking")
	Content     string   // full markdown body after frontmatter
	Triggers    []string // extracted from description
	FilePath    string   // absolute path to SKILL.md (empty for bundled)
}

// Loader manages skill loading and indexing.
type Loader struct {
	skills   []Skill
	byName   map[string]*Skill
	byPlugin map[string][]*Skill
}

// frontmatter holds the YAML fields we care about.
type frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// NewLoader creates a loader and reads skills from the given path.
// The path should point to a CoworkPlugins root directory with structure:
//
//	plugins/{plugin-name}/skills/{skill-name}/SKILL.md
//
// If path is empty or doesn't exist, falls back to bundled skills.
func NewLoader(skillsPath string) (*Loader, error) {
	l := &Loader{
		byName:   make(map[string]*Skill),
		byPlugin: make(map[string][]*Skill),
	}

	if skillsPath != "" {
		info, err := os.Stat(skillsPath)
		if err == nil && info.IsDir() {
			if loadErr := l.loadFromFilesystem(skillsPath); loadErr != nil {
				return nil, fmt.Errorf("load skills from %s: %w", skillsPath, loadErr)
			}
			if len(l.skills) > 0 {
				return l, nil
			}
		}
	}

	// Fallback to bundled skills
	if err := l.loadFromBundled(); err != nil {
		return nil, fmt.Errorf("load bundled skills: %w", err)
	}
	return l, nil
}

// loadFromFilesystem walks the CoworkPlugins directory structure.
func (l *Loader) loadFromFilesystem(root string) error {
	var skipped int
	pluginsDir := filepath.Join(root, "plugins")
	info, err := os.Stat(pluginsDir)
	if err != nil || !info.IsDir() {
		// If there's no plugins/ subdir, try root directly as plugins dir
		pluginsDir = root
	}

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("read plugins directory: %w", err)
	}

	for _, pluginEntry := range entries {
		if !pluginEntry.IsDir() {
			continue
		}
		pluginName := pluginEntry.Name()
		skillsDir := filepath.Join(pluginsDir, pluginName, "skills")

		skillEntries, err := os.ReadDir(skillsDir)
		if err != nil {
			continue // plugin might not have skills dir
		}

		for _, skillEntry := range skillEntries {
			if !skillEntry.IsDir() {
				continue
			}
			skillFile := filepath.Join(skillsDir, skillEntry.Name(), "SKILL.md")
			if _, err := os.Stat(skillFile); err != nil {
				continue
			}

			data, err := os.ReadFile(skillFile)
			if err != nil {
				continue
			}

			skill, err := parseSkillFile(string(data), pluginName, skillFile)
			if err != nil {
				skipped++
				fmt.Fprintf(os.Stderr, "WARN: skills loader: skipping %s: %v\n", skillFile, err)
				continue // skip files that fail to parse
			}

			l.addSkill(skill)
		}
	}

	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "WARN: skills loader: skipped %d skills due to parse errors\n", skipped)
	}

	return nil
}

// addSkill indexes a skill by name and plugin.
func (l *Loader) addSkill(s Skill) {
	l.skills = append(l.skills, s)
	idx := len(l.skills) - 1
	l.byName[s.Name] = &l.skills[idx]
	l.byPlugin[s.Plugin] = append(l.byPlugin[s.Plugin], &l.skills[idx])
}

// parseSkillFile parses a SKILL.md file into a Skill struct.
func parseSkillFile(content, plugin, filePath string) (Skill, error) {
	fm, body, err := parseFrontmatter(content)
	if err != nil {
		return Skill{}, err
	}
	if fm.Name == "" {
		return Skill{}, fmt.Errorf("SKILL.md missing 'name' in frontmatter")
	}

	triggers := extractTriggers(fm.Description)

	return Skill{
		Name:        fm.Name,
		Description: fm.Description,
		Plugin:      plugin,
		Content:     body,
		Triggers:    triggers,
		FilePath:    filePath,
	}, nil
}

// parseFrontmatter splits a SKILL.md file into YAML frontmatter and body.
// It handles the common case where SKILL.md description fields contain
// unquoted colons, which standard YAML would reject.
func parseFrontmatter(content string) (frontmatter, string, error) {
	var fm frontmatter

	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return fm, content, fmt.Errorf("no frontmatter delimiter found")
	}

	// Find the second ---
	rest := content[3:]
	// Skip leading newline after first ---
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	endIdx := strings.Index(rest, "\n---")
	if endIdx == -1 {
		// Handle case where --- is at the very start (empty frontmatter)
		if strings.HasPrefix(rest, "---") {
			return fm, strings.TrimLeft(rest[3:], "\r\n"), nil
		}
		return fm, content, fmt.Errorf("no closing frontmatter delimiter found")
	}

	yamlContent := rest[:endIdx]
	body := rest[endIdx+4:] // skip \n---
	// Trim leading newline from body
	body = strings.TrimLeft(body, "\r\n")

	// Try standard YAML parsing first
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		// Fall back to line-by-line parsing for files with unquoted colons
		fm = parseFrontmatterLines(yamlContent)
	}

	return fm, body, nil
}

// parseFrontmatterLines is a fallback parser that handles YAML-like frontmatter
// where values contain unquoted colons (common in CoworkPlugins SKILL.md files).
func parseFrontmatterLines(content string) frontmatter {
	var fm frontmatter
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Find the first colon-space separator
		idx := strings.Index(line, ": ")
		if idx == -1 {
			// Check for colon at end of line (key with empty value)
			if strings.HasSuffix(line, ":") {
				continue
			}
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+2:])
		switch key {
		case "name":
			fm.Name = value
		case "description":
			fm.Description = value
		}
	}
	return fm
}

// extractTriggers parses the description field for trigger phrases.
func extractTriggers(description string) []string {
	var triggers []string
	lower := strings.ToLower(description)

	// Look for trigger phrase patterns
	patterns := []string{
		"trigger phrases:",
		"trigger:",
		"use when ",
		"use this when ",
	}

	for _, pattern := range patterns {
		idx := strings.Index(lower, pattern)
		if idx == -1 {
			continue
		}

		// Extract text after the pattern
		after := description[idx+len(pattern):]

		// For "trigger phrases:" pattern, look for quoted phrases
		if strings.HasPrefix(pattern, "trigger") {
			// Split on comma or single-quote delimited phrases
			parts := strings.Split(after, ",")
			for _, part := range parts {
				part = strings.Trim(part, " \t\n\r'\".")
				if part != "" && len(part) > 3 {
					triggers = append(triggers, strings.ToLower(part))
				}
			}
		} else {
			// For "use when" patterns, extract the clause
			// Find the end of the sentence (period or end of string)
			endIdx := strings.IndexAny(after, ".\n")
			if endIdx == -1 {
				endIdx = len(after)
			}
			clause := strings.TrimSpace(after[:endIdx])
			// Split on commas or "or"
			parts := strings.Split(clause, ",")
			for _, part := range parts {
				subParts := strings.Split(part, " or ")
				for _, sub := range subParts {
					sub = strings.TrimSpace(sub)
					if sub != "" && len(sub) > 3 {
						triggers = append(triggers, strings.ToLower(sub))
					}
				}
			}
		}
	}

	return triggers
}

// GetByName returns a skill by its exact name.
func (l *Loader) GetByName(name string) (*Skill, error) {
	s, ok := l.byName[name]
	if !ok {
		return nil, fmt.Errorf("skill not found: %s", name)
	}
	return s, nil
}

// ListByPlugin returns all skills grouped by plugin name.
func (l *Loader) ListByPlugin() map[string][]*Skill {
	return l.byPlugin
}

// AllSkills returns a copy of all loaded skills.
func (l *Loader) AllSkills() []Skill {
	return l.skills
}

// Count returns the total number of loaded skills.
func (l *Loader) Count() int {
	return len(l.skills)
}

// PluginNames returns all plugin names in sorted order.
func (l *Loader) PluginNames() []string {
	names := make([]string, 0, len(l.byPlugin))
	for name := range l.byPlugin {
		names = append(names, name)
	}
	// Simple sort
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return names
}
