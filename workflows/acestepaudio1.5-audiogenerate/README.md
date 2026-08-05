# AceStep Audio 1.5 Audio Generate

This package contains the installed `acestepaudio1.5-audiogenerate` ComfyUI
workflow and an args file for its three useful generation controls:

- `length` — song length in seconds (`duration` is an equivalent alias)
- `audio_prompt` — the AceStep style/tags prompt
- `conditioning` — the lyrics or other conditioning text

The length is applied to the shared `Song Duration` primitive, so it updates
both the empty audio latent and the AceStep text encoder.

## Args-file example

From the repository root:

~~~powershell
go-comfy-cli run --workflow .\workflows\acestepaudio1.5-audiogenerate\acestepaudio1.5-audiogenerate.json --args-file .\workflows\acestepaudio1.5-audiogenerate\acestepaudio1.5-audiogenerate.args.yaml --set 'length=8' --set 'audio_prompt=upbeat synthwave with punchy drums and a memorable bassline' --set 'conditioning=[Intro]\nInstrumental\n\n[Verse]\nNeon lights across the night' --output-folder .\out\acestepaudio1.5-audiogenerate
~~~

The workflow expects `ace_step_1.5_turbo_aio.safetensors` in ComfyUI's
`models/checkpoints` directory. It writes MP3 and FLAC outputs through the two
save-audio nodes already present in the packaged workflow.
