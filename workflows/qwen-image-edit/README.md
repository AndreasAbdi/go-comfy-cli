# Qwen Image Edit

Files:

- qwen-image-edit.json — saved ComfyUI workflow
- qwen-image-edit.args.yaml — aliases for input, reference, and prompt values
- source.jpg — bundled input image
- reference.jpg — bundled reference image

Local file references are resolved from the current working directory, then
ComfyUI's input directory. From the repository root, use the explicit paths
below for the bundled images.

The checked-in workflow contains an input image and a reference image slot. Set
both aliases when you want the source images to be explicit.

## Args-file example: input image, reference image, and prompt

~~~powershell
go-comfy-cli run --workflow .\workflows\qwen-image-edit\qwen-image-edit.json --args-file .\workflows\qwen-image-edit\qwen-image-edit.args.yaml --set 'input_image=.\workflows\qwen-image-edit\source.jpg' --set 'reference_image=.\workflows\qwen-image-edit\reference.jpg' --set 'positive_prompt=Replace the subject clothing with a bright yellow raincoat while preserving the face, pose, and composition' --output-folder .\out\qwen-image-edit-args
~~~

## JSON replacement example: input image, reference image, and prompt

~~~powershell
go-comfy-cli run --workflow .\workflows\qwen-image-edit\qwen-image-edit.json --replace-json '.nodes[] | select(.id == 41) | .widgets_values[0]::".\\workflows\\qwen-image-edit\\source.jpg"' --replace-json '.nodes[] | select(.id == 83) | .widgets_values[0]::".\\workflows\\qwen-image-edit\\reference.jpg"' --replace-json '.nodes[] | select(.id == 170) | .widgets_values[0]::"Replace the subject clothing with a bright yellow raincoat while preserving the face, pose, and composition"' --output-folder .\out\qwen-image-edit-replace
~~~
