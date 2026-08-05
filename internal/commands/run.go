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
	"go-comfy-cli/internal/workflows"
)

const defaultComfyUIURL = "http://127.0.0.1:8000"

func runCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run a saved ComfyUI workflow",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "named",
				Usage:    "workflow name, with or without the .json extension",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "url",
				Usage:   "ComfyUI server URL (defaults to COMFYUI_URL or the local Desktop server)",
				EnvVars: []string{"COMFYUI_URL"},
				Value:   defaultComfyUIURL,
			},
		},
		Action: runWorkflow,
	}
}

func runWorkflow(ctx *cli.Context) error {
	dir, err := workflows.DefaultDir()
	if err != nil {
		return fmt.Errorf("find user home directory: %w", err)
	}

	registry := workflows.NewRegistry(dir)
	workflow, err := registry.GetWorkflow(ctx.String("named"))
	if err != nil {
		return fmt.Errorf("get workflow %q: %w", ctx.String("named"), err)
	}

	prompt, err := workflows.Prompt(workflow)
	if err != nil {
		return fmt.Errorf("prepare workflow %q: %w", workflow.Name, err)
	}

	client := comfyClient{
		baseURL: strings.TrimRight(ctx.String("url"), "/"),
		http:    &http.Client{Timeout: 30 * time.Second},
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

	result["prompt_id"] = promptID
	result["workflow"] = workflow.Name
	encoder := json.NewEncoder(ctx.App.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
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
