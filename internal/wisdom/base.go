package wisdom

import (
	"fmt"
	"sort"
	"strings"
)

// Seed represents a Dojo Seed Patch
type Seed struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Triggers    string `json:"triggers"`
}

// Resource represents a Dojo documentation resource
type Resource struct {
	Name        string
	Description string
	Content     string
}

// SearchResult represents a search result from the wisdom base
type SearchResult struct {
	Type        string  `json:"type"` // "seed", "principle", "resource"
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Relevance   float64 `json:"relevance"`
	Snippet     string  `json:"snippet"`
}

// Base manages the Dojo wisdom base
type Base struct {
	seeds      []Seed
	resources  []Resource
	principles string
}

// NewBase creates a new wisdom base with all Dojo knowledge
func NewBase() *Base {
	return &Base{
		seeds:      getSeeds(),
		resources:  getResources(),
		principles: getPrinciples(),
	}
}

// Search performs a semantic search on the wisdom base
func (b *Base) Search(query string) []SearchResult {
	query = strings.ToLower(query)
	results := []SearchResult{}

	// Search seeds
	for _, seed := range b.seeds {
		relevance := calculateRelevance(query, seed.Name, seed.Description, seed.Content)
		if relevance > 0.1 {
			results = append(results, SearchResult{
				Type:        "seed",
				Name:        seed.Name,
				Description: seed.Description,
				Relevance:   relevance,
				Snippet:     getSnippet(seed.Content, query),
			})
		}
	}

	// Search resources
	for _, resource := range b.resources {
		relevance := calculateRelevance(query, resource.Name, resource.Description, resource.Content)
		if relevance > 0.1 {
			results = append(results, SearchResult{
				Type:        "resource",
				Name:        resource.Name,
				Description: resource.Description,
				Relevance:   relevance,
				Snippet:     getSnippet(resource.Content, query),
			})
		}
	}

	// Search principles
	relevance := calculateRelevance(query, "principles", "Core Dojo Principles", b.principles)
	if relevance > 0.1 {
		results = append(results, SearchResult{
			Type:        "principle",
			Name:        "Core Dojo Principles",
			Description: "The three foundational principles of Dojo",
			Relevance:   relevance,
			Snippet:     getSnippet(b.principles, query),
		})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Relevance > results[j].Relevance })
	return results
}

// GetSeed retrieves a specific seed by name
func (b *Base) GetSeed(name string) (*Seed, error) {
	for _, seed := range b.seeds {
		if seed.Name == name {
			return &seed, nil
		}
	}
	return nil, fmt.Errorf("seed not found: %s", name)
}

// ListSeeds returns all available seeds
func (b *Base) ListSeeds() []Seed {
	return b.seeds
}

// GetPrinciples returns the core Dojo principles
func (b *Base) GetPrinciples() string {
	return b.principles
}

// GetResource retrieves a specific resource by name
func (b *Base) GetResource(name string) (string, error) {
	for _, resource := range b.resources {
		if resource.Name == name {
			return resource.Content, nil
		}
	}
	return "", fmt.Errorf("resource not found: %s", name)
}

// ListResources returns all available resources
func (b *Base) ListResources() []Resource {
	return b.resources
}

// Helper functions

func calculateRelevance(query, name, description, content string) float64 {
	query = strings.ToLower(query)
	name = strings.ToLower(name)
	description = strings.ToLower(description)
	content = strings.ToLower(content)

	score := 0.0

	// Name match is highly relevant
	if strings.Contains(name, query) {
		score += 1.0
	}

	// Description match is moderately relevant
	if strings.Contains(description, query) {
		score += 0.5
	}

	// Content match is less relevant but still counts
	if strings.Contains(content, query) {
		score += 0.2
	}

	// Keyword matching
	keywords := strings.Fields(query)
	for _, keyword := range keywords {
		if len(keyword) < 3 {
			continue
		}
		if strings.Contains(name, keyword) {
			score += 0.3
		}
		if strings.Contains(description, keyword) {
			score += 0.2
		}
		if strings.Contains(content, keyword) {
			score += 0.1
		}
	}

	return score
}

func getSnippet(content, query string) string {
	query = strings.ToLower(query)
	content = strings.ToLower(content)

	// Find the first occurrence of the query
	index := strings.Index(content, query)
	if index == -1 {
		// If exact query not found, try first keyword
		keywords := strings.Fields(query)
		for _, keyword := range keywords {
			if len(keyword) < 3 {
				continue
			}
			index = strings.Index(content, keyword)
			if index != -1 {
				break
			}
		}
	}

	if index == -1 {
		// No match found, return first 200 chars
		if len(content) > 200 {
			return content[:200] + "..."
		}
		return content
	}

	// Extract snippet around the match
	start := index - 100
	if start < 0 {
		start = 0
	}
	end := index + 100
	if end > len(content) {
		end = len(content)
	}

	snippet := content[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return snippet
}
