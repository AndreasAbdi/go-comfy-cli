# Qwen Image Edit

Files:

- qwen-image-edit.json — saved ComfyUI workflow
- qwen-image-edit.args.yaml — aliases for input, reference, and prompt values

The checked-in workflow contains an input image and a reference image slot. Set
both aliases when you want the source images to be explicit.

## Args-file example: input image, reference image, and prompt

~~~powershell
go-comfy-cli run --workflow .\workflows\qwen-image-edit\qwen-image-edit.json --args-file .\workflows\qwen-image-edit\qwen-image-edit.args.yaml --set 'input_image=a4f99e6ea2d1960a5e2b48c5f5193bac.jpg' --set 'reference_image=220205-Music-Core-Twitter-Update-Happy-Birthday-Kim-Minju-documents-1.jpeg' --set 'positive_prompt=Replace the subject clothing with a bright yellow raincoat while preserving the face, pose, and composition' --output-folder .\out\qwen-image-edit-args
~~~

## JSON replacement example: input image, reference image, and prompt

~~~powershell
go-comfy-cli run --workflow .\workflows\qwen-image-edit\qwen-image-edit.json --replace-json '.nodes[] | select(.id == 41) | .widgets_values[0]::"a4f99e6ea2d1960a5e2b48c5f5193bac.jpg"' --replace-json '.nodes[] | select(.id == 83) | .widgets_values[0]::"220205-Music-Core-Twitter-Update-Happy-Birthday-Kim-Minju-documents-1.jpeg"' --replace-json '.nodes[] | select(.id == 170) | .widgets_values[0]::"Replace the subject clothing with a bright yellow raincoat while preserving the face, pose, and composition"' --output-folder .\out\qwen-image-edit-replace
~~~
