# Anima

Files:

- anima.json — saved ComfyUI workflow
- anima.args.yaml — friendly aliases for the prompt fields

Run from the repository root after building go-comfy-cli.

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\anima\anima.json --args-file .\workflows\anima\anima.args.yaml --set 'positive_prompt=a cheerful chibi girl holding a small yellow umbrella, clean anime illustration' --set 'negative_prompt=blurry, text, watermark' --output-folder .\out\anima-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\anima\anima.json --replace-json '.nodes[] | select(.id == 90) | .widgets_values[0]::"A cheerful chibi girl holding a small yellow umbrella, clean anime illustration"' --output-folder .\out\anima-replace
~~~
