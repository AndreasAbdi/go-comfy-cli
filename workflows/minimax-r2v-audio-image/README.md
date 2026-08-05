# MiniMax R2V Audio + Image

Files:

- `minimax-r2v-audio-image.json` — reimported ComfyUI reference-to-video workflow
- `minimax-r2v-audio-image.args.yaml` — aliases for the reference audio, reference image, input prompt, and duration
- `minimax-r2v-audio-image-prompt.md` — guide-compliant sample prompt using “The earth is round, but the sky is blue.”

The workflow uses `MiniMaxH3ReferenceToVideo`. The reference image is wired
to `ref_images.ref_image_0`, the reference audio is wired to
`ref_audios.ref_audio_0`, and the prompt refers to them as `<Picture 1>`
and `<Audio 1>`. Duration is supplied in seconds and converted by the workflow
to a valid MiniMax H3 frame count.

Local media paths are resolved from the current working directory and then
uploaded to ComfyUI when needed. Markdown prompt files are read as prompt
text.

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v-audio-image\minimax-r2v-audio-image.json --args-file .\workflows\minimax-r2v-audio-image\minimax-r2v-audio-image.args.yaml --set 'reference_audio=.\path\voice-reference.wav' --set 'reference_image=.\path\reference.png' --set 'prompt=.\workflows\minimax-r2v-audio-image\minimax-r2v-audio-image-prompt.md' --set 'duration=5' --output-folder .\out\minimax-r2v-audio-image-args
~~~

The bundled prompt follows the full-reference format in
`reference-docs/minimax-prompt-structure-guide.md`: it contains the six
ordered sections, stable reference labels, retention markers, shot-by-shot
description, and dialogue inside `<d>[English] ...</d>`.

The required MiniMax H3 model files are documented in the saved workflow's
notes. Prepare those models in ComfyUI before running the sample.
