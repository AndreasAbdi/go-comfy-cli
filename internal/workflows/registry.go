package workflows

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Workflow is a named workflow file and its JSON definition.
// Definition is nil for entries returned by List and populated by Get.
type Workflow struct {
	Name       string
	Path       string
	Definition json.RawMessage
}

// Registry provides named access to the workflow files in one directory.
type Registry struct {
	dir string
}

func NewRegistry(dir string) *Registry {
	return &Registry{dir: dir}
}

// DefaultDir returns the ComfyUI Desktop workflow directory for the current user.
func DefaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Documents", "ComfyUI", "user", "default", "workflows"), nil
}

// List returns the workflow files in the registry, sorted by name.
func (r *Registry) List() ([]Workflow, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}

	result := make([]Workflow, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		result = append(result, Workflow{
			Name: entry.Name(),
			Path: filepath.Join(r.dir, entry.Name()),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})
	return result, nil
}

// GetWorkflow retrieves and validates a named workflow definition.
// The .json suffix is optional, but path-like names are rejected intentionally;
// --file is reserved for a later command extension.
func (r *Registry) GetWorkflow(name string) (Workflow, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Workflow{}, errors.New("workflow name is required")
	}
	if filepath.Base(name) != name || strings.ContainsAny(name, `/\\`) {
		return Workflow{}, fmt.Errorf("workflow name %q is not a name", name)
	}
	if !strings.EqualFold(filepath.Ext(name), ".json") {
		name += ".json"
	}

	entries, err := r.List()
	if err != nil {
		return Workflow{}, err
	}
	for _, entry := range entries {
		if !strings.EqualFold(entry.Name, name) {
			continue
		}

		definition, err := os.ReadFile(entry.Path)
		if err != nil {
			return Workflow{}, err
		}
		definition = bytes.TrimSpace(definition)
		if len(definition) == 0 {
			return Workflow{}, errors.New("workflow definition is empty")
		}
		var object map[string]any
		if err := json.Unmarshal(definition, &object); err != nil {
			return Workflow{}, fmt.Errorf("parse workflow JSON: %w", err)
		}
		if object == nil {
			return Workflow{}, errors.New("workflow definition must be a JSON object")
		}

		return Workflow{
			Name:       entry.Name,
			Path:       entry.Path,
			Definition: json.RawMessage(definition),
		}, nil
	}

	return Workflow{}, fmt.Errorf("workflow %q not found in %q", name, r.dir)
}

// Get is a short alias for GetWorkflow.
func (r *Registry) Get(name string) (Workflow, error) {
	return r.GetWorkflow(name)
}
