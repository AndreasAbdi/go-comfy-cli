# MiniMax Audio Reference to Video

Files:

- minimax-r2v-audio.json - workkflow file. 
- minimax-r2v-audio.args.yaml - aliases for audio, prompt, and duration
- minimax-r2v-audio-prompt.md - a test prompt. 

## Run the fixed workflow as saved

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v-audio\minimax-r2v-audio.json --output-folder .\out\minimax-r2v-audio-default
~~~

## Args-file example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v-audio\minimax-r2v-audio.json --args-file .\workflows\minimax-r2v-audio\minimax-r2v-audio.args.yaml --set 'audio=.\workflows\minimax-r2v-audio\hey_come_on_you_promised_to_take_me_shopping.wav' --set 'prompt=.\workflows\minimax-r2v-audio\minimax-r2v-audio-prompt.md' --set 'duration=5' --output-folder .\out\minimax-r2v-audio-args
~~~

## JSON replacement example

~~~powershell
go-comfy-cli run --workflow .\workflows\minimax-r2v-audio\minimax-r2v-audio.json --replace-json '.nodes[] | select(.id == 148) | .widgets_values[0]::"test.wav"' --replace-json '.nodes[] | select(.title == "Input Text (Prompt)") | .widgets_values[0]::"A chibi girl pouts and says the exact line from the audio reference"' --output-folder .\out\minimax-r2v-audio-replace
~~~
