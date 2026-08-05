package workflows

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListReturnsSortedJSONFiles(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"zeta.json", "alpha.JSON", "notes.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "nested.json"), 0o700); err != nil {
		t.Fatal(err)
	}

	got, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"alpha.JSON", "zeta.json"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}

func TestListReturnsDirectoryError(t *testing.T) {
	_, err := List(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("List() returned nil error for a missing directory")
	}
}
