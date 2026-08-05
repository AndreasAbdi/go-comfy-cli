package workflows

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultDir returns the ComfyUI Desktop workflow directory for the current user.
func DefaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Documents", "ComfyUI", "user", "default", "workflows"), nil
}

// List returns the workflow filenames in dir, sorted alphabetically.
// ComfyUI Desktop stores workflows as JSON files in this directory.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		result = append(result, entry.Name())
	}

	sort.Strings(result)
	return result, nil
}
