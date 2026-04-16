// Package gateway provides a lightweight HTTP+SSE client for the Dojo Genesis AgenticGateway.
// It uses only stdlib — no CLI module dependency.
package gateway

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client talks to an AgenticGateway instance.
type Client struct {
	base       string
	token      string
	http       *http.Client
	streamHTTP *http.Client // no timeout — for SSE streaming
}

// New creates a Client targeting baseURL with the given auth token.
func New(baseURL, token string) *Client {
	return &Client{
		base:  strings.TrimRight(baseURL, "/"),
		token: token,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		streamHTTP: &http.Client{},
	}
}

// BaseURL returns the configured gateway base URL.
func (c *Client) BaseURL() string {
	return c.base
}

// ─── Health ───────────────────────────────────────────────────────────────────

// HealthResponse is the /health response.
type HealthResponse struct {
	Status        string            `json:"status"`
	Version       string            `json:"version"`
	Timestamp     string            `json:"timestamp"`
	Uptime        string            `json:"uptime,omitempty"`
	Providers     map[string]string `json:"providers"`
	UptimeSeconds int64             `json:"uptime_seconds"`
}

// IsOnline calls GET /health and returns true if the gateway is reachable and healthy.
func (c *Client) IsOnline(ctx context.Context) bool {
	h, err := c.Health(ctx)
	if err != nil {
		return false
	}
	return h.Status != ""
}

