package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunUploadsMediaFromCurrentWorkingDirectory(t *testing.T) {
	workflowDir := t.TempDir()
	workingDir := t.TempDir()
	workflowPath := filepath.Join(workflowDir, "workflow.json")
	workflow := `{"nodes":[{"id":1,"type":"LoadImage","inputs":[{"name":"image","type":"COMBO","widget":{"name":"image"},"link":null}],"outputs":[],"widgets_values":["reference.png","image"]}],"links":[]}`
	if err := os.WriteFile(workflowPath, []byte(workflow), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "reference.png"), []byte("packaged-image"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workingDir, "reference.png"), []byte("working-directory-image"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(workingDir)

	uploads := 0
	var prompt map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/upload/image":
			uploads++
			if err := request.ParseMultipartForm(1 << 20); err != nil {
				t.Fatal(err)
			}
			file, _, err := request.FormFile("image")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			contents, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}
			if string(contents) != "working-directory-image" {
				t.Fatalf("uploaded contents = %q", contents)
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"name":"reference.png","subfolder":"go-comfy-cli"}`))
		case "/prompt":
			var payload struct {
				Prompt map[string]any `json:"prompt"`
			}
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			prompt = payload.Prompt
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"prompt_id":"prompt-1","node_errors":{}}`))
		case "/history/prompt-1":
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"prompt-1":{"status":{"status_str":"success","completed":true}}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	app := NewApp()
	app.Writer = io.Discard
	app.ErrWriter = io.Discard
	if err := app.Run([]string{
		"go-comfy-cli", "run",
		"--workflow", workflowPath,
		"--url", server.URL,
	}); err != nil {
		t.Fatal(err)
	}
	if uploads != 1 {
		t.Fatalf("upload count = %d, want 1", uploads)
	}
	node := prompt["1"].(map[string]any)
	inputs := node["inputs"].(map[string]any)
	if inputs["image"] != "go-comfy-cli/reference.png" {
		t.Fatalf("prompt image = %#v", inputs["image"])
	}
}
