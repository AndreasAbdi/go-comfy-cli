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
// The .json suffix is optional, but path-like names are rejected intentionally.
func (r *Registry) GetWorkflow(name string) (Workflow, error) {
	filename, err := workflowFilename(name)
	if err != nil {
		return Workflow{}, err
	}

	entries, err := r.List()
	if err != nil {
		return Workflow{}, err
	}
	for _, entry := range entries {
		if !strings.EqualFold(entry.Name, filename) {
			continue
		}

		return loadWorkflowFile(entry.Path, entry.Name)
	}

	return Workflow{}, fmt.Errorf("workflow %q not found in %q", filename, r.dir)
}

// LoadWorkflow reads and validates a workflow definition from an explicit
// path. Unlike GetWorkflow, the path does not need to be inside the user's
// ComfyUI workflow directory.
func LoadWorkflow(path string) (Workflow, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return Workflow{}, errors.New("workflow file path is required")
	}
	return loadWorkflowFile(path, filepath.Base(path))
}

func loadWorkflowFile(path, name string) (Workflow, error) {
	definition, err := os.ReadFile(path)
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
		Name:       name,
		Path:       path,
		Definition: json.RawMessage(definition),
	}, nil
}

// Get is a short alias for GetWorkflow.
func (r *Registry) Get(name string) (Workflow, error) {
	return r.GetWorkflow(name)
}

// Upload validates and writes a named workflow definition to the registry.
// The .json suffix is optional, matching GetWorkflow.
func (r *Registry) Upload(name string, definition []byte) (Workflow, error) {
	filename, err := workflowFilename(name)
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

	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return Workflow{}, fmt.Errorf("create workflow directory: %w", err)
	}
	path := filepath.Join(r.dir, filename)
	temporary, err := os.CreateTemp(r.dir, ".go-comfy-cli-workflow-*")
	if err != nil {
		return Workflow{}, fmt.Errorf("create temporary workflow file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporary.Write(definition); err != nil {
		_ = temporary.Close()
		return Workflow{}, fmt.Errorf("write workflow definition: %w", err)
	}
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return Workflow{}, fmt.Errorf("set workflow permissions: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return Workflow{}, fmt.Errorf("close temporary workflow file: %w", err)
	}

	// Windows does not replace an existing file with os.Rename, so remove the
	// destination only after the complete temporary file has been written.
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Workflow{}, fmt.Errorf("replace workflow %q: %w", filename, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return Workflow{}, fmt.Errorf("install workflow %q: %w", filename, err)
	}

	return Workflow{
		Name:       filename,
		Path:       path,
		Definition: json.RawMessage(definition),
	}, nil
}

func workflowFilename(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("workflow name is required")
	}
	if name == "." || name == ".." || filepath.Base(name) != name || strings.ContainsAny(name, `/\\`) {
		return "", fmt.Errorf("workflow name %q is not a name", name)
	}
	if !strings.EqualFold(filepath.Ext(name), ".json") {
		name += ".json"
	}
	return name, nil
}
