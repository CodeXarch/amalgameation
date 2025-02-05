DEFAULT_GOAL := build

PHONY: fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	env GOOS=js GOARCH=wasm go build -o amalgameation.wasm thehumandroid.org/amalgameation 
