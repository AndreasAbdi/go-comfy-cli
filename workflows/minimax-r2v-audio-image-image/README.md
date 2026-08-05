# MiniMax R2V Audio + First/Last Image

Files:

- `minimax-r2v-audio-image-image.json` — saved ComfyUI reference-to-video workflow
- `minimax-r2v-audio-image-image.args.yaml` — aliases for the reference audio, first image, last image, prompt, and duration
- `minimax-r2v-audio-image-image-prompt.md` — guide-compliant two-keyframe prompt

The workflow uses `MiniMaxH3ReferenceToVideo` with the first image connected to
`ref_images.ref_image_0` and the end image connected to
`ref_images.ref_image_1`. In the prompt these are `<Picture 1>` and
`<Picture 2>`, respectively. The audio is connected to
`ref_audios.ref_audio_0` and is referenced as `<Audio 1>`.

Local media paths are resolved from the current working directory and then
ComfyUI's input directory. Markdown prompt files are read as prompt text.

## Args-file example

From the repository root, with the two images and WAV already in ComfyUI's
input directory:

~~~powershell
& .\go-comfy-cli.exe run --workflow .\workflows\minimax-r2v-audio-image-image\minimax-r2v-audio-image-image.json --args-file .\workflows\minimax-r2v-audio-image-image\minimax-r2v-audio-image-image.args.yaml --set 'reference_audio=voice-reference.wav' --set 'first_image=first-frame.png' --set 'last_image=last-frame.png' --set 'prompt=.\workflows\minimax-r2v-audio-image-image\minimax-r2v-audio-image-image-prompt.md' --set 'duration=5' --output-folder .\out\minimax-r2v-audio-image-image-args
~~~

If no separate end frame is needed, pass the same image to both `first_image`
and `last_image`. The workflow converts duration in seconds to a valid MiniMax
H3 frame count.

## JSON replacement example

~~~powershell
& .\go-comfy-cli.exe run --workflow .\workflows\minimax-r2v-audio-image-image\minimax-r2v-audio-image-image.json --replace-json '.nodes[] | select(.id == 148) | .widgets_values[0]::"voice-reference.wav"' --replace-json '.nodes[] | select(.id == 149) | .widgets_values[0]::"first-frame.png"' --replace-json '.nodes[] | select(.id == 150) | .widgets_values[0]::"last-frame.png"' --output-folder .\out\minimax-r2v-audio-image-image-replace
~~~

The required MiniMax H3 model files are documented in the saved workflow's
notes. Prepare those models in ComfyUI before running the sample.

