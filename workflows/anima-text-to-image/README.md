# Anima

Files:

- anima-text-to-image.json — saved ComfyUI workflow
- anima-text-to-image.args.yaml — friendly aliases for the prompt fields

Run from the repository root after building go-comfy-cli.

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\anima-text-to-image\anima-text-to-image.json --args-file .\workflows\anima-text-to-image\anima-text-to-image.args.yaml --set 'positive_prompt=a cheerful chibi girl holding a small yellow umbrella, clean anime illustration' --set 'negative_prompt=blurry, text, watermark' --output-folder .\out\anima-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\anima-text-to-image\anima-text-to-image.json --replace-json '.nodes[] | select(.id == 90) | .widgets_values[0]::"A cheerful chibi girl holding a small yellow umbrella, clean anime illustration"' --output-folder .\out\anima-replace
~~~
