package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go-comfy-cli/internal/media"
	"go-comfy-cli/internal/workflows"
)

const defaultComfyUIURL = "http://127.0.0.1:8000"

func runCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run a saved ComfyUI workflow",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"named"},
				Usage:   "workflow name, with or without the .json extension",
			},
			&cli.StringFlag{
				Name:  "workflow",
				Usage: "path to a workflow JSON file",
			},
			&cli.StringSliceFlag{
				Name:  "replace-json",
				Usage: "replace a jq-selected value: SELECTOR::JSON_VALUE (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:  "set",
				Usage: "set an args-file alias: ALIAS=VALUE or ALIAS::VALUE (repeatable)",
			},
			&cli.StringFlag{
				Name:  "args-file",
				Usage: "YAML alias mapping used by --set",
			},
			&cli.StringFlag{
				Name:    "url",
				Usage:   "ComfyUI server URL (defaults to COMFYUI_URL or the local Desktop server)",
				EnvVars: []string{"COMFYUI_URL"},
				Value:   defaultComfyUIURL,
			},
			&cli.StringFlag{
				Name:  "output-folder",
				Usage: "copy completed output files to this directory",
			},
		},
		Action: runWorkflow,
	}
}

func runWorkflow(ctx *cli.Context) error {
	workflow, err := selectedWorkflow(ctx)
	if err != nil {
		return err
	}

	client := comfyClient{
		baseURL: strings.TrimRight(ctx.String("url"), "/"),
		http:    &http.Client{Timeout: 10 * time.Minute},
	}
	uploader, err := media.NewUploader(client.baseURL, client.http)
	if err != nil {
		return fmt.Errorf("prepare workflow media: %w", err)
	}

	operations, err := replacementOperations(ctx)
	if err != nil {
		return err
	}
	if len(operations) > 0 {
		operations, err = uploader.PrepareOperations(ctx.Context, operations)
		if err != nil {
			return fmt.Errorf("prepare workflow media: %w", err)
		}
		workflow.Definition, err = workflows.Transpile(workflow.Definition, operations)
		if err != nil {
			return fmt.Errorf("transform workflow %q: %w", workflow.Name, err)
		}
	}

	prompt, err := workflows.Prompt(workflow)
	if err != nil {
		return fmt.Errorf("prepare workflow %q: %w", workflow.Name, err)
	}
	prompt, err = uploader.PreparePrompt(ctx.Context, prompt)
	if err != nil {
		return fmt.Errorf("prepare workflow media: %w", err)
	}

	promptID, nodeErrors, err := client.QueuePrompt(ctx.Context, prompt)
	if err != nil {
		return fmt.Errorf("queue workflow %q: %w", workflow.Name, err)
	}
	if len(nodeErrors) > 0 {
		data, marshalErr := json.Marshal(nodeErrors)
		if marshalErr != nil {
			return fmt.Errorf("workflow %q rejected with node errors: %v", workflow.Name, nodeErrors)
		}
		return fmt.Errorf("workflow %q rejected with node errors: %s", workflow.Name, data)
	}

	fmt.Fprintf(ctx.App.Writer, "queued %s\n", promptID)

	result, err := client.Wait(ctx.Context, promptID)
	if err != nil {
		return fmt.Errorf("run workflow %q: %w", workflow.Name, err)
	}
	if outputFolder := strings.TrimSpace(ctx.String("output-folder")); outputFolder != "" {
		downloaded, err := client.DownloadOutputs(ctx.Context, result["outputs"], outputFolder)
		if err != nil {
			return fmt.Errorf("save workflow outputs: %w", err)
		}
		result["downloaded_outputs"] = downloaded
	}

	result["prompt_id"] = promptID
	result["workflow"] = workflow.Name
	encoder := json.NewEncoder(ctx.App.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func selectedWorkflow(ctx *cli.Context) (workflows.Workflow, error) {
	return selectWorkflow(ctx.String("name"), ctx.String("workflow"))
}

func selectWorkflow(rawName, rawPath string) (workflows.Workflow, error) {
	name := strings.TrimSpace(rawName)
	path := strings.TrimSpace(rawPath)
	switch {
	case name == "" && path == "":
		return workflows.Workflow{}, errors.New("exactly one of --name/--named or --workflow is required")
	case name != "" && path != "":
		return workflows.Workflow{}, errors.New("--name/--named and --workflow are mutually exclusive")
	case path != "":
		workflow, err := workflows.LoadWorkflow(path)
		if err != nil {
			return workflows.Workflow{}, fmt.Errorf("load workflow file %q: %w", path, err)
		}
		return workflow, nil
	default:
		dir, err := workflows.DefaultDir()
		if err != nil {
			return workflows.Workflow{}, fmt.Errorf("find user home directory: %w", err)
		}
		workflow, err := workflows.NewRegistry(dir).GetWorkflow(name)
		if err != nil {
			return workflows.Workflow{}, fmt.Errorf("get workflow %q: %w", name, err)
		}
		return workflow, nil
	}
}

func replacementOperations(ctx *cli.Context) ([]workflows.Operation, error) {
	operations := make([]workflows.Operation, 0)
	for _, raw := range ctx.StringSlice("replace-json") {
		operation, err := workflows.ParseJSONReplacement(raw)
		if err != nil {
			return nil, fmt.Errorf("--replace-json: %w", err)
		}
		operations = append(operations, operation)
	}

	assignments := ctx.StringSlice("set")
	argsFile := strings.TrimSpace(ctx.String("args-file"))
	if len(assignments) > 0 && argsFile == "" {
		return nil, errors.New("--set requires --args-file")
	}
	if len(assignments) == 0 && argsFile != "" {
		return nil, errors.New("--args-file requires at least one --set assignment")
	}
	if len(assignments) > 0 {
		mapping, err := workflows.LoadMapping(argsFile)
		if err != nil {
			return nil, err
		}
		aliasOperations, err := mapping.AliasOperations(assignments)
		if err != nil {
			return nil, fmt.Errorf("--set: %w", err)
		}
		operations = append(operations, aliasOperations...)
	}

	return operations, nil
}

type comfyClient struct {
	baseURL string
	http    *http.Client
}

type queueResponse struct {
	PromptID   string         `json:"prompt_id"`
	NodeErrors map[string]any `json:"node_errors"`
	Error      any            `json:"error"`
}

func (c comfyClient) QueuePrompt(ctx context.Context, prompt map[string]any) (string, map[string]any, error) {
	payload := map[string]any{
		"prompt":    prompt,
		"client_id": "go-comfy-cli",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/prompt", strings.NewReader(string(body)))
	if err != nil {
		return "", nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	var result queueResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("decode response (%s): %w", response.Status, err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		if len(result.NodeErrors) > 0 {
			data, marshalErr := json.Marshal(result.NodeErrors)
			if marshalErr == nil {
				return "", result.NodeErrors, fmt.Errorf("server returned %s: node errors: %s", response.Status, data)
			}
		}
		return "", nil, fmt.Errorf("server returned %s: %v", response.Status, result.Error)
	}
	if result.PromptID == "" {
		return "", result.NodeErrors, errors.New("server response did not include prompt_id")
	}

	return result.PromptID, result.NodeErrors, nil
}

func (c comfyClient) Wait(ctx context.Context, promptID string) (map[string]any, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		result, done, err := c.history(ctx, promptID)
		if err != nil {
			return nil, err
		}
		if done {
			return result, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c comfyClient) history(ctx context.Context, promptID string) (map[string]any, bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/history/"+promptID, nil)
	if err != nil {
		return nil, false, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, false, fmt.Errorf("history request returned %s", response.Status)
	}

	var history map[string]map[string]any
	if err := json.NewDecoder(response.Body).Decode(&history); err != nil {
		return nil, false, fmt.Errorf("decode history response: %w", err)
	}
	record, ok := history[promptID]
	if !ok {
		return nil, false, nil
	}

	status, _ := record["status"].(map[string]any)
	if statusString, _ := status["status_str"].(string); statusString == "error" {
		return nil, false, fmt.Errorf("ComfyUI reported an execution error: %v", status)
	}
	if completed, _ := status["completed"].(bool); completed {
		result := map[string]any{
			"status": status["status_str"],
		}
		if outputs, ok := record["outputs"]; ok {
			result["outputs"] = outputs
		}
		return result, true, nil
	}

	return record, false, nil
}
