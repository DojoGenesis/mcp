package dojo

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	tmpDir := t.TempDir()
	h, err := NewHandler("", tmpDir, nil)
	if err != nil {
		t.Fatalf("NewHandler returned error: %v", err)
	}
	return h
}

func newCallToolRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

// extractText is a test helper that pulls the text string from a CallToolResult.
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

// --- NewHandler Tests ---

func TestNewHandler(t *testing.T) {
	h := newTestHandler(t)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.wisdomBase == nil {
		t.Error("handler has nil wisdomBase")
	}
	if h.skillsLoader == nil {
		t.Error("handler has nil skillsLoader")
	}
	if h.decisionWriter == nil {
		t.Error("handler has nil decisionWriter")
	}
}

func TestNewHandler_WithSkillsPath(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := NewHandler("/nonexistent/path", tmpDir, nil)
	if err != nil {
		t.Fatalf("NewHandler with nonexistent skills path should not error (falls back to bundled): %v", err)
	}
	if h == nil {
		t.Fatal("handler is nil")
	}
}

func TestUnmarshalArgs(t *testing.T) {
	args := map[string]interface{}{
		"name":  "test",
		"query": "search term",
	}
	var dest struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}
	err := unmarshalArgs(args, &dest)
	if err != nil {
		t.Fatalf("unmarshalArgs returned error: %v", err)
	}
	if dest.Name != "test" {
		t.Errorf("name = %q, want %q", dest.Name, "test")
	}
	if dest.Query != "search term" {
		t.Errorf("query = %q, want %q", dest.Query, "search term")
	}
}

func TestUnmarshalArgs_Nil(t *testing.T) {
	var dest struct {
		Name string `json:"name"`
	}
	err := unmarshalArgs(nil, &dest)
	if err != nil {
		t.Fatalf("unmarshalArgs with nil returned error: %v", err)
	}
	if dest.Name != "" {
		t.Errorf("name = %q, expected empty", dest.Name)
	}
}

// --- Tool 1: dojo.scout ---

func TestHandleScout(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"situation": "Should we rewrite the backend in Go or keep Node.js?",
	})

	result, err := h.handleScout(ctx, req)
	if err != nil {
		t.Fatalf("handleScout returned error: %v", err)
	}
	if result.IsError {
		t.Error("handleScout returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Strategic Scout") {
		t.Error("output does not contain 'Strategic Scout' header")
	}
	if !strings.Contains(text, "Step 1") {
		t.Error("output does not contain 'Step 1'")
	}
	if !strings.Contains(text, "Step 4") {
		t.Error("output does not contain 'Step 4'")
	}
	if !strings.Contains(text, "rewrite the backend") {
		t.Error("output does not contain the situation text")
	}
}

func TestHandleScout_EmptySituation(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"situation": "",
	})

	result, err := h.handleScout(ctx, req)
	if err != nil {
		t.Fatalf("handleScout returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("handleScout should return error for empty situation")
	}
}

// --- Tool 2: dojo.invoke_skill ---

func TestHandleInvokeSkill_Found(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "strategic-scout",
	})

	result, err := h.handleInvokeSkill(ctx, req)
	if err != nil {
		t.Fatalf("handleInvokeSkill returned error: %v", err)
	}
	if result.IsError {
		t.Error("handleInvokeSkill returned error result for existing skill")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Skill: strategic-scout") {
		t.Error("output does not contain skill name")
	}
	if !strings.Contains(text, "Plugin:") {
		t.Error("output does not contain plugin label")
	}
}

func TestHandleInvokeSkill_NotFound(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "nonexistent-skill-xyz",
	})

	result, err := h.handleInvokeSkill(ctx, req)
	if err != nil {
		t.Fatalf("handleInvokeSkill returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("handleInvokeSkill should return error for nonexistent skill")
	}
	text := extractText(t, result)
	if !strings.Contains(strings.ToLower(text), "not found") {
		t.Errorf("error text should mention 'not found', got: %s", text[:100])
	}
}

func TestHandleInvokeSkill_NotFoundWithSuggestion(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "strategic",
	})

	result, err := h.handleInvokeSkill(ctx, req)
	if err != nil {
		t.Fatalf("handleInvokeSkill returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("should return error")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Did you mean") {
		t.Error("error should suggest similar skills")
	}
}

func TestHandleInvokeSkill_Empty(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"name": "",
	})

	result, err := h.handleInvokeSkill(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("should return error for empty name")
	}
}

// --- Tool 3: dojo.search_skills ---

func TestHandleSearchSkills_Found(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"query": "debugging",
	})

	result, err := h.handleSearchSkills(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	if result.IsError {
		t.Error("returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Skills matching") {
		t.Error("output does not contain 'Skills matching' header")
	}
	if !strings.Contains(text, "debug") {
		t.Error("output does not contain 'debug'")
	}
}

func TestHandleSearchSkills_NoMatch(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"query": "xyzzynonexistent12345",
	})

	result, err := h.handleSearchSkills(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(strings.ToLower(text), "no skills found") {
		t.Errorf("should say 'no skills found', got: %s", text)
	}
}

func TestHandleSearchSkills_SpecQuery(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"query": "write a spec",
	})

	result, err := h.handleSearchSkills(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Skills matching") {
		t.Error("output does not contain results header")
	}
}

// --- Tool 4: dojo.apply_seed ---

func TestHandleApplySeed_Known(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"seed_name": "three_tiered_governance",
		"situation": "Designing governance for a new AI project",
	})

	result, err := h.handleApplySeed(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	if result.IsError {
		t.Error("returned error result for known seed")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Applying Seed") {
		t.Error("output does not contain 'Applying Seed' header")
	}
	if !strings.Contains(text, "three_tiered_governance") {
		t.Error("output does not contain seed name")
	}
	if !strings.Contains(text, "Designing governance") {
		t.Error("output does not contain situation")
	}
	if !strings.Contains(text, "Reflection Questions") {
		t.Error("output does not contain reflection questions")
	}
}

