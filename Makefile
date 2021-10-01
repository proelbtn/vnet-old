all: build

build: vnet

vnet: cmd/vnet/main.go
	go build -o vnet $^

.PHONY: build test vnet
