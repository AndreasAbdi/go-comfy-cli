# what? 

This project is a simple golang wrapper that enables customers to invoke ComfyUI with a variety of workflows in a fairly straightforward way. 

Please see the ./README.md for detail on intended behavior and usage.

## technology choices: 
1. golang
2. urfave

## system structure
The package is roughly structured with all the internal logic under internal

1. internal/commands -> points to the CLI commands
2. internal/workflows -> utility functions around workflow args parsing, downloading, transformation, etc. 
3. internal/models -> utilities for downloading models for a given model file.
4. internal/media -> files for parsing and uploading media to the comfyui installation for parsing. 

## prepackaged workflows
under workflows, we have prepackaged workflows that customers can mess around with. 

1. anima -> used for image generation
2. qwen-image-edit -> used for image editing
3. minimax-reference-video -> used for generating based on a reference image
4. minimax-audio-reference -> used for generating based on some reference audio.

## instructions for prompt construction.

certain models like minimax-reference-video require that your prompts be constructed in a specific way. Please ensure that you MUST use the reference structure for generating the prompt.

the reference doc is in reference-docs/minimax-prompt-structure-guide.md for that specifically.

## ci

The Ci is github actions. 

The release is via goreleaser. 

The release releases against windows/linux/mac. 
Download directly via the golang page, or use an installer command. 
