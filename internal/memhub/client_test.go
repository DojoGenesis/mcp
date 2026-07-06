package memhub

import (
	"context"
	"os"
	"testing"
)

func TestClampLimit(t *testing.T) {
	cases := []struct{ n, def, max, want int }{
		{0, 8, 25, 8},
		{-3, 8, 25, 8},
		{5, 8, 25, 5},
		{100, 8, 25, 25},
	}
	for _, tc := range cases {
		if got := clampLimit(tc.n, tc.def, tc.max); got != tc.want {
			t.Errorf("clampLimit(%d,%d,%d) = %d, want %d", tc.n, tc.def, tc.max, got, tc.want)
		}
	}
}

func TestEscapeLike(t *testing.T) {
	if got := escapeLike(`100%_done\`); got != `100\%\_done\\` {
		t.Errorf("escapeLike = %q", got)
	}
}

// TestLiveHub runs real queries against a Memory Hub when
// DOJO_TEST_MEMORY_DB_URL is set (e.g. through the dojo-bridge tunnel).
// Skipped otherwise, so CI stays hermetic.
func TestLiveHub(t *testing.T) {
	url := os.Getenv("DOJO_TEST_MEMORY_DB_URL")
	if url == "" {
		t.Skip("DOJO_TEST_MEMORY_DB_URL not set")
	}
	ctx := context.Background()
	c, err := New(ctx, url)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()
	if err := c.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}

	hits, err := c.SearchMemories(ctx, "night shift", "", 5)
	if err != nil {
		t.Fatalf("SearchMemories: %v", err)
	}
	if len(hits) == 0 {
		t.Fatal("SearchMemories(\"night shift\") returned 0 rows — hub content missing?")
	}
	t.Logf("search top hit: %s (%s) rank=%.4f", hits[0].Slug, hits[0].Type, hits[0].Rank)

	got, err := c.GetMemory(ctx, hits[0].Slug)
	if err != nil {
		t.Fatalf("GetMemory(%q): %v", hits[0].Slug, err)
	}
	if got.Body == "" {
		t.Fatalf("GetMemory(%q): empty body", hits[0].Slug)
	}

	recent, err := c.RecentMemories(ctx, "", 5)
	if err != nil {
		t.Fatalf("RecentMemories: %v", err)
	}
	if len(recent) == 0 {
		t.Fatal("RecentMemories returned 0 rows")
	}
}
