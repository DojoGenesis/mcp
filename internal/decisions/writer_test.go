package decisions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWriter_DefaultPath(t *testing.T) {
	tmpDir := t.TempDir()
	defaultPath := filepath.Join(tmpDir, "decisions")

	w, err := NewWriter(defaultPath)
	if err != nil {
		t.Fatalf("NewWriter returned error: %v", err)
	}
	if w.BasePath() != defaultPath {
		t.Errorf("BasePath = %q, want %q", w.BasePath(), defaultPath)
	}

	// Directory should exist
	info, err := os.Stat(defaultPath)
	if err != nil {
		t.Fatalf("decision directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("decision path is not a directory")
	}
}

func TestNewWriter_EmptyPath(t *testing.T) {
	// Save and restore working directory
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	w, err := NewWriter("")
	if err != nil {
		t.Fatalf("NewWriter with empty path returned error: %v", err)
	}
	if w.BasePath() != "./decisions" {
		t.Errorf("BasePath = %q, want %q", w.BasePath(), "./decisions")
	}
}

func TestLogDecision_WritesFile(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := NewWriter(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	fp, err := w.LogDecision(
		"Use PostgreSQL for persistence",
		"We need a database for storing user data.",
		"We will use PostgreSQL because it is reliable and well-supported.",
		"We need to add PostgreSQL to our deployment stack.",
	)
	if err != nil {
		t.Fatalf("LogDecision returned error: %v", err)
	}
	if fp == "" {
		t.Fatal("LogDecision returned empty filepath")
	}

	// File should exist
	data, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("cannot read ADR file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# ADR: Use PostgreSQL for persistence") {
		t.Error("ADR does not contain title")
	}
	if !strings.Contains(content, "## Context") {
		t.Error("ADR does not contain Context section")
	}
	if !strings.Contains(content, "## Decision") {
		t.Error("ADR does not contain Decision section")
	}
	if !strings.Contains(content, "## Consequences") {
		t.Error("ADR does not contain Consequences section")
	}
	if !strings.Contains(content, "Accepted") {
		t.Error("ADR does not contain status")
	}
}

func TestLogDecision_FilenameSafe(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := NewWriter(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	fp, err := w.LogDecision(
		"Use Go + YAML for Config!!! @#$%",
		"ctx", "dec", "cons",
	)
	if err != nil {
		t.Fatal(err)
	}

	filename := filepath.Base(fp)
	// Should not contain special characters
	if strings.ContainsAny(filename, "!@#$%^&*()") {
		t.Errorf("filename contains special characters: %q", filename)
	}
	// Should end with .md
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("filename does not end with .md: %q", filename)
	}
}

func TestSlugify_Basic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Use PostgreSQL", "use-postgresql"},
		{"Special!@#$Characters", "special-characters"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"", "untitled"},
		{"a--b", "a-b"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSlugify_MaxLength(t *testing.T) {
	long := strings.Repeat("a", 100)
	result := slugify(long)
	if len(result) > 50 {
		t.Errorf("slugify of long string produced %d chars, want <= 50", len(result))
	}
}

func TestSlugify_NoTrailingHyphen(t *testing.T) {
	// This title when truncated to 50 chars could end on a hyphen
	input := strings.Repeat("word-", 20)
	result := slugify(input)
	if strings.HasSuffix(result, "-") {
		t.Errorf("slugify produced trailing hyphen: %q", result)
	}
}

func TestFormatADR_Structure(t *testing.T) {
	result := formatADR("Test Title", "Test Context", "Test Decision", "Test Consequences", "2026-04-05")

	if !strings.HasPrefix(result, "# ADR: Test Title") {
		t.Error("ADR does not start with title")
	}
	if !strings.Contains(result, "**Date:** 2026-04-05") {
		t.Error("ADR does not contain date")
	}
	if !strings.Contains(result, "**Status:** Accepted") {
		t.Error("ADR does not contain status")
	}
	if !strings.Contains(result, "## Context\n\nTest Context") {
		t.Error("ADR does not contain context section with content")
	}
	if !strings.Contains(result, "## Decision\n\nTest Decision") {
		t.Error("ADR does not contain decision section with content")
	}
	if !strings.Contains(result, "## Consequences\n\nTest Consequences") {
		t.Error("ADR does not contain consequences section with content")
	}
}

func TestLogDecision_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	w, err := NewWriter(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	fp1, err := w.LogDecision("Decision One", "ctx", "dec", "cons")
	if err != nil {
		t.Fatal(err)
	}
	fp2, err := w.LogDecision("Decision Two", "ctx", "dec", "cons")
	if err != nil {
		t.Fatal(err)
	}

	if fp1 == fp2 {
		t.Error("two decisions produced the same filepath")
	}

	// Both files should exist
	for _, fp := range []string{fp1, fp2} {
		if _, err := os.Stat(fp); err != nil {
			t.Errorf("ADR file does not exist: %s", fp)
		}
	}
}
