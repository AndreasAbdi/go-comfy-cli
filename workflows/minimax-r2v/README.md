# MiniMax Reference Video to Video

Files:

- minimax-r2v.json — saved ComfyUI workflow
- minimax-r2v.args.yaml — aliases for the reference video and prompt
- power-headshake.md — the long prompt used by the example

## Args-file example

The .mp4 and .md values are resolved automatically. Media files are uploaded
when they are outside ComfyUI's input directory; Markdown files are read as
prompt text.

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v\minimax-r2v.json --args-file .\workflows\minimax-r2v\minimax-r2v.args.yaml --set 'video=chainsaw-man-power-ezgif.com-gif-to-mp4-converter.mp4' --set 'prompt=.\workflows\minimax-r2v\power-headshake.md' --output-folder .\out\minimax-r2v-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v\minimax-r2v.json --replace-json '.nodes[] | select(.type == "LoadVideo") | .widgets_values[0]::"chainsaw-man-power-ezgif.com-gif-to-mp4-converter.mp4"' --replace-json '.nodes[] | select(.title == "Input Text (Prompt)") | .widgets_values[0]::"Makima performs a short irritated finger-wagging gesture beside her head in a rough manga animation style"' --output-folder .\out\minimax-r2v-replace
~~~
