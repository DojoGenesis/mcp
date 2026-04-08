---
name: multi-surface-strategy
description: >
  Guide the design of a coherent multi-surface product strategy where each surface
  has a unique, complementary role. Use when planning a product across multiple
  surfaces (desktop, mobile, web, API, CLI), when adding a new surface to an
  existing product, or when a multi-surface product feels fragmented. Produces a
  strategy document with surfaces, jobs-to-be-done, feature map, handoff design,
  business model, and phasing plan.
---

# Multi-Surface Product Strategy Skill

**Version:** 1.0
**Purpose:** To guide the design of a coherent multi-surface product strategy where each surface (e.g., desktop, mobile, web, API) has a unique, complementary role.

---

## I. The Philosophy: Complement, Don't Compete

In a multi-surface world, the biggest mistake is to build the same product on every device. A desktop app, a mobile app, and a web app should not be clones of each other. They should be **complementary surfaces**, each optimized for the unique context in which it will be used.

Three principles guide this skill:

1. **Surfaces are for contexts, not devices.** Users don't think "I need the mobile version." They think "I need to capture this idea quickly while I'm walking." Frame each surface by its context of use, not by its device type.

2. **The handoff is the feature.** The most magical part of a multi-surface strategy is the seamless handoff between surfaces. This is what makes the whole greater than the sum of its parts. Invest heavily in sync architecture and handoff UX.

3. **Asymmetry is a feature.** Not every surface needs every feature. The best multi-surface strategies give each surface a reason to exist that the others can't replicate. Symmetrical strategies (same features everywhere) are wasteful and boring.

---

## II. When to Use This Skill

- When planning a new product that will exist on multiple surfaces
- When adding a new surface (e.g., a mobile app) to an existing product
- When a multi-surface product feels fragmented or confusing
- During a strategic review of a product line
- After using `product-positioning` to identify unique value propositions for each surface
- When deciding the rollout order for multiple surfaces

**When NOT to use:**
- When the product genuinely only needs one surface
- When the question is about positioning (use `product-positioning` first)
- When the surfaces are already well-defined and the question is about execution

---

## III. The Workflow

This is a 6-step workflow for designing a multi-surface product strategy.

### Step 1: Identify the Surfaces

**Goal:** List all current and potential product surfaces.

**Actions:**
1. List existing surfaces (web app, desktop app, etc.)
2. Brainstorm potential new surfaces (mobile app, browser extension, CLI, API, voice interface)
3. Consider unconventional surfaces that match user contexts
4. Ask the user which surfaces they're considering and why

**Output:** A complete list of surfaces to evaluate.

### Step 2: Define the Context and Job-to-be-Done for Each Surface

**Goal:** For each surface, define the primary context and the job users will hire it to do.

**Actions:**
1. For each surface, define the context: when, where, and why a user reaches for this surface
2. Write a job-to-be-done statement: "When I'm [context], I want to [action] so I can [outcome]."
3. Ask: "What is this surface uniquely good at that the others aren't?"
4. Validate that each surface has a genuinely unique job. If two surfaces have the same job, one may be unnecessary.

**Output:** A context definition and job-to-be-done statement for each surface.

**Example:**
| Surface | Context | Job-to-be-Done |
|---------|---------|---------------|
| Desktop | At my desk, sustained focus session | "When I'm in deep work mode, I want to orchestrate complex multi-step tasks so I can build something sophisticated" |
| Mobile | Walking, commuting, between meetings | "When I'm on the go, I want to quickly capture ideas and check status so I can stay connected without breaking flow" |
| Web | Any device, first encounter | "When I'm discovering the product, I want to try it without installing so I can evaluate before committing" |

### Step 3: Map Features to Surfaces

**Goal:** Map features to the surface where they best fit, based on the job-to-be-done.

**Actions:**
1. Create a table with surfaces as columns and features as rows
2. For each feature, assign: **Primary** (this is where the feature lives), **Secondary** (available but not optimized), or **Not Available** (deliberately excluded)
3. Identify features that should be surface-exclusive
4. Validate that each surface has at least one Primary feature the others don't have

**Output:** A feature-to-surface mapping table.

**Example:**
| Feature | Desktop | Mobile | Web |
|---------|---------|--------|-----|
| Complex multi-step orchestration | Primary | -- | -- |
| Quick idea capture | Secondary | Primary | Secondary |
| Status monitoring / dashboards | Secondary | Primary | -- |
| Deep configuration | Primary | -- | Secondary |
| Onboarding tutorial | -- | -- | Primary |
| Artifact rendering | Primary | Secondary | Primary |

**Key principle:** If a feature is Primary on every surface, it's not differentiated. Asymmetry is a feature.

### Step 4: Design the Handoffs

**Goal:** Design seamless transitions between surfaces.

**Actions:**
1. Define the sync architecture (cloud-based, real-time, eventual consistency)
2. For each pair of surfaces, identify the handoff moment:
   - What triggers the handoff?
   - What state transfers?
   - What does the UX moment feel like?
3. Design cross-surface awareness (notifications, "Continue on Desktop" prompts)
4. Consider what happens when surfaces are used offline or out of sync

**Output:** A handoff design document.

**Key insight:** The handoff is the feature. A "Continue on Desktop" button on mobile is the moment where multi-surface strategy becomes tangible to users. Invest in making this feel seamless.

### Step 5: Define the Business Model

**Goal:** Define pricing and packaging for the multi-surface strategy.

**Actions:**
1. Decide which surfaces are free, paid, or part of a subscription
2. Consider surface-specific pricing that reflects different value propositions
3. Define what's included in each tier
4. Consider whether surfaces are bundled or sold separately

