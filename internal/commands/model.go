package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"go-comfy-cli/internal/models"
	"go-comfy-cli/internal/workflows"
)

func modelCommand() *cli.Command {
	return &cli.Command{
		Name:  "model",
		Usage: "Manage ComfyUI models",
		Subcommands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Download models referenced by a workflow",
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
					workflowDirectoryFlag(),
					&cli.StringFlag{
						Name:    "models-directory",
						Aliases: []string{"models-dir"},
						EnvVars: []string{"COMFYUI_MODELS_DIRECTORY"},
						Usage:   "ComfyUI models directory (defaults to the Desktop models directory)",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "list model downloads without making network or filesystem changes",
					},
				},
				Action: downloadModels,
			},
		},
	}
}

func downloadModels(ctx *cli.Context) error {
	workflow, err := selectedModelWorkflow(ctx)
	if err != nil {
		return err
	}

	references, err := models.Extract(workflow.Definition)
	if err != nil {
		return fmt.Errorf("parse model metadata in workflow %q: %w", workflow.Name, err)
	}
	if len(references) == 0 {
		return fmt.Errorf("workflow %q does not contain any downloadable models", workflow.Name)
	}

	modelsDirectory := strings.TrimSpace(ctx.String("models-directory"))
	if modelsDirectory == "" {
		modelsDirectory, err = models.DefaultDirectory()
		if err != nil {
			return fmt.Errorf("find ComfyUI models directory: %w", err)
		}
	}

	results, err := (models.Downloader{
		ModelsDirectory: modelsDirectory,
		DryRun:          ctx.Bool("dry-run"),
	}).Download(ctx.Context, references)
	if err != nil {
		return fmt.Errorf("download models for workflow %q: %w", workflow.Name, err)
	}

	for _, result := range results {
		status := "downloaded"
		if result.Skipped {
			status = "skipped"
		} else if ctx.Bool("dry-run") {
			status = "would download"
		}
		fmt.Fprintf(ctx.App.Writer, "%s %s -> %s\n", status, result.Model.Name, result.Path)
	}
	return nil
}

func selectedModelWorkflow(ctx *cli.Context) (workflows.Workflow, error) {
	name := strings.TrimSpace(ctx.String("name"))
	path := strings.TrimSpace(ctx.String("workflow"))
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
		dir, err := workflowDirectory(ctx)
		if err != nil {
			return workflows.Workflow{}, err
		}
		workflow, err := workflows.NewRegistry(dir).GetWorkflow(name)
		if err != nil {
			return workflows.Workflow{}, fmt.Errorf("get workflow %q: %w", name, err)
		}
		return workflow, nil
	}
}
