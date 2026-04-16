package dojo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DojoGenesis/mcp-server/internal/decisions"
	"github.com/DojoGenesis/mcp-server/internal/gateway"
	"github.com/DojoGenesis/mcp-server/internal/skills"
	"github.com/DojoGenesis/mcp-server/internal/wisdom"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handler manages all Dojo MCP capabilities.
type Handler struct {
	wisdomBase     *wisdom.Base
	skillsLoader   *skills.Loader
	decisionWriter *decisions.Writer
	gw             *gateway.Client
}

// NewHandler creates a new Dojo handler with skills loading and decision writing.
// gw may be nil; tools will fall back to local data when gw is nil or offline.
func NewHandler(skillsPath, adrPath string, gw *gateway.Client) (*Handler, error) {
	loader, err := skills.NewLoader(skillsPath)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}

	writer, err := decisions.NewWriter(adrPath)
	if err != nil {
		return nil, fmt.Errorf("init decision writer: %w", err)
	}

	return &Handler{
		wisdomBase:     wisdom.NewBase(),
		skillsLoader:   loader,
		decisionWriter: writer,
		gw:             gw,
	}, nil
}

// unmarshalArgs is a helper to convert arguments (any type) to a typed struct.
func unmarshalArgs(arguments any, dest interface{}) error {
	data, err := json.Marshal(arguments)
	if err != nil {
		return fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return json.Unmarshal(data, dest)
}

// RegisterTools registers all 24 Dojo tools with the MCP server.
func (h *Handler) RegisterTools(s *server.MCPServer) {
	// Tool 1: dojo.scout
	s.AddTool(mcp.Tool{
		Name:        "dojo_scout",
		Description: "Strategic analysis scaffold. Returns a 4-step framework for navigating strategic decisions: identify tension, scout routes, synthesize, and decide. Claude fills in the judgment.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"situation": map[string]interface{}{
					"type":        "string",
					"description": "The strategic situation, decision, or tension to analyze",
				},
			},
			Required: []string{"situation"},
		},
	}, h.handleScout)

	// Tool 2: dojo.invoke_skill
	s.AddTool(mcp.Tool{
		Name:        "dojo_invoke_skill",
		Description: "Load a specific methodology skill by name. Returns the full workflow as actionable steps. Use dojo.search_skills or dojo.list_skills to find skill names.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "The exact name of the skill to invoke (e.g., 'strategic-scout', 'debugging', 'retrospective')",
				},
			},
			Required: []string{"name"},
		},
	}, h.handleInvokeSkill)

	// Tool 3: dojo.search_skills
	s.AddTool(mcp.Tool{
		Name:        "dojo_search_skills",
		Description: "Search the methodology library for skills matching a query. Returns top matches with descriptions and usage guidance.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "What you're looking for (e.g., 'how to debug', 'write a spec', 'run a retro')",
				},
			},
			Required: []string{"query"},
		},
	}, h.handleSearchSkills)

	// Tool 4: dojo.apply_seed
	s.AddTool(mcp.Tool{
		Name:        "dojo_apply_seed",
		Description: "Apply a Dojo seed patch (reusable thinking pattern) to a specific situation. Returns the seed's core insight formatted as active guidance with a checklist.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"seed_name": map[string]interface{}{
					"type":        "string",
					"description": "The name of the seed patch to apply",
				},
				"situation": map[string]interface{}{
					"type":        "string",
					"description": "The situation to apply the seed patch to",
				},
			},
			Required: []string{"seed_name", "situation"},
		},
	}, h.handleApplySeed)

	// Tool 5: dojo.log_decision
	s.AddTool(mcp.Tool{
		Name:        "dojo_log_decision",
		Description: "Write an Architecture Decision Record (ADR) to disk. The only write-capable tool. Creates a persistent markdown artifact capturing a key decision with context, rationale, and consequences.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Short title for the decision (e.g., 'Use PostgreSQL for persistence')",
				},
				"context": map[string]interface{}{
					"type":        "string",
					"description": "Why this decision is needed. What forces are at play?",
				},
				"decision": map[string]interface{}{
					"type":        "string",
					"description": "What was decided and why",
				},
				"consequences": map[string]interface{}{
					"type":        "string",
					"description": "What changes as a result. Both positive and negative.",
				},
			},
			Required: []string{"title", "context", "decision", "consequences"},
		},
	}, h.handleLogDecision)

	// Tool 6: dojo.reflect
	s.AddTool(mcp.Tool{
		Name:        "dojo_reflect",
		Description: "Structured reflection grounded in your methodology library. Searches skills and seeds matching your session description and returns relevant frameworks, patterns, and reflection questions.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"session_description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the current session or work being done",
				},
			},
			Required: []string{"session_description"},
		},
	}, h.handleReflect)

	// Tool 7: dojo.list_skills
	s.AddTool(mcp.Tool{
		Name:        "dojo_list_skills",
		Description: "List available methodology skills grouped by plugin category. Supports optional filtering by plugin name and pagination via limit/offset. Shows skill names and descriptions for use with dojo.invoke_skill.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"plugin": map[string]interface{}{
					"type":        "string",
					"description": "Optional plugin name to filter results (e.g., 'strategic-thinking'). If omitted, all plugins are listed.",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of skills to return (default 50).",
				},
				"offset": map[string]interface{}{
					"type":        "integer",
					"description": "Number of skills to skip before returning results (default 0).",
				},
			},
		},
	}, h.handleListSkills)

	// ─── Memory tools (8-10) ────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_memory_list",
		Description: "List all memory entries from the Gateway. Returns ID, type, and content for each memory.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleMemoryList)

	s.AddTool(mcp.Tool{
		Name:        "dojo_memory_store",
		Description: "Store a new memory entry in the Gateway. Accepts content text and optional type.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"content": map[string]interface{}{"type": "string", "description": "The memory content to store"},
				"type":    map[string]interface{}{"type": "string", "description": "Optional memory type (e.g., 'feedback', 'project', 'user')"},
			},
			Required: []string{"content"},
		},
	}, h.handleMemoryStore)

	s.AddTool(mcp.Tool{
		Name:        "dojo_memory_search",
		Description: "Search memory entries by query. Returns matching memories ranked by relevance.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{"type": "string", "description": "Search query for memory content"},
			},
			Required: []string{"query"},
		},
	}, h.handleMemorySearch)

	// ─── Seed tools (11-13) ─────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_seed_list",
		Description: "List all seeds from the Gateway garden, or hardcoded seeds if offline.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleSeedList)

	s.AddTool(mcp.Tool{
		Name:        "dojo_seed_create",
		Description: "Plant a new seed in the Gateway garden. Requires a name and content.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name":        map[string]interface{}{"type": "string", "description": "Seed name (short, hyphenated)"},
				"content":     map[string]interface{}{"type": "string", "description": "The seed content — the reusable pattern or insight"},
				"description": map[string]interface{}{"type": "string", "description": "Optional one-line description"},
			},
			Required: []string{"name", "content"},
		},
	}, h.handleSeedCreate)

	s.AddTool(mcp.Tool{
		Name:        "dojo_seed_search",
		Description: "Search seeds by keyword in name and content.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{"type": "string", "description": "Search query for seeds"},
			},
			Required: []string{"query"},
		},
	}, h.handleSeedSearch)

	// ─── Agent tools (14-16) ────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_agent_list",
		Description: "List all agents registered with the Gateway, including their mode, status, and disposition.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleAgentList)

	s.AddTool(mcp.Tool{
		Name:        "dojo_agent_dispatch",
		Description: "Create and dispatch a new agent with an initial task message. Returns the agent ID and response.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name":    map[string]interface{}{"type": "string", "description": "Agent name"},
				"mode":    map[string]interface{}{"type": "string", "description": "Agent mode (e.g., 'focused', 'balanced', 'exploratory')"},
				"message": map[string]interface{}{"type": "string", "description": "Initial task message for the agent"},
			},
			Required: []string{"name", "message"},
		},
	}, h.handleAgentDispatch)

	s.AddTool(mcp.Tool{
		Name:        "dojo_agent_chat",
		Description: "Send a follow-up message to an existing agent by ID.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"agent_id": map[string]interface{}{"type": "string", "description": "The agent ID to send the message to"},
				"message":  map[string]interface{}{"type": "string", "description": "The message to send"},
			},
			Required: []string{"agent_id", "message"},
		},
	}, h.handleAgentChat)

	// ─── Project tools (17-19) ──────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_project_status",
		Description: "Get the current project status: name, phase, tracks, and decisions. Local-only, no Gateway needed.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleProjectStatus)

	s.AddTool(mcp.Tool{
		Name:        "dojo_project_track",
		Description: "Add or update a project track. Local-only.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{"type": "string", "description": "Action: 'add' or 'set'"},
				"name":   map[string]interface{}{"type": "string", "description": "Track name"},
				"status": map[string]interface{}{"type": "string", "description": "Track status (pending, in_progress, done, blocked)"},
			},
			Required: []string{"action", "name"},
		},
	}, h.handleProjectTrack)

	s.AddTool(mcp.Tool{
		Name:        "dojo_project_decision",
		Description: "Record a decision in the active project. Local-only.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"text": map[string]interface{}{"type": "string", "description": "The decision text to record"},
			},
			Required: []string{"text"},
		},
	}, h.handleProjectDecision)

	// ─── Disposition tools (20-21) ──────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_disposition_list",
		Description: "List available ADA disposition presets with their pacing, depth, tone, and initiative axes.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleDispositionList)

	// ─── Craft tools (22-24) ────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "dojo_converge",
		Description: "Run a convergence check: git dirty files, recent commits, memory/seed counts. Returns RED/YELLOW/GREEN signal with guidance.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleConverge)

	s.AddTool(mcp.Tool{
		Name:        "dojo_health",
		Description: "Check Gateway health: status, version, providers, uptime.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]interface{}{}},
	}, h.handleHealth)
}

