package dojo

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DojoGenesis/mcp/internal/memhub"
	"github.com/mark3labs/mcp-go/mcp"
)

// fakeHub implements Hub for tests.
type fakeHub struct {
	entries []memhub.Entry
	err     error
}

func (f *fakeHub) SearchMemories(_ context.Context, query, typ string, limit int) ([]memhub.Entry, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []memhub.Entry
	for _, e := range f.entries {
		if typ != "" && e.Type != typ {
			continue
		}
		if strings.Contains(strings.ToLower(e.Name+" "+e.Body), strings.ToLower(query)) {
			out = append(out, e)
		}
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakeHub) GetMemory(_ context.Context, slug string) (*memhub.Entry, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, e := range f.entries {
		if e.Slug == slug {
			return &e, nil
		}
	}
	return nil, memhub.ErrNotFound
}

func (f *fakeHub) RecentMemories(_ context.Context, typ string, limit int) ([]memhub.Entry, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []memhub.Entry
	for _, e := range f.entries {
		if typ == "" || e.Type == typ {
			out = append(out, e)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func newHubHandler(t *testing.T, hub Hub) *Handler {
	t.Helper()
	h, err := NewHandler("", t.TempDir(), nil, "")
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	if hub != nil {
		h.SetHub(hub)
	}
	return h
}

func callTool(t *testing.T, fn func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error), args map[string]any) *mcp.CallToolResult {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	res, err := fn(context.Background(), req)
	if err != nil {
		t.Fatalf("handler returned protocol error: %v", err)
	}
	if res == nil {
		t.Fatal("handler returned nil result")
	}
	return res
}

func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := mcp.AsTextContent(c); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String()
}

var hubFixtures = []memhub.Entry{
	{Slug: "project-nightshift-family", Name: "Night Shift Family", Type: "project",
		Description: "V1 live", Snippet: "hard-wall trial", Body: "Night Shift V1 body", Updated: "2026-07-05"},
	{Slug: "reference-dojo-bridge", Name: "dojo-bridge hub", Type: "reference",
		Description: "bridge wiring", Snippet: "gateway 7340", Body: "bridge body", Updated: "2026-07-06"},
}

func TestSearchMemoryHub_NotConfigured(t *testing.T) {
	h := newHubHandler(t, nil)
	res := callTool(t, h.handleSearchMemoryHub, map[string]any{"query": "night shift"})
	if !res.IsError {
		t.Fatal("want error result when hub unset")
	}
	if !strings.Contains(resultText(t, res), "not configured") {
		t.Fatalf("want not-configured text, got: %s", resultText(t, res))
	}
}

func TestSearchMemoryHub_ReturnsRankedList(t *testing.T) {
	h := newHubHandler(t, &fakeHub{entries: hubFixtures})
	res := callTool(t, h.handleSearchMemoryHub, map[string]any{"query": "night shift"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	for _, want := range []string{"project-nightshift-family", "Night Shift Family", "dojo_get_memory"} {
		if !strings.Contains(text, want) {
			t.Errorf("result missing %q:\n%s", want, text)
		}
	}
}

func TestSearchMemoryHub_EmptyQueryRejected(t *testing.T) {
	h := newHubHandler(t, &fakeHub{entries: hubFixtures})
	res := callTool(t, h.handleSearchMemoryHub, map[string]any{"query": "  "})
	if !res.IsError {
		t.Fatal("want error for empty query")
	}
}

func TestGetMemoryHub_FullBody(t *testing.T) {
	h := newHubHandler(t, &fakeHub{entries: hubFixtures})
	res := callTool(t, h.handleGetMemoryHub, map[string]any{"slug": "project-nightshift-family"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(text, "Night Shift V1 body") {
		t.Fatalf("want full body, got: %s", text)
	}
}

func TestGetMemoryHub_NotFound(t *testing.T) {
	h := newHubHandler(t, &fakeHub{entries: hubFixtures})
	res := callTool(t, h.handleGetMemoryHub, map[string]any{"slug": "nope"})
	if !res.IsError {
		t.Fatal("want error result for unknown slug")
	}
	if !strings.Contains(resultText(t, res), "nope") {
		t.Fatalf("error should name the slug: %s", resultText(t, res))
	}
}

func TestRecentMemoriesHub_TypeFilter(t *testing.T) {
	h := newHubHandler(t, &fakeHub{entries: hubFixtures})
	res := callTool(t, h.handleRecentMemoriesHub, map[string]any{"type": "reference"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !strings.Contains(text, "reference-dojo-bridge") || strings.Contains(text, "project-nightshift-family") {
		t.Fatalf("type filter not applied:\n%s", text)
	}
}

func TestSearchMemoryHub_BackendError(t *testing.T) {
	h := newHubHandler(t, &fakeHub{err: fmt.Errorf("connection refused")})
	res := callTool(t, h.handleSearchMemoryHub, map[string]any{"query": "x y"})
	if !res.IsError {
		t.Fatal("want error result on backend failure")
	}
	if !strings.Contains(resultText(t, res), "connection refused") {
		t.Fatalf("error should surface cause: %s", resultText(t, res))
	}
}

func TestDispatch_NoGateway(t *testing.T) {
	h := newHubHandler(t, nil)
	res := callTool(t, h.handleDispatch, map[string]any{"prompt": "hello"})
	if !res.IsError {
		t.Fatal("want error result with nil gateway")
	}
	if !strings.Contains(resultText(t, res), "Gateway is not configured") {
		t.Fatalf("unexpected error text: %s", resultText(t, res))
	}
}

func TestDispatch_EmptyPromptRejected(t *testing.T) {
	h := newHubHandler(t, nil)
	res := callTool(t, h.handleDispatch, map[string]any{"prompt": ""})
	if !res.IsError {
		t.Fatal("want error result for empty prompt")
	}
}
