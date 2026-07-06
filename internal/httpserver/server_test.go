package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DojoGenesis/mcp/internal/authz"
	"github.com/DojoGenesis/mcp/internal/dojo"
	"github.com/mark3labs/mcp-go/server"
)

const testKey = "0123456789abcdef0123456789abcdef"

// newTestHandler wires the real stack the way cmd/server does in HTTP mode:
// dojo handler (bundled skills, temp ADR dir, nil gateway) + gate middleware
// + streamable transport + auth routes.
func newTestHandler(t *testing.T, dispatchLabels map[string]bool) http.Handler {
	t.Helper()

	handler, err := dojo.NewHandler("", t.TempDir(), nil, "")
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	s := server.NewMCPServer("dojo-mcp-server", "test",
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
		server.WithToolHandlerMiddleware(dojo.GateMiddleware(authz.NewLimiter(60))),
	)
	handler.RegisterTools(s)
	handler.RegisterResources(s)

	streamable := server.NewStreamableHTTPServer(s,
		server.WithStateLess(true),
		server.WithHeartbeatInterval(30*time.Second),
		server.WithHTTPContextFunc(IdentityContextFunc),
	)

	keys, err := authz.ParseKeys("test:" + testKey)
	if err != nil {
		t.Fatalf("ParseKeys: %v", err)
	}

	return NewHandler(Options{
		Version:        "test",
		Keys:           keys,
		DispatchLabels: dispatchLabels,
		MCP:            streamable,
	})
}

func postJSON(t *testing.T, ts *httptest.Server, path, bearer, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	return resp
}

const initializeBody = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}`

func TestHealth_NoAuthRequired(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("/health status = %d, want 200", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" || body["service"] != "dojo-mcp-server" {
		t.Fatalf("unexpected health body: %v", body)
	}
}

func TestMCP_401WithoutKey(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil))
	defer ts.Close()

	resp := postJSON(t, ts, "/mcp", "", initializeBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("keyless /mcp status = %d, want 401", resp.StatusCode)
	}
	if got := resp.Header.Get("WWW-Authenticate"); !strings.Contains(got, "Bearer") {
		t.Fatalf("WWW-Authenticate = %q, want Bearer challenge", got)
	}
}

func TestMCP_401WithWrongKey(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil))
	defer ts.Close()

	resp := postJSON(t, ts, "/mcp", "wrong-key-material-0000000000000", initializeBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong-key /mcp status = %d, want 401", resp.StatusCode)
	}
}

func TestMCP_InitializeWithKey(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil))
	defer ts.Close()

	resp := postJSON(t, ts, "/mcp", testKey, initializeBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("authed initialize status = %d, want 200", resp.StatusCode)
	}
	raw := readAll(t, resp)
	if !strings.Contains(raw, "dojo-mcp-server") {
		t.Fatalf("initialize response missing serverInfo: %s", truncate(raw, 300))
	}
}

func TestMCP_PathKeyFallback(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil))
	defer ts.Close()

	resp := postJSON(t, ts, "/mcp/k/"+testKey, "", initializeBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("path-key initialize status = %d, want 200", resp.StatusCode)
	}

	resp2 := postJSON(t, ts, "/mcp/k/not-the-key-000000000000000000", "", initializeBody)
	defer resp2.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("bad path-key status = %d, want 401", resp2.StatusCode)
	}
}

// TestDispatchClassDenied proves the full identity chain: HTTP auth →
// request context → HTTPContextFunc → tool middleware. A key that is not on
// the dispatch allowlist reaches the tool but gets a scoped denial.
func TestDispatchClassDenied(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, nil)) // no dispatch labels
	defer ts.Close()

	callBody := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"dojo_agent_chat","arguments":{"agent_id":"a","message":"hi"}}}`
	resp := postJSON(t, ts, "/mcp", testKey, callBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("tools/call transport status = %d, want 200", resp.StatusCode)
	}
	raw := readAll(t, resp)
	if !strings.Contains(raw, "not authorized for dispatch-class") {
		t.Fatalf("want dispatch denial in tool result, got: %s", truncate(raw, 400))
	}
}

// TestDispatchClassAllowedLabel verifies an allowlisted key passes the gate
// (and then fails at the gateway layer, since no gateway is configured —
// which proves the gate itself was the thing that opened).
func TestDispatchClassAllowedLabel(t *testing.T) {
	ts := httptest.NewServer(newTestHandler(t, map[string]bool{"test": true}))
	defer ts.Close()

	callBody := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"dojo_agent_chat","arguments":{"agent_id":"a","message":"hi"}}}`
	resp := postJSON(t, ts, "/mcp", testKey, callBody)
	defer resp.Body.Close() //nolint:errcheck // test cleanup; close error non-actionable
	raw := readAll(t, resp)
	if strings.Contains(raw, "not authorized for dispatch-class") {
		t.Fatalf("allowlisted key was denied: %s", truncate(raw, 400))
	}
	if !strings.Contains(raw, "Gateway is not configured") {
		t.Fatalf("expected gateway-not-configured tool error, got: %s", truncate(raw, 400))
	}
}

func readAll(t *testing.T, resp *http.Response) string {
	t.Helper()
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			break
		}
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
