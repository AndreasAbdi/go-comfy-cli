package commands

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDownloadOutputsWritesFilesAndPreservesSubfolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/view" {
			t.Fatalf("request path = %q, want /view", request.URL.Path)
		}
		if request.URL.Query().Get("filename") != "result.mp4" || request.URL.Query().Get("subfolder") != "video" || request.URL.Query().Get("type") != "output" {
			t.Fatalf("request query = %v", request.URL.Query())
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("video-bytes"))
	}))
	defer server.Close()

	folder := filepath.Join(t.TempDir(), "out")
	client := comfyClient{baseURL: server.URL, http: server.Client()}
	paths, err := client.DownloadOutputs(context.Background(), map[string]any{
		"node": map[string]any{
			"videos": []any{
				map[string]any{"filename": "result.mp4", "subfolder": "video", "type": "output"},
			},
		},
	}, folder)
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(folder, "video", "result.mp4")
	if !reflect.DeepEqual(paths, []string{wantPath}) {
		t.Fatalf("DownloadOutputs() = %v, want %v", paths, []string{wantPath})
	}
	contents, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "video-bytes" {
		t.Fatalf("output contents = %q", contents)
	}
}

func TestOutputPathRejectsTraversal(t *testing.T) {
	_, err := outputPath(t.TempDir(), comfyOutputReference{Filename: "..\\outside.mp4", Type: "output"})
	if err == nil {
		t.Fatal("outputPath() accepted traversal")
	}
}

func TestDownloadWorkflowCommandWritesDefinitionToStdout(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "example.json"), []byte(`{"nodes":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var output strings.Builder
	app := NewApp()
	app.Writer = &output
	app.ErrWriter = &output
	if err := app.Run([]string{"go-comfy-cli", "workflow", "download", "--dir", dir, "example"}); err != nil {
		t.Fatal(err)
	}
	if output.String() != `{"nodes":[]}` {
		t.Fatalf("download output = %q", output.String())
	}
}

func TestUploadWorkflowCommandReadsStdin(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "workflows")
	app := NewApp()
	app.Reader = strings.NewReader(`{"nodes":[]}`)
	app.Writer = io.Discard
	app.ErrWriter = io.Discard
	if err := app.Run([]string{"go-comfy-cli", "workflow", "upload", "--dir", dir, "from-stdin"}); err != nil {
		t.Fatal(err)
	}
	definition, err := os.ReadFile(filepath.Join(dir, "from-stdin.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(definition) != `{"nodes":[]}` {
		t.Fatalf("uploaded definition = %q", definition)
	}
}

func TestUploadWorkflowCommandReadsInputFile(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "workflows")
	inputFile := filepath.Join(root, "workflow.json")
	if err := os.WriteFile(inputFile, []byte(`{"nodes":[{"id":1}]}`), 0o600); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.Writer = io.Discard
	app.ErrWriter = io.Discard
	if err := app.Run([]string{"go-comfy-cli", "workflow", "upload", "--dir", dir, "--input-file", inputFile, "from-file"}); err != nil {
		t.Fatal(err)
	}
	definition, err := os.ReadFile(filepath.Join(dir, "from-file.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(definition) != `{"nodes":[{"id":1}]}` {
		t.Fatalf("uploaded definition = %q", definition)
	}
}
