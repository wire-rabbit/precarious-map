.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
		go fmt ./...

vet: fmt
		go vet ./...

test: vet
		go test ./...

build: test
		go build