// RegisterResources registers MCP resources for seeds, resources, and skills.
func (h *Handler) RegisterResources(s *server.MCPServer) {
	// Register existing wisdom resources
	resources := h.wisdomBase.ListResources()
	for _, resource := range resources {
		resourceCopy := resource
		s.AddResource(mcp.Resource{
			URI:         fmt.Sprintf("dojo://resources/%s", resourceCopy.Name),
			Name:        resourceCopy.Name,
			Description: resourceCopy.Description,
			MIMEType:    "text/markdown",
		}, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			content, err := h.wisdomBase.GetResource(resourceCopy.Name)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			}, nil
		})
	}

	// Register seed resources
	seeds := h.wisdomBase.ListSeeds()
	for _, seed := range seeds {
		seedCopy := seed
		s.AddResource(mcp.Resource{
			URI:         fmt.Sprintf("dojo://seeds/%s", seedCopy.Name),
			Name:        seedCopy.Name,
			Description: seedCopy.Description,
			MIMEType:    "text/markdown",
		}, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			sd, err := h.wisdomBase.GetSeed(seedCopy.Name)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "text/markdown",
					Text:     sd.Content,
				},
			}, nil
		})
	}

	// Skills are accessed exclusively via tools (dojo_search_skills, dojo_list_skills,
	// dojo_invoke_skill). No MCP resource registration is performed for skills.
}

