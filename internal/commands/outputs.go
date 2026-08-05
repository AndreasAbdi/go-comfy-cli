package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type comfyOutputReference struct {
	Filename  string
	Subfolder string
	Type      string
}

// DownloadOutputs copies the files reported by ComfyUI history into folder
// and returns their local paths. ComfyUI subfolders are preserved below it.
func (c comfyClient) DownloadOutputs(ctx context.Context, outputs any, folder string) ([]string, error) {
	folder, err := filepath.Abs(folder)
	if err != nil {
		return nil, fmt.Errorf("resolve output folder: %w", err)
	}
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return nil, fmt.Errorf("create output folder %q: %w", folder, err)
	}

	references := make(map[string]comfyOutputReference)
	collectOutputReferences(outputs, references)
	ordered := make([]comfyOutputReference, 0, len(references))
	for _, reference := range references {
		ordered = append(ordered, reference)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return outputReferenceKey(ordered[i]) < outputReferenceKey(ordered[j])
	})

	paths := make([]string, 0, len(ordered))
	for _, reference := range ordered {
		path, err := outputPath(folder, reference)
		if err != nil {
			return nil, err
		}
		if err := c.downloadOutput(ctx, reference, path); err != nil {
			return nil, fmt.Errorf("download output %q: %w", reference.Filename, err)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func collectOutputReferences(value any, references map[string]comfyOutputReference) {
	switch typed := value.(type) {
	case map[string]any:
		filename, _ := typed["filename"].(string)
		filename = strings.TrimSpace(filename)
		if filename != "" {
			subfolder, _ := typed["subfolder"].(string)
			fileType, _ := typed["type"].(string)
			fileType = strings.TrimSpace(fileType)
			if fileType == "" {
				fileType = "output"
			}
			reference := comfyOutputReference{
				Filename:  filename,
				Subfolder: strings.TrimSpace(subfolder),
				Type:      fileType,
			}
			references[outputReferenceKey(reference)] = reference
		}
		for _, item := range typed {
			collectOutputReferences(item, references)
		}
	case []any:
		for _, item := range typed {
			collectOutputReferences(item, references)
		}
	}
}

func outputReferenceKey(reference comfyOutputReference) string {
	return strings.Join([]string{reference.Type, reference.Subfolder, reference.Filename}, "\x00")
}

func outputPath(folder string, reference comfyOutputReference) (string, error) {
	filename := filepath.FromSlash(reference.Filename)
	if filename == "" || filename == "." || filename == ".." || filepath.IsAbs(filename) {
		return "", fmt.Errorf("invalid ComfyUI output filename %q", reference.Filename)
	}
	subfolder := filepath.FromSlash(reference.Subfolder)
	relative := filepath.Clean(filepath.Join(subfolder, filename))
	if relative == "." || relative == ".." || filepath.IsAbs(relative) || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid ComfyUI output path %q/%q", reference.Subfolder, reference.Filename)
	}

	path := filepath.Join(folder, relative)
	check, err := filepath.Rel(folder, path)
	if err != nil || check == ".." || strings.HasPrefix(check, ".."+string(filepath.Separator)) || filepath.IsAbs(check) {
		return "", fmt.Errorf("ComfyUI output escapes output folder: %q/%q", reference.Subfolder, reference.Filename)
	}
	return path, nil
}

func (c comfyClient) downloadOutput(ctx context.Context, reference comfyOutputReference, destination string) error {
	query := url.Values{}
	query.Set("filename", reference.Filename)
	if reference.Subfolder != "" {
		query.Set("subfolder", reference.Subfolder)
	}
	query.Set("type", reference.Type)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/view?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("server returned %s", response.Status)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create output subfolder: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".go-comfy-cli-output-*")
	if err != nil {
		return fmt.Errorf("create temporary output file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := io.Copy(temporary, response.Body); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("save output: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary output file: %w", err)
	}
	if err := os.Remove(destination); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("replace output file: %w", err)
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return fmt.Errorf("install output file: %w", err)
	}
	return nil
}
