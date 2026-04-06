package skills

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed all:bundled
var bundledFS embed.FS

// loadFromBundled loads skills from the embedded filesystem.
func (l *Loader) loadFromBundled() error {
	return fs.WalkDir(bundledFS, "bundled", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "SKILL.md" {
			return nil
		}

		data, err := bundledFS.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		// Extract plugin name from path: bundled/{plugin}/{skill}/SKILL.md
		plugin := extractPluginFromEmbedPath(path)

		skill, err := parseSkillFile(string(data), plugin, "")
		if err != nil {
			return nil // skip unparseable files
		}

		l.addSkill(skill)
		return nil
	})
}

// extractPluginFromEmbedPath extracts the plugin name from an embedded path.
// Path format: bundled/{plugin}/{skill}/SKILL.md
func extractPluginFromEmbedPath(path string) string {
	// Normalize separators
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")
	// Expected: ["bundled", "plugin-name", "skill-name", "SKILL.md"]
	if len(parts) >= 3 {
		return parts[1]
	}
	return "unknown"
}