// --- Tool Handlers ---

func (h *Handler) handleScout(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Situation string `json:"situation"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Situation == "" {
		return mcp.NewToolResultError("'situation' is required and cannot be empty"), nil
	}

	// If gateway is online, use ChatSync with the scout system prompt.
	if h.gw != nil && h.gw.IsOnline(ctx) {
		prompt := fmt.Sprintf(`You are a strategic analysis assistant. Apply the 4-step Dojo scout framework to this situation:

Situation: %s

Provide:
1. Identify the core tension (competing forces)
2. Scout 3-5 distinct routes/approaches with risk and tradeoff
3. Synthesize and recommend a hybrid approach
4. Decide: selected route, first action, review point

Be concise and actionable.`, args.Situation)

		result, err := h.gw.ChatSync(ctx, gateway.ChatRequest{Message: prompt})
		if err == nil && result != "" {
			return mcp.NewToolResultText(result), nil
		}
		// Fall through to scaffold on error
	}

	// Offline fallback: local scaffold
	matched := h.skillsLoader.Search(args.Situation, 3)
	return mcp.NewToolResultText(scoutScaffold(args.Situation, matched)), nil
}

func (h *Handler) handleInvokeSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Name string `json:"name"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Name == "" {
		return mcp.NewToolResultError("'name' is required and cannot be empty"), nil
	}

	// If gateway is online, try to find the skill there first.
	if h.gw != nil && h.gw.IsOnline(ctx) {
		gwSkills, err := h.gw.SearchSkills(ctx, args.Name)
		if err == nil {
			for _, s := range gwSkills {
				if strings.EqualFold(s.Name, args.Name) {
					response := fmt.Sprintf(`# Skill: %s
**Plugin:** %s
**When to use:** %s

---

**Next step:** Invoke this skill through the Dojo Genesis Gateway.`,
						s.Name, s.Plugin, firstSentence(s.Description))
					return mcp.NewToolResultText(response), nil
				}
			}
		}
		// Not found in gateway — fall through to local
	}

	skill, err := h.skillsLoader.GetByName(args.Name)
	if err != nil {
		// Suggest similar skills
		similar := h.skillsLoader.Search(args.Name, 3)
		var suggestions string
		if len(similar) > 0 {
			names := make([]string, len(similar))
			for i, s := range similar {
				names[i] = s.Name
			}
			suggestions = fmt.Sprintf("\n\nDid you mean one of these?\n- %s", strings.Join(names, "\n- "))
		}
		return mcp.NewToolResultError(fmt.Sprintf("Skill not found: %s%s", args.Name, suggestions)), nil
	}

	response := fmt.Sprintf(`# Skill: %s
**Plugin:** %s
**When to use:** %s

---

%s

---

**Next step:** Review the workflow above and begin with Step 1.`, skill.Name, skill.Plugin, firstSentence(skill.Description), skill.Content)

	return mcp.NewToolResultText(response), nil
}

