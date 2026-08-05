package workflows

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAceStepAudioPackageAliasesReachPromptInputs(t *testing.T) {
	packageDir := filepath.Join("..", "..", "workflows", "acestepaudio1.5-audiogenerate")
	workflowPath := filepath.Join(packageDir, "acestepaudio1.5-audiogenerate.json")
	argsPath := filepath.Join(packageDir, "acestepaudio1.5-audiogenerate.args.yaml")

	definition, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	mapping, err := LoadMapping(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	operations, err := mapping.AliasOperations([]string{
		"length=8",
		"audio_prompt=upbeat synthwave with punchy drums",
		"conditioning=[Intro]\\nInstrumental",
	})
	if err != nil {
		t.Fatal(err)
	}
	transpiled, err := Transpile(definition, operations)
	if err != nil {
		t.Fatal(err)
	}

	prompt, err := Prompt(Workflow{Definition: transpiled})
	if err != nil {
		t.Fatal(err)
	}

	latentInputs := prompt["98"].(map[string]any)["inputs"].(map[string]any)
	if latentInputs["seconds"] != float64(8) {
		t.Fatalf("latent seconds = %#v, want 8", latentInputs["seconds"])
	}
	encoderInputs := prompt["94"].(map[string]any)["inputs"].(map[string]any)
	if encoderInputs["tags"] != "upbeat synthwave with punchy drums" {
		t.Fatalf("audio prompt = %#v", encoderInputs["tags"])
	}
	if encoderInputs["lyrics"] != "[Intro]\\nInstrumental" {
		t.Fatalf("conditioning = %#v", encoderInputs["lyrics"])
	}
	if encoderInputs["bpm"] != float64(190) || encoderInputs["duration"] != float64(8) {
		t.Fatalf("AceStep numeric inputs = %#v", encoderInputs)
	}
	if encoderInputs["timesignature"] != "4" || encoderInputs["language"] != "en" || encoderInputs["keyscale"] != "E minor" {
		t.Fatalf("AceStep combo inputs = %#v", encoderInputs)
	}

	if _, err := json.Marshal(prompt); err != nil {
		t.Fatalf("flattened prompt is not JSON-serializable: %v", err)
	}
}
