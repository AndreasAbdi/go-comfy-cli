# tv-girl
use comfyui, a tv, a cron, and an llm to trigger girl on tv. 

# what? 

This project is roughly a set of instructions and code that an agent can run against to generate images on your TV via chromecast. 

# How to use
1. clone repo
2. point an agent to this repo, ask it to use it. 

# required:
1. some agentic harness (agent)
2. comfyui installation with ability to run a minimax h3 (RTX 4090/5090/etc)
3. television with chromecast install
4. some runtime daemon that can execute

# references: 
1. https://github.com/erkstruwe/chromecast-cli
2. https://docs.comfy.org/tutorials/video/minimax/minimax-h3#minimax-h3-reference-to-video-r2v

## go-comfy-cli

The repository includes a small Go CLI for discovering ComfyUI Desktop workflows.

```powershell
go run . workflows list
```

By default, it lists `.json` workflows from:

```text
<user home>\Documents\ComfyUI\user\default\workflows
```

Use `--dir` (or `-d`) to list another workflow directory:

```powershell
go run . workflows list --dir C:\path\to\workflows
```

Run a named workflow through the local ComfyUI Desktop API:

```powershell
go run . run --named minimax-i2v
```

The default server is `http://127.0.0.1:8000`; override it with `--url` or
`COMFYUI_URL`. The current run command intentionally supports named workflows
only. A future `--file` option is reserved for direct workflow-file execution.

Raw JQ replacements use `SELECTOR::JSON_VALUE` and can be repeated:

```powershell
go run . run --name minimax-i2v `
  --replace-json '.nodes[] | select(.id == 114) | .widgets_values[0]::"image.png"'
```

For friendlier arguments, create an args mapping file:

```yaml
version: 1
aliases:
  image:
    selector: '.nodes[] | select(.type == "LoadImage") | .widgets_values[0]'
    type: string
    cardinality: many
  output:
    selector: '.nodes[] | select(.type == "SaveVideo") | .widgets_values[0]'
    type: string
    cardinality: one
```

Then apply aliases with `--set`:

```powershell
go run . run --name minimax-i2v `
  --args-file minimax-map.yaml `
  --set 'image=transparent_rgb_gaming_mouse.png' `
  --set 'output=video/result'
```

Alias values are typed from the mapping (`string`, `number`, `integer`,
`boolean`, or `json`). Indexed aliases can define `indexed_selector` and use
the `alias[0]` syntax; `alias[]` applies the replacement to every match.

Local file values are prepared automatically. Existing files under the local
ComfyUI `input` directory are referenced by their input-relative name; other
media files are uploaded through `/upload/image` and replaced with the returned
ComfyUI input path. Existing `.md`, `.markdown`, and `.txt` files are read into
the replacement value, which is useful for long prompt aliases.

For example, the included R2V mapping can replace the reference video with a
local file path:

```powershell
go run . run --name minimax-r2v `
  --args-file examples/minimax-r2v.args.yaml `
  --set 'video=C:\path\to\reference.mp4'
```
