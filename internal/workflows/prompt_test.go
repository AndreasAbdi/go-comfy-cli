package workflows

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPromptPassesThroughAPIFormat(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{"1":{"class_type":"Test","inputs":{"value":"hello"}}}`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"1": map[string]any{"class_type": "Test", "inputs": map[string]any{"value": "hello"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Prompt() = %#v, want %#v", got, want)
	}
}

func TestPromptFlattensEmbeddedSubgraph(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{
      "nodes": [
        {"id": 10, "type": "subgraph-id", "inputs": [{"name": "text", "type": "STRING", "widget": {"name": "text"}}], "widgets_values": ["hello"]},
        {"id": 20, "type": "Sink", "inputs": [{"name": "value", "type": "STRING", "link": 3}], "widgets_values": []}
      ],
      "links": [[3, 10, 0, 20, 0, "STRING"]],
      "definitions": {"subgraphs": [{
        "id": "subgraph-id",
        "inputs": [{"name": "text", "linkIds": [1]}],
        "outputs": [{"name": "value", "linkIds": [2]}],
        "nodes": [
          {"id": 1, "type": "Source", "inputs": [{"name": "text", "type": "STRING", "link": 1}], "widgets_values": []}
        ],
        "links": [
          [1, -10, 0, 1, 0, "STRING"],
          [2, 1, 0, -20, 0, "STRING"]
        ]
      }]}
    }`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"1":  map[string]any{"class_type": "Source", "inputs": map[string]any{"text": "hello"}},
		"20": map[string]any{"class_type": "Sink", "inputs": map[string]any{"value": []any{"1", 0}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Prompt() = %#v, want %#v", got, want)
	}
}
