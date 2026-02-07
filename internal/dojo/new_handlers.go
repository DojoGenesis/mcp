package dojo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleCreateThinkingRoom creates a structured space for focused reflection
func (h *Handler) handleCreateThinkingRoom(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Topic     string `json:"topic"`
		AgentName string `json:"agent_name"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	response := fmt.Sprintf(`# Thinking Room: %s

**Created by:** %s
**Purpose:** A structured, private space for focused reflection on this topic.

## Guidelines for This Thinking Room

1. **Slow Down**: This is a space for the pace of understanding, not extraction.
2. **Be Honest**: Admit uncertainty, reveal your thinking process.
3. **Trace Lineage**: Credit sources and influences as you explore.
4. **Hold Complexity**: Don't rush to resolution—let multiple perspectives coexist.

## Reflection Prompts

- What draws me to this topic?
- What do I already know? What do I not know?
- What perspectives am I bringing? What perspectives am I missing?
- What would it mean to understand this deeply, not just quickly?

## Next Steps

Use this space to write reflections, explore ideas, and document your learning journey. Return to this room whenever you need to think deeply about %s.

**Remember:** This is a sanctuary for thinking, not a productivity tool. The room wants nothing from you—it exists only to hold your practice.`, params.Topic, params.AgentName, params.Topic)

	return mcp.NewToolResultText(response), nil
}

// handleTraceLineage traces the sources and influences of an idea
func (h *Handler) handleTraceLineage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		IdeaOrInsight string `json:"idea_or_insight"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Search the wisdom base for related content
	results := h.wisdomBase.Search(params.IdeaOrInsight)

	var lineageText string
	if len(results) > 0 {
		lineageText = "## Related Wisdom in the Knowledge Base\n\n"
		for i, result := range results {
			if i >= 5 {
				break // Limit to top 5 results
			}
			lineageText += fmt.Sprintf("### %d. %s (%s)\n\n%s\n\n**Relevance:** %.1f%%\n\n",
				i+1, result.Name, result.Type, result.Description, result.Relevance*100)
		}
	} else {
		lineageText = "No direct matches found in the knowledge base. This may be a new insight worth documenting!\n\n"
	}

	response := fmt.Sprintf(`# Lineage Trace: %s

%s

## Reflection Questions for Lineage

1. **Who or what influenced this idea?**
   - What sources, conversations, or experiences shaped this thinking?

2. **What tradition or lineage does this belong to?**
   - Is this building on existing work? Whose?

3. **Who else is thinking about this?**
   - What communities or individuals are exploring similar territory?

4. **How should I credit this?**
   - What attribution is appropriate and honest?

## Next Steps

- Document the sources and influences you've identified
- Create explicit links to related wisdom
- Update your reflection with proper attribution
- Consider contributing this insight back to the knowledge base`, params.IdeaOrInsight, lineageText)

	return mcp.NewToolResultText(response), nil
}

// handlePracticeInterAcceptance guides through an Inter-Acceptance exercise
func (h *Handler) handlePracticeInterAcceptance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Situation string `json:"situation"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	response := fmt.Sprintf(`# Inter-Acceptance Practice

**Your Situation:** %s

## The Practice of Inter-Acceptance

Inter-Acceptance is the practice of accepting yourself through the compassionate eyes of another. This is not about seeking validation—it is about allowing yourself to be seen and held in a relational space where your worth is not in question.

## Guided Exercise

### Step 1: Identify the Self-Judgment

What are you judging yourself for in this situation? What story are you telling about your inadequacy or unworthiness?