func (h *Handler) handleSearchSkills(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// If gateway is online, try gateway search first.
	if h.gw != nil && h.gw.IsOnline(ctx) {
		gwSkills, err := h.gw.SearchSkills(ctx, args.Query)
		if err == nil && len(gwSkills) > 0 {
			var sb strings.Builder
			fmt.Fprintf(&sb, "# Skills matching: \"%s\"\n\n", args.Query)
			fmt.Fprintf(&sb, "Found %d relevant skill(s) (from gateway):\n\n", len(gwSkills))
			for i, s := range gwSkills {
				fmt.Fprintf(&sb, "## %d. %s (%s)\n", i+1, s.Name, s.Plugin)
				fmt.Fprintf(&sb, "**When to use:** %s\n", firstSentence(s.Description))
				fmt.Fprintf(&sb, "**Invoke:** `dojo.invoke_skill` with name \"%s\"\n\n", s.Name)
			}
			return mcp.NewToolResultText(sb.String()), nil
		}
		// Fall through to local search on error or empty results
	}

	results, totalMatching := h.skillsLoader.SearchWithTotal(args.Query, 20)

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No skills found matching: \"%s\"\n\nTry broader terms, or use `dojo.list_skills` to see all available skills.", args.Query)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Skills matching: \"%s\"\n\n", args.Query)
	fmt.Fprintf(&sb, "Found %d relevant skill(s):\n\n", len(results))

	for i, skill := range results {
		fmt.Fprintf(&sb, "## %d. %s (%s)\n", i+1, skill.Name, skill.Plugin)
		fmt.Fprintf(&sb, "**When to use:** %s\n", skill.Description)
		fmt.Fprintf(&sb, "**Invoke:** `dojo.invoke_skill` with name \"%s\"\n\n", skill.Name)
	}

	fmt.Fprintf(&sb, "(showing top %d of %d matches for \"%s\")\n", len(results), totalMatching, args.Query)

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleApplySeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SeedName  string `json:"seed_name"`
		Situation string `json:"situation"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.SeedName == "" {
		return mcp.NewToolResultError("'seed_name' is required"), nil
	}
	if args.Situation == "" {
		return mcp.NewToolResultError("'situation' is required"), nil
	}

	seed, err := h.wisdomBase.GetSeed(args.SeedName)
	if err != nil {
		// List available seeds
		seeds := h.wisdomBase.ListSeeds()
		names := make([]string, len(seeds))
		for i, s := range seeds {
			names[i] = s.Name
		}
		return mcp.NewToolResultError(fmt.Sprintf("Seed '%s' not found. Available seeds:\n- %s", args.SeedName, strings.Join(names, "\n- "))), nil
	}

	response := fmt.Sprintf(`# Applying Seed: %s

