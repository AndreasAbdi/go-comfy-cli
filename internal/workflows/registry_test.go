package workflows

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRegistryListsSortedJSONFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"zeta.json", "alpha.JSON", "notes.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(`{"name":"test"}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "nested.json"), 0o700); err != nil {
		t.Fatal(err)
	}

	got, err := NewRegistry(dir).List()
	if err != nil {
		t.Fatal(err)
	}
	gotNames := make([]string, 0, len(got))
	for _, workflow := range got {
		gotNames = append(gotNames, workflow.Name)
	}

	want := []string{"alpha.JSON", "zeta.json"}
	if !reflect.DeepEqual(gotNames, want) {
		t.Fatalf("List() = %v, want %v", gotNames, want)
	}
}

func TestRegistryGetsWorkflowByNameWithoutExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.json")
	if err := os.WriteFile(path, []byte(`{"nodes":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}

	workflow, err := NewRegistry(dir).Get("example")
	if err != nil {
		t.Fatal(err)
	}
	if workflow.Name != "example.json" || workflow.Path != path {
		t.Fatalf("Get() returned %+v", workflow)
	}
	if string(workflow.Definition) != `{"nodes":[]}` {
		t.Fatalf("Get() returned definition %s", workflow.Definition)
	}
}

func TestRegistryRejectsPathLikeName(t *testing.T) {
	if _, err := NewRegistry(t.TempDir()).Get("..\\example"); err == nil {
		t.Fatal("Get() accepted a path-like workflow name")
	}
}
