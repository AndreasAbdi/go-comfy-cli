.PHONY: build test fmt

default: build test fmt

build: 
	go build ./...

test: 
	go test ./...

fmt: 
	go fmt ./...

