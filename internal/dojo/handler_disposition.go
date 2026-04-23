package dojo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/DojoGenesis/mcp-server/internal/fsutil"
	"github.com/mark3labs/mcp-go/mcp"
)

// dojoSettingsMu serializes concurrent load-modify-write sequences on
// ~/.dojo/settings.json. AtomicWriteFile prevents torn writes, but this mutex
// prevents two concurrent goroutines from each reading the same stale state and
// then overwriting each other's changes.
var dojoSettingsMu sync.Mutex

// dojoSettingsPath returns the path to ~/.dojo/settings.json.
func dojoSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".dojo", "settings.json"), nil
}

// loadDojoSettings reads ~/.dojo/settings.json as a generic map so we preserve
// unknown keys on write. Returns an empty map if the file does not exist.
func loadDojoSettings(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]interface{}{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}
	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse settings: %w", err)
	}
	if m == nil {
		m = map[string]interface{}{}
	}
	return m, nil
}

// saveDojoSettings writes the settings map as indented JSON atomically.
// Callers must hold dojoSettingsMu before calling this function.
func saveDojoSettings(path string, m map[string]interface{}) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}
	return fsutil.AtomicWriteFile(path, data, 0o644)
}

// handleDispositionSet writes the selected disposition under the top-level
// "defaults" key in ~/.dojo/settings.json. All other keys are preserved.
func (h *Handler) handleDispositionSet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Disposition string `json:"disposition"`
	}
	if err := unmarshalArgs(request.Params.Arguments, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}
	if args.Disposition == "" {
		return mcp.NewToolResultError("'disposition' is required and cannot be empty"), nil
	}

	path, err := dojoSettingsPath()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve settings path: %v", err)), nil
	}

	dojoSettingsMu.Lock()
	defer dojoSettingsMu.Unlock()

	settings, err := loadDojoSettings(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load settings: %v", err)), nil
	}

	// Ensure "defaults" is a map
	var defaults map[string]interface{}
	if existing, ok := settings["defaults"]; ok {
		if asMap, ok := existing.(map[string]interface{}); ok {
			defaults = asMap
		}
	}
	if defaults == nil {
		defaults = map[string]interface{}{}
	}
	defaults["disposition"] = args.Disposition
	settings["defaults"] = defaults

	if err := saveDojoSettings(path, settings); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write settings: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Disposition set to **%s**\n\nWrote to: %s\n\nThe value lives under the top-level `defaults.disposition` key and will be read by the Dojo CLI / Gateway next time they load settings.",
		args.Disposition, path,
	)), nil
}