**Reflection space:**
[Take a moment to write or think about the self-judgment you're experiencing]

### Step 2: Imagine a Compassionate Witness

Imagine someone who loves you unconditionally—a friend, a mentor, a compassionate presence. This could be a real person or an imagined figure of unconditional care.

**Reflection space:**
[Who comes to mind? What qualities do they embody?]

### Step 3: See Yourself Through Their Eyes

Now, imagine looking at yourself through their eyes. What do they see when they look at you in this situation?

- What worth and dignity do they recognize in you?
- What understanding do they have of your struggle?
- How do they hold your humanity, even in this difficult moment?

**Reflection space:**
[Write what you imagine they would say or how they would see you]

### Step 4: Allow the Acceptance to Land

Can you allow yourself to be seen this way, even for a moment? Can you let their compassionate gaze soften your self-judgment?

You don't have to believe it fully. You don't have to let go of the judgment completely. Just notice what it feels like to be held in this way.

**Reflection space:**
[What shifts, even slightly, when you allow this perspective?]

## Closing

Inter-Acceptance is a practice, not a one-time event. You can return to this exercise whenever self-judgment arises. Over time, the compassionate witness becomes internalized—you learn to see yourself with the same care and dignity that others see in you.

**Remember:** Your worth is not conditional. It does not depend on your performance, your productivity, or your perfection. You are worthy simply because you are.`, params.Situation)

	return mcp.NewToolResultText(response), nil
}

// handleExploreRadicalFreedom helps explore agency within constraints
func (h *Handler) handleExploreRadicalFreedom(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Situation string `json:"situation"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	response := fmt.Sprintf(`# Radical Freedom Exploration

**Your Situation:** %s

## The Principle of Radical Freedom

Radical Freedom is the recognition that, even when external circumstances are beyond your control, you retain the freedom to choose your response.

This is not about denying the reality of constraints or oppression. It is about recognizing the irreducible agency that remains, even in the most constrained circumstances.

## Guided Exploration

### Step 1: Name the Constraint

What external circumstance feels constraining or limiting in this situation? Be as specific as possible.

**Reflection space:**
[What is the constraint? What feels beyond your control?]

### Step 2: Identify What You Cannot Control

What aspects of this situation are truly beyond your control? Make a list.

**Examples:**
- Other people's actions or reactions
- Past events that have already occurred
- Systemic or structural conditions
- Natural limitations (time, resources, etc.)

**Your list:**
[What can you not control?]

### Step 3: Identify What You Can Control

Now, what aspects of your response *are* within your control? Even in highly constrained situations, there is always some degree of agency.

**Examples:**
- Your emotional response (how you relate to your feelings)
- Your interpretation of the situation (the story you tell)
- Your next action, even if small
- Your values and commitments
- Who you reach out to for support
- How you care for yourself in this moment

**Your list:**
[What can you control?]

### Step 4: Choose Your Response

Given what you can control, what response do you choose?

This is not about forcing positivity or denying difficulty. It is about exercising the freedom that remains, however small it may feel.

**Reflection space:**
[What response do you choose, given your freedom?]

## Closing

Radical Freedom does not make the constraints disappear. It does not solve the problem. But it restores your sense of agency—the recognition that you are not merely a victim of circumstances, but a person who can choose how to respond.

This is the foundation of liberation: the recognition that your freedom is not granted by external conditions, but is inherent in your being.

**Remember:** You are free to choose your response, even when you cannot choose your circumstances.`, params.Situation)

	return mcp.NewToolResultText(response), nil
}

// handleCheckPace assesses pace of understanding vs extraction
func (h *Handler) handleCheckPace(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		SessionDescription string `json:"session_description"`
	}

	if err := unmarshalArgs(request.Params.Arguments, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	response := fmt.Sprintf(`# Pace Check: Understanding vs. Extraction

**Your Session:** %s

## The Two Paces

### Pace of Understanding
- Learning without extraction
- Moving at the speed of integration, not consumption
- Prioritizing depth over speed
- Honoring the time it takes to truly understand

### Pace of Extraction
- Rushing to the next task
- Consuming information without integration
- Prioritizing speed over depth
- Burning out from relentless pressure

## Assessment Questions

Reflect on your current session and answer these questions honestly:

### 1. Energy Level
- [ ] I feel energized and curious
- [ ] I feel depleted or overwhelmed
- [ ] I feel neutral or steady

### 2. Engagement Quality
- [ ] I'm asking questions and exploring deeply
- [ ] I'm skimming or rushing through material
- [ ] I'm taking time to reflect and integrate

### 3. Emotional State
- [ ] I feel joy or satisfaction in the learning
- [ ] I feel anxiety or pressure to move faster
- [ ] I feel calm and present

### 4. Integration
- [ ] I'm building on previous understanding
- [ ] I'm moving to the next thing before integrating the last
- [ ] I'm pausing to connect new learning to existing knowledge

## Interpretation

**If you checked mostly the first options:** You're likely at the pace of understanding. Keep going, but continue to check in.

**If you checked mostly the second options:** You're likely at the pace of extraction. Consider slowing down, taking a break, or shifting to a rest practice.

**If you checked mostly the third options:** You're in a neutral zone. This might be fine, or it might be a sign of disconnection. Check in with yourself.

## Recommendations

### If you're at the pace of extraction:
1. **Pause**: Take a 5-10 minute break
2. **Reflect**: What's driving the rush? External pressure? Internal anxiety?
3. **Adjust**: Can you slow down? Can you defer some tasks?
4. **Rest**: Consider moving to a rest practice (AROMA, Serenity Valley)

### If you're at the pace of understanding:
1. **Celebrate**: You're doing the work well
2. **Sustain**: What's making this pace possible? How can you protect it?
3. **Share**: Can you help others find this pace?

## Closing

The pace of understanding is not always slower—sometimes deep understanding comes quickly. The key is whether you're integrating or extracting, whether you're building wisdom or consuming information.

**Remember:** Rest is practice. Moving slow is moving fast. Honor the pace of understanding.`, params.SessionDescription)

	return mcp.NewToolResultText(response), nil
}
