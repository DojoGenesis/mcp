package dojo

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.wisdomBase == nil {
		t.Error("NewHandler created handler with nil wisdomBase")
	}
}

func TestUnmarshalArgs(t *testing.T) {
	args := map[string]interface{}{
		"name":  "test_seed",
		"query": "governance",
	}

	var dest struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}

	err := unmarshalArgs(args, &dest)
	if err != nil {
		t.Fatalf("unmarshalArgs returned error: %v", err)
	}
	if dest.Name != "test_seed" {
		t.Errorf("unmarshalArgs name = %q, want %q", dest.Name, "test_seed")
	}
	if dest.Query != "governance" {
		t.Errorf("unmarshalArgs query = %q, want %q", dest.Query, "governance")
	}
}

func TestUnmarshalArgs_Invalid(t *testing.T) {
	// nil map should still work (results in zero-value struct)
	var dest struct {
		Name string `json:"name"`
	}
	err := unmarshalArgs(nil, &dest)
	if err != nil {
		t.Fatalf("unmarshalArgs with nil map returned error: %v", err)
	}
	if dest.Name != "" {
		t.Errorf("unmarshalArgs with nil map set name to %q, expected empty", dest.Name)
	}

	// Bad type: passing a number where a string is expected should not error
	// at the unmarshal level (JSON is flexible), but the value will be wrong type.
	badArgs := map[string]interface{}{
		"name": 12345,
	}
	var dest2 struct {
		Name string `json:"name"`
	}
	// json.Unmarshal into string from number should fail
	err = unmarshalArgs(badArgs, &dest2)
	if err != nil {
		// This is acceptable: strict type mismatch
		return
	}
	// If no error, the name should be empty or the string representation
	// Either outcome is acceptable for this test
}

func newCallToolRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestHandleReflect_AllModes(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	modes := []string{"mirror", "scout", "gardener", "implementation"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			req := newCallToolRequest(map[string]interface{}{
				"situation":    "Should we refactor the codebase?",
				"perspectives": []interface{}{"maintainability", "cost", "risk"},
				"mode":         mode,
			})

			result, err := h.handleReflect(ctx, req)
			if err != nil {
				t.Fatalf("handleReflect(%s) returned error: %v", mode, err)
			}
			if result == nil {
				t.Fatalf("handleReflect(%s) returned nil result", mode)
			}
			if result.IsError {
				t.Errorf("handleReflect(%s) returned error result", mode)
			}
			text := extractText(t, result)
			if text == "" {
				t.Errorf("handleReflect(%s) returned empty text", mode)
			}
			// Each mode should include the mode name in uppercase in the output
			expectedMode := strings.ToUpper(mode)
			if !strings.Contains(text, expectedMode) {
				t.Errorf("handleReflect(%s) output does not contain %q", mode, expectedMode)
			}
		})
	}
}

func TestHandleReflect_UnknownMode(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"situation":    "Test situation",
		"perspectives": []interface{}{"perspective1"},
		"mode":         "nonexistent_mode",
	})

	result, err := h.handleReflect(ctx, req)
	if err != nil {
		t.Fatalf("handleReflect(unknown) returned Go error: %v", err)
	}
	if result == nil {
		t.Fatal("handleReflect(unknown) returned nil result")
	}
	text := extractText(t, result)
	if !strings.Contains(strings.ToLower(text), "unknown") {
		t.Errorf("handleReflect(unknown mode) should mention 'unknown' in output, got: %s", text)
	}
}

func TestHandleSearchWisdom(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"query": "governance",
	})

	result, err := h.handleSearchWisdom(ctx, req)
	if err != nil {
		t.Fatalf("handleSearchWisdom returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleSearchWisdom returned nil result")
	}
	if result.IsError {
		t.Error("handleSearchWisdom returned error result")
	}
	text := extractText(t, result)
	if text == "" {
		t.Error("handleSearchWisdom returned empty text")
	}
	// Should contain JSON array with results
	if !strings.Contains(text, "governance") {
		t.Error("handleSearchWisdom for 'governance' does not contain 'governance' in results")
	}
}

func TestHandleGetSeed_Found(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "three_tiered_governance",
	})

	result, err := h.handleGetSeed(ctx, req)
	if err != nil {
		t.Fatalf("handleGetSeed returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleGetSeed returned nil result")
	}
	if result.IsError {
		t.Error("handleGetSeed returned error result for existing seed")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "three_tiered_governance") {
		t.Error("handleGetSeed result does not contain seed name")
	}
}

