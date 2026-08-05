package media

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalPathResolverFindsInputReference(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	if err := os.MkdirAll(filepath.Join(inputDir, "video"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(inputDir, "video", "reference.mp4")
	if err := os.WriteFile(path, []byte("video"), 0o600); err != nil {
		t.Fatal(err)
	}

	resolver := NewLocalPathResolver(inputDir, filepath.Join(root, "working"))
	resolved, found, err := resolver.Resolve("video/reference.mp4")
	if err != nil {
		t.Fatal(err)
	}
	if !found || resolved.Path != path || resolved.InputReference != "video/reference.mp4" {
		t.Fatalf("Resolve() = %#v, found %v", resolved, found)
	}
}

func TestLocalPathResolverLeavesMissingValuesUnchanged(t *testing.T) {
	resolver := NewLocalPathResolver(t.TempDir())
	_, found, err := resolver.Resolve("not-a-local-file.mp4")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("Resolve() found a missing path")
	}
}

func TestLocalPathResolverLeavesLongLiteralValuesUntouched(t *testing.T) {
	requested := strings.Repeat("a prompt with words, ", 20)
	resolver := NewLocalPathResolver(t.TempDir(), t.TempDir())

	resolved, found, err := resolver.Resolve(requested)
	if err != nil {
		t.Fatal(err)
	}
	if found || resolved != (LocalPath{}) {
		t.Fatalf("Resolve() = %#v, %v; want an unresolved literal", resolved, found)
	}
}
