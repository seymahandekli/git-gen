all: build install-tools

build:
	@echo Code generation
	@go generate ./...
	@echo Building binary
	@go build ./cmd/git-gen/

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