// Health fetches GET /health.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var h HealthResponse
	if err := c.get(ctx, "/health", &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ─── Memory ───────────────────────────────────────────────────────────────────

// Memory represents an entry from /v1/memory.
type Memory struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type memoriesResponse struct {
	Memories []Memory `json:"memories"`
	Total    int      `json:"total"`
}

// Memories fetches GET /v1/memory.
func (c *Client) Memories(ctx context.Context) ([]Memory, error) {
	var r memoriesResponse
	if err := c.get(ctx, "/v1/memory", &r); err != nil {
		return nil, err
	}
	return r.Memories, nil
}

// StoreMemoryRequest is the body for POST /v1/memory.
type StoreMemoryRequest struct {
	Content string `json:"content"`
	Type    string `json:"type,omitempty"`
}

// StoreMemory creates a new memory entry.
func (c *Client) StoreMemory(ctx context.Context, req StoreMemoryRequest) (*Memory, error) {
	var wrapper struct {
		Memory Memory `json:"memory"`
	}
	if err := c.post(ctx, "/v1/memory", req, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Memory, nil
}

// SearchMemories performs a semantic search across memories.
func (c *Client) SearchMemories(ctx context.Context, query string) ([]Memory, error) {
	body := map[string]string{"query": query}
	var r memoriesResponse
	if err := c.post(ctx, "/v1/memory/search", body, &r); err != nil {
		return nil, err
	}
	return r.Memories, nil
}

// ─── Seeds ────────────────────────────────────────────────────────────────────

// Seed is a single entry in the /v1/seeds response.
type Seed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type seedsEnvelope struct {
	Seeds []Seed `json:"seeds"`
}

// Seeds fetches GET /v1/seeds.
func (c *Client) Seeds(ctx context.Context) ([]Seed, error) {
	var r seedsEnvelope
	if err := c.get(ctx, "/v1/seeds", &r); err != nil {
		return nil, err
	}
	return r.Seeds, nil
}

// CreateSeedRequest is the body for POST /v1/seeds.
type CreateSeedRequest struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
}

// CreateSeed posts a new seed to /v1/seeds.
func (c *Client) CreateSeed(ctx context.Context, req CreateSeedRequest) (*Seed, error) {
	var wrapper struct {
		Seed Seed `json:"seed"`
	}
	if err := c.post(ctx, "/v1/seeds", req, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Seed, nil
}

// ─── Agents ───────────────────────────────────────────────────────────────────

// Agent is a single entry in the /v1/gateway/agents response.
type Agent struct {
	ID          string `json:"agent_id"`
	Name        string `json:"name"`
	Mode        string `json:"active_mode"`
	Status      string `json:"status"`
	Disposition string `json:"disposition"`
}

type agentsEnvelope struct {
	Agents []Agent `json:"agents"`
	Total  int     `json:"total"`
}

// Agents fetches GET /v1/gateway/agents.
func (c *Client) Agents(ctx context.Context) ([]Agent, error) {
	var r agentsEnvelope
	if err := c.get(ctx, "/v1/gateway/agents", &r); err != nil {
		return nil, err
	}
	return r.Agents, nil
}

// CreateAgentRequest is the body for POST /v1/gateway/agents.
type CreateAgentRequest struct {
	Name        string `json:"name,omitempty"`
	Mode        string `json:"active_mode,omitempty"`
	Disposition string `json:"disposition,omitempty"`
}

// CreateAgent creates a new agent in the gateway.
func (c *Client) CreateAgent(ctx context.Context, req CreateAgentRequest) (*Agent, error) {
	var r Agent
	if err := c.post(ctx, "/v1/gateway/agents", req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ─── Chat (streaming SSE) ─────────────────────────────────────────────────────

// ChatRequest mirrors the /v1/chat request body.
type ChatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	Stream    bool   `json:"stream"`
}

// SSEChunk is a parsed line from the SSE stream.
type SSEChunk struct {
	Event string
	Data  string
}

// ChatStream opens a streaming POST /v1/chat and calls onChunk for each SSE event.
// Returns when the stream ends or ctx is cancelled.
func (c *Client) ChatStream(ctx context.Context, req ChatRequest, onChunk func(SSEChunk)) error {
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/v1/chat", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.streamHTTP.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gateway returned %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	return parseSSE(resp.Body, onChunk)
}

// ChatSync wraps ChatStream, collects all delta text, and returns the concatenated result.
// Uses a 120-second timeout via context.
func (c *Client) ChatSync(ctx context.Context, req ChatRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var sb strings.Builder
	err := c.ChatStream(ctx, req, func(chunk SSEChunk) {
		// Try to parse delta text from JSON data field
		text := extractDeltaText(chunk.Data)
		sb.WriteString(text)
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

// AgentChatSync sends a message to a specific agent and collects the full SSE response.
func (c *Client) AgentChatSync(ctx context.Context, agentID string, message string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	body := map[string]interface{}{
		"message": message,
		"stream":  true,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	path := "/v1/gateway/agents/" + agentID + "/chat"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.streamHTTP.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gateway returned %d: %s", resp.StatusCode, strings.TrimSpace(string(rb)))
	}

	var sb strings.Builder
	if err := parseSSE(resp.Body, func(chunk SSEChunk) {
		sb.WriteString(extractDeltaText(chunk.Data))
	}); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// ─── Skills ───────────────────────────────────────────────────────────────────

// Skill is a single entry in the /api/skills response.
type Skill struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Plugin      string   `json:"plugin,omitempty"`
	Triggers    []string `json:"triggers,omitempty"`
}

type skillsEnvelope struct {
	Skills []Skill `json:"skills"`
	Total  int     `json:"total"`
}

// Skills fetches all skills from the gateway.
func (c *Client) Skills(ctx context.Context) ([]Skill, error) {
	const pageSize = 100
	var all []Skill
	offset := 0
	for {
		var r skillsEnvelope
		path := fmt.Sprintf("/api/skills?limit=%d&offset=%d", pageSize, offset)
		if err := c.get(ctx, path, &r); err != nil {
			return nil, err
		}
		all = append(all, r.Skills...)
		if len(r.Skills) == 0 || (r.Total > 0 && len(all) >= r.Total) {
			break
		}
		offset += len(r.Skills)
	}
	return all, nil
}

// SearchSkills searches skills by query string.
func (c *Client) SearchSkills(ctx context.Context, query string) ([]Skill, error) {
	const pageSize = 100
	var all []Skill
	offset := 0
	for {
		var r skillsEnvelope
		path := fmt.Sprintf("/api/skills?q=%s&limit=%d&offset=%d", url.QueryEscape(query), pageSize, offset)
		if err := c.get(ctx, path, &r); err != nil {
			return nil, err
		}
		all = append(all, r.Skills...)
		if len(r.Skills) == 0 || (r.Total > 0 && len(all) >= r.Total) {
			break
		}
		offset += len(r.Skills)
	}
	return all, nil
}

// ─── Orchestration ────────────────────────────────────────────────────────────

// OrchestrateStep is a single step in an orchestration request.
type OrchestrateStep struct {
	Name  string         `json:"name"`
	Tool  string         `json:"tool"`
	Input map[string]any `json:"input,omitempty"`
}

// OrchestrateRequest is the body for POST /v1/gateway/orchestrate.
type OrchestrateRequest struct {
	Task  string            `json:"task"`
	Steps []OrchestrateStep `json:"steps,omitempty"`
}

// OrchestrationStatus is returned by POST /v1/gateway/orchestrate.
type OrchestrationStatus struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
}

// Orchestrate submits an orchestration request (async — returns immediately).
func (c *Client) Orchestrate(ctx context.Context, req OrchestrateRequest) (*OrchestrationStatus, error) {
	var r OrchestrationStatus
	if err := c.post(ctx, "/v1/gateway/orchestrate", req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// DAGNode is a single node in the DAG status response.
type DAGNode struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// DAGStatus is the response from GET /v1/gateway/orchestrate/:id/dag.
type DAGStatus struct {
	ExecutionID string    `json:"execution_id"`
	Status      string    `json:"status"`
	Nodes       []DAGNode `json:"nodes,omitempty"`
}

// OrchestrationDAG polls the DAG state for an execution.
func (c *Client) OrchestrationDAG(ctx context.Context, executionID string) (*DAGStatus, error) {
	var r DAGStatus
	if err := c.get(ctx, "/v1/gateway/orchestrate/"+executionID+"/dag", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ─── SSE parser ───────────────────────────────────────────────────────────────

// parseSSE reads an SSE stream and calls onChunk for each data line.
// Per SSE spec, the event field resets on blank lines (event dispatch boundary).
// Buffer is raised to 1MB to handle large JSON payloads.
func parseSSE(r io.Reader, onChunk func(SSEChunk)) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var event string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "[DONE]" {
				return nil
			}
			onChunk(SSEChunk{Event: event, Data: data})
		case line == "":
			// SSE dispatch boundary — reset event for the next block
			event = ""
		}
	}
	return scanner.Err()
}

// extractDeltaText tries to extract delta text from a JSON SSE data payload.
// Gateway SSE data is typically: {"type":"delta","content":"..."}
// Falls back to returning data as-is if it's not JSON.
func extractDeltaText(data string) string {
	var msg struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Delta   string `json:"delta"`
		Text    string `json:"text"`
	}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		// Not JSON — return raw data
		return data
	}
	if msg.Content != "" {
		return msg.Content
	}
	if msg.Delta != "" {
		return msg.Delta
	}
	if msg.Text != "" {
		return msg.Text
	}
	return ""
}

// ─── HTTP helpers ─────────────────────────────────────────────────────────────

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gateway %s returned %d: %s", path, resp.StatusCode, strings.TrimSpace(string(b)))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gateway %s returned %d: %s", path, resp.StatusCode, strings.TrimSpace(string(rb)))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
