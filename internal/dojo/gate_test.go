package dojo

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/DojoGenesis/mcp/internal/authz"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// silenceGateLogs redirects the package-level logger for the duration of a
// test so GateMiddleware's log.Printf calls don't clutter test output, then
// restores the previous output on cleanup.
func silenceGateLogs(t *testing.T) {
	t.Helper()
	prev := log.Writer()
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(prev) })
}

// newGateRequest builds a mcp.CallToolRequest naming the given tool, the way
// the handler_test.go / handler_hub_test.go helpers build requests.
func newGateRequest(tool string) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Name = tool
	return req
}

// recordingNext returns a server.ToolHandlerFunc that records whether it was
// invoked and echoes back the given result/error.
func recordingNext(res *mcp.CallToolResult, err error) (server.ToolHandlerFunc, *bool) {
	called := false
	fn := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		called = true
		return res, err
	}
	return fn, &called
}

// gateResultText pulls the text out of a CallToolResult the way
// handler_hub_test.go's resultText does.
func gateResultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := mcp.AsTextContent(c); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String()
}

// --- 1. Dispatch-denied branch ---

func TestGateMiddleware_DispatchDenied(t *testing.T) {
	silenceGateLogs(t)

	limiter := authz.NewLimiter(60)
	next, called := recordingNext(mcp.NewToolResultText("should not be reached"), nil)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), "test", false)
	req := newGateRequest("dojo_dispatch")

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("GateMiddleware returned Go error: %v", err)
	}
	if *called {
		t.Fatal("next handler was invoked despite dispatch being denied")
	}
	if res == nil || !res.IsError {
		t.Fatal("want an error result for dispatch-denied identity")
	}
	text := gateResultText(t, res)
	if !strings.Contains(text, "not authorized for dispatch-class") {
		t.Errorf("error text should explain dispatch denial, got: %s", text)
	}
	if !strings.Contains(text, "DOJO_DISPATCH_ALLOWED_LABELS") {
		t.Errorf("error text should mention the remediation env var, got: %s", text)
	}
}

// --- 2. Rate-limited branch ---

func TestGateMiddleware_RateLimited(t *testing.T) {
	silenceGateLogs(t)

	// perMin=1 gives burst=1: the bucket starts full with exactly one token,
	// so the first call succeeds deterministically and the second is
	// guaranteed to be denied without any sleeps or wall-clock waits.
	limiter := authz.NewLimiter(1)
	next, called := recordingNext(mcp.NewToolResultText("ok"), nil)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), "test", true)
	req := newGateRequest("dojo_dispatch")

	// First call consumes the only token and should reach next.
	if _, err := gated(ctx, req); err != nil {
		t.Fatalf("first call returned Go error: %v", err)
	}
	if !*called {
		t.Fatal("first call should have reached next (burst=1 token available)")
	}

	// Reset the recorder for the second call.
	*called = false

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("second call returned Go error: %v", err)
	}
	if *called {
		t.Fatal("next handler was invoked despite the limiter being exhausted")
	}
	if res == nil || !res.IsError {
		t.Fatal("want an error result once the limiter is exhausted")
	}
	text := gateResultText(t, res)
	if !strings.Contains(text, "rate limit") {
		t.Errorf("error text should mention rate limiting, got: %s", text)
	}
}

// --- 3. Dispatch-class allowed path ---

func TestGateMiddleware_DispatchAllowedPath(t *testing.T) {
	silenceGateLogs(t)

	limiter := authz.NewLimiter(60)
	want := mcp.NewToolResultText("inner result")
	next, called := recordingNext(want, nil)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), "test", true)
	req := newGateRequest("dojo_dispatch")

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !*called {
		t.Fatal("next handler was not invoked for an allowed dispatch-class call")
	}
	if res != want {
		t.Fatalf("result was not passed through unchanged: got %#v, want %#v", res, want)
	}
}

// --- 4. Non-dispatch tool bypasses the gate entirely ---

func TestGateMiddleware_NonDispatchToolBypassesGate(t *testing.T) {
	silenceGateLogs(t)

	// perMin=1 -> burst=1; exhaust the single token up front so any dispatch-
	// class call under this label would be denied. The identity is also
	// dispatch-denied. A non-dispatch tool must ignore both.
	limiter := authz.NewLimiter(1)
	label := "exhausted"
	if !limiter.Allow(label) {
		t.Fatal("setup: expected the first Allow call to succeed and consume the token")
	}

	want := mcp.NewToolResultText("health ok")
	next, called := recordingNext(want, nil)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), label, false)
	req := newGateRequest("dojo_health")

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !*called {
		t.Fatal("next handler was not invoked for a non-dispatch-class tool")
	}
	if res != want {
		t.Fatalf("result was not passed through unchanged: got %#v, want %#v", res, want)
	}
}

