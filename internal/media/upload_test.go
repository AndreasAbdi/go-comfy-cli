package media

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-comfy-cli/internal/workflows"
)

func TestPrepareOperationsInlinesTextAndUploadsLocalFiles(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	workingDir := filepath.Join(root, "working")
	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(workingDir, 0o755); err != nil {
		t.Fatal(err)
	}

	videoPath := filepath.Join(workingDir, "reference.mp4")
	if err := os.WriteFile(videoPath, []byte("video-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	promptPath := filepath.Join(workingDir, "prompt.md")
	if err := os.WriteFile(promptPath, []byte("a long prompt\nwith two lines"), 0o600); err != nil {
		t.Fatal(err)
	}

	uploads := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/upload/image" {
			http.NotFound(writer, request)
			return
		}
		uploads++
		if err := request.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm() = %v", err)
		}
		file, header, err := request.FormFile("image")
		if err != nil {
			t.Fatalf("FormFile() = %v", err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("ReadAll() = %v", err)
		}
		if header.Filename != "reference.mp4" || string(data) != "video-bytes" {
			t.Fatalf("uploaded file = %q, %q", header.Filename, data)
		}
		if request.FormValue("type") != "input" || request.FormValue("subfolder") != defaultUploadSubfolder || request.FormValue("overwrite") != "true" {
			t.Fatalf("upload fields = type %q, subfolder %q, overwrite %q", request.FormValue("type"), request.FormValue("subfolder"), request.FormValue("overwrite"))
		}
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]string{
			"name":      header.Filename,
			"subfolder": defaultUploadSubfolder,
			"type":      "input",
		})
	}))
	defer server.Close()

	resolver := NewLocalPathResolver(inputDir, workingDir)
	uploader := NewUploaderWithResolver(server.URL, server.Client(), resolver)
	operations := []workflows.Operation{
		{Selector: ".video", Value: "./reference.mp4"},
		{Selector: ".prompt", Value: "./prompt.md"},
	}
	prepared, err := uploader.PrepareOperations(context.Background(), operations)
	if err != nil {
		t.Fatal(err)
	}
	if prepared[0].Value != "go-comfy-cli/reference.mp4" {
		t.Fatalf("video replacement = %#v", prepared[0].Value)
	}
	if prepared[1].Value != "a long prompt\nwith two lines" {
		t.Fatalf("prompt replacement = %#v", prepared[1].Value)
	}
	if uploads != 1 {
		t.Fatalf("upload count = %d, want 1", uploads)
	}
}

func TestPrepareValueUsesExistingComfyUIInputReference(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(inputDir, "existing.mp4")
	if err := os.WriteFile(path, []byte("already uploaded"), 0o600); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		t.Fatalf("existing input file was uploaded through %s", request.URL.Path)
	}))
	defer server.Close()

	uploader := NewUploaderWithResolver(server.URL, server.Client(), NewLocalPathResolver(inputDir))
	value, err := uploader.PrepareValue(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if value != "existing.mp4" {
		t.Fatalf("PrepareValue() = %#v", value)
	}
}

func TestUploadContentTypeIsDerivedFromExtension(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "reference.mp4")
	if err := os.WriteFile(path, []byte("video"), 0o600); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mediaType := request.Header.Get("Content-Type")
		if !strings.HasPrefix(mediaType, "multipart/form-data;") {
			t.Fatalf("Content-Type = %q", mediaType)
		}
		reader, err := request.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		for {
			part, nextErr := reader.NextPart()
			if nextErr == io.EOF {
				break
			}
			if nextErr != nil {
				t.Fatal(nextErr)
			}
			if part.FormName() == "image" && part.Header.Get("Content-Type") != "video/mp4" {
				t.Fatalf("file Content-Type = %q", part.Header.Get("Content-Type"))
			}
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"name":"reference.mp4","subfolder":"go-comfy-cli"}`))
	}))
	defer server.Close()

	uploader := NewUploaderWithResolver(server.URL, server.Client(), NewLocalPathResolver(filepath.Join(root, "input"), root))
	if _, err := uploader.PrepareValue(context.Background(), path); err != nil {
		t.Fatal(err)
	}
}
