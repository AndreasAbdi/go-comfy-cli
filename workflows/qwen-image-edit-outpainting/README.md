# Qwen Image Edit Outpainting

This sample is the saved qwen-image-edit-outpainting ComfyUI workflow,
packaged for use with go-comfy-cli.

The args file exposes:

- input_image — source image to outpaint
- reference_text — text describing the desired edit or scene
- target_width and target_height — output canvas dimensions in pixels

## Single image

From the repository root:

~~~powershell
.\go-comfy-cli.exe run --workflow .\workflows\qwen-image-edit-outpainting\qwen-image-edit-outpainting.json --args-file .\workflows\qwen-image-edit-outpainting\qwen-image-edit-outpainting.args.yaml --set 'input_image=.\input\Makima.png' --set 'reference_text=complete the inpainting of the image such that it matches the center image' --set 'target_width=864' --set 'target_height=480' --output-folder .\out\qwen-image-edit-outpainting-480p-corrected
~~~

## Batch every input image

The following PowerShell loop runs the same workflow once per supported image
in input and writes the completed files to out\qwen-image-edit-outpainting-480p-corrected.

~~~powershell
$output = '.\out\qwen-image-edit-outpainting-480p-corrected'
New-Item -ItemType Directory -Force $output | Out-Null
Get-ChildItem '.\input' -File | Where-Object { $_.Extension -match '^\.(jpg|jpeg|png|webp|jfif)$' } | ForEach-Object {
  $imageOutput = Join-Path $output $_.BaseName
  .\go-comfy-cli.exe run --workflow .\workflows\qwen-image-edit-outpainting\qwen-image-edit-outpainting.json --args-file .\workflows\qwen-image-edit-outpainting\qwen-image-edit-outpainting.args.yaml --set "input_image=$($_.FullName)" --set 'reference_text=complete the inpainting of the image such that it matches the center image' --set 'target_width=864' --set 'target_height=480' --output-folder $imageOutput
  if ($LASTEXITCODE -ne 0) { throw "outpainting failed for $($_.Name)" }
}
~~~

The reimported workflow has its internal scale node bypassed, so the output
uses the requested 864×480 target canvas directly.