// --- 5. No-identity / stdio local mode ---

func TestGateMiddleware_NoIdentityDefaultsAllowed(t *testing.T) {
	silenceGateLogs(t)

	limiter := authz.NewLimiter(60)
	want := mcp.NewToolResultText("dispatched")
	next, called := recordingNext(want, nil)
	gated := GateMiddleware(limiter)(next)

	// No authz.WithIdentity call at all — plain background context, as in
	// stdio local mode.
	ctx := context.Background()
	req := newGateRequest("dojo_dispatch")

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("returned Go error: %v", err)
	}
	if !*called {
		t.Fatal("next handler was not invoked when no identity is present (DispatchAllowed should default true)")
	}
	if res != want {
		t.Fatalf("result was not passed through unchanged: got %#v, want %#v", res, want)
	}
}

func TestGateMiddleware_NoIdentityStillRateLimitedUnderLocalLabel(t *testing.T) {
	silenceGateLogs(t)

	// burst=1: the first no-identity dispatch call consumes the "local"
	// bucket's only token; the second must be rate limited even though
	// DispatchAllowed defaults true.
	limiter := authz.NewLimiter(1)
	next, called := recordingNext(mcp.NewToolResultText("ok"), nil)
	gated := GateMiddleware(limiter)(next)

	ctx := context.Background()
	req := newGateRequest("dojo_dispatch")

	if _, err := gated(ctx, req); err != nil {
		t.Fatalf("first call returned Go error: %v", err)
	}
	if !*called {
		t.Fatal("first no-identity call should have reached next")
	}
	*called = false

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("second call returned Go error: %v", err)
	}
	if *called {
		t.Fatal("next handler was invoked despite the local-label limiter being exhausted")
	}
	if res == nil || !res.IsError {
		t.Fatal("want an error result once the local-label limiter is exhausted")
	}
}

// --- 6. DispatchClassTools membership ---

func TestDispatchClassTools_Membership(t *testing.T) {
	want := []string{"dojo_dispatch", "dojo_agent_dispatch", "dojo_agent_chat"}

	if len(DispatchClassTools) != len(want) {
		t.Fatalf("DispatchClassTools has %d entries, want exactly %d: %v", len(DispatchClassTools), len(want), DispatchClassTools)
	}
	for _, tool := range want {
		if !DispatchClassTools[tool] {
			t.Errorf("DispatchClassTools missing expected entry %q", tool)
		}
	}
	if DispatchClassTools["dojo_health"] {
		t.Error("DispatchClassTools should not contain non-dispatch tool dojo_health")
	}
}

// --- 7. Passthrough of next's error and of an IsError result ---

func TestGateMiddleware_PassthroughGoError(t *testing.T) {
	silenceGateLogs(t)

	limiter := authz.NewLimiter(60)
	wantErr := errors.New("boom from inner handler")
	next, called := recordingNext(nil, wantErr)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), "test", true)
	req := newGateRequest("dojo_dispatch")

	res, err := gated(ctx, req)
	if !*called {
		t.Fatal("next handler was not invoked")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("error was not passed through unchanged: got %v, want %v", err, wantErr)
	}
	if res != nil {
		t.Fatalf("expected nil result alongside the passthrough error, got: %#v", res)
	}
}

func TestGateMiddleware_PassthroughIsErrorResult(t *testing.T) {
	silenceGateLogs(t)

	limiter := authz.NewLimiter(60)
	want := mcp.NewToolResultError("inner tool-level failure")
	next, called := recordingNext(want, nil)
	gated := GateMiddleware(limiter)(next)

	ctx := authz.WithIdentity(context.Background(), "test", true)
	req := newGateRequest("dojo_dispatch")

	res, err := gated(ctx, req)
	if err != nil {
		t.Fatalf("returned unexpected Go error: %v", err)
	}
	if !*called {
		t.Fatal("next handler was not invoked")
	}
	if res != want {
		t.Fatalf("IsError result was not passed through unchanged: got %#v, want %#v", res, want)
	}
	if !res.IsError {
		t.Fatal("sanity check: want's IsError should be true")
	}
}
