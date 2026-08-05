// Package models extracts and installs model metadata embedded in ComfyUI
// workflow definitions.
package models

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Model is a downloadable model reference from a workflow's models section.
type Model struct {
	Name      string
	URL       string
	Directory string
	SHA256    string
	Size      int64
}

// Result describes one model installation attempt.
type Result struct {
	Model   Model
	Path    string
	Skipped bool
}

// Extract finds model references anywhere in a workflow JSON document. The
// current ComfyUI workflow format stores them under node properties, including
// nodes nested in definitions.subgraphs, so this intentionally walks every
// object and array rather than depending on one particular workflow shape.
func Extract(data []byte) ([]Model, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var document any
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("parse workflow JSON: %w", err)
	}

	result := make([]Model, 0)
	byDestination := make(map[string]int)
	if err := walk(document, &result, byDestination); err != nil {
		return nil, err
	}
	return result, nil
}

func walk(value any, result *[]Model, byDestination map[string]int) error {
	switch value := value.(type) {
	case []any:
		for _, child := range value {
			if err := walk(child, result, byDestination); err != nil {
				return err
			}
		}
	case map[string]any:
		if rawModels, ok := value["models"]; ok {
			entries, ok := rawModels.([]any)
			if !ok {
				return errors.New(`workflow "models" sections must be arrays`)
			}
			for index, rawModel := range entries {
				model, err := parseModel(rawModel)
				if err != nil {
					return fmt.Errorf("models[%d]: %w", index, err)
				}

				key := strings.ToLower(model.Directory) + "\x00" + strings.ToLower(model.Name)
				if previous, exists := byDestination[key]; exists {
					if (*result)[previous].URL != model.URL {
						return fmt.Errorf("model %q in directory %q has conflicting URLs", model.Name, model.Directory)
					}
					continue
				}
				byDestination[key] = len(*result)
				*result = append(*result, model)
			}
		}

		keys := make([]string, 0, len(value))
		for key := range value {
			if key != "models" {
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)
		for _, key := range keys {
			child := value[key]
			if err := walk(child, result, byDestination); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseModel(value any) (Model, error) {
	object, ok := value.(map[string]any)
	if !ok {
		return Model{}, errors.New("model entry must be an object")
	}

	name := stringField(object, "relative_path", "path", "name", "filename")
	if name == "" {
		return Model{}, errors.New("model entry has no name")
	}

	modelURL := stringField(object, "url", "download_url")
	if modelURL == "" {
		return Model{}, fmt.Errorf("model %q has no URL", name)
	}

	directory := stringField(object, "directory", "category", "save_path")
	if directory == "" {
		return Model{}, fmt.Errorf("model %q has no destination directory", name)
	}

	model := Model{
		Name:      name,
		URL:       modelURL,
		Directory: directory,
		SHA256:    strings.TrimSpace(stringField(object, "sha256", "hash")),
	}
	if rawSize, ok := object["size"]; ok {
		size, err := modelSize(rawSize)
		if err != nil {
			return Model{}, fmt.Errorf("model %q has invalid size: %w", name, err)
		}
		model.Size = size
	}
	return model, nil
}

func stringField(object map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := object[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func modelSize(value any) (int64, error) {
	var encoded string
	switch value := value.(type) {
	case json.Number:
		encoded = string(value)
	case string:
		encoded = strings.TrimSpace(value)
	default:
		return 0, errors.New("expected an integer")
	}

	size, err := strconv.ParseInt(encoded, 10, 64)
	if err != nil || size < 0 {
		return 0, errors.New("expected a non-negative integer")
	}
	return size, nil
}

// DefaultDirectory returns the standard ComfyUI Desktop model root.
func DefaultDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Documents", "ComfyUI", "models"), nil
}

// Downloader installs models sequentially into a ComfyUI model root.
type Downloader struct {
	ModelsDirectory string
	HTTPClient      *http.Client
	DryRun          bool
}

// Download installs each model in order. Existing regular files are left in
// place, making repeated provisioning runs safe and resumable.
func (d Downloader) Download(ctx context.Context, references []Model) ([]Result, error) {
	root := strings.TrimSpace(d.ModelsDirectory)
	if root == "" {
		return nil, errors.New("models directory is required")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve models directory: %w", err)
	}

	client := d.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	results := make([]Result, 0, len(references))
	for index, reference := range references {
		path, err := modelPath(root, reference)
		if err != nil {
			return nil, fmt.Errorf("model %d (%q): %w", index+1, reference.Name, err)
		}

		info, statErr := os.Stat(path)
		switch {
		case statErr == nil:
			if !info.Mode().IsRegular() {
				return nil, fmt.Errorf("model %q destination %q is not a regular file", reference.Name, path)
			}
			results = append(results, Result{Model: reference, Path: path, Skipped: true})
			continue
		case !errors.Is(statErr, os.ErrNotExist):
			return nil, fmt.Errorf("inspect model %q destination: %w", reference.Name, statErr)
		}

		if d.DryRun {
			results = append(results, Result{Model: reference, Path: path})
			continue
		}

		if err := d.downloadOne(ctx, client, reference, path); err != nil {
			return nil, fmt.Errorf("download model %q: %w", reference.Name, err)
		}
		results = append(results, Result{Model: reference, Path: path})
	}
	return results, nil
}

func modelPath(root string, reference Model) (string, error) {
	directory := normalizeModelPath(reference.Directory)
	name := normalizeModelPath(reference.Name)
	if directory == "" || name == "" {
		return "", errors.New("model directory and name are required")
	}
	if isAbsoluteModelPath(directory) || isAbsoluteModelPath(name) {
		return "", errors.New("model destination must be relative")
	}

	categoryRoot := filepath.Join(root, directory)
	if !pathWithinModelRoot(root, categoryRoot) {
		return "", errors.New("model destination escapes the models directory")
	}
	destination := filepath.Join(categoryRoot, name)
	if !pathWithinModelRoot(categoryRoot, destination) {
		return "", errors.New("model destination escapes the models directory")
	}
	return destination, nil
}

func normalizeModelPath(value string) string {
	return filepath.FromSlash(strings.ReplaceAll(strings.TrimSpace(value), `\`, "/"))
}

func isAbsoluteModelPath(value string) bool {
	normalized := strings.ReplaceAll(value, `\`, "/")
	if filepath.IsAbs(value) || filepath.VolumeName(value) != "" || strings.HasPrefix(normalized, "/") {
		return true
	}
	return len(normalized) >= 2 && isASCIIAlphaModelPath(normalized[0]) && normalized[1] == ':'
}

func isASCIIAlphaModelPath(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}

func pathWithinModelRoot(root, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}

func (d Downloader) downloadOne(ctx context.Context, client *http.Client, reference Model, destination string) error {
	parsed, err := url.Parse(reference.URL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("invalid download URL %q", reference.URL)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(destination), ".go-comfy-cli-model-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, reference.URL, nil)
	if err != nil {
		_ = temporary.Close()
		return fmt.Errorf("create request: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		_ = temporary.Close()
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		_ = response.Body.Close()
		_ = temporary.Close()
		return fmt.Errorf("server returned %s", response.Status)
	}

	hasher := sha256.New()
	bytesWritten, copyErr := io.Copy(io.MultiWriter(temporary, hasher), response.Body)
	closeResponseErr := response.Body.Close()
	closeTemporaryErr := temporary.Close()
	if copyErr != nil {
		return fmt.Errorf("write download: %w", copyErr)
	}
	if closeResponseErr != nil {
		return fmt.Errorf("close response: %w", closeResponseErr)
	}
	if closeTemporaryErr != nil {
		return fmt.Errorf("close temporary file: %w", closeTemporaryErr)
	}
	if reference.Size > 0 && bytesWritten != reference.Size {
		return fmt.Errorf("downloaded %d bytes, expected %d", bytesWritten, reference.Size)
	}
	if expected := strings.ToLower(strings.TrimSpace(reference.SHA256)); expected != "" {
		actual := hex.EncodeToString(hasher.Sum(nil))
		if actual != expected {
			return fmt.Errorf("sha256 %s does not match expected %s", actual, expected)
		}
	}

	if err := os.Rename(temporaryPath, destination); err != nil {
		return fmt.Errorf("install downloaded file: %w", err)
	}
	return nil
}
