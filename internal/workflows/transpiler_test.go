package workflows

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseJSONReplacement(t *testing.T) {
	operation, err := ParseJSONReplacement(`.settings.count::8`)
	if err != nil {
		t.Fatal(err)
	}
	if operation.Selector != ".settings.count" || operation.Value != float64(8) {
		t.Fatalf("ParseJSONReplacement() = %#v", operation)
	}
}

func TestApplyOperationsUpdatesAllJQMatches(t *testing.T) {
	document := map[string]any{
		"nodes": []any{
			map[string]any{"id": float64(1), "type": "LoadImage", "widgets_values": []any{"one.png"}},
			map[string]any{"id": float64(2), "type": "LoadImage", "widgets_values": []any{"two.png"}},
		},
	}

	operation, err := ParseJSONReplacement(`.nodes[] | select(.type == "LoadImage") | .widgets_values[0]::"updated.png"`)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ApplyOperations(document, []Operation{operation})
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]any{
		"nodes": []any{
			map[string]any{"id": float64(1), "type": "LoadImage", "widgets_values": []any{"updated.png"}},
			map[string]any{"id": float64(2), "type": "LoadImage", "widgets_values": []any{"updated.png"}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ApplyOperations() = %#v, want %#v", got, want)
	}
}

func TestApplyOperationsEnforcesCardinality(t *testing.T) {
	operation := Operation{Selector: ".nodes[]", Value: "updated", Cardinality: "one"}
	_, err := ApplyOperations(map[string]any{"nodes": []any{map[string]any{}, map[string]any{}}}, []Operation{operation})
	if err == nil {
		t.Fatal("ApplyOperations() accepted multiple matches for a one-valued operation")
	}
}

func TestTranspileReturnsJSON(t *testing.T) {
	data := json.RawMessage(`{"value":"before"}`)
	operation, err := ParseJSONReplacement(`.value::"after"`)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Transpile(data, []Operation{operation})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"value":"after"}` {
		t.Fatalf("Transpile() = %s", got)
	}
}
