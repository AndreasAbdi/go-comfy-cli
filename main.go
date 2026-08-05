package main

import (
	"fmt"
	"os"

	"go-comfy-cli/internal/commands"
)

func main() {
	app := commands.NewApp()
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
