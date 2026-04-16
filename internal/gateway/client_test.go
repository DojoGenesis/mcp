package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ─── Constructor / accessors ──────────────────────────────────────────────────

func TestNew_CreatesValidClient(t *testing.T) {
	c := New("http://localhost:7340", "my-token")
	if c == nil {
		t.Fatal("New returned nil")
	}
	if c.http == nil {
		t.Error("http client is nil")
	}
	if c.streamHTTP == nil {
		t.Error("streamHTTP client is nil")
	}
}

func TestBaseURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"http://localhost:7340", "http://localhost:7340"},
		// Trailing slash should be stripped
		{"http://localhost:7340/", "http://localhost:7340"},
		{"http://gateway.example.com:9000//", "http://gateway.example.com:9000"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			c := New(tt.input, "")
			if got := c.BaseURL(); got != tt.want {
				t.Errorf("BaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ─── IsOnline ─────────────────────────────────────────────────────────────────

func TestIsOnline_ReturnsFalseWhenNoServerRunning(t *testing.T) {
	// Port 59999 should not have anything listening in CI/dev.
	c := New("http://localhost:59999", "")
	if c.IsOnline(context.Background()) {
		t.Error("IsOnline should return false when no server is running")
	}
}

func TestIsOnline_ReturnsTrueWhenServerHealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok", Version: "1.0"})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	if !c.IsOnline(context.Background()) {
		t.Error("IsOnline should return true when gateway responds with a non-empty status")
	}
}

// ─── Health ───────────────────────────────────────────────────────────────────

func TestHealth_ReturnsResponseOnSuccess(t *testing.T) {
	want := HealthResponse{
		Status:        "ok",
		Version:       "2.0.0",
		Timestamp:     "2026-04-15T00:00:00Z",
		UptimeSeconds: 3600,
		Providers:     map[string]string{"anthropic": "active"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	got, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("Health() unexpected error: %v", err)
	}
	if got.Status != want.Status {
		t.Errorf("Status: got %q, want %q", got.Status, want.Status)
	}
	if got.Version != want.Version {
		t.Errorf("Version: got %q, want %q", got.Version, want.Version)
	}
	if got.UptimeSeconds != want.UptimeSeconds {
		t.Errorf("UptimeSeconds: got %d, want %d", got.UptimeSeconds, want.UptimeSeconds)
	}
}

func TestHealth_ErrorOn500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.Health(context.Background())
	if err == nil {
		t.Fatal("Health() should return an error when server responds 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status code 500, got: %v", err)
	}
}

func TestHealth_ErrorOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{not valid json"))
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.Health(context.Background())
	if err == nil {
		t.Fatal("Health() should return an error on invalid JSON response")
	}
}

// ─── Memories ─────────────────────────────────────────────────────────────────

func TestMemories_ReturnsList(t *testing.T) {
	payload := `{"memories":[{"id":"1","content":"test","type":"user"}],"total":1}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/memory" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	mems, err := c.Memories(context.Background())
	if err != nil {
		t.Fatalf("Memories() unexpected error: %v", err)
	}
	if len(mems) != 1 {
		t.Fatalf("expected 1 memory, got %d", len(mems))
	}
	if mems[0].ID != "1" {
		t.Errorf("ID: got %q, want %q", mems[0].ID, "1")
	}
	if mems[0].Content != "test" {
		t.Errorf("Content: got %q, want %q", mems[0].Content, "test")
	}
	if mems[0].Type != "user" {
		t.Errorf("Type: got %q, want %q", mems[0].Type, "user")
	}
}

func TestMemories_ErrorOn500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.Memories(context.Background())
	if err == nil {
		t.Fatal("Memories() should return an error on 500")
	}
}

// ─── Seeds ────────────────────────────────────────────────────────────────────

func TestSeeds_ReturnsList(t *testing.T) {
	payload := `{"seeds":[{"id":"1","name":"test-seed","content":"pattern","description":""}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/seeds" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	seeds, err := c.Seeds(context.Background())
	if err != nil {
		t.Fatalf("Seeds() unexpected error: %v", err)
	}
	if len(seeds) != 1 {
		t.Fatalf("expected 1 seed, got %d", len(seeds))
	}
	if seeds[0].ID != "1" {
		t.Errorf("ID: got %q, want %q", seeds[0].ID, "1")
	}
	if seeds[0].Name != "test-seed" {
		t.Errorf("Name: got %q, want %q", seeds[0].Name, "test-seed")
	}
	if seeds[0].Content != "pattern" {
		t.Errorf("Content: got %q, want %q", seeds[0].Content, "pattern")
	}
}

func TestSeeds_ErrorOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json at all"))
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.Seeds(context.Background())
	if err == nil {
		t.Fatal("Seeds() should return an error on invalid JSON")
	}
}

// ─── ChatSync / SSE ───────────────────────────────────────────────────────────

func TestChatSync_CollectsDeltaText(t *testing.T) {
	// The mock server emits three SSE lines then [DONE].
	sseBody := strings.Join([]string{
		`data: {"delta":"hello "}`,
		``,
		`data: {"delta":"world"}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, sseBody)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	result, err := c.ChatSync(context.Background(), ChatRequest{Message: "hi"})
	if err != nil {
		t.Fatalf("ChatSync() unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("ChatSync() = %q, want %q", result, "hello world")
	}
}

func TestChatSync_ErrorOn500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, err := c.ChatSync(context.Background(), ChatRequest{Message: "hi"})
	if err == nil {
		t.Fatal("ChatSync() should return an error when server responds 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention 500, got: %v", err)
	}
}

func TestChatSync_EmptyResponseWhenNoDeltas(t *testing.T) {
	// SSE stream with no delta fields — result should be empty.
	sseBody := "data: [DONE]\n\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, sseBody)
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	result, err := c.ChatSync(context.Background(), ChatRequest{Message: "hi"})
	if err != nil {
		t.Fatalf("ChatSync() unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("ChatSync() = %q, want empty string", result)
	}
}

// ─── Authorization header ─────────────────────────────────────────────────────

func TestClient_SendsAuthorizationHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}))
	defer srv.Close()

	c := New(srv.URL, "bearer-xyz")
	_, _ = c.Health(context.Background())

	want := "Bearer bearer-xyz"
	if gotAuth != want {
		t.Errorf("Authorization header: got %q, want %q", gotAuth, want)
	}
}

func TestClient_NoAuthorizationHeaderWhenTokenEmpty(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}))
	defer srv.Close()

	c := New(srv.URL, "")
	_, _ = c.Health(context.Background())

	if gotAuth != "" {
		t.Errorf("Authorization header should be absent for empty token, got %q", gotAuth)
	}
}
