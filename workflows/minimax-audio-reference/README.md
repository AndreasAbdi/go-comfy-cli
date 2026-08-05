# MiniMax Audio Reference to Video

Files:

- minimax-audio-reference.json — saved workflow with the standalone audio reference
  wired to ref_audios.ref_audio_0
- minimax-audio-reference.args.yaml — aliases for audio, prompt, and duration
- minimax-audio-reference-prompt.md — the chibi performance prompt
- shopping-promise.wav — bundled audio reference

The saved workflow is runnable without graph patch replacements. Its default
audio is the bundled recording, “Hey, you promised to take me shopping,” and
its embedded prompt describes a chibi girl pouting.
Local file references are resolved from the current working directory, then
ComfyUI's input directory. From the repository root, pass the bundled WAV's
explicit path when using it as a local upload.

## Run the fixed workflow as saved

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-audio-reference\minimax-audio-reference.json --args-file .\workflows\minimax-audio-reference\minimax-audio-reference.args.yaml --set 'audio=.\workflows\minimax-audio-reference\shopping-promise.wav' --output-folder .\out\minimax-audio-reference-default
~~~

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-audio-reference\minimax-audio-reference.json --args-file .\workflows\minimax-audio-reference\minimax-audio-reference.args.yaml --set 'audio=.\workflows\minimax-audio-reference\shopping-promise.wav' --set 'prompt=.\workflows\minimax-audio-reference\minimax-audio-reference-prompt.md' --set 'duration=5' --output-folder .\out\minimax-audio-reference-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-audio-reference\minimax-audio-reference.json --replace-json '.nodes[] | select(.id == 148) | .widgets_values[0]::".\\workflows\\minimax-audio-reference\\shopping-promise.wav"' --replace-json '.nodes[] | select(.title == "Input Text (Prompt)") | .widgets_values[0]::"A chibi girl pouts and says the exact line from the audio reference"' --output-folder .\out\minimax-audio-reference-replace
~~~