func TestHandleApplySeed_Unknown(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"seed_name": "nonexistent_seed_xyz",
		"situation": "any",
	})

	result, err := h.handleApplySeed(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("should return error for unknown seed")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "not found") {
		t.Error("error should mention 'not found'")
	}
	if !strings.Contains(text, "Available seeds") {
		t.Error("error should list available seeds")
	}
}

func TestHandleApplySeed_EmptyFields(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"seed_name": "",
		"situation": "test",
	})

	result, err := h.handleApplySeed(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("should return error for empty seed_name")
	}
}

// --- Tool 5: dojo.log_decision ---

func TestHandleLogDecision(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"title":        "Use Go for the MCP server",
		"context":      "We need a reliable MCP server.",
		"decision":     "We chose Go for performance.",
		"consequences": "Team needs Go expertise.",
	})

	result, err := h.handleLogDecision(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	if result.IsError {
		t.Error("returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "ADR written to:") {
		t.Error("output does not confirm ADR was written")
	}

	// Verify file exists
	parts := strings.SplitN(text, "ADR written to: ", 2)
	if len(parts) == 2 {
		fp := strings.TrimSpace(parts[1])
		if _, err := os.Stat(fp); err != nil {
			t.Errorf("ADR file does not exist: %s", fp)
		}
	}
}

func TestHandleLogDecision_MissingFields(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"missing title", map[string]interface{}{"title": "", "context": "c", "decision": "d", "consequences": "e"}},
		{"missing context", map[string]interface{}{"title": "t", "context": "", "decision": "d", "consequences": "e"}},
		{"missing decision", map[string]interface{}{"title": "t", "context": "c", "decision": "", "consequences": "e"}},
		{"missing consequences", map[string]interface{}{"title": "t", "context": "c", "decision": "d", "consequences": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newCallToolRequest(tt.args)
			result, err := h.handleLogDecision(ctx, req)
			if err != nil {
				t.Fatalf("returned Go error: %v", err)
			}
			if !result.IsError {
				t.Error("should return error for missing field")
			}
		})
	}
}

// --- Tool 6: dojo.reflect ---

func TestHandleReflect(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"session_description": "Building a strategic architecture for the new platform",
	})

	result, err := h.handleReflect(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	if result.IsError {
		t.Error("returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Reflection") {
		t.Error("output does not contain 'Reflection' header")
	}
	if !strings.Contains(text, "Reflection Questions") {
		t.Error("output does not contain reflection questions")
	}
}

func TestHandleReflect_EmptyDescription(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{
		"session_description": "",
	})

	result, err := h.handleReflect(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !result.IsError {
		t.Error("should return error for empty description")
	}
}

func TestHandleReflect_WithMatches(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	// "debugging" should match the debugging skill
	req := newCallToolRequest(map[string]interface{}{
		"session_description": "debugging a complex performance issue in production",
	})

	result, err := h.handleReflect(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Relevant Skills") {
		t.Error("output should contain matched skills section")
	}
}

// --- Tool 7: dojo.list_skills ---

func TestHandleListSkills(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	req := newCallToolRequest(map[string]interface{}{})

	result, err := h.handleListSkills(ctx, req)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}
	if result.IsError {
		t.Error("returned error result")
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Dojo Genesis Skills") {
		t.Error("output does not contain header")
	}
	if !strings.Contains(text, "strategic-scout") {
		t.Error("output does not contain known skill 'strategic-scout'")
	}
	if !strings.Contains(text, "invoke_skill") {
		t.Error("output does not contain usage instruction")
	}
}

// --- Scaffold Tests ---

func TestFirstSentence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello world. More text.", "Hello world."},
		{"Short", "Short"},
		{"# Heading\n\nActual content. More.", "Actual content."},
		{"**Bold:** Real text. More.", "Real text."},
	}
	for _, tt := range tests {
		t.Run(tt.input[:min(20, len(tt.input))], func(t *testing.T) {
			result := firstSentence(tt.input)
			if result != tt.expected {
				t.Errorf("firstSentence(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractFirstStep(t *testing.T) {
	content := `# Workflow

## Steps

1. Identify the problem clearly.
2. Reproduce the issue.`

	step := extractFirstStep(content)
	if step == "" {
		t.Error("extractFirstStep returned empty string")
	}
	if !strings.Contains(step, "Identify the problem") {
		t.Errorf("extractFirstStep = %q, expected to contain 'Identify the problem'", step)
	}
}

func TestScoutScaffold(t *testing.T) {
	result := scoutScaffold("Should we use microservices?", nil)
	if !strings.Contains(result, "Strategic Scout") {
		t.Error("scaffold missing header")
	}
	if !strings.Contains(result, "microservices") {
		t.Error("scaffold missing situation")
	}
	if !strings.Contains(result, "Step 1") {
		t.Error("scaffold missing Step 1")
	}
	if !strings.Contains(result, "Step 4") {
		t.Error("scaffold missing Step 4")
	}
}

func TestScoutScaffold_WithMatchedSkills(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	// Use a situation that should match skills
	req := newCallToolRequest(map[string]interface{}{
		"situation": "strategic architecture decision for the platform",
	})

	result, err := h.handleScout(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Strategic Scout") {
		t.Error("scaffold missing header")
	}
	// Should contain methodology section if skills matched
	if !strings.Contains(text, "strategic") {
		t.Error("output should reference strategic content")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