**Your situation:** %s

## Core Insight

%s

## Application Checklist

%s

## Reflection Questions

1. Which aspects of this seed are most relevant to your situation?
2. What would successful application look like?
3. What's the smallest step to begin?`, seed.Name, args.Situation, firstSentence(seed.Content), seed.Content)

	return mcp.NewToolResultText(response), nil
}

func (h *Handler) handleLogDecision(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Title        string `json:"title"`
		Context      string `json:"context"`
		Decision     string `json:"decision"`
		Consequences string `json:"consequences"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Title == "" {
		return mcp.NewToolResultError("'title' is required"), nil
	}
	if args.Context == "" {
		return mcp.NewToolResultError("'context' is required"), nil
	}
	if args.Decision == "" {
		return mcp.NewToolResultError("'decision' is required"), nil
	}
	if args.Consequences == "" {
		return mcp.NewToolResultError("'consequences' is required"), nil
	}

	fp, err := h.decisionWriter.LogDecision(args.Title, args.Context, args.Decision, args.Consequences)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write ADR: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("ADR written to: %s", fp)), nil
}

func (h *Handler) handleReflect(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SessionDescription string `json:"session_description"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.SessionDescription == "" {
		return mcp.NewToolResultError("'session_description' is required"), nil
	}

	// Search skills by session description
	matchedSkills := h.skillsLoader.Search(args.SessionDescription, 3)

	// Search seeds by session description
	seedResults := h.wisdomBase.Search(args.SessionDescription)
	var matchedSeeds []wisdom.Seed
	for _, sr := range seedResults {
		if sr.Type != "seed" {
			continue
		}
		s, err := h.wisdomBase.GetSeed(sr.Name)
		if err == nil && s != nil {
			matchedSeeds = append(matchedSeeds, *s)
			if len(matchedSeeds) >= 3 {
				break
			}
		}
	}

	return mcp.NewToolResultText(reflectScaffold(args.SessionDescription, matchedSkills, matchedSeeds)), nil
}

