package decisions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Writer handles writing ADR files to disk.
type Writer struct {
	basePath string
}

// NewWriter creates a writer. Creates basePath directory if it doesn't exist.
func NewWriter(basePath string) (*Writer, error) {
	if basePath == "" {
		basePath = "./decisions"
	}
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create ADR directory: %w", err)
	}
	return &Writer{basePath: basePath}, nil
}

// LogDecision writes an ADR file and returns the filepath.
func (w *Writer) LogDecision(title, context, decision, consequences string) (string, error) {
	date := time.Now().Format("2006-01-02")
	slug := slugify(title)
	filename := fmt.Sprintf("%s_%s.md", date, slug)
	fp := filepath.Join(w.basePath, filename)

	content := formatADR(title, context, decision, consequences, date)

	if err := os.WriteFile(fp, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write ADR: %w", err)
	}
	return fp, nil
}

// BasePath returns the configured ADR output directory.
func (w *Writer) BasePath() string {
	return w.basePath
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// slugify converts a title into a safe filename slug.
func slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 50 {
		s = s[:50]
		// Don't end on a hyphen
		s = strings.TrimRight(s, "-")
	}
	if s == "" {
		s = "untitled"
	}
	return s
}

// formatADR produces the markdown content for an ADR file.
func formatADR(title, context, decision, consequences, date string) string {
	return fmt.Sprintf(`# ADR: %s

**Date:** %s
**Status:** Accepted

## Context

%s

## Decision

%s

## Consequences

%s
`, title, date, context, decision, consequences)
}
