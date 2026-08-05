package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go-comfy-cli/workflows"
)

func newApp() *cli.App {
	return &cli.App{
		Name:      "go-comfy-cli",
		Usage:     "List ComfyUI Desktop workflows",
		Version:   "0.1.0",
		HelpName:  "go-comfy-cli",
		UsageText: "go-comfy-cli workflows list",
		Commands: []*cli.Command{
			{
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

	entries, err := workflows.List(dir)
	if err != nil {
		return fmt.Errorf("list workflows in %q: %w", dir, err)
	}

	for _, entry := range entries {
		fmt.Fprintln(ctx.App.Writer, entry)
	}

	return nil
}

func main() {
	app := newApp()
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if exitErr, ok := err.(cli.ExitCoder); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
