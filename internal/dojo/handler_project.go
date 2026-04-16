package dojo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ProjectFile represents the project.json structure for a tracked project.
type ProjectFile struct {
	Name      string    `json:"name"`
	Phase     string    `json:"phase"`
	Tracks    []Track   `json:"tracks"`
	Decisions []string  `json:"decisions"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Track represents a single work track within a project.
type Track struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// projectsDir returns the ~/.dojo/projects/ path.
func projectsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".dojo", "projects")
	}
	return filepath.Join(home, ".dojo", "projects")
}

// findActiveProject returns the first project directory found under ~/.dojo/projects/.
// Returns "", nil if no project exists.
func findActiveProject() (string, error) {
	dir := projectsDir()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read projects dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			pf := filepath.Join(dir, e.Name(), "project.json")
			if _, err := os.Stat(pf); err == nil {
				return filepath.Join(dir, e.Name()), nil
			}
		}
	}
	return "", nil
}

// loadProject reads project.json from the given project directory.
func loadProject(projectDir string) (*ProjectFile, error) {
	data, err := os.ReadFile(filepath.Join(projectDir, "project.json"))
	if err != nil {
		return nil, fmt.Errorf("read project.json: %w", err)
	}
	var pf ProjectFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parse project.json: %w", err)
	}
	return &pf, nil
}

// saveProject writes project.json to the given project directory.
func saveProject(projectDir string, pf *ProjectFile) error {
	pf.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal project.json: %w", err)
	}
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}
	return os.WriteFile(filepath.Join(projectDir, "project.json"), data, 0644)
}

func (h *Handler) handleProjectStatus(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectDir, err := findActiveProject()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to scan projects: %v", err)), nil
	}
	if projectDir == "" {
		return mcp.NewToolResultText(
			"No project found in ~/.dojo/projects/\n\n" +
				"Use `dojo_project_track` with action \"add\" to start tracking a project,\n" +
				"or create ~/.dojo/projects/{name}/project.json manually.",
		), nil
	}

	pf, err := loadProject(projectDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load project: %v", err)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# Project: %s\n\n", pf.Name)
	if pf.Phase != "" {
		fmt.Fprintf(&sb, "**Phase:** %s\n\n", pf.Phase)
	}

	if len(pf.Tracks) > 0 {
		fmt.Fprintf(&sb, "## Tracks\n\n")
		fmt.Fprintf(&sb, "| Track | Status |\n")
		fmt.Fprintf(&sb, "|-------|--------|\n")
		for _, t := range pf.Tracks {
			fmt.Fprintf(&sb, "| %s | %s |\n", t.Name, t.Status)
		}
		fmt.Fprintf(&sb, "\n")
	} else {
		fmt.Fprintf(&sb, "No tracks defined yet.\n\n")
	}

	if len(pf.Decisions) > 0 {
		fmt.Fprintf(&sb, "## Decisions (%d)\n\n", len(pf.Decisions))
		for i, d := range pf.Decisions {
			fmt.Fprintf(&sb, "%d. %s\n", i+1, d)
		}
		fmt.Fprintf(&sb, "\n")
	}

	fmt.Fprintf(&sb, "_Project file: %s_\n", filepath.Join(projectDir, "project.json"))

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handler) handleProjectTrack(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Action string `json:"action"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Action == "" {
		return mcp.NewToolResultError("'action' is required (add or set)"), nil
	}
	if args.Name == "" {
		return mcp.NewToolResultError("'name' is required"), nil
	}

	action := strings.ToLower(args.Action)
	if action != "add" && action != "set" {
		return mcp.NewToolResultError("'action' must be 'add' or 'set'"), nil
	}

	projectDir, err := findActiveProject()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to scan projects: %v", err)), nil
	}

	// Create a default project if none exists
	if projectDir == "" {
		projectDir = filepath.Join(projectsDir(), "default")
		pf := &ProjectFile{
			Name:   "default",
			Phase:  "active",
			Tracks: []Track{},
		}
		if err := saveProject(projectDir, pf); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create default project: %v", err)), nil
		}
	}

	pf, err := loadProject(projectDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load project: %v", err)), nil
	}

	status := args.Status
	if status == "" {
		status = "in-progress"
	}

	switch action {
	case "add":
		// Check for duplicate
		for _, t := range pf.Tracks {
			if strings.EqualFold(t.Name, args.Name) {
				return mcp.NewToolResultError(fmt.Sprintf("Track '%s' already exists. Use action 'set' to update its status.", args.Name)), nil
			}
		}
		pf.Tracks = append(pf.Tracks, Track{Name: args.Name, Status: status})

	case "set":
		found := false
		for i, t := range pf.Tracks {
			if strings.EqualFold(t.Name, args.Name) {
				pf.Tracks[i].Status = status
				found = true
				break
			}
		}
		if !found {
			return mcp.NewToolResultError(fmt.Sprintf("Track '%s' not found. Use action 'add' to create it.", args.Name)), nil
		}
	}

	if err := saveProject(projectDir, pf); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to save project: %v", err)), nil
	}

	verb := "Added"
	if action == "set" {
		verb = "Updated"
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s track **%s** → status: **%s**\n\nProject: %s", verb, args.Name, status, pf.Name)), nil
}

func (h *Handler) handleProjectDecision(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Text string `json:"text"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Text == "" {
		return mcp.NewToolResultError("'text' is required"), nil
	}

	projectDir, err := findActiveProject()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to scan projects: %v", err)), nil
	}

	// Create a default project if none exists
	if projectDir == "" {
		projectDir = filepath.Join(projectsDir(), "default")
		pf := &ProjectFile{
			Name:  "default",
			Phase: "active",
		}
		if err := saveProject(projectDir, pf); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create default project: %v", err)), nil
		}
	}

	pf, err := loadProject(projectDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load project: %v", err)), nil
	}

	pf.Decisions = append(pf.Decisions, args.Text)

	if err := saveProject(projectDir, pf); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to save project: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Decision recorded (#%d):\n\n> %s\n\nProject: %s (%d total decisions)",
		len(pf.Decisions), args.Text, pf.Name, len(pf.Decisions),
	)), nil
}
