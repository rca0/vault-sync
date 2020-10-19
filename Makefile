build:
	CGO_ENABLED=0 GO111MODULES=on go build .
.PHONY: build

fmt:
	gofmt -w .
.PHONY: fmt

.DEFAULT_GOAL := build