func TestHandleGetSeed_NotFound(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "nonexistent_seed_xyz",
	})

	result, err := h.handleGetSeed(ctx, req)
	if err != nil {
		t.Fatalf("handleGetSeed returned Go error: %v", err)
	}
	if result == nil {
		t.Fatal("handleGetSeed returned nil result")
	}
	if !result.IsError {
		t.Error("handleGetSeed should return error result for nonexistent seed")
	}
	text := extractText(t, result)
	if !strings.Contains(strings.ToLower(text), "not found") {
		t.Errorf("handleGetSeed error text should mention 'not found', got: %s", text)
	}
}

func TestHandleListSeeds(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{})

	result, err := h.handleListSeeds(ctx, req)
	if err != nil {
		t.Fatalf("handleListSeeds returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleListSeeds returned nil result")
	}
	if result.IsError {
		t.Error("handleListSeeds returned error result")
	}
	text := extractText(t, result)
	if text == "" {
		t.Error("handleListSeeds returned empty text")
	}
	// Should list at least one known seed
	if !strings.Contains(text, "three_tiered_governance") {
		t.Error("handleListSeeds output does not contain expected seed name")
	}
}

func TestHandleGetPrinciples(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{})

	result, err := h.handleGetPrinciples(ctx, req)
	if err != nil {
		t.Fatalf("handleGetPrinciples returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleGetPrinciples returned nil result")
	}
	if result.IsError {
		t.Error("handleGetPrinciples returned error result")
	}
	text := extractText(t, result)
	if text == "" {
		t.Error("handleGetPrinciples returned empty text")
	}
	if !strings.Contains(text, "Beginner") {
		t.Error("handleGetPrinciples does not mention Beginner's Mind")
	}
}

func TestHandleCreateThinkingRoom(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"topic":      "system architecture",
		"agent_name": "TestAgent",
	})

	result, err := h.handleCreateThinkingRoom(ctx, req)
	if err != nil {
		t.Fatalf("handleCreateThinkingRoom returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleCreateThinkingRoom returned nil result")
	}
	if result.IsError {
		t.Error("handleCreateThinkingRoom returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "system architecture") {
		t.Error("handleCreateThinkingRoom output does not contain the topic")
	}
	if !strings.Contains(text, "TestAgent") {
		t.Error("handleCreateThinkingRoom output does not contain the agent name")
	}
	if !strings.Contains(text, "Thinking Room") {
		t.Error("handleCreateThinkingRoom output does not contain 'Thinking Room' header")
	}
}

func TestHandleTraceLineage(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"idea_or_insight": "governance multiplies velocity",
	})

	result, err := h.handleTraceLineage(ctx, req)
	if err != nil {
		t.Fatalf("handleTraceLineage returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleTraceLineage returned nil result")
	}
	if result.IsError {
		t.Error("handleTraceLineage returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Lineage Trace") {
		t.Error("handleTraceLineage output does not contain 'Lineage Trace' header")
	}
	if !strings.Contains(text, "governance multiplies velocity") {
		t.Error("handleTraceLineage output does not contain the idea")
	}
}

func TestHandleCheckPace(t *testing.T) {
	h := NewHandler()
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"session_description": "Rapidly building three features in parallel without pausing to reflect",
	})

	result, err := h.handleCheckPace(ctx, req)
	if err != nil {
		t.Fatalf("handleCheckPace returned error: %v", err)
	}
	if result == nil {
		t.Fatal("handleCheckPace returned nil result")
	}
	if result.IsError {
		t.Error("handleCheckPace returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Pace Check") {
		t.Error("handleCheckPace output does not contain 'Pace Check' header")
	}
	if !strings.Contains(text, "Understanding") {
		t.Error("handleCheckPace output does not mention 'Understanding'")
	}
	if !strings.Contains(text, "Extraction") {
		t.Error("handleCheckPace output does not mention 'Extraction'")
	}
}

// extractText is a test helper that pulls the text string from a CallToolResult.
// The result Content is []interface{} where each item is a TextContent.
func extractText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		return ""
	}
	tc, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("result.Content[0] is %T, not mcp.TextContent", result.Content[0])
	}
	return tc.Text
}
