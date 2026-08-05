# go-comfy-cli

go-comfy-cli is a CLI to generate images, audio, video from your command line, rather than relying on the comfy UI interface.

## Prerequisites

- ComfyUI Desktop running locally.


## Build and install

### Install a release on Windows

~~~powershell
irm https://github.com/portpowered/you-agent-factory/releases/latest/download/install.ps1 | iex
~~~

Or download from the latest release.

### Download models

If a workflow's models are not installed yet, download the models referenced by
an explicit or named workflow:

~~~powershell
go-comfy-cli model download --workflow .\workflows\anima-text-to-image\anima-text-to-image.json
go-comfy-cli model download --name anima
~~~

## Command reference

| Command | Function |
| --- | --- |
| go-comfy-cli workflow list | List JSON workflows in ComfyUI Desktop's user workflow directory |
| go-comfy-cli workflow download <name> | Write a named workflow JSON definition to stdout |
| go-comfy-cli workflow upload <name> --input-file <file.json> | Upload a workflow JSON definition, so that it can be invoked from the comfy UI website on localhost:8000 |
| go-comfy-cli model download --workflow <file.json> | Download models listed in a workflow |
| go-comfy-cli model download --name <name> | Download models listed in a named workflow |
| go-comfy-cli run --name <name> | Run a workflow by name from the Desktop workflow directory |
| go-comfy-cli run --workflow <file.json> | Run a workflow from an explicit JSON path |
| go-comfy-cli run --args-file <file.yaml> --set ALIAS=VALUE | Apply aliases from a mapping file |
| go-comfy-cli run --replace-json SELECTOR::JSON_VALUE | Apply one or more raw jq replacements on the workflow JSON |

### Run options

1. Exactly one of --name/--named or --workflow is required.
2. Both --args-file/--set and --replace-json can be used with either workflow selector.
3. Use --output-folder to copy completed output files locally.

## Workflow packages and examples

Each workflow keeps its JSON, args mapping, prompt inputs, and runnable
examples together:

- [Anima](workflows/anima-text-to-image/README.md)
- [Qwen Image Edit](workflows/qwen-image-edit/README.md)
- [MiniMax Image to Video](workflows/minimax-image-to-video/README.md)
- [MiniMax Reference Video to Video](workflows/minimax-reference-video/README.md)
- [MiniMax Audio Reference to Video](workflows/minimax-audio-reference/README.md)

All examples assume the CLI has already been built and that the command is
being run from this repository's root directory. Local file references are
resolved from the current working directory, then ComfyUI's input directory.
For media kept beside a workflow, pass the explicit `./workflows/...` path.

## Args files and JSON replacements

An args file maps a friendly alias to a [jq selector](https://jqlang.org/manual/):

~~~yaml
version: 1
aliases:
  positive_prompt:
    selector: '.nodes[] | select(.id == 90) | .widgets_values[0]'
    type: string
    cardinality: one
~~~

Use the mapping with --set:

~~~powershell
go-comfy-cli run --workflow .\workflows\anima-text-to-image\anima-text-to-image.json --args-file .\workflows\anima-text-to-image\anima-text-to-image.args.yaml --set 'positive_prompt=a bright anime portrait'
~~~

Raw replacements use SELECTOR::JSON_VALUE and can be repeated:

~~~powershell
go-comfy-cli run --workflow .\workflows\anima-text-to-image\anima-text-to-image.json --replace-json '.nodes[] | select(.id == 90) | .widgets_values[0]::"a bright anime portrait"'
~~~

Local media values are resolved before the workflow is submitted.

1. Existing files under ComfyUI's input directory are referenced by input-relative name;
2. other image, video, and audio files are uploaded through /upload/image.
3. Markdown and text files are read as literal string content, which is useful for long prompts.

Example:
~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-audio-reference\minimax-audio-reference.json --args-file .\workflows\minimax-audio-reference\minimax-audio-reference.args.yaml --set 'audio=.\workflows\minimax-audio-reference\shopping-promise.wav' --set 'prompt=.\workflows\minimax-audio-reference\minimax-audio-reference-prompt.md' --set 'duration=5' --output-folder .\out\minimax-audio-reference-args
~~~

## Development checks

~~~powershell
go test ./...
go fmt ./...
go vet ./...
go build ./...
~~~

### Build binary specficially

~~~powershell
go build -o go-comfy-cli.exe .
~~~

On Windows, place go-comfy-cli.exe on PATH, or invoke it as
.\go-comfy-cli.exe from the directory containing the binary.


## License

MIT
