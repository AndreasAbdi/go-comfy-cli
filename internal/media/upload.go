package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"go-comfy-cli/internal/workflows"
)

const defaultUploadSubfolder = "go-comfy-cli"

// Uploader prepares local file values for a ComfyUI workflow.
type Uploader struct {
	baseURL       string
	http          *http.Client
	resolver      *LocalPathResolver
	uploadFolder  string
	uploadedFiles map[string]string
}

type uploadResponse struct {
	Name      string `json:"name"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// NewUploader creates a media uploader using the standard local path roots.
func NewUploader(baseURL string, httpClient *http.Client) (*Uploader, error) {
	resolver, err := NewDefaultLocalPathResolver()
	if err != nil {
		return nil, fmt.Errorf("create local path resolver: %w", err)
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Uploader{
		baseURL:       strings.TrimRight(baseURL, "/"),
		http:          httpClient,
		resolver:      resolver,
		uploadFolder:  defaultUploadSubfolder,
		uploadedFiles: make(map[string]string),
	}, nil
}

// NewUploaderWithResolver is useful for tests and callers with a nonstandard
// ComfyUI input directory.
func NewUploaderWithResolver(baseURL string, httpClient *http.Client, resolver *LocalPathResolver) *Uploader {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Uploader{
		baseURL:       strings.TrimRight(baseURL, "/"),
		http:          httpClient,
		resolver:      resolver,
		uploadFolder:  defaultUploadSubfolder,
		uploadedFiles: make(map[string]string),
	}
}

// PrepareOperations resolves local file values before a workflow is
// transpiled. Text and Markdown files are read as literal string content;
// other local files are uploaded and replaced with a ComfyUI input reference.
func (u *Uploader) PrepareOperations(ctx context.Context, operations []workflows.Operation) ([]workflows.Operation, error) {
	prepared := make([]workflows.Operation, len(operations))
	copy(prepared, operations)
	for index := range prepared {
		value, err := u.PrepareValue(ctx, prepared[index].Value)
		if err != nil {
			return nil, fmt.Errorf("prepare operation %d value: %w", index+1, err)
		}
		prepared[index].Value = value
	}
	return prepared, nil
}

// PrepareValue resolves local files recursively in JSON-compatible values.
func (u *Uploader) PrepareValue(ctx context.Context, value any) (any, error) {
	switch typed := value.(type) {
	case string:
		return u.prepareString(ctx, typed)
	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			prepared, err := u.PrepareValue(ctx, item)
			if err != nil {
				return nil, err
			}
			result[index] = prepared
		}
		return result, nil
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			prepared, err := u.PrepareValue(ctx, item)
			if err != nil {
				return nil, err
			}
			result[key] = prepared
		}
		return result, nil
	default:
		return value, nil
	}
}

func (u *Uploader) prepareString(ctx context.Context, value string) (string, error) {
	resolved, found, err := u.resolver.Resolve(value)
	if err != nil {
		return "", fmt.Errorf("resolve local file %q: %w", value, err)
	}
	if !found {
		return value, nil
	}

	if isTextFile(resolved.Path) {
		contents, err := os.ReadFile(resolved.Path)
		if err != nil {
			return "", fmt.Errorf("read text file %q: %w", resolved.Path, err)
		}
		return string(contents), nil
	}
	if resolved.InputReference != "" {
		return resolved.InputReference, nil
	}

	key := filepath.Clean(resolved.Path)
	if uploaded, ok := u.uploadedFiles[key]; ok {
		return uploaded, nil
	}

	reference, err := u.upload(ctx, resolved.Path)
	if err != nil {
		return "", err
	}
	u.uploadedFiles[key] = reference
	return reference, nil
}

func (u *Uploader) upload(ctx context.Context, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file %q for upload: %w", path, err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writeUploadFields(writer, u.uploadFolder); err != nil {
		return "", err
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, escapeFilename(filepath.Base(path))))
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return "", fmt.Errorf("create upload part: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("read file %q for upload: %w", path, err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("finish upload request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.baseURL+"/upload/image", &body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := u.http.Do(request)
	if err != nil {
		return "", fmt.Errorf("upload %q: %w", path, err)
	}
	defer response.Body.Close()

	var result uploadResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode upload response (%s): %w", response.Status, err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("upload %q returned %s", path, response.Status)
	}
	if strings.TrimSpace(result.Name) == "" {
		return "", errors.New("upload response did not include a file name")
	}

	if strings.TrimSpace(result.Subfolder) == "" {
		return filepath.ToSlash(result.Name), nil
	}
	return filepath.ToSlash(filepath.Join(result.Subfolder, result.Name)), nil
}

func writeUploadFields(writer *multipart.Writer, subfolder string) error {
	fields := map[string]string{
		"type":      "input",
		"subfolder": subfolder,
		"overwrite": "true",
	}
	for _, key := range []string{"type", "subfolder", "overwrite"} {
		if err := writer.WriteField(key, fields[key]); err != nil {
			return fmt.Errorf("write upload field %q: %w", key, err)
		}
	}
	return nil
}

func escapeFilename(filename string) string {
	return strings.NewReplacer(`\`, `_`, `"`, `_`, "\r", `_`, "\n", `_`).Replace(filename)
}

func isTextFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown", ".txt":
		return true
	default:
		return false
	}
}
