package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go-comfy-cli/internal/workflows"
)

// NewApp builds the go-comfy-cli command tree.
func NewApp() *cli.App {
	return &cli.App{
		Name:      "go-comfy-cli",
		Usage:     "Run and inspect ComfyUI Desktop workflows",
		Version:   "0.1.0",
		HelpName:  "go-comfy-cli",
		UsageText: "go-comfy-cli workflows list",
		Commands: []*cli.Command{
			workflowCommand(),
			runCommand(),
		},
	}
}

func workflowCommand() *cli.Command {
	return &cli.Command{
		Name:  "workflows",
		Usage: "Manage ComfyUI workflows",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List saved workflow JSON files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Usage:   "workflow directory (defaults to ComfyUI Desktop in the user's home directory)",
					},
				},
				Action: listWorkflows,
			},
		},
	}
}

func listWorkflows(ctx *cli.Context) error {
	dir := ctx.String("dir")
	if dir == "" {
		var err error
		dir, err = workflows.DefaultDir()
		if err != nil {
			return fmt.Errorf("find user home directory: %w", err)
		}
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

// Run is separated from main so it can be exercised by callers and tests.
func Run(args []string) error {
	app := NewApp()
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	return app.Run(args)
}
