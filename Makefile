all: build

build:

test:
	gotest -v ./...

.PHONY: build test