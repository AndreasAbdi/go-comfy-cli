package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go-comfy-cli/internal/workflows"
)

// NewApp builds the go-comfy-cli command tree.
func NewApp() *cli.App {
	return &cli.App{
		Name:                      "go-comfy-cli",
		Usage:                     "Run and inspect ComfyUI Desktop workflows",
		Version:                   "0.1.0",
		HelpName:                  "go-comfy-cli",
		UsageText:                 "go-comfy-cli workflow list",
		DisableSliceFlagSeparator: true,
		Commands: []*cli.Command{
			workflowCommand(),
			runCommand(),
		},
	}
}

func workflowCommand() *cli.Command {
	return &cli.Command{
		Name:    "workflow",
		Aliases: []string{"workflows"},
		Usage:   "Manage ComfyUI workflows",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List saved workflow JSON files",
				Flags: []cli.Flag{
					workflowDirectoryFlag(),
				},
				Action: listWorkflows,
			},
			{
				Name:      "download",
				Usage:     "Write a saved workflow JSON definition to stdout",
				ArgsUsage: "<name>",
				Flags: []cli.Flag{
					workflowDirectoryFlag(),
				},
				Action: downloadWorkflow,
			},
			{
				Name:      "upload",
				Usage:     "Upload a workflow JSON definition from stdin or a file",
				ArgsUsage: "<name>",
				Flags: []cli.Flag{
					workflowDirectoryFlag(),
					&cli.StringFlag{
						Name:    "input-file",
						Aliases: []string{"i"},
						Usage:   "read the workflow definition from this file instead of stdin",
					},
				},
				Action: uploadWorkflow,
			},
		},
	}
}

func workflowDirectoryFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Usage:   "workflow directory (defaults to ComfyUI Desktop in the user's home directory)",
	}
}

func workflowDirectory(ctx *cli.Context) (string, error) {
	dir := strings.TrimSpace(ctx.String("dir"))
	if dir != "" {
		return dir, nil
	}
	dir, err := workflows.DefaultDir()
	if err != nil {
		return "", fmt.Errorf("find user home directory: %w", err)
	}
	return dir, nil
}

func listWorkflows(ctx *cli.Context) error {
	dir, err := workflowDirectory(ctx)
	if err != nil {
		return err
	}

	registry := workflows.NewRegistry(dir)
	entries, err := registry.List()
	if err != nil {
		return fmt.Errorf("list workflows in %q: %w", dir, err)
	}

	for _, entry := range entries {
		fmt.Fprintln(ctx.App.Writer, entry.Name)
	}

	return nil
}

func downloadWorkflow(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("workflow download requires exactly one name")
	}
	dir, err := workflowDirectory(ctx)
	if err != nil {
		return err
	}

	workflow, err := workflows.NewRegistry(dir).GetWorkflow(ctx.Args().First())
	if err != nil {
		return fmt.Errorf("get workflow %q: %w", ctx.Args().First(), err)
	}
	if _, err := ctx.App.Writer.Write(workflow.Definition); err != nil {
		return fmt.Errorf("write workflow %q: %w", workflow.Name, err)
	}
	return nil
}

func uploadWorkflow(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("workflow upload requires exactly one name")
	}

	var reader io.Reader = ctx.App.Reader
	inputFile := strings.TrimSpace(ctx.String("input-file"))
	if inputFile != "" && inputFile != "-" {
		file, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("open input file %q: %w", inputFile, err)
		}
		defer file.Close()
		reader = file
	}
	if reader == nil {
		reader = os.Stdin
	}

	definition, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read workflow definition: %w", err)
	}
	dir, err := workflowDirectory(ctx)
	if err != nil {
		return err
	}
	workflow, err := workflows.NewRegistry(dir).Upload(ctx.Args().First(), definition)
	if err != nil {
		return fmt.Errorf("upload workflow %q: %w", ctx.Args().First(), err)
	}

	fmt.Fprintln(ctx.App.Writer, workflow.Path)
	return nil
}

// Run is separated from main so it can be exercised by callers and tests.
func Run(args []string) error {
	app := NewApp()
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	return app.Run(args)
}
