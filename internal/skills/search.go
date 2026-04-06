package skills

import (
	"sort"
	"strings"
)

// searchResult pairs a skill with its relevance score.
type searchResult struct {
	Skill     Skill
	Relevance float64
}

// Search returns skills matching a query, ranked by relevance.
//
// Matching algorithm:
//  1. Exact name match -> relevance 1.0
//  2. Name contains query -> relevance 0.8
//  3. Trigger phrase match -> relevance 0.7
//  4. Description contains query -> relevance 0.5
//  5. Content contains query keyword -> relevance 0.2
//
// Returns at most maxResults skills, sorted by relevance descending.
func (l *Loader) Search(query string, maxResults int) []Skill {
	if query == "" {
		// Empty query returns all skills (up to maxResults)
		if maxResults <= 0 || maxResults > len(l.skills) {
			return l.skills
		}
		return l.skills[:maxResults]
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))
	keywords := strings.Fields(queryLower)

	var results []searchResult

	for _, skill := range l.skills {
		score := calculateSkillRelevance(queryLower, keywords, skill)
		if score > 0.1 {
			results = append(results, searchResult{
				Skill:     skill,
				Relevance: score,
			})
		}
	}

	// Sort by relevance descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})

	// Limit results
	if maxResults > 0 && len(results) > maxResults {
		results = results[:maxResults]
	}

	skills := make([]Skill, len(results))
	for i, r := range results {
		skills[i] = r.Skill
	}
	return skills
}

// calculateSkillRelevance scores a single skill against the query.
func calculateSkillRelevance(queryLower string, keywords []string, skill Skill) float64 {
	nameLower := strings.ToLower(skill.Name)
	descLower := strings.ToLower(skill.Description)
	contentLower := strings.ToLower(skill.Content)

	score := 0.0

	// 1. Exact name match
	if nameLower == queryLower {
		score += 1.0
	}

	// 2. Name contains query
	if strings.Contains(nameLower, queryLower) {
		score += 0.8
	}

	// 3. Trigger phrase match
	for _, trigger := range skill.Triggers {
		if strings.Contains(trigger, queryLower) || strings.Contains(queryLower, trigger) {
			score += 0.7
			break
		}
	}

	// 4. Description contains query
	if strings.Contains(descLower, queryLower) {
		score += 0.5
	}

	// 5. Content contains query keyword
	if strings.Contains(contentLower, queryLower) {
		score += 0.2
	}

	// Keyword matching (individual words)
	for _, keyword := range keywords {
		if len(keyword) < 3 {
			continue
		}
		if strings.Contains(nameLower, keyword) {
			score += 0.3
		}
		if strings.Contains(descLower, keyword) {
			score += 0.15
		}
		// Check triggers for keyword match
		for _, trigger := range skill.Triggers {
			if strings.Contains(trigger, keyword) {
				score += 0.1
				break
			}
		}
	}

	return score
}
