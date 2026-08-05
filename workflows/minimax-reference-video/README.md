# MiniMax Reference Video to Video

Files:

- minimax-reference-video.json — saved ComfyUI workflow
- minimax-reference-video.args.yaml — aliases for the reference video and prompt
- power-prompt.md — the long prompt used by the example
- power-reference.mp4 — bundled reference video

Local file references are resolved from the current working directory, then
ComfyUI's input directory. From the repository root, use the explicit path for
the bundled video:

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-reference-video\minimax-reference-video.json --args-file .\workflows\minimax-reference-video\minimax-reference-video.args.yaml --set 'video=.\workflows\minimax-reference-video\power-reference.mp4' --output-folder .\out\minimax-reference-video-default
~~~

## Args-file example

The .mp4 and .md values are resolved automatically. Media files are uploaded
when they are outside ComfyUI's input directory; Markdown files are read as
prompt text.

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-reference-video\minimax-reference-video.json --args-file .\workflows\minimax-reference-video\minimax-reference-video.args.yaml --set 'video=.\workflows\minimax-reference-video\power-reference.mp4' --set 'prompt=.\workflows\minimax-reference-video\power-prompt.md' --output-folder .\out\minimax-reference-video-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-reference-video\minimax-reference-video.json --replace-json '.nodes[] | select(.type == "LoadVideo") | .widgets_values[0]::".\\workflows\\minimax-reference-video\\power-reference.mp4"' --replace-json '.nodes[] | select(.title == "Input Text (Prompt)") | .widgets_values[0]::"Makima performs a short irritated finger-wagging gesture beside her head in a rough manga animation style"' --output-folder .\out\minimax-reference-video-replace
~~~
