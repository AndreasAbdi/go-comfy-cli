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

func TestPromptResolvesPrimitiveNode(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{
      "nodes": [
        {"id": 10, "type": "PrimitiveNode", "inputs": [], "outputs": [{"name": "STRING", "type": "STRING", "links": [1]}], "widgets_values": ["hello"]},
        {"id": 20, "type": "Sink", "inputs": [{"name": "value", "type": "STRING", "link": 1}], "widgets_values": []}
      ],
      "links": [[1, 10, 0, 20, 0, "STRING"]]
    }`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"20": map[string]any{"class_type": "Sink", "inputs": map[string]any{"value": "hello"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Prompt() = %#v, want %#v", got, want)
	}
}

func TestPromptResolvesTypedPrimitiveNodes(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{
      "nodes": [
        {"id": 10, "type": "PrimitiveStringMultiline", "inputs": [], "outputs": [{"name": "STRING", "type": "STRING", "links": [1]}], "widgets_values": ["hello"]},
        {"id": 11, "type": "PrimitiveFloat", "inputs": [], "outputs": [{"name": "FLOAT", "type": "FLOAT", "links": [2]}], "widgets_values": [5]},
        {"id": 20, "type": "Sink", "inputs": [{"name": "text", "type": "STRING", "link": 1}, {"name": "duration", "type": "FLOAT", "link": 2}], "widgets_values": []}
      ],
      "links": [[1, 10, 0, 20, 0, "STRING"], [2, 11, 0, 20, 1, "FLOAT"]]
    }`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"20": map[string]any{"class_type": "Sink", "inputs": map[string]any{"text": "hello", "duration": float64(5)}},
	}
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

func TestPromptSkipsKSamplerSeedControlWidget(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{
      "nodes": [{
        "id": 1,
        "type": "KSampler",
        "inputs": [
          {"name": "seed", "type": "INT", "widget": {"name": "seed"}, "link": null},
          {"name": "steps", "type": "INT", "widget": {"name": "steps"}, "link": null},
          {"name": "cfg", "type": "FLOAT", "widget": {"name": "cfg"}, "link": null},
          {"name": "sampler_name", "type": "COMBO", "widget": {"name": "sampler_name"}, "link": null},
          {"name": "scheduler", "type": "COMBO", "widget": {"name": "scheduler"}, "link": null},
          {"name": "denoise", "type": "FLOAT", "widget": {"name": "denoise"}, "link": null}
        ],
        "widgets_values": [123, "randomize", 30, 4, "euler", "simple", 1]
      }]
    }`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	inputs := got["1"].(map[string]any)["inputs"].(map[string]any)
	if inputs["sampler_name"] != "euler" || inputs["scheduler"] != "simple" || inputs["denoise"] != float64(1) {
		t.Fatalf("KSampler inputs = %#v", inputs)
	}
}

func TestPromptResolvesBypassedNode(t *testing.T) {
	workflow := Workflow{Definition: json.RawMessage(`{
      "nodes": [
        {"id": 1, "type": "Source", "inputs": [], "outputs": [{"name": "IMAGE", "type": "IMAGE", "links": [1]}], "widgets_values": []},
        {"id": 2, "type": "ImageScaleToTotalPixels", "mode": 4, "inputs": [{"name": "image", "type": "IMAGE", "link": 1}], "outputs": [{"name": "IMAGE", "type": "IMAGE", "links": [2]}], "widgets_values": []},
        {"id": 3, "type": "Sink", "inputs": [{"name": "image", "type": "IMAGE", "link": 2}], "widgets_values": []}
      ],
      "links": [[1, 1, 0, 2, 0, "IMAGE"], [2, 2, 0, 3, 0, "IMAGE"]]
    }`)}

	got, err := Prompt(workflow)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"1": map[string]any{"class_type": "Source", "inputs": map[string]any{}},
		"3": map[string]any{"class_type": "Sink", "inputs": map[string]any{"image": []any{"1", 0}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Prompt() = %#v, want %#v", got, want)
	}
}
