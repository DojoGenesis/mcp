package dojo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleDispositionSet_CreatesFileWhenMissing(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	h := newTestHandler(t)
	req := newCallToolRequest(map[string]interface{}{"disposition": "focused"})

	result, err := h.handleDispositionSet(context.Background(), req)
	if err != nil {
		t.Fatalf("handleDispositionSet returned error: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "focused") {
		t.Errorf("result text missing 'focused': %s", text)
	}

	path := filepath.Join(tmpHome, ".dojo", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings.json was not written: %v", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("settings.json invalid: %v", err)
	}
	defaults, ok := settings["defaults"].(map[string]interface{})
	if !ok {
		t.Fatalf("settings.defaults is not a map: %T", settings["defaults"])
	}
	if defaults["disposition"] != "focused" {
		t.Errorf("defaults.disposition = %v, want 'focused'", defaults["disposition"])
	}
}

func TestHandleDispositionSet_PreservesExistingKeys(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Seed an existing settings.json with unrelated keys
	dir := filepath.Join(tmpHome, ".dojo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	seedPath := filepath.Join(dir, "settings.json")
	seed := map[string]interface{}{
		"gateway": map[string]interface{}{
			"url":   "https://gateway.example.com",
			"token": "keep-me",
		},
		"defaults": map[string]interface{}{
			"provider": "anthropic",
			"model":    "claude-opus-4",
		},
	}
	data, _ := json.MarshalIndent(seed, "", "  ")
	if err := os.WriteFile(seedPath, data, 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	h := newTestHandler(t)
	req := newCallToolRequest(map[string]interface{}{"disposition": "exploratory"})
	if _, err := h.handleDispositionSet(context.Background(), req); err != nil {
		t.Fatalf("handleDispositionSet returned error: %v", err)
	}

	out, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read after set: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("parse after set: %v", err)
	}

	gw, ok := got["gateway"].(map[string]interface{})
	if !ok || gw["token"] != "keep-me" {
		t.Errorf("gateway block not preserved: %#v", got["gateway"])
	}
	defaults, ok := got["defaults"].(map[string]interface{})
	if !ok {
		t.Fatalf("defaults not a map: %T", got["defaults"])
	}
	if defaults["provider"] != "anthropic" {
		t.Errorf("defaults.provider not preserved: %v", defaults["provider"])
	}
	if defaults["model"] != "claude-opus-4" {
		t.Errorf("defaults.model not preserved: %v", defaults["model"])
	}
	if defaults["disposition"] != "exploratory" {
		t.Errorf("defaults.disposition = %v, want 'exploratory'", defaults["disposition"])
	}
}

func TestHandleDispositionSet_RejectsEmpty(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	h := newTestHandler(t)
	req := newCallToolRequest(map[string]interface{}{"disposition": ""})
	result, err := h.handleDispositionSet(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatalf("expected error result for empty disposition, got %+v", result)
	}
}
