package workflows

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadMappingAndBuildAliasOperations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mapping.yaml")
	mappingYAML := `version: 1
aliases:
  prompt:
    selector: '.nodes[] | select(.type == "Prompt") | .widgets_values[0]'
    type: string
    cardinality: one
  duration:
    selector: '.nodes[] | select(.type == "Duration") | .widgets_values[0]'
    type: number
    cardinality: one
  image:
    selector: '.nodes[] | select(.type == "LoadImage") | .widgets_values[0]'
    indexed_selector: '.nodes[$index].widgets_values[0]'
    type: string
    cardinality: many
`
	if err := os.WriteFile(path, []byte(mappingYAML), 0o600); err != nil {
		t.Fatal(err)
	}

	mapping, err := LoadMapping(path)
	if err != nil {
		t.Fatal(err)
	}
	operations, err := mapping.AliasOperations([]string{"prompt=A new prompt", "duration=8", "image[1]=second.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(operations) != 3 {
		t.Fatalf("AliasOperations() returned %d operations", len(operations))
	}
	if operations[0].Selector != mapping.Aliases["prompt"].Selector || operations[0].Value != "A new prompt" {
		t.Fatalf("prompt operation = %#v", operations[0])
	}
	number, numberOK := operations[1].Value.(json.Number)
	if !numberOK || number.String() != "8" || operations[1].Cardinality != "one" {
		t.Fatalf("duration operation = %#v", operations[1])
	}
	if operations[2].Selector != mapping.Aliases["image"].IndexedSelector || operations[2].Variables["index"] != 1 {
		t.Fatalf("image operation = %#v", operations[2])
	}
}

func TestAliasOperationsApplyIndexedSelector(t *testing.T) {
	mapping := Mapping{
		Version: 1,
		Aliases: map[string]Alias{
			"image": {
				Selector:        ".nodes[] | select(.type == \"LoadImage\") | .widgets_values[0]",
				IndexedSelector: ".nodes[$index].widgets_values[0]",
				Type:            "string",
				Cardinality:     "many",
			},
		},
	}
	operations, err := mapping.AliasOperations([]string{"image[1]=second.png"})
	if err != nil {
		t.Fatal(err)
	}
	document := map[string]any{"nodes": []any{
		map[string]any{"type": "LoadImage", "widgets_values": []any{"first.png"}},
		map[string]any{"type": "LoadImage", "widgets_values": []any{"old.png"}},
	}}
	got, err := ApplyOperations(document, operations)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"nodes": []any{
		map[string]any{"type": "LoadImage", "widgets_values": []any{"first.png"}},
		map[string]any{"type": "LoadImage", "widgets_values": []any{"second.png"}},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ApplyOperations() = %#v, want %#v", got, want)
	}
}
