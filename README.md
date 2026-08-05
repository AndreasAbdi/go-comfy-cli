# go-comfy-cli

go-comfy-cli runs comfy UI from the CLI on an existing Comfy UI installed.

# Prerequisites

- ComfyUI Desktop running locally.
- The models and custom nodes required by the selected workflow have to be preinstalled. 

# Build and install

## ez build

### install
windows: 
```powershell
irm https://github.com/portpowered/you-agent-factory/releases/latest/download/install.ps1 | iex
```

or just download the binaries from the latest release file. 

### running
#### install models if not downloaded
```
go-comfy-cli model download --workflow .\workflows\anima\anima.json
```
#### run the command

```
go-comfy-cli run --workflow  ./workflows/anima/anima.json`
  --args-file ./workflows/anima/anima.json`
  --set 'positive_prompt=a picture of a banana' `
```

### manual build
~~~powershell
go build -o go-comfy-cli.exe .
~~~

On Windows, place go-comfy-cli.exe on PATH, or invoke it as
.\go-comfy-cli.exe from the directory containing the binary.



## Command reference

| Command | Function |
| --- | --- |
| go-comfy-cli workflow list | List JSON workflows in comfyUI's user workflow directory by the name |
| go-comfy-cli workflow download <name> | Write a named workflow JSON definition to stdout |
| go-comfy-cli workflow upload <name> --input-file <file.json> | Upload a workflow JSON definition |
| go-comfy-cli model download --workflow <file.json> | Download models listed in a workflow |
| go-comfy-cli model download --name <name> | Download models listed in a named workflow |
| go-comfy-cli run --name <name> | Run a workflow by name from the Desktop workflow directory |
| go-comfy-cli run --workflow <file.json> | Run a workflow from an explicit JSON path |
| go-comfy-cli run --args-file <file.yaml> --set ALIAS=VALUE | Apply aliases from a mapping file |
| go-comfy-cli run --replace-json 'SELECTOR::JSON_VALUE' | Apply one or more raw jq replacements on the json template file |

### function notes
1. Exactly one of --name/--named or --workflow is required. 
2. Both --args-file/--set and --replace-json can be used with either workflow selector. 
3. Use --output-folder to copy completed output files locally.


## Workflow packages and examples

Each workflow keeps its JSON, args mapping, prompt inputs, and runnable
examples together:

- [Anima](workflows/anima/README.md)
- [Qwen Image Edit](workflows/qwen-image-edit/README.md)
- [MiniMax Image to Video](workflows/minimax-i2v/README.md)
- [MiniMax Reference Video to Video](workflows/minimax-r2v/README.md)
- [MiniMax Audio Reference to Video](workflows/minimax-r2v-audio/README.md)

## Args files and JSON replacements

An args file maps a friendly alias to a jq selector:

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
go-comfy-cli run --workflow .\workflows\anima\anima.json --args-file .\workflows\anima\anima.args.yaml --set 'positive_prompt=a bright anime portrait'
~~~

Raw replacements use SELECTOR::JSON_VALUE and can be repeated:

~~~powershell
go-comfy-cli run --workflow .\workflows\anima\anima.json --replace-json '.nodes[] | select(.id == 90) | .widgets_values[0]::"a bright anime portrait"'
~~~

Local media values are resolved before the workflow is submitted. 

Existing files under ComfyUI's input directory are referenced by input-relative name. 
If you point the thing at a file, then it will download the file into your server for such things as: 
1. markdown files and text files: .md, .markdown, and .txt
2. image files (img)
3. video files (mp4)
4. audio files (wav)

You usually use the media files for image to image workflows and what not. 

Example: 

```powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v-audio\minimax-r2v-audio.json --args-file .\workflows\minimax-r2v-audio\minimax-r2v-audio.args.yaml --set 'audio=.\workflows\minimax-r2v-audio\hey_come_on_you_promised_to_take_me_shopping.wav' --set 'prompt=.\workflows\minimax-r2v-audio\minimax-r2v-audio-prompt.md' --set 'duration=5' --output-folder .\out\minimax-r2v-audio-args
```

## Development checks

~~~powershell
go test ./...
go fmt ./...
go vet ./...
go build ./...
~~~

## License

MIT
