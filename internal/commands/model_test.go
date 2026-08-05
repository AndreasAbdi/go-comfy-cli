package commands

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestModelDownloadCommandDryRun(t *testing.T) {
	workflowPath := filepath.Join(t.TempDir(), "workflow.json")
	workflow := `{"definitions":{"subgraphs":[{"nodes":[{"properties":{"models":[{"name":"model.safetensors","url":"https://example.test/model","directory":"diffusion_models"}]}}]}]}}`
	if err := os.WriteFile(workflowPath, []byte(workflow), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsDirectory := filepath.Join(t.TempDir(), "models")
	var output strings.Builder
	app := NewApp()
	app.Writer = &output
	app.ErrWriter = io.Discard
	if err := app.Run([]string{
		"go-comfy-cli", "model", "download",
		"--workflow", workflowPath,
		"--models-directory", modelsDirectory,
		"--dry-run",
	}); err != nil {
		t.Fatal(err)
	}

	want := "would download model.safetensors -> " + filepath.Join(modelsDirectory, "diffusion_models", "model.safetensors")
	if !strings.Contains(output.String(), want) {
		t.Fatalf("output = %q, want it to contain %q", output.String(), want)
	}
	if _, err := os.Stat(modelsDirectory); !os.IsNotExist(err) {
		t.Fatalf("dry run created models directory: %v", err)
	}
}