**Output:** A business model document.

**Key insight:** Symmetrical pricing (same price for all surfaces) doesn't reflect their different value propositions. Consider asymmetric pricing that matches asymmetric value.

### Step 6: Plan the Phasing

**Goal:** Determine the build order for surfaces.

**Actions:**
1. Identify which surface is the core (build first, prove value)
2. Define the sequence for adding surfaces, based on user need and engineering capacity
3. Set milestones for each surface launch
4. Plan how to validate each surface's unique job before investing in the next

**Output:** A phased rollout plan.

**Key insight:** Building all surfaces simultaneously is expensive and risky. Start with one core surface, prove the value, then expand strategically.

---

## IV. Best Practices

### 1. Surfaces are for Contexts, Not Devices

**Why:** "Desktop app" vs. "mobile app" is a device-centric framing. "Deep work surface" vs. "on-the-go surface" is a context-centric framing. The latter produces better strategy.

**How:** Never name a surface by its device. Name it by its context.

### 2. The Handoff is the Feature

**Why:** Users don't care about your sync architecture. They care that when they switch from phone to laptop, their work is there.

**How:** Design handoffs as first-class features, not infrastructure afterthoughts. Prototype the handoff moment before building the surface.

### 3. Design for Asymmetry

**Why:** If every surface has every feature, users don't see the value in having multiple surfaces. Asymmetry creates complementary value.

**How:** For each feature, ask "which surface is this feature's natural home?" If the answer is "all of them," you need a sharper definition of each surface's job.

### 4. Start with One Surface, Expand Strategically

**Why:** Multi-surface development is expensive. Proving value on one surface reduces risk.

**How:** Launch the core surface first. Validate its job-to-be-done. Then add surfaces that complement it.

### 5. Simplicity on Each Surface

**Why:** Feature bloat on any single surface kills the magic. Each surface should be ruthlessly focused on its core job.

**How:** If a feature doesn't serve the surface's primary job-to-be-done, it's a Secondary or Not Available. Resist the temptation to add everything everywhere.

---

## V. Quality Checklist

Before delivering, verify:

- [ ] Have you identified all potential product surfaces?
- [ ] Does each surface have a unique context and job-to-be-done?
- [ ] Have you produced a feature-to-surface mapping table?
- [ ] Does at least one feature per surface have Primary status that no other surface has?
- [ ] Have you designed handoff moments between each pair of surfaces?
- [ ] Have you defined a sustainable business model?
- [ ] Have you planned a phased rollout (not all surfaces at once)?
- [ ] Is the strategy document saved as a file?

---

## VI. Common Pitfalls

### Pitfall 1: Building the Same Product on Every Surface

**Problem:** Users don't see the value in multiple surfaces if they're all clones.

**Solution:** Define unique jobs-to-be-done for each surface. Map features asymmetrically.

### Pitfall 2: Ignoring the Handoff Experience

**Problem:** Users get frustrated when they can't seamlessly move between surfaces.

**Solution:** Invest in sync architecture and handoff UX from day one. The handoff is the feature.

### Pitfall 3: Feature Bloat on Every Surface

**Problem:** Trying to add every feature to every surface leads to complexity and dilutes each surface's purpose.

**Solution:** Be ruthless about Primary vs. Secondary vs. Not Available. Each surface should focus on its core job.

### Pitfall 4: Building All Surfaces Simultaneously

**Problem:** Splits engineering focus, delays all surfaces, increases coordination costs.

**Solution:** Start with one core surface, prove its value, then expand strategically.

### Pitfall 5: Symmetrical Pricing

**Problem:** Charging the same for all surfaces doesn't reflect their different value propositions.

**Solution:** Consider asymmetric pricing. The deep-work surface might command a premium; the discovery surface might be free.

---

## VII. Example: Desktop + Mobile + Web Strategy

**Surfaces Identified:** Desktop (Electron/Tauri), Mobile (PWA), Web (existing)

**Jobs-to-Be-Done:**
| Surface | Context | Job |
|---------|---------|-----|
| Desktop | At desk, sustained focus | Deep work, complex multi-agent orchestration, sustained focus sessions |
| Mobile | On the go, between meetings | Quick task capture, status checks, lightweight orchestration |
| Web | Any device, first encounter | Discovery, onboarding, lightweight access without installation |

**Feature Map:**
| Feature | Desktop | Mobile | Web |
|---------|---------|--------|-----|
| Complex orchestration | Primary | -- | -- |
| Quick capture | Secondary | Primary | Secondary |
| Status monitoring | Secondary | Primary | -- |
| Deep configuration | Primary | -- | Secondary |
| Onboarding | -- | -- | Primary |

**Handoff Design:**
- Mobile "Continue on Desktop" button when task complexity exceeds mobile capability
- Desktop "Quick Share" to generate a web URL for sharing artifacts
- Cloud sync via backend API, real-time status updates across surfaces

**Business Model:**
- Desktop: Core product ($20/month)
- Mobile: Premium tier (separate subscription, 4-6 weeks after desktop)
- Web: Free tier for discovery and onboarding

**Phasing:** Desktop first (core value), Web second (discovery funnel), Mobile third (retention and on-the-go access).

**Outcome:** A clear, asymmetric strategy where each surface has a unique value proposition and reason to exist.

---

## VIII. Related Skills

- **`product-positioning`** -- Use first to identify the unique value of each surface
- **`strategic-scout`** -- Use to explore multiple routes for multi-surface strategy
- **`iterative-scouting`** -- Use to refine the strategy based on feedback loops
