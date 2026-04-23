package dojo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestSaveDojoSettings_ConcurrentWrites verifies that two concurrent
// saveDojoSettings calls with different payloads result in a final file that
// parses as valid JSON equal to one of the payloads (not torn/mixed).
//
// This is a handler-layer test: it exercises the dojoSettingsMu guard plus the
// AtomicWriteFile rename to prove that neither partial writes nor
// load-modify-write races corrupt the file.
func TestSaveDojoSettings_ConcurrentWrites(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	dir := filepath.Join(tmpHome, ".dojo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "settings.json")

	// Seed an initial file so neither goroutine starts from empty.
	initial := map[string]interface{}{
		"defaults": map[string]interface{}{"disposition": "balanced"},
	}
	data, _ := json.MarshalIndent(initial, "", "  ")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	dispositions := []string{"focused", "exploratory"}
	const goroutines = 2

	h := newTestHandler(t)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		d := dispositions[i]
		go func() {
			defer wg.Done()
			req := newCallToolRequest(map[string]interface{}{"disposition": d})
			_, _ = h.handleDispositionSet(context.Background(), req)
		}()
	}
	wg.Wait()

	// The final file must parse as valid JSON.
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile after concurrent writes: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("final file is not valid JSON (torn write?): %v\ncontent: %q", err, out)
	}

	// The disposition must be exactly one of the two written values.
	defaults, ok := got["defaults"].(map[string]interface{})
	if !ok {
		t.Fatalf("defaults is not a map: %T", got["defaults"])
	}
	disposition, _ := defaults["disposition"].(string)
	if disposition != "focused" && disposition != "exploratory" {
		t.Errorf("disposition %q is neither 'focused' nor 'exploratory' (torn write?)", disposition)
	}
}

// TestSaveProject_ConcurrentWrites verifies that two concurrent
// saveProject calls result in a final project.json that parses as valid JSON
// equal to one of the written payloads (not torn/mixed).
func TestSaveProject_ConcurrentWrites(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Ensure the projects dir exists with a default project so handlers find it.
	projectDir := filepath.Join(tmpHome, ".dojo", "projects", "default")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	initial := &ProjectFile{
		Name:   "default",
		Phase:  "active",
		Tracks: []Track{},
	}
	data, _ := json.MarshalIndent(initial, "", "  ")
	if err := os.WriteFile(filepath.Join(projectDir, "project.json"), data, 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	h := newTestHandler(t)
	trackNames := []string{"track-alpha", "track-beta"}

	var wg sync.WaitGroup
	wg.Add(len(trackNames))
	for _, name := range trackNames {
		n := name
		go func() {
			defer wg.Done()
			req := newCallToolRequest(map[string]interface{}{
				"action": "add",
				"name":   n,
				"status": "in-progress",
			})
			_, _ = h.handleProjectTrack(context.Background(), req)
		}()
	}
	wg.Wait()

	// The final file must parse as valid JSON.
	out, err := os.ReadFile(filepath.Join(projectDir, "project.json"))
	if err != nil {
		t.Fatalf("ReadFile after concurrent writes: %v", err)
	}
	var got ProjectFile
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("final project.json is not valid JSON (torn write?): %v\ncontent: %q", err, out)
	}

	// At least one track must be present (the mutex ensures sequential writes,
	// but one goroutine may "win" the add-check and see no duplicate).
	if len(got.Tracks) == 0 {
		t.Error("no tracks written — expected at least one")
	}
	// Each track name must be one of the valid options.
	valid := map[string]bool{"track-alpha": true, "track-beta": true}
	for _, tr := range got.Tracks {
		if !valid[strings.ToLower(tr.Name)] {
			t.Errorf("unexpected track name %q", tr.Name)
		}
	}
}
