package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func TestExtractFindsNestedAndDeduplicatesModels(t *testing.T) {
	data := []byte(`{
        "nodes": [{"properties": {"models": [
            {"name":"base.safetensors","url":"https://example.test/base","directory":"diffusion_models"},
            {"name":"clip.safetensors","url":"https://example.test/clip","directory":"text_encoders"}
        ]}}],
        "definitions": {"subgraphs": [{"nodes": [{"properties": {"models": [
            {"name":"base.safetensors","url":"https://example.test/base","directory":"diffusion_models"}
        ]}}]}]}
    }`)

	got, err := Extract(data)
	if err != nil {
		t.Fatal(err)
	}
	want := []Model{
		{Name: "base.safetensors", URL: "https://example.test/base", Directory: "diffusion_models"},
		{Name: "clip.safetensors", URL: "https://example.test/clip", Directory: "text_encoders"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Extract() = %#v, want %#v", got, want)
	}
}

func TestDownloaderDownloadsSequentiallyAndSkipsExistingFiles(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mu.Lock()
		requests = append(requests, request.URL.Path)
		mu.Unlock()
		_, _ = writer.Write([]byte("content-" + filepath.Base(request.URL.Path)))
	}))
	defer server.Close()

	root := t.TempDir()
	references := []Model{
		{Name: "first.safetensors", URL: server.URL + "/first", Directory: "checkpoints"},
		{Name: "second.safetensors", URL: server.URL + "/second", Directory: "loras"},
	}
	downloader := Downloader{ModelsDirectory: root, HTTPClient: server.Client()}
	results, err := downloader.Download(context.Background(), references)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].Skipped || results[1].Skipped {
		t.Fatalf("results = %#v", results)
	}
	if got, want := requests, []string{"/first", "/second"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("request order = %v, want %v", got, want)
	}
	for _, reference := range references {
		path := filepath.Join(root, reference.Directory, reference.Name)
		if data, err := os.ReadFile(path); err != nil || len(data) == 0 {
			t.Fatalf("read %q: %v", path, err)
		}
	}

	results, err = downloader.Download(context.Background(), references)
	if err != nil {
		t.Fatal(err)
	}
	if !results[0].Skipped || !results[1].Skipped {
		t.Fatalf("second run results = %#v", results)
	}
	if len(requests) != 2 {
		t.Fatalf("second run made network requests: %v", requests)
	}
}

func TestDownloaderRejectsDestinationEscape(t *testing.T) {
	_, err := (Downloader{ModelsDirectory: t.TempDir(), DryRun: true}).Download(context.Background(), []Model{{
		Name: "..\\outside.safetensors", URL: "https://example.test/model", Directory: "loras",
	}})
	if err == nil {
		t.Fatal("Download() accepted a destination outside the models directory")
	}
}
