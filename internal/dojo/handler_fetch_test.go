package dojo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newFetchHandler(t *testing.T, hub Hub) (*Handler, string) {
	t.Helper()
	adrDir := t.TempDir()
	h, err := NewHandler("", adrDir, nil, "")
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	if hub != nil {
		h.SetHub(hub)
	}
	return h, adrDir
}

func TestFetch_RequiresQueryOrID(t *testing.T) {
	h, _ := newFetchHandler(t, nil)
	res := callTool(t, h.handleFetch, map[string]any{})
	if !res.IsError {
		t.Fatal("want error when neither query nor id given")
	}
}

func TestFetch_QueryMergesStores(t *testing.T) {
	adrName := "2026-07-05_use-postgres-for-memory.md"
	h, adrDir := newFetchHandler(t, &fakeHub{entries: hubFixtures})
	if err := os.WriteFile(filepath.Join(adrDir, adrName), []byte("# ADR: memory decision\n\npostgres"), 0644); err != nil {
		t.Fatalf("write adr: %v", err)
	}

	res := callTool(t, h.handleFetch, map[string]any{"query": "memory postgres", "limit": 15})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(text, "adr:"+adrName) {
		t.Errorf("ADR filename hit missing:\n%s", text)
	}
	if !strings.Contains(text, "dojo store") {
		t.Errorf("header missing:\n%s", text)
	}
}

func TestFetch_MemoryStoreSkippedWhenNoHub(t *testing.T) {
	h, _ := newFetchHandler(t, nil)
	res := callTool(t, h.handleFetch, map[string]any{"query": "debugging"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(text, "memory store skipped") {
		t.Errorf("want skip note when hub unset:\n%s", text)
	}
	// Bundled skills should still produce hits for "debugging".
	if !strings.Contains(text, "[skill]") {
		t.Errorf("want skill hits:\n%s", text)
	}
}

func TestFetch_ByID_Skill(t *testing.T) {
	h, _ := newFetchHandler(t, nil)
	// Find a bundled skill name deterministically.
	all := h.skillsLoader.AllSkills()
	if len(all) == 0 {
		t.Skip("no bundled skills")
	}
	name := all[0].Name

	res := callTool(t, h.handleFetch, map[string]any{"id": "skill:" + name})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	if !strings.Contains(resultText(t, res), "# Skill: "+name) {
		t.Fatalf("full skill content missing:\n%s", truncateStr(resultText(t, res), 200))
	}
}

func TestFetch_ByID_MemoryAndSeed(t *testing.T) {
	h, _ := newFetchHandler(t, &fakeHub{entries: hubFixtures})

	res := callTool(t, h.handleFetch, map[string]any{"id": "memory:reference-dojo-bridge"})
	if res.IsError {
		t.Fatalf("memory fetch error: %s", resultText(t, res))
	}
	if !strings.Contains(resultText(t, res), "bridge body") {
		t.Fatalf("memory body missing: %s", resultText(t, res))
	}

	seeds := h.wisdomBase.ListSeeds()
	if len(seeds) == 0 {
		t.Skip("no bundled seeds")
	}
	res2 := callTool(t, h.handleFetch, map[string]any{"id": "seed:" + seeds[0].Name})
	if res2.IsError {
		t.Fatalf("seed fetch error: %s", resultText(t, res2))
	}
}

func TestFetch_ByID_ADRTraversalBlocked(t *testing.T) {
	h, _ := newFetchHandler(t, nil)
	for _, bad := range []string{"adr:../secrets.md", "adr:..\\up.md", "adr:notmarkdown.txt", "adr:a/b.md"} {
		res := callTool(t, h.handleFetch, map[string]any{"id": bad})
		if !res.IsError {
			t.Errorf("id %q: want error, got success", bad)
		}
	}
}

func TestFetch_ByID_UnknownStore(t *testing.T) {
	h, _ := newFetchHandler(t, nil)
	res := callTool(t, h.handleFetch, map[string]any{"id": "warehouse:thing"})
	if !res.IsError || !strings.Contains(resultText(t, res), "unknown store") {
		t.Fatalf("want unknown-store error, got: %s", resultText(t, res))
	}
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
