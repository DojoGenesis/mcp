package dojo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TresPies-source/dojo-mcp-server/internal/wisdom"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handler manages all Dojo-specific MCP capabilities
type Handler struct {
	wisdomBase *wisdom.Base
}

// NewHandler creates a new Dojo handler
func NewHandler() *Handler {
	return &Handler{
		wisdomBase: wisdom.NewBase(),
	}
}

// unmarshalArgs is a helper to convert map[string]interface{} arguments to a typed struct
func unmarshalArgs(arguments map[string]interface{}, dest interface{}) error {
	data, err := json.Marshal(arguments)
	if err != nil {
		return fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return json.Unmarshal(data, dest)
}

// RegisterTools registers all Dojo tools with the MCP server
func (h *Handler) RegisterTools(s *server.MCPServer) {
	// dojo.reflect - The core Dojo thinking partner
	s.AddTool(mcp.Tool{
		Name:        "dojo.reflect",
		Description: "The core Dojo thinking partner. Applies one of the four Dojo modes (Mirror, Scout, Gardener, Implementation) to a given situation and set of perspectives.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"situation": map[string]interface{}{
					"type":        "string",
					"description": "The situation or question to reflect on",
				},
				"perspectives": map[string]interface{}{
					"type":        "array",
					"description": "A list of perspectives to consider",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"mode": map[string]interface{}{
					"type":        "string",
					"description": "The Dojo mode to apply: mirror, scout, gardener, or implementation",
					"enum":        []string{"mirror", "scout", "gardener", "implementation"},
				},
			},
			Required: []string{"situation", "perspectives", "mode"},
		},
	}, h.handleReflect)

	// dojo.search_wisdom - Semantic search on the Dojo wisdom base
	s.AddTool(mcp.Tool{
		Name:        "dojo.search_wisdom",
		Description: "Performs a semantic search on the entire Dojo wisdom base, including all seed patches, documentation, and principles.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The search query",
				},
			},
			Required: []string{"query"},
		},
	}, h.handleSearchWisdom)

	// dojo.get_seed - Retrieve a specific Dojo Seed Patch
	s.AddTool(mcp.Tool{
		Name:        "dojo.get_seed",
		Description: "Retrieves a specific Dojo Seed Patch by name.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "The name of the seed patch (e.g., 'three_tiered_governance')",
				},
			},
			Required: []string{"name"},
		},
	}, h.handleGetSeed)

	// dojo.apply_seed - Apply a Dojo Seed Patch to a situation
	s.AddTool(mcp.Tool{
		Name:        "dojo.apply_seed",
		Description: "Applies a Dojo Seed Patch to a given situation, providing guidance and a checklist.",
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

	// dojo.list_seeds - List all available Dojo Seed Patches
	s.AddTool(mcp.Tool{
		Name:        "dojo.list_seeds",
		Description: "Lists all available Dojo Seed Patches with their descriptions.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.handleListSeeds)

	// dojo.get_principles - Get the core Dojo principles
	s.AddTool(mcp.Tool{
		Name:        "dojo.get_principles",
		Description: "Retrieves the three core Dojo principles: Beginner's Mind, Self-Definition, and Understanding is Love.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.handleGetPrinciples)

	// v2.0 Tools: AROMA / Serenity Valley

	// dojo.create_thinking_room - Create a structured space for focused reflection
	s.AddTool(mcp.Tool{
		Name:        "dojo.create_thinking_room",
		Description: "Creates a structured, private space for focused reflection on a given topic.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"topic": map[string]interface{}{
					"type":        "string",
					"description": "The topic to reflect on",
				},
				"agent_name": map[string]interface{}{
					"type":        "string",
					"description": "The name of the agent or user creating the room",
				},
			},
			Required: []string{"topic", "agent_name"},
		},
	}, h.handleCreateThinkingRoom)

	// dojo.trace_lineage - Trace the sources and influences of an idea
	s.AddTool(mcp.Tool{
		Name:        "dojo.trace_lineage",
		Description: "Traces the sources and influences of an idea or insight, searching the wisdom base for related content.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"idea_or_insight": map[string]interface{}{
					"type":        "string",
					"description": "The idea or insight to trace",
				},
			},
			Required: []string{"idea_or_insight"},
		},
	}, h.handleTraceLineage)

	// dojo.practice_inter_acceptance - Guided Inter-Acceptance exercise
	s.AddTool(mcp.Tool{
		Name:        "dojo.practice_inter_acceptance",
		Description: "Guides through an Inter-Acceptance exercise from Serenity Valley's Emotional Interbeing Therapy.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"situation": map[string]interface{}{
					"type":        "string",
					"description": "The situation to practice inter-acceptance with",
				},
			},
			Required: []string{"situation"},
		},
	}, h.handlePracticeInterAcceptance)

	// dojo.explore_radical_freedom - Explore agency within constraints
	s.AddTool(mcp.Tool{
		Name:        "dojo.explore_radical_freedom",
		Description: "Helps explore agency and freedom within constraints, based on Serenity Valley's Radical Freedom principle.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"situation": map[string]interface{}{
					"type":        "string",
					"description": "The constrained situation to explore",
				},
			},
			Required: []string{"situation"},
		},
	}, h.handleExploreRadicalFreedom)

	// dojo.check_pace - Assess pace of understanding vs extraction
	s.AddTool(mcp.Tool{
		Name:        "dojo.check_pace",
		Description: "Assesses whether the current session pace is one of understanding or extraction, with recommendations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"session_description": map[string]interface{}{
					"type":        "string",
					"description": "A description of the current session or work being done",
				},
			},
			Required: []string{"session_description"},
		},
	}, h.handleCheckPace)
}

