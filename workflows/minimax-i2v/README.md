# MiniMax Image to Video

Files:

- minimax-i2v.json — saved ComfyUI workflow
- minimax-i2v.args.yaml — aliases for the first-frame image and prompt

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-i2v\minimax-i2v.json --args-file .\workflows\minimax-i2v\minimax-i2v.args.yaml --set 'input_image=transparent_rgb_gaming_mouse.png' --set 'positive_prompt=A slow product turntable of the transparent gaming mouse, with blue and amber studio lighting and crisp reflections' --output-folder .\out\minimax-i2v-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-i2v\minimax-i2v.json --replace-json '.nodes[] | select(.id == 114) | .widgets_values[0]::"transparent_rgb_gaming_mouse.png"' --replace-json '.nodes[] | select(.id == 105) | .widgets_values[0]::"A slow product turntable of the transparent gaming mouse, with blue and amber studio lighting and crisp reflections"' --output-folder .\out\minimax-i2v-replace
~~~