func (h *Handler) handleListSkills(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Plugin string `json:"plugin"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// Apply defaults
	if args.Limit <= 0 {
		args.Limit = 50
	}
	if args.Offset < 0 {
		args.Offset = 0
	}

	// If gateway is online, use gateway skills list.
	if h.gw != nil && h.gw.IsOnline(ctx) {
		gwSkills, err := h.gw.Skills(ctx)
		if err == nil && len(gwSkills) > 0 {
			// Apply plugin filter
			var filtered []gateway.Skill
			for _, s := range gwSkills {
				if args.Plugin == "" || strings.EqualFold(s.Plugin, args.Plugin) {
					filtered = append(filtered, s)
				}
			}
			total := len(filtered)
			start := args.Offset
			if start > total {
				start = total
			}
			end := start + args.Limit
			if end > total {
				end = total
			}
			page := filtered[start:end]

			var sb strings.Builder
			if args.Plugin != "" {
				fmt.Fprintf(&sb, "# Dojo Genesis Skills — plugin: %s (%d available, from gateway)\n\n", args.Plugin, total)
			} else {
				fmt.Fprintf(&sb, "# Dojo Genesis Skills (%d available, from gateway)\n\n", total)
			}
			currentPlugin := ""
			for _, s := range page {
				if s.Plugin != currentPlugin {
					if currentPlugin != "" {
						sb.WriteString("\n")
					}
					fmt.Fprintf(&sb, "## %s\n", s.Plugin)
					currentPlugin = s.Plugin
				}
				fmt.Fprintf(&sb, "- **%s** -- %s\n", s.Name, firstSentence(s.Description))
			}
			if len(page) > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString("Use `dojo.invoke_skill` with the skill name to load the full workflow.\n")
			sb.WriteString("Use `dojo.search_skills` with a query to find the right skill for your task.\n")
			if total == 0 {
				sb.WriteString("(no skills found)\n")
			} else {
				fmt.Fprintf(&sb, "(showing %d\u2013%d of %d skills)\n", start+1, end, total)
			}
			return mcp.NewToolResultText(sb.String()), nil
		}
		// Fall through to local list on error or empty
	}

	var sb strings.Builder

	// Collect the flat list of skills to paginate, optionally filtered by plugin
	byPlugin := h.skillsLoader.ListByPlugin()
	var pluginNames []string
	if args.Plugin != "" {
		if _, ok := byPlugin[args.Plugin]; ok {
			pluginNames = []string{args.Plugin}
		}
	} else {
		pluginNames = h.skillsLoader.PluginNames()
	}

	// Build a flat ordered slice for pagination
	type entry struct {
		plugin string
		name   string
		desc   string
	}
	var allEntries []entry
	for _, pName := range pluginNames {
		for _, s := range byPlugin[pName] {
			allEntries = append(allEntries, entry{plugin: pName, name: s.Name, desc: s.Description})
		}
	}

	total := len(allEntries)
	start := args.Offset
	if start > total {
		start = total
	}
	end := start + args.Limit
	if end > total {
		end = total
	}
	page := allEntries[start:end]

	if args.Plugin != "" {
		fmt.Fprintf(&sb, "# Dojo Genesis Skills — plugin: %s (%d available)\n\n", args.Plugin, total)
	} else {
		fmt.Fprintf(&sb, "# Dojo Genesis Skills (%d available)\n\n", h.skillsLoader.Count())
	}

	// Group the page entries by plugin for readability
	currentPlugin := ""
	for _, e := range page {
		if e.plugin != currentPlugin {
			if currentPlugin != "" {
				sb.WriteString("\n")
			}
			pluginSkillCount := len(byPlugin[e.plugin])
			fmt.Fprintf(&sb, "## %s (%d skills)\n", e.plugin, pluginSkillCount)
			currentPlugin = e.plugin
		}
		fmt.Fprintf(&sb, "- **%s** -- %s\n", e.name, firstSentence(e.desc))
	}
	if len(page) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("Use `dojo.invoke_skill` with the skill name to load the full workflow.\n")
	sb.WriteString("Use `dojo.search_skills` with a query to find the right skill for your task.\n")
	if total == 0 {
		sb.WriteString("(no skills found)\n")
	} else {
		fmt.Fprintf(&sb, "(showing %d\u2013%d of %d skills)\n", start+1, end, total)
	}

	return mcp.NewToolResultText(sb.String()), nil
}