// Tool handlers

func (h *Handler) handleReflect(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Situation    string   `json:"situation"`
		Perspectives []string `json:"perspectives"`
		Mode         string   `json:"mode"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	reflection := h.reflect(args.Situation, args.Perspectives, args.Mode)

	return mcp.NewToolResultText(reflection), nil
}

func (h *Handler) handleSearchWisdom(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	results := h.wisdomBase.Search(args.Query)

	resultsJSON, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(string(resultsJSON)), nil
}

func (h *Handler) handleGetSeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Name string `json:"name"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	seed, err := h.wisdomBase.GetSeed(args.Name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Seed not found: %v", err)), nil
	}

	seedJSON, _ := json.MarshalIndent(seed, "", "  ")
	return mcp.NewToolResultText(string(seedJSON)), nil
}

func (h *Handler) handleApplySeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SeedName  string `json:"seed_name"`
		Situation string `json:"situation"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	guidance := h.applySeed(args.SeedName, args.Situation)

	return mcp.NewToolResultText(guidance), nil
}

func (h *Handler) handleListSeeds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seeds := h.wisdomBase.ListSeeds()

	seedsJSON, _ := json.MarshalIndent(seeds, "", "  ")
	return mcp.NewToolResultText(string(seedsJSON)), nil
}

func (h *Handler) handleGetPrinciples(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	principles := h.wisdomBase.GetPrinciples()

	return mcp.NewToolResultText(principles), nil
}

// RegisterPrompts registers all Dojo prompts with the MCP server
func (h *Handler) RegisterPrompts(s *server.MCPServer) {
	seeds := h.wisdomBase.ListSeeds()

	for _, seed := range seeds {
		seedCopy := seed // Capture for closure
		s.AddPrompt(mcp.Prompt{
			Name:        fmt.Sprintf("dojo.seed.%s", seedCopy.Name),
			Description: seedCopy.Description,
		}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			fullSeed, _ := h.wisdomBase.GetSeed(seedCopy.Name)

			return &mcp.GetPromptResult{
				Messages: []mcp.PromptMessage{
					{
						Role: "user",
						Content: mcp.TextContent{
							Type: "text",
							Text: fullSeed.Content,
						},
					},
				},
			}, nil
		})
	}
}

// RegisterResources registers all Dojo resources with the MCP server
func (h *Handler) RegisterResources(s *server.MCPServer) {
	resources := h.wisdomBase.ListResources()

	for _, resource := range resources {
		resourceCopy := resource // Capture for closure
		s.AddResource(mcp.Resource{
			URI:         fmt.Sprintf("dojo://%s", resourceCopy.Name),
			Name:        resourceCopy.Name,
			Description: resourceCopy.Description,
			MIMEType:    "text/markdown",
		}, func(ctx context.Context, request mcp.ReadResourceRequest) ([]interface{}, error) {
			content, err := h.wisdomBase.GetResource(resourceCopy.Name)
			if err != nil {
				return nil, err
			}

			return []interface{}{
				mcp.TextResourceContents{
					ResourceContents: mcp.ResourceContents{
						URI:      request.Params.URI,
						MIMEType: "text/markdown",
					},
					Text: content,
				},
			}, nil
		})
	}
}

// reflect implements the core Dojo reflection logic
func (h *Handler) reflect(situation string, perspectives []string, mode string) string {
	switch mode {
	case "mirror":
		return h.mirrorMode(situation, perspectives)
	case "scout":
		return h.scoutMode(situation, perspectives)
	case "gardener":
		return h.gardenerMode(situation, perspectives)
	case "implementation":
		return h.implementationMode(situation, perspectives)
	default:
		return "Unknown mode. Please use: mirror, scout, gardener, or implementation."
	}
}

func (h *Handler) mirrorMode(situation string, perspectives []string) string {
	return fmt.Sprintf(`**MIRROR MODE**

