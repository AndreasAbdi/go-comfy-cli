# what? 

experiment to validate that the new Minimax ref models can be usable for every day usage. 

# arch: 
1. wraps the go-chromecast CLI
2. python wrapper/cli to comfyui install with workflow for minimax h3. 
3. instruction prompt guidelines for prompt structure. 
4. workflow json schemas for reference 2 video, image 2 video, text 2 video
5. raw data samples for generation (image ref, baseline audio ref, )

# theory
- run agent, default exec via AGENTS.md/CLAUDE.md, 
- instructions (read above references, guidance from AGENTS.md to retrieve constructed references, use it to generate the video from ref, analyze output, validate outputs, repeat until visually correct)

## theoretical problems
1. context window for video understanding is expensive, should likely use subagents for artifact analysis/validation rather than direct visualization. 
2. long durational execution means that harness waits are going to be expensive. (480p at 5 mins)
3. durational execution is expensive


# notes on implementation

## image refs: 
1. videos/refs/images must be loaded into the appropriate directory, or alternatively via the CLI.
C:\Users\<user>\AppData\Local\Comfy-Desktop\ComfyUI-Shared\input

## code refs: 
- https://github.com/Comfy-Org/ComfyUI/blob/6f7cd7fceaaf60d2669b554936394a7412c6fde5/server.py#L4 (Comfy YUI server)
- https://github.com/StableCanvas/comfyui-client/blob/main/packages/cli/readme.md (API server client )
- comfy ui CLI (https://github.com/Comfy-Org/comfy-cli)
- comfyui openclaw skils: https://github.com/HuangYuChuh/ComfyUI_Skills_OpenClaw
- https://github.com/HuangYuChuh/ComfyUI_Skill_CLI/tree/main, agentic variant
- comfyui APi documentation https://docs.comfy.org/development/comfyui-server/comms_routes

# general notes: 
1. as of today, there is no one definitive CLI for generating the comfy ui images/video from a comfyui installation. 
1.1. there are a few different mechanisms for generation via the comfyui CLI and random other comfyui skills but nothing properly wraps it for long term video generation and tools; you sort of need to piece it together. 

how to use today: 
1. call CLI to upload your reference content
2. TODO: figure out how to upload audio/video content (lol: its just the images API)
3. invoke the CLI to runt he workflow definition with the reference images/names

GOTCHAS: 
1. none fo the existing CLIs document how to upload user content or whatever..
2. none of the existing CLIs document how to pass in args or how to format execution for a given workflow

## WTF: 
- the image API is the one used for general content uploading, how strange. 
- upload via the image API for botht the images and videos. 

## where is the included files for input folder? 

C:\Users\<user>\Documents\ComfyUI\input
C:\Users\<workflows>\Documents\ComfyUI\user\default\workflows


## asdasd


uh bunch of docs in the 
Comfy UI files are located in usually
```text
<user home>\Documents\ComfyUI\user\default\workflows
```


Download a workflow definition to stdout:

```powershell
go run . workflow download minimax-i2v > minimax-i2v.json
```

Upload a workflow definition from stdin or from a file. The `.json` suffix is
added when it is omitted:

```powershell
Get-Content .\minimax-i2v.json -Raw | go run . workflow upload minimax-i2v
go run . workflow upload minimax-i2v --input-file .\minimax-i2v.json
```

Both commands accept `--dir` when working with a non-default workflow
directory.

Checked-in workflow bundles are available in [`workflows/`](workflows/). Each
workflow keeps its JSON definition, args mapping, and examples together:

- [`workflows/anima/`](workflows/anima/) contains `anima.json` and `anima.args.yaml`
- [`workflows/qwen-image-edit/`](workflows/qwen-image-edit/) contains `qwen-image-edit.json` and `qwen-image-edit.args.yaml`
- [`workflows/minimax-i2v/`](workflows/minimax-i2v/) contains `minimax-i2v.json`
- [`workflows/minimax-r2v/`](workflows/minimax-r2v/) contains `minimax-r2v.json`, `minimax-r2v.args.yaml`, and the `power-headshake.md` example
- [`workflows/minimax-r2v-audio/`](workflows/minimax-r2v-audio/) contains `minimax-r2v-audio.json`, `minimax-r2v-audio.args.yaml`, and `minimax-r2v-audio-prompt.md`

Run a named workflow through the local ComfyUI Desktop API:

```powershell
go run . run --named minimax-i2v
```

Copy completed files into a local directory with `--output-folder`. Any
ComfyUI output subfolder is preserved below it, and the JSON result includes
the resulting paths in `downloaded_outputs`:

```powershell
go run . run --name minimax-i2v --output-folder .\out
```

The default server is `http://127.0.0.1:8000`; override it with `--url` or
`COMFYUI_URL`. The current run command intentionally supports named workflows
only. A future `--file` option is reserved for direct workflow-file execution.
