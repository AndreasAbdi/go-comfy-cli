# MiniMax Image to Video

Files:

- minimax-image-to-video.json — saved ComfyUI workflow
- minimax-image-to-video.args.yaml — aliases for the first-frame image and prompt
- gaming-mouse.png — bundled first-frame reference image

File references are resolved relative to the directory where the CLI is
executed, followed by ComfyUI's input directory.

## Args-file example: short-name input

If `gaming-mouse.png` is in the current directory or already in ComfyUI's
input directory, use its short name:

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-image-to-video\minimax-image-to-video.json --args-file .\workflows\minimax-image-to-video\minimax-image-to-video.args.yaml --set 'input_image=gaming-mouse.png' --set 'positive_prompt=A slow product turntable of the transparent gaming mouse, with blue and amber studio lighting and crisp reflections' --output-folder .\out\minimax-image-to-video-args
~~~

When using the bundled PNG from the repository root, pass its explicit path:

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-image-to-video\minimax-image-to-video.json --args-file .\workflows\minimax-image-to-video\minimax-image-to-video.args.yaml --set 'input_image=.\workflows\minimax-image-to-video\gaming-mouse.png' --set 'positive_prompt=A slow product turntable of the transparent gaming mouse, with blue and amber studio lighting and crisp reflections' --output-folder .\out\minimax-image-to-video-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-image-to-video\minimax-image-to-video.json --replace-json '.nodes[] | select(.id == 114) | .widgets_values[0]::".\\workflows\\minimax-image-to-video\\gaming-mouse.png"' --replace-json '.nodes[] | select(.id == 105) | .widgets_values[0]::"A slow product turntable of the transparent gaming mouse, with blue and amber studio lighting and crisp reflections"' --output-folder .\out\minimax-image-to-video-replace
~~~