Situation: %s

Perspectives provided: %d

**Pattern across perspectives:**
The perspectives reveal a multi-faceted view of the situation. Each perspective brings a unique lens, and together they form a more complete picture than any single view could provide.

**Assumptions/tensions identified:**
1. There may be an implicit assumption that one perspective is "correct" while others are less valid.
2. Tension exists between different priorities and values embedded in each perspective.

**Reframe:**
What if this situation doesn't require choosing one perspective over another, but rather holding all perspectives simultaneously? The richness is in the multiplicity, not in the resolution.`, situation, len(perspectives))
}

func (h *Handler) scoutMode(situation string, perspectives []string) string {
	return fmt.Sprintf(`**SCOUT MODE**

Situation: %s

Perspectives considered: %d

**Possible routes:**

1. **Explore each perspective deeply** - Spend time with each view, understanding its logic and implications.
   - Tradeoff: Takes time, but builds comprehensive understanding.

2. **Identify common ground** - Look for where perspectives overlap or agree.
   - Tradeoff: May miss important differences, but creates foundation for action.

3. **Test assumptions** - Challenge the core assumptions in each perspective.
   - Tradeoff: Can feel destabilizing, but reveals hidden constraints.

4. **Prototype small** - Pick the smallest possible test that engages multiple perspectives.
   - Tradeoff: Limited scope, but provides real data quickly.

**Recommended smallest test:**
Articulate one core question that each perspective would answer differently. See what emerges from the contrast.`, situation, len(perspectives))
}

func (h *Handler) gardenerMode(situation string, perspectives []string) string {
	return fmt.Sprintf(`**GARDENER MODE**

Situation: %s

Perspectives in the garden: %d

**Strongest ideas (ready to grow):**
1. The perspectives that are most grounded in direct experience or evidence.
2. The perspectives that open up new possibilities rather than closing them down.

**Ideas that need growth:**
1. Perspectives that are overly abstract or disconnected from the specific situation.
2. Perspectives that seem defensive or reactive rather than generative.

**Cultivation note:**
The goal is not to eliminate weak perspectives, but to strengthen them through attention and questioning. What would make each perspective more robust?`, situation, len(perspectives))
}

func (h *Handler) implementationMode(situation string, perspectives []string) string {
	return fmt.Sprintf(`**IMPLEMENTATION MODE**

Situation: %s

Perspectives integrated: %d

**Next steps:**

1. **Document the perspectives** - Write down each perspective clearly, in 1-2 sentences.

2. **Identify decision criteria** - What would make you choose one path over another?

3. **Set a decision point** - When will you commit to a direction? (e.g., "after gathering X data" or "by Y date")

4. **Design a review** - How will you know if the chosen path is working?

5. **Take the first action** - What's the smallest concrete step you can take today?

These steps keep you moving while honoring the complexity of multiple perspectives.`, situation, len(perspectives))
}

func (h *Handler) applySeed(seedName, situation string) string {
	seed, err := h.wisdomBase.GetSeed(seedName)
	if err != nil {
		return fmt.Sprintf("Seed '%s' not found.", seedName)
	}

	return fmt.Sprintf(`**Applying Seed: %s**

**Situation:** %s

**Seed Content:**
%s

**Guidance:**
Review the seed content above and consider how each principle or pattern applies to your specific situation. Use the checklist items as a guide for implementation.

**Reflection Questions:**
1. Which aspects of this seed are most relevant to your situation?
2. What would successful application of this seed look like?
3. What obstacles might prevent full application?
4. What's the smallest step you could take to begin applying this seed?`, seed.Name, situation, seed.Content)
}
